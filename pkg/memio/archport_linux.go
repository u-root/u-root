// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)

package memio

import (
	"fmt"
	"sync"
	"syscall"
)

// ArchPort is used for architectural access to a port, instead of file system
// level access. On the x86, this means direct, in-line in[bwl]/out[bwl]
// instructions, requiring an iopl system call. On other architectures,
// it may require special mmap setup.
type ArchPort struct{}

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
func (a *ArchPort) In(addr uint16, data UintN) error {
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

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func (a *ArchPort) Out(addr uint16, data UintN) error {
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

// Close implements close
func (a *ArchPort) Close() error {
	return nil
}

// ArchIn is deprecated. Only here to keep compatibility
func ArchIn(addr uint16, data UintN) error {
	archport := &ArchPort{}
	return archport.In(addr, data)
}

// ArchOut is deprecated. Only here to keep compatibility
func ArchOut(addr uint16, data UintN) error {
	archport := &ArchPort{}
	return archport.Out(addr, data)
}
