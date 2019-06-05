// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"
)

// TestMountKexec runs an init which mounts a filesystem and kexecs a kernel.
func TestMountKexec(t *testing.T) {
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	// Create the CPIO and start QEMU.
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/kexec/uinit",
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/mount",
			"github.com/u-root/u-root/cmds/core/kexec",
		},
	})
	defer cleanup()

	if err := q.Expect("KEXECCOUNTER=0"); err != nil {
		t.Fatal(err)
	}
	if err := q.Expect("KEXECCOUNTER=1"); err != nil {
		t.Fatal(err)
	}
}
