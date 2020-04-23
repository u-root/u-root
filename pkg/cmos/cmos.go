// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmos

import (
	"github.com/u-root/u-root/pkg/memio"
)

func Read(reg uint64, data memio.UintN) error {
        regVal := memio.Uint8(reg)
        if err := memio.Out(0x70, &regVal); err != nil {
                return err
        }
        return memio.In(0x71, data)
}

func Write(reg uint64, data memio.UintN) error {
        regVal := memio.Uint8(reg)
        if err := memio.Out(0x70, &regVal); err != nil {
                return err
        }
        return memio.Out(0x71, data)
}

