// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/u-root/u-root/pkg/boot/image"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

func kexecLoadRISCVImage(kernel, ramfs *os.File, cmdline string, dtb io.ReaderAt, reservedRanges kexec.Ranges) (*kimage, error) {
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
	return kexecLoadRISCVImageMM(mm, kernel, ramfs, fdt, cmdline)
}

func kexecLoadRISCVImageMM(mm kexec.MemoryMap, kernel, ramfs *os.File, fdt *dt.FDT, cmdline string) (*kimage, error) {
	kmem := &kexec.Memory{
		Phys: mm,
	}

	img := &kimage{}

	kernelBuf, cleanup, err := getFile(kernel)
	if err != nil {
		return nil, fmt.Errorf("failed to get kernel contents: %w", err)
	}
	img.cleanup = append(img.cleanup, cleanup)

	kImage, err := image.ParseRISCVFromBytes(kernelBuf)
	if err != nil {
		return nil, fmt.Errorf("parse riscv64 Image from bytes: %w", err)
	}

	// The Image is loaded text_offset bytes from a 2MB aligned base address.
	// image_size is the memory the kernel needs, which is larger than the
	// file itself as it covers the BSS.
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

	// Unlike arm64, no trampoline is needed. The riscv64 kernel expects the
	// hartid in a0 and the device tree address in a1, and neither can be
	// set from here: the hartid is only known at run time. Instead the
	// currently running kernel does that work in machine_kexec(), which
	// locates the FDT by scanning the loaded segments for a valid device
	// tree header and reads the hartid via cpuid_to_hartid_map(), then its
	// relocation code sets a0 and a1 before jumping. So the kernel itself
	// is the entry point.
	img.entry = kernelRange.Start
	img.segments = kmem.Segments
	Debug("Entry: %#x", img.entry)
	return img, nil
}
