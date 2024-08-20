// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for i8042 io for Linux.
package main

const (
	i8042Data = 0x60 /* data port */

	i8042Status   = 0x64 /* status port */
	i8042Inready  = 0x01 /*  input character ready */
	i8042Outbusy  = 0x02 /*  output busy */
	i8042Sysflag  = 0x04 /*  system flag */
	i8042Cmddata  = 0x08 /*  cmd==0 data==1 */
	i8042Inhibit  = 0x10 /*  keyboard/mouse inhibited */
	i8042Minready = 0x20 /*  mouse character ready */
	i8042Rtimeout = 0x40 /*  general timeout */
	i8042Parity   = 0x80

	i8042Cmd   = 0x64 /* command port (write only) */
	i8042Nscan = 128
)

type i8042 struct {
	data, status int64
}

func openi8042() (*i8042, error) {
	if err := openPort(); err != nil {
		return nil, err
	}

	return &i8042{data: i8042Data, status: i8042Status}, nil
}

func (u *i8042) OK(bit uint8) (bool, error) {
	var rdy [1]byte
	if _, err := portFile.ReadAt(rdy[:], u.status); err != nil {
		return false, err
	}
	return rdy[0]&bit != 0, nil
}

func (u *i8042) io(b []byte, bit uint8, f func([]byte) error) (int, error) {
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

func (u *i8042) Read(b []byte) (int, error) {
	return u.io(b, i8042Inready, func(b []byte) error {
		_, err := portFile.ReadAt(b[:1], u.data)
		return err
	})
}
