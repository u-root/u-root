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

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/zimage"
	"github.com/u-root/u-root/pkg/dt"
)

const (
	// zImageTextOffset is where the kernel is placed relative to the start
	// of the region found for it. (arch/arm/boot/compressed/vmlinux.lds.S
	// exports it as TEXT_OFFSET.)
	zImageTextOffset = 0x8000

	// zImageExtraSpace is the stack (4K) and malloc space (64K) the zImage
	// uses while decompressing, which its own length does not cover.
	zImageExtraSpace = 0x11000

	// zImageAlignSize is the alignment for the region the kernel is placed
	// in. A page is enough; the decompressor relocates itself as needed.
	zImageAlignSize = 0x1000
)

// zImageMemSize returns the amount of memory to reserve for a zImage.
//
// A zImage is self-extracting and its header does not say how large the kernel
// becomes, so the space that must be kept free is derived from the kernel size
// entry in its extension table: the decompressed kernel is edata_size +
// bss_size, and while decompressing, the compressed image sits past _edata of
// the decompressed kernel, so at least edata_size + the compressed length is
// needed as well.
//
// If the entry is absent -- it is optional, and older kernels do not carry one
// -- fall back to five times the compressed length, as kexec-tools does,
// assuming a compression ratio of at most four plus room for the image itself.
func zImageMemSize(kernelBuf []byte) (uint, error) {
	z, err := zimage.Parse(bytes.NewReader(kernelBuf))
	if err != nil {
		return 0, fmt.Errorf("parse zImage: %w", err)
	}

	// The zImage needs its stack and malloc space on top of its own length.
	compressed := uint(len(kernelBuf)) + zImageExtraSpace

	sizeAddr, bssSize, err := z.GetKernelSizes()
	if err != nil {
		Debug("zImage has no usable kernel size entry (%v), assuming a compression ratio of 4", err)
		return uint(len(kernelBuf)) * 5, nil
	}

	// sizeAddr is an offset into the image at which the size of the
	// decompressed kernel up to _edata is stored.
	if uint64(sizeAddr)+4 > uint64(len(kernelBuf)) {
		return 0, fmt.Errorf("zImage kernel size is at %#x, past the end of the %#x byte image", sizeAddr, len(kernelBuf))
	}
	edataSize := uint(binary.LittleEndian.Uint32(kernelBuf[sizeAddr:]))

	size := max(edataSize+uint(bssSize), edataSize+compressed)
	Debug("zImage decompresses to %#x bytes (text+data %#x, bss %#x), reserving %#x", edataSize+uint(bssSize), edataSize, bssSize, size)
	return size, nil
}

func kexecLoadZImage(kernel, ramfs *os.File, cmdline string, dtb io.ReaderAt, reservedRanges kexec.Ranges) (*kimage, error) {
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
	return kexecLoadZImageMM(mm, kernel, ramfs, fdt, cmdline)
}

func kexecLoadZImageMM(mm kexec.MemoryMap, kernel, ramfs *os.File, fdt *dt.FDT, cmdline string) (*kimage, error) {
	kmem := &kexec.Memory{
		Phys: mm,
	}

	img := &kimage{}

	kernelBuf, cleanup, err := getFile(kernel)
	if err != nil {
		return nil, fmt.Errorf("failed to get kernel contents: %w", err)
	}
	img.cleanup = append(img.cleanup, cleanup)

	memSize, err := zImageMemSize(kernelBuf)
	if err != nil {
		return nil, err
	}

	// The kernel is placed TEXT_OFFSET into the region, and memSize must
	// stay free so that the decompressed kernel and its BSS do not land on
	// the initrd or the device tree.
	kernelRange, err := kmem.AddKexecSegmentExplicit(kernelBuf, memSize, zImageTextOffset, zImageAlignSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errKernelSegmentFailed, err)
	}
	Debug("Added %#x byte (size %#x) kernel at %s with offset %#x", len(kernelBuf), memSize, kernelRange, zImageTextOffset)

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

	// As on riscv64, no trampoline is needed. The arm kernel expects r0=0,
	// r1=machine type and r2=the device tree, and none of those can be set
	// from here. The running kernel does it instead: machine_kexec_prepare()
	// scans the loaded segments for a device tree header to find r2, and
	// relocate_kernel.S sets all three before jumping. So the kernel itself
	// is the entry point, as it also is in kexec-tools.
	img.entry = kernelRange.Start
	img.segments = kmem.Segments
	Debug("Entry: %#x", img.entry)
	return img, nil
}
