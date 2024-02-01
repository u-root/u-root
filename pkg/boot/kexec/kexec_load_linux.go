// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
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

// kexecSegment defines kernel memory layout.
type kexecSegment struct {
	// Buf points to a buffer in user space.
	Buf Range

	// Phys is a physical address of kernel.
	Phys Range
}

func (s Segment) toKexecSegment() kexecSegment {
	if s.Buf == nil {
		return kexecSegment{
			Buf:  Range{Start: 0, Size: 0},
			Phys: s.Phys,
		}
	}
	return kexecSegment{
		Buf: Range{
			Start: uintptr((unsafe.Pointer(&s.Buf[0]))),
			Size:  uint(len(s.Buf)),
		},
		Phys: s.Phys,
	}
}

func (segs Segments) toKexecSegments() []kexecSegment {
	var ks []kexecSegment
	for _, seg := range segs {
		ks = append(ks, seg.toKexecSegment())
	}
	return ks
}

// rawLoad is a wrapper around kexec_load(2) syscall.
// Preconditions:
// - segments must not overlap
// - segments must be full pages
func rawLoad(entry uintptr, segments Segments, flags uint64) error {
	ks := segments.toKexecSegments()
	_, _, errno := unix.Syscall6(
		unix.SYS_KEXEC_LOAD,
		entry,
		uintptr(len(segments)),
		uintptr(unsafe.Pointer(&ks[0])),
		uintptr(flags),
		0, 0)
	// segments (and all the buffers therein) may have gotten freed after
	// evaluating Syscall6 arguments, but before the syscall actually
	// happens.
	for _, seg := range segments {
		runtime.KeepAlive(seg.Buf)
	}
	runtime.KeepAlive(segments)
	if errno != 0 {
		return ErrKexec{
			Entry:    entry,
			Segments: segments,
			Flags:    flags,
			Errno:    errno,
		}
	}
	return nil
}
