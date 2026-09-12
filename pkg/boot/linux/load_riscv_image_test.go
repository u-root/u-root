// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"errors"

	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

func riscvMemFDT() *dt.FDT {
	return &dt.FDT{
		RootNode: dt.NewNode("/", dt.WithChildren(
			dt.NewNode("chosen"),
			dt.NewNode("test memory", dt.WithProperty(
				dt.PropertyString("device_type", "memory"),
				dt.PropertyRegion("reg", 0x100000, 0xf00000),
			)),
		)),
	}
}

func riscvMM(t *testing.T) kexec.MemoryMap {
	t.Helper()
	mm, err := kexec.MemoryMapFromFDT(riscvMemFDT())
	if err != nil {
		t.Fatal(err)
	}
	return mm
}

// The riscv64 kernel is entered directly rather than through a trampoline: the
// running kernel finds the FDT among the segments and sets a0/a1 itself, so
// there is nothing for us to place in front of it.
func TestKexecLoadRISCVImageEntryIsKernel(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	img, err := kexecLoadRISCVImageMM(riscvMM(t), kernel, nil, riscvMemFDT(), "")
	if err != nil {
		t.Fatalf("kexecLoadRISCVImageMM = %v, want nil", err)
	}
	defer img.clean()

	var kernelStart uintptr
	for _, s := range img.segments {
		if s.Phys.Size >= 0xb429b8 {
			kernelStart = s.Phys.Start
		}
	}
	if kernelStart == 0 {
		t.Fatalf("no kernel-sized segment in %v", img.segments)
	}
	if img.entry != kernelStart {
		t.Errorf("entry = %#x, want kernel start %#x", img.entry, kernelStart)
	}

	// The kernel must land at text_offset from a 2MB aligned base.
	if got := img.entry % kernelAlignSize; got != 0x200000%kernelAlignSize {
		t.Errorf("entry %#x is not correctly aligned (mod %#x = %#x)", img.entry, kernelAlignSize, got)
	}
}

// image_size covers the BSS, so the reserved memory must follow the header's
// image_size rather than the size of the file on disk.
func TestKexecLoadRISCVImageReservesImageSize(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	fi, err := kernel.Stat()
	if err != nil {
		t.Fatal(err)
	}

	img, err := kexecLoadRISCVImageMM(riscvMM(t), kernel, nil, riscvMemFDT(), "")
	if err != nil {
		t.Fatalf("kexecLoadRISCVImageMM = %v, want nil", err)
	}
	defer img.clean()

	var found bool
	for _, s := range img.segments {
		if s.Phys.Start == img.entry {
			found = true
			if uint(s.Phys.Size) < 0xb429b8 {
				t.Errorf("kernel segment size %#x, want at least image_size %#x", s.Phys.Size, 0xb429b8)
			}
			if int64(s.Phys.Size) <= fi.Size() {
				t.Errorf("kernel segment size %#x should exceed file size %#x", s.Phys.Size, fi.Size())
			}
		}
	}
	if !found {
		t.Error("no segment at entry")
	}
}

func TestKexecLoadRISCVImageRejectsArm64(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/Image")
	defer kernel.Close()

	if _, err := kexecLoadRISCVImageMM(riscvMM(t), kernel, nil, riscvMemFDT(), ""); err == nil {
		t.Error("kexecLoadRISCVImageMM(arm64 Image) = nil, want error")
	}
}

func TestKexecLoadRISCVImageCmdlineAndInitramfs(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()
	ramfs := createFile(t, []byte("ramfs"))
	defer ramfs.Close()

	fdt := riscvMemFDT()
	img, err := kexecLoadRISCVImageMM(riscvMM(t), kernel, ramfs, fdt, "console=ttyS0")
	if err != nil {
		t.Fatalf("kexecLoadRISCVImageMM = %v, want nil", err)
	}
	defer img.clean()

	chosen, _ := fdt.NodeByName("chosen")
	if chosen == nil {
		t.Fatal("no chosen node")
	}
	for _, want := range []string{"bootargs", "linux,initrd-start", "linux,initrd-end"} {
		if _, found := chosen.LookProperty(want); !found {
			t.Errorf("chosen is missing %q", want)
		}
	}
}

func TestKexecLoadRISCVImageNoMemory(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	_, err := kexecLoadRISCVImageMM(kexec.MemoryMap{}, kernel, nil, riscvMemFDT(), "")
	if err == nil {
		t.Error("kexecLoadRISCVImageMM(empty mm) = nil, want error")
	}
	if !errors.Is(err, errKernelSegmentFailed) {
		t.Errorf("err = %v, want %v", err, errKernelSegmentFailed)
	}
}

// Exercise the outer entry point, which reads and parses the FDT before
// handing off. Passing a dtb keeps it off /sys/firmware/fdt.
func TestKexecLoadRISCVImageWithDTB(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	img, err := kexecLoadRISCVImage(kernel, nil, "console=ttyS0", fdtReader(t, riscvMemFDT()), nil)
	if err != nil {
		t.Fatalf("kexecLoadRISCVImage = %v, want nil", err)
	}
	defer img.clean()

	if img.entry == 0 {
		t.Error("entry is 0")
	}
	if len(img.segments) == 0 {
		t.Error("no segments")
	}
}

func TestKexecLoadRISCVImageBadDTB(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	if _, err := kexecLoadRISCVImage(kernel, nil, "", bytes.NewReader([]byte("not an fdt")), nil); err == nil {
		t.Error("kexecLoadRISCVImage(bad dtb) = nil, want error")
	}
}

// Reserved ranges must be kept away from the segments we allocate.
func TestKexecLoadRISCVImageReservations(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../image/testdata/riscv64_header.bin")
	defer kernel.Close()

	// The FDT needs enough RAM to fit the kernel's image_size alongside the
	// reserved hole.
	fdt := &dt.FDT{
		RootNode: dt.NewNode("/", dt.WithChildren(
			dt.NewNode("chosen"),
			dt.NewNode("test memory", dt.WithProperty(
				dt.PropertyString("device_type", "memory"),
				dt.PropertyRegion("reg", 0x100000, 0x8000000),
			)),
		)),
	}
	reserved := kexec.Ranges{{Start: 0x100000, Size: 0x400000}}
	img, err := kexecLoadRISCVImage(kernel, nil, "", fdtReader(t, fdt), reserved)
	if err != nil {
		t.Fatalf("kexecLoadRISCVImage = %v, want nil", err)
	}
	defer img.clean()

	for _, s := range img.segments {
		for _, r := range reserved {
			if s.Phys.Start < r.Start+uintptr(r.Size) && r.Start < s.Phys.Start+uintptr(s.Phys.Size) {
				t.Errorf("segment %v overlaps reserved range %v", s.Phys, r)
			}
		}
	}
}
