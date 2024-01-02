// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ulog

import (
	"os"
	"strings"
	"testing"
)

func TestKernelLog(t *testing.T) {
	// This is an integration test run in QEMU.
	// Cannot use guest.SkipIfNotInVM here as vmtest depends on ulog --
	// import cycles in tests are forbidden.
	if os.Getenv("VMTEST_IN_GUEST") != "1" {
		t.Skipf("Test only runs in VM")
	}

	// do something.
	KernelLog.Printf("haha %v", "oh foobar")

	want := "haha oh foobar"
	b := make([]byte, 1024)
	n, err := KernelLog.Read(b)
	if err != nil {
		t.Fatalf("Could not read from kernel log: %v", err)
	}
	if got := string(b[:n]); strings.Contains(got, want) {
		t.Errorf("kernel log read = %v (len %d), want it to include %v", got, n, want)
	}
}
