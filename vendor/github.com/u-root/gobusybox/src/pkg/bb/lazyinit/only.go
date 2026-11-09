// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

// eagerNotSelected is why a package that could have been deferred is not.
const eagerNotSelected = "is not selected by -lazy-init"

// Option narrows which of a busybox's packages a Plan defers.
type Option func(*Plan)

// DeferOnly defers the packages the patterns select and the packages that
// import them, and nothing else.
//
// Without it a Plan defers everything it can, which is the question "what
// could this busybox recover". DeferOnly is the answer to the different
// question a build asks: which packages is this build moving. The patterns
// are the dial, and a package they select gives up its package-level
// variables and its init together.
//
// A pattern is a whole import path, or a path holding "..." to reach a
// subtree. Matching says which packages one selects.
//
// A pattern that selects no package in this busybox, or only packages that
// are eager for a reason of their own, is ignored here rather than reported.
// Deferred says which it was, so the caller can refuse a build whose
// -lazy-init did nothing; a silent no-op is exactly the failure this work has
// to avoid.
//
// # Why the importers come too
//
// An eager package initializes before any deferred one, so it may not import
// one. Naming a package therefore defers every package that can reach it, not
// just the one selected. That closure can never meet a package that has to
// stay eager: the eager set New computes is already closed downward, so a
// package importing a deferrable one is deferrable itself.
func DeferOnly(patterns []string) Option {
	return func(p *Plan) {
		deferred := map[string]bool{}
		importers := p.importers()
		for _, pattern := range patterns {
			for _, path := range p.Matching(pattern) {
				if !p.Deferred(path) {
					continue
				}
				markImporters(path, importers, deferred)
			}
		}

		for path := range p.pkgs {
			// A command is not a library and is rewritten whatever
			// the patterns say, so it is never eager. Nothing
			// imports a main package, so leaving commands out
			// cannot leave an eager package importing a lazy one.
			if deferred[path] || p.Eager(path) || p.isCommand(path) {
				continue
			}
			p.eager[path] = eagerNotSelected
		}
	}
}

// importers maps each package's import path to the paths importing it
// directly, which is the import graph with every edge reversed.
func (p *Plan) importers() map[string][]string {
	rev := map[string][]string{}
	for path, pkg := range p.pkgs {
		for _, imp := range sortedImports(pkg) {
			rev[imp] = append(rev[imp], path)
		}
	}
	return rev
}

// markImporters records path and everything that can reach it through
// imports.
func markImporters(path string, importers map[string][]string, seen map[string]bool) {
	if seen[path] {
		return
	}
	seen[path] = true
	for _, importer := range importers[path] {
		markImporters(importer, importers, seen)
	}
}

// isCommand reports whether this import path is one of the busybox's
// commands.
func (p *Plan) isCommand(pkgPath string) bool {
	for _, cmd := range p.cmds {
		if cmd == pkgPath {
			return true
		}
	}
	return false
}

// Deferred reports whether the package at this import path is a library in
// this busybox whose initialization waits for a command that imports it.
//
// It is not the negation of Eager. A path naming no package here is neither,
// and neither is a command: a command is rewritten whatever the patterns say,
// its initialization already deferred to its own dispatch by registeredInit.
// CannotDefer tells the three apart.
func (p *Plan) Deferred(pkgPath string) bool {
	_, known := p.pkgs[pkgPath]
	return known && !p.Eager(pkgPath) && !p.isCommand(pkgPath)
}

// CannotDefer says why the package at this import path will not have its
// initialization deferred, and is empty when it will.
//
// The three refusals read differently because they send the reader to
// different places: a path that is not here is usually a typo, a command is a
// misunderstanding of what the flag is for, and an eager package is a rule
// doing its job.
func (p *Plan) CannotDefer(pkgPath string) string {
	if _, known := p.pkgs[pkgPath]; !known {
		return "is not imported by any command in this busybox"
	}
	if p.isCommand(pkgPath) {
		return "is a command, whose initialization already waits for its own dispatch"
	}
	return p.EagerReason(pkgPath)
}
