// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package purgatory

// Purgatory abstracts a executable kexec purgatory in golang.
//
// What is a purgatory ? It is a binary object that runs between
// two kernels. See more https://lwn.net/Articles/582711/.
//
// We currently short circuit by generating purgatory executable
// via non-go toolchain. See generation logic from genpurg.go and
// finall generated golang purgatories code in asm.go
//
// See doc.go for more reading.
type Purgatory struct {
	// Name is a human readable alis to this purgatory executable.
	Name string
	// Hexdump is a hexdump of the executabl.
	Hexdump string
	// Code is the executable code in bytes slice that will be loaded
	// in memory as a kexec segment during kexec load.
	Code []byte
}
