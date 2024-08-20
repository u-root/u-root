// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains common variables for port IO devices for Linux
package main

import (
	"os"
)

var portFile *os.File

func openPort() (err error) {
	portFile, err = os.OpenFile("/dev/port", os.O_RDWR, 0)
	return err
}
