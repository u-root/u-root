// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"

	"github.com/u-root/u-root/pkg/boot/image"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

type kimage struct {
	segments kexec.Segments
	entry    uintptr
	cleanup  []func() error
}

func (k kimage) clean() {
	for _, c := range k.cleanup {
		if err := c(); err != nil {
			Debug("Failure: %v", err)
		}
	}
}

const (
	kernelAlignSize = 1 << 21 // 2 MB.
)

func mmap(f *os.File) (data []byte, ummap func() error, err error) {
	s, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("stat error: %w", err)
	}
	if s.Size() == 0 {
		return nil, nil, fmt.Errorf("cannot mmap zero-len file")
	}
	d, err := unix.Mmap(int(f.Fd()), 0, int(s.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, nil, fmt.Errorf("mmap failed: %w", err)
	}

	ummap = func() error {
		if err := unix.Munmap(d); err != nil {
			return fmt.Errorf("failed to unmap %s: %v", f.Name(), err)
		}
		return nil
	}
	return d, ummap, nil
}

// sanitizeFDT cleanups boot param properties from chosen node of the given FDT.
func sanitizeFDT(fdt *dt.FDT) (*dt.Node, error) {
	// Clear old entries in case we've already been through kexec to get
	// to this instance of runtime.
	chosen, _ := fdt.NodeByName("chosen")
	if chosen == nil {
		return nil, fmt.Errorf("no /chosen node in device tree")
	}
	for _, property := range []string{"linux,elfcorehdr", "linux,usable-memory-range", "kaslr-seed", "rng-seed", "linux,initrd-start", "linux,initrd-end"} {
		chosen.RemoveProperty(property)
	}

	return chosen, nil
}

func kexecLoadImage(kernel, ramfs *os.File, cmdline string, opts KexecOptions) (*kimage, error) {
	fdt, err := dt.LoadFDT(opts.DTB)
	if err != nil {
		return nil, fmt.Errorf("loadFDT(%s) = %v", opts.DTB, err)
	}
	Debug("Loaded FDT: %s", fdt)

	// Prepare segments.
	Debug("Try parsing memory map...")
	mm, err := kexec.MemoryMapFromFDT(fdt)
	if err != nil {
		return nil, fmt.Errorf("MemoryMapFromFDT(%v): %v", fdt, err)
	}
	Debug("Mem map: \n%+v", mm)
	return kexecLoadImageMM(mm, kernel, ramfs, fdt, cmdline, opts)
}

func kexecLoadImageMM(mm kexec.MemoryMap, kernel, ramfs *os.File, fdt *dt.FDT, cmdline string, opts KexecOptions) (*kimage, error) {
	kmem := &kexec.Memory{
		Phys: mm,
	}

	img := &kimage{}
	// Load kernel.
	var kernelBuf []byte
	var err error
	if opts.MmapKernel {
		Debug("Mmapping kernel to virtual buffer...")
		var cleanup func() error
		kernelBuf, cleanup, err = mmap(kernel)
		if err != nil {
			return nil, fmt.Errorf("mmap kernel: %v", err)
		}
		img.cleanup = append(img.cleanup, cleanup)
	} else {
		Debug("Read kernel from file ...")
		kernelBuf, err = uio.ReadAll(kernel)
		if err != nil {
			return nil, fmt.Errorf("read kernel from file: %v", err)
		}
	}

	kImage, err := image.ParseFromBytes(kernelBuf)
	if err != nil {
		return nil, fmt.Errorf("parse arm64 Image from bytes: %v", err)
	}

	kernelRange, err := kmem.AddKexecSegmentExplicit(kernelBuf, uint(kImage.Header.ImageSize), uint(kImage.Header.TextOffset), kernelAlignSize)
	if err != nil {
		return nil, fmt.Errorf("add kernel segment: %v", err)
	}

	Debug("Added %#x byte (size %#x) kernel at %s with offset %#x with alignment %#x", len(kernelBuf), kImage.Header.ImageSize, kernelRange, kImage.Header.TextOffset, kernelAlignSize)

	chosen, err := sanitizeFDT(fdt)
	if err != nil {
		return nil, fmt.Errorf("sanitizeFDT(%v) = %v", fdt, err)
	}
	Debug("FDT after sanitization: %s", fdt)

	var ramfsBuf []byte
	if ramfs != nil {
		if opts.MmapRamfs {
			Debug("Mmap ramfs file to virtual buffer...")
			var cleanup func() error
			ramfsBuf, cleanup, err = mmap(ramfs)
			if err != nil {
				return nil, fmt.Errorf("mmap ramfs: %v", err)
			}
			img.cleanup = append(img.cleanup, cleanup)
		} else {
			Debug("Read ramfs from file...")
			ramfsBuf, err = uio.ReadAll(ramfs)
			if err != nil {
				return nil, fmt.Errorf("read ramfs from file: %v", err)
			}
		}

		// NOTE(10000TB): This need be placed after kernel by convention.
		ramfsRange, err := kmem.AddKexecSegment(ramfsBuf)
		if err != nil {
			return nil, fmt.Errorf("add initramfs segment: %v", err)
		}
		Debug("Added %d byte initramfs at %s", len(ramfsBuf), ramfsRange)

		ramfsStart := make([]byte, 8)
		binary.BigEndian.PutUint64(ramfsStart, uint64(ramfsRange.Start))
		chosen.UpdateProperty("linux,initrd-start", ramfsStart)
		ramfsEnd := make([]byte, 8)
		binary.BigEndian.PutUint64(ramfsEnd, uint64(ramfsRange.Start)+uint64(ramfsRange.Size))
		chosen.UpdateProperty("linux,initrd-end", ramfsEnd)
	}

	Debug("Kernel cmdline to append: %s", cmdline)
	if len(cmdline) > 0 {
		cmdlineBuf := append([]byte(cmdline), byte(0))
		chosen.UpdateProperty("bootargs", cmdlineBuf)
	} else {
		chosen.RemoveProperty("bootargs")
	}

	var dtbBuffer bytes.Buffer
	if _, err := fdt.Write(&dtbBuffer); err != nil {
		return nil, fmt.Errorf("flattening device tree: %v", err)
	}
	dtbBuf := dtbBuffer.Bytes()
	dtbRange, err := kmem.AddKexecSegment(dtbBuf)
	if err != nil {
		return nil, fmt.Errorf("add device tree segment: %w", err)
	}
	Debug("Added %d byte device tree at %s", len(dtbBuf), dtbRange)

	// Trampoline.
	//
	// We need a trampoline to pass the DTB to the kernel; because
	// we'll use this code as our entry point, it also needs to know
	// the real entry point to kernel.
	//
	// TODO(10000TB): this assumes a little endian kernel, support
	// big endian if needed per flag.
	kernelEntry := kernelRange.Start
	dtbBase := dtbRange.Start

	var trampoline [10]uint32
	// Instruction encoding per
	// "Arm Architecture Reference Manual Armv8, for Armv8-A architecture
	// profile" [ ARM DDI 0487E.a (ID070919) ]
	trampoline[0] = 0x580000c4 // ldr x4, #0x18 (PC relative: trampoline[6 and 7])
	trampoline[1] = 0x580000e0 // ldr x0, #0x1c (PC relative: trampoline[8 and 9])
	// Zero out x1, x2, x3
	trampoline[2] = 0xaa1f03e1 // mov x1, xzr
	trampoline[3] = 0xaa1f03e2 // mov x2, xzr
	trampoline[4] = 0xaa1f03e3 // mov x3, xzr
	// Branch register / Jump to instruction from x4.
	trampoline[5] = 0xd61f0080 // br  x4

	trampoline[6] = uint32(uint64(kernelEntry) & 0xffffffff)
	trampoline[7] = uint32(uint64(kernelEntry) >> 32)
	trampoline[8] = uint32(uint64(dtbBase) & 0xffffffff)
	trampoline[9] = uint32(uint64(dtbBase) >> 32)

	var trampolineBuffer bytes.Buffer
	if err := binary.Write(&trampolineBuffer, binary.LittleEndian, trampoline); err != nil {
		return nil, fmt.Errorf("make trampoline: %v", err)
	}
	Debug("trampoline bytes %x", trampolineBuffer.Bytes())
	trampolineRange, err := kmem.AddKexecSegment(trampolineBuffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("add trampoline segment: %v", err)
	}
	Debug("Added %d byte trampoline at %s", len(trampolineBuffer.Bytes()), trampolineRange)

	/* Load it */
	img.entry = trampolineRange.Start
	img.segments = kmem.Segments
	Debug("Entry: %#x", img.entry)
	return img, nil
}
