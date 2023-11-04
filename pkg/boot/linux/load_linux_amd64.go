// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/boot/bzimage"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/purgatory"
	"github.com/u-root/u-root/pkg/uio"
)

const (
	bootParams = "/sys/kernel/boot_params/data"
)

// KexecLoad loads a bzImage-formated Linux kernel file as the to-be-kexeced
// kernel with the given ramfs file and cmdline string.
//
// It uses the kexec_load system call.
func KexecLoad(kernel, ramfs *os.File, cmdline string, opts KexecOptions) error {
	bzimage.Debug = Debug

	// A collection of vars used for processing the kernel for kexec
	var err error
	// bzimage is the deserialized bzImage from the kernel
	// io.ReaderAt.
	var bzimg bzimage.BzImage
	// kmem is a struct holding kexec segments.
	//
	// It has routines to work with physical memory
	// ranges.
	var kmem *kexec.Memory
	// TODO(10000TB): construct default params in go.
	//
	// boot_params directory is x86 specific. So for now, following code only
	// works on x86.
	// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-kernel-boot_params
	bp, err := os.ReadFile("/sys/kernel/boot_params/data")
	if err != nil {
		return fmt.Errorf("reading boot_param data: %w", err)
	}
	var lp = &bzimage.LinuxParams{}
	if err := lp.UnmarshalBinary(bp); err != nil {
		return fmt.Errorf("unmarshaling header: %w", err)
	}

	kb, err := uio.ReadAll(kernel)
	if err != nil {
		return fmt.Errorf("reading Linux kernel into memory: %w", err)
	}
	if err := bzimg.UnmarshalBinary(kb); err != nil {
		return fmt.Errorf("parsing bzImage Linux kernel: %w", err)
	}

	if len(bzimg.KernelCode) < 1024 {
		return fmt.Errorf("kernel code size smaller than 1024 bytes: %d", len(bzimg.KernelCode))
	}

	kelf, err := bzimg.ELF()
	if err != nil {
		return fmt.Errorf("getting ELF from bzImage: %w", err)
	}
	kernelEntry := uintptr(kelf.Entry)
	Debug("kernelEntry: %v", kernelEntry)

	// Prepare segments.
	kmem = &kexec.Memory{}
	Debug("Try parsing memory map...")
	// TODO(10000TB): refactor this call into initialization of
	// kexec.Memory, as it does not depend on specific boot.
	if err := kmem.ParseMemoryMap(); err != nil {
		return fmt.Errorf("parse memory map: %v", err)
	}

	var relocatableKernel bool
	if bzimg.Header.Protocolversion < 0x0205 {
		return fmt.Errorf("bzImage boot protocol earlier thatn 2.05 is not supported currently: %v", bzimg.Header.Protocolversion)
	}
	relocatableKernel = bzimg.Header.RelocatableKernel != 0
	// Only protected mode is currently supported.
	// In protected mode, kernel need be relocatable, or it will need to fall
	// to real mode executing.
	if !relocatableKernel {
		return errors.New("non-relocateable Kernels are not supported")
	}
	if _, err := kmem.LoadElfSegments(bytes.NewReader(bzimg.KernelCode)); err != nil {
		return fmt.Errorf("loading kernel ELF segments: %w", err)
	}

	var ramfsRange kexec.Range
	if ramfs != nil {
		ramfsContents, err := io.ReadAll(ramfs)
		if err != nil {
			return fmt.Errorf("unable to read initramfs: %w", err)
		}
		if ramfsRange, err = kmem.AddKexecSegment(ramfsContents); err != nil {
			return fmt.Errorf("add initramfs segment: %w", err)
		}
		Debug("Added %d byte initramfs at %s", len(ramfsContents), ramfsRange)
		lp.Initrdstart = uint32(ramfsRange.Start)
		lp.Initrdsize = uint32(ramfsRange.Size)
	}

	Debug("Kernel cmdline to append: %s", cmdline)
	if len(cmdline) > 0 {
		var cmdlineRange kexec.Range
		Debug("Add cmdline: %s", cmdline)

		// Cmdline must be null-terminated.
		cmdlineBytes := []byte(cmdline + "\x00")
		if cmdlineRange, err = kmem.AddKexecSegment(cmdlineBytes); err != nil {
			return fmt.Errorf("add cmdline segment: %v", err)
		}
		Debug("Added %d byte of cmdline at %s", len(cmdlineBytes), cmdlineRange)
		lp.CLPtr = uint32(cmdlineRange.Start)      // 2.02+
		lp.CmdLineSize = uint32(cmdlineRange.Size) // 2.06+
	}

	// The kernel is a bzImage kernel if the protocol >= 2.00 and the 0x01
	// bit (LOAD_HIGH) in the loadflags field is set.
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
	setupRange, err := kmem.AddPhysSegment(
		linuxParam,
		kexec.RangeFromInterval(
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

	// Verify purgatory loads higher than the parameters.
	// TODO(10000TB): if rel_addr < setupRange.Start then return error.

	// Load purgatory.
	purgatoryEntry, err := purgatory.Load(kmem, kernelEntry, setupRange.Start)
	if err != nil {
		return err
	}
	Debug("purgatory entry: %v", purgatoryEntry)

	// Load it.
	if err := kexec.Load(purgatoryEntry, kmem.Segments, 0); err != nil {
		return fmt.Errorf("kexec load(%v, %v, %d): %w", purgatoryEntry, kmem.Segments, 0, err)
	}
	return nil
}
