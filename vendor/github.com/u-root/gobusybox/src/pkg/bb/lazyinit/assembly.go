// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

// middleDot and divisionSlash are what Go assembly writes where a symbol name
// has a "." and an import path has a "/".
//
// The assembler's lexer does not accept the ASCII characters in an identifier,
// so github.com/u-root/u-root/pkg/x.F is spelled
// github·com∕u-root∕u-root∕pkg∕x·F. See https://go.dev/doc/asm, "Symbols", and
// cmd/asm/internal/lex.Make, which does this substitution.
const (
	middleDot     = "\u00b7"
	divisionSlash = "\u2215"
)

// toolchainHeaders are the headers the toolchain supplies, whose contents are
// known not to name a Go symbol. Any other include is a file this has not
// read, and a symbol reference could be hiding in it.
var toolchainHeaders = map[string]bool{
	"textflag.h":    true,
	"funcdata.h":    true,
	"go_asm.h":      true,
	"go_tls.h":      true,
	"abi_amd64.h":   true,
	"asm_amd64.h":   true,
	"asm_arm64.h":   true,
	"asm_ppc64x.h":  true,
	"asm_riscv64.h": true,
}

// symbolRun is the characters an assembly symbol may hold.
const symbolRun = `[\w\x{00b7}\x{2215}]*`

var (
	// asmComment matches both comment forms in one pass. Stripping them in
	// two passes lets a "/*" inside one line comment pair with a "*/"
	// inside a later one and swallow every instruction between them.
	asmComment = regexp.MustCompile(`(?m://.*$)|(?s:/\*.*?\*/)`)

	// The assembler's preprocessor allows space after the "#" and needs
	// none before the quote, so "# define" and `#include"x.h"` are both
	// real directives.
	asmInclude = regexp.MustCompile(`(?m)^\s*#\s*include\b\s*[<"]?([^>"\s]*)[>"]?`)
	asmDefine  = regexp.MustCompile(`(?m)^\s*#\s*define\b`)

	// A symbol is a run of identifier characters, the middle dot and the
	// division slash, and nothing else: the assembler rejects ASCII ".",
	// "/" and "-" inside one. It is a reference only when it holds a
	// middle dot, which is what separates a symbol from its package.
	asmSymbol  = regexp.MustCompile(symbolRun + `\x{00b7}` + symbolRun)
	asmDefines = regexp.MustCompile(`(?m)^\s*(TEXT|GLOBL|DATA)\s+(` + symbolRun + `)`)
)

// asmName turns an assembly symbol into the name the linker uses, which is
// what //go:linkname writes and what linknameTarget resolves.
func asmName(s string) string {
	return strings.NewReplacer(middleDot, ".", divisionSlash, "/").Replace(s)
}

// stripComments removes both comment forms, leaving a space in place of each
// so that adjacent tokens do not fuse, and keeping the line count so that a
// directive cannot be hidden by joining it to the line above.
func stripComments(src string) string {
	return asmComment.ReplaceAllStringFunc(src, func(m string) string {
		return " " + strings.Repeat("\n", strings.Count(m, "\n"))
	})
}

// asmScan reports why a package's assembly forbids deferring its
// initialization, and the symbols its assembly names in other packages.
//
// The hazard is that the linker binds a symbol by name rather than through an
// import edge, so assembly can reach a package's state without anything
// having called its hook. Two shapes do it: reading a Go symbol this
// package's own assembly does not define, and naming a symbol in another
// package. Everything else an assembly file does is reached by calling a
// function, which goes through the import graph like any other call.
//
// Refuses unless provably safe. Assembly this cannot read in full keeps the
// package eager, so an unreadable file costs laziness rather than
// correctness. That is why #define refuses: the assembler has a
// preprocessor, and a macro can put a symbol reference in a line that does
// not appear to hold one.
//
// Every file is read even once a refusal is certain, because the returned
// symbols are edges the import graph does not have, exactly as a
// two-argument //go:linkname is, and the caller makes their packages eager
// the same way. Stopping at the first refusal would drop the symbols found
// after it, leaving the package the assembly reads into deferred and moving
// the hazard rather than closing it.
func asmScan(files []string) (refusal string, symbols []string) {
	for _, file := range files {
		if path.Ext(file) != ".s" {
			continue
		}
		b, err := os.ReadFile(file)
		if err != nil {
			refusal = cmp(refusal, fmt.Sprintf("has an assembly file that could not be read: %v", err))
			continue
		}
		src := stripComments(string(b))

		if asmDefine.MatchString(src) {
			refusal = cmp(refusal, "has an assembly file using #define")
		}
		for _, m := range asmInclude.FindAllStringSubmatch(src, -1) {
			if m[1] == "" || strings.Contains(m[1], "/") || !toolchainHeaders[m[1]] {
				refusal = cmp(refusal, "has an assembly file including "+m[1])
			}
		}

		// The symbols this file defines. A Go function whose body is
		// assembly is declared in the Go source and implemented here,
		// so reaching it means calling it. A method carries its
		// receiver, so ·T·M defines "T.M" rather than "T".
		defined := map[string]bool{}
		for _, m := range asmDefines.FindAllStringSubmatch(src, -1) {
			if name, ok := strings.CutPrefix(asmName(m[2]), "."); ok {
				defined[name] = true
			}
		}
		for _, m := range asmSymbol.FindAllString(src, -1) {
			name := asmName(m)
			if local, ok := strings.CutPrefix(name, "."); ok {
				if !defined[local] {
					refusal = cmp(refusal, "has an assembly file reading "+middleDot+local)
				}
				continue
			}
			symbols = append(symbols, name)
			refusal = cmp(refusal, "has an assembly file naming "+m)
		}
	}
	return refusal, symbols
}

// cmp keeps the first refusal, so the reason does not depend on which file of
// a package happened to be read last.
func cmp(have, now string) string {
	if have != "" {
		return have
	}
	return now
}
