// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var memPath = "/dev/mem"

var pageSize = int64(syscall.Getpagesize())

type syscalls interface {
	Mmap(int, int64, int, int, int) ([]byte, error)
	Munmap([]byte) error
}

type calls struct{}

func (c *calls) Mmap(fd int, page int64, mapSize int, prot int, callid int) ([]byte, error) {
	return syscall.Mmap(fd, page, mapSize, prot, callid)
}

func (c *calls) Munmap(mem []byte) error {
	return syscall.Munmap(mem)
}

// MMap is a struct containing an os.File and an interface to system calls to manage mapped files.
type MMap struct {
	*os.File
	syscalls
}

// mmap aligns the address and maps multiple pages when needed.
func (m *MMap) mmap(f *os.File, addr int64, size int64, prot int) (mem []byte, offset int64, err error) {
	page := addr &^ (pageSize - 1)
	offset = addr - page
	mapSize := offset + size
	mem, err = m.Mmap(int(f.Fd()), int64(page), int(mapSize), prot, syscall.MAP_SHARED)
	return
}

// ReadAt reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls. On arm platforms, this uses mmap.
func (m *MMap) ReadAt(addr int64, data UintN) error {
	mem, offset, err := m.mmap(m.File, addr, data.Size(), syscall.PROT_READ)
	if err != nil {
		return fmt.Errorf("reading %#x/%d: %w", addr, data.Size(), err)
	}
	defer m.Munmap(mem)

	// MMIO makes this a bit tricky. Reads must be conducted in one load
	// operation. Review the generated assembly to make sure.
	if err := data.read(unsafe.Pointer(&mem[offset])); err != nil {
		return fmt.Errorf("reading %#x/%d: %w", addr, data.Size(), err)
	}
	return nil
}

// WriteAt writes data to physical memory at address addr. On x86 platforms, this
// uses the seek+read syscalls. On arm platforms, this uses mmap.
func (m *MMap) WriteAt(addr int64, data UintN) error {
	mem, offset, err := m.mmap(m.File, addr, data.Size(), syscall.PROT_WRITE)
	if err != nil {
		return err
	}
	defer m.Munmap(mem)

	// MMIO makes this a bit tricky. Writes must be conducted in one store
	// operation. Review the generated assembly to make sure.
	return data.write(unsafe.Pointer(&mem[offset]))
}

// Close implements Close.
func (m *MMap) Close() error {
	return m.File.Close()
}

// NewMMap returns an Mmap for a file (usually a device) passed as a string.
func NewMMap(path string) (*MMap, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &MMap{
		File:     f,
		syscalls: &calls{},
	}, nil
}

// Read is deprecated. Still here for compatibility.
// Use NewMMap() and the interface function instead.
func Read(addr int64, data UintN) error {
	mmap, err := NewMMap(memPath)
	if err != nil {
		return err
	}
	defer mmap.Close()
	return mmap.ReadAt(addr, data)
}

// Write is deprecated. Still here for compatibility.
// Use NewMMap() and the interface function instead.
func Write(addr int64, data UintN) error {
	mmap, err := NewMMap(memPath)
	if err != nil {
		return err
	}
	defer mmap.Close()
	return mmap.WriteAt(addr, data)
}
