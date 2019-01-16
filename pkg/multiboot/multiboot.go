// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package multiboot implements basic primitives
// to load multiboot kernels as defined in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html.
package multiboot

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot/internal/trampoline"
	"github.com/u-root/u-root/pkg/ubinary"
)

const bootloader = "u-root kexec"

var (
	// Debug can be set to, e.g, log.Printf for debugging
	// For now, we are leaving it on. Once we have done enough
	// testing, we can turn it off.
	Debug = log.Printf // func(string, ...interface{}) {}
)

// Multiboot defines parameters for working with multiboot kernels.
type Multiboot struct {
	mem kexec.Memory

	file    string
	modules []string

	cmdLine    string
	bootloader string

	// trampoline is a path to an executable blob, which contains a trampoline segment.
	// Trampoline sets machine to a specific state defined by multiboot v1 spec.
	// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.
	trampoline string

	header Header

	// infoAddr is a pointer to multiboot info.
	infoAddr uintptr
	// kernelEntry is a pointer to entry point of kernel.
	kernelEntry uintptr
	// EntryPoint is a pointer to trampoline.
	EntryPoint uintptr

	info          Info
	loadedModules []Module
}

var (
	rangeTypes = map[kexec.RangeType]uint32{
		kexec.RangeRAM:      1,
		kexec.RangeDefault:  2,
		kexec.RangeACPI:     3,
		kexec.RangeNVS:      4,
		kexec.RangeReserved: 2,
	}
	PageSize = os.Getpagesize()
)

var sizeofMemoryMap = uint(binary.Size(MemoryMap{}))

// MemoryMap represents a reserved range of memory passed via the Multiboot Info header.
type MemoryMap struct {
	// Size is the size of the associated structure in bytes.
	Size uint32
	// BaseAddr is the starting address.
	BaseAddr uint64
	// Length is the size of the memory region in bytes.
	Length uint64
	// Type is the variety of address range represented.
	Type uint32
}

type memoryMaps []MemoryMap

// Probe checks if file is multiboot v1 kernel.
func Probe(file string) error {
	b, err := readFile(file)
	if err != nil {
		return err
	}
	kernel := &kernelReader{buf: b}
	_, err = parseHeader(kernel)
	return err
}

// New returns a new Multiboot instance.
func New(file, cmdLine, trampoline string, modules []string) *Multiboot {
	return &Multiboot{
		file:       file,
		modules:    modules,
		cmdLine:    cmdLine,
		trampoline: trampoline,
		bootloader: bootloader,
		mem:        kexec.Memory{},
	}
}

// Load loads and parses multiboot information from m.file.
func (m *Multiboot) Load(debug bool) error {
	log.Printf("Parsing file %v", m.file)
	b, err := readFile(m.file)
	if err != nil {
		return err
	}
	kernel := kernelReader{buf: b}
	log.Println("Parsing Multiboot Header")
	if m.header, err = parseHeader(&kernel); err != nil {
		return fmt.Errorf("Error parsing headers: %v", err)
	}

	log.Printf("Getting kernel entry point")
	if m.kernelEntry, err = getEntryPoint(kernel); err != nil {
		return fmt.Errorf("Error getting kernel entry point: %v", err)
	}

	log.Printf("Parsing ELF segments")
	if err := m.mem.LoadElfSegments(kernel); err != nil {
		return fmt.Errorf("Error loading ELF segments: %v", err)
	}

	log.Printf("Parsing memory map")
	if err := m.mem.ParseFromMemmap(); err != nil {
		return fmt.Errorf("Error parsing memory map: %v", err)
	}

	log.Printf("Preparing Multiboot Info")
	if m.infoAddr, err = m.addInfo(); err != nil {
		return fmt.Errorf("Error preparing Multiboot Info: %v", err)
	}

	log.Printf("Adding trampoline")
	if m.EntryPoint, err = m.addTrampoline(); err != nil {
		return fmt.Errorf("Error adding trampoline: %v", err)
	}

	if debug {
		info, err := m.Description()
		if err != nil {
			log.Printf("%v cannot create debug info: %v", DebugPrefix, err)
		}
		log.Printf("%v %v", DebugPrefix, info)
	}

	return nil
}

// ACPI adds an ACPI segment to be added to the existing ACPI
// tables. This will currently only work for one new ACPI table.  If
// there is ever a need (unlikely!) we can add support for 2 or more
// later.  N.B. Multiboot2 lets us name the new RSDP, which would have
// been nice.  Multiboot1 only mentions the APM table, which is
// surprising; I had no idea it was that old.
func (m *Multiboot) ACPI(n string) error {
	xtra, err := acpi.RawFromFile(n)
	if err != nil {
		return err
	}
	return m.ACPITable(xtra)
}

func (m *Multiboot) ACPITable(t ...acpi.Tabler) error {
	// NewSDT won't work on some linux kernels that limit reading
	// above 1m. So we have to read the rsdp, which seems ok; then cons up
	// a new SDT from scratch, since there is not one in /sys.
	_, r, err := acpi.GetRSDP()
	if err != nil {
		return err
	}

	s, err := acpi.NewSDT()
	if err != nil {
		return err
	}
	s.Base = r.Base()
	// We can't use the table pointers in SDT as Linux seems to want to deny
	// access. So we have an SDT full of pointers we can't use. Awesome!
	// And we can't read the SDT. I'm just chock full of good news today.
	// acpi.RawTables reads from files in /sys
	rr, err := acpi.RawTables()
	if err != nil {
		return err
	}
	rr = append(rr, t...)

	log.Printf("Calling SDT Marshal with %d tables", len(rr))
	b, err := s.MarshalAll(rr...)
	if err != nil {
		return err
	}
	log.Printf("Marshal'ed %d bytes for ACPI segment", len(b))
	a := kexec.NewSegment(b, kexec.Range{Start: uintptr(s.Base), Size: uint(len(b))})
	if err := m.EnsureSize(a.Phys, kexec.RangeACPI); err != nil {
		// Last chance: maybe you have buggy crap bios and it marks it
		// incorrectly.
		if err := m.EnsureSize(a.Phys, kexec.RangeReserved); err != nil {
			return err
		}
	}
	m.mem.Segments = append(m.mem.Segments, a)
	return nil
}

func getEntryPoint(r io.ReaderAt) (uintptr, error) {
	f, err := elf.NewFile(r)
	if err != nil {
		return 0, err
	}
	return uintptr(f.Entry), err
}

func (m *Multiboot) addInfo() (addr uintptr, err error) {
	iw, err := m.newMultibootInfo()
	if err != nil {
		return 0, err
	}
	infoSize, err := iw.size()
	if err != nil {
		return 0, err
	}

	addr, err = m.mem.FindSpace(infoSize)
	if err != nil {
		return 0, err
	}

	d, err := iw.marshal(addr)
	if err != nil {
		return 0, err
	}
	m.info = iw.Info

	addr, err = m.mem.AddKexecSegment(d)
	if err != nil {
		return 0, err
	}
	return addr, nil
}

func (m Multiboot) memoryMap() memoryMaps {
	var ret memoryMaps
	for _, r := range m.mem.Phys {
		typ, ok := rangeTypes[r.Type]
		if !ok {
			typ = rangeTypes[kexec.RangeDefault]
		}
		v := MemoryMap{
			// Size is really used for skipping to the next pair.
			Size:     uint32(sizeofMemoryMap) - 4,
			BaseAddr: uint64(r.Start),
			Length:   uint64(r.Size) + 1,
			Type:     typ,
		}
		ret = append(ret, v)
	}
	return ret
}

// EnsureSize ensures m.mem.Phys can contain a kexec.Range of a kexec.RangeType.
// The only use at present is ensuring that our ACPI tables live safely
// in an ACPI or Reserved area.
// If not, it will grow the memmap region at the expense of the following
// region. The following region must be kexec.RangeRAM
// There are a few rules:
// Start must be in the region
//   i.e. for this implementation we don't grow down.
//   This is a limitation but I doubt a serious one
// Type must match exactly
// Region must be growing or at least not shrinking -- if shrinking, but it fits, we leave it alone.
// There must be a following map entry of type kexec.RangeRam
// which we will shrink, and it must be at least one page in size.
//   We still do not handle what happens when ACPI grows into the ACPINVS range.
//   Since ACPI NVS can be pointed to be (e.g.) SMM code, for S3 and S4, moving it
//   around is not really an option. In future, we may have to move the ACPI table,
//   or make it discontiguous, if we run out of memory. This need can be acommodated
//   in the ACPI function.
func (m Multiboot) EnsureSize(r kexec.Range, typ kexec.RangeType) error {
	Debug("find %#x %v", r, typ)
	for i, mr := range m.mem.Phys {
		Debug("check %#x", mr)
		if r.Start < mr.Start {
			Debug("start %#x too low for %#x", r.Start, mr.Start)
			continue
		}
		if r.Start > mr.Start+uintptr(mr.Size) {
			Debug("start %#x too high for %#x", r.Start, mr.Start+uintptr(mr.Size))
			continue
		}

		if mr.Type != typ {
			Debug("type wrong")
			continue
		}
		if mr.Start+uintptr(mr.Size) >= r.Start+uintptr(r.Size) {
			Debug("size is ok")
			return nil
		}
		Debug("need to adjust")
		if i == len(m.mem.Phys)-1 {
			return fmt.Errorf("Growmap: can't grow last segment element")
		}
		n := m.mem.Phys[i+1]
		adjust := uint(r.Start + uintptr(r.Size+uint(PageSize)-1-mr.Size) & ^uintptr(PageSize-1))
		if n.Size < adjust {
			return fmt.Errorf("Growmap: next segment len is %d, must be at least %d", mr.Size, adjust)
		}
		Debug("adjust it")
		mr.Size += adjust
		n.Start += uintptr(adjust)
		n.Size -= adjust
		return nil
	}
	return fmt.Errorf("Can not ensure %#x type %v: no room", r, typ)
}

func (m *Multiboot) addMmap() (addr uintptr, size uint, err error) {
	mmap := m.memoryMap()
	d, err := mmap.marshal()
	if err != nil {
		return 0, 0, err
	}
	addr, err = m.mem.AddKexecSegment(d)
	if err != nil {
		return 0, 0, err
	}
	return addr, uint(len(mmap)) * sizeofMemoryMap, nil
}

func (m Multiboot) memoryBoundaries() (lower, upper uint32) {
	const M1 = 1048576
	const K640 = 640 * 1024
	for _, r := range m.mem.Phys {
		if r.Type != kexec.RangeRAM {
			continue
		}
		end := uint32(r.Start) + uint32(r.Size)
		// Lower memory starts at address 0, and upper memory starts at address 1 megabyte.
		// The maximum possible value for lower memory is 640 kilobytes.
		// The value returned for upper memory is maximally the address of the first upper memory hole minus 1 megabyte.
		// It is not guaranteed to be this value.
		if r.Start <= K640 && end > lower {
			lower = end
		}
		if r.Start <= M1 && end > upper+M1 {
			upper = end - M1
		}
	}
	return
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func (m *Multiboot) newMultibootInfo() (*infoWrapper, error) {
	mmapAddr, mmapSize, err := m.addMmap()
	if err != nil {
		return nil, err
	}
	var info Info
	if m.header.Flags&flagHeaderMemoryInfo != 0 {
		lower, upper := m.memoryBoundaries()
		info = Info{
			Flags:      flagInfoMemMap | flagInfoMemory,
			MemLower:   min(uint32(lower>>10), 0xFFFFFFFF),
			MemUpper:   min(uint32(upper>>10), 0xFFFFFFFF),
			MmapLength: uint32(mmapSize),
			MmapAddr:   uint32(mmapAddr),
		}
	}

	if len(m.modules) > 0 {
		modAddr, err := m.addModules()
		if err != nil {
			return nil, err
		}
		info.Flags |= flagInfoMods
		info.ModsAddr = uint32(modAddr)
		info.ModsCount = uint32(len(m.modules))
	}

	info.CmdLine = sizeofInfo
	info.BootLoaderName = sizeofInfo + uint32(len(m.cmdLine)) + 1
	info.Flags |= flagInfoCmdLine | flagInfoBootLoaderName
	return &infoWrapper{
		Info:           info,
		CmdLine:        m.cmdLine,
		BootLoaderName: m.bootloader,
	}, nil
}

// Segments returns kexec.Segments, where all the multiboot related
// information is stored.
func (m Multiboot) Segments() []kexec.Segment {
	return m.mem.Segments
}

// marshal writes out the exact bytes expected by the multiboot info header
// specified in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format.
func (m memoryMaps) marshal() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, ubinary.NativeEndian, m)
	return buf.Bytes(), err
}

func (m *Multiboot) addTrampoline() (entry uintptr, err error) {
	// Trampoline setups the machine registers to desired state
	// and executes the loaded kernel.
	d, err := trampoline.Setup(m.trampoline, m.infoAddr, m.kernelEntry)
	if err != nil {
		return 0, err
	}

	addr, err := m.mem.AddKexecSegment(d)
	if err != nil {
		return 0, err
	}

	return addr, nil
}
