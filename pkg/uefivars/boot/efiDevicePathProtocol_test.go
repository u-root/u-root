// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"errors"
	"testing"

	"github.com/u-root/u-root/pkg/uefivars"
)

var boot7 = []byte{
	0x01, 0x00, 0x00, 0x00, 0x5e, 0x00, 0x55, 0x00, 0x45, 0x00, 0x46, 0x00,
	0x49, 0x00, 0x20, 0x00, 0x4f, 0x00, 0x53, 0x00, 0x00, 0x00, 0x04, 0x01,
	0x2a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xcd, 0x5c,
	0x63, 0x81, 0x4f, 0x1b, 0x3f, 0x4d, 0xb7, 0xb8, 0xf7, 0x8a, 0x5b, 0x02,
	0x9f, 0x35, 0x02, 0x02, 0x04, 0x04, 0x30, 0x00, 0x5c, 0x00, 0x45, 0x00,
	0x46, 0x00, 0x49, 0x00, 0x5c, 0x00, 0x42, 0x00, 0x4f, 0x00, 0x4f, 0x00,
	0x54, 0x00, 0x5c, 0x00, 0x42, 0x00, 0x4f, 0x00, 0x4f, 0x00, 0x54, 0x00,
	0x58, 0x00, 0x36, 0x00, 0x34, 0x00, 0x2e, 0x00, 0x45, 0x00, 0x46, 0x00,
	0x49, 0x00, 0x00, 0x00, 0x7f, 0xff, 0x04, 0x00, 0x00, 0x00, 0x42, 0x4f,
}

var boot8 = []byte{
	0x01, 0x00, 0x00, 0x00, // attributes
	0x5e, 0x00, // path length
	0x55, 0x00, 0x45, 0x00, 0x46, 0x00, 0x49, 0x00, 0x20, 0x00, 0x4f, 0x00, 0x53, 0x00, 0x00, 0x00, // description
	0x04, 0x01, // node type, subtype (Media, HD)
	0x2a, 0x00, // node len
	0x02, 0x00, 0x00, 0x00, // partition number
	0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // partition start
	0x00, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // partition size
	0xde, 0xad, 0xef, 0xde, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // partition signature
	0x01, 0x01, // partition format, signature type
	0x04, 0x04,
	0x30, 0x00, 0x5c, 0x00, 0x45, 0x00,
	0x46, 0x00, 0x49, 0x00, 0x5c, 0x00, 0x42, 0x00, 0x4f, 0x00, 0x4f, 0x00,
	0x54, 0x00, 0x5c, 0x00, 0x42, 0x00, 0x4f, 0x00, 0x4f, 0x00, 0x54, 0x00,
	0x58, 0x00, 0x36, 0x00, 0x34, 0x00, 0x2e, 0x00, 0x45, 0x00, 0x46, 0x00,
	0x49, 0x00, 0x00, 0x00, 0x7f, 0xff, 0x04, 0x00, 0x00, 0x00, 0x42, 0x4f,
}

// func ParseFilePathList(in []byte) (EfiDevicePathProtocolList, error)
func TestParseFilePathList(t *testing.T) {
	// When this test runs, you will see log entries like
	// "Skipping loop0: open /dev/loop0: permission denied"
	// These entries are safe to ignore, unless you ran as root (!) in which
	// case the devices ought to be readable.
	tests := []struct {
		name              string
		e                 uefivars.EfiVar
		wantpath          string
		wantpathlen       int
		wantstr           string
		wantdesc          string
		wantResolveError  []error
		wantResolveOutput []string
	}{
		{
			name: "GPT file path test",
			e: uefivars.EfiVar{
				UUID: BootUUID,
				Name: "Boot0007",
				Data: boot7,
			},
			// same as efibootmgr output, except using forward slashes
			wantpath:          "HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI)",
			wantpathlen:       2,
			wantstr:           `Boot0007: attrs=0x1, desc="UEFI OS", path=HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`,
			wantdesc:          `Boot0007: attrs=0x1, desc="UEFI OS", path=HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`,
			wantResolveError:  []error{ErrNotFound, nil},
			wantResolveOutput: []string{"", "/EFI/BOOT/BOOTX64.EFI"},
		},
		{
			name: "MBR file path test",
			e: uefivars.EfiVar{
				UUID: BootUUID,
				Name: "Boot0008",
				Data: boot8,
			},
			// same as efibootmgr output, except using forward slashes
			wantpath:          "HD(2,MBR,deef-adde,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI)",
			wantpathlen:       2,
			wantstr:           `Boot0008: attrs=0x1, desc="UEFI OS", path=HD(2,MBR,deef-adde,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`,
			wantdesc:          `Boot0008: attrs=0x1, desc="UEFI OS", path=HD(2,MBR,deef-adde,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`,
			wantResolveError:  []error{ErrNotFound, nil},
			wantResolveOutput: []string{"", "/EFI/BOOT/BOOTX64.EFI"},
		},
	}

	for _, tc := range tests {
		b := BootVar(tc.e)

		gotpath := b.FilePathList.String()
		if gotpath != tc.wantpath {
			t.Errorf("mismatch\nwant %q\n got %q", tc.wantpath, gotpath)
		}

		gotstr := b.String()
		if tc.wantstr != gotstr {
			t.Errorf("mismatch\nwant %s\n got %s", tc.wantstr, gotstr)
		}

		gotdesc := b.String()
		if gotdesc != tc.wantdesc {
			t.Errorf("mismatch\nwant %s\n got %s", tc.wantdesc, gotdesc)
		}

		if len(b.FilePathList) != tc.wantpathlen {
			t.Errorf("mismatch\nwant pathlen %d\n got pathlen %d", tc.wantpathlen, len(b.FilePathList))
		}

		for n, p := range b.FilePathList {
			r, err := p.Resolver()
			if err != nil {
				if tc.wantResolveError[n] == nil {
					t.Errorf("mismatch\nwant no error, got %s", err)
				}
				if !errors.Is(err, tc.wantResolveError[n]) {
					t.Errorf("mismatch\nwant error %s, got %s", tc.wantResolveError[n], err)
				}
			} else {
				if tc.wantResolveError[n] != nil {
					t.Errorf("mismatch\nwant %s, got no error", tc.wantResolveError[n])
				}
				if r.String() != tc.wantResolveOutput[n] {
					t.Errorf("mismatch\nwant %s, got %s", tc.wantResolveOutput[n], r.String())
				}
			}
		}
	}
}
