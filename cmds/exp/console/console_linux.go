// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
)

func console(serial string) (io.Reader, io.Writer, error) {
	in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

	// This switch is kind of hokey, true, but it's also quite convenient for users.
	switch {
	// A raw IO port for serial console
	case []byte(serial)[0] == '0':
		u, err := openUART(serial)
		if err != nil {
			return nil, nil, fmt.Errorf("console exits: sorry, can't get a uart: %w", err)
		}
		in, out = u, u
	case serial == "i8042":
		u, err := openi8042()
		if err != nil {
			return nil, nil, fmt.Errorf("console exits: sorry, can't get an i8042: %w", err)
		}
		in, out = u, os.Stdout
	case serial == "stdio":
	default:
		return nil, nil, fmt.Errorf("console must be one of stdio, i8042, or an IO port with a leading 0 (e.g. 0x3f8)")
	}

	return in, out, nil
}
