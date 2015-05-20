// bb converts standalone u-root tools to shell builtins.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
)

func x(a string, b string,) {
}

func main() {
	// src is the input for which we want to inspect the AST.
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
func main(s string) {
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.FuncDecl:
			s = fmt.Sprintf("%v", reflect.TypeOf(x.Type.Params.List[0].Type))
			if x.Name.Name == "main" {
				x.Name.Name = "cat"
				x.Type.Params.List = append(x.Type.Params.List, &ast.Field{Names: []*ast.Ident{&ast.Ident{Name:"a"}}, Type: &ast.Ident{Name: "string",}})
			}
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

}
