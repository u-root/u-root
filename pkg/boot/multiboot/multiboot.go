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
	"github.com/u-root/u-root/pkg/ubinary"
	"github.com/u-root/u-root/pkg/uio"
)

const bootloader = "u-root kexec"

// Module describe a module by a ReaderAt and a `CmdLine`
type Module struct {
	Module  io.ReaderAt
	Name    string
	CmdLine string
}

// Modules is a range of module with a Closer interface
type Modules []Module

// multiboot defines parameters for working with multiboot kernels.
type multiboot struct {
	modules []Module
	cmdLine string
}

var (
	rangeTypes = map[kexec.RangeType]uint32{
		kexec.RangeRAM:      1,
		kexec.RangeDefault:  2,
		kexec.RangeACPI:     3,
		kexec.RangeNVS:      4,
		kexec.RangeReserved: 2,
	}
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
	err := binary.Write(&buf, ubinary.NativeEndian, m)
	return buf.Bytes(), err
}

// elems adds mutiboot info elements describing the memory map of the system.
func (m memoryMaps) elems() []elem {
	var e []elem
	for _, mm := range m {
		e = append(e, &mutibootMemRange{
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

// Probe checks if `kernel` is multiboot v1 or mutiboot kernel.
func Probe(kernel io.ReaderAt) error {
	r := tryGzipFilter(kernel)
	_, err := parseHeader(uio.Reader(r))
	if err == ErrHeaderNotFound {
		_, err = parseMutiHeader(uio.Reader(r))
	}
	return err
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
	kernel = tryGzipFilter(kernel)
	for i, mod := range modules {
		modules[i].Module = tryGzipFilter(mod.Module)
	}
	args, err := load(kernel, cmdline, modules, debug, ibft)
	if err != nil {
		return err
	}
	if err := kexec.Load(args.EntryPoint, args.Segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}
	return nil
}

// OpenModules open modules as files and fill a range of `Module` struct
//
// Each module is a path followed by optional command-line arguments, e.g.
// []string{"./module arg1 arg2", "./module2 arg3 arg4"}.
func OpenModules(cmds []string) (Modules, error) {
	modules := make([]Module, len(cmds))
	for i, cmd := range cmds {
		modules[i].CmdLine = cmd
		name := strings.Fields(cmd)[0]
		modules[i].Name = name
		f, err := os.Open(name)
		if err != nil {
			// TODO close already open files
			return nil, fmt.Errorf("error opening module %v: %v", name, err)
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
			CmdLine: cmd,
			Name:    name,
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

type kexecArgs struct {
	Segments   kexec.Segments
	EntryPoint uintptr
}

// load loads decompressed mu(l)tiboot kernel and modules into userspace
// memory, and generates the mu(l)tiboot info structure, memory map, and iBFT
// to be passed to the kernel upon execution.
//
// load returns a set of segments (relocations) that say which piece of
// userspace memory must be moved to which place in physical memory in order to
// jump to the mu(l)tiboot kernel.
func load(kernel io.ReaderAt, cmdLine string, modules []Module, debug bool, ibft *ibft.IBFT) (*kexecArgs, error) {
	m := &multiboot{
		modules: modules,
		cmdLine: cmdLine,
	}
	var mem kexec.Memory

	log.Println("Parsing multiboot header")
	// TODO: the kernel is opened like 4 separate times here. Just open it
	// once and pass it around.

	var header imageType
	multibootHeader, err := parseHeader(uio.Reader(kernel))
	if err == nil {
		header = multibootHeader
	} else if err == ErrHeaderNotFound {
		var mutibootHeader *mutibootHeader
		// We don't even need the header at the moment. Just need to
		// know it's there. Everything that matters is in the ELF.
		mutibootHeader, err = parseMutiHeader(uio.Reader(kernel))
		header = mutibootHeader
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %v", err)
	}
	log.Printf("Found %s image", header.name())

	log.Printf("Getting kernel entry point")
	kernelEntry, err := getEntryPoint(kernel)
	if err != nil {
		return nil, fmt.Errorf("error getting kernel entry point: %v", err)
	}
	log.Printf("Kernel entry point at %#x", kernelEntry)

	log.Printf("Parsing ELF segments")
	if err := mem.LoadElfSegments(kernel); err != nil {
		return nil, fmt.Errorf("error loading ELF segments: %v", err)
	}

	log.Printf("Parsing memory map")
	if err := mem.ParseMemoryMap(); err != nil {
		return nil, fmt.Errorf("error parsing memory map: %v", err)
	}

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
		r, err := mem.FindSpaceAndReserve(uint(len(ibuf)), allowedRange)
		if err != nil {
			return nil, fmt.Errorf("reserving space for the iBFT in %s failed: %v", allowedRange, err)
		}
		log.Printf("iBFT was allocated at %s: %#v", r, ibft)
		mem.Segments.Insert(kexec.NewSegment(ibuf, r))
	}

	log.Printf("Preparing %s info", header.name())
	infoAddr, err := header.addInfo(&mem, m)
	if err != nil {
		return nil, fmt.Errorf("error preparing %s info: %v", header.name(), err)
	}
	log.Printf("Info structure at %#x", infoAddr)

	log.Printf("Adding trampoline")
	entryPoint, err := addTrampoline(&mem, header.bootMagic(), infoAddr, kernelEntry)
	if err != nil {
		return nil, fmt.Errorf("error adding trampoline: %v", err)
	}
	log.Printf("Trampoline entry point at %#x", entryPoint)

	return &kexecArgs{
		Segments:   mem.Segments,
		EntryPoint: entryPoint,
	}, nil
}

func getEntryPoint(r io.ReaderAt) (uintptr, error) {
	f, err := elf.NewFile(r)
	if err != nil {
		return 0, err
	}
	return uintptr(f.Entry), err
}

// addInfo collects and adds mutiboot (without L!) into the segments.
//
// The format is described in the structs in
// https://github.com/vmware/esx-boot/blob/master/include/mutiboot.h
//
// It includes a memory map and a list of modules.
func (*mutibootHeader) addInfo(mem *kexec.Memory, m *multiboot) (addr uintptr, err error) {
	var mi mutibootInfo

	mi.elems = append(mi.elems, memoryMap(mem).elems()...)
	cmdLinePtrs, mods, err := loadModulesAndStrings(mem, m.modules, m.cmdLine)
	if err != nil {
		return 0, err
	}
	mi.elems = append(mi.elems, mods.elems()...)
	mi.cmdline = uint64(cmdLinePtrs[0])

	r, err := mem.AddKexecSegment(mi.marshal())
	if err != nil {
		return 0, err
	}
	return r.Start, nil
}

func memoryMap(mem *kexec.Memory) memoryMaps {
	var ret memoryMaps
	for _, r := range mem.Phys {
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
	log.Printf("Memory map:\n%v", ret)
	return ret
}

func (*header) memoryBoundaries(mem *kexec.Memory) (lower, upper uint32) {
	const M1 = 1048576
	const K640 = 640 * 1024
	for _, r := range mem.Phys {
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

func (*header) addMmapInfo(mem *kexec.Memory) (addr uintptr, size uint, err error) {
	mmap := memoryMap(mem)
	d, err := mmap.marshal()
	if err != nil {
		return 0, 0, err
	}
	r, err := mem.AddKexecSegment(d)
	if err != nil {
		return 0, 0, err
	}
	return r.Start, uint(len(mmap)) * sizeofMemoryMap, nil
}

// addMultibootInfo puts the multiboot info structure into memory.
//
// addInfo collects and adds multiboot info into the relocations/segments.
//
// addInfo marshals out everything required for
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format
// which is a memory map; a list of module structures, pointed to by mods_addr
// and mods_count; and the multiboot info structure itself.
//
// - adds the memory map into memory,
// - adds modules into memory + kernel command-line + bootloader name,
// - adds pointers to modules into memory,
// - then adds multiboot info structure with pointers to the above into mem.
//
// If contiguous memory is available, you'll end up with a layout like this:
//
//  1-2 pages memory map
//  1-2 pages strings (cmdlines, bootloader name)
//  page-aligned modules, one by one
//  1-2 pages with multiboot_mod_list pointers to modules, their size, and their command-lines
//  1 page multiboot_info struct, pointing to mod_list, strings, and memory map
func (h *header) addInfo(mem *kexec.Memory, m *multiboot) (uintptr, error) {
	// Add memory map into memory.
	mmapAddr, mmapSize, err := h.addMmapInfo(mem)
	if err != nil {
		return 0, err
	}
	log.Printf("Memory map at %#x", mmapAddr)

	// Refer to the memory map in the info struct.
	var inf info
	if h.Flags&flagHeaderMemoryInfo != 0 {
		lower, upper := h.memoryBoundaries(mem)
		inf = info{
			Flags:      flagInfoMemMap | flagInfoMemory,
			MemLower:   min(uint32(lower>>10), 0xFFFFFFFF),
			MemUpper:   min(uint32(upper>>10), 0xFFFFFFFF),
			MmapLength: uint32(mmapSize),
			MmapAddr:   uint32(mmapAddr),
		}
	}

	// Load bootloader name, kernel command-line, and modules and their
	// cmdlines into memory.
	extraData := []string{bootloader, m.cmdLine}
	extraPtrs, loaded, err := loadModulesAndStrings(mem, m.modules, extraData...)
	if err != nil {
		return 0, err
	}
	log.Printf("Module / cmdline area: %#x", extraPtrs[0])

	// This loads pointers to modules + cmdline into memory, to be pointed
	// to by the info struct.
	if len(loaded) > 0 {
		b, err := loaded.marshal()
		if err != nil {
			return 0, err
		}
		modRange, err := mem.AddKexecSegment(b)
		if err != nil {
			return 0, err
		}
		log.Printf("Module info area: %#x / %d modules", modRange.Start, len(m.modules))

		inf.Flags |= flagInfoMods
		inf.ModsAddr = uint32(modRange.Start)
		inf.ModsCount = uint32(len(m.modules))
	}

	// Now just fix up some pointers...
	inf.CmdLine = uint32(extraPtrs[1])
	inf.BootLoaderName = uint32(extraPtrs[0])
	inf.Flags |= flagInfoCmdLine | flagInfoBootLoaderName

	log.Printf("Multiboot info: %#v", inf)

	// Add the info struct itself.
	mbInfo, err := inf.marshal()
	if err != nil {
		return 0, err
	}
	mbRange, err := mem.AddKexecSegment(mbInfo)
	if err != nil {
		return 0, err
	}

	if debug {
		info, err := h.description(mem, loaded, mbInfo)
		if err != nil {
			log.Printf("%v cannot create debug info: %v", DebugPrefix, err)
		}
		log.Printf("%v %v", DebugPrefix, info)
	}

	return mbRange.Start, nil
}

func addTrampoline(mem *kexec.Memory, magic, infoAddr, kernelEntry uintptr) (entry uintptr, err error) {
	// Trampoline should be a part of current binary.
	p, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("cannot find current executable path: %v", err)
	}
	pathToSelf, err := filepath.EvalSymlinks(p)
	if err != nil {
		return 0, fmt.Errorf("cannot eval symlinks for %v: %v", p, err)
	}

	// Trampoline setups the machine registers to desired state
	// and executes the loaded kernel.
	d, err := trampoline.Setup(pathToSelf, magic, infoAddr, kernelEntry)
	if err != nil {
		return 0, err
	}

	r, err := mem.AddKexecSegment(d)
	if err != nil {
		return 0, err
	}
	return r.Start, nil
}
