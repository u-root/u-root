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
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
)

// Per-package templates.
//
// TODO: Generate the cmd_xx.go template using the AST, and use a function-name
// generator to avoid collisions.
const (
	cmdFunc = `package main

import "github.com/u-root/u-root/bbsh/cmds/{{.CmdName}}"

func _forkbuiltin_{{.CmdName}}(c *Command) error {
	{{.CmdName}}.Main()
	return nil
}

func {{.CmdName}}Init() {
	addForkBuiltIn("{{.CmdName}}", _forkbuiltin_{{.CmdName}})
	{{.CmdName}}.Init()
}
`
)

// init.go
const initGo = `package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
        "syscall"

	"github.com/u-root/u-root/pkg/uroot/util"
)

func init() {
	// Init for the command that is being run.
	if f, ok := initMap[filepath.Base(os.Args[0])]; ok {
		f()
	}

	// Yeah, it's weird that this stuff happens in bbsh's init.
	// We only do it if we actually *are* init. But merging this into bbsh
	// should happen some other way. This should just be in bbsh.Main?
	//
	// Maybe init should just be a normal bb command, just like the others,
	// except with a special symlink to /init?

	if os.Args[0] != "/init" {
		// Skipping root file system setup since we are not /init.
		return
	}

	if os.Getpid() != 1 {
		// Skipping root file system setup since /init is not pid 1.
		return
	}

	// Set up root file system. Mount stuff.
	util.Rootfs()

	// Figure out which init to run.
	//
	// Spawn the first shell. We had been running the shell as pid 1 but
	// that makes control tty stuff messy. We think.
	for _, v := range []string{"/inito", "/bbin/uinit", "/bbin/rush"} {
		cmd := exec.Command(v)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
		log.Printf("Run %v", cmd)
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
	}

	// This will drop us into a rush prompt, since this is the init for rush.
	// That's a nice fallback for when everything goes wrong. 
	return
}
`

// Commands to skip building in bb mode. init and rush should be obvious
// builtin and script we skip as we have no toolchain in this mode.
var skip = map[string]struct{}{
	"builtin": {},
	"init":    {},
	"rush":    {},
	"script":  {},
}

type bbBuilder struct {
	opts    BuildOpts
	bbshDir string
	initMap string
	af      ArchiveFiles

	importer types.Importer
}

// BBBuild is an implementation of Build for the busybox-like u-root initramfs.
//
// BBBuild rewrites the source files of the packages given to create one
// busybox-like binary containing all commands in `opts.Packages`.
func BBBuild(opts BuildOpts) (ArchiveFiles, error) {
	urootPkg, err := opts.Env.Package("github.com/u-root/u-root")
	if err != nil {
		return ArchiveFiles{}, err
	}

	bbshDir := filepath.Join(urootPkg.Dir, "bbsh")
	// Blow bbsh away before trying to re-create it.
	if err := os.RemoveAll(bbshDir); err != nil {
		return ArchiveFiles{}, err
	}
	if err := os.MkdirAll(bbshDir, 0755); err != nil {
		return ArchiveFiles{}, err
	}

	builder := &bbBuilder{
		opts:     opts,
		bbshDir:  bbshDir,
		initMap:  "package main\n\nvar initMap = map[string]func() {",
		af:       NewArchiveFiles(),
		importer: importer.For("source", nil),
	}

	// Move and rewrite package files.
	for _, pkg := range opts.Packages {
		if _, ok := skip[filepath.Base(pkg)]; ok {
			continue
		}

		if err := builder.moveCommand(pkg); err != nil {
			return ArchiveFiles{}, err
		}
	}

	// Create init.go.
	if err := ioutil.WriteFile(filepath.Join(builder.bbshDir, "init.go"), []byte(initGo), 0644); err != nil {
		return ArchiveFiles{}, err
	}

	// Move rush shell over.
	p, err := opts.Env.Package("github.com/u-root/u-root/cmds/rush")
	if err != nil {
		return ArchiveFiles{}, err
	}

	if err := filepath.Walk(p.Dir, func(name string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filepath.Join(builder.bbshDir, fi.Name()), b, 0644)
	}); err != nil {
		return ArchiveFiles{}, err
	}

	// Map init functions.
	builder.initMap += "\n}"
	if err := ioutil.WriteFile(filepath.Join(builder.bbshDir, "initmap.go"), []byte(builder.initMap), 0644); err != nil {
		return ArchiveFiles{}, err
	}

	// Compile rush + commands to /bbin/rush.
	rushPath := filepath.Join(opts.TempDir, "rush")
	if err := opts.Env.Build("github.com/u-root/u-root/bbsh", rushPath, golang.BuildOpts{}); err != nil {
		return ArchiveFiles{}, err
	}
	if err := builder.af.AddFile(rushPath, "bbin/rush"); err != nil {
		return ArchiveFiles{}, err
	}

	// Symlink from /init to rush.
	if err := builder.af.AddRecord(cpio.Symlink("init", "/bbin/rush")); err != nil {
		return ArchiveFiles{}, err
	}
	return builder.af, nil
}

type CommandTemplate struct {
	CmdName string
}

type Package struct {
	name string

	pkg      *build.Package
	fset     *token.FileSet
	ast      *ast.Package
	typeInfo types.Info
	types    *types.Package

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

func (p *Package) CommandTemplate() CommandTemplate {
	return CommandTemplate{
		CmdName: p.name,
	}
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
			// TODO: eval this or find a post-AST source for this
			// info.
			importAliases[strings.Trim(impt.Path.Value, "\"")] = impt.Name.Name
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

func (p *Package) writeTemplate(path string, text string) error {
	var b bytes.Buffer
	t := template.Must(template.New("uroot").Parse(text))
	if err := t.Execute(&b, p.CommandTemplate()); err != nil {
		return fmt.Errorf("couldn't write template %q: %v", path, err)
	}

	return writeGoFile(path, b.Bytes())
}

func getPackage(opts BuildOpts, importPath string, importer types.Importer) (*Package, error) {
	p, err := opts.Env.Package(importPath)
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
		return nil, nil
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

	var files []*ast.File
	for _, f := range pp.ast.Files {
		files = append(files, f)
	}
	// Type-check the package before we continue. We need types to rewrite
	// some statements.
	conf := types.Config{
		Importer: importer,

		// We only need global declarations' types.
		IgnoreFuncBodies: true,
	}
	tpkg, err := conf.Check(pp.pkg.ImportPath, pp.fset, files, &pp.typeInfo)
	if err != nil {
		return nil, fmt.Errorf("type checking failed: %v", err)
	}
	pp.types = tpkg
	return pp, nil
}

func (b *bbBuilder) moveCommand(pkgPath string) error {
	p, err := getPackage(b.opts, pkgPath, b.importer)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	pkgDir := filepath.Join(b.bbshDir, "cmds", p.name)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return err
	}

	// This init holds all variable initializations.
	varInit := &ast.FuncDecl{
		Name: p.nextInit(),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	var mainPath string
	var hasMain bool
	for filePath, sourceFile := range p.ast.Files {
		if hasMainFile := p.rewriteFile(sourceFile); hasMainFile {
			hasMain = true
			mainPath = filePath
		}
	}
	if !hasMain {
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

	// We could add these statements to any of the package files. We choose
	// the one that contains Main() to guarantee reproducibility of the
	// same bbsh binary.
	p.ast.Files[mainPath].Decls = append(p.ast.Files[mainPath].Decls, varInit, p.init)

	for filePath, sourceFile := range p.ast.Files {
		path := filepath.Join(pkgDir, filepath.Base(filePath))
		if err := writeFile(path, p.fset, sourceFile); err != nil {
			return err
		}
	}

	cmdPath := filepath.Join(b.bbshDir, fmt.Sprintf("cmd_%s.go", p.name))
	if err := p.writeTemplate(cmdPath, cmdFunc); err != nil {
		return err
	}

	b.initMap += "\n\t\"" + p.name + "\": " + p.name + "Init,"

	// Add a symlink to our bbsh.
	return b.af.AddRecord(cpio.Symlink(filepath.Join("bbin", p.name), "/bbin/rush"))
}
