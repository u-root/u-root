// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for io for Linux.
package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

func in(f IoIntf, addr uint64, data interface{}) error {
	ps := uint64(syscall.Getpagesize())
	page := (addr & ^(ps - 1))
	offset := addr - page
	mem, err := syscall.Mmap(int(f.Fd()), int64(page), int(ps), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("in: mmap failed for address 0x%x: %v", addr, err)
	}

	p := unsafe.Pointer(&mem[offset])

	switch data := data.(type) {
	case *byte:
		*data = *(*byte)(p)
	case *uint16:
		*data = *(*uint16)(p)
	case *uint32:
		*data = *(*uint32)(p)
	case *uint64:
		*data = *(*uint64)(p)
	default:
		return fmt.Errorf("in: internal error, got unsupported type")
	}
	err = syscall.Munmap(mem)
	if err != nil {
		return fmt.Errorf("in: failed to munmap: %v", err)
	}
	return nil
}

func out(f IoIntf, addr uint64, data interface{}) error {
	ps := uint64(syscall.Getpagesize())
	page := (addr & ^(ps - 1))
	offset := addr - page
	mem, err := syscall.Mmap(int(f.Fd()), int64(page), int(ps), syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("out: mmap failed for address 0x%x: %v", addr, err)
	}

	p := unsafe.Pointer(&mem[offset])

	switch data := data.(type) {
	case *byte:
		*(*byte)(p) = *data
	case *uint16:
		*(*uint16)(p) = *data
	case *uint32:
		*(*uint32)(p) = *data
	case *uint64:
		*(*uint64)(p) = *data
	default:
		return fmt.Errorf("out: internal error, got unsupported type")
	}
	err = syscall.Munmap(mem)
	if err != nil {
		return fmt.Errorf("out: failed to munmap: %v", err)
	}
	return nil
}

func inp(f IoIntf, addr uint64, data interface{}) error {
	_, err := f.Seek(int64(addr), 0)
	if err != nil {
		return fmt.Errorf("in: bad address %v: %v", addr, err)
	}
	return binary.Read(f, binary.LittleEndian, data)
}

func outp(f IoIntf, addr uint64, data interface{}) error {
	_, err := f.Seek(int64(addr), 0)
	if err != nil {
		return fmt.Errorf("out: bad address %v: %v", addr, err)
	}
	return binary.Write(f, binary.LittleEndian, data)
}
