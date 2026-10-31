// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

func zImageFDT() *dt.FDT {
	return &dt.FDT{
		RootNode: dt.NewNode("/", dt.WithChildren(
			dt.NewNode("chosen"),
			dt.NewNode("test memory", dt.WithProperty(
				dt.PropertyString("device_type", "memory"),
				dt.PropertyRegion("reg", 0x80000000, 0x10000000),
			)),
		)),
	}
}

func zImageMM(t *testing.T) kexec.MemoryMap {
	t.Helper()
	mm, err := kexec.MemoryMapFromFDT(zImageFDT())
	if err != nil {
		t.Fatal(err)
	}
	return mm
}

// A zImage is self-extracting, so the space reserved for it has to cover the
// decompressed kernel and its BSS rather than the size of the file.
func TestZImageMemSize(t *testing.T) {
	Debug = t.Logf

	kernelBuf, err := os.ReadFile("../zimage/testdata/zImage")
	if err != nil {
		t.Fatal(err)
	}

	got, err := zImageMemSize(kernelBuf)
	if err != nil {
		t.Fatalf("zImageMemSize = %v, want nil", err)
	}
	if got <= uint(len(kernelBuf)) {
		t.Errorf("reserved %#x for a %#x byte image, want more than the file itself", got, len(kernelBuf))
	}
	// The compressed image sits past _edata while unpacking, so at minimum
	// _edata plus the image and its scratch space must be free.
	if want := uint(len(kernelBuf)) + zImageExtraSpace; got < want {
		t.Errorf("reserved %#x, want at least %#x", got, want)
	}
}

// Without a kernel size entry there is nothing to compute from, so fall back
// to a guess rather than refusing the image.
func TestZImageMemSizeNoEntry(t *testing.T) {
	Debug = t.Logf

	kernelBuf, err := os.ReadFile("../zimage/testdata/zImage")
	if err != nil {
		t.Fatal(err)
	}
	// Break the table magic so the entry is not found.
	buf := append([]byte(nil), kernelBuf...)
	binary.LittleEndian.PutUint32(buf[0x24+16:], 0)

	got, err := zImageMemSize(buf)
	if err != nil {
		t.Fatalf("zImageMemSize = %v, want nil", err)
	}
	if want := uint(len(buf)) * 5; got != want {
		t.Errorf("fallback reserved %#x, want %#x", got, want)
	}
}

func TestZImageMemSizeNotAZImage(t *testing.T) {
	if _, err := zImageMemSize([]byte("this is not a zImage")); err == nil {
		t.Error("zImageMemSize(garbage) = nil, want error")
	}
}

// The kernel is entered directly: r0/r1/r2 are set by the running kernel, so
// there is no trampoline to place in front of it.
func TestKexecLoadZImageEntryIsKernel(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	img, err := kexecLoadZImageMM(zImageMM(t), kernel, nil, zImageFDT(), "")
	if err != nil {
		t.Fatalf("kexecLoadZImageMM = %v, want nil", err)
	}
	defer img.clean()

	var kernelSeg *kexec.Segment
	for i := range img.segments {
		if img.segments[i].Phys.Start == img.entry {
			kernelSeg = &img.segments[i]
		}
	}
	if kernelSeg == nil {
		t.Fatalf("no segment starts at entry %#x: %v", img.entry, img.segments)
	}

	// The kernel goes TEXT_OFFSET into an aligned region.
	if got := (img.entry - zImageTextOffset) % zImageAlignSize; got != 0 {
		t.Errorf("entry %#x is not TEXT_OFFSET into an aligned region (remainder %#x)", img.entry, got)
	}
}

// The reserved space must exceed the file, or the decompressed kernel would
// land on whatever was placed after it.
func TestKexecLoadZImageReservesDecompressedSize(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	fi, err := kernel.Stat()
	if err != nil {
		t.Fatal(err)
	}

	img, err := kexecLoadZImageMM(zImageMM(t), kernel, nil, zImageFDT(), "")
	if err != nil {
		t.Fatalf("kexecLoadZImageMM = %v, want nil", err)
	}
	defer img.clean()

	kernelBuf, err := os.ReadFile("../zimage/testdata/zImage")
	if err != nil {
		t.Fatal(err)
	}
	want, err := zImageMemSize(kernelBuf)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, s := range img.segments {
		if s.Phys.Start == img.entry {
			found = true
			if int64(s.Phys.Size) <= fi.Size() {
				t.Errorf("kernel segment size %#x should exceed the file size %#x", s.Phys.Size, fi.Size())
			}
			if uint(s.Phys.Size) < want {
				t.Errorf("kernel segment reserves %#x, want at least the decompressed size %#x", s.Phys.Size, want)
			}
		}
	}
	if !found {
		t.Errorf("no segment starts at entry %#x", img.entry)
	}
}

func TestKexecLoadZImageCmdlineAndInitramfs(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()
	ramfs := createFile(t, []byte("ramfs"))
	defer ramfs.Close()

	fdt := zImageFDT()
	img, err := kexecLoadZImageMM(zImageMM(t), kernel, ramfs, fdt, "console=ttyO0")
	if err != nil {
		t.Fatalf("kexecLoadZImageMM = %v, want nil", err)
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

func TestKexecLoadZImageNoMemory(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	if _, err := kexecLoadZImageMM(kexec.MemoryMap{}, kernel, nil, zImageFDT(), ""); err == nil {
		t.Error("kexecLoadZImageMM(empty mm) = nil, want error")
	}
}

// Exercise the outer entry point, which reads and parses the FDT before
// handing off. Passing a dtb keeps it off /sys/firmware/fdt.
func TestKexecLoadZImageWithDTB(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	img, err := kexecLoadZImage(kernel, nil, "console=ttyO0", fdtReader(t, zImageFDT()), nil)
	if err != nil {
		t.Fatalf("kexecLoadZImage = %v, want nil", err)
	}
	defer img.clean()

	if img.entry == 0 {
		t.Error("entry is 0")
	}
	if len(img.segments) == 0 {
		t.Error("no segments")
	}
}

func TestKexecLoadZImageBadDTB(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	if _, err := kexecLoadZImage(kernel, nil, "", bytes.NewReader([]byte("not an fdt")), nil); err == nil {
		t.Error("kexecLoadZImage(bad dtb) = nil, want error")
	}
}

// Reserved ranges must be kept away from the segments we allocate.
func TestKexecLoadZImageReservations(t *testing.T) {
	Debug = t.Logf

	kernel := openFile(t, "../zimage/testdata/zImage")
	defer kernel.Close()

	reserved := kexec.Ranges{{Start: 0x80000000, Size: 0x1000000}}
	img, err := kexecLoadZImage(kernel, nil, "", fdtReader(t, zImageFDT()), reserved)
	if err != nil {
		t.Fatalf("kexecLoadZImage = %v, want nil", err)
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
