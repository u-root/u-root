// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
)

// Commands to skip building in bb mode.
var skip = map[string]struct{}{
	"bb": {},
}

// BuildBusybox builds a busybox of the given Go packages.
//
// pkgs is a list of Go import paths. If nil is returned, binaryPath will hold
// the busybox-style binary.
func BuildBusybox(env golang.Environ, pkgs []string, binaryPath string) error {
	urootPkg, err := env.Package("github.com/u-root/u-root")
	if err != nil {
		return err
	}

	bbDir := filepath.Join(urootPkg.Dir, "bb")
	// Blow bb away before trying to re-create it.
	if err := os.RemoveAll(bbDir); err != nil {
		return err
	}
	if err := os.MkdirAll(bbDir, 0755); err != nil {
		return err
	}

	var bbPackages []string
	// Move and rewrite package files.
	importer := importer.For("source", nil)
	for _, pkg := range pkgs {
		if _, ok := skip[path.Base(pkg)]; ok {
			continue
		}

		// TODO: use bbDir to derive import path below or vice versa.
		if err := RewritePackage(env, pkg, filepath.Join(bbDir, "cmds"), "github.com/u-root/u-root/pkg/bb", importer); err != nil {
			return err
		}

		bbPackages = append(bbPackages, fmt.Sprintf("github.com/u-root/u-root/bb/cmds/%s", path.Base(pkg)))
	}

	// Create bb main.go.
	if err := CreateBBMainSource(env, importer, "github.com/u-root/u-root/pkg/bb/cmd", bbPackages, bbDir); err != nil {
		return err
	}

	// Compile bb.
	return env.Build("github.com/u-root/u-root/bb", binaryPath, golang.BuildOpts{})
}

// CreateBBMainSource creates a bb Go command that imports all given pkgs.
//
// - Takes code from templatePkg, which must be ONE Go main() file.
// - For each pkg in pkgs, add
//     import _ "pkg"
//   to templatePkg main() file
// - Write source file out to destDir.
func CreateBBMainSource(env golang.Environ, importer types.Importer, templatePkg string, pkgs []string, destDir string) error {
	bb, err := getPackage(env, templatePkg, importer)
	if err != nil {
		return err
	}
	if bb == nil {
		return fmt.Errorf("bb cmd template missing")
	}
	if len(bb.ast.Files) != 1 {
		return fmt.Errorf("bb cmd template is supposed to only have one file")
	}
	for _, pkg := range pkgs {
		for _, sourceFile := range bb.ast.Files {
			// Add side-effect import to bb binary so init registers itself.
			//
			// import _ "pkg"
			astutil.AddNamedImport(bb.fset, sourceFile, "_", pkg)
			break
		}
	}

	// Write bb main binary out.
	for filePath, sourceFile := range bb.ast.Files {
		path := filepath.Join(destDir, filepath.Base(filePath))
		if err := writeFile(path, bb.fset, sourceFile); err != nil {
			return err
		}
		break
	}
	return nil
}

// BBBuild is an implementation of Build for the busybox-like u-root initramfs.
//
// BBBuild rewrites the source files of the packages given to create one
// busybox-like binary containing all commands in `opts.Packages`.
func BBBuild(af ArchiveFiles, opts BuildOpts) error {
	// Build the busybox binary.
	bbPath := filepath.Join(opts.TempDir, "bb")
	if err := BuildBusybox(opts.Env, opts.Packages, bbPath); err != nil {
		return err
	}

	binDir := opts.TargetDir("bbin")

	// Build initramfs.
	if err := af.AddFile(bbPath, path.Join(binDir, "bb")); err != nil {
		return err
	}

	// Add symlinks for included commands to initramfs.
	for _, pkg := range opts.Packages {
		if _, ok := skip[path.Base(pkg)]; ok {
			continue
		}

		// Add a symlink /bbin/{cmd} -> /bbin/bb to our initramfs.
		if err := af.AddRecord(cpio.Symlink(filepath.Join(binDir, path.Base(pkg)), "bb")); err != nil {
			return err
		}
	}

	// Symlink from /init to busybox init.
	return af.AddRecord(cpio.Symlink("init", path.Join(binDir, "init")))
}

// Package is a Go package.
//
// It holds AST, type, file, and Go package information about a Go package.
type Package struct {
	name string

	pkg         *build.Package
	fset        *token.FileSet
	ast         *ast.Package
	typeInfo    types.Info
	types       *types.Package
	sortedFiles []*ast.File

	// initCount keeps track of what the next init's index should be.
	initCount uint

	// init is the cmd.Init function that calls all other InitXs in the
	// right order.
	init *ast.FuncDecl

	// initAssigns is a map of assignment expression -> assignment
	// statement.
	//
	// types.Info.InitOrder keeps track of Initializations by Lhs name and
	// Rhs ast.Expr.  We reparent the Rhs in assignment statements in InitX
	// functions, so we use the Rhs as an easy key here.
	// types.Info.InitOrder + initAssigns can then easily be used to derive
	// the order of AssignStmts.
	//
	// The key Expr must also be the AssignStmt.Rhs[0].
	initAssigns map[ast.Expr]*ast.AssignStmt
}

func (p *Package) nextInit() *ast.Ident {
	i := ast.NewIdent(fmt.Sprintf("Init%d", p.initCount))
	p.init.Body.List = append(p.init.Body.List, &ast.ExprStmt{X: &ast.CallExpr{Fun: i}})
	p.initCount++
	return i
}

// TODO:
// - write an init name generator, in case InitN is already taken.
// - also rewrite all non-Go-stdlib dependencies.
func (p *Package) rewriteFile(f *ast.File) bool {
	hasMain := false

	// Change the package name declaration from main to the command's name.
	f.Name = ast.NewIdent(p.name)

	// Map of fully qualified package name -> imported alias in the file.
	importAliases := make(map[string]string)
	for _, impt := range f.Imports {
		if impt.Name != nil {
			importPath, err := strconv.Unquote(impt.Path.Value)
			if err != nil {
				panic(err)
			}
			importAliases[importPath] = impt.Name.Name
		}
	}

	// When the types.TypeString function translates package names, it uses
	// this function to map fully qualified package paths to a local alias,
	// if it exists.
	qualifier := func(pkg *types.Package) string {
		name, ok := importAliases[pkg.Path()]
		if ok {
			return name
		}
		// When referring to self, don't use any package name.
		if pkg == p.types {
			return ""
		}
		return pkg.Name()
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// We only care about vars.
			if d.Tok != token.VAR {
				break
			}
			for _, spec := range d.Specs {
				s := spec.(*ast.ValueSpec)
				if s.Values == nil {
					continue
				}

				// Add an assign statement for these values to init.
				for i, name := range s.Names {
					p.initAssigns[s.Values[i]] = &ast.AssignStmt{
						Lhs: []ast.Expr{name},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{s.Values[i]},
					}
				}

				// Add the type of the expression to the global
				// declaration instead.
				if s.Type == nil {
					typ := p.typeInfo.Types[s.Values[0]]
					s.Type = ast.NewIdent(types.TypeString(typ.Type, qualifier))
				}
				s.Values = nil
			}

		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.Name == "main" {
				d.Name.Name = "Main"
				hasMain = true
			}
			if d.Recv == nil && d.Name.Name == "init" {
				d.Name = p.nextInit()
			}
		}
	}

	// Now we change any import names attached to package declarations. We
	// just upcase it for now; it makes it easy to look in bbsh for things
	// we changed, e.g. grep -r bbsh Import is useful.
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// import") {
				c.Text = "// Import" + c.Text[9:]
			}
		}
	}
	return hasMain
}

func writeFile(path string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("error formatting Go file %q: %v", path, err)
	}
	return writeGoFile(path, buf.Bytes())
}

func writeGoFile(path string, code []byte) error {
	// Fix up imports.
	opts := imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	code, err := imports.Process("commandline", code, &opts)
	if err != nil {
		return fmt.Errorf("bad parse while processing imports %q: %v", path, err)
	}

	if err := ioutil.WriteFile(path, code, 0644); err != nil {
		return fmt.Errorf("error writing Go file to %q: %v", path, err)
	}
	return nil
}

func getPackage(env golang.Environ, importPath string, importer types.Importer) (*Package, error) {
	p, err := env.Package(importPath)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(p.Dir)
	if !p.IsCommand() {
		return nil, fmt.Errorf("package %q is not a command and cannot be included in bb", name)
	}

	fset := token.NewFileSet()
	pars, err := parser.ParseDir(fset, p.Dir, func(fi os.FileInfo) bool {
		// Only parse Go files that match build tags of this package.
		for _, name := range p.GoFiles {
			if name == fi.Name() {
				return true
			}
		}
		return false
	}, parser.ParseComments)
	if err != nil {
		log.Printf("can't parsedir %q: %v", p.Dir, err)
		return nil, err
	}

	pp := &Package{
		name: name,
		pkg:  p,
		fset: fset,
		ast:  pars[p.Name],
		typeInfo: types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
		initAssigns: make(map[ast.Expr]*ast.AssignStmt),
	}

	// This Init will hold calls to all other InitXs.
	pp.init = &ast.FuncDecl{
		Name: ast.NewIdent("Init"),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	// The order of types.Info.InitOrder depends on this list of files
	// always being passed to conf.Check in the same order.
	filenames := make([]string, 0, len(pp.ast.Files))
	for name := range pp.ast.Files {
		filenames = append(filenames, name)
	}
	sort.Strings(filenames)

	pp.sortedFiles = make([]*ast.File, 0, len(pp.ast.Files))
	for _, name := range filenames {
		pp.sortedFiles = append(pp.sortedFiles, pp.ast.Files[name])
	}
	// Type-check the package before we continue. We need types to rewrite
	// some statements.
	conf := types.Config{
		Importer: importer,

		// We only need global declarations' types.
		IgnoreFuncBodies: true,
	}
	tpkg, err := conf.Check(pp.pkg.ImportPath, pp.fset, pp.sortedFiles, &pp.typeInfo)
	if err != nil {
		return nil, fmt.Errorf("type checking failed: %v", err)
	}
	pp.types = tpkg
	return pp, nil
}

// RewritePackage rewrites pkgPath to be bb-mode compatible, where destDir is
// the file system destination of the written files and bbImportPath is the Go
// import path of the bb package to register with.
func RewritePackage(env golang.Environ, pkgPath, destDir, bbImportPath string, importer types.Importer) error {
	p, err := getPackage(env, pkgPath, importer)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	pkgDir := filepath.Join(destDir, p.name)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return err
	}

	// This init holds all variable initializations.
	//
	// func Init0() {}
	varInit := &ast.FuncDecl{
		Name: p.nextInit(),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	var mainFile *ast.File
	for _, sourceFile := range p.sortedFiles {
		if hasMainFile := p.rewriteFile(sourceFile); hasMainFile {
			mainFile = sourceFile
		}
	}
	if mainFile == nil {
		return os.RemoveAll(pkgDir)
	}

	// Add variable initializations to Init0 in the right order.
	for _, initStmt := range p.typeInfo.InitOrder {
		a, ok := p.initAssigns[initStmt.Rhs]
		if !ok {
			return fmt.Errorf("couldn't find init assignment %s", initStmt)
		}
		varInit.Body.List = append(varInit.Body.List, a)
	}

	// import bb "bbImportPath"
	astutil.AddNamedImport(p.fset, mainFile, "bb", bbImportPath)

	// func init() {
	//   bb.Register(p.name, Init, Main)
	// }
	bbRegisterInit := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.CallExpr{
					Fun: ast.NewIdent("bb.Register"),
					Args: []ast.Expr{
						// name=
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote(p.name),
						},
						// init=
						ast.NewIdent("Init"),
						// main=
						ast.NewIdent("Main"),
					},
				}},
			},
		},
	}

	// We could add these statements to any of the package files. We choose
	// the one that contains Main() to guarantee reproducibility of the
	// same bbsh binary.
	mainFile.Decls = append(mainFile.Decls, varInit, p.init, bbRegisterInit)

	// Write all files out.
	for filePath, sourceFile := range p.ast.Files {
		path := filepath.Join(pkgDir, filepath.Base(filePath))
		if err := writeFile(path, p.fset, sourceFile); err != nil {
			return err
		}
	}
	return nil
}
