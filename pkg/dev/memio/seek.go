// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!arm,!arm64

package memio

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/ubinary"
)

var memPath = "/dev/mem"

func pathRead(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Read(f, ubinary.NativeEndian, data)
}

// Read reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls. On arm platforms, this uses mmap.
func Read(addr int64, data UintN) error {
	return pathRead(memPath, addr, data)
}

func pathWrite(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Write(f, ubinary.NativeEndian, data)
}

// Write writes data to physical memory at address addr. On x86 platforms, this
// uses the seek+read syscalls. On arm platforms, this uses mmap.
func Write(addr int64, data UintN) error {
	return pathWrite(memPath, addr, data)
}
