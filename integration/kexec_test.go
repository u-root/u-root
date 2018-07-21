// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"
)

// TestMountKExec runs an init which mounts a filesystem and kexecs a kernel.
func TestMountKExec(t *testing.T) {
	// Create the CPIO and start QEMU.
	tmpDir, q := testWithQEMU(t, "kexec", []string{})
	defer cleanup(t, tmpDir, q)

	if err := q.Expect("KEXECCOUNTER=1"); err != nil {
		t.Fatal(err)
	}
}
