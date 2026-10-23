// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func console(serial string) (io.Reader, io.Writer, error) {

	// for all systems we have interest in, a 16-bit port address means we are in
	// /dev/port space.
	_, err16 := strconv.ParseUint(serial, 0, 16)

	// This switch is kind of hokey, true, but it's also quite convenient for users.
	// All they have to do to pick a devices is set an address or a name.
	switch {
	// A raw IO port for serial console
	case serial[0] == '0' && err16 == nil:
		u, err := openUART(serial)
		if err != nil {
			return nil, nil, fmt.Errorf("console exits: sorry, can't get a uart: %w", err)
		}
		return u, u, nil
	case serial[0] == '0':
		u, err := openMemUART(serial)
		if err != nil {
			return nil, nil, fmt.Errorf("console exits: sorry, can't get a uart: %w", err)
		}
		return u, u, nil
	case serial == "i8042":
		u, err := openi8042()
		if err != nil {
			return nil, nil, fmt.Errorf("console exits: sorry, can't get an i8042: %w", err)
		}
		return u, os.Stdout, nil
	case serial == "stdio":
	default:
		return nil, nil, fmt.Errorf("console must be one of stdio, i8042, or an IO port with a leading 0 (e.g. 0x3f8)")
	}

	return io.Reader(os.Stdin), io.Writer(os.Stdout), nil
}
