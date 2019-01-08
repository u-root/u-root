// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!arm,!arm64

package io

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/ubinary"
)

var memPath = "/dev/mem"

func pathRead(path string, addr int64, data interface{}) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, os.SEEK_SET); err != nil {
		return err
	}
	return binary.Read(f, ubinary.NativeEndian, data)
}

// Read reads data from physical memory at address addr. data must be one of:
// *uint8, *uint16, *uint32, or *uint64. On x86 platforms, this uses the
// seek+read syscalls. On arm platforms, this uses mmap.
func Read(addr int64, data interface{}) error {
	switch data.(type) {
	case *uint8, *uint16, *uint32, *uint64:
	default:
		return fmt.Errorf("cannot read type %T", data)
	}
	return pathRead(memPath, addr, data)
}

func pathWrite(path string, addr int64, data interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, os.SEEK_SET); err != nil {
		return err
	}
	return binary.Write(f, ubinary.NativeEndian, data)
}

// Write writes data to physical memory at address addr. data must be one of:
// uint8, uint16, uint32, or uint64. On x86 platforms, this uses the seek+read
// syscalls. On arm platforms, this uses mmap.
func Write(addr int64, data interface{}) error {
	switch data.(type) {
	case uint8, uint16, uint32, uint64:
	default:
		return fmt.Errorf("cannot write type %T", data)
	}
	return pathWrite(memPath, addr, data)
}
