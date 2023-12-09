// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
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

// func ParseFilePathList(in []byte) (EfiDevicePathProtocolList, error)
func TestParseFilePathList(t *testing.T) {
	// When this test runs, you will see log entries like
	// "Skipping loop0: open /dev/loop0: permission denied"
	// These entries are safe to ignore, unless you ran as root (!) in which
	// case the devices ought to be readable.
	e := uefivars.EfiVar{
		UUID: BootUUID,
		Name: "Boot0007",
		Data: boot7,
	}
	b := BootVar(e)

	// same as efibootmgr output, except using forward slashes
	wantpath := "HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI)"
	gotpath := b.FilePathList.String()
	if gotpath != wantpath {
		t.Errorf("mismatch\nwant %q\n got %q", wantpath, gotpath)
	}
	wantstr := `Boot0007: attrs=0x1, desc="UEFI OS", path=HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`
	gotstr := b.String()
	if wantstr != gotstr {
		t.Errorf("mismatch\nwant %s\n got %s", wantstr, gotstr)
	}

	wantdesc := `Boot0007: attrs=0x1, desc="UEFI OS", path=HD(1,GPT,81635ccd-1b4f-4d3f-b7b8-f78a5b029f35,0x40,0xf000)/File(/EFI/BOOT/BOOTX64.EFI), opts=00e4bd82`
	gotdesc := b.String()
	if gotdesc != wantdesc {
		t.Errorf("mismatch\nwant %s\n got %s", wantdesc, gotdesc)
	}
	expectedOutput := "described device not found\n/EFI/BOOT/BOOTX64.EFI\n"
	resolveFailed := false
	var output string
	for _, p := range b.FilePathList {
		r, err := p.Resolver()
		if err != nil {
			resolveFailed = true
			output += err.Error() + "\n"
		} else {
			output += r.String() + "\n"
		}
	}
	if !resolveFailed {
		t.Error("resolve should fail - the chances of a device matching the guid are infinitesimally small")
	}
	if output != expectedOutput {
		t.Errorf("\nwant %s\n got %s", expectedOutput, output)
	}
}
