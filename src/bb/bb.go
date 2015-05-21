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
)

const cmdFunc = `func {{.CmdName}}(c *Command) (err error) {
save := os.Args
defer func() {
os.Args = save
        if r := recover(); r != nil {
            err = r.(error)
        }
return
    }()
os.Args = append([]string{c.cmd}, c.argv...)
{{.CmdName}}_main()
return
}

func init() {
	addBuiltIn("{{.CmdName}}", {{.CmdName}})
}
`

var config struct {
	CmdName string
}

func main() {
	src := `package main

func main() {
os.Exit(1)
panic(1)
}
`
	config.CmdName = "c"
	flag.Parse()
	a := flag.Args()
	os.Args = []string{"hi", "there"}
	if len(a) > 0 {
		b, err := ioutil.ReadFile(a[0])
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		src = string(b)
		// assume it ends in .go. Not much point otherwise.
		n := path.Base(a[0])
		config.CmdName = n[:len(n)-3]
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and change all instances of main()
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if false {
				fmt.Printf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))
			}
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("%v_main", config.CmdName)
			}
			// Append a return.
			x.Body.List = append(x.Body.List, &ast.ReturnStmt{})

			case *ast.CallExpr:
			fmt.Fprintf(os.Stderr, "%v %v\n", reflect.TypeOf(n), n)
			switch z := x.Fun.(type) {
			case *ast.SelectorExpr:
				// somebody tell me how to do this.
				sel := fmt.Sprintf("%v", z.X)
				if  sel == "os" && z.Sel.Name == "Exit" {
					x.Fun = &ast.Ident{Name: "panic"}
				}
			}
		}
		return true
	})
	
	if true {
		ast.Fprint(os.Stderr, fset, f, nil)
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	t := template.Must(template.New("cmdFunc").Parse(cmdFunc))
	var b bytes.Buffer
	if err := t.Execute(&b, config); err != nil {
		log.Fatalf("spec %v: %v\n", cmdFunc, err)
	}
	fmt.Printf("%v\n", b.String())

}
