// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import "strings"

// InOut is a stack-like interface used for IO.
// We are no longer sure we need it.
type InOut interface {
	// Push one or more strings onto the InOut
	Push(...string)
	// Pop a string fom the InOut.
	Pop() string
	// Pop all strings from the Inout.
	PopAll() []string
	// ReadAll implements io.ReadAll for an InOut.
	ReadAll() ([]byte, error)
	// Write implements io.Write for an InOut
	Write([]byte) (int, error)
}

// Line is used to implement an InOut based on an array of strings.
type Line struct {
	L []string
}

// NewLine returns an empty Line struct
func NewLine() InOut {
	return &Line{}
}

// Push implements Push for a Line
func (l *Line) Push(s ...string) {
	l.L = append(l.L, s...)
}

// Pop implements Pop for a Line
func (l *Line) Pop() (s string) {
	if len(l.L) == 0 {
		return s
	}
	s, l.L = l.L[len(l.L)-1], l.L[:len(l.L)-1]
	return s
}

// PopAll implements PopAll for a Line
func (l *Line) PopAll() (s []string) {
	s, l.L = l.L, []string{}
	return
}

// ReadAll implements ReadAll for a Line. There are no errors.
func (l *Line) ReadAll() ([]byte, error) {
	return []byte(strings.Join(l.PopAll(), "")), nil
}

// Write implements Write for a Line. There are no errors.
func (l *Line) Write(b []byte) (int, error) {
	l.Push(string(b))
	return len(b), nil
}
