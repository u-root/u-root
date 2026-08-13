// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !amd64 && !386

package acpi

import (
	"testing"
)

func TestRSDPGetters(t *testing.T) {
	if len(rsdpgetters) != 1 {
		t.Fatalf("expected exactly 1 RSDP getter on non-x86 architectures, got %d", len(rsdpgetters))
	}
}
