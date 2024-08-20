// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for UART io for Linux.
package main

import (
	"strconv"
)

const (
	status = 5
	inRdy  = 1
	outRdy = 0x40
)

type uart struct {
	data, status int64
}

func openUART(comPort string) (*uart, error) {
	port, err := strconv.ParseUint(comPort, 0, 16)
	if err != nil {
		return nil, err
	}

	if err := openPort(); err != nil {
		return nil, err
	}

	return &uart{data: int64(port), status: int64(port + status)}, nil
}

func (u *uart) OK(bit uint8) (bool, error) {
	var rdy [1]byte
	if _, err := portFile.ReadAt(rdy[:], u.status); err != nil {
		return false, err
	}
	return rdy[0]&bit != 0, nil
}

func (u *uart) io(b []byte, bit uint8, f func([]byte) error) (int, error) {
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

func (u *uart) Write(b []byte) (int, error) {
	return u.io(b, outRdy, func(b []byte) error {
		_, err := portFile.WriteAt(b[:1], u.data)
		return err
	})
}

func (u *uart) Read(b []byte) (int, error) {
	return u.io(b, inRdy, func(b []byte) error {
		_, err := portFile.ReadAt(b[:1], u.data)
		return err
	})
}
