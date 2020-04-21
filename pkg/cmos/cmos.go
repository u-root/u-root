// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package cmos lets you read and write to cmos registers while doing basic checks on valid register selections.

package cmos

import (
	"fmt"

	"github.com/u-root/u-root/pkg/memio"
)

const (
	cmosRegPort  = 0x70
	cmosDataPort = 0x71
)

// Read reads a register reg from CMOS into data.
func Read(reg int64, data memio.UintN) error {
	regVal := memio.Uint8(reg)
	if regVal%128 < 14 {
		return fmt.Errorf("byte %d is inside the range 0-13 which is reserved for RTC", regVal)
	}
	if err := memio.Out(cmosRegPort, &regVal); err != nil {
		return err
	}
	return memio.In(cmosDataPort, data)
}

// Write writes value data into CMOS register reg.
func Write(reg int64, data memio.UintN) error {
	regVal := memio.Uint8(reg)
	if regVal%128 < 14 {
		return fmt.Errorf("byte %d is inside the range 0-13 which is reserved for RTC", regVal)
	}
	if err := memio.Out(cmosRegPort, &regVal); err != nil {
		return err
	}
	return memio.Out(cmosDataPort, data)
}
