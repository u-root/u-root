// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bbinternal is the internal API for busybox builds.
//
// It contains exported functions that are not for user consumption and not
// stable.
package bbinternal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/u-root/uio/cp"
)

// The Go spec defines the following grammar:
//
//	identifier = letter { letter | unicode_digit } .
//
// See also https://golang.org/ref/spec#Identifiers
var pnameRegex = regexp.MustCompile("[^a-zA-Z0-9_]+")

// ParseAST parses the given files for a package named name.
//
// Only files with a matching package statement will be part of the AST
// returned.
func ParseAST(name string, files []string) (*token.FileSet, []*ast.File, []string, error) {
	fset := token.NewFileSet()
	astFiles := make(map[string]*ast.File)
	for _, path := range files {
		if src, err := parser.ParseFile(fset, path, nil, parser.ParseComments); err == nil && src.Name.Name == name {
			astFiles[path] = src
		} else if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse AST in file %q: %v", path, err)
		}
	}

	// Did we parse anything?
	if len(astFiles) == 0 {
		return nil, nil, nil, fmt.Errorf("no valid `main` package files found in %v", files)
	}

	// The order of types.Info.InitOrder depends on this list of files
	// always being passed to conf.Check in the same order.
	sort.Strings(files)

	sortedFiles := make([]*ast.File, 0, len(astFiles))
	parsedFiles := make([]string, 0, len(astFiles))
	for _, name := range files {
		if f, ok := astFiles[name]; ok {
			sortedFiles = append(sortedFiles, f)
			parsedFiles = append(parsedFiles, name)
		}
	}
	return fset, sortedFiles, parsedFiles, nil
}

// CreateBBMainSource creates a bb Go command main.go that imports all given
// pkgs and writes the command to destDir.
//
// fset and files must be parsed bb template main.go, usually ./bbmain/cmd/main.go.
func CreateBBMainSource(fset *token.FileSet, files []*ast.File, pkgs []string, destDir string) error {
	if len(files) != 1 {
		return fmt.Errorf("bb cmd template is supposed to only have one file")
	}

	for _, pkg := range pkgs {
		astutil.AddNamedImport(fset, files[0], "_", pkg)
	}
	return writeFiles(destDir, fset, files)
}

// Package is a Go package.
type Package struct {
	// Name is the executable command name.
	//
	// In the standard Go tool chain, this is usually the base name of the
	// directory containing its source files.
	Name string

	// Pkg is the actual data about the package.
	Pkg *packages.Package

	// library is set for a package rewritten to initialize when a command
	// that imports it is dispatched, rather than for a command.
	//
	// A library keeps its package name and its import path, because
	// everything importing it goes on importing it by that path. A command
	// is renamed to a package so that its main becomes an ordinary
	// function, which is why the two rewrites differ at all.
	library bool

	// hookName is the exported function that runs a library's
	// initialization, empty for a command.
	hookName string

	// hookDone is the package-level boolean that makes the second and
	// later calls to the hook return immediately, empty for a command.
	hookDone string

	// initCount keeps track of what the next init's index should be.
	initCount uint

	// initCalls names the busyboxInitN functions that have to be called,
	// in the order Go would have run them: the hoisted package-level
	// variables first, then each original init body in file order.
	initCalls []string

	// mainFuncName is the name for the renamed main().
	mainFuncName string

	// init is the registeredInit function that calls busyboxInit0 and the
	// renamed bodies of the original init functions, in the right order.
	init *ast.FuncDecl

	// addedImports records, per file, the names this rewrite has given to
	// imports it added. go/types' scopes describe the source as loaded and
	// cannot know about them.
	addedImports map[*ast.File]map[string]bool

	// initAssigns is a map of assignment expression -> busyboxInitN function
	// call statement.
	//
	// That busyboxInitN should contain the assignment statement for the
	// appropriate assignment expression.
	//
	// types.Info.InitOrder keeps track of Initializations by Lhs name and
	// Rhs ast.Expr.  We reparent the Rhs in assignment statements in
	// busyboxInitN functions, so we use the Rhs as an easy key here.
	// types.Info.InitOrder + initAssigns can then easily be used to derive
	// the order of Stmts in the "real" init.
	//
	// The key Expr must also be the AssignStmt.Rhs[0].
	initAssigns map[ast.Expr]ast.Stmt
}

// NewPackage creates a new Package based on an existing packages.Package.
func NewPackage(name string, p *packages.Package) *Package {
	pp := &Package{
		// Name is the executable name.
		Name:        path.Base(name),
		Pkg:         p,
		initAssigns: make(map[ast.Expr]ast.Stmt),
	}

	pp.mainFuncName = pp.newFunctionName("registeredMain")

	// This will hold calls to busyboxInit0 and to each renamed original init.
	pp.init = &ast.FuncDecl{
		Name: ast.NewIdent(pp.newFunctionName("registeredInit")),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}
	return pp
}

// NewLibrary creates a Package for a library whose initialization is deferred
// to the commands that import it.
//
// Its hook is exported because the packages that call it are other packages.
// The name is chosen by the same collision check that picks registeredMain,
// since a library is arbitrary third-party source that may already have a
// function of that name.
func NewLibrary(p *packages.Package) *Package {
	pp := &Package{
		Name:        p.Name,
		Pkg:         p,
		library:     true,
		initAssigns: make(map[ast.Expr]ast.Stmt),
	}
	pp.hookName = pp.newFunctionName("BusyboxInit")
	pp.hookDone = pp.newFunctionName("busyboxInitDone")
	return pp
}

// HookName is the exported function that runs this library's initialization.
func (p *Package) HookName() string {
	return p.hookName
}

func (p *Package) nextInit(addToCallList bool) *ast.Ident {
	nextInitName := fmt.Sprintf("busyboxInit%d", p.initCount)
	for p.funcNameTaken(nextInitName) {
		p.initCount++
		nextInitName = fmt.Sprintf("busyboxInit%d", p.initCount)
	}
	i := ast.NewIdent(nextInitName)
	if addToCallList {
		// A library has no registeredInit: its hook is generated from
		// initCalls, which both kinds record.
		if p.init != nil {
			p.init.Body.List = append(p.init.Body.List, &ast.ExprStmt{X: &ast.CallExpr{Fun: i}})
		}
		p.initCalls = append(p.initCalls, nextInitName)
	}
	p.initCount++
	return i
}

// addVarInit moves one global initializer into a busybox init of its own,
// appended to f, and records the call to it in initAssigns.
//
// lhs holds every name the value initializes: one for `var a = f()`, and all
// of them for `var a, b = f()`, which must be assigned in a single statement
// because the value can only be evaluated once. value is the key Rewrite
// looks up to place the call in init order, so it must be the same expression
// go/types recorded, not a copy.
func (p *Package) addVarInit(f *ast.File, lhs []ast.Expr, value ast.Expr) {
	varInit := &ast.FuncDecl{
		Name: p.nextInit(false),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: lhs,
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{value},
				},
			},
		},
	}
	p.initAssigns[value] = &ast.ExprStmt{X: &ast.CallExpr{Fun: varInit.Name}}
	f.Decls = append(f.Decls, varInit)
}

// dropComments removes the comment groups lying between from and to.
//
// go/printer interleaves a file's comments by position, so a comment whose
// expression has moved is emitted wherever its old position falls, which for
// a generated declaration with no position of its own is any token at all.
// The ones in a moved region are the only comments the rewrite displaces, and
// generated source is the one place losing them costs nothing.
func dropComments(f *ast.File, from, to token.Pos) {
	kept := f.Comments[:0]
	for _, group := range f.Comments {
		if group.Pos() >= from && group.End() <= to {
			continue
		}
		kept = append(kept, group)
	}
	f.Comments = kept
}

// movedFrom is where the part of a var spec that hoisting takes away begins:
// after the names, or after the type when the spec writes one.
//
// The = and every value after it move into the generated functions, so the
// comments anywhere in that region are displaced, not only those inside one
// value. A comment between two values of `var a, b = x, y` sits inside
// neither.
func movedFrom(s *ast.ValueSpec) token.Pos {
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}

// importName finds the package path to import, given the go/types pkg path.
//
// E.g. go/types uses the fully vendored name of a package, such as
// github.com/u-root/u-root/vendor/golang.org/x/sys/unix. importName returns
// the name that should appear in the import statement for this package, which
// is golang.org/x/sys/unix.
//
// Since the only way this happens is through an implicit import, we know that
// somewhere in the dependency tree this package must exist, so we visit all
// dependencies looking for the go/types package path looking for a valid
// package import path statement.
func importName(p *packages.Package, typePkgPath string) string {
	var importName string
	// Go through all dependent packages.
	packages.Visit([]*packages.Package{p}, func(p *packages.Package) bool {
		// Yeah, packages.Visit already goes through all imports -- but
		// it does not give us the index of the p.Imports map, which is
		// the "import paths appearing in the package's Go source
		// files".
		for name, pkg := range p.Imports {
			if pkg.PkgPath == typePkgPath {
				importName = name
				return false
			}
		}
		return true
	}, nil)
	if len(importName) > 0 {
		return importName
	}
	if spl := strings.Split(typePkgPath, "/vendor/"); len(spl) > 1 {
		return spl[len(spl)-1]
	}
	// It doesn't appear. We'll go import it.
	return typePkgPath
}

// pkgImportNameTaken checks whether name would conflict with any existing
// imports in f or variable/const/func declarations in p.
//
// Import statements may conflict with import statements in other files in
// the same package.
func (p *Package) pkgImportNameTaken(name string, f *ast.File) bool {
	// package scope is all variable, const, and func names
	if p.Pkg.Types.Scope().Lookup(name) != nil {
		return true
	}

	// file scope is imports. Only imports from this file can conflict.
	// Imports in other files have no effect.
	if p.Pkg.TypesInfo.Scopes[f].Lookup(name) != nil {
		return true
	}

	// Scopes are go/types' record of the source as it was loaded, so they
	// know nothing of an import this rewrite has already added. Two
	// inferred variable types naming different packages that call
	// themselves the same thing -- a/foo and b/foo, which is ordinary --
	// would otherwise both be offered "foo":
	//
	//	lib.go:4:2: other declaration of foo
	//	lib.go:5:2: "example.com/b/foo" imported and not used
	return p.addedImports[f][name]
}

// recordImport notes a name this rewrite has given to an import in f, so that
// the next one in the same file does not take it as well.
func (p *Package) recordImport(name string, f *ast.File) {
	if p.addedImports == nil {
		p.addedImports = make(map[*ast.File]map[string]bool)
	}
	if p.addedImports[f] == nil {
		p.addedImports[f] = make(map[string]bool)
	}
	p.addedImports[f][name] = true
}

// funcNameTaken checks whether name would conflict with any
// import/variable/const/func declarations in all of p.
//
// Variable/const/func names may not conflict with import statements in
// other files of the same package!
func (p *Package) funcNameTaken(name string) bool {
	// package scope is all variable, const, and func names
	if p.Pkg.Types.Scope().Lookup(name) != nil {
		return true
	}

	// file scope is all imports
	for _, file := range p.Pkg.Syntax {
		if p.Pkg.TypesInfo.Scopes[file].Lookup(name) != nil {
			return true
		}
	}
	return false
}

// newFunctionName returns an unused function name in p with the prefix name.
func (p *Package) newFunctionName(name string) string {
	var i int
	proposed := name
	for p.funcNameTaken(proposed) {
		proposed = fmt.Sprintf("%s%d", name, i)
		i++
	}
	return proposed
}

// newImportName returns an unused import name in f/p with the prefix name.
func (p *Package) newImportName(name string, f *ast.File) string {
	var i int
	proposed := name
	for p.pkgImportNameTaken(proposed, f) {
		proposed = fmt.Sprintf("%s%d", name, i)
		i++
	}
	return proposed
}

// PackageName is teh name of the rewritten Go package.
func (p *Package) PackageName() string {
	return "bb" + pnameRegex.ReplaceAllString(p.Name, "")
}

func (p *Package) rewriteFile(f *ast.File) bool {
	hasMain := false

	// Change the package name declaration from main to the command's name.
	// Remove all non-alphanumeric characters except for underscore and ensure
	// starting with a letter. There are more valid identifiers though.
	//
	// A library keeps its own name and its own import path. Everything that
	// imports it names it by that path, and those import statements are not
	// rewritten; only a command is renamed, because a command's package is
	// main and two of those cannot sit in one binary.
	if !p.library {
		f.Name.Name = p.PackageName()
	}

	// Map of fully qualified package name -> imported alias in the file.
	importAliases := make(map[string]string)
	unaliasedImports := make(map[string]struct{})
	for _, impt := range f.Imports {
		importPath, err := strconv.Unquote(impt.Path.Value)
		if err != nil {
			panic(err)
		}

		if impt.Name != nil {
			importAliases[importPath] = impt.Name.Name
		} else {
			// We do not record the name of the package, because we
			// do not know it. However, `qualifier` will know it
			// because it's passed in.
			unaliasedImports[importPath] = struct{}{}
		}
	}

	// When the types.TypeString function translates package names, it uses
	// this function to map fully qualified package paths to a local alias,
	// if it exists.
	qualifier := func(pkg *types.Package) string {
		// pkg.Path() = fully vendored package name.
		// importPath = package name as it appears in `import` statement.
		importPath := importName(p.Pkg, pkg.Path())

		if name, ok := importAliases[importPath]; ok {
			return name
		}
		if _, ok := unaliasedImports[importPath]; ok {
			// The package name is NOT the last component of its path.
			return pkg.Name()
		}
		// When referring to self, don't use any package name.
		if pkg == p.Pkg.Types {
			return ""
		}

		// This type is not imported in this file yet.
		//
		// This typically happens when a derived global import uses a
		// type that was previously only implicitly used.
		//
		// E.g. if we call func Foo() *log.Logger like this:
		//
		//   var l = pkg.Foo()
		//
		// Then it's possible the `log` package was not referred to at
		// all previously, and we now need to add an import for log.
		importAlias := p.newImportName(pkg.Name(), f)
		astutil.AddNamedImport(p.Pkg.Fset, f, importAlias, importPath)
		// Make sure we do not add this import twice, and that the next
		// package calling itself the same thing is given another name.
		importAliases[importPath] = importAlias
		p.recordImport(importAlias, f)

		return importAlias
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// We only care about vars.
			if d.Tok != token.VAR {
				break
			}
			var specs []ast.Spec
			for _, spec := range d.Specs {
				s := spec.(*ast.ValueSpec)
				if s.Values == nil {
					specs = append(specs, s)
					continue
				}

				// A comment in the region that moves keeps the
				// position it had while the functions built
				// for it have none, so go/printer can flush
				// it between a function's name and its
				// parameter list. Before s.Values is cleared,
				// which is what s.End() reads.
				dropComments(f, movedFrom(s), s.End())

				// A var spec either gives one value per name,
				// or one multi-valued expression for all of
				// them -- `var a, b = f()` and the comma-ok
				// forms.
				//
				// go/types groups its initializers the same
				// way, and initAssigns is keyed by the value
				// expression, so create one busybox init per
				// value rather than one per name. Indexing the
				// values by the name's index only holds in the
				// first case.
				if len(s.Values) == len(s.Names) {
					for i, name := range s.Names {
						p.addVarInit(f, []ast.Expr{name}, s.Values[i])
					}
				} else {
					lhs := make([]ast.Expr, len(s.Names))
					for i, name := range s.Names {
						lhs[i] = name
					}
					p.addVarInit(f, lhs, s.Values[0])
				}
				s.Values = nil

				// Add the type to the global declaration
				// instead, since the value has moved away.
				if s.Type != nil {
					// Already written out, and shared by
					// every name in the spec.
					specs = append(specs, s)
					continue
				}
				// Inferred types need not agree across the
				// names -- `var f, err = os.Open(p)` -- and a
				// spec holds only one type, so give each name
				// its own spec. format.Node parenthesizes the
				// declaration when there is more than one.
				for _, name := range s.Names {
					// `var _ = f()` declares nothing, so
					// there is nothing to leave behind:
					// the assignment keeps the side
					// effect and the blank needs no
					// declaration to be assigned to.
					// Writing one out is at best noise
					// and at worst unbuildable, since the
					// inferred type may be one this file
					// cannot name.
					if name.Name == "_" {
						continue
					}
					typ := p.Pkg.TypesInfo.Defs[name].Type()
					specs = append(specs, &ast.ValueSpec{
						Names: []*ast.Ident{name},
						Type:  ast.NewIdent(types.TypeString(typ, qualifier)),
					})
				}
			}
			d.Specs = specs

		case *ast.FuncDecl:
			if !p.library && d.Recv == nil && d.Name.Name == "main" {
				d.Name.Name = p.mainFuncName
				hasMain = true
			}
			if d.Recv == nil && d.Name.Name == "init" {
				d.Name = p.nextInit(true)
			}
		}
	}

	// Now we change any import names attached to package declarations. We
	// just upcase it for now; it makes it easy to look in bbsh for things
	// we changed, e.g. grep -r bbsh Import is useful.
	//
	// TODO(chrisko): We don't have to do this anymore.
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// import") {
				c.Text = "// Import" + c.Text[9:]
			}
		}
	}
	return hasMain
}

// Hook is a deferred package's initialization, as the generated code reaches
// it: the path to import, and the exported function to call.
type Hook struct {
	// ImportPath is the package's import path, which is also where the
	// generated tree vendors it.
	ImportPath string

	// FuncName is the package's exported hook.
	FuncName string
}

// hoistInit moves p's package-level variable initializers and its init bodies
// into the numbered functions the rewrite emits as busyboxInitN, and returns
// the one holding the variables.
//
// The variables come first and the init bodies after, which is the order Go
// runs them in. Each variable's initializer is left in the file it was
// written in, because an expression may name an import and imports are scoped
// to a file; only the calls to them are gathered here, and a call names
// nothing but a package-scope function.
func (p *Package) hoistInit() (*ast.FuncDecl, *ast.File, error) {
	// This init holds all variable initializations.
	//
	// func busyboxInit0() {}
	varInit := &ast.FuncDecl{
		Name: p.nextInit(true),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	var mainFile *ast.File
	for _, sourceFile := range p.Pkg.Syntax {
		if hasMainFile := p.rewriteFile(sourceFile); hasMainFile {
			mainFile = sourceFile
		}
	}
	if !p.library && mainFile == nil {
		return nil, nil, fmt.Errorf("no main function found in package %q", p.Pkg.PkgPath)
	}

	// Add variable initializations to Init0 in the right order.
	for _, initStmt := range p.Pkg.TypesInfo.InitOrder {
		a, ok := p.initAssigns[initStmt.Rhs]
		if !ok {
			return nil, nil, fmt.Errorf("couldn't find init assignment %s", initStmt)
		}
		varInit.Body.List = append(varInit.Body.List, a)
	}
	return varInit, mainFile, nil
}

// RewriteLibrary rewrites p into destDir as a package whose initialization
// runs when a command that imports it is dispatched.
//
// deps are the hooks of p's direct deferred imports, which p's own hook calls
// before running anything of its own.
func (p *Package) RewriteLibrary(destDir string, deps []Hook) error {
	if !p.library {
		return fmt.Errorf("package %q was not created as a library", p.Pkg.PkgPath)
	}

	varInit, _, err := p.hoistInit()
	if err != nil {
		return err
	}

	// The gathered calls name only package-scope functions, so this may go
	// in any of the package's files. The first is chosen because
	// go/packages reports them in a fixed order, and two builds of one
	// input have to generate the same tree.
	if len(p.Pkg.Syntax) == 0 {
		return fmt.Errorf("package %q has no Go files to rewrite", p.Pkg.PkgPath)
	}
	p.Pkg.Syntax[0].Decls = append(p.Pkg.Syntax[0].Decls, varInit)

	if err := WritePkg(p.Pkg, destDir); err != nil {
		return err
	}
	return p.writeHook(destDir, deps)
}

// writeHook writes destDir/bbinit.go, holding p's hook and the boolean that
// keeps it to one run.
//
// It is a file of its own rather than more declarations in one of p's, so
// that the import statements the hook needs do not have to be merged into a
// file that may already import some of those paths under names of its own.
// It is also where someone debugging a generated tree looks first.
func (p *Package) writeHook(destDir string, deps []Hook) error {
	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated by gobusybox for %s. DO NOT EDIT.\n\n", p.Pkg.PkgPath)
	fmt.Fprintf(&b, "package %s\n\n", p.Pkg.Name)

	aliases := p.importAliases(&b, deps)

	fmt.Fprintf(&b, "// %s records that %s has run, so that the second and later\n", p.hookDone, p.hookName)
	fmt.Fprintf(&b, "// calls return immediately. A package reached by two paths through the\n")
	fmt.Fprintf(&b, "// import graph is initialized once, as Go initializes it once.\n")
	fmt.Fprintf(&b, "var %s bool\n\n", p.hookDone)

	fmt.Fprintf(&b, "// %s runs this package's initialization: the hooks of its direct\n", p.hookName)
	fmt.Fprintf(&b, "// imports first, then its own package-level variables, then the bodies of\n")
	fmt.Fprintf(&b, "// its init functions.\n")
	fmt.Fprintf(&b, "func %s() {\n", p.hookName)
	fmt.Fprintf(&b, "\tif %s {\n\t\treturn\n\t}\n", p.hookDone)
	fmt.Fprintf(&b, "\t%s = true\n", p.hookDone)
	for i, dep := range deps {
		fmt.Fprintf(&b, "\t%s.%s()\n", aliases[i], dep.FuncName)
	}
	for _, call := range p.initCalls {
		fmt.Fprintf(&b, "\t%s()\n", call)
	}
	fmt.Fprintf(&b, "}\n")

	return writeGoFile(filepath.Join(destDir, p.generatedFileName("bbinit")), []byte(b.String()))
}

// importAliases writes the import block for deps and returns the alias chosen
// for each, in the same order.
//
// Every import is aliased. The hook calls a package by that alias, so letting
// the path decide the name would mean knowing what each package calls itself,
// which is not always the last element of its path.
func (p *Package) importAliases(b *strings.Builder, deps []Hook) []string {
	if len(deps) == 0 {
		return nil
	}
	aliases := make([]string, len(deps))
	b.WriteString("import (\n")
	for i, dep := range deps {
		aliases[i] = p.newFunctionName(fmt.Sprintf("bbdep%d", i))
		fmt.Fprintf(b, "\t%s %s\n", aliases[i], strconv.Quote(dep.ImportPath))
	}
	b.WriteString(")\n\n")
	return aliases
}

// generatedFileName returns a file name in destDir that none of p's own
// files already uses.
//
// A package may have a file called bbinit.go of its own -- it is an ordinary
// name and nothing reserves it. Writing over it loses whatever was in it, and
// quietly: WritePkg has already written the package's real files by the time
// this runs, and the hoisted variable initializers are appended to whichever
// file sorts first, which "bbinit.go" very often does. The result compiles
// only by accident, and when it does not the error is `undefined:
// busyboxInit0`, which names neither the file nor the collision.
func (p *Package) generatedFileName(want string) string {
	taken := make(map[string]bool, len(p.Pkg.GoFiles)+len(p.Pkg.OtherFiles))
	for _, f := range p.Pkg.GoFiles {
		taken[filepath.Base(f)] = true
	}
	for _, f := range p.Pkg.OtherFiles {
		taken[filepath.Base(f)] = true
	}

	name := want + ".go"
	for i := 0; taken[name]; i++ {
		name = fmt.Sprintf("%s%d.go", want, i)
	}
	return name
}

// WritePkg writes p's files into destDir.
func WritePkg(p *packages.Package, destDir string) error {
	// TODO(hugelgupf):
	// - join errors
	// - seems a bit late to check for these errors, but works for now --
	//   should check when these packages are queried? first used?
	// - test
	if len(p.Errors) > 0 {
		return p.Errors[0]
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, fp := range p.OtherFiles {
		if err := cp.Copy(fp, filepath.Join(destDir, filepath.Base(fp))); err != nil {
			return fmt.Errorf("copy failed: %v", err)
		}
	}

	// The package directory, which all GoFiles share.
	pkgDir := filepath.Dir(p.GoFiles[0])

	written := writtenGoFiles(p.Fset, p.Syntax)
	for _, fp := range p.EmbedFiles {
		// pkg.go.dev/embed documentation: "The patterns are
		// interpreted relative to the package directory containing the
		// source file. The path separator is a forward slash, even on
		// Windows systems. Patterns may not contain ‘.’ or ‘..’ or
		// empty path elements, nor may they begin or end with a
		// slash."
		//
		// This means that the file must be a descendant of the package
		// directory and we can assume that all EmbedFiles share a base
		// path.
		relPath, err := filepath.Rel(pkgDir, fp)
		if err != nil {
			return err
		}
		// A package may embed one of its own Go files, as
		// golang.org/x/tools/go/analysis/passes/inline embeds doc.go.
		// writeFiles writes those from p.Syntax, so copying the
		// original here would only put a read-only module cache file
		// where the rewritten one goes.
		if written[relPath] {
			continue
		}
		os.MkdirAll(filepath.Join(destDir, filepath.Dir(relPath)), 0755)
		if err := cp.Copy(fp, filepath.Join(destDir, relPath)); err != nil {
			return fmt.Errorf("copy failed: %v", err)
		}
	}

	return writeFiles(destDir, p.Fset, p.Syntax)
}

func writeFiles(destDir string, fset *token.FileSet, files []*ast.File) error {
	// Write all files out.
	for _, file := range files {
		name := fset.File(file.Package).Name()

		path := filepath.Join(destDir, filepath.Base(name))
		if err := writeFile(path, fset, file); err != nil {
			return err
		}
	}
	return nil
}

// writtenGoFiles returns the names writeFiles writes, relative to destDir.
func writtenGoFiles(fset *token.FileSet, files []*ast.File) map[string]bool {
	names := make(map[string]bool, len(files))
	for _, file := range files {
		names[filepath.Base(fset.File(file.Package).Name())] = true
	}
	return names
}

// Rewrite rewrites p into destDir as a bb package, rewriting its init and main
// functions.
//
// bbImportPath is the import path to use for bbmain. deps are the hooks of
// the command's direct deferred imports, empty unless this build defers
// something.
func (p *Package) Rewrite(destDir, bbImportPath string, deps []Hook) error {
	varInit, mainFile, err := p.hoistInit()
	if err != nil {
		return err
	}

	// The deferred packages are initialized before anything of this
	// command's, because the command's own variables may read them. Go had
	// that order already: an imported package is finished in full before
	// the importing package's variables are initialized.
	//
	// The call is prepended rather than appended for the same reason, and
	// it names nothing but a package-scope function, so the import it
	// needs can live in the generated file that declares it.
	deferredInit := p.deferredInitName(deps)
	if deferredInit != "" {
		p.init.Body.List = append([]ast.Stmt{
			&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent(deferredInit)}},
		}, p.init.Body.List...)
	}

	// import bbmain "bbImportPath"
	importName := p.newImportName("bbmain", mainFile)
	astutil.AddNamedImport(p.Pkg.Fset, mainFile, importName, bbImportPath)

	// func init() {
	//   bbmain.Register("p.name", Init, Main)
	// }
	bbRegisterSelf := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.CallExpr{
					Fun: ast.NewIdent(fmt.Sprintf("%s.Register", importName)),
					Args: []ast.Expr{
						// name=
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote(p.Name),
						},
						// init=
						ast.NewIdent(p.init.Name.Name),
						// main=
						ast.NewIdent(p.mainFuncName),
					},
				}},
			},
		},
	}

	mainFile.Decls = append(mainFile.Decls, varInit, p.init, bbRegisterSelf)

	if err := WritePkg(p.Pkg, destDir); err != nil {
		return err
	}
	return p.writeDeferredInit(destDir, deferredInit, deps)
}

// deferredInitName is the name of the generated function that runs this
// command's deferred imports, empty when there are none -- which is every
// command in a busybox built without -lazy-init.
func (p *Package) deferredInitName(deps []Hook) string {
	if len(deps) == 0 {
		return ""
	}
	return p.newFunctionName("busyboxInitDeferred")
}

// writeDeferredInit writes destDir/bbinit.go, holding the function that runs
// the initialization of the deferred packages this command imports.
func (p *Package) writeDeferredInit(destDir, name string, deps []Hook) error {
	if name == "" {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated by gobusybox for %s. DO NOT EDIT.\n\n", p.Pkg.PkgPath)
	fmt.Fprintf(&b, "package %s\n\n", p.PackageName())

	aliases := p.importAliases(&b, deps)

	fmt.Fprintf(&b, "// %s runs the initialization of the packages this command\n", name)
	fmt.Fprintf(&b, "// imports whose initialization was deferred. Their own hooks reach the rest\n")
	fmt.Fprintf(&b, "// of the graph, so only the direct imports are named here.\n")
	fmt.Fprintf(&b, "func %s() {\n", name)
	for i, dep := range deps {
		fmt.Fprintf(&b, "\t%s.%s()\n", aliases[i], dep.FuncName)
	}
	fmt.Fprintf(&b, "}\n")

	return writeGoFile(filepath.Join(destDir, p.generatedFileName("bbinit")), []byte(b.String()))
}

func writeFile(path string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("error formatting Go file %q: %v", path, err)
	}
	return writeGoFile(path, buf.Bytes())
}

func writeGoFile(path string, code []byte) error {
	// Format the file. Do not fix up imports, as we only moved code around
	// within files.
	opts := imports.Options{
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: true,
	}
	code, err := imports.Process("commandline", code, &opts)
	if err != nil {
		return fmt.Errorf("bad parse while processing imports %q: %v", path, err)
	}

	if err := os.WriteFile(path, code, 0644); err != nil {
		return fmt.Errorf("error writing Go file to %q: %v", path, err)
	}
	return nil
}
