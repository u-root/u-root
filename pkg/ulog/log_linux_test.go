// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ulog

import (
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestKernelLog(t *testing.T) {
	// This is an integration test run in QEMU.
	testutil.SkipIfNotRoot(t)

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
