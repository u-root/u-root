// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for UART io for Linux.
package main

import (
	"strconv"

	"github.com/u-root/u-root/pkg/memio"
)

type MemUART struct {
	data, status int64
	m            *memio.MMap
}

func openMemUART(comPort string) (*MemUART, error) {
	port, err := strconv.ParseUint(comPort, 0, 64)
	if err != nil {
		return nil, err
	}

	m, err := memio.NewMMap("/dev/mem")
	if err != nil {
		return nil, err
	}

	if err := openPort(); err != nil {
		return nil, err
	}

	return &MemUART{data: int64(port), status: int64(port + 0x14), m: m}, nil
}

func (u *MemUART) OK(bit uint8) (bool, error) {
	var rdy memio.Uint32
	if err := u.m.ReadAt(u.status, memio.UintN(&rdy)); err != nil {
		return false, err
	}
	return uint8(rdy)&bit != 0, nil
}

func (u *MemUART) io(b []byte, bit uint8, f func([]byte) error) (int, error) {
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

func (u *MemUART) Write(b []byte) (int, error) {
	const outRdy = 0x40
	return u.io(b, outRdy, func(b []byte) error {
		val := memio.Uint32(uint32(b[0]))
		err := u.m.WriteAt(u.data, &val)
		return err
	})
}

func (u *MemUART) Read(b []byte) (int, error) {
	const inRdy = 1
	return u.io(b, inRdy, func(b []byte) error {
		val := memio.Uint32(0)
		err := u.m.ReadAt(u.data, &val)
		b[0] = uint8(val)
		return err
	})
}
