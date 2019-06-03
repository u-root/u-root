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
	"path/filepath"

	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot/internal/trampoline"
	"github.com/u-root/u-root/pkg/ubinary"
)

const bootloader = "u-root kexec"

// multiboot defines parameters for working with multiboot kernels.
type multiboot struct {
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

// MemoryMap represents a reserved range of memory passed via the multiboot Info header.
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

type MemoryMaps []MemoryMap

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

// newMB returns a new multiboot instance.
func newMB(file, cmdLine string, modules []string) (*multiboot, error) {
	// Trampoline should be a part of current binary.
	p, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("Cannot find current executable path: %v", err)
	}
	trampoline, err := filepath.EvalSymlinks(p)
	if err != nil {
		return nil, fmt.Errorf("Cannot eval symlinks for %v: %v", p, err)
	}

	return &multiboot{
		file:       file,
		modules:    modules,
		cmdLine:    cmdLine,
		trampoline: trampoline,
		bootloader: bootloader,
		mem:        kexec.Memory{},
	}, nil
}

// Load parses and loads a multiboot kernel `file` using kexec_load.
//
// debug turns on debug logging.
//
// Load can set up an arbitrary number of modules, and takes care of the
// multiboot info structure, including the memory map.
//
// After Load is called, kexec.Reboot() is ready to be called any time to stop
// Linux and execute the loaded kernel.
func Load(debug bool, file, cmdline string, modules []string) error {
	m, err := newMB(file, cmdline, modules)
	if err != nil {
		return err
	}
	if err := m.load(debug); err != nil {
		return err
	}
	if err := kexec.Load(m.EntryPoint, m.Segments(), 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}
	return nil
}

// load loads and parses multiboot information from m.file.
func (m *multiboot) load(debug bool) error {
	log.Printf("Parsing file %v", m.file)
	b, err := readFile(m.file)
	if err != nil {
		return err
	}
	kernel := kernelReader{buf: b}
	log.Println("Parsing multiboot Header")
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
	if err := m.mem.ParseMemoryMap(); err != nil {
		return fmt.Errorf("Error parsing memory map: %v", err)
	}

	log.Printf("Preparing multiboot Info")
	if m.infoAddr, err = m.addInfo(); err != nil {
		return fmt.Errorf("Error preparing multiboot Info: %v", err)
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

func getEntryPoint(r io.ReaderAt) (uintptr, error) {
	f, err := elf.NewFile(r)
	if err != nil {
		return 0, err
	}
	return uintptr(f.Entry), err
}

func (m *multiboot) addInfo() (addr uintptr, err error) {
	iw, err := m.newmultibootInfo()
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

func (m multiboot) memoryMap() MemoryMaps {
	var ret MemoryMaps
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

func (m *multiboot) addMmap() (addr uintptr, size uint, err error) {
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

func (m multiboot) memoryBoundaries() (lower, upper uint32) {
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

func (m *multiboot) newmultibootInfo() (*infoWrapper, error) {
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
func (m multiboot) Segments() []kexec.Segment {
	return m.mem.Segments
}

// marshal writes out the exact bytes expected by the multiboot info header
// specified in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format.
func (m MemoryMaps) marshal() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, ubinary.NativeEndian, m)
	return buf.Bytes(), err
}

func (m *multiboot) addTrampoline() (entry uintptr, err error) {
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
