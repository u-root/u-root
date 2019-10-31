// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"os"
)

func Open(devnum int) (*Ipmi, error) {
	d := fmt.Sprintf("/dev/ipmi%d", devnum)

	f, err := os.OpenFile(d, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &Ipmi{File: f}, nil
}
