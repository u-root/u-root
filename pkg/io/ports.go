// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386

package io

import (
	"fmt"
)

const portPath = "/dev/port"

// In reads data from the x86 port at address addr. data must be one of:
// *uint8, *uint16, or *uint32.
func In(addr uint16, data interface{}) error {
	switch data.(type) {
	case *uint8, *uint16, *uint32:
	default:
		return fmt.Errorf("cannot read port type %T", data)
	}
	return pathRead(portPath, int64(addr), data)
}

// Out writes data to the x86 port at address addr. data must be one of: uint8,
// uint16, or uint32.
func Out(addr uint16, data interface{}) error {
	switch data.(type) {
	case uint8, uint16, uint32:
	default:
		return fmt.Errorf("cannot write port type %T", data)
	}
	return pathWrite(portPath, int64(addr), data)
}
