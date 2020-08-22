// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386

package memio

import (
	"fmt"
)

const portPath = "/dev/port"

// In reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func In(addr uint16, data UintN) error {
	if _, ok := data.(*Uint64); ok {
		return fmt.Errorf("port data must be 8, 16 or 32 bits")
	}
	return pathRead(portPath, int64(addr), data)
}

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func Out(addr uint16, data UintN) error {
	if _, ok := data.(*Uint64); ok {
		return fmt.Errorf("port data must be 8, 16 or 32 bits")
	}
	return pathWrite(portPath, int64(addr), data)
}
