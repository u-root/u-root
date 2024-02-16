// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package guest has functions for use in tests running in VM guests.
package guest

import (
	"os"
	"testing"
)

// SkipIfNotInVM skips the test if it is not running in a vmtest-started VM.
//
// The presence of VMTEST_IN_GUEST=1 env var (which can be passed on the
// kernel commandline, using qemu.WithVmtestIdent) is used to determine this.
func SkipIfNotInVM(t testing.TB) {
	if os.Getenv("VMTEST_IN_GUEST") != "1" {
		t.Skip("Skipping test -- must be run inside vmtest VM")
	}
}

// SkipIfInVM skips the test if it is running in a vmtest-started VM.
//
// The presence of VMTEST_IN_GUEST=1 env var (which can be passed on the
// kernel commandline, using qemu.WithVmtestIdent) is used to determine this.
func SkipIfInVM(t testing.TB) {
	if os.Getenv("VMTEST_IN_GUEST") != "1" {
		t.Skip("Skipping test -- must be run inside vmtest VM")
	}
}
