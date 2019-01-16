// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
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
	"unsafe"

	"github.com/u-root/u-root/pkg/dt"
)

var pageMask = uint(os.Getpagesize() - 1)

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

// Disjunct returns true if r and r2 do not overlap.
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

func alignUp(p uint) uint {
	return (p + pageMask) &^ pageMask
}

func alignUpPtr(p uintptr) uintptr {
	return uintptr(alignUp(uint(p)))
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
	orig := s.Phys.Start
	// Find the page address of the starting point.
	s.Phys.Start = s.Phys.Start &^ uintptr(pageMask)

	diff := orig - s.Phys.Start
	// Round up to page size.
	s.Phys.Size = alignUp(s.Phys.Size + uint(diff))

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

var memmapRoot = "/sys/firmware/memmap/"

// ParseFromMemmap reads the firmware provided memory map from
// /sys/firmware/memmap.
func (m *Memory) ParseFromMemmap() error {
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

	if err := filepath.Walk(memmapRoot, walker); err != nil {
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

// ParseFromDeviceTree reads the memory map from the device tree.
func (m *Memory) ParseFromDeviceTree(fdt *dt.FDT) error {
	// TODO: Does not support a number of this: regions, multiple memory nodes,
	//       #address-cells, #size-cells, ...
	memoryNode := fdt.RootNode.GetNode("memory")
	if memoryNode == nil {
		return errors.New("device tree missing /memory node")
	}
	regProp := memoryNode.GetProperty("reg")
	if memoryNode == nil {
		return errors.New("device tree missing /memory/reg property")
	}
	array, _ := regProp.AsPropEncodedArray()
	if len(array)%8 != 0 {
		return errors.New("device tree /memory/reg property not a multiple of two cells")
	}

	// The reg property is an array of (address, value) pair cells.
	b := bytes.NewBuffer(array)
	for b.Len() > 0 {
		var address, size uint32
		binary.Read(b, binary.BigEndian, &address)
		binary.Read(b, binary.BigEndian, &size)

		m.Phys = append(m.Phys, TypedAddressRange{
			Range: Range{
				Start: uintptr(address),
				Size:  uint(size),
			},
			Type: RangeRAM,
		})
	}
	return nil
}

// RangeType defines type of a TypedAddressRange based on the Linux
// kernel string provided by firmware memory map.
type RangeType string

const (
	RangeRAM      RangeType = "System RAM"
	RangeDefault            = "Default"
	RangeACPI               = "ACPI Tables"
	RangeNVS                = "ACPI Non-volatile Storage"
	RangeReserved           = "Reserved"
)

// Memory provides routines to work with physical memory ranges.
type Memory struct {
	Phys     []TypedAddressRange
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
	sz = alignUp(sz)
	ranges := m.availableRAM()
	for _, r := range ranges {
		// don't use memory below 1M, just in case.
		// TODO: not true on non-x86 platforms
		if uint(r.Start)+r.Size < 1048576 {
			continue
		}
		if r.Size >= sz {
			return r.Start, nil
		}
	}
	return 0, ErrNotEnoughSpace
}

// ContainedInAvailableMemory return true if the given range appears in the
// available memory.
func (m Memory) ContainedInAvailableMemory(r Range) bool {
	for _, r2 := range m.availableRAM() {
		if r2.Range.IsSupersetOf(r) {
			return true
		}
	}
	return false
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

// AddKexecSegment adds d to a new kexec segment.
func (m *Memory) AddKexecSegment(d []byte) (addr uintptr, err error) {
	size := uint(len(d))
	start, err := m.FindSpace(size)
	if err != nil {
		return 0, err
	}
	m.addKexecSegment(start, d)
	return start, nil
}

// AddKexecSegmentPhys adds a kexec segment at the physical address.
func (m *Memory) AddKexecSegmentPhys(d []byte, phys uintptr) error {
	size := uint(len(d))
	if !m.ContainedInAvailableMemory(Range{Start: phys, Size: size}) {
		return fmt.Errorf(
			"trying to add segment (size=%#x, phys=%#x) outside of availabe memory",
			size, phys)
	}
	m.addKexecSegment(phys, d)
	return nil
}

// availableRAM subtracts physical ranges of kexec segments from
// RAM segments of TypedAddressRange aligning range beginnings
// to a page boundary.
//
// E.g if page size is 4K and RAM segments are
//            [{start:0 size:8192} {start:8192 size:8000}]
// and kexec segments are
//            [{start:40 size:50} {start:8000 size:2000}]
// result should be
//            [{start:0 size:40} {start:4096 end:8000 - 4096}]
func (m Memory) availableRAM() (avail []TypedAddressRange) {
	type point struct {
		// x is a point coordinate on an axis.
		x uintptr
		// start is true if the point is the beginning of segment.
		start bool
		// ram is true if the point is part of a RAM segment.
		ram bool
	}
	// points stores starting and ending points of segments
	// sorted by coordinate.
	var points []point
	addPoint := func(r Range, ram bool) {
		points = append(points,
			point{x: r.Start, start: true, ram: ram},
			point{x: r.Start + uintptr(r.Size) - 1, start: false, ram: ram},
		)
	}

	for _, s := range m.Phys {
		if s.Type == RangeRAM {
			addPoint(s.Range, true)
		}
	}
	for _, s := range m.Segments {
		addPoint(s.Phys, false)
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].x < points[j].x
	})

	add := func(start, end uintptr, ramRange, kexecRange bool) {
		if !ramRange || kexecRange {
			return
		}
		start = alignUpPtr(start)
		if start >= end {
			return
		}
		avail = append(avail, TypedAddressRange{
			Range: Range{
				Start: start,
				Size:  uint(end-start) + 1,
			},
			Type: RangeRAM,
		})
	}

	var start uintptr
	var ramRange bool
	var kexecRange bool
	for _, p := range points {
		switch {
		case p.start && p.ram:
			start = p.x
		case p.start && !p.ram:
			if start != p.x {
				add(start, p.x-1, ramRange, kexecRange)
			}
		case !p.start && p.ram:
			add(start, p.x, ramRange, kexecRange)
		case !p.start && !p.ram:
			if ramRange {
				start = p.x + 1
			}
		}

		if p.ram {
			ramRange = p.start
		} else {
			kexecRange = p.start
		}
	}

	return avail
}
