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

        // spawn the first shell. We had been running the shell as pid 1
        // but that makes control tty stuff messy. We think.
        cloneFlags := uintptr(0)
	for _, v := range []string{"/inito", "/ubin/uinit", "/ubin/rush"} {
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
	bbsetupGo = `
package {{.CmdName}}

	import "flag"

	var {{.CmdName}}flag = flag.NewFlagSet("{{.CmdName}}", flag.ExitOnError)
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
	var pos token.Position
	isMain := false
	c.Init = ""
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		default:
			debug("%v %v\n", reflect.TypeOf(n), n)
		// This is rather gross. Arguably, so is the way that
		// Go has embedded build information in comments
		// ... this import comment attachment to a package
		// came in 1.4, a few years ago, and it only just bit
		// us with one file in upspin. So we go with gross.
		case *ast.Ident:
			// we assume the first identifier is the package id
			if !pos.IsValid() {
				pos = fset.Position(x.Pos())
				debug("Ident at %v", pos)
			}
		case *ast.File:
			x.Name.Name = c.CmdName
		case *ast.FuncDecl:
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("Main")
				// Append a return.
				x.Body.List = append(x.Body.List, &ast.ReturnStmt{})
				isMain = true
			}
			if x.Recv == nil && x.Name.Name == "init" {
				x.Name.Name = fmt.Sprintf("Init")
				c.Init = c.CmdName + ".Init()"
			}

		// rewrite use of the flag package.
		// The flag package uses a global variable to contain all flags.
		// This works poorly with the busybox mode, so as part of turning
		// commands into packages, we rewrite their use of flags to use a
		// package-private FlagSet.
		// in bbsetup.go, we add a var for the package flagset with params
		// (packagename, os.ExitOnError)
		// The ExitOnError may be a mistake, we'll have to see, since
		// packages are written to assume it returns. But it sure
		// is much handier since all our code currently follows
		// flag.Usage with os.Exit(1)
		// We rewrite arguments for calls to flag.Parse from () to
		// (os.Args[1:])
		// We rewrite all other uses of flag. to "commandname"+flag.
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
						i.Name = c.CmdName + "flag"
						debug("Found a call to flag.Parse")
						x.Args = []ast.Expr{&ast.Ident{Name: "os.Args[1:]"}}
					case "Flag":
					default:
						i.Name = c.CmdName + "flag"
					}
				}
			}

		}
		return true
	})

	// Now we change any import names attached to package declarations.
	// We just upcase it for now; it makes it easy to look in bbsh
	// for things we changed, e.g. grep -r bbsh Import is useful.
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			l := fset.Position(c.Pos()).Line
			if l != pos.Line {
				continue
			}
			if c.Text[0:9] == "// import" {
				c.Text = "// Import" + c.Text[9:]
			}
		}
	}

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
		var b bytes.Buffer
		// create the bbsetup.go file with all the random stuff we need
		t := template.Must(template.New("setup").Parse(bbsetupGo))
		if err := t.Execute(&b, c); err != nil {
			log.Fatalf("spec %v: %v\n", bbsetupGo, err)
		}
		// we don't yet needs imports.Process, but we'll see.
		if err := ioutil.WriteFile(filepath.Join(config.Bbsh, "cmds", c.CmdName, "bbsetup.go"), b.Bytes(), 0444); err != nil {
			log.Fatalf("%v\n", err)
		}
		b.Reset()
		// Write the file to interface to the command package.
		t = template.Must(template.New("cmdFunc").Parse(cmdFunc))
		if err := t.Execute(&b, c); err != nil {
			log.Fatalf("spec %v: %v\n", cmdFunc, err)
		}
		fullCode, err := imports.Process("commandline", []byte(b.Bytes()), &opts)
		if err != nil {
			log.Fatalf("Main commandline imports: bad parse: '%v': %v", out, err)
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
	p, err := parser.ParseDir(fset, c.FullPath, nil, parser.ParseComments)
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
	// In the bb case, the commands are built. In some cases, we want to
	// specify init= to be a u-root command on boot. Hence, it now makes sense
	// to have the ubin directory populated on boot, not by /init.
	l := filepath.Join(config.Bbsh, "ubin", c.CmdName)
	if err := os.Symlink("/init", l); err != nil {
		log.Fatalf("Symlinking %v -> /init: %v", l, err)
	}
}
func main() {
	var err error

	doConfig()

	if err := os.MkdirAll(filepath.Join(config.Bbsh, "ubin"), 0755); err != nil {
		log.Fatalf("%v", err)
	}

	if len(flag.Args()) > 0 {
		cmdlist = flag.Args()
	}
	config.Commands = []Command{}

	var numCmds int
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

			if len(g) == 0 {
				log.Println("=================================================")
				log.Printf("Warning: %v matched no paths", v)
				log.Println("=================================================")
			}
			for i := range g {
				c := Command{Gopath: gp}
				c.CmdPath, err = filepath.Rel(gp, g[i])
				if err != nil {
					log.Fatalf("Can't take rel path of %v from %v? %v", g[i], gp, err)
				}

				config.Commands = append(config.Commands, c)
				numCmds++
			}
		}
	}

	if numCmds == 0 {
		log.Print("=======================================================")
		log.Print("Warning: ZERO commands were added; check your arguments")
		log.Print("=======================================================")
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

	rush := filepath.Join(config.Bbsh, "ubin", "rush")
	if err := os.Symlink("/init", rush); err != nil {
		log.Printf("Symlink /init to %v: %v", rush, err)
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
