// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for UART io for Linux.
package main

import (
	"os"
)

const (
	data   = 0x3f8
	status = 0x3fd
	inRdy  = 1
	outRdy = 0x40
)

type uart struct{}

var port *os.File

func openUART() (err error) {
	port, err = os.OpenFile("/dev/port", os.O_RDWR, 0)
	return
}

func (uart) OK(bit uint8) (bool, error) {
	var rdy [1]byte
	if _, err := port.ReadAt(rdy[:], status); err != nil {
		return false, err
	}
	return rdy[0]&bit != 0, nil
}

func (u uart) io(b []byte, bit uint8, f func([]byte) error) (int, error) {
	var (
		err error
		amt int
	)

	for amt = 0; amt < len(b); {
		rdy, err := u.OK(bit)
		if err != nil {
			break
		}
		if !rdy {
			continue
		}
		if err = f(b[amt:]); err != nil {
			break
		}
		amt = amt + 1
	}
	return amt, err
}

func (u uart) Write(b []byte) (int, error) {
	return u.io(b, outRdy, func(b []byte) error {
		_, err := port.WriteAt(b[:1], data)
		return err
	})
}

func (u uart) Read(b []byte) (int, error) {
	return u.io(b, inRdy, func(b []byte) error {
		_, err := port.ReadAt(b[:1], data)
		return err
	})
}
