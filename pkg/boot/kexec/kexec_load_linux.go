// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/bzimage"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_INITRD_ADDR_MAX  uint = 0x37FFFFFF
	DEFAULT_BZIMAGE_ADDR_MAX uint = 0x37FFFFFF
)

type purgatory struct {
	name    string
	hexdump string
	code    []byte
}

const defaultPurgatory = "to32bit_3000"

var (
	// Debug is called to print out verbose debug info.
	//
	// Set this to appropriate output stream for display
	// of useful debug info.
	Debug        = log.Printf // func(string, ...interface{}) {}
	curPurgatory = purgatories[defaultPurgatory]
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

// rawLoad is a wrapper around kexec_load(2) syscall.
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

// LoadBzImage loads the given kernel file as to-be-kexeced kernel with
// the given ramfs file and cmdline string.
func KexecLoad(kernel, ramfs io.ReaderAt, cmdline string) error {
	bzimage.Debug = Debug

	/* A collection of vars used for processing the kernel for kexec */

	var err error
	// b is the deserialized bzImage from the kernel
	// io.ReaderAt.
	var b bzimage.BzImage
	// kmem is a struct holding kexec segments.
	//
	// It has routines to work with physical memory
	// ranges.
	var kmem *Memory
	// realMode holds setup code.
	//
	// If not executing in real mode, it will be just
	// be setup header.
	var realMode *LinuxParamHeader

	var kernel16SizeNeeded uint

	var regs32 Entry32Regs

	Debug("Start LoadBzImage...")
	Debug("Try decompressing kernel...")
	kb, err := uio.ReadAll(util.TryGzipFilter(kernel))
	if err != nil {
		return err
	}
	Debug("Try parsing bzImage...")
	if err := b.UnmarshalBinary(kb); err != nil {
		return err
	}
	Debug("Done parsing bzImage.")

	if len(b.KernelCode) < 1024 {
		return fmt.Errorf("kernel code size smaller than 1024 bytes: %d", len(b.KernelCode))
	}

	// TODO(10000TB): add initramfs by a io.ReaderAt.
	//Debug("Try adding initramfs...")
	//if err := b.AddInitRAMFS(ramfs.Name()); err != nil {
	//	return err
	//}

	Debug("Try get ELF from bzImage...")
	kelf, err := b.ELF()
	if err != nil {
		return err
	}
	Debug("Done.")
	kernelEntry := uintptr(kelf.Entry)
	Debug("kernelEntry: %v", kernelEntry)

	// Prepare segments.
	kmem = &Memory{}
	Debug("Try parsing memory map...")
	// TODO(10000TB): refactor this call into initialization of
	// kexec.Memory, as it does not depend on specific boot.
	if err := kmem.ParseMemoryMap(); err != nil {
		return fmt.Errorf("parse memory map: %v", err)
	}

	var relocatableKernel bool
	if b.Header.Protocolversion >= 0x0205 {
		relocatableKernel = b.Header.RelocatableKernel != 0
		Debug("bzImage is relocatable")
	}

	Debug("Loading purgatory...")

	/* Load the trampoline.
	 *
	 * This must load at a higher address than the argument/parameter
	 * segment or the kernel will stomp it's gdt.
	 *
	 * x86_64 purgatory code has got relocations type R_X86_64_32S
	 * that means purgatory got to be loaded within first 2G otherwise
	 * overflow takes place while applying relocations.
	 */

	var purgatoryEntry uintptr

	// Only protected mode is currently supported.
	// In protected mode, kernel need be relocatable, or it will need to fall
	// to real mode executing.

	if !relocatableKernel {
		return errors.New("real mode executing is not supported currently")
	}

	Debug("purgatory entry: %v", purgatoryEntry)

	/* The argument/parameter segment */

	// Assume executing protected mode.
	kernel16SizeNeeded = uint(b.KernelOffset) // Kernel16 size.
	if kernel16SizeNeeded < 4096 {
		kernel16SizeNeeded = 4096
	}

	if err := kmem.LoadElfSegments(bytes.NewReader(b.KernelCode)); err != nil {
		return err
	}

	// TODO(10000TB): Insert cmdline.
	// cmdlineLen := len(cmdline) + 1
	// cmdlineR, err := kmem.FindSpace(cmdlineLen)
	//
	// Do we just modify the bzImage with the cmdline passed in
	// Or we need to allocate a new segments for the cmdline
	// and insert it into segments list to be loaded ?
	//
	// kernelCmdlineLen := 0

	/* Setup segment.
	 *
	 * Currently, only protected mode is implemented. So only copy setup header.
	 */
	realMode = getLinuxParamHeader(&b) // No real mode is executing.
	Debug("got setup header: %v", realMode)

	// No support for kexec on crash.

	var setupRange Range
	// The kernel is a bzImage kernel if the protocol >= 2.00 and the 0x01
	// bit (LOAD_HIGH) in the loadflags field is set
	// TODO(10000TB): check on loadflags.
	if b.Header.Protocolversion >= 0x0200 {
		// TODO(10000TB): Free mem hole start end aligns by
		// max(16, pagesize).
		setupRange, err = kmem.AddPhysSegment(
			realMode.ToBytes(),
			RangeFromInterval(
				uintptr(0x3000),
				uintptr(640*1024),
			),
			// TODO(10000TB): evaluate if we need to provide  option to
			// reserve from end.
			//
			// Our go code defaults to pick up a mem block of requested
			// size from beginning, e.g.
			//
			//   [Range.Start, Range.Start+memsz)
			//
			// Kexec userspace use the range from end, e.g.
			//
			//   [Range.end-memsz+1, Range.end)
			//
		)
		if err != nil {
			return fmt.Errorf("add real mode data and cmdline: %v", err)
		}
	} else {
		// TODO(10000TB): add support.
		return fmt.Errorf("boot protocol version: %v is currently not supported", b.Header.Protocolversion)
	}

	Debug("Loaded real mode data and cmdline at: %v", setupRange)

	/* Verify purgatory loads higher than the parameters. */
	// TODO(10000TB): if rel_addr < setupRange.Start then return error.

	/* Main kernel segment.
	 */
	var mainKernelRange Range
	if b.Header.Protocolversion >= 0x0205 && relocatableKernel {
		// kernelAlign := b.Header.KernelAlignment
		//
		// TODO(10000TB): apply kernel alignment.
		kernel32MaxAddr := DEFAULT_BZIMAGE_ADDR_MAX

		if kernel32MaxAddr > uint(b.Header.InitrdAddrMax) {
			kernel32MaxAddr = uint(b.Header.InitrdAddrMax)
		}

		mainKernelRange, err = kmem.AddPhysSegment(
			// TODO(10000TB): kernel align + reserve from
			// end: [Range.End-memsz+1, Range.end)
			b.KernelCode,
			RangeFromInterval(0x100000, uintptr(kernel32MaxAddr)),
		)
		if err != nil {
			return fmt.Errorf("load main kexec segment: %v", err)
		}
	} else {
		// TODO(10000TB): Impl support for earlier boot protocols.
		return fmt.Errorf("bzImage boot protocol earlier thatn 2.05 is not supported currently: %v", b.Header.Protocolversion)
	}

	Debug("Loaded 32bit kernel at %v", mainKernelRange)

	/* Tell current kernel what is going on.
	 */
	if err = SetupLinuxBootloaderParameters(realMode, kmem, ramfs, kernel16SizeNeeded, uint(setupRange.Start), cmdline); err != nil {
		return fmt.Errorf("setup linux bootloader params: %v", err)
	}

	/*
	 * Initialize the 32bit start information.
	 */
	regs32.Eax = 0 /* Unused */
	regs32.Ebx = 0 /* 0 == boot not AP processor start */
	regs32.Ecx = 0 /* Unused */
	regs32.Edx = 0 /* Unused */
	regs32.Edi = 0 /* unused */
	// TODO(10000TB): Add support.
	// regs32.esp = elf_rel_get_addr(&info->rhdr, "stack_end"); /* stack, unused */
	regs32.Ebp = 0 /* unused */

	// actual parameters ...
	regs32.Esi = uint32(setupRange.Start)      /* kernel parameters */
	regs32.Eip = uint32(mainKernelRange.Start) /* kernel entry point */

	cmdLineEnd := int(setupRange.Start) + int(kernel16SizeNeeded) + len(cmdline) - 1
	Debug("cmdLineEnd: %v", cmdLineEnd)

	if err = SetupLinuxSystemParameters(realMode); err != nil {
		return fmt.Errorf("setup linux system params: %v", err)
	}

	binary.LittleEndian.PutUint64(curPurgatory.code[8:], uint64(mainKernelRange.Start))
	binary.LittleEndian.PutUint64(curPurgatory.code[16:], uint64(setupRange.Start))
	if purgatoryEntry, err = ELFLoad(kmem, curPurgatory.code, 0x3000, 0x7fffffff, -1, 0); err != nil {
		return err
	}

	/* Load it */
	if err = Load(purgatoryEntry, kmem.Segments, 0); err != nil {
		return fmt.Errorf("Kexec Load(%v, %v, %d) = %v", purgatoryEntry, kmem.Segments, 0, err)
	}

	return nil
}

// SelectPurgatory picks a purgatory, returning an error if none is found
func SelectPurgator(name string) error {
	p, ok := purgatories[name]
	if !ok {
		var s []string
		for i := range purgatories {
			s = append(s, i)
		}
		return fmt.Errorf("%s: no such purgatory, try one of %v", name, s)

	}
	curPurgatory = p
	return nil
}
