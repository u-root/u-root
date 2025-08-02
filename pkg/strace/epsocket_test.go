// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"golang.org/x/sys/unix"
)

func TestGetAddress(t *testing.T) {
	for _, tt := range []struct {
		name string
		addr []byte
		f    *FullAddress
		err  error
	}{
		{name: "empty", addr: []byte{}, f: nil, err: io.EOF},
		{name: "too short", addr: []byte{unix.AF_UNIX, 0, 0}, f: nil, err: unix.EINVAL},
		{name: "bad family", addr: []byte{0, 2, 'h', 'i', 0}, f: nil, err: unix.ENOTSUP},
		{name: "unix", addr: []byte{unix.AF_UNIX, 0, 'h', 'i', 0}, f: &FullAddress{Addr: "hi"}, err: nil},
		{name: "unix no null", addr: []byte{unix.AF_UNIX, 0, 'h', 'i'}, f: nil, err: unix.EINVAL},
		{name: "unix ENAMETOOLONG", addr: (&[unix.PathMax * 2]byte{unix.AF_UNIX, 0, 'h', 'i', 0})[:], f: nil, err: unix.ENAMETOOLONG},
		{name: "IP4", addr: []byte{unix.AF_INET, 0, 13, 14, 1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0}, f: &FullAddress{Addr: "1.2.3.4", Port: 3342}, err: nil},
		{name: "IP4short", addr: []byte{unix.AF_INET, 0, 13}, f: nil, err: unix.EFAULT},
		{name: "IP6", addr: []byte{unix.AF_INET6, 0, 13, 14, 0xde, 0xad, 0xbe, 0xef, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0xa, 0xb, 0xc, 0xd}, f: &FullAddress{Addr: "1:203:405:607:809:a0b:c0d:e0f", Port: 3342}, err: nil},
		{name: "IP6short", addr: []byte{unix.AF_INET6, 0, 13, 14}, f: nil, err: unix.EFAULT},
	} {
		f, err := GetAddress(tt.addr)
		if !errors.Is(err, tt.err) {
			t.Errorf("%s:got err %v, want %v", tt.name, err, tt.err)
			continue
		}
		if !reflect.DeepEqual(f, tt.f) {
			t.Errorf("%s:got FullAddress %s, want %s", tt.name, f.String(), tt.f.String())
		}
	}
}

func TestFullAddressStringer(t *testing.T) {
	var a *FullAddress
	s := a.String()
	if s != ":" {
		t.Errorf("nil String: got %q, want %q", s, ":")
	}
	a = &FullAddress{Addr: "unix", Port: 0xdead}
	s = a.String()
	if s != "unix:0xdead" {
		t.Errorf("String: got %q, want %q", s, "unix:0xdead")
	}
}
