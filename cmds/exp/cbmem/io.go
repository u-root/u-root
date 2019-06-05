// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"io"
	"log"
)

// readOneSize reads an entry of any type. This Size variant is for
// the console log only, though we know of no case in which it is
// larger than 1M. We really want the SectionReader as a way to ReadAt
// for the binary.Read. Any meaningful limit will be enforced by the kernel.
func readOneSize(r io.ReaderAt, i interface{}, o int64, n int64) {
	err := binary.Read(io.NewSectionReader(r, o, n), binary.LittleEndian, i)
	if err != nil {
		log.Fatalf("Trying to read section for %T: %v", r, err)
	}
}

// readOneSize reads an entry of any type, limited to 64K.
func readOne(r io.ReaderAt, i interface{}, o int64) {
	readOneSize(r, i, o, 65536)
}
