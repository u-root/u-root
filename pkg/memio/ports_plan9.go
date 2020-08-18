// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package memio

import (
	"fmt"
)

// In reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func In(addr uint16, data UintN) error {
	switch data.(type) {
	case *Uint32:
		return pathRead("#P/iol", int64(addr), data)
	case *Uint16:
		return pathRead("#P/iow", int64(addr), data)
	case *Uint8:
		return pathRead("#P/iob", int64(addr), data)
	}
	return fmt.Errorf("port data must be 8, 16 or 32 bits")
}

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func Out(addr uint16, data UintN) error {
	switch data.(type) {
	case *Uint32:
		return pathWrite("#P/iol", int64(addr), data)
	case *Uint16:
		return pathWrite("#P/iow", int64(addr), data)
	case *Uint8:
		return pathWrite("#P/iob", int64(addr), data)
	}
	return fmt.Errorf("port data must be 8, 16 or 32 bits")
}
