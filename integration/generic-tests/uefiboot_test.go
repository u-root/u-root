// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!race

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	expect "github.com/google/goexpect"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestUefiboot tests uefiboot commmands to boot to uefishell.
func TestUefiBoot(t *testing.T) {
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	payload := "UEFIPAYLOAD.fd"
	src := fmt.Sprintf("/home/circleci/%v", payload)
	if tk := os.Getenv("UROOT_TEST_UEFIPAYLOAD_DIR"); len(tk) > 0 {
		src = filepath.Join(tk, payload)
	}

	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		t.Skipf("UEFI payload image is not found: %s\n Usage: uefiboot <payload>", src)
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: []string{"uefiboot /dev/sda"},
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				qemu.IDEBlockDevice{File: src},
				qemu.ArbitraryArgs{"-machine", "q35"},
				qemu.ArbitraryArgs{"-m", "2048"},
			},
		},
	})
	defer cleanup()

	// Edk2 debug mode will print PROGRESS CODE. We will want to make sure
	// payload is booting to uefi shell correctly.
	if out, err := q.ExpectBatch([]expect.Batcher{
		// Finish booting.
		&expect.BExp{R: "PROGRESS CODE: V02070003"},
		// Last code before booting to UEFI Shell
		&expect.BExp{R: "PROGRESS CODE: V03058001"},
	}, 50*time.Second); err != nil {
		t.Fatalf("VM output did not match expectations: %v", out)
	}
}
