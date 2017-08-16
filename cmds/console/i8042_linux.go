// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for i8042 io for Linux.
package main

const (
	Data = 0x60 /* data port */

	Status   = 0x64 /* status port */
	Inready  = 0x01 /*  input character ready */
	Outbusy  = 0x02 /*  output busy */
	Sysflag  = 0x04 /*  system flag */
	Cmddata  = 0x08 /*  cmd==0 data==1 */
	Inhibit  = 0x10 /*  keyboard/mouse inhibited */
	Minready = 0x20 /*  mouse character ready */
	Rtimeout = 0x40 /*  general timeout */
	Parity   = 0x80

	Cmd   = 0x64 /* command port (write only) */
	Nscan = 128
)

type i8042 struct {
	data, status int64
}

func openi8042() (*i8042, error) {
	if err := openPort(); err != nil {
		return nil, err
	}

	return &i8042{data: Data, status: Status}, nil
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
	return u.io(b, Inready, func(b []byte) error {
		_, err := portFile.ReadAt(b[:1], u.data)
		return err
	})
}
