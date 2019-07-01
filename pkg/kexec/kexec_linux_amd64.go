// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// FileLoad loads the given kernel as the new kernel with the given ramfs and
// cmdline.
//
// The kexec_file_load(2) syscall is x86-64 bit only.
func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	var flags int
	var ramfsfd int
	if ramfs != nil {
		ramfsfd = int(ramfs.Fd())
	} else {
		flags |= unix.KEXEC_FILE_NO_INITRAMFS
	}

	if err := unix.KexecFileLoad(int(kernel.Fd()), ramfsfd, cmdline, flags); err != nil {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = %v", kernel.Fd(), ramfsfd, cmdline, flags, err)
	}
	return nil
}

type Range struct {
	Start uintptr
	Size  uint
}

func (r Range) Overlaps(r2 Range) bool {
	return r.Start < (r2.Start+uintptr(r2.Size)) && r2.Start < (r.Start+uintptr(r.Size))
}

// True if r2 \in r.
func (r Range) IsSupersetOf(r2 Range) bool {
	return r.Start <= r2.Start && (r.Start+uintptr(r.Size)) >= (r2.Start+uintptr(r2.Size))
}

func (r Range) Disjunct(r2 Range) bool {
	return !r.Overlaps(r2)
}

type Segment struct {
	Buf  Range
	Phys Range
}

var Delete = errors.New("delete that shit")

func (s *Segment) TryMerge(s2 Segment) error {
	// The world is fine.
	if s.Phys.Disjunct(s2.Phys) {
		return nil
	}

	// s can swallow s2 completely.
	if s.Phys.IsSupersetOf(s2.Phys) {
		return Delete
	}

	// s and s2 overlap somewhat.
	s.Phys.Size = uint(s2.Phys.Start-s.Phys.Start) + s2.Phys.Size
	s.Buf.Size = uint(s2.Buf.Start-s.Buf.Start) + s2.Buf.Size
	return Delete
}

const PAGE_MASK = 4095

// Adjust fixes s to the kexec_load preconditions.
//
// s's physical addresses must be multiples of the page size.
//
// E.g. if page size is 0x1000:
// Segment {
//   Buf:  {Start: 0x1011, Size: 0x1022}
//   Phys: {Start: 0x2011, Size: 0x1022}
// }
// has to become
// Segment {
//   Buf:  {Start: 0x1000, Size: 0x1033}
//   Phys: {Start: 0x2000, Size: 0x2000}
// }
func Adjust(s Segment) Segment {
	orig := s.Phys.Start
	// Find the page address of the starting point.
	s.Phys.Start = s.Phys.Start &^ PAGE_MASK

	diff := orig - s.Phys.Start
	// Round up to page size.
	s.Phys.Size = (s.Phys.Size + uint(diff) + PAGE_MASK) &^ PAGE_MASK

	if s.Buf.Size > 0 {
		s.Buf.Start -= diff
		s.Buf.Size += uint(diff)
	}
	return s
}

// Dedup merges segments in segs as much as possible.
func Dedup(segs []Segment) []Segment {
	var s []Segment
	sort.Slice(segs, func(i, j int) bool { return segs[i].Phys.Start < segs[j].Phys.Start })

	for _, seg := range segs {
		doIt := true
		for i := range s {
			if err := s[i].TryMerge(seg); err == Delete {
				doIt = false
			}
		}
		if doIt {
			s = append(s, seg)
		}
	}
	return s
}

// Load loads the given segments into memory to be executed on a kexec-reboot.
//
// It is assumed that segments is made up of the next kernel's code and text
// segments, and that `entry` is the next kernel's entry point.
//
// As it is, assumes that the next kernel has a 64bit entry point (no
// trampoline).
func Load(entry uintptr, segments []Segment, flags uint64) error {
	for i, s := range segments {
		segments[i] = Adjust(segments[i])
		s = segments[i]
		log.Printf("virt: %#x + %#x | phys: %#x + %#x", s.Buf.Start, s.Buf.Size, s.Phys.Start, s.Phys.Size)
	}

	log.Printf("between adjust and dedup")

	segments = Dedup(segments)
	for _, s := range segments {
		log.Printf("virt: %#x + %#x | phys: %#x + %#x", s.Buf.Start, s.Buf.Size, s.Phys.Start, s.Phys.Size)
	}
	log.Printf("entry point: %#v", entry)
	return rawLoad(entry, segments, flags)
}

type ErrKexec struct {
	Entry    uintptr
	Segments []Segment
	Flags    uint64
	Errno    syscall.Errno
}

func (e ErrKexec) Error() string {
    return "<ErrKexec>"
}

//
//
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
