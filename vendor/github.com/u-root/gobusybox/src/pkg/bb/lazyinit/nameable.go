// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// refusedUnnameableType is why a package whose hoisted variables cannot be
// declared is left alone.
const refusedUnnameableType = "has a package-level variable whose inferred type cannot be named outside the package that defines it"

// hasUnnameableVarType reports whether any of p's package-level variables has
// an inferred type that the rewrite could not write down.
//
// Hoisting a variable splits `var V = expr` into a declaration and an
// assignment, and the declaration has to state the type that was inferred.
// When that type is unexported and belongs to another package there is no way
// to state it:
//
//	package other
//	type client struct{ ... }
//	func New() *client
//
//	package lib
//	var V = other.New()     // rewritten as: var V *other.client
//
// which the compiler refuses with "name client not exported by package
// other". There is no spelling that works -- Go has no way to declare a
// variable of an unnameable type and assign it later -- so the only answers
// are to leave the package's initialization where it is, or to emit source
// that does not build. This picks the first, and the build reports it like
// any other refusal.
//
// Only variables with no written type are at risk. One that states its own
// type keeps it, and a type the source could already name is a type the
// rewrite can name.
//
// This is not a defect of deferring: `bbinternal.Package.Rewrite` has always
// hoisted a command's variables the same way and fails the same way on the
// same input. What changes is the chance of meeting it, because commands are
// hand-written and libraries are whatever the module graph drags in. See #9.
func hasUnnameableVarType(p *packages.Package) bool {
	if p.TypesInfo == nil || p.Types == nil {
		return false
	}
	for _, f := range p.Syntax {
		for _, decl := range f.Decls {
			d, ok := decl.(*ast.GenDecl)
			if !ok || d.Tok != token.VAR {
				continue
			}
			for _, spec := range d.Specs {
				s, ok := spec.(*ast.ValueSpec)
				// A spec with no values is not hoisted, and
				// one with a written type keeps it.
				if !ok || s.Values == nil || s.Type != nil {
					continue
				}
				for _, name := range s.Names {
					// `var _ = f()` declares nothing, so
					// the rewrite leaves no declaration
					// behind and no type has to be
					// written. go/types does record an
					// object for the blank identifier, so
					// this is a decision rather than a
					// nil check -- though it is that too,
					// since Defs is documented to hold
					// nil for identifiers that denote no
					// object.
					if name.Name == "_" {
						continue
					}
					def := p.TypesInfo.Defs[name]
					if def == nil {
						continue
					}
					if unnameableOutside(def.Type(), p.Types, map[types.Type]bool{}) {
						return true
					}
				}
			}
		}
	}
	return false
}

// unnameableOutside reports whether writing t as source inside home would have
// to name something unexported that belongs to another package.
//
// It walks the type constructors a declaration spells out and stops at every
// named type, because that is where printing stops too: `[]other.Thing` names
// Thing and says nothing about its fields. An anonymous struct or interface
// is the exception, since printing one does spell out its members.
func unnameableOutside(t types.Type, home *types.Package, seen map[types.Type]bool) bool {
	if t == nil || seen[t] {
		return false
	}
	seen[t] = true

	switch t := t.(type) {
	case *types.Named:
		obj := t.Obj()
		if obj == nil {
			return false
		}
		if obj.Pkg() != nil && obj.Pkg() != home && !obj.Exported() {
			return true
		}
		// A generic type's arguments are printed too.
		for i := 0; i < t.TypeArgs().Len(); i++ {
			if unnameableOutside(t.TypeArgs().At(i), home, seen) {
				return true
			}
		}
		return false

	case *types.Pointer:
		return unnameableOutside(t.Elem(), home, seen)
	case *types.Slice:
		return unnameableOutside(t.Elem(), home, seen)
	case *types.Array:
		return unnameableOutside(t.Elem(), home, seen)
	case *types.Chan:
		return unnameableOutside(t.Elem(), home, seen)
	case *types.Map:
		return unnameableOutside(t.Key(), home, seen) || unnameableOutside(t.Elem(), home, seen)

	case *types.Signature:
		return tupleIsUnnameable(t.Params(), home, seen) || tupleIsUnnameable(t.Results(), home, seen)

	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			if unnameableOutside(t.Field(i).Type(), home, seen) {
				return true
			}
		}
		return false

	case *types.Interface:
		for i := 0; i < t.NumMethods(); i++ {
			if unnameableOutside(t.Method(i).Type(), home, seen) {
				return true
			}
		}
		return false
	}
	return false
}

// tupleIsUnnameable is unnameableOutside over a function's parameters or
// results.
func tupleIsUnnameable(tup *types.Tuple, home *types.Package, seen map[types.Type]bool) bool {
	if tup == nil {
		return false
	}
	for i := 0; i < tup.Len(); i++ {
		if unnameableOutside(tup.At(i).Type(), home, seen) {
			return true
		}
	}
	return false
}
