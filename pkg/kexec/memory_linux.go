// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"debug/elf"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Range represents a contiguous uintptr interval [Start, Start+Size).
type Range struct {
	// Start is the inclusive start of the range.
	Start uintptr
	// Size is the size of the range.
	// Start+Size is the exclusive end of the range.
	Size uint
}

// Overlaps returns true if r and r2 overlap.
func (r Range) Overlaps(r2 Range) bool {
	return r.Start < (r2.Start+uintptr(r2.Size)) && r2.Start < (r.Start+uintptr(r.Size))
}

// IsSupersetOf returns true if r2 in r.
func (r Range) IsSupersetOf(r2 Range) bool {
	return r.Start <= r2.Start && (r.Start+uintptr(r.Size)) >= (r2.Start+uintptr(r2.Size))
}

// Disjunct returns true if r and r2 does not overlap.
func (r Range) Disjunct(r2 Range) bool {
	return !r.Overlaps(r2)
}

func (r Range) toSlice() []byte {
	var data []byte

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = r.Start
	sh.Len = int(r.Size)
	sh.Cap = int(r.Size)

	return data
}

// pool stores byte slices pointed by the pointers Segments.Buf to
// prevent underlying arrays to be collected by garbage collector.
var pool [][]byte

// Segment defines kernel memory layout.
type Segment struct {
	// Buf is a buffer in user space.
	Buf Range
	// Phys is a physical address of kernel.
	Phys Range
}

// NewSegment creates new Segment.
// Segments should be created using NewSegment method to prevent
// data pointed by Segment.Buf to be collected by garbage collector.
func NewSegment(buf []byte, phys Range) Segment {
	pool = append(pool, buf)
	return Segment{
		Buf: Range{
			Start: uintptr((unsafe.Pointer(&buf[0]))),
			Size:  uint(len(buf)),
		},
		Phys: phys,
	}
}

func (s Segment) String() string {
	return fmt.Sprintf("(virt: %#x + %#x | phys: %#x + %#x)", s.Buf.Start, s.Buf.Size, s.Phys.Start, s.Phys.Size)
}

func ptrToSlice(ptr uintptr, size int) []byte {
	var data []byte

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = ptr
	sh.Len = size
	sh.Cap = size

	return data
}

func (s *Segment) tryMerge(s2 Segment) (ok bool) {
	if s.Phys.Disjunct(s2.Phys) {
		return false
	}

	// Virtual memory ranges should never overlap,
	// concatenate ranges.
	a := s.Buf.toSlice()
	b := s2.Buf.toSlice()
	c := append(a, b...)

	phys := s.Phys
	// s1 and s2 overlap somewhat.
	if !s.Phys.IsSupersetOf(s2.Phys) {
		phys.Size = uint(s2.Phys.Start-s.Phys.Start) + s2.Phys.Size
	}

	*s = NewSegment(c, phys)
	return true
}

// AlignPhys fixes s to the kexec_load preconditions.
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
func AlignPhys(s Segment) Segment {
	pageMask := uint(os.Getpagesize() - 1)
	orig := s.Phys.Start
	// Find the page address of the starting point.
	s.Phys.Start = s.Phys.Start &^ uintptr(pageMask)

	diff := orig - s.Phys.Start
	// Round up to page size.
	s.Phys.Size = (s.Phys.Size + uint(diff) + pageMask) &^ pageMask

	if s.Buf.Start < diff {
		panic("cannot have virtual memory address within first page")
	}
	s.Buf.Start -= diff

	if s.Buf.Size > 0 {
		s.Buf.Size += uint(diff)
	}
	return s
}

// Dedup merges segments in segs as much as possible.
func Dedup(segs []Segment) []Segment {
	var s []Segment
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].Phys.Start == segs[j].Phys.Start {
			// let segs[i] be the superset of segs[j]
			return segs[i].Phys.Size > segs[j].Phys.Size
		}
		return segs[i].Phys.Start < segs[j].Phys.Start
	})

	for _, seg := range segs {
		doIt := true
		for i := range s {
			if merged := s[i].tryMerge(seg); merged {
				doIt = false
				break
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
// segments, and that `entry` is the entry point, either kernel entry point or trampoline.
func Load(entry uintptr, segments []Segment, flags uint64) error {
	for i := range segments {
		segments[i] = AlignPhys(segments[i])
	}

	segments = Dedup(segments)
	ok := false
	for _, s := range segments {
		ok = ok || (s.Phys.Start <= entry && entry < s.Phys.Start+uintptr(s.Phys.Size))
	}
	if !ok {
		return fmt.Errorf("entry point %#v is not covered by any segment", entry)
	}

	return rawLoad(entry, segments, flags)
}

// ErrKexec is the error type returned kexec.
// It describes entry point, flags, errno and kernel layout.
type ErrKexec struct {
	Entry    uintptr
	Segments []Segment
	Flags    uint64
	Errno    syscall.Errno
}

func (e ErrKexec) Error() string {
	return fmt.Sprintf("entry %x, flags %x, errno %d, segments %v", e.Entry, e.Flags, e.Errno, e.Segments)
}

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

// LoadElfSegments loads loadable ELF segments.
func (m *Memory) LoadElfSegments(r io.ReaderAt) error {
	f, err := elf.NewFile(r)
	if err != nil {
		return err
	}

	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}
		d := make([]byte, p.Filesz)
		n, err := r.ReadAt(d, int64(p.Off))
		if err != nil {
			return err
		}
		if n < len(d) {
			return fmt.Errorf("not all data of the segment was read")
		}
		s := NewSegment(d, Range{
			Start: uintptr(p.Paddr),
			Size:  uint(p.Memsz),
		})

		m.Segments = append(m.Segments, s)
	}
	return nil
}

var memoryMapRoot = "/sys/firmware/memmap/"

// ParseMemoryMap reads firmware provided memory map
// from /sys/firmware/memmap.
func (m *Memory) ParseMemoryMap() error {
	type memRange struct {
		// start and addresses are inclusive
		start, end uintptr
		typ        RangeType
	}

	ranges := make(map[string]memRange)
	walker := func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		const (
			// file names
			start = "start"
			end   = "end"
			typ   = "type"
		)

		base := path.Base(name)
		if base != start && base != end && base != typ {
			return fmt.Errorf("unexpected file %q", name)
		}
		dir := path.Dir(name)

		b, err := ioutil.ReadFile(name)
		if err != nil {
			return fmt.Errorf("error reading file %q: %v", name, err)
		}

		data := strings.TrimSpace(string(b))
		r := ranges[dir]
		if base == typ {
			r.typ = RangeType(data)
			ranges[dir] = r
			return nil
		}

		v, err := strconv.ParseUint(data, 0, 64)
		if err != nil {
			return err
		}
		switch base {
		case start:
			r.start = uintptr(v)
		case end:
			r.end = uintptr(v)
		}
		ranges[dir] = r
		return nil
	}

	if err := filepath.Walk(memoryMapRoot, walker); err != nil {
		return err
	}

	for _, r := range ranges {
		m.Phys = append(m.Phys, TypedAddressRange{
			Range: Range{
				Start: r.start,
				Size:  uint(r.end - r.start),
			},
			Type: r.typ,
		})
	}
	sort.Slice(m.Phys, func(i, j int) bool {
		return m.Phys[i].Start < m.Phys[j].Start
	})
	return nil
}

// RangeType defines type of a TypedAddressRange based on the Linux
// kernel string provided by firmware memory map.
type RangeType string

const (
	RangeRAM     RangeType = "System RAM"
	RangeDefault           = "Default"
	RangeNVACPI            = "ACPI Non-volatile Storage"
	RangeACPI              = "ACPI Tables"
	RangeNVS               = "Reserved"
)

// Memory provides routines to work with physical memory ranges.
type Memory struct {
	Phys []TypedAddressRange

	Segments []Segment
}

// TypedAddressRange represents range of physical memory.
type TypedAddressRange struct {
	Range
	Type RangeType
}

var ErrNotEnoughSpace = errors.New("not enough space")

// FindSpace returns pointer to the physical memory,
// where array of size sz can be stored during next
// AddKexecSegment call.
func (m Memory) FindSpace(sz uint) (start uintptr, err error) {
	pageSize := uint(os.Getpagesize())
	sz = (sz + pageSize - 1) &^ (pageSize - 1)
	ranges := m.availableRAM()
	for _, r := range ranges {
		// don't use memory below 1M, just in case.
		if uint(r.Start)+r.Size < 1048576 {
			continue
		}
		if r.Size >= sz {
			return r.Start, nil
		}
	}
	return 0, ErrNotEnoughSpace
}

func (m *Memory) addKexecSegment(addr uintptr, d []byte) {
	s := NewSegment(d, Range{
		Start: addr,
		Size:  uint(len(d)),
	})
	s = AlignPhys(s)
	m.Segments = append(m.Segments, s)
	sort.Slice(m.Segments, func(i, j int) bool {
		return m.Segments[i].Phys.Start < m.Segments[j].Phys.Start
	})
}

// AddKexecSegment adds d to a new kexec segment
func (m *Memory) AddKexecSegment(d []byte) (addr uintptr, err error) {
	size := uint(len(d))
	start, err := m.FindSpace(size)
	if err != nil {
		return 0, err
	}
	m.addKexecSegment(start, d)
	return start, nil
}

// availableRAM subtracts physycal ranges of kexec segments from
// RAM segments of TypedAddressRange.
//
// E.g if RAM segments are
//            [{start:0 size:100} {start:200 size:100}]
// and kexec segments are
//            [{start:0 size:50} {start:100 size:100} {start:250 size:50}]
// result should be
//            [{start:50 size:50} {start:200 end:250}]
func (m Memory) availableRAM() []TypedAddressRange {
	var ret []TypedAddressRange
	var i int
	for _, a := range m.Phys {
		if a.Type != RangeRAM {
			continue
		}

		start := a.Start
		end := a.Start + uintptr(a.Size)
		for ; i < len(m.Segments); i++ {
			b := m.Segments[i]
			if end < b.Phys.Start {
				break
			}
			if start < b.Phys.Start {
				ret = append(ret, TypedAddressRange{
					Range: Range{
						Start: start,
						Size:  uint(b.Phys.Start - start),
					},
					Type: a.Type,
				})
			}
			start = b.Phys.Start + uintptr(b.Phys.Size)
		}

		if start >= end {
			continue
		}

		ret = append(ret, TypedAddressRange{
			Range: Range{
				Start: start,
				Size:  uint(end - start),
			},
			Type: a.Type,
		})
	}
	return ret
}
