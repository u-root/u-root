// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !riscv64

package universalpayload

func discoverPlatformSerial() (useMmio uint8, stride uint8, baud uint32, base uint64, ok bool) {
	return 0, 0, 0, 0, false
}
