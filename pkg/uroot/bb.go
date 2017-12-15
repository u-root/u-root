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
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
)

// Per-package templates.
const (
	cmdFunc = `package main

import "github.com/u-root/u-root/bbsh/cmds/{{.CmdName}}"

func _forkbuiltin_{{.CmdName}}(c *Command) (err error) {
	{{.CmdName}}.Main()
	return
}

func {{.CmdName}}Init() {
	addForkBuiltIn("{{.CmdName}}", _forkbuiltin_{{.CmdName}})
	{{.Init}}
}
`

	bbsetupGo = `package {{.CmdName}}

import "flag"

var {{.CmdName}}flag = flag.NewFlagSet("{{.CmdName}}", flag.ExitOnError)
`
)

// init.go
const initGo = `package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
        "syscall"

	"github.com/u-root/u-root/pkg/uroot/util"
)

func usage () {
	n := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s:\n", n)
	flag.VisitAll(func(f *flag.Flag) {
		if !strings.HasPrefix(f.Name, n+".") {
			return
		}
		fmt.Fprintf(os.Stderr, "\tFlag %s: '%s', Default %v, Value %v\n", f.Name[len(n)+1:], f.Usage, f.Value, f.DefValue)
	})
}

func init() {
	flag.Usage = usage
	// This getpid adds a bit of cost to each invocation (not much really)
	// but it allows us to merge init and sh. The 600K we save is worth it.
	// Figure out which init to run. We must always do this.

	// log.Printf("init: os is %v, initMap %v", filepath.Base(os.Args[0]), initMap)
	// we use filepath.Base in case they type something like ./cmd
	if f, ok := initMap[filepath.Base(os.Args[0])]; ok {
		//log.Printf("run the Init function for %v: run %v", os.Args[0], f)
		f()
	}

	if os.Args[0] != "/init" {
		//log.Printf("Skipping root file system setup since we are not /init")
		return
	}
	if os.Getpid() != 1 {
		//log.Printf("Skipping root file system setup since /init is not pid 1")
		return
	}
	util.Rootfs()

        // spawn the first shell. We had been running the shell as pid 1
        // but that makes control tty stuff messy. We think.
        cloneFlags := uintptr(0)
	for _, v := range []string{"/inito", "/bbin/uinit", "/bbin/rush"} {
		cmd := exec.Command(v)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		var fd int
		cmd.SysProcAttr = &syscall.SysProcAttr{Ctty: fd, Setctty: true, Setsid: true, Cloneflags: cloneFlags}
		log.Printf("Run %v", cmd)
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
		// only the first init needs its own PID space.
		cloneFlags = 0
	}

	// This will drop us into a rush prompt, since this is the init for rush.
	// That's a nice fallback for when everything goes wrong. 
	return
}
`

// Commands to skip building in bb mode. init and rush should be obvious
// builtin and script we skip as we have no toolchain in this mode.
var skip = map[string]struct{}{
	"builtin": struct{}{},
	"init":    struct{}{},
	"rush":    struct{}{},
	"script":  struct{}{},
}

type bbBuilder struct {
	opts    BuildOpts
	bbshDir string
	initMap string
	af      ArchiveFiles
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
		opts:    opts,
		bbshDir: bbshDir,
		initMap: "package main\n\nvar initMap = map[string]func() {",
		af:      NewArchiveFiles(),
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
	Gopath  string
	CmdName string
	Init    string
}

type Package struct {
	name string
	pkg  *build.Package
	fset *token.FileSet
	ast  *ast.Package

	initCount uint
	init      string
}

func (p *Package) CommandTemplate() CommandTemplate {
	return CommandTemplate{
		Gopath:  p.pkg.Root,
		CmdName: p.name,
		Init:    p.init,
	}
}

func (p *Package) rewriteFile(opts BuildOpts, f *ast.File) bool {
	// Inspect the AST and change all instances of main()
	var pos token.Position
	hasMain := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		// This is rather gross. Arguably, so is the way that Go has
		// embedded build information in comments ... this import
		// comment attachment to a package came in 1.4, a few years
		// ago, and it only just bit us with one file in upspin. So we
		// go with gross.
		case *ast.Ident:
			// We assume the first identifier is the package id.
			if !pos.IsValid() {
				pos = p.fset.Position(x.Pos())
			}

		case *ast.File:
			x.Name.Name = p.name

		case *ast.FuncDecl:
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("Main")
				// Append a return.
				x.Body.List = append(x.Body.List, &ast.ReturnStmt{})
				hasMain = true
			}
			if x.Recv == nil && x.Name.Name == "init" {
				x.Name.Name = fmt.Sprintf("Init%d", p.initCount)
				p.init += fmt.Sprintf("%s.Init%d()\n", p.name, p.initCount)
				p.initCount++
			}

		// Rewrite use of the flag package.
		//
		// The flag package uses a global variable to contain all
		// flags. This works poorly with the busybox mode, as flags may
		// conflict, so as part of turning commands into packages, we
		// rewrite their use of flags to use a package-private FlagSet.
		//
		// bbsetup.go contains a var for the package flagset with
		// params (packagename, os.ExitOnError).
		//
		// We rewrite arguments for calls to flag.Parse from () to
		// (os.Args[1:]). We rewrite all other uses of 'flag.' to
		// '"commandname"+flag.'.
		case *ast.CallExpr:
			switch s := x.Fun.(type) {
			case *ast.SelectorExpr:
				switch i := s.X.(type) {
				case *ast.Ident:
					if i.Name != "flag" {
						break
					}
					switch s.Sel.Name {
					case "Parse":
						i.Name = p.name + "flag"
						//debug("Found a call to flag.Parse")
						x.Args = []ast.Expr{&ast.Ident{Name: "os.Args[1:]"}}
					case "Flag":
					default:
						i.Name = p.name + "flag"
					}
				}
			}

		}
		return true
	})

	// Now we change any import names attached to package declarations.  We
	// just upcase it for now; it makes it easy to look in bbsh for things
	// we changed, e.g. grep -r bbsh Import is useful.
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			l := p.fset.Position(c.Pos()).Line
			if l != pos.Line {
				continue
			}
			if c.Text[0:9] == "// import" {
				c.Text = "// Import" + c.Text[9:]
			}
		}
	}
	return hasMain
}

func writeFile(path string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("error formating: %v", err)
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
	fullCode, err := imports.Process("commandline", code, &opts)
	if err != nil {
		return fmt.Errorf("bad parse %q: %v", string(code), err)
	}

	if err := ioutil.WriteFile(path, fullCode, 0644); err != nil {
		return fmt.Errorf("error writing to %q: %v", path, err)
	}
	return nil
}

func (p *Package) writeTemplate(path string, text string) error {
	var b bytes.Buffer
	t := template.Must(template.New("uroot").Parse(text))
	if err := t.Execute(&b, p.CommandTemplate()); err != nil {
		return fmt.Errorf("spec %v: %v", text, err)
	}

	return writeGoFile(path, b.Bytes())
}

func getPackage(opts BuildOpts, importPath string) (*Package, error) {
	p, err := opts.Env.Package(importPath)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(p.Dir)
	if !p.IsCommand() {
		return nil, fmt.Errorf("package %q is not a command and cannot be included in bb", name)
	}

	fset := token.NewFileSet()
	pars, err := parser.ParseDir(fset, p.Dir, nil, parser.ParseComments)
	if err != nil {
		log.Printf("can't parsedir %q: %v", p.Dir, err)
		return nil, nil
	}

	return &Package{
		pkg:  p,
		fset: fset,
		ast:  pars[p.Name],
		name: name,
	}, nil
}

func (b *bbBuilder) moveCommand(pkgPath string) error {
	p, err := getPackage(b.opts, pkgPath)
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

	var hasMain bool
	for filePath, sourceFile := range p.ast.Files {
		if hasMainFile := p.rewriteFile(b.opts, sourceFile); hasMainFile {
			hasMain = true
		}

		path := filepath.Join(pkgDir, filepath.Base(filePath))
		if err := writeFile(path, p.fset, sourceFile); err != nil {
			return err
		}
	}

	if !hasMain {
		return os.RemoveAll(pkgDir)
	}

	bbsetupPath := filepath.Join(b.bbshDir, "cmds", p.name, "bbsetup.go")
	if err := p.writeTemplate(bbsetupPath, bbsetupGo); err != nil {
		return err
	}

	cmdPath := filepath.Join(b.bbshDir, fmt.Sprintf("cmd_%s.go", p.name))
	if err := p.writeTemplate(cmdPath, cmdFunc); err != nil {
		return err
	}

	b.initMap += "\n\t\"" + p.name + "\": " + p.name + "Init,"

	// Add a symlink to our bbsh.
	return b.af.AddRecord(cpio.Symlink(filepath.Join("bbin", p.name), "/bbin/rush"))
}
