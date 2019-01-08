// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,arm linux,arm64

package io

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	pageSize = int64(syscall.Getpagesize())
	memPath  = "/dev/mem"
)

// mmap aligns the address and maps multiple pages when needed.
func mmap(f *os.File, addr int64, size int64, prot int) (mem []byte, offset int64, err error) {
	if addr+size <= addr {
		return nil, 0, fmt.Errorf("invalid address for size %#x", size)
	}
	page := addr &^ (pageSize - 1)
	offset = addr - page
	mapSize := offset + size
	mem, err = syscall.Mmap(int(f.Fd()), int64(page),
		int(mapSize), prot, syscall.MAP_SHARED)
	return
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

	f, err := os.OpenFile(memPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	size := reflect.Indirect(reflect.ValueOf(data)).Type().Size()
	mem, offset, err := mmap(f, addr, int64(size), syscall.PROT_READ)
	if err != nil {
		return err
	}
	defer syscall.Munmap(mem)

	// MMIO makes this a bit tricky. Reads must be conducted in one load
	// operation. Review the generated assembly to make sure.
	p := unsafe.Pointer(&mem[offset])
	switch data := data.(type) {
	case *uint8:
		*data = *(*uint8)(p)
	case *uint16:
		*data = *(*uint16)(p)
	case *uint32:
		*data = *(*uint32)(p)
	case *uint64:
		// Warning: On arm, this uses two ldr's rather than ldrd.
		*data = *(*uint64)(p)
	}
	return nil
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

	f, err := os.OpenFile(memPath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	size := reflect.TypeOf(data).Size()
	mem, offset, err := mmap(f, addr, int64(size), syscall.PROT_WRITE)
	if err != nil {
		return err
	}
	defer syscall.Munmap(mem)

	// MMIO makes this a bit tricky. Writes must be conducted in one store
	// operation. Review the generated assembly to make sure.
	p := unsafe.Pointer(&mem[offset])
	switch data := data.(type) {
	case uint8:
		*(*uint8)(p) = data
	case uint16:
		*(*uint16)(p) = data
	case uint32:
		*(*uint32)(p) = data
	case uint64:
		// Warning: On arm, this uses two str's rather than strd.
		*(*uint64)(p) = data
	}
	return nil
}
