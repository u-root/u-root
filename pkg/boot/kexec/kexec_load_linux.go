// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/bzimage"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_INITRD_ADDR_MAX  uint = 0x37FFFFFF
	DEFAULT_BZIMAGE_ADDR_MAX uint = 0x37FFFFFF
	bootParams                    = "/sys/kernel/boot_params/data"
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
func KexecLoad(kernel io.ReaderAt, ramfs io.Reader, cmdline string) error {
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

	bp, err := ioutil.ReadFile("/sys/kernel/boot_params/data")
	if err != nil {
		return fmt.Errorf("Reading boot_param data: %w", err)
	}
	var lp = &bzimage.LinuxParams{}
	if err := lp.UnmarshalBinary(bp); err != nil {
		return fmt.Errorf("Unmarshaling header: %w", err)
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
			return fmt.Errorf("Adding initramfs segment: %v", err)
		}
		Debug("Added %d byte initramfs at %s", len(b), ramfsRange)
		lp.Initrdstart = uint32(ramfsRange.Start)
		lp.Initrdsize = uint32(ramfsRange.Size)
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

	var setupRange Range
	// The kernel is a bzImage kernel if the protocol >= 2.00 and the 0x01
	// bit (LOAD_HIGH) in the loadflags field is set
	// TODO(10000TB): check on loadflags.
	linuxParam, err := lp.MarshalBinary()
	if err != nil {
		return fmt.Errorf("Re-marshaling header: %w", err)
	}

	// TODO(10000TB): Free mem hole start end aligns by
	// max(16, pagesize).
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

	/* Main kernel segment.
	 */
	var purgatoryEntry uintptr
	if purgatoryEntry, err = PurgeLoad(kmem, curPurgatory.code, kernelEntry, setupRange.Start); err != nil {
		return err
	}
	Debug("purgatory entry: %v", purgatoryEntry)

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
