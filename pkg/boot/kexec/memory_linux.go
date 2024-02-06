// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/u-root/u-root/pkg/align"
)

var pageMask = uint(os.Getpagesize() - 1)

// ErrNotEnoughSpace is returned by the FindSpace family of functions if no
// range is large enough to accommodate the request.
var ErrNotEnoughSpace = fmt.Errorf("not enough space to allocate bytes")

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

// RangeFromInclusiveInterval returns a Range representing [start, last].
func RangeFromInclusiveInterval(start, last uintptr) Range {
	return Range{
		Start: start,
		Size:  uint(last - start + 1),
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

// Last returns last uintptr inside the interval.
func (r Range) Last() uintptr {
	return r.Start + uintptr(r.Size) - 1
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

// WithStart returns a range that begins at start and ends at r.End().
func (r Range) WithStart(start uintptr) Range {
	switch {
	case r.Start > start:
		return Range{Start: start, Size: r.Size + uint(r.Start-start)}
	case r.Start == start:
		return Range{Start: start, Size: 0}
	default:
		return Range{Start: start, Size: r.Size - uint(start-r.Start)}
	}
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

// MaxAddr is the highest address in a 64bit address space.
const MaxAddr = ^uintptr(0)

// FindSpaceAbove finds a continuous piece of sz points within Ranges and
// returns a space.Start >= minAddr.
func (rs Ranges) FindSpaceAbove(sz uint, minAddr uintptr) (space Range, err error) {
	return rs.FindSpaceIn(sz, RangeFromInterval(minAddr, MaxAddr))
}

// FindSpaceIn finds a continuous piece of sz points within Ranges and returns
// a Range where space.Start >= limit.Start, with space.End() < limit.End().
func (rs Ranges) FindSpaceIn(sz uint, limit Range) (space Range, err error) {
	return rs.FindSpace(sz, WithinRange(limit))
}

type findSpaceOptions struct {
	limit      Range
	size       uint
	startAlign uint
}

// FindOptioner is a config option for FindSpace.
type FindOptioner func(o *findSpaceOptions)

// WithMinimumAddr requires FindSpace to return a range with an address above
// minAddr.
func WithMinimumAddr(minAddr uintptr) FindOptioner {
	return func(o *findSpaceOptions) {
		o.limit.Start = minAddr
		o.limit.Size -= uint(minAddr)
	}
}

// WithinRange requires FindSpace to return a range within the limit.
func WithinRange(limit Range) FindOptioner {
	return func(o *findSpaceOptions) {
		o.limit = limit
	}
}

// WithAlignment requires FindSpace to return a range with an address and size
// aligned to alignSize.
func WithAlignment(alignSize uint) FindOptioner {
	return func(o *findSpaceOptions) {
		o.size = align.Up(o.size, alignSize)
		o.startAlign = alignSize
	}
}

// WithStartAlignment requires FindSpace to return a range with an address
// aligned to alignSize.
func WithStartAlignment(alignSize uint) FindOptioner {
	return func(o *findSpaceOptions) {
		o.startAlign = alignSize
	}
}

// FindSpace finds a continuous piece of sz points within Ranges and the given
// options and returns the Range pointing to it.
func (rs Ranges) FindSpace(sz uint, opts ...FindOptioner) (Range, error) {
	o := &findSpaceOptions{
		limit: RangeFromInterval(0, MaxAddr),
		size:  sz,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.startAlign != 0 && !align.IsAligned(o.limit.Start, uintptr(o.startAlign)) {
		o.limit = o.limit.WithStart(align.Up(o.limit.Start, uintptr(o.startAlign)))
	}
	for _, r := range rs {
		if o.startAlign != 0 && !align.IsAligned(r.Start, uintptr(o.startAlign)) {
			r = r.WithStart(align.Up(r.Start, uintptr(o.startAlign)))
		}
		if overlap := r.Intersect(o.limit); overlap != nil && overlap.Size >= o.size {
			return Range{Start: overlap.Start, Size: o.size}, nil
		}
	}
	return Range{}, fmt.Errorf("%w: %#x bytes", ErrNotEnoughSpace, sz)
}

// Sort sorts ranges by their start point.
func (rs Ranges) Sort() {
	sort.Slice(rs, func(i, j int) bool {
		if rs[i].Start == rs[j].Start {
			// let rs[i] be the superset of rs[j]
			return rs[i].Size > rs[j].Size
		}

		return rs[i].Start < rs[j].Start
	})
}

// Segment defines kernel memory layout.
type Segment struct {
	// Buf is a buffer to map to Phys in kexec.
	Buf []byte

	// Phys is a physical address of kernel.
	Phys Range
}

// NewSegment creates new Segment.
// Segments should be created using NewSegment method to prevent
// data pointed by Segment.Buf to be collected by garbage collector.
func NewSegment(buf []byte, phys Range) Segment {
	return Segment{
		Buf:  buf,
		Phys: phys,
	}
}

// SegmentEqual returns whether s and t point at the same physical region and
// contain the same data.
func SegmentEqual(s, t Segment) bool {
	return s.Phys == t.Phys && bytes.Equal(s.Buf, t.Buf)
}

func (s Segment) String() string {
	return fmt.Sprintf("(phys: %s, buffer: size %#x)", s.Phys, len(s.Buf))
}

// AlignAndMerge adjusts segs to the preconditions of kexec_load.
//
// Pre-conditions: segs physical ranges are disjoint.
// Post-conditions: segs physical start addresses & size aligned to page size.
func AlignAndMerge(segs Segments) (Segments, error) {
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].Phys.Start == segs[j].Phys.Start {
			// let segs[i] be the superset of segs[j]
			return segs[i].Phys.Size > segs[j].Phys.Size
		}
		return segs[i].Phys.Start < segs[j].Phys.Start
	})

	// We index 0 below.
	if len(segs) == 0 {
		return segs, nil
	}

	// Physical ranges may not overlap.
	//
	// Overlapping ranges could be allowed if the corresponding
	// intersecting buffer ranges contain the same bytes. TBD whether
	// that's needed.
	for i := 0; i < len(segs)-1; i++ {
		if segs[i].Phys.Overlaps(segs[i+1].Phys) {
			return nil, fmt.Errorf("segment %s and %s overlap in the physical space", segs[i], segs[i+1])
		}
	}

	// Since segments' physical ranges are guaranteed to be disjoint, the
	// only condition under which they overlap is if an aligned Phys.Start
	// address overlaps. In that case, merge the segments.
	//
	// The sorting guarantees we can step through linearly.
	var newSegs Segments
	newSegs = append(newSegs, AlignPhysStart(segs[0]))

	for i := 1; i < len(segs); i++ {
		cand := AlignPhysStart(segs[i])

		// Does the aligned segment overlap with the previous
		// segment? We'll have to merge the (unaligned) segment with
		// the last segment.
		if cand.Phys.Overlaps(newSegs[len(newSegs)-1].Phys) {
			if ok := newSegs[len(newSegs)-1].mergeDisjoint(segs[i]); !ok {
				// This should be impossible as long as
				// mergeDisjoint and Overlaps have a matching
				// definition of overlapping.
				return nil, fmt.Errorf("could not merge disjoint segments")
			}
		} else {
			newSegs = append(newSegs, cand)
		}
	}

	// Align the sizes. This is guaranteed to still produce disjoint
	// physical ranges.
	for i := range newSegs {
		// If we adjust the size up, we must curtail the buffer.
		//
		// If Phys.Size = 1 and Buf.Size = 8, the caller only expected
		// 1 byte to go into the physical range.
		//
		// We don't need to deal with the inverse, because kexec_load
		// will fill the remainder of the segment with zeros anyway
		// when buf.Size < phys.Size.
		newSegs[i].Buf = newSegs[i].realBufTruncate()
		newSegs[i].Phys.Size = align.UpPage(newSegs[i].Phys.Size)
	}
	return newSegs, nil
}

// realBufPad adjusts s.Buf.Size = s.Phys.Size. Buf will either gain some zeros
// or be truncated.
func (s Segment) realBufPad() []byte {
	switch {
	case uint(len(s.Buf)) == s.Phys.Size:
		return s.Buf

	case uint(len(s.Buf)) < s.Phys.Size:
		// Pad Buf.
		return append(s.Buf, make([]byte, int(s.Phys.Size-uint(len(s.Buf))))...)

	case uint(len(s.Buf)) > s.Phys.Size:
		// Truncate Buf.
		return s.Buf[:s.Phys.Size]
	}
	return nil
}

// realBufTruncate adjusts s.Buf.Size = s.Phys.Size, except when Buf is smaller
// than Phys. Buf will either remain the same or be truncated.
func (s Segment) realBufTruncate() []byte {
	if uint(len(s.Buf)) > s.Phys.Size {
		return s.Buf[:s.Phys.Size]
	}
	return s.Buf
}

func (s *Segment) mergeDisjoint(s2 Segment) bool {
	if s.Phys.Overlaps(s2.Phys) {
		return false
	}
	// Must be s < s2
	if s.Phys.Start > s2.Phys.Start {
		return false
	}

	a := s.realBufPad()
	// Second half can drop the extra padded zeroes.
	b := s2.realBufTruncate()
	diffSize := s2.Phys.Start - s.Phys.End()
	// Zeros for the middle.
	buf := append(a, make([]byte, int(diffSize))...)
	buf = append(buf, b...)

	phys := s.Phys
	phys.Size += uint(diffSize) + s2.Phys.Size
	*s = NewSegment(buf, phys)
	return true
}

// Align aligns start and size by the given alignSize.
//
// The resulting range is guaranteed to be superset of r.
func (r Range) Align(alignSize uint) Range {
	if alignSize == 0 {
		return r
	}
	s := align.Down(r.Start, uintptr(alignSize))
	// Empty range remains empty.
	if r.Size == 0 {
		return Range{Start: s}
	}
	return Range{
		Start: s,
		Size:  align.Up(r.Size+uint(r.Start-s), alignSize),
	}
}

// AlignPage aligns start and size by page size.
//
// The resulting range is guaranteed to be superset of r.
func (r Range) AlignPage() Range {
	s := align.DownPage(r.Start)
	// Empty range remains empty.
	if r.Size == 0 {
		return Range{Start: s}
	}
	return Range{
		Start: s,
		Size:  align.UpPage(r.Size + uint(r.Start-s)),
	}
}

// AlignPhysStart aligns s.Phys.Start to the page size. AlignPhysStart does not
// align the size of the segment.
func AlignPhysStart(s Segment) Segment {
	orig := s.Phys.Start
	// Find the page address of the starting point.
	s.Phys.Start = s.Phys.Start &^ uintptr(pageMask)
	diff := orig - s.Phys.Start
	s.Phys.Size = s.Phys.Size + uint(diff)

	s.Buf = append(make([]byte, diff), s.Buf...)
	return s
}

// Segments is a collection of segments.
type Segments []Segment

func (segs Segments) String() string {
	var s strings.Builder
	for _, seg := range segs {
		s.WriteString(seg.String())
		s.WriteString("\n")
	}
	return s.String()
}

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

// Phys returns all physical address ranges.
func (segs Segments) Phys() Ranges {
	var r Ranges
	for _, s := range segs {
		r = append(r, s.Phys)
	}
	return r
}

// SegmentsEqual returns whether the contents of all segments are the same,
// while pointing to the same physical memory region.
func SegmentsEqual(s, t Segments) bool {
	if len(s) != len(t) {
		return false
	}
	for i := range s {
		if !SegmentEqual(s[i], t[i]) {
			return false
		}
	}
	return true
}

// IsSupersetOf checks whether all segments in o are present in s and contain
// the same buffer content.
func (segs Segments) IsSupersetOf(o Segments) error {
	for _, seg := range o {
		size := min(seg.Phys.Size, uint(len(seg.Buf)))
		if size == 0 {
			continue
		}
		r := Range{Start: seg.Phys.Start, Size: size}
		buf := segs.GetPhys(r)
		if buf == nil {
			return fmt.Errorf("phys %s not found", r)
		}
		if !bytes.Equal(buf, seg.Buf[:size]) {
			return fmt.Errorf("phys %s contains different bytes", r)
		}
	}
	return nil
}

// GetPhys gets the buffer corresponding to the physical address range r.
func (segs Segments) GetPhys(r Range) []byte {
	for _, seg := range segs {
		if seg.Phys.IsSupersetOf(r) {
			offset := r.Start - seg.Phys.Start
			// TODO: This could be out of range.
			buf := seg.Buf[int(offset) : int(offset)+int(r.Size)]
			return buf
		}
	}
	return nil
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
func (m *Memory) LoadElfSegments(r io.ReaderAt) (Object, error) {
	f, err := ObjectNewFile(r)
	if err != nil {
		return nil, err
	}

	for _, p := range f.Progs() {
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
				return nil, err
			}
			if n < len(d) {
				return nil, fmt.Errorf("not all data of the segment was read")
			}
		}
		// TODO(hugelgupf): check if this is within availableRAM??
		s := NewSegment(d, Range{
			Start: uintptr(p.Paddr),
			Size:  uint(p.Memsz),
		})
		m.Segments.Insert(s)
	}
	return f, nil
}

// M1 is 1 Megabyte in bits.
const M1 = 1 << 20

// FindSpace returns pointer to the physical memory, where array of size sz can
// be stored during next AddKexecSegment call.
//
// Align up to at least a page size if alignSizeBytes is smaller.
func (m Memory) FindSpace(sz, alignSizeBytes uint) (Range, error) {
	if alignSizeBytes == 0 {
		alignSizeBytes = uint(os.Getpagesize())
	}
	sz = align.Up(sz, alignSizeBytes)

	// Don't use memory below 1M, just in case.
	return m.AvailableRAM().FindSpaceAbove(sz, M1)
}

// ReservePhys reserves page-aligned sz bytes in the physical memmap within
// the given limit address range.
func (m *Memory) ReservePhys(sz uint, limit Range) (Range, error) {
	sz = align.UpPage(sz)

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
	r, err := m.FindSpace(uint(len(d)), uint(os.Getpagesize()))
	if err != nil {
		return Range{}, err
	}
	m.Segments.Insert(NewSegment(d, r))
	return r, nil
}

// AddKexecSegmentExplicit adds d to a new kexec segment, but allows asking
// for extra space, secifying alignment size, and setting text_offset.
func (m *Memory) AddKexecSegmentExplicit(d []byte, sz, offset, alignSizeBytes uint) (Range, error) {
	r, err := m.AvailableRAM().FindSpace(offset+sz, WithAlignment(alignSizeBytes))
	if err != nil {
		return Range{}, err
	}
	r.Start = uintptr(uint(r.Start) + offset)
	m.Segments.Insert(NewSegment(d, r))
	return r, nil
}

// AvailableRAM returns page-aligned unused regions of RAM.
//
// AvailableRAM takes all RAM-marked pages in the memory map and subtracts the
// kexec segments already allocated. RAM segments begin at a page boundary.
//
// E.g if page size is 4K and RAM segments are
//
//	[{start:0 size:8192} {start:8192 size:8000}]
//
// and kexec segments are
//
//	[{start:40 size:50} {start:8000 size:2000}]
//
// result should be
//
//	[{start:0 size:40} {start:4096 end:8000 - 4096}]
func (m Memory) AvailableRAM() Ranges {
	ram := m.Phys.RAM()

	// Remove all points we've already reserved from available RAM.
	for _, s := range m.Segments {
		ram = ram.Minus(s.Phys)
	}

	// Only return Ranges starting at an aligned size.
	var alignedRanges Ranges
	for _, r := range ram {
		alignedStart := uintptr(align.UpPage(uint(r.Start)))
		if alignedStart < r.End() {
			alignedRanges = append(alignedRanges, Range{
				Start: alignedStart,
				Size:  r.Size - uint(alignedStart-r.Start),
			})
		}
	}
	return alignedRanges
}
