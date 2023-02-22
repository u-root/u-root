// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"os"
)

// OpenPath opens a channel to an IPMI device by path (e.g., /dev/ipmi0).
func OpenPath(path string) (*IPMI, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &IPMI{dev: newDev(f)}, nil
}

// Open opens a channel to an IPMI device by device number (e.g., 0 for /dev/ipmi{devnum}).
func Open(devnum int) (*IPMI, error) {
	return OpenPath(fmt.Sprintf("/dev/ipmi%d", devnum))
}
