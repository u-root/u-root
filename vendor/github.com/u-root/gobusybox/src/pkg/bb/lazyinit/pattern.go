// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazyinit

import (
	"regexp"
	"slices"
	"strings"
)

// wildcard is what a pattern uses to reach a subtree, spelled as the go
// command spells it. See `go help packages`.
const wildcard = "..."

// IsPattern reports whether a -lazy-init value selects a subtree rather than
// naming one package.
//
// This is not findpkg.IsGlob, which answers the same shape of question about
// the other pattern language here. Commands are selected with path.Match
// globs, whose "*" stops at a slash. -lazy-init selects packages already
// resolved into the busybox, so the only wildcard worth having is the one
// that crosses slashes.
func IsPattern(s string) bool {
	return strings.Contains(s, wildcard)
}

// matches turns a pattern into the test of whether one import path is in it.
//
// The rule is the go command's, from `go help packages`: "..." matches any
// string, including the empty one and ones containing slashes, and a trailing
// "/..." matches the path without it too, so grpc/... covers grpc itself.
//
// The go command's other special case, which keeps a wildcard out of a
// "vendor" path element, is left out. These are resolved import paths, so a
// vendor element here names a standard library package vendored into GOROOT,
// and those are eager whatever a pattern says.
func matches(pattern string) func(string) bool {
	// QuoteMeta first, so that a path's dots match dots and nothing else.
	// The wildcard is put back below, after the suffix case has had it.
	re := regexp.QuoteMeta(pattern)
	if trimmed, ok := strings.CutSuffix(re, `/\.\.\.`); ok {
		re = trimmed + `(/\.\.\.)?`
	}
	re = strings.ReplaceAll(re, `\.\.\.`, `.*`)

	// MustCompile cannot fail: every character of the pattern was escaped
	// by QuoteMeta, and the rest of the expression is written here.
	return regexp.MustCompile(`^` + re + `$`).MatchString
}

// Matching returns the import paths of this busybox's libraries that pattern
// selects, in import path order.
//
// Commands are left out. A subtree holding commands and libraries together is
// the ordinary shape of a repository, and a command is rewritten whatever a
// pattern says, so naming one as unmovable would be noise rather than news.
// Naming a command exactly is still an error, which CannotDefer reports.
func (p *Plan) Matching(pattern string) []string {
	in := matches(pattern)
	var out []string
	for path := range p.pkgs {
		if in(path) && !p.isCommand(path) {
			out = append(out, path)
		}
	}
	slices.Sort(out)
	return out
}
