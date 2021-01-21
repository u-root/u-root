// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"os"
)

// Open a channel to an IPMI device /dev/ipmi{devnum}.
func Open(devnum int) (*IPMI, error) {
	d := fmt.Sprintf("/dev/ipmi%d", devnum)

	f, err := os.OpenFile(d, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &IPMI{File: f}, nil
}
