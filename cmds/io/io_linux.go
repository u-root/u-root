// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for io for Linux.
package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

func in(f io.ReadSeeker, addr uint64, data interface{}) error {
	_, err := f.Seek(int64(addr), 0)
	if err != nil {
		return fmt.Errorf("in: bad address %v: %v", addr, err)
	}
	return binary.Read(f, binary.LittleEndian, data)
}

func out(f io.WriteSeeker, addr uint64, data interface{}) error {
	_, err := f.Seek(int64(addr), 0)
	if err != nil {
		return fmt.Errorf("in: bad address %v: %v", addr, err)
	}
	return binary.Write(f, binary.LittleEndian, data)
}
