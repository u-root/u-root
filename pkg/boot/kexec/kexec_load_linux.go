// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/bzimage"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_INITRD_ADDR_MAX  uint = 0x37FFFFFF
	DEFAULT_BZIMAGE_ADDR_MAX uint = 0x37FFFFFF
	bootParams                    = "/sys/kernel/boot_params/data"
)

const defaultPurgatory = "default"

var (
	// Debug is called to print out verbose debug info.
	//
	// Set this to appropriate output stream for display
	// of useful debug info.
	Debug        = log.Printf // func(string, ...interface{}) {}
	curPurgatory = linux.Purgatories[defaultPurgatory]
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

// KexecLoad loads the given kernel file as to-be-kexeced kernel with
// the given ramfs file and cmdline string. It uses the kexec "classic"
// system call, i.e. memory segments + entry point.
func KexecLoad(kernel, ramfs *os.File, cmdline string) error {
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
	// TODO(10000TB): construct default params in go.
	//
	// boot_params directory is x86 specific. So for now, following code only
	// works on x86.
	// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-kernel-boot_params
	bp, err := ioutil.ReadFile("/sys/kernel/boot_params/data")
	if err != nil {
		return fmt.Errorf("reading boot_param data: %w", err)
	}
	var lp = &bzimage.LinuxParams{}
	if err := lp.UnmarshalBinary(bp); err != nil {
		return fmt.Errorf("unmarshaling header: %w", err)
	}
	Debug("Start LoadBzImage...")
	Debug("Try decompressing kernel...")
	kb, err := uio.ReadAll(kernel)
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
	if b.Header.Protocolversion < 0x0205 {
		return fmt.Errorf("bzImage boot protocol earlier thatn 2.05 is not supported currently: %v", b.Header.Protocolversion)
	}
	relocatableKernel = b.Header.RelocatableKernel != 0
	// Only protected mode is currently supported.
	// In protected mode, kernel need be relocatable, or it will need to fall
	// to real mode executing.
	if !relocatableKernel {
		return errors.New("non-relocateable Kernels are not supported")
	}
	Debug("Loading purgatory...")
	if _, err := kmem.LoadElfSegments(bytes.NewReader(b.KernelCode)); err != nil {
		return err
	}

	var ramfsRange Range
	if ramfs != nil {
		b, err := ioutil.ReadAll(ramfs)
		if err != nil {
			return fmt.Errorf("unable to read initramfs: %w", err)
		}
		if ramfsRange, err = kmem.AddKexecSegment(b); err != nil {
			return fmt.Errorf("add initramfs segment: %v", err)
		}
		Debug("Added %d byte initramfs at %s", len(b), ramfsRange)
		lp.Initrdstart = uint32(ramfsRange.Start)
		lp.Initrdsize = uint32(ramfsRange.Size)
	}

	Debug("Kernel cmdline to append: %s", cmdline)
	if len(cmdline) > 0 {
		var cmdlineRange Range
		Debug("Add cmdline: %s", cmdline)
		cmdlineBytes := []byte(cmdline)
		if cmdlineRange, err = kmem.AddKexecSegment(cmdlineBytes); err != nil {
			return fmt.Errorf("add cmdline segment: %v", err)
		}
		Debug("Added %d byte of cmdline at %s", len(cmdlineBytes), cmdlineRange)
		lp.CLPtr = uint32(cmdlineRange.Start)      // 2.02+
		lp.CmdLineSize = uint32(cmdlineRange.Size) // 2.06+
	}

	var setupRange Range
	// The kernel is a bzImage kernel if the protocol >= 2.00 and the 0x01
	// bit (LOAD_HIGH) in the loadflags field is set
	// TODO(10000TB): check on loadflags.
	linuxParam, err := lp.MarshalBinary()
	if err != nil {
		return fmt.Errorf("re-marshaling header: %w", err)
	}

	// TODO(10000TB): Free mem hole start end aligns by
	// max(16, pagesize).
	//
	// Push alignment logic to kexec memory functions, e.g. a similar
	// function to FindSpace.
	setupRange, err = kmem.AddPhysSegment(
		linuxParam,
		RangeFromInterval(
			uintptr(0x90000),
			uintptr(len(linuxParam)),
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

	Debug("Loaded real mode data and cmdline at: %v", setupRange)

	/* Verify purgatory loads higher than the parameters. */
	// TODO(10000TB): if rel_addr < setupRange.Start then return error.

	/* Main kernel segment */
	var purgatoryEntry uintptr
	if purgatoryEntry, err = PurgeLoad(kmem, curPurgatory.Code, kernelEntry, setupRange.Start); err != nil {
		return err
	}
	Debug("purgatory entry: %v", purgatoryEntry)

	/* Load it */
	if err = Load(purgatoryEntry, kmem.Segments, 0); err != nil {
		return fmt.Errorf("kexec Load(%v, %v, %d) = %v", purgatoryEntry, kmem.Segments, 0, err)
	}

	return nil
}

// SelectPurgatory picks a purgatory, returning an error if none is found
func SelectPurgator(name string) error {
	p, ok := linux.Purgatories[name]
	if !ok {
		var s []string
		for i := range linux.Purgatories {
			s = append(s, i)
		}
		return fmt.Errorf("%s: no such purgatory, try one of %v", name, s)

	}
	curPurgatory = p
	return nil
}
