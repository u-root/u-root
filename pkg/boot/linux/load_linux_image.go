// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/u-root/u-root/pkg/boot/image"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
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

var errNoChosenNode = fmt.Errorf("no /chosen node in device tree")

// sanitizeFDT cleanups boot param properties from chosen node of the given FDT.
func sanitizeFDT(fdt *dt.FDT) (*dt.Node, error) {
	// Clear old entries in case we've already been through kexec to get
	// to this instance of runtime.
	chosen, _ := fdt.NodeByName("chosen")
	if chosen == nil {
		return nil, errNoChosenNode
	}
	for _, property := range []string{"linux,elfcorehdr", "linux,usable-memory-range", "kaslr-seed", "rng-seed", "linux,initrd-start", "linux,initrd-end"} {
		chosen.RemoveProperty(property)
	}

	return chosen, nil
}

var ErrMemmapEmpty = errors.New("memory map is empty or contains no information about system RAM")

func kexecLoadImage(kernel, ramfs *os.File, cmdline string, dtb io.ReaderAt, reservedRanges kexec.Ranges) (*kimage, error) {
	var fdt *dt.FDT
	var err error
	// We want to fail when a user-supplied FDT is not parseable, not
	// implicitly fall back to some other FDT. Avoid the dt.LoadFDT API.
	if dtb != nil {
		fdt, err = dt.ReadFDT(io.NewSectionReader(dtb, 0, math.MaxInt64))
	} else {
		fdt, err = dt.ReadFile("/sys/firmware/fdt")
	}
	if err != nil {
		return nil, fmt.Errorf("read FDT = %w", err)
	}
	Debug("Loaded FDT: %s", fdt)

	// Prepare segments.
	Debug("Try parsing memory map...")
	mm, err := kexec.MemoryMapFromFDT(fdt)
	if err != nil {
		return nil, fmt.Errorf("memoryMapFromFDT(%v): %w", fdt, err)
	}
	Debug("Mem map: \n%+v", mm)
	if len(mm.RAM()) == 0 {
		return nil, ErrMemmapEmpty
	}
	for _, r := range reservedRanges {
		mm.Insert(kexec.TypedRange{Range: r, Type: kexec.RangeReserved})
	}
	return kexecLoadImageMM(mm, kernel, ramfs, fdt, cmdline)
}

var (
	errKernelSegmentFailed     = errors.New("failed to add kernel segment")
	errInitramfsSegmentFailed  = errors.New("failed to add initramfs segment")
	errDTBSegmentFailed        = errors.New("failed to add DTB segment")
	errTrampolineSegmentFailed = errors.New("failed to add trampolineSegment")
)

func kexecLoadImageMM(mm kexec.MemoryMap, kernel, ramfs *os.File, fdt *dt.FDT, cmdline string) (*kimage, error) {
	kmem := &kexec.Memory{
		Phys: mm,
	}

	img := &kimage{}

	// Load kernel.
	kernelBuf, cleanup, err := getFile(kernel)
	if err != nil {
		return nil, fmt.Errorf("failed to get kernel contents: %w", err)
	}
	img.cleanup = append(img.cleanup, cleanup)

	kImage, err := image.ParseFromBytes(kernelBuf)
	if err != nil {
		return nil, fmt.Errorf("parse arm64 Image from bytes: %w", err)
	}

	// "The Image must be placed text_offset bytes from a 2MB aligned base
	// address anywhere in usable system RAM and called there."
	// (arm64/booting.rst)
	//
	// "At least image_size bytes from the start of the image must be free
	// for use by the kernel." (arm64/booting.rst)
	//
	// TODO: support versions below v4.6?
	//
	// "NOTE: versions prior to v4.6 cannot make use of memory below the
	// physical offset of the Image so it is recommended that the Image be
	// placed as close as possible to the start of system RAM."
	// (arm64/booting.rst)
	kernelRange, err := kmem.AddKexecSegmentExplicit(kernelBuf, uint(kImage.Header.ImageSize), uint(kImage.Header.TextOffset), kernelAlignSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errKernelSegmentFailed, err)
	}

	Debug("Added %#x byte (size %#x) kernel at %s with offset %#x with alignment %#x", len(kernelBuf), kImage.Header.ImageSize, kernelRange, kImage.Header.TextOffset, kernelAlignSize)

	chosen, err := sanitizeFDT(fdt)
	if err != nil {
		return nil, fmt.Errorf("sanitizeFDT(%v) = %w", fdt, err)
	}
	Debug("FDT after sanitization: %s", fdt)

	if ramfs != nil {
		ramfsBuf, cleanup, err := getFile(ramfs)
		if err != nil {
			return nil, fmt.Errorf("failed to get initramfs contents: %w", err)
		}
		img.cleanup = append(img.cleanup, cleanup)

		// NOTE(10000TB): This need be placed after kernel by convention.
		//
		// "If an initrd/initramfs is passed to the kernel at boot, it
		// must reside entirely within a 1 GB aligned physical memory
		// window of up to 32 GB in size that fully covers the kernel
		// Image as well." (arm64/booting.rst)
		ramfsRange, err := kmem.AddKexecSegment(ramfsBuf)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errInitramfsSegmentFailed, err)
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
		return nil, fmt.Errorf("flattening device tree: %w", err)
	}
	dtbBuf := dtbBuffer.Bytes()
	dtbRange, err := kmem.AddKexecSegment(dtbBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDTBSegmentFailed, err)
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
		return nil, fmt.Errorf("make trampoline: %w", err)
	}
	Debug("trampoline bytes %x", trampolineBuffer.Bytes())
	trampolineRange, err := kmem.AddKexecSegment(trampolineBuffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errTrampolineSegmentFailed, err)
	}
	Debug("Added %d byte trampoline at %s", len(trampolineBuffer.Bytes()), trampolineRange)

	/* Load it */
	img.entry = trampolineRange.Start
	img.segments = kmem.Segments
	Debug("Entry: %#x", img.entry)
	return img, nil
}
