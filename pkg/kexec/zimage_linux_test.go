// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type kexecLoadArgs struct {
	entry    uintptr
	segments []Segment
	flags    uint64
}

func (k *kexecLoadArgs) String() string {
	s := fmt.Sprintf("  entry: %#v\n", k.entry)
	for i, seg := range k.segments {
		s += fmt.Sprintf("  segments[%d] = %s\n", i, seg)
	}
	s += fmt.Sprintf("  flags: %#v", k.flags)
	return s
}

func TestZImageLoad(t *testing.T) {
	// Prevent failure when run on an arch where page size is not 4KiB.
	oldPageMask := pageMask
	pageMask = 4096 - 1
	defer func() { pageMask = oldPageMask }()

	for _, tt := range []struct {
		name                               string
		kernelFile, initramfsFile, fdtFile string
		want                               *kexecLoadArgs
	}{
		{
			name:          "rpi3",
			kernelFile:    "testdata/zImage",
			initramfsFile: "testdata/fake_initramfs.cpio",
			fdtFile:       "../dt/testdata/rpi_fdt.dtb",
			want: &kexecLoadArgs{
				entry: 0x8000,
				segments: []Segment{
					// Kernel
					{Buf: Range{Size: 0xeac88}, Phys: Range{Start: 0, Size: 0xeb000}},
					// Initramfs
					{Buf: Range{Size: 0xf}, Phys: Range{Start: 0xeb000, Size: 0x1000}},
					// Device tree
					{Buf: Range{Size: 0x6877}, Phys: Range{Start: 0xed000, Size: 0x7000}},
				},
				flags: 0,
			},
		},
		{
			name:       "rpi3_noinitramfs",
			kernelFile: "testdata/zImage",
			fdtFile:    "../dt/testdata/rpi_fdt.dtb",
			want: &kexecLoadArgs{
				entry: 0x8000,
				segments: []Segment{
					// Same as above but skips initramfs.
					{Buf: Range{Size: 0xeac88}, Phys: Range{Start: 0, Size: 0xeb000}},
					{Buf: Range{Size: 0x6833}, Phys: Range{Start: 0xec000, Size: 0x7000}},
				},
				flags: 0,
			},
		},
		{
			name:          "qemu",
			kernelFile:    "testdata/zImage",
			initramfsFile: "testdata/fake_initramfs.cpio",
			fdtFile:       "../dt/testdata/qemu_fdt.dtb",
			want: &kexecLoadArgs{
				entry:    0,           // TODO
				segments: []Segment{}, // TODO
				flags:    0,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: The memory node in the qemu device tree is not yet supported.
			if tt.name == "qemu" {
				t.Skip()
			}

			// Read kernel.
			kernel, err := os.Open(tt.kernelFile)
			if err != nil {
				t.Fatal(err)
			}
			defer kernel.Close()

			// Read device tree.
			var initramfs *os.File
			if tt.initramfsFile != "" {
				initramfs, err = os.Open(tt.initramfsFile)
				if err != nil {
					t.Fatal(err)
				}
				defer initramfs.Close()
			}

			// Read device tree.
			fdt, err := os.Open(tt.fdtFile)
			if err != nil {
				t.Fatal(err)
			}
			defer fdt.Close()

			got := &kexecLoadArgs{}

			opts := DryrunOpts()
			opts.Kernel = kernel
			opts.Initramfs = initramfs
			opts.CmdLine = "hello"
			opts.DTB = fdt
			opts.LoadSyscall = func(entry uintptr, segments []Segment, flags uint64) error {
				got.entry = entry
				got.segments = segments
				got.flags = flags
				return nil
			}

			if err := ZImageLoad(opts); err != nil {
				t.Fatalf("processZImage error: %v", err)
			}

			// Virtual addresses change every run, so clear before comparing.
			for i := range got.segments {
				got.segments[i].Buf.Start = 0
			}

			// Validate the segments.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got kexec_load:\n%v\nwant kexec_load:\n%v", got, tt.want)
			}
		})
	}
}
