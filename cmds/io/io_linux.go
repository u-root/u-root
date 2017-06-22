// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for io for Linux.
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func doio(a16 uint64, f func(*os.File) error) error {
	port, err := os.OpenFile("/dev/port", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer port.Close()
	_, err = port.Seek(int64(a16), 0)
	if err != nil {
		return fmt.Errorf("in: bad address %v: %v", a16, err)
	}
	return f(port)
}

func in(a16 uint64, data interface{}) error {
	return doio(a16, func(port *os.File) error {
		return binary.Read(port, binary.LittleEndian, data)
	})
}

func out(a16 uint64, data interface{}) error {
	return doio(a16, func(port *os.File) error {
		return binary.Write(port, binary.LittleEndian, data)
	})

}
