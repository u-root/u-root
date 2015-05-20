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
	"reflect"
)

func x(a string, b string,) {
}

func main() {
	flag.Parse()
	a := flag.Args()
	b, err := ioutil.ReadFile(a[0])
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	// src is the input for which we want to inspect the AST.
	src := string(b)

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and change all instances of main()
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if false {fmt.Printf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))}
			if x.Name.Name == "main" {
				x.Name.Name = "cat"
				x.Type.Params.List = []*ast.Field{ &ast.Field{Names: []*ast.Ident{&ast.Ident{Name:"a"}}, Type: &ast.Ident{Name: "string",}}}
			}
			// we're done.
			return false
		}
		return true
	})
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

}
