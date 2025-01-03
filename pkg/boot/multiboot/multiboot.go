// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package multiboot implements bootloading multiboot kernels as defined by
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html.
//
// Package multiboot crafts kexec segments that can be used with the kexec_load
// system call.
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
	"strings"

	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/multiboot/internal/trampoline"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/uio/uio"
)

const bootloader = "u-root kexec"

// Module describe a module by a ReaderAt and a `Cmdline`
type Module struct {
	Module  io.ReaderAt
	Cmdline string
}

// Name returns the first field of the cmdline, if there is one.
func (m Module) Name() string {
	f := strings.Fields(m.Cmdline)
	if len(f) > 0 {
		return f[0]
	}
	return ""
}

// Modules is a range of module with a Closer interface
type Modules []Module

// multiboot defines parameters for working with multiboot kernels.
type multiboot struct {
	mem kexec.Memory

	kernel  io.ReaderAt
	modules []Module

	cmdLine    string
	bootloader string

	// trampoline is a path to an executable blob, which contains a trampoline segment.
	// The trampoline sets the machine to a specific state defined by multiboot v1 spec.
	// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.
	trampoline string

	// EntryPoint is a pointer to trampoline.
	entryPoint uintptr

	info          info
	loadedModules modules
}

var rangeTypes = map[kexec.RangeType]uint32{
	kexec.RangeRAM:      1,
	kexec.RangeDefault:  2,
	kexec.RangeACPI:     3,
	kexec.RangeNVS:      4,
	kexec.RangeReserved: 2,
}

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

// String returns a readable representation of a MemoryMap entry.
func (m MemoryMap) String() string {
	return fmt.Sprintf("[0x%x, 0x%x) (len: 0x%x, size: 0x%x, type: %d)", m.BaseAddr, m.BaseAddr+m.Length, m.Length, m.Size, m.Type)
}

type memoryMaps []MemoryMap

// marshal writes out the exact bytes expected by the multiboot info header
// specified in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format.
func (m memoryMaps) marshal() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.NativeEndian, m)
	return buf.Bytes(), err
}

// elems adds esxBootInfo info elements describing the memory map of the system.
func (m memoryMaps) elems() []elem {
	var e []elem
	for _, mm := range m {
		e = append(e, &esxBootInfoMemRange{
			startAddr: mm.BaseAddr,
			length:    mm.Length,
			memType:   mm.Type,
		})
	}
	return e
}

// String returns a new-line-separated representation of the entire memory map.
func (m memoryMaps) String() string {
	var s []string
	for _, mm := range m {
		s = append(s, mm.String())
	}
	return strings.Join(s, "\n")
}

// Probe checks if `kernel` is multiboot v1 or esxBootInfo kernel.
// If the `kernel` is gzip'ed, it will decompress it.
// Only Gzip decmpression is supported at present.
func Probe(kernel io.ReaderAt) error {
	r := util.TryGzipFilter(kernel)
	_, err := parseHeader(uio.Reader(r))
	if err == ErrHeaderNotFound {
		_, err = parseMutiHeader(uio.Reader(r))
	}
	return err
}

// newMB returns a new multiboot instance.
func newMB(kernel io.ReaderAt, cmdLine string, modules []Module) (*multiboot, error) {
	// Trampoline should be a part of current binary.
	p, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot find current executable path: %w", err)
	}
	trampoline, err := filepath.EvalSymlinks(p)
	if err != nil {
		return nil, fmt.Errorf("cannot eval symlinks for %v: %w", p, err)
	}

	return &multiboot{
		kernel:     kernel,
		modules:    modules,
		cmdLine:    cmdLine,
		trampoline: trampoline,
		bootloader: bootloader,
		mem:        kexec.Memory{},
	}, nil
}

// Load parses and loads a multiboot `kernel` using kexec_load.
//
// debug turns on debug logging.
//
// Load can set up an arbitrary number of modules, and takes care of the
// multiboot info structure, including the memory map.
//
// After Load is called, kexec.Reboot() is ready to be called any time to stop
// Linux and execute the loaded kernel.
func Load(debug bool, kernel io.ReaderAt, cmdline string, modules []Module, ibft *ibft.IBFT) error {
	entryPoint, segments, err := PrepareLoad(debug, kernel, cmdline, modules, ibft)
	if err != nil {
		return err
	}
	if err := kexec.Load(entryPoint, segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %w", err)
	}
	return nil
}

// PrepareLoad parses and loads a multiboot `kernel` ready for kexec_load. It
// returns an entry point value and segments.
//
// Load can set up an arbitrary number of modules, and takes care of the
// multiboot info structure, including the memory map.
func PrepareLoad(debug bool, kernel io.ReaderAt, cmdline string, modules []Module, ibft *ibft.IBFT) (uintptr, kexec.Segments, error) {
	kernel = util.TryGzipFilter(kernel)
	for i, mod := range modules {
		modules[i].Module = util.TryGzipFilter(mod.Module)
	}

	m, err := newMB(kernel, cmdline, modules)
	if err != nil {
		return 0, nil, err
	}
	if err := m.load(debug, ibft); err != nil {
		return 0, nil, err
	}
	return m.entryPoint, m.mem.Segments, nil
}

// OpenModules open modules as files and fill a range of `Module` struct
//
// Each module is a path followed by optional command-line arguments, e.g.
// []string{"./module arg1 arg2", "./module2 arg3 arg4"}.
func OpenModules(cmds []string) (Modules, error) {
	modules := make([]Module, len(cmds))
	for i, cmd := range cmds {
		modules[i].Cmdline = cmd
		name := strings.Fields(cmd)[0]
		f, err := os.Open(name)
		if err != nil {
			// TODO close already open files
			return nil, fmt.Errorf("error opening module %v: %w", name, err)
		}
		modules[i].Module = f
	}
	return modules, nil
}

// LazyOpenModules assigns modules to be opened as files.
//
// Each module is a path followed by optional command-line arguments, e.g.
// []string{"./module arg1 arg2", "./module2 arg3 arg4"}.
func LazyOpenModules(cmds []string) Modules {
	modules := make([]Module, 0, len(cmds))
	for _, cmd := range cmds {
		name := strings.Fields(cmd)[0]
		modules = append(modules, Module{
			Cmdline: cmd,
			Module:  uio.NewLazyFile(name),
		})
	}
	return modules
}

// Close closes all Modules ReaderAt implementing the io.Closer interface
func (m Modules) Close() error {
	// poor error handling inspired from uio.multiCloser
	var allErr error
	for _, mod := range m {
		if c, ok := mod.Module.(io.Closer); ok {
			if err := c.Close(); err != nil {
				allErr = err
			}
		}
	}
	return allErr
}

// load loads and parses multiboot information from m.kernel.
func (m *multiboot) load(debug bool, ibft *ibft.IBFT) error {
	var err error
	log.Println("Parsing multiboot header")
	// TODO: the kernel is opened like 4 separate times here. Just open it
	// once and pass it around.

	var header imageType
	multibootHeader, err := parseHeader(uio.Reader(m.kernel))
	if err == nil {
		header = multibootHeader
	} else if err == ErrHeaderNotFound {
		var esxBootInfoHeader *esxBootInfoHeader
		// We don't even need the header at the moment. Just need to
		// know it's there. Everything that matters is in the ELF.
		esxBootInfoHeader, err = parseMutiHeader(uio.Reader(m.kernel))
		header = esxBootInfoHeader
	}
	if err != nil {
		return fmt.Errorf("error parsing headers: %w", err)
	}
	log.Printf("Found %s image", header.name())

	log.Printf("Getting kernel entry point")
	kernelEntry, err := getEntryPoint(m.kernel)
	if err != nil {
		return fmt.Errorf("error getting kernel entry point: %w", err)
	}
	log.Printf("Kernel entry point at %#x", kernelEntry)

	log.Printf("Parsing ELF segments")
	if _, err := m.mem.LoadElfSegments(m.kernel); err != nil {
		return fmt.Errorf("error loading ELF segments: %w", err)
	}

	log.Printf("Parsing memory map")
	memmap, err := kexec.MemoryMapFromSysfsMemmap()
	if err != nil {
		return fmt.Errorf("error parsing memory map: %w", err)
	}
	m.mem.Phys = memmap

	// Insert the iBFT now, since nothing else has been allocated and this
	// is the most restricted allocation we're gonna have to make.
	if ibft != nil {
		ibuf := ibft.Marshal()

		// The iBFT may sit between 512K and 1M in physical memory. The
		// loaded OS finds it by scanning this region.
		allowedRange := kexec.Range{
			Start: 0x80000,
			Size:  0x80000,
		}
		r, err := m.mem.ReservePhys(uint(len(ibuf)), allowedRange)
		if err != nil {
			return fmt.Errorf("reserving space for the iBFT in %s failed: %w", allowedRange, err)
		}
		log.Printf("iBFT was allocated at %s: %#v", r, ibft)
		m.mem.Segments.Insert(kexec.NewSegment(ibuf, r))
	}

	log.Printf("Preparing %s info", header.name())
	infoAddr, err := header.addInfo(m)
	if err != nil {
		return fmt.Errorf("error preparing %s info: %w", header.name(), err)
	}

	log.Printf("Adding trampoline")
	if m.entryPoint, err = m.addTrampoline(header.bootMagic(), infoAddr, kernelEntry); err != nil {
		return fmt.Errorf("error adding trampoline: %w", err)
	}
	log.Printf("Trampoline entry point at %#x", m.entryPoint)

	if debug {
		info, err := m.description()
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

// addInfo collects and adds multiboot info into the relocations/segments.
//
// addInfo marshals out everything required for
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format
// which is a memory map; a list of module structures, pointed to by mods_addr
// and mods_count; and the multiboot info structure itself.
func (h *header) addInfo(m *multiboot) (addr uintptr, err error) {
	iw, err := h.newMultibootInfo(m)
	if err != nil {
		return 0, err
	}
	infoSize, err := iw.size()
	if err != nil {
		return 0, err
	}

	r, err := m.mem.FindSpace(infoSize, uint(os.Getpagesize()))
	if err != nil {
		return 0, err
	}

	d, err := iw.marshal(r.Start)
	if err != nil {
		return 0, err
	}
	m.info = iw.info

	m.mem.Segments.Insert(kexec.NewSegment(d, r))
	return r.Start, nil
}

// addInfo collects and adds esxBootInfo (without L!) into the segments.
//
// The format is described in the structs in
// https://github.com/vmware/esx-boot/blob/master/include/esxbootinfo.h
//
// It includes a memory map and a list of modules.
func (*esxBootInfoHeader) addInfo(m *multiboot) (addr uintptr, err error) {
	var mi esxBootInfoInfo

	mi.elems = append(mi.elems, m.memoryMap().elems()...)
	mods, err := m.loadModules()
	if err != nil {
		return 0, err
	}
	mi.elems = append(mi.elems, mods.elems()...)

	// This marshals the esxBootInfo info with cmdline = 0. We're gonna append
	// the cmdline, so we must know the size of the marshaled stuff first
	// to be able to point to it.
	//
	// TODO: find a better place to put the cmdline so we don't do this
	// bullshit.
	b := mi.marshal()

	// string + null-terminator
	cmdlineLen := len(m.cmdLine) + 1

	memRange, err := m.mem.FindSpace(uint(len(b)+cmdlineLen), uint(os.Getpagesize()))
	if err != nil {
		return 0, err
	}
	mi.cmdline = uint64(memRange.Start + uintptr(len(b)))

	// Re-marshal, now that the cmdline is set.
	b = mi.marshal()
	b = append(b, []byte(m.cmdLine)...)
	b = append(b, 0)
	m.mem.Segments.Insert(kexec.NewSegment(b, memRange))
	return memRange.Start, nil
}

func (m multiboot) memoryMap() memoryMaps {
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
			Length:   uint64(r.Size),
			Type:     typ,
		}
		ret = append(ret, v)
	}
	log.Printf("Memory map: %v", ret)
	return ret
}

// addMmap adds a multiboot-marshaled memory map in memory.
func (m *multiboot) addMmap() (addr uintptr, size uint, err error) {
	mmap := m.memoryMap()
	d, err := mmap.marshal()
	if err != nil {
		return 0, 0, err
	}
	r, err := m.mem.AddKexecSegment(d)
	if err != nil {
		return 0, 0, err
	}
	return r.Start, uint(len(mmap)) * sizeofMemoryMap, nil
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

func (h *header) newMultibootInfo(m *multiboot) (*infoWrapper, error) {
	mmapAddr, mmapSize, err := m.addMmap()
	if err != nil {
		return nil, err
	}
	var inf info
	if h.Flags&flagHeaderMemoryInfo != 0 {
		lower, upper := m.memoryBoundaries()
		inf = info{
			Flags:      flagInfoMemMap | flagInfoMemory,
			MemLower:   min(uint32(lower>>10), 0xFFFFFFFF),
			MemUpper:   min(uint32(upper>>10), 0xFFFFFFFF),
			MmapLength: uint32(mmapSize),
			MmapAddr:   uint32(mmapAddr),
		}
	}

	if len(m.modules) > 0 {
		modAddr, err := m.addMultibootModules()
		if err != nil {
			return nil, err
		}
		inf.Flags |= flagInfoMods
		inf.ModsAddr = uint32(modAddr)
		inf.ModsCount = uint32(len(m.modules))
	}

	return &infoWrapper{
		info:           inf,
		Cmdline:        m.cmdLine,
		BootLoaderName: m.bootloader,
	}, nil
}

func (m *multiboot) addTrampoline(magic, infoAddr, kernelEntry uintptr) (entry uintptr, err error) {
	// Trampoline setups the machine registers to desired state
	// and executes the loaded kernel.
	d, err := trampoline.Setup("", magic, infoAddr, kernelEntry)
	if err != nil {
		return 0, err
	}

	r, err := m.mem.AddKexecSegment(d)
	if err != nil {
		return 0, err
	}
	return r.Start, nil
}
