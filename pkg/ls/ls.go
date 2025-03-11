// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ls implements formatting tools to list files like the Linux ls tool.
package ls

import (
	"fmt"
	"regexp"
)

// Matches characters which would interfere with ls's formatting.
var unprintableRe = regexp.MustCompile("[[:cntrl:]\n]")

// PrintableName returns a printable file name.
func (fi FileInfo) PrintableName() string {
	return unprintableRe.ReplaceAllLiteralString(fi.Name, "?")
}

// Stringer provides a consistent way to format FileInfo.
type Stringer interface {
	// FileString formats a FileInfo.
	FileString(fi FileInfo) string
}

// NameStringer is a Stringer implementation that just prints the name.
type NameStringer struct{}

// FileString implements Stringer.FileString and just returns fi's name.
func (ns NameStringer) FileString(fi FileInfo) string {
	return fi.PrintableName()
}

// QuotedStringer is a Stringer that returns the file name surrounded by qutoes
// with escaped control characters.
type QuotedStringer struct{}

// FileString returns the name surrounded by quotes with escaped control characters.
func (qs QuotedStringer) FileString(fi FileInfo) string {
	return fmt.Sprintf("%#v", fi.Name)
}

// LongStringer is a Stringer that returns the file info formatted in `ls -l`
// long format.
type LongStringer struct {
	Human bool
	Name  Stringer
}
