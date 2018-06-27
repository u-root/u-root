// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ed is a simple line-oriented editor
//
// Synopsis:
//     dd
//
// Description:
//
// Options:
package main

import "io"

type Editor interface {
	Dot() int
	Move(int)
	Range() (int, int)
	Replace([]byte, int, int) (int, error)
	Read(io.Reader, int, int) (int, error)
	Write(io.Writer, int, int) (int, error)
	Sub(string, string, string, int, int) error
	Print(io.Writer, int, int) (int, error)
	IsDirty() bool
	Dirty(bool)
	Equal(f Editor) error
}
