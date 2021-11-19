// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)
// +build linux,amd64 linux,386

package memio

import (
	"fmt"
	"sync"
	"syscall"
)

const portPath = "/dev/port"

var ioplError struct {
	sync.Once
	err error
}

func iopl() error {
	ioplError.Do(func() {
		ioplError.err = syscall.Iopl(3)
	})
	return ioplError.err
}

// In reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func In(addr uint16, data UintN) error {
	if _, ok := data.(*Uint8); !ok {
		return fmt.Errorf("/dev/port data must be 8 bits on Linux")
	}
	return pathRead(portPath, int64(addr), data)
}

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func Out(addr uint16, data UintN) error {
	if _, ok := data.(*Uint8); !ok {
		return fmt.Errorf("/dev/port data must be 8 bits on Linux")
	}
	return pathWrite(portPath, int64(addr), data)
}

func archInl(uint16) uint32
func archInw(uint16) uint16
func archInb(uint16) uint8

// ArchIn reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func ArchIn(addr uint16, data UintN) error {
	if err := iopl(); err != nil {
		return err
	}

	switch p := data.(type) {
	case *Uint32:
		*p = Uint32(archInl(addr))
	case *Uint16:
		*p = Uint16(archInw(addr))
	case *Uint8:
		*p = Uint8(archInb(addr))
	default:
		return fmt.Errorf("port data must be 8, 16 or 32 bits")
	}
	return nil
}

func archOutl(uint16, uint32)
func archOutw(uint16, uint16)
func archOutb(uint16, uint8)

// ArchOut writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func ArchOut(addr uint16, data UintN) error {
	if err := iopl(); err != nil {
		return err
	}

	switch p := data.(type) {
	case *Uint32:
		archOutl(addr, uint32(*p))
	case *Uint16:
		archOutw(addr, uint16(*p))
	case *Uint8:
		archOutb(addr, uint8(*p))
	default:
		return fmt.Errorf("port data must be 8, 16 or 32 bits")
	}
	return nil
}
