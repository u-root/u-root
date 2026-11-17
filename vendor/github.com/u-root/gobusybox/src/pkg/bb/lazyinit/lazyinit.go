// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lazyinit decides which of a busybox's packages may have their
// initialization deferred to the command that imports them, and in what order
// the deferred ones must run.
//
// It decides only. Nothing here rewrites source or writes a file; the result
// is the input to that work. See docs/adr/0004-library-init-runs-for-the-command-that-needs-it.md.
package lazyinit

import (
	"path"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Reasons laziness is refused, as they appear in the build's report. Phrased
// to name what was found rather than what it implies, because the reader is
// looking for it in their own source.
const (
	refusedObject   = "has a .syso object file"
	refusedLinkname = "has a //go:linkname directive"
)

// .S and .sx are assembly too, and are deliberately not looked at:
// go/packages reports them in OtherFiles only for a package that uses cgo,
// which this never sees, because makebb disables cgo and gomodbb inherits
// that. Measured on Go 1.27 with CGO_ENABLED=1, a package carrying empty.S
// and empty.sx reports OtherFiles=[] for both. Reading them would be a check
// that cannot fire rather than a guard.

// scanResult is what one package's own source says about deferring its
// initialization.
type scanResult struct {
	// refusal names what was found that forbids deferring this package, or
	// is empty if nothing was.
	refusal string
	// linknameSymbols are the remote symbols named by two-argument
	// //go:linkname directives in this package, verbatim. They are edges the
	// import graph does not have, so they are collected even though this
	// package is already refused: it is the *target* that the closure cannot
	// otherwise reach. Plan.linknameTarget turns a symbol into the package
	// that holds it.
	linknameSymbols []string
	// asmSymbols are the remote symbols this package's assembly names,
	// with the assembler's middle dot and division slash turned back into
	// "." and "/". Like linknameSymbols they are edges the import graph
	// does not have, and for the same reason they are collected even when
	// this package is already refused: it is the target the closure cannot
	// otherwise reach. Plan.linknameTarget turns a symbol into the package
	// that holds it, which is the same question for both.
	asmSymbols []string
}

// scan reads one package's files for the two things that let another package
// reach this one without importing it.
//
// Assembly is read from OtherFiles and directives from Syntax, both of which
// go/packages has already filtered by build constraints. That filtering is the
// point: a directive in a file excluded from this build cannot bind anything
// in it, and counting it would make packages eager on every platform for the
// sake of one.
func scan(p *packages.Package) scanResult {
	var r scanResult

	// Checked first, because it is the refusal a reader is least likely to
	// guess from their own source.
	if hasUnnameableVarType(p) {
		r.refusal = refusedUnnameableType
	}

	// A .syso is refused outright, being a precompiled object with no
	// source to read. Assembly is judged by what it does, because refusing
	// every package carrying a .s file costs far more than the hazard it
	// guards against. See ADR 0004.
	for _, f := range p.OtherFiles {
		if path.Ext(f) == ".syso" {
			r.refusal = refusedObject
			break
		}
	}
	// Read the assembly even when the package is already refused. The
	// symbols are about the packages it *names*, which nothing else
	// reaches, so skipping the scan would leave them deferred.
	refusal, symbols := asmScan(p.OtherFiles)
	r.asmSymbols = symbols
	if r.refusal == "" {
		r.refusal = refusal
	}

	for _, f := range p.Syntax {
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				symbol, ok := linknameSymbol(c.Text)
				if !ok {
					continue
				}
				r.refusal = refusedLinkname
				if symbol != "" && !slices.Contains(r.linknameSymbols, symbol) {
					r.linknameSymbols = append(r.linknameSymbols, symbol)
				}
			}
		}
	}
	return r
}

// linknameSymbol reports whether a comment is a //go:linkname directive, and
// returns the remote symbol it names if it is the two-argument form.
//
// The directive is only a directive when it is the whole of the comment, so
// prose mentioning it is not one.
//
// The symbol is returned as written. Splitting it into a package and a name
// cannot be done here: both halves may contain dots, and which dot separates
// them is decided by what packages exist, not by the text. A method's symbol
// carries its receiver -- example.com/victim.(*T).Ready -- so the final dot
// is the wrong one, and an import path contains dots of its own, so the first
// one is too. Plan.linknameTarget does it against the loaded package paths.
func linknameSymbol(text string) (string, bool) {
	if !strings.HasPrefix(text, "//go:linkname") {
		return "", false
	}
	fields := strings.Fields(text)
	if len(fields) == 0 || fields[0] != "//go:linkname" {
		return "", false
	}
	if len(fields) < 3 {
		return "", true
	}
	return fields[2], true
}
