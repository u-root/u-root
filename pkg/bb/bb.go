// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bb builds one busybox-like binary out of many Go command sources.
//
// This allows you to take two Go commands, such as Go implementations of `sl`
// and `cowsay` and compile them into one binary, callable like `./bb sl` and
// `./bb cowsay`.
//
// Which command is invoked is determined by `argv[0]` or `argv[1]` if
// `argv[0]` is not recognized.
//
// Under the hood, bb implements a Go source-to-source transformation on pure
// Go code. This AST transformation does the following:
//
//   - Takes a Go command's source files and rewrites them into Go package files
//     without global side effects.
//   - Writes a `main.go` file with a `main()` that calls into the appropriate Go
//     command package based on `argv[0]`.
//
// Principally, the AST transformation moves all global side-effects into
// callable package functions. E.g. `main` becomes `Main`, each `init` becomes
// `InitN`, and global variable assignments are moved into their own `InitN`.
package bb

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/golang"
)

// BuildBusybox builds a busybox of the given Go packages.
//
// pkgs is a list of Go import paths. If nil is returned, binaryPath will hold
// the busybox-style binary.
func BuildBusybox(env golang.Environ, pkgs []string, noStrip bool, binaryPath string) (nerr error) {
	tmpDir, err := ioutil.TempDir("", "bb-")
	if err != nil {
		return err
	}
	defer func() {
		if nerr != nil {
			log.Printf("Preserving bb temporary directory at %s due to error", tmpDir)
		} else {
			os.RemoveAll(tmpDir)
		}
	}()

	// INB4: yes, this *is* too clever. It's because Go modules are too
	// clever. Sorry.
	//
	// Inevitably, we will build commands across multiple modules, e.g.
	// u-root and u-bmc each have their own go.mod, but will get built into
	// one busybox executable.
	//
	// Each u-bmc and u-root command need their respective go.mod
	// dependencies, so we'll preserve their module file.
	//
	// However, we *also* need to still allow building with GOPATH and
	// vendoring. The structure we build *has* to also resemble a
	// GOPATH-based build.
	//
	// The easiest thing to do is to make the directory structure
	// GOPATH-compatible, and use go.mods to replace modules with the local
	// directories.
	//
	// So pretend GOPATH=tmp.
	//
	// Structure we'll build:
	//
	//   tmp/src/bb
	//   tmp/src/bb/main.go
	//      import "<module1>/cmd/foo"
	//      import "<module2>/cmd/bar"
	//
	//      func main()
	//
	// The only way to refer to other Go modules locally is to replace them
	// with local paths in a top-level go.mod:
	//
	//   tmp/go.mod
	//      package bb.u-root.com
	//
	//	replace <module1> => ./src/<module1>
	//	replace <module2> => ./src/<module2>
	//
	// Because GOPATH=tmp, the source should be in $GOPATH/src, just to
	// accommodate both build systems.
	//
	//   tmp/src/<module1>
	//   tmp/src/<module1>/go.mod
	//   tmp/src/<module1>/cmd/foo/main.go
	//
	//   tmp/src/<module2>
	//   tmp/src/<module2>/go.mod
	//   tmp/src/<module2>/cmd/bar/main.go

	bbDir := filepath.Join(tmpDir, "src/bb")
	if err := os.MkdirAll(bbDir, 0755); err != nil {
		return err
	}
	pkgDir := filepath.Join(tmpDir, "src")

	// Collect all packages that we need to actually re-write.
	var fpkgs []string
	seenPackages := make(map[string]struct{})
	for _, pkg := range pkgs {
		basePkg := path.Base(pkg)
		if _, ok := seenPackages[basePkg]; ok {
			return fmt.Errorf("failed to build with bb: found duplicate pkgs %s", basePkg)
		}
		seenPackages[basePkg] = struct{}{}

		fpkgs = append(fpkgs, pkg)
	}

	// Ask go about all the packages in one batch for dependency caching.
	ps, err := NewPackages(env, fpkgs...)
	if err != nil {
		return fmt.Errorf("finding packages failed: %v", err)
	}

	var bbImports []string
	for _, p := range ps {
		destination := filepath.Join(pkgDir, p.Pkg.PkgPath)
		if err := p.Rewrite(destination, "github.com/u-root/u-root/pkg/bb/bbmain"); err != nil {
			return fmt.Errorf("rewriting %q failed: %v", p.Pkg.PkgPath, err)
		}
		bbImports = append(bbImports, p.Pkg.PkgPath)
	}

	bb, err := NewPackages(env, "github.com/u-root/u-root/pkg/bb/bbmain/cmd")
	if err != nil {
		return err
	}
	if len(bb) == 0 {
		return fmt.Errorf("bb cmd template missing")
	}

	// Add bb to the list of packages that need their dependencies.
	mainPkgs := append(ps, bb[0])

	// Module-enabled Go programs resolve their dependencies in one of two ways:
	//
	// - locally, if the dependency is *in* the module
	// - remotely, if it is outside of the module
	//
	// I.e. if the module is github.com/u-root/u-root,
	//
	// - local: github.com/u-root/u-root/pkg/uio
	// - remote: github.com/hugelgupf/p9/p9
	//
	// For remote dependencies, we copy the go.mod into the temporary directory.
	// For local dependencies, we copy all dependency packages' files over.
	var depPkgs, modulePaths []string
	for _, p := range mainPkgs {
		// Find all dependency packages that are *within* module boundaries for this package.
		//
		// writeDeps also copies the go.mod into the right place.
		mods, modulePath, err := writeDeps(env, pkgDir, p.Pkg)
		if err != nil {
			return fmt.Errorf("resolving dependencies for %q failed: %v", p.Pkg.PkgPath, err)
		}
		depPkgs = append(depPkgs, mods...)
		if len(modulePath) > 0 {
			modulePaths = append(modulePaths, modulePath)
		}
	}

	// Create bb main.go.
	if err := CreateBBMainSource(bb[0].Pkg, bbImports, bbDir); err != nil {
		return fmt.Errorf("creating bb main() file failed: %v", err)
	}

	// Copy local dependency packages into temporary module directories.
	deps, err := NewPackages(env, depPkgs...)
	if err != nil {
		return err
	}
	for _, p := range deps {
		if err := p.Write(filepath.Join(pkgDir, p.Pkg.PkgPath)); err != nil {
			return err
		}
	}

	// Add local replace rules for all modules we're compiling.
	//
	// This is the only way to locally reference another modules'
	// repository. Otherwise, go'll try to go online to get the source.
	//
	// The module name is something that'll never be online, lest Go
	// decides to go on the internet.
	if len(modulePaths) == 0 {
		env.GOPATH = tmpDir
		// Compile bb.
		return env.Build("bb", binaryPath, golang.BuildOpts{NoStrip: noStrip})
	}

	content := `module bb.u-root.com`
	for _, mpath := range modulePaths {
		content += fmt.Sprintf("\nreplace %s => ./src/%s\n", mpath, mpath)
	}
	if err := ioutil.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(content), 0755); err != nil {
		return err
	}

	// Compile bb.
	return env.BuildDir(bbDir, binaryPath, golang.BuildOpts{NoStrip: noStrip})
}

func writeDeps(env golang.Environ, pkgDir string, p *packages.Package) ([]string, string, error) {
	listp, err := env.FindOne(p.PkgPath)
	if err != nil {
		return nil, "", err
	}

	var deps []string
	if listp.Module != nil {
		if err := os.MkdirAll(filepath.Join(pkgDir, listp.Module.Path), 0755); err != nil {
			return nil, "", err
		}

		// Use the module file for all outside dependencies.
		if err := cp.Copy(listp.Module.GoMod, filepath.Join(pkgDir, listp.Module.Path, "go.mod")); err != nil {
			return nil, "", err
		}

		// Collect all "local" dependency packages, to be copied into
		// the temporary directory structure later.
		for _, dep := range listp.Deps {
			// Is this a dependency within the module?
			if strings.HasPrefix(dep, listp.Module.Path) {
				deps = append(deps, dep)
			}
		}
		return deps, listp.Module.Path, nil
	}

	// If modules are not enabled, we need a copy of *ALL*
	// non-standard-library dependencies in the temporary directory.
	for _, dep := range listp.Deps {
		// First component of package path contains a "."?
		//
		// Poor man's standard library test.
		firstComp := strings.SplitN(dep, "/", 2)
		if strings.Contains(firstComp[0], ".") {
			deps = append(deps, dep)
		}
	}
	return deps, "", nil
}

// CreateBBMainSource creates a bb Go command that imports all given pkgs.
//
// p must be the bb template.
//
// - For each pkg in pkgs, add
//     import _ "pkg"
//   to astp's first file.
// - Write source file out to destDir.
func CreateBBMainSource(p *packages.Package, pkgs []string, destDir string) error {
	if len(p.Syntax) != 1 {
		return fmt.Errorf("bb cmd template is supposed to only have one file")
	}
	for _, pkg := range pkgs {
		// Add side-effect import to bb binary so init registers itself.
		//
		// import _ "pkg"
		astutil.AddNamedImport(p.Fset, p.Syntax[0], "_", pkg)
	}

	return writeFiles(destDir, p.Fset, p.Syntax)
}

// Package is a Go package.
type Package struct {
	// Name is the executable command name.
	//
	// In the standard Go tool chain, this is usually the base name of the
	// directory containing its source files.
	Name string

	// Pkg is the actual data about the package.
	Pkg *packages.Package

	// initCount keeps track of what the next init's index should be.
	initCount uint

	// init is the cmd.Init function that calls all other InitXs in the
	// right order.
	init *ast.FuncDecl

	// initAssigns is a map of assignment expression -> InitN function call
	// statement.
	//
	// That InitN should contain the assignment statement for the
	// appropriate assignment expression.
	//
	// types.Info.InitOrder keeps track of Initializations by Lhs name and
	// Rhs ast.Expr.  We reparent the Rhs in assignment statements in InitN
	// functions, so we use the Rhs as an easy key here.
	// types.Info.InitOrder + initAssigns can then easily be used to derive
	// the order of Stmts in the "real" init.
	//
	// The key Expr must also be the AssignStmt.Rhs[0].
	initAssigns map[ast.Expr]ast.Stmt
}

// NewPackages collects package metadata about all named packages.
func NewPackages(env golang.Environ, names ...string) ([]*Package, error) {
	ps, err := loadPkgs(env, names...)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %v: %v", names, err)
	}
	var ips []*Package
	for _, p := range ps {
		ips = append(ips, NewPackage(path.Base(p.PkgPath), p))
	}
	return ips, nil
}

func loadPkgs(env golang.Environ, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedFiles | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedCompiledGoFiles,
		Env:  append(os.Environ(), env.Env()...),
	}
	return packages.Load(cfg, patterns...)
}

// NewPackage creates a new Package based on an existing packages.Package.
func NewPackage(name string, p *packages.Package) *Package {
	pp := &Package{
		// Name is the executable name.
		Name:        path.Base(name),
		Pkg:         p,
		initAssigns: make(map[ast.Expr]ast.Stmt),
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
	return pp
}

func (p *Package) nextInit(addToCallList bool) *ast.Ident {
	i := ast.NewIdent(fmt.Sprintf("Init%d", p.initCount))
	if addToCallList {
		p.init.Body.List = append(p.init.Body.List, &ast.ExprStmt{X: &ast.CallExpr{Fun: i}})
	}
	p.initCount++
	return i
}

// TODO:
// - write an init name generator, in case InitN is already taken.
func (p *Package) rewriteFile(f *ast.File) bool {
	hasMain := false

	// Change the package name declaration from main to the command's name.
	f.Name.Name = p.Name

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
		if pkg == p.Pkg.Types {
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

				// For each assignment, create a new init
				// function, and place it in the same file.
				for i, name := range s.Names {
					varInit := &ast.FuncDecl{
						Name: p.nextInit(false),
						Type: &ast.FuncType{
							Params:  &ast.FieldList{},
							Results: nil,
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Lhs: []ast.Expr{name},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{s.Values[i]},
								},
							},
						},
					}
					// Add a call to the new init func to
					// this map, so they can be added to
					// Init0() in the correct init order
					// later.
					p.initAssigns[s.Values[i]] = &ast.ExprStmt{X: &ast.CallExpr{Fun: varInit.Name}}
					f.Decls = append(f.Decls, varInit)
				}

				// Add the type of the expression to the global
				// declaration instead.
				if s.Type == nil {
					typ := p.Pkg.TypesInfo.Types[s.Values[0]]
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
				d.Name = p.nextInit(true)
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

// Write writes p into destDir.
func (p *Package) Write(destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, fp := range p.Pkg.OtherFiles {
		if err := cp.Copy(fp, filepath.Join(destDir, filepath.Base(fp))); err != nil {
			return err
		}
	}

	return writeFiles(destDir, p.Pkg.Fset, p.Pkg.Syntax)
}

func writeFiles(destDir string, fset *token.FileSet, files []*ast.File) error {
	// Write all files out.
	for _, file := range files {
		name := fset.File(file.Package).Name()

		path := filepath.Join(destDir, filepath.Base(name))
		if err := writeFile(path, fset, file); err != nil {
			return err
		}
	}
	return nil
}

// Rewrite rewrites p into destDir as a bb package using bbImportPath for the
// bb implementation.
func (p *Package) Rewrite(destDir, bbImportPath string) error {
	// This init holds all variable initializations.
	//
	// func Init0() {}
	varInit := &ast.FuncDecl{
		Name: p.nextInit(true),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	var mainFile *ast.File
	for _, sourceFile := range p.Pkg.Syntax {
		if hasMainFile := p.rewriteFile(sourceFile); hasMainFile {
			mainFile = sourceFile
		}
	}
	if mainFile == nil {
		return fmt.Errorf("no main function found in package %q", p.Pkg.PkgPath)
	}

	// Add variable initializations to Init0 in the right order.
	for _, initStmt := range p.Pkg.TypesInfo.InitOrder {
		a, ok := p.initAssigns[initStmt.Rhs]
		if !ok {
			return fmt.Errorf("couldn't find init assignment %s", initStmt)
		}
		varInit.Body.List = append(varInit.Body.List, a)
	}

	// import bb "bbImportPath"
	astutil.AddNamedImport(p.Pkg.Fset, mainFile, "bb", bbImportPath)

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
							Value: strconv.Quote(p.Name),
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

	return p.Write(destDir)
}

func writeFile(path string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("error formatting Go file %q: %v", path, err)
	}
	return writeGoFile(path, buf.Bytes())
}

func writeGoFile(path string, code []byte) error {
	// Format the file. Do not fix up imports, as we only moved code around
	// within files.
	opts := imports.Options{
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: true,
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
