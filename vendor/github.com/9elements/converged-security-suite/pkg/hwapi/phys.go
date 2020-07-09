// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !arm,!arm64

package hwapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

var memPaths = [...]string{"/dev/fmem", "/dev/mem"}

// UintN is a wrapper around uint types and provides a few io-related
// functions.
type UintN interface {
	// Return size in bytes.
	Size() int64

	// Return string formatted in hex.
	String() string

	// Read from given address with native endianess.
	read(addr unsafe.Pointer) error

	// Write to given address with native endianess.
	write(addr unsafe.Pointer) error
}

// Uint8 is a wrapper around uint8.
type Uint8 uint8

// Uint16 is a wrapper around uint16.
type Uint16 uint16

// Uint32 is a wrapper around uint32.
type Uint32 uint32

// Uint64 is a wrapper around uint64.
type Uint64 uint64

// Size of uint8 is 1.
func (u *Uint8) Size() int64 {
	return 1
}

// Size of uint16 is 2.
func (u *Uint16) Size() int64 {
	return 2
}

// Size of uint32 is 4.
func (u *Uint32) Size() int64 {
	return 4
}

// Size of uint64 is 8.
func (u *Uint64) Size() int64 {
	return 8
}

// String formats a uint8 in hex.
func (u *Uint8) String() string {
	return fmt.Sprintf("%#02x", *u)
}

// String formats a uint16 in hex.
func (u *Uint16) String() string {
	return fmt.Sprintf("%#04x", *u)
}

// String formats a uint32 in hex.
func (u *Uint32) String() string {
	return fmt.Sprintf("%#08x", *u)
}

// String formats a uint64 in hex.
func (u *Uint64) String() string {
	return fmt.Sprintf("%#016x", *u)
}

func (u *Uint8) read(addr unsafe.Pointer) error {
	*u = Uint8(*(*uint8)(addr)) // TODO: rewrite in Go assembly for ARM
	return nil                  // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint16) read(addr unsafe.Pointer) error {
	*u = Uint16(*(*uint16)(addr)) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint32) read(addr unsafe.Pointer) error {
	*u = Uint32(*(*uint32)(addr)) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint64) read(addr unsafe.Pointer) error {
	// Warning: On arm, this uses two ldr's rather than ldrd.
	*u = Uint64(*(*uint64)(addr)) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint8) write(addr unsafe.Pointer) error {
	*(*uint8)(addr) = uint8(*u) // TODO: rewrite in Go assembly for ARM
	return nil                  // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint16) write(addr unsafe.Pointer) error {
	*(*uint16)(addr) = uint16(*u) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint32) write(addr unsafe.Pointer) error {
	*(*uint32)(addr) = uint32(*u) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func (u *Uint64) write(addr unsafe.Pointer) error {
	// Warning: On arm, this uses two str's rather than strd.
	*(*uint64)(addr) = uint64(*u) // TODO: rewrite in Go assembly for ARM
	return nil                    // TODO: catch misalign, segfault, sigbus, ...
}

func pathRead(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekCurrent); err != nil {
		return err
	}
	return binary.Read(f, binary.LittleEndian, data)
}

func selectDevMem() (string, error) {
	if len(memPaths) == 0 {
		return "", fmt.Errorf("Internal error: no /dev/mem device specified")
	}

	for _, p := range memPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("No suitable /dev/mem device found. Tried %#v", memPaths)
}

// ReadPhys reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls.
func (t TxtAPI) ReadPhys(addr int64, data UintN) error {
	devMem, err := selectDevMem()
	if err != nil {
		return err
	}

	return pathRead(devMem, addr, data)
}

// ReadPhysBuf reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls.
func (t TxtAPI) ReadPhysBuf(addr int64, buf []byte) error {
	devMem, err := selectDevMem()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(devMem, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekCurrent); err != nil {
		return err
	}
	return binary.Read(f, binary.LittleEndian, buf)
}

func pathWrite(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekCurrent); err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, data)
}

// WritePhys writes data to physical memory at address addr. On x86 platforms, this
// uses the seek+read syscalls.
func (t TxtAPI) WritePhys(addr int64, data UintN) error {
	devMem, err := selectDevMem()
	if err != nil {
		return err
	}

	return pathWrite(devMem, addr, data)
}
