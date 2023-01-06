// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"runtime"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// TestRSDP tests whether any method for getting an RSDP works.
func TestRSDP(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	// Our QEMU aarch64 does not boot via UEFI, so RSDP only works on x86.
	if runtime.GOARCH != "amd64" {
		t.Skip("RSDP in QEMU only available on amd64 for now")
	}

	_, err := GetRSDP()
	if err != nil {
		t.Fatalf("GetRSDP: got %v, want nil", err)
	}
}
