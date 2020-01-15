// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"debug/elf"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

var pageMask = uint(os.Getpagesize() - 1)

// ErrNotEnoughSpace is returned by the FindSpace family of functions if no
// range is large enough to accommodate the request.
type ErrNotEnoughSpace struct {
	Size uint
}

func (e ErrNotEnoughSpace) Error() string {
	return fmt.Sprintf("not enough space to allocate %#x bytes", e.Size)
}

// Range represents a contiguous uintptr interval [Start, Start+Size).
type Range struct {
	// Start is the inclusive start of the range.
	Start uintptr

	// Size is the number of elements in the range.
	//
	// Start+Size is the exclusive end of the range.
	Size uint
}

// RangeFromInterval returns a Range representing [start, end).
func RangeFromInterval(start, end uintptr) Range {
	return Range{
		Start: start,
		Size:  uint(end - start),
	}
}

// String returns [Start, Start+Size) as a string.
func (r Range) String() string {
	return fmt.Sprintf("[%#x, %#x)", r.Start, r.End())
}

// End returns the uintptr *after* the end of the interval.
func (r Range) End() uintptr {
	return r.Start + uintptr(r.Size)
}

// Adjacent returns true if r and r2 do not overlap, but are immediately next
// to each other.
func (r Range) Adjacent(r2 Range) bool {
	return r2.End() == r.Start || r.End() == r2.Start
}

// Contains returns true iff p is in the interval described by r.
func (r Range) Contains(p uintptr) bool {
	return r.Start <= p && p < r.End()
}

func min(a, b uintptr) uintptr {
	if a < b {
		return a
	}
	return b
}

func max(a, b uintptr) uintptr {
	if a > b {
		return a
	}
	return b
}

// Intersect returns the continuous range of points common to r and r2 if there
// is one.
func (r Range) Intersect(r2 Range) *Range {
	if !r.Overlaps(r2) {
		return nil
	}
	i := RangeFromInterval(max(r.Start, r2.Start), min(r.End(), r2.End()))
	return &i
}

// Minus removes all points in r2 from r.
func (r Range) Minus(r2 Range) []Range {
	var result []Range
	if r.Contains(r2.Start) && r.Start != r2.Start {
		result = append(result, Range{
			Start: r.Start,
			Size:  uint(r2.Start - r.Start),
		})
	}
	if r.Contains(r2.End()) && r.End() != r2.End() {
		result = append(result, Range{
			Start: r2.End(),
			Size:  uint(r.End() - r2.End()),
		})
	}
	// Neither end was in r?
	//
	// Either r is a subset of r2 and r disappears completely, or they are
	// completely disjunct.
	if len(result) == 0 && r.Disjunct(r2) {
		result = append(result, r)
	}
	return result
}

// Overlaps returns true if r and r2 overlap.
func (r Range) Overlaps(r2 Range) bool {
	return r.Start < r2.End() && r2.Start < r.End()
}

// IsSupersetOf returns true if r2 in r.
func (r Range) IsSupersetOf(r2 Range) bool {
	return r.Start <= r2.Start && r.End() >= r2.End()
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

// Ranges is a list of non-overlapping ranges.
type Ranges []Range

// Minus removes all points in r from all ranges in rs.
func (rs Ranges) Minus(r Range) Ranges {
	var ram Ranges
	for _, oldRange := range rs {
		ram = append(ram, oldRange.Minus(r)...)
	}
	return ram
}

// FindSpace finds a continguous piece of sz points within Ranges and returns
// the Range pointing to it.
func (rs Ranges) FindSpace(sz uint) (space Range, err error) {
	return rs.FindSpaceAbove(sz, 0)
}

// MaxAddr is the highest address in a 64bit address space.
const MaxAddr = ^uintptr(0)

// FindSpaceAbove finds a continguous piece of sz points within Ranges and
// returns a space.Start >= minAddr.
func (rs Ranges) FindSpaceAbove(sz uint, minAddr uintptr) (space Range, err error) {
	return rs.FindSpaceIn(sz, RangeFromInterval(minAddr, MaxAddr))
}

// FindSpaceIn finds a continguous piece of sz points within Ranges and returns
// a Range where space.Start >= limit.Start, with space.End() < limit.End().
func (rs Ranges) FindSpaceIn(sz uint, limit Range) (space Range, err error) {
	for _, r := range rs {
		if overlap := r.Intersect(limit); overlap != nil && overlap.Size >= sz {
			return Range{Start: overlap.Start, Size: sz}, nil
		}
	}
	return Range{}, ErrNotEnoughSpace{Size: sz}
}

// Sort sorts ranges by their start point.
func (rs Ranges) Sort() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Start < rs[j].Start
	})
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
	if buf == nil {
		return Segment{
			Buf: Range{
				Start: 0,
				Size:  0,
			},
			Phys: phys,
		}
	}
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
	return fmt.Sprintf("(virt: %s, phys: %s)", s.Buf, s.Phys)
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

	if s.Buf.Start < diff && diff > 0 {
		panic("cannot have virtual memory address within first page")
	}
	s.Buf.Start -= diff

	if s.Buf.Size > 0 {
		s.Buf.Size += uint(diff)
	}
	return s
}

// Segments is a collection of segments.
type Segments []Segment

// PhysContains returns whether p exists in any of segs' physical memory
// ranges.
func (segs Segments) PhysContains(p uintptr) bool {
	for _, s := range segs {
		if s.Phys.Contains(p) {
			return true
		}
	}
	return false
}

// Insert inserts s assuming it does not overlap with an existing segment.
func (segs *Segments) Insert(s Segment) {
	*segs = append(*segs, s)
	segs.sort()
}

func (segs Segments) sort() {
	sort.Slice(segs, func(i, j int) bool {
		return segs[i].Phys.Start < segs[j].Phys.Start
	})
}

// Dedup deduplicates overlapping and merges adjacent segments in segs.
func Dedup(segs Segments) Segments {
	var s Segments
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

// Memory provides routines to work with physical memory ranges.
type Memory struct {
	// Phys defines the layout of physical memory.
	//
	// Phys is used to tell loaded operating systems what memory is usable
	// as RAM, and what memory is reserved (for ACPI or other reasons).
	Phys MemoryMap

	// Segments are the segments used to load a new operating system.
	//
	// Each segment also contains a physical memory region it maps to.
	Segments Segments
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

		var d []byte
		// Only load segment if there are some data. The kexec call will zero out the rest of the buffer (all of it if Filesz=0):
		// | bufsz bytes are copied from the source buffer to the target kernel buffer. If bufsz is less than memsz, then the excess bytes in the kernel buffer are zeroed out.
		// http://man7.org/linux/man-pages/man2/kexec_load.2.html
		if p.Filesz != 0 {
			d = make([]byte, p.Filesz)
			n, err := r.ReadAt(d, int64(p.Off))
			if err != nil {
				return err
			}
			if n < len(d) {
				return fmt.Errorf("not all data of the segment was read")
			}
		}
		// TODO(hugelgupf): check if this is within availableRAM??
		s := NewSegment(d, Range{
			Start: uintptr(p.Paddr),
			Size:  uint(p.Memsz),
		})
		m.Segments.Insert(s)
	}
	return nil
}

// ParseMemoryMap reads firmware provided memory map from /sys/firmware/memmap.
func (m *Memory) ParseMemoryMap() error {
	p, err := ParseMemoryMap()
	if err != nil {
		return err
	}
	m.Phys = p
	return nil
}

var memoryMapRoot = "/sys/firmware/memmap/"

// ParseMemoryMap reads firmware provided memory map from /sys/firmware/memmap.
func ParseMemoryMap() (MemoryMap, error) {
	return internalParseMemoryMap(memoryMapRoot)
}

func internalParseMemoryMap(memoryMapDir string) (MemoryMap, error) {
	type memRange struct {
		// start and end addresses are inclusive
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
			typ, ok := sysfsToRangeType[data]
			if !ok {
				log.Printf("Sysfs file %q contains unrecognized memory map type %q, defaulting to Reserved", name, data)
				r.typ = RangeReserved
			} else {
				r.typ = typ
			}
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

	if err := filepath.Walk(memoryMapDir, walker); err != nil {
		return nil, err
	}

	var phys []TypedRange
	for _, r := range ranges {
		// Range's end address is exclusive, while Linux's sysfs prints
		// the end address inclusive.
		//
		// E.g. sysfs will contain
		//
		// start: 0x100, end: 0x1ff
		//
		// while we represent
		//
		// start: 0x100, size: 0x100.
		phys = append(phys, TypedRange{
			Range: RangeFromInterval(r.start, r.end+1),
			Type:  r.typ,
		})
	}
	sort.Slice(phys, func(i, j int) bool {
		return phys[i].Start < phys[j].Start
	})
	return phys, nil
}

// M1 is 1 Megabyte in bits.
const M1 = 1 << 20

// FindSpace returns pointer to the physical memory, where array of size sz can
// be stored during next AddKexecSegment call.
func (m Memory) FindSpace(sz uint) (Range, error) {
	// Allocate full pages.
	sz = alignUp(sz)

	// Don't use memory below 1M, just in case.
	return m.AvailableRAM().FindSpaceAbove(sz, M1)
}

// ReservePhys reserves page-aligned sz bytes in the physical memmap within
// the given limit address range.
func (m *Memory) ReservePhys(sz uint, limit Range) (Range, error) {
	sz = alignUp(sz)

	r, err := m.AvailableRAM().FindSpaceIn(sz, limit)
	if err != nil {
		return Range{}, err
	}

	m.Phys.Insert(TypedRange{
		Range: r,
		Type:  RangeReserved,
	})
	return r, nil
}

// AddPhysSegment reserves len(d) bytes in the physical memmap within limit and
// adds a kexec segment with d in that range.
func (m *Memory) AddPhysSegment(d []byte, limit Range) (Range, error) {
	r, err := m.ReservePhys(uint(len(d)), limit)
	if err != nil {
		return Range{}, err
	}
	m.Segments.Insert(NewSegment(d, r))
	return r, nil
}

// AddKexecSegment adds d to a new kexec segment
func (m *Memory) AddKexecSegment(d []byte) (Range, error) {
	r, err := m.FindSpace(uint(len(d)))
	if err != nil {
		return Range{}, err
	}
	m.Segments.Insert(NewSegment(d, r))
	return r, nil
}

// AvailableRAM returns page-aligned unused regions of RAM.
//
// AvailableRAM takes all RAM-marked pages in the memory map and subtracts the
// kexec segments already allocated. RAM segments begin at a page boundary.
//
// E.g if page size is 4K and RAM segments are
//            [{start:0 size:8192} {start:8192 size:8000}]
// and kexec segments are
//            [{start:40 size:50} {start:8000 size:2000}]
// result should be
//            [{start:0 size:40} {start:4096 end:8000 - 4096}]
func (m Memory) AvailableRAM() Ranges {
	ram := m.Phys.FilterByType(RangeRAM)

	// Remove all points in Segments from available RAM.
	for _, s := range m.Segments {
		ram = ram.Minus(s.Phys)
	}

	// Only return Ranges starting at an aligned size.
	var alignedRanges Ranges
	for _, r := range ram {
		alignedStart := alignUpPtr(r.Start)
		if alignedStart < r.End() {
			alignedRanges = append(alignedRanges, Range{
				Start: alignedStart,
				Size:  r.Size - uint(alignedStart-r.Start),
			})
		}
	}
	return alignedRanges
}

// RangeType defines type of a TypedRange based on the Linux
// kernel string provided by firmware memory map.
type RangeType string

// These are the range types we know Linux uses.
const (
	RangeRAM      RangeType = "System RAM"
	RangeDefault  RangeType = "Default"
	RangeACPI     RangeType = "ACPI Tables"
	RangeNVS      RangeType = "ACPI Non-volatile Storage"
	RangeReserved RangeType = "Reserved"
)

// String implements fmt.Stringer.
func (r RangeType) String() string {
	return string(r)
}

var sysfsToRangeType = map[string]RangeType{
	"System RAM":                RangeRAM,
	"Default":                   RangeDefault,
	"ACPI Tables":               RangeACPI,
	"ACPI Non-volatile Storage": RangeNVS,
	"Reserved":                  RangeReserved,
	"reserved":                  RangeReserved,
}

// TypedRange represents range of physical memory.
type TypedRange struct {
	Range
	Type RangeType
}

func (tr TypedRange) String() string {
	return fmt.Sprintf("{addr: %s, type: %s}", tr.Range, tr.Type)
}

// MemoryMap defines the layout of physical memory.
//
// MemoryMap defines which ranges in memory are usable RAM and which are
// reserved for various reasons.
type MemoryMap []TypedRange

// FilterByType only returns ranges of the given typ.
func (m MemoryMap) FilterByType(typ RangeType) Ranges {
	var rs Ranges
	for _, tr := range m {
		if tr.Type == typ {
			rs = append(rs, tr.Range)
		}
	}
	return rs
}

func (m MemoryMap) sort() {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Start < m[j].Start
	})
}

// Insert a new TypedRange into the memory map, removing chunks of other ranges
// as necessary.
//
// Assumes that TypedRange is a valid range -- no checking.
func (m *MemoryMap) Insert(r TypedRange) {
	var newMap MemoryMap

	// Remove points in r from all existing physical ranges.
	for _, q := range *m {
		split := q.Range.Minus(r.Range)
		for _, r2 := range split {
			newMap = append(newMap, TypedRange{Range: r2, Type: q.Type})
		}
	}

	newMap = append(newMap, r)
	newMap.sort()
	*m = newMap
}
