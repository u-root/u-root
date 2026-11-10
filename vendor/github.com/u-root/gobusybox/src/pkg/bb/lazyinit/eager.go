// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

import (
	"maps"
	"slices"

	"golang.org/x/tools/go/packages"
)

// Reasons a package is eager that are not about its own source.
const (
	eagerStandardLibrary = "is a standard library package"
	eagerLinknameTarget  = "is named by a //go:linkname directive in another package"
	eagerAssemblyTarget  = "is named by an assembly file in another package"
	eagerImportedByEager = "is imported by an eager package"
)

// Plan is which of a busybox's packages initialize at program start and which
// wait for a command that imports them.
//
// The eager set is computed as a closure rather than decided one package at a
// time, because an eager package initializes before any deferred one and so
// may not import a deferred one. See TestNoEagerPackageImportsALazyOne, which
// is the property the rest of this exists to keep.
type Plan struct {
	// pkgs is every package reachable from the commands, by import path.
	pkgs map[string]*packages.Package
	// cmds are the import paths of the commands, in the order given.
	cmds []string
	// eager maps an eager package's import path to why it is eager. A path
	// absent from this map is deferred.
	eager map[string]string
}

// New decides, for the busybox built from cmds, which packages initialize at
// program start.
//
// With no options every package that can be deferred is, which is what the
// busybox could recover rather than what any particular build moves. A build
// passes DeferOnly to narrow that to the packages it was asked for.
func New(cmds []*packages.Package, opts ...Option) *Plan {
	p := &Plan{
		pkgs:  map[string]*packages.Package{},
		eager: map[string]string{},
	}
	for _, c := range cmds {
		p.cmds = append(p.cmds, c.PkgPath)
	}
	for _, c := range cmds {
		p.collect(c)
	}

	// Roots, in two kinds. A package's own source can refuse laziness, and a
	// two-argument //go:linkname in some *other* package can refuse it on
	// its behalf -- that is an edge the import graph does not have, so the
	// closure below could never reach the target from the puller.
	roots := map[string]string{}
	// In import path order. A target named by two packages keeps the
	// first reason recorded for it, so iterating the map directly would
	// make the reported reason depend on Go's randomized map order.
	for _, path := range slices.Sorted(maps.Keys(p.pkgs)) {
		pkg := p.pkgs[path]
		if isStandardLibrary(pkg) {
			roots[path] = eagerStandardLibrary
			continue
		}
		s := scan(pkg)
		// A refusal says "this package cannot be rewritten to initialize
		// later", and so makes it an eager root, which drags everything
		// it imports eager with it. That cannot be said of a command.
		// A command is rewritten whatever this decides -- its main
		// becomes registeredMain and its initialization registeredInit,
		// called at dispatch -- so it never initializes at program start
		// and the premise the closure rests on does not hold for it.
		//
		// Letting it hold cost the whole busybox: one command carrying
		// an assembly file made every library it reached eager, which
		// for a busybox of any size is all of them.
		//
		// Its linkname *targets* still join the roots, below. Those are
		// about a package the command may never import, which nothing
		// else would initialize.
		if s.refusal != "" && !p.isCommand(path) {
			roots[path] = s.refusal
		}
		// An assembly file naming another package reaches it without an
		// import edge, exactly as a two-argument //go:linkname does, so
		// the target joins the roots the same way. Refusing only the
		// package holding the assembly would leave the one it reads
		// into deferred and move the hazard rather than close it.
		for _, symbol := range s.asmSymbols {
			// Resolved the same way a linkname symbol is, because
			// it is the same question: which package owns this
			// name. A target outside this busybox cannot be
			// deferred by it anyway, so naming one is not an error.
			target := p.linknameTarget(symbol)
			if target == "" {
				continue
			}
			if _, taken := roots[target]; !taken {
				roots[target] = eagerAssemblyTarget
			}
		}
		for _, symbol := range s.linknameSymbols {
			// A target outside this busybox cannot be deferred by it
			// anyway, so naming one is not an error.
			target := p.linknameTarget(symbol)
			if target == "" {
				continue
			}
			if _, taken := roots[target]; !taken {
				roots[target] = eagerLinknameTarget
			}
		}
	}

	// Record every root's own reason before closing downward. A root is
	// often reachable from another root -- most of the standard library is
	// -- and the walk would then blame whichever importer happened to
	// arrive first, which is map iteration order and so differs between
	// runs of one input.
	for path, reason := range roots {
		p.eager[path] = reason
	}

	// Close downward: everything an eager package imports is eager too.
	walked := map[string]bool{}
	for path := range roots {
		p.markEager(path, walked)
	}

	// Narrowing runs last, on a settled eager set. DeferOnly needs to know
	// which packages could be deferred before it can say which of them are.
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// collect walks the import graph from pkg, recording every package once.
func (p *Plan) collect(pkg *packages.Package) {
	if _, seen := p.pkgs[pkg.PkgPath]; seen {
		return
	}
	p.pkgs[pkg.PkgPath] = pkg
	// By the package rather than by the key: the key is how pkg's source
	// spells the import, and sortedImports says why that is not the same
	// thing. Order does not matter here, because this only fills a map.
	for _, imp := range pkg.Imports {
		p.collect(imp)
	}
}

// markEager carries eagerness from path down to everything it imports.
//
// A package that has not already been given a reason is blamed on its
// importer. A package that has one keeps it, so that a package eager on its
// own account -- standard library, or a refusal, or named by a //go:linkname
// -- says so rather than blaming an importer that merely reached it first.
// That is why New records the roots' reasons before calling this.
//
// walked is separate from the reasons because the two stop at different
// points: a root already carries a reason but its imports still have to be
// visited.
func (p *Plan) markEager(path string, walked map[string]bool) {
	if walked[path] {
		return
	}
	walked[path] = true
	if _, already := p.eager[path]; !already {
		p.eager[path] = eagerImportedByEager
	}

	pkg, ok := p.pkgs[path]
	if !ok {
		return
	}
	for _, imp := range sortedImports(pkg) {
		p.markEager(imp, walked)
	}
}

// linknameTarget returns the package in this busybox that holds the symbol a
// two-argument //go:linkname names, or "" if no package here does.
//
// The symbol is a package path and a name joined by a dot, but both halves may
// contain dots of their own, so where the split falls is decided by which
// packages exist rather than by the text. example.com/victim.(*T).Ready is a
// method on a type in example.com/victim: its final dot separates the method
// name and its first separates nothing. So each dot is tried from the right
// and the first prefix that is a package wins, which is the longest one --
// a package path may itself end in a dotted segment, gopkg.in/yaml.v3 being
// the ordinary case, and gopkg.in/yaml may be in the same busybox.
//
// An unresolved symbol is not an error. x/sys/unix's directives name the
// standard library, which is eager already, and nothing requires a target to
// be in this build at all.
func (p *Plan) linknameTarget(symbol string) string {
	for i := len(symbol) - 1; i >= 0; i-- {
		if symbol[i] != '.' {
			continue
		}
		if _, ok := p.pkgs[symbol[:i]]; ok {
			return symbol[:i]
		}
	}
	return ""
}

// Eager reports whether the package at this import path initializes at program
// start rather than when a command that imports it is dispatched.
func (p *Plan) Eager(pkgPath string) bool {
	_, ok := p.eager[pkgPath]
	return ok
}

// EagerReason says why a package is eager, for the report at the end of a
// build. It is empty for a deferred package.
func (p *Plan) EagerReason(pkgPath string) string {
	return p.eager[pkgPath]
}

// sortedImports returns the paths of pkg's direct imports in import path
// order, so that two builds of one input agree.
//
// The path of each imported package, not the key it is filed under.
// go/packages keys Imports by "import paths appearing in the package's Go
// source files", and an import does not always resolve to the package of the
// same name: the standard library is vendored into GOROOT, so net/http's
// "golang.org/x/net/http2/hpack" resolves to the package whose path is
// "vendor/golang.org/x/net/http2/hpack". Loading net/http on Go 1.27 produces
// 28 edges like it.
//
// Everything else here -- pkgs, eager, Hooks, CallOrder -- is keyed by
// package path, so a key reaching any of them names the wrong package or none
// at all. It names the wrong one whenever the busybox also holds the module
// copy of a vendored path, which any command importing golang.org/x/net/...
// brings in. bbinternal's importName:194 is the same distinction from the
// other side: it has the package path and needs the key to write an import
// statement.
func sortedImports(pkg *packages.Package) []string {
	out := make([]string, 0, len(pkg.Imports))
	for _, imp := range pkg.Imports {
		out = append(out, imp.PkgPath)
	}
	slices.Sort(out)
	return out
}

// isStandardLibrary reports whether pkg is a standard library package, which
// gobusybox never rewrites.
//
// Standard library packages belong to no module and every other package
// belongs to one, so this is the go tool's own answer rather than a guess at
// one. go/packages fills Module in because findpkg asks for
// packages.NeedModule.
//
// The guess to avoid is whether the first element of the import path contains
// a dot. Module paths are not required to contain one -- `module nodot` is
// legal for anything not fetched by its path -- so that test takes every
// package in such a module for standard library. bb.collectDeps made exactly
// that mistake, and 7b24c96 replaced it with this; see
// TestAModuleWithoutADotIsNotStandardLibrary.
//
// A package that failed to load also has no module, but findpkg rejects those
// before a Plan is built, so the distinction does not arise here.
//
// No standard library package imports a vendored one, so making the whole of
// it eager drags nothing else along.
func isStandardLibrary(pkg *packages.Package) bool {
	return pkg.Module == nil
}
