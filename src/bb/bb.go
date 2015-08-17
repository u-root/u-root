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
// GOPATH=/home/rminnich/projects/u-root/u-root/src/bb/bbsh:/home/rminnich/projects/u-root/u-root
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
	"path"
	"path/filepath"
	"reflect"
	"text/template"

	"golang.org/x/tools/imports"
)

const (
	cmdFunc = `package main
import "{{.CmdName}}"
func _builtin_{{.CmdName}}(c *Command) (err error) {
save := *flag.CommandLine
defer func() {
*flag.CommandLine = save
        if r := recover(); r != nil {
            err = errors.New(fmt.Sprintf("%v", r))
        }
return
    }()
flag.CommandLine.Init(c.cmd, flag.PanicOnError)
os.Args = fixArgs("{{.CmdName}}", append([]string{c.cmd}, c.argv...))
{{.CmdName}}.Main()
return
}

func init() {
	addBuiltIn("{{.CmdName}}", _builtin_{{.CmdName}})
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
	initFunc = `
package main
import "uroot"

func init() {
	uroot.Rootfs()
	return
}
`
)

func debugPrint(f string, s ...interface{}) {
	log.Printf(f, s...)
}

func nodebugPrint(f string, s ...interface{}) {
}

const cmds = "src/cmds"

var (
	debug      = nodebugPrint
	defaultCmd = []string{
		"cat",
		"cmp",
		"comm",
		"cp",
		"date",
		"dd",
		"dhcp",
		"dmesg",
		"echo",
		"freq",
		"grep",
		"ip",
		"kexec",
		"ls",
		"mkdir",
		"mount",
		"netcat",
		"ping",
		"printenv",
		"rm",
		"seq",
		"srvfiles",
		"tcz",
		"uname",
		"uniq",
		"unshare",
		"wc",
		"wget",
	}

	fixFlag = map[string]bool{
		"Bool":        true,
		"BoolVar":     true,
		"Duration":    true,
		"DurationVar": true,
		"Float64":     true,
		"Float64Var":  true,
		"Int":         true,
		"Int64":       true,
		"Int64Var":    true,
		"IntVar":      true,
		"String":      true,
		"StringVar":   true,
		"Uint":        true,
		"Uint64":      true,
		"Uint64Var":   true,
		"UintVar":     true,
		"Var":         true,
	}
	dumpAST = flag.Bool("D", false, "Dump the AST")
)

var config struct {
	Args     []string
	CmdName  string
	FullPath string
	Src      string
	Uroot    string
	Cwd      string
	Bbsh     string

	Goroot    string
	Gosrcroot string
	Arch      string
	Goos      string
	Gopath    string
	TempDir   string
	Go        string
	Debug     bool
	Fail      bool
}

func oneFile(dir, s string, fset *token.FileSet, f *ast.File) error {
	// Inspect the AST and change all instances of main()
	isMain := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			x.Name.Name = config.CmdName
		case *ast.FuncDecl:
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("Main")
				// Append a return.
				x.Body.List = append(x.Body.List, &ast.ReturnStmt{})
				isMain = true
			}

		case *ast.CallExpr:
			debug("%v %v\n", reflect.TypeOf(n), n)
			switch z := x.Fun.(type) {
			case *ast.SelectorExpr:
				// somebody tell me how to do this.
				sel := fmt.Sprintf("%v", z.X)
				if sel == "os" && z.Sel.Name == "Exit" {
					x.Fun = &ast.Ident{Name: "panic"}
				}
				if sel == "log" && z.Sel.Name == "Fatal" {
					x.Fun = &ast.Ident{Name: "panic"}
				}
				if sel == "log" && z.Sel.Name == "Fatalf" {
					nx := *x
					nx.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name = "fmt"
					nx.Fun.(*ast.SelectorExpr).Sel.Name = "Sprintf"
					x.Fun = &ast.Ident{Name: "panic"}
					x.Args = []ast.Expr{&nx}
					return false
				}
				if sel == "flag" && fixFlag[z.Sel.Name] {
					switch zz := x.Args[0].(type) {
					case *ast.BasicLit:
						zz.Value = "\"" + config.CmdName + "." + zz.Value[1:]
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

	of := path.Join(dir, path.Base(s))
	if err := ioutil.WriteFile(of, []byte(fullCode), 0666); err != nil {
		log.Fatalf("%v\n", err)
	}

	// fun: must write the file first so the import fixup works :-)
	if isMain {
		// Write the file to interface to the command package.
		t := template.Must(template.New("cmdFunc").Parse(cmdFunc))
		var b bytes.Buffer
		if err := t.Execute(&b, config); err != nil {
			log.Fatalf("spec %v: %v\n", cmdFunc, err)
		}
		fullCode, err := imports.Process("commandline", []byte(b.Bytes()), &opts)
		if err != nil {
			log.Fatalf("bad parse: '%v': %v", out, err)
		}
		if err := ioutil.WriteFile(path.Join("bbsh", "cmd_"+config.CmdName+".go"), fullCode, 0444); err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	return nil
}

func oneCmd() {
	// Create the directory for the package.
	// For now, ./src/<package name>
	packageDir := path.Join("bbsh", "src", config.CmdName)
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		log.Fatalf("Can't create target directory: %v", err)
	}
	fset := token.NewFileSet()
	config.FullPath = path.Join(config.Uroot, cmds, config.CmdName)
	p, err := parser.ParseDir(fset, config.FullPath, nil, 0)
	if err != nil {
		panic(err)
	}

	for _, f := range p {
		for n, v := range f.Files {
			oneFile(packageDir, n, fset, v)
		}
	}
}
func main() {
	var err error
	doConfig()

	for _, v := range config.Args {
		// Yes, gross. Fix me.
		config.CmdName = v
		oneCmd()
	}

	if err := ioutil.WriteFile(path.Join(config.Bbsh, "fixargs.go"), []byte(fixArgs), 0644); err != nil {
		log.Fatalf("%v\n", err)
	}

	if err := ioutil.WriteFile(path.Join(config.Bbsh, "init.go"), []byte(initFunc), 0644); err != nil {
		log.Fatalf("%v\n", err)
	}
	// copy all shell files
	err = filepath.Walk(path.Join(config.Uroot, cmds, "sh"), func(name string, fi os.FileInfo, err error) error {
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
		if err := ioutil.WriteFile(path.Join(config.Bbsh, fi.Name()), b, 0644); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal("%v", err)
	}

	buildinit()
	ramfs()
}
