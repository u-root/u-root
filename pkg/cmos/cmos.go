// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386

package cmos

import (
	"github.com/u-root/u-root/pkg/memio"
)

const (
	cmosRegPort  = 0x70
	cmosDataPort = 0x71
)

// Read reads a register reg from CMOS into data.
func Read(reg memio.Uint8, data memio.UintN) error {
	if err := memio.Out(cmosRegPort, &reg); err != nil {
		return err
	}
	return memio.In(cmosDataPort, data)
}

// Write writes value data into CMOS register reg.
func Write(reg memio.Uint8, data memio.UintN) error {
	if err := memio.Out(cmosRegPort, &reg); err != nil {
		return err
	}
	return memio.Out(cmosDataPort, data)
}
