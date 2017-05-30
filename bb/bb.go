// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bb converts standalone u-root tools to shell builtins.
// It copies and converts a set of u-root utilities into a directory called bbsh.
// It assumes nothing; all files it needs are always copied, no matter what
// is in bbsh.
// bb needs to know where the uroot you are using is so it can find command source.
// UROOT=/home/rminnich/projects/u-root/u-root/
// bb needs to know the arch:
// GOARCH=amd64
// bb needs to know where the tools are, and they are in two places, the place it created them
// and the place where packages live:
// GOPATH=/home/rminnich/projects/u-root/u-root/bb/bbsh:/home/rminnich/projects/u-root/u-root
// bb needs to have a GOROOT
// GOROOT=/home/rminnich/projects/u-root/go1.5/go/
// There are no defaults.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"golang.org/x/tools/imports"
)

const (
	cmdFunc = `package main
import "github.com/u-root/u-root/bb/bbsh/cmds/{{.CmdName}}"
func _forkbuiltin_{{.CmdName}}(c *Command) (err error) {
os.Args = fixArgs("{{.CmdName}}", append([]string{c.cmd}, c.argv...))
{{.CmdName}}.Main()
return
}

func {{.CmdName}}Init() {
	addForkBuiltIn("{{.CmdName}}", _forkbuiltin_{{.CmdName}})
	{{.Init}}
}
`
	fixArgs = `
package main

func fixArgs(cmd string, args[]string) (s []string) {
	for _, v := range args {
		if v[0] == '-' {
			v = "-" + cmd + "." + v[1:]
		}
		s = append(s, v)
	}
	return
}
`
	initGo = `
package main
import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
        "syscall"

	"github.com/u-root/u-root/uroot"
)

func usage () {
	n := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s:\n", n)
	flag.VisitAll(func(f *flag.Flag) {
		if ! strings.HasPrefix(f.Name, n+".") {
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
	uroot.Rootfs()

	for n := range initMap {
		t := filepath.Join("/ubin", n)
		if err := os.Symlink("/init", t); err != nil {
			log.Printf("Symlink /init to %v: %v", t, err)
		}
	}
	if err := os.Symlink("/init", "/ubin/rush"); err != nil {
		log.Printf("Symlink /init to %v: %v", "/ubin/rush", err)
	}
        // spawn the first shell. We had been running the shell as pid 1
        // but that makes control tty stuff messy. We think.
        cloneFlags := uintptr(0)
	for _, v := range []string{"/inito", "/ubin/uinit", "/ubin/rush"} {
		cmd := exec.Command(v)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		var fd int
		cons, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			log.Printf("can't open /dev/tty: %v", err)
		} else {
			fd = int(cons.Fd())
			log.Printf("#### setting 0, 1, 2 to opened tty fd is %v", cons.Fd())
			cmd.Stdin, cmd.Stdout, cmd.Stderr = cons, cons, cons
		}
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
)

func debugPrint(f string, s ...interface{}) {
	log.Printf(f, s...)
}

func nodebugPrint(f string, s ...interface{}) {
}

const cmds = "cmds"

var (
	debug   = nodebugPrint
	cmdlist = []string{
		"src/github.com/u-root/u-root/cmds/*",
	}

	// fixFlag tells by existence if an argument needs to be fixed.
	// The value tells which argument.
	fixFlag = map[string]int{
		"Bool":        0,
		"BoolVar":     1,
		"Duration":    0,
		"DurationVar": 1,
		"Float64":     0,
		"Float64Var":  1,
		"Int":         0,
		"Int64":       0,
		"Int64Var":    1,
		"IntVar":      1,
		"String":      0,
		"StringVar":   1,
		"Uint":        0,
		"Uint64":      0,
		"Uint64Var":   1,
		"UintVar":     1,
		"Var":         1,
	}
	dumpAST = flag.Bool("D", false, "Dump the AST")
	initMap = "package main\nvar initMap = map[string] func() {\n"

	// commands to skip. init and rush should be obvious
	// builtin and script we skip as we have no toolchain in this mode.
	skip = map[string]bool{
		"builtin": true,
		"init":    true,
		"rush":    true,
		"script":  true,
	}
)

type Command struct {
	Gopath   string
	CmdName  string
	CmdPath  string
	Init     string
	FullPath string
}

var config struct {
	Commands []Command
	Src      string
	Cwd      string
	Bbsh     string

	Goroot    string
	Gosrcroot string
	Arch      string
	Goos      string
	// GOPATH is several paths, separated by :
	// We require that the first element be the
	// basic path that works with u-root.
	Gopath  string
	Gopaths []string
	TempDir string
	Go      string
	Debug   bool
	Fail    bool
}

func oneFile(c Command, dir, s string, fset *token.FileSet, f *ast.File) error {
	// Inspect the AST and change all instances of main()
	isMain := false
	c.Init = ""
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			x.Name.Name = c.CmdName
		case *ast.FuncDecl:
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("Main")
				// Append a return.
				x.Body.List = append(x.Body.List, &ast.ReturnStmt{})
				isMain = true
			}
			if x.Name.Name == "init" {
				x.Name.Name = fmt.Sprintf("Init")
				c.Init = c.CmdName + ".Init()"
			}

		case *ast.CallExpr:
			debug("%v %v\n", reflect.TypeOf(n), n)
			switch z := x.Fun.(type) {
			case *ast.SelectorExpr:
				// somebody tell me how to do this.
				sel := fmt.Sprintf("%v", z.X)
				// TODO: Need to have fixFlag and fixFlagVar
				// as the Var variation has name in the SECOND argument.
				if sel == "flag" {
					if ix, ok := fixFlag[z.Sel.Name]; ok {
						switch zz := x.Args[ix].(type) {
						case *ast.BasicLit:
							zz.Value = "\"" + c.CmdName + "." + zz.Value[1:]
						}
					}
				}
			}
		}
		return true
	})

	if *dumpAST {
		ast.Fprint(os.Stderr, fset, f, nil)
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	debug("%s", buf.Bytes())
	out := string(buf.Bytes())

	// fix up any imports. We may have forced the issue
	// with os.Args
	opts := imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	fullCode, err := imports.Process("commandline", []byte(out), &opts)
	if err != nil {
		log.Fatalf("bad parse: '%v': %v", out, err)
	}

	of := filepath.Join(dir, filepath.Base(s))
	if err := ioutil.WriteFile(of, []byte(fullCode), 0666); err != nil {
		log.Fatalf("%v\n", err)
	}

	// fun: must write the file first so the import fixup works :-)
	if isMain {
		// Write the file to interface to the command package.
		t := template.Must(template.New("cmdFunc").Parse(cmdFunc))
		var b bytes.Buffer
		if err := t.Execute(&b, c); err != nil {
			log.Fatalf("spec %v: %v\n", cmdFunc, err)
		}
		fullCode, err := imports.Process("commandline", []byte(b.Bytes()), &opts)
		if err != nil {
			log.Fatalf("bad parse: '%v': %v", out, err)
		}
		if err := ioutil.WriteFile(filepath.Join(config.Bbsh, "cmd_"+c.CmdName+".go"), fullCode, 0444); err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	return nil
}

func oneCmd(c Command) {
	// Create the directory for the package.
	// For now, ./cmds/<package name>
	packageDir := filepath.Join(config.Bbsh, "cmds", c.CmdName)
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		log.Fatalf("Can't create target directory: %v", err)
	}

	fset := token.NewFileSet()
	c.FullPath = filepath.Join(c.Gopath, c.CmdPath)
	p, err := parser.ParseDir(fset, c.FullPath, nil, 0)
	if err != nil {
		log.Printf("Can't Parsedir %v, %v", c.FullPath, err)
		return
	}

	for _, f := range p {
		for n, v := range f.Files {
			oneFile(c, packageDir, n, fset, v)
		}
	}
	initMap += "\n\t\"" + c.CmdName + "\":" + c.CmdName + "Init,"
}
func main() {
	var err error

	doConfig()

	if err := os.MkdirAll(config.Bbsh, 0755); err != nil {
		log.Fatalf("%v", err)
	}

	if len(flag.Args()) > 0 {
		cmdlist = flag.Args()
	}
	config.Commands = []Command{}

	for _, v := range cmdlist {
		debug("Check %v", v)
		for _, gp := range config.Gopaths {
			v = filepath.Join(gp, v)
			g, err := filepath.Glob(v)
			debug("v %v globs to %v, err %v", v, g, err)
			if err != nil {
				debug("tried to match path %v and cmd %v, failed", gp, v)
				continue
			}

			for i := range g {
				c := Command{Gopath: gp}
				c.CmdPath, err = filepath.Rel(gp, g[i])
				if err != nil {
					log.Fatalf("Can't take rel path of %v from %v? %v", g[i], gp, err)
				}

				config.Commands = append(config.Commands, c)
			}
		}
	}

	debug("config.Commands is %v", config.Commands)
	for _, c := range config.Commands {
		if skip[filepath.Base(c.CmdPath)] {
			continue
		}
		c.CmdName = filepath.Base(c.CmdPath)
		oneCmd(c)
	}

	if err := ioutil.WriteFile(filepath.Join(config.Bbsh, "init.go"), []byte(initGo), 0644); err != nil {
		log.Fatalf("%v\n", err)
	}
	// copy all shell files

	err = filepath.Walk(filepath.Join(config.Gopath, "src/github.com/u-root/u-root/cmds/rush"), func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(config.Bbsh, fi.Name()), b, 0644); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(config.Bbsh, "fixargs.go"), []byte(fixArgs), 0644); err != nil {
		log.Fatalf("%v\n", err)
	}

	initMap += "\n}"
	if err := ioutil.WriteFile(filepath.Join(config.Bbsh, "initmap.go"), []byte(initMap), 0644); err != nil {
		log.Fatalf("%v\n", err)
	}

	buildinit()
	ramfs()
}
