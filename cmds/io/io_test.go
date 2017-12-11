// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

type op struct {
	name   string
	addr   uint64
	size   int
	retval []byte
	val    interface{}
	want   uint64
	outerr error
	inerr  error
}

func (o op) Seek(addr int64, whence int) (int64, error) {
	if whence != 0 {
		return -1, fmt.Errorf("Fix your seek methed")
	}
	if addr == 9 && whence == 0 {
		return -1, fmt.Errorf("Seek: illegal seek to %d", addr)
	}
	if o.addr != uint64(addr) {
		return -1, fmt.Errorf("Seek: asked to seek to %v, want %v", addr, o.addr)
	}
	return addr, nil
}

func (o op) Read(b []byte) (int, error) {
	// check equality right here.
	return len(o.retval), o.inerr
}

func (o op) Write(b []byte) (int, error) {
	if len(b) > 8 {
		return -1, fmt.Errorf("Bad write size: %d bytes", len(b))
	}
	// check equality right here.
	return len(b), o.outerr
}

// checkError checks two error cases to make sure they match
// in some reasonable way, since there are four cases ...
func checkError(msg string, got, want error) error {
	if got == nil && want == nil {
		return nil
	}
	if got != nil && want != nil && got.Error() == want.Error() {
		return nil
	}
	return fmt.Errorf("%s: Got %v, want %v", msg, got, want)
}

func TestIO(t *testing.T) {
	var ops = []op{
		{name: "Test bad seek", addr: 9, size: -1, outerr: fmt.Errorf("in: bad address 9: Seek: illegal seek to 9")},
		{name: "Write and Read byte", val: uint8(1), want: 1, retval: []byte{1}},
		{name: "Write and Read 16 bits", val: uint16(0x12), want: 0x12, retval: []byte{1, 2}},
		{name: "Write and Read 32 bits", val: uint32(0x1234), want: 0x1234, retval: []byte{1, 2, 3, 4}},
		{name: "Write and Read 64 bits", val: uint64(0x12345678), want: 0x12345678, retval: []byte{1, 2, 3, 4, 5, 6, 7, 8}},
	}

	for i, o := range ops {
		err := out(o, o.addr, o.val)
		if err := checkError(o.name, err, o.outerr); err != nil {
			t.Errorf("%v", err)
		}
		if i == 0 {
			continue
		}
		err = in(o, o.addr, o.val)
		if err := checkError(o.name, err, o.inerr); err != nil {
			t.Errorf("%v", err)
		}
		var val uint64
		switch oval := o.val.(type) {
		case uint8:
			val = uint64(oval)
		case uint16:
			val = uint64(oval)
		case uint32:
			val = uint64(oval)
		case uint64:
			val = oval
		default:
			t.Fatalf("Can't handle %T for %v command", t, o)
		}

		if val != o.want {
			t.Errorf("Write and read: got %v, want %v", val, o.want)
		}
	}
}
