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
	"reflect"
	"text/template"
)

const cmdFunc = `func {{.CmdName}}(cmd string, args ...string) {
save := os.Args
os.Args = append([]string{cmd}, args...)
{{.CmdName}}_main()
os.Args = save
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
fmt.Printf("%v\n", os.Args)
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
FIXHERE path.base ec.
		config.CmdName = a[0][:len(a[0])-3]
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}
	
	// Inspect the AST and change all instances of main()
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			// placeholder in case we want to add something.
			//fmt.Printf("-->file! %v d %v   ..>\n", x, reflect.TypeOf(x.Decls))
			return true
		case *ast.FuncDecl:
			if false {fmt.Printf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))}
			if x.Name.Name == "main" {
				x.Name.Name = fmt.Sprintf("%vmain", config.CmdName)
				// don't rewrite the param list now, but leave this here so we remember
				// how it's done.
				if false {
					x.Type.Params.List = []*ast.Field{ 
						&ast.Field{Names: []*ast.Ident{&ast.Ident{Name:"cmd"}}, Type: &ast.Ident{Name: "string",}},
						&ast.Field{Names: []*ast.Ident{&ast.Ident{Name:"args"}}, Type: &ast.Ident{Name: "...string",}},				}
				}
			}
			
			return true
			// someday, we'll know how to change all instances of os.Args to something else. Someday.
		// case *ast.SelectorExpr:
		// 	if true {
		// 		fmt.Printf("%v\n", x)//reflect.TypeOf(x.Type.Params.List[0].Type))
		// 		fmt.Printf("%v %v %v\n", x.X, reflect.TypeOf(x.X), x.X.(*ast.Ident).Name)//reflect.TypeOf(x.Type.Params.List[0].Type))
		// 		fmt.Printf("%v\n", reflect.TypeOf(x))//reflect.TypeOf(x.Type.Params.List[0].Type))
		// 	}
		// 	// This is hoky.Need to check for whether this is a package. One thing at a time.
		// 	if x.X.(*ast.Ident).Name != "os" {
		// 		return true
		// 	}
		// 	if x.Sel.Name != "Arg" {
		// 	}
		// 	fmt.Printf("Got a hit on os.arg\n")
		// 	x.X.(*ast.Ident).Name = "main"
		// 	fmt.Printf("can set %v\n", reflect.ValueOf(x).CanSet())
		// 	return true
   // 123  .  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
   // 124  .  .  .  .  .  .  .  .  .  0: *ast.SelectorExpr {
   // 125  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   // 126  .  .  .  .  .  .  .  .  .  .  .  NamePos: src.go:18:9
   // 127  .  .  .  .  .  .  .  .  .  .  .  Name: "os"
   // 128  .  .  .  .  .  .  .  .  .  .  .  Obj: nil
   // 129  .  .  .  .  .  .  .  .  .  .  }
   // 130  .  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
   // 131  .  .  .  .  .  .  .  .  .  .  .  NamePos: src.go:18:12
   // 132  .  .  .  .  .  .  .  .  .  .  .  Name: "Args"
   // 133  .  .  .  .  .  .  .  .  .  .  .  Obj: nil
   // 134  .  .  .  .  .  .  .  .  .  .  }
   // 135  .  .  .  .  .  .  .  .  .  }
   // 136  .  .  .  .  .  .  .  .  }

		}
		return true
	})

	if true {ast.Fprint(os.Stderr, fset, f, nil)}
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
