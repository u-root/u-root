// bb converts standalone u-root tools to shell builtins.
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
	"reflect"
	"text/template"

	"golang.org/x/tools/imports"
)

const cmdFunc = `package main
import "{{.CmdName}}"
func _builtin_{{.CmdName}}(c *Command) (err error) {
save := os.Args
defer func() {
os.Args = save
        if r := recover(); r != nil {
            err = errors.New(fmt.Sprintf("%v", r))
        }
return
    }()
os.Args = append([]string{c.cmd}, c.argv...)
{{.CmdName}}.Main()
return
}


func init() {
	addBuiltIn("{{.CmdName}}", _builtin_{{.CmdName}})
}
`

var (
	defaultCmd = []string{
		"cat",
		"cmp",
		"comm",
		"cp",
		"date",
		"dmesg",
		"echo",
		"freq",
		"grep",
		"ip",
		"ls",
		"mkdir",
		"mount",
		"netcat",
		"ping",
		"printenv",
		"rm",
		"seq",
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
	dumpAST = flag.Bool("d", false, "Dump the AST")
)

var config struct {
	Args     []string
	CmdName  string
	FullPath string
	Src      string
}

func oneFile(dir, s string, fset *token.FileSet, f *ast.File) error {
	// Inspect the AST and change all instances of main()
	isMain := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			x.Name.Name = config.CmdName
		case *ast.FuncDecl:
			if false {
				fmt.Printf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))
			}
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("Main")
				// Append a return.
				x.Body.List = append(x.Body.List, &ast.ReturnStmt{})
				isMain = true
			}

		case *ast.CallExpr:
			fmt.Fprintf(os.Stderr, "%v %v\n", reflect.TypeOf(n), n)
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
	fmt.Printf("%s", buf.Bytes())
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
	if err := os.MkdirAll(packageDir, 0666); err != nil {
		log.Fatalf("Can't create target directory: %v", err)
	}
	fset := token.NewFileSet()
	config.FullPath = path.Join(os.Getenv("UROOT"), "src/cmds", config.CmdName)
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
	flag.Parse()
	config.Args = flag.Args()
	if len(config.Args) == 0 {
		config.Args = defaultCmd
	}
	for _, v := range config.Args {
		// Yes, gross. Fix me.
		config.CmdName = v
		oneCmd()
	}
}
