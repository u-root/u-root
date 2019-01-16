// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/u-root/u-root/pkg/dt"
	"github.com/u-root/u-root/pkg/zimage"
	"golang.org/x/sys/unix"
)

const (
	// 1024 with the null character
	maxCommandLineSize = 1023

	// These values are hard-coded into the kernel.
	stackSize  = 0x1000
	heapSize   = 0x10000
	textOffset = 0x8000
)

// ZImageLoad loads a zImage formatted Linux kernel for kexec. If DTB is nil,
// it defaults to /sys/firmware/fdt. Uses the kexec_load(2) syscall.
func ZImageLoad(opts *LinuxOpts) error {
	// Physical layout:
	//     1. Empty memory (0x8000 bytes). This puts the kernel's entry point
	//        at the expected 0x8000.
	//     2. zImage. The entry point is at offset 0x0.
	//     3. Empty memory for the zImage's BSS (memory used by the
	//        decompressor). The size of this can be found in the ZImage's
	//        header.
	//     4. Empty memory to decompress the kernel into. The size can be found
	//        in the zImage.
	//     5. Initramfs (page aligned).
	//     6. Empty page. Prevents the Initramfs and DT from combining into one
	//        segment.
	//     7. Device tree (page aligned). Must be in its own segment.
	// The initramfs (#5) and DT (#7) are their own segments. The zImage (#2)
	// segment is padded with #1, #2 and #4.

	// Detect image type.
	zImage, err := zimage.Parse(opts.Kernel)
	if err != nil {
		return &ProbeError{err}
	}
	if zImage.Header.Start != 0 {
		return errors.New("zImage must be position independent")
	}
	edataSize, kernelBSSSize, err := zImage.GetKernelSizes(opts.Kernel)
	if err != nil {
		return fmt.Errorf("kernel does not support size extension: %v", err)
	}

	// Read kernel into RAM.
	if _, err := opts.Kernel.Seek(0, io.SeekStart); err != nil {
		return err
	}
	zImageData, err := ioutil.ReadAll(opts.Kernel)
	if err != nil {
		return fmt.Errorf("error reading kernel: %v", err)
	}
	if zImage.Header.End > uint32(len(zImageData)) {
		return fmt.Errorf("zImage kernel is truncated: %d > %d",
			zImage.Header.End, len(zImageData))
	}
	zImageData = zImageData[:zImage.Header.End] // Remove trailing data

	// Read and parse the device tree.
	if opts.DTB == nil {
		f, err := os.Open("/sys/firmware/fdt")
		if err != nil {
			return fmt.Errorf("error opening /sys/firmware/fdt: %v", err)
		}
		defer f.Close()
		opts.DTB = f
	}
	fdt, err := dt.ReadFDT(opts.DTB)
	if err != nil {
		return fmt.Errorf("error parsing fdt: %v", err)
	}
	mem := Memory{}
	if err := mem.ParseFromDeviceTree(fdt); err != nil {
		return err
	}
	chosen := fdt.RootNode.CreateNode("chosen")

	// Validate command line.
	if len(opts.CmdLine) > maxCommandLineSize {
		return fmt.Errorf("cmdline is too long, %d > %d bytes",
			len(opts.CmdLine), maxCommandLineSize)
	}
	chosen.CreateProperty("bootargs").SetString(opts.CmdLine)

	// Create segment for kernel.
	// TODO: Padding can be avoided by creating a "hole" function which can take
	//       a minimum offset + physical size constraint.
	kernelData := append(append(
		make([]byte, textOffset),
		zImageData...),
		make([]byte, min(stackSize+heapSize, edataSize+kernelBSSSize))...)
	// TODO: As a optimization, prefer to place the kernel above 32MiB to avoid
	//       relocating the kernel before decompression.
	kernelStart, err := mem.AddKexecSegment(kernelData)

	// Read initramfs into RAM (optional).
	initramfsStart := kernelStart + uintptr(alignUp(uint(len(kernelData))))
	var initramfsSize uint
	if opts.Initramfs != nil {
		initramfs, err := ioutil.ReadAll(opts.Initramfs)
		if err != nil {
			return fmt.Errorf("error reading initramfs: %v", err)
		}

		if err := mem.AddKexecSegmentPhys(initramfs, initramfsStart); err != nil {
			return err
		}
		initramfsSize = uint(len(initramfs))

		// TODO: Check #cell-size before doing uint32 casts.
		chosen.CreateProperty("linux,initrd-start").SetU32(uint32(initramfsStart))
		chosen.CreateProperty("linux,initrd-end").SetU32(
			uint32(initramfsStart) + uint32(initramfsSize))
	} else {
		// The FDT might already contain an initrd pointer from the previous
		// boot. Delete because it has likely been repurposed by the kernel.
		chosen.DeleteProperty("linux,initrd-start")
		chosen.DeleteProperty("linux,initrd-end")
	}

	// Create segment for DTB.
	dtb := &bytes.Buffer{}
	if _, err := fdt.Write(dtb); err != nil {
		return err
	}
	// Include a buffer page before the DTB.
	dtbStart := initramfsStart + uintptr(alignUp(alignUp(uint(initramfsSize))+1))
	fmt.Printf("%#x\n", dtbStart)
	if err = mem.AddKexecSegmentPhys(dtb.Bytes(), dtbStart); err != nil {
		return err
	}

	if opts.LoadSyscall == nil {
		opts.LoadSyscall = RawLoad
	}
	return opts.LoadSyscall(kernelStart+textOffset, mem.Segments, unix.KEXEC_ARCH_DEFAULT)
}

func min(x uint32, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}
