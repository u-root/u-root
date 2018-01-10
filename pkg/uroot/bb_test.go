package uroot

import (
	"go/ast"
	"go/importer"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
)

func TestBBBuild(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := BuildOpts{
		Env: golang.Default(),
		Packages: []string{
			"github.com/u-root/u-root/pkg/uroot/test/foo",
			"github.com/u-root/u-root/cmds/rush",
		},
		TempDir: dir,
	}
	af := NewArchiveFiles()
	if err := BBBuild(af, opts); err != nil {
		t.Error(err)
	}

	var mustContain = []string{
		"init",
		"bbin/rush",
		"bbin/foo",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q", name)
		}
	}

}

func findFile(filemap map[string]*ast.File, basename string) *ast.File {
	for name, f := range filemap {
		if filepath.Base(name) == basename {
			return f
		}
	}
	return nil
}

func TestPackageRewriteFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	p, err := getPackage(golang.Default(), "github.com/u-root/u-root/pkg/uroot/test/foo", importer.For("source", nil))
	if err != nil {
		t.Fatal(err)
	}

	f := findFile(p.ast.Files, "foo.go")
	if f == nil {
		t.Fatalf("file not found in files: %v", p.ast.Files)
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

	hasMain := p.rewriteFile(f)
	if !hasMain {
		t.Fatalf("foo.go should have main")
	}

	// Add variable initializations to Init0 in the right order.
	for _, initStmt := range p.typeInfo.InitOrder {
		a, ok := p.initAssigns[initStmt.Rhs]
		if !ok {
			t.Fatalf("couldn't find init assignment %s", initStmt)
		}
		varInit.Body.List = append(varInit.Body.List, a)
	}

	f.Decls = append(f.Decls, varInit)

	// Change the package name back to main.
	f.Name = ast.NewIdent("main")

	// Add init.
	f.Decls = append(f.Decls, p.init)
	f.Decls = append(f.Decls, &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("Init")}},
				&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("Main")}},
			},
		},
	})

	d, err := ioutil.TempDir("", "foo")
	if err != nil {
		t.Fatal(err)
	}

	if err := writeFile(filepath.Join(d, "foo.go"), p.fset, f); err != nil {
		t.Fatal(err)
	}

	if err := golang.Default().BuildDir(d, filepath.Join(d, "foo"), golang.BuildOpts{}); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(filepath.Join(d, "foo"))
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("foo failed: %v %v", string(o), err)
	}
}
