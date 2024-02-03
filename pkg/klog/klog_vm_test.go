// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package klog

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestIntegration(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	vmtest.RunGoTestsInVM(t, []string{"github.com/u-root/u-root/pkg/klog"},
		vmtest.WithVMOpt(vmtest.WithQEMUFn(qemu.WithVMTimeout(time.Minute))),
	)
}

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
