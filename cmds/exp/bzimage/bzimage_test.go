// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/bzimage"
)

var (
	testdataPath = "../../../pkg/boot/bzimage/testdata/"
	// Test BzImage we are not using yet.
	b = bzimage.BzImage{
		BootCode:     []byte{1, 2, 3, 4},
		KernelCode:   []byte{5, 6, 7, 8},
		KernelBase:   0x100000,
		KernelOffset: 620,
		Header: bzimage.LinuxHeader{
			MBRCode: [192]byte{
				0xea, 0x05, 0x00, 0xc0, 0x07, 0x8c, 0xc8, 0x8e,
				0xd8, 0x8e, 0xc0, 0x8e, 0xd0, 0x31, 0xe4, 0xfb,
				0xfc, 0xbe, 0x2d, 0x00, 0xac, 0x20, 0xc0, 0x74,
				0x09, 0xb4, 0x0e, 0xbb, 0x07, 0x00, 0xcd, 0x10,
				0xeb, 0xf2, 0x31, 0xc0, 0xcd, 0x16, 0xcd, 0x19,
				0xea, 0xf0, 0xff, 0x00, 0xf0, 0x55, 0x73, 0x65,
				0x20, 0x61, 0x20, 0x62, 0x6f, 0x6f, 0x74, 0x20,
				0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x0d,
				0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
				0x20, 0x64, 0x69, 0x73, 0x6b, 0x20, 0x61, 0x6e,
				0x64, 0x20, 0x70, 0x72, 0x65, 0x73, 0x73, 0x20,
				0x61, 0x6e, 0x79, 0x20, 0x6b, 0x65, 0x79, 0x20,
				0x74, 0x6f, 0x20, 0x72, 0x65, 0x62, 0x6f, 0x6f,
				0x74, 0x2e, 0x2e, 0x2e, 0x0d, 0x0a, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			ExtRamdiskImage:     0x00,
			ExtRamdiskSize:      0x00,
			ExtCmdlinePtr:       0x00,
			SetupSects:          0x1e,
			RootFlags:           0x01,
			Syssize:             0xb51d,
			RAMSize:             0x00,
			Vidmode:             0xffff,
			RootDev:             0x00,
			Bootsectormagic:     0xaa55,
			Jump:                0x66eb,
			HeaderMagic:         [4]byte{0x48, 0x64, 0x72, 0x53},
			Protocolversion:     0x20d,
			RealModeSwitch:      0x00,
			StartSys:            0x1000,
			Kveraddr:            0x3140,
			TypeOfLoader:        0x00,
			Loadflags:           0x01,
			Setupmovesize:       0x8000,
			Code32Start:         0x100000,
			RamdiskImage:        0x00,
			RamdiskSize:         0x00,
			BootSectKludge:      [4]uint8{},
			Heapendptr:          0x5320,
			ExtLoaderVer:        0x00,
			ExtLoaderType:       0x00,
			Cmdlineptr:          0x00,
			InitrdAddrMax:       0x7fffffff,
			Kernelalignment:     0x200000,
			RelocatableKernel:   0x00,
			MinAlignment:        0x15,
			XLoadFlags:          0x01,
			CmdLineSize:         0x7ff,
			HardwareSubArch:     0x00,
			HardwareSubArchData: 0x00,
			PayloadOffset:       0x255,
			PayloadSize:         0x9532c,
			SetupData:           0x00,
			PrefAddress:         0x1000000,
			InitSize:            0x6e0000,
			HandoverOffset:      0x00,
		},
	}
	uskip   = len("2018/08/10 21:20:42 ")
	jsonVer = `{
	"Release": "4.12.7",
	"Version": "#6 Fri Aug 10 14:47:18 PDT 2018",
	"Builder": "rminnich@uroot",
	"BuildNum": 6,
	"BuildTime": "2018-08-10T14:47:18`
	// The rest of this is too sensitive to formatting
	// on the various CI systems, this is enough.
)

func TestRun(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name    string
		args    []string
		debug   bool
		jsonOut bool
		want    string
	}{
		{
			name: "too big initramfs",
			args: []string{"initramfs", filepath.Join(testdataPath, "bzImage"), filepath.Join(testdataPath, "init.cpio"), "zz/zz/zz"},
			want: "new initramfs is 1536 bytes, won't fit in 480 byte old one",
		},
		{
			name: "Bad output file",
			args: []string{"initramfs", filepath.Join(testdataPath, "bzImage"), "/dev/null", "zz/zz/zz"},
			want: "open zz/zz/zz: no such file or directory",
		},
		{
			name: "correct initramfs test",
			args: []string{"initramfs", filepath.Join(testdataPath, "bzImage"), "/dev/null", filepath.Join(tmpdir, "zz")},
		},
		{
			name: "no args",
			args: []string{},
		},
		{
			name: "dump",
			args: []string{"dump", filepath.Join(testdataPath, "bzImage")},
			want: "MBRCode:0xea0500c0078cc88ed88ec08ed031e4fbfcbe2d00ac20c07409b40ebb0700cd10ebf231c0cd16cd19eaf0ff00f0557365206120626f6f74206c6f616465722e0d0a0a52656d6f7665206469736b20616e6420707265737320616e79206b657920746f207265626f6f742e2e2e0d0a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\nExtRamdiskImage:0x00\nExtRamdiskSize:0x00\nExtCmdlinePtr:0x00\nO:0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffff\nSetupSects:0x1e\nRootFlags:0x01\nSyssize:0xb51d\nRAMSize:0x00\nVidmode:0xffff\nRootDev:0x00\nBootsectormagic:0xaa55\nJump:0x66eb\nHeaderMagic:0x48647253\nProtocolversion:0x20d\nRealModeSwitch:0x00\nStartSys:0x1000\nKveraddr:0x3140\nTypeOfLoader:0x00\nLoadflags:0x01\nSetupmovesize:0x8000\nCode32Start:0x100000\nRamdiskImage:0x00\nRamdiskSize:0x00\nBootSectKludge:0x00000000\nHeapendptr:0x5320\nExtLoaderVer:0x00\nExtLoaderType:0x00\nCmdlineptr:0x00\nInitrdAddrMax:0x7fffffff\nKernelalignment:0x200000\nRelocatableKernel:0x00\nMinAlignment:0x15\nXLoadFlags:0x01\nCmdLineSize:0x7ff\nHardwareSubArch:0x00\nHardwareSubArchData:0x00\nPayloadOffset:0x255\nPayloadSize:0x9532c\nSetupData:0x00\nPrefAddress:0x1000000\nInitSize:0x6e0000\nHandoverOffset:0x00\n",
		},
		{
			name: "initramfs too many args",
			args: []string{"initramfs", "a", "b", "c", "too many"},
		},
		{
			name: "initramfs with bad input file",
			args: []string{"initramfs", "a", "b", "c"},
			want: "open a: no such file or directory",
		},
		{
			name: "initramfs with bad initramfs file",
			args: []string{"initramfs", filepath.Join(testdataPath, "bzImage"), "b", "c"},
			want: "open b: no such file or directory",
		},
		{
			name: "kernel version",
			args: []string{"ver", filepath.Join(testdataPath, "bzImage")},
			want: "4.12.7 (rminnich@uroot) #6 Fri Aug 10 14:47:18 PDT 2018\n",
		},
		{
			name:    "kernel version with jsonOut",
			args:    []string{"ver", filepath.Join(testdataPath, "bzImage")},
			jsonOut: true,
			want:    strings.ReplaceAll(jsonVer, "\t", "    "),
		},
		{
			name: "cfg with wrong image file",
			args: []string{"cfg", filepath.Join(testdataPath, "bzImage")},
			want: "embedded config not found",
		},
		{
			name: "diff",
			args: []string{"diff", filepath.Join(testdataPath, "bzImage"), filepath.Join(testdataPath, "bzimage-64kurandominitramfs")},
			want: "MBRCode:0xea0500c0078cc88ed88ec08ed031e4fbfcbe2d00ac20c07409b40ebb0700cd10ebf231c0cd16cd19eaf0ff00f0557365206120626f6f74206c6f616465722e0d0a0a52656d6f7665206469736b20616e6420707265737320616e79206b657920746f207265626f6f742e2e2e0d0a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 != MBRCode:0x4d5aea0700c0078cc88ed88ec08ed031e4fbfcbe4000ac20c07409b40ebb0700cd10ebf231c0cd16cd19eaf0ff00f00000000000000000000000000082000000557365206120626f6f74206c6f616465722e0d0a0a52656d6f7665206469736b20616e6420707265737320616e79206b657920746f207265626f6f742e2e2e0d0a005045000064860400000000000000000001000000a00006020b0202143050120000000000d0ad610010420000000200000000000000000000200000002000O:0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffff != O:0x0000000000000000740000020000000000000a000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002e73657475700000e03d000000020000e03d000000020000000000000000000000000000200050602e72656c6f63000020000000e03f000020000000e03f0000000000000000000000000000400010422e7465787400000030121200004000003012120000400000000000000000000000000000200050602e62737300000000d0ad6100305212000000000000000000000000000000000000000000800000c8000000000000000000000000000000000000000000ffffSetupSects:0x1e != SetupSects:0x1fSyssize:0xb51d != Syssize:0x12123RelocatableKernel:0x00 != RelocatableKernel:0x01XLoadFlags:0x01 != XLoadFlags:0x0bPayloadOffset:0x255 != PayloadOffset:0x3b4PayloadSize:0x9532c != PayloadSize:0xf0da4InitSize:0x6e0000 != InitSize:0x740000HandoverOffset:0x00 != HandoverOffset:0x190",
		},
		{
			name: "diff with wrong input file",
			args: []string{"diff", filepath.Join(testdataPath, "bzImage"), filepath.Join(tmpdir, "filedoesnotexist")},
			want: fmt.Sprintf("open %s: no such file or directory", filepath.Join(tmpdir, "filedoesnotexist")),
		},
		{
			name: "copy with err in write",
			args: []string{"copy", filepath.Join(testdataPath, "bzImage"), tmpdir},
			want: fmt.Sprintf("writing %s: open %s: is a directory", tmpdir, tmpdir),
		},
		{
			name: "extract",
			args: []string{"extract", filepath.Join(testdataPath, "bzImage"), tmpdir},
			want: "ramfs is 480 bytes",
		},
		{
			name:  "debug on",
			args:  []string{},
			debug: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*debug = tt.debug
			*jsonOut = tt.jsonOut
			buf := &bytes.Buffer{}
			if got := run(buf, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("run() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if !strings.Contains(buf.String(), tt.want) {
					t.Errorf("run() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}
