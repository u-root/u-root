// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,arm linux,arm64

package memio

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	pageSize = int64(syscall.Getpagesize())
)

// mmap aligns the address and maps multiple pages when needed.
func mmap(f *os.File, addr int64, size int64, prot int) (mem []byte, offset int64, err error) {
	if addr+size <= addr {
		return nil, 0, fmt.Errorf("invalid address for size %#x", size)
	}
	page := addr &^ (pageSize - 1)
	offset = addr - page
	mapSize := offset + size
	mem, err = syscall.Mmap(int(f.Fd()), int64(page), int(mapSize), prot, syscall.MAP_SHARED)
	return
}

// Read reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls. On arm platforms, this uses mmap.
func Read(addr int64, data UintN) error {
	f, err := os.OpenFile(memPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	mem, offset, err := mmap(f, addr, data.Size(), syscall.PROT_READ)
	if err != nil {
		return fmt.Errorf("Reading %#x/%d: %v", addr, data.Size(), err)
	}
	defer syscall.Munmap(mem)

	// MMIO makes this a bit tricky. Reads must be conducted in one load
	// operation. Review the generated assembly to make sure.
	if err := data.read(unsafe.Pointer(&mem[offset])); err != nil {
		return fmt.Errorf("Reading %#x/%d: %v", addr, data.Size(), err)
	}
	return nil
}

// Write writes data to physical memory at address addr. On x86 platforms, this
// uses the seek+read syscalls. On arm platforms, this uses mmap.
func Write(addr int64, data UintN) error {
	f, err := os.OpenFile(memPath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	mem, offset, err := mmap(f, addr, data.Size(), syscall.PROT_WRITE)
	if err != nil {
		return err
	}
	defer syscall.Munmap(mem)

	// MMIO makes this a bit tricky. Writes must be conducted in one store
	// operation. Review the generated assembly to make sure.
	return data.write(unsafe.Pointer(&mem[offset]))
}
