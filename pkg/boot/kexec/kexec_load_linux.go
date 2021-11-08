// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/linux"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_INITRD_ADDR_MAX  uint = 0x37FFFFFF
	DEFAULT_BZIMAGE_ADDR_MAX uint = 0x37FFFFFF
	bootParams                    = "/sys/kernel/boot_params/data"
)

const defaultPurgatory = "to32bit_3000"

var (
	// Debug is called to print out verbose debug info.
	//
	// Set this to appropriate output stream for display
	// of useful debug info.
	Debug        = log.Printf // func(string, ...interface{}) {}
	curPurgatory = linux.Purgatories[defaultPurgatory]
)

// Load loads the given segments into memory to be executed on a kexec-reboot.
//
// It is assumed that segments is made up of the next kernel's code and text
// segments, and that `entry` is the entry point, either kernel entry point or trampoline.
//
// Load will align segments to page boundaries and deduplicate overlapping ranges.
func Load(entry uintptr, segments Segments, flags uint64) error {
	segments, err := AlignAndMerge(segments)
	if err != nil {
		return fmt.Errorf("could not align segments: %w", err)
	}

	if !segments.PhysContains(entry) {
		return fmt.Errorf("entry point %#v is not contained by any segment", entry)
	}
	return rawLoad(entry, segments, flags)
}

// ErrKexec is returned by Load if the kexec failed. It describes entry point,
// flags, errno and kernel layout.
type ErrKexec struct {
	Entry    uintptr
	Segments []Segment
	Flags    uint64
	Errno    syscall.Errno
}

// Error implements error.
func (e ErrKexec) Error() string {
	return fmt.Sprintf("kexec_load(entry=%#x, segments=%s, flags %#x) = errno %s", e.Entry, e.Segments, e.Flags, e.Errno)
}

// rawLoad is a wrapper around kexec_load(2) syscall.
// Preconditions:
// - segments must not overlap
// - segments must be full pages
func rawLoad(entry uintptr, segments []Segment, flags uint64) error {
	if _, _, errno := unix.Syscall6(
		unix.SYS_KEXEC_LOAD,
		entry,
		uintptr(len(segments)),
		uintptr(unsafe.Pointer(&segments[0])),
		uintptr(flags),
		0, 0); errno != 0 {
		return ErrKexec{
			Entry:    entry,
			Segments: segments,
			Flags:    flags,
			Errno:    errno,
		}
	}
	return nil
}

// SelectPurgatory picks a purgatory, returning an error if none is found
func SelectPurgator(name string) error {
	p, ok := linux.Purgatories[name]
	if !ok {
		var s []string
		for i := range linux.Purgatories {
			s = append(s, i)
		}
		return fmt.Errorf("%s: no such purgatory, try one of %v", name, s)

	}
	curPurgatory = p
	return nil
}
