// Copyright 2015-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	gover "go/version"
	"strings"
)

// normalizeGoVersion puts a version into the "go1.2.3" form go/version
// requires, so that go.mod directives ("1.26.6") and `go version` output
// ("go1.26.6") can be handled alike.
//
// Input that is not a Go version is returned unchanged, and go/version treats
// it as invalid rather than guessing.
func normalizeGoVersion(v string) string {
	if v == "" || strings.HasPrefix(v, "go") {
		return v
	}
	return "go" + v
}

// CompareGoVersions compares two Go versions, returning -1, 0 or 1.
//
// Accepted forms are go.mod directives ("1.21", "1.26.6", "1.21rc1") and
// `go version` output ("go1.26.6").
//
// Ordering is the go tool's own, via go/version, which is not the intuitive
// one. A language version sorts *below* its release candidates and releases:
//
//	1.26 < 1.26rc1 < 1.26.0 < 1.26.6
//
// and pre-release counters are numeric, so 1.26rc2 < 1.26rc10. Comparing
// components by hand gets all three of those backwards, and the result is a
// generated go.work directive too low for a member module, which is the
// failure this helper exists to prevent.
//
// Invalid versions sort below valid ones and equal to each other.
func CompareGoVersions(a, b string) int {
	return gover.Compare(normalizeGoVersion(a), normalizeGoVersion(b))
}

// MajorMinorGoVersion returns the major.minor language version of a Go
// version, e.g. "1.22" for both "go1.22.4" and "1.22.4".
//
// fallback is returned if v is not a Go version, or carries no minor component
// ("1" is not a language version, though go/version accepts it).
func MajorMinorGoVersion(v, fallback string) string {
	lang := strings.TrimPrefix(gover.Lang(normalizeGoVersion(v)), "go")
	if !strings.Contains(lang, ".") {
		return fallback
	}
	return lang
}
