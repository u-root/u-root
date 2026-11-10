// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

import "slices"

// Hooks returns the import paths whose initialization pkgPath's own hook has
// to run before its body, which is its direct deferred imports in import path
// order.
//
// Direct imports only. A package may always import what it already imports, so
// calling them is legal by construction; reaching further is not. Another
// module's internal/ packages are where that breaks, and inside vendor/ the
// visibility rule still applies:
//
//	vendor/example.com/other/other.go:6:2: use of internal package
//	example.com/lib/internal/hidden not allowed
//
// grpc, protobuf and x/net all keep most of themselves under internal/, so
// this is the ordinary case rather than a corner. The recursion reaches the
// rest.
func (p *Plan) Hooks(pkgPath string) []string {
	pkg, ok := p.pkgs[pkgPath]
	if !ok {
		return nil
	}
	var out []string
	for _, imp := range sortedImports(pkg) {
		if !p.Eager(imp) {
			out = append(out, imp)
		}
	}
	return out
}

// CallOrder returns every package a dispatched command initializes, in the
// order the recursive hook calls reach them.
//
// It is a depth-first walk, with each package's direct imports taken in import
// path order, emitting a package after everything it imports. Taking the
// imports in a fixed order is what makes two builds of one input agree.
//
// It is not the order Go itself would have used, and cannot be made into it.
// The specification fixes a total order rather than only a partial one:
// "Given the list of all packages, sorted by import path, in each step the
// first uninitialized package in the list for which all imported packages (if
// any) are already initialized is initialized." A hook may call only the hooks
// of its own direct imports, per Hooks above, so a package's deferred imports
// are always initialized contiguously with it, and Go's rule interleaves them
// instead. TestCallOrderDiffersFromGosPackageInitializationOrder holds the
// smallest graph that tells the two apart, together with the order Go really
// produces for it.
//
// Both are orders in which every package follows its imports, which is what a
// program may rely on without saying so. Deferring initialization past program
// start, which is the whole feature, is a larger change to when a package's
// init runs than this is.
//
// Eager packages are absent: their init ran at program start, before any
// command was dispatched. So this is what the command's initialization costs,
// and the packages another command needs are not in it.
//
// Nothing generated needs this list -- the generated source is the graph walk
// in Hooks, not a list. It is the walk made explicit so that it can be
// asserted against the import graph, which a recursion spread across generated
// files cannot be.
func (p *Plan) CallOrder(cmdPath string) []string {
	if _, ok := p.pkgs[cmdPath]; !ok {
		return nil
	}
	var order []string
	seen := map[string]bool{}

	var visit func(path string)
	visit = func(path string) {
		if seen[path] {
			return
		}
		seen[path] = true
		for _, imp := range p.Hooks(path) {
			visit(imp)
		}
		order = append(order, path)
	}
	for _, imp := range p.Hooks(cmdPath) {
		visit(imp)
	}
	return slices.Clip(order)
}
