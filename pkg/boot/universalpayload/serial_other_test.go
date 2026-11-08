// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !riscv64

package universalpayload

import "testing"

func TestDiscoverPlatformSerialUnsupported(t *testing.T) {
	mmio, stride, baud, base, ok := discoverPlatformSerial()
	if ok || mmio != 0 || stride != 0 || baud != 0 || base != 0 {
		t.Fatalf("discoverPlatformSerial() = mmio %d stride %d baud %d base %#x ok %v, want unset", mmio, stride, baud, base, ok)
	}
}
