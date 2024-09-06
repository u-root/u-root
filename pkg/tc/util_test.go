// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl_test

import (
	"errors"
	"strconv"
	"testing"

	trafficctl "github.com/u-root/u-root/pkg/tc"
)

func TestParseHandle(t *testing.T) {
	for _, tt := range []struct {
		name string
		arg  string
		exp  uint32
	}{
		{
			name: "Handle_1",
			arg:  "1:",
			exp:  1 << 16,
		},
		{
			name: "Handle_1",
			arg:  "FFFF:",
			exp:  0xFFFF << 16,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ret, err := trafficctl.ParseHandle(tt.arg)
			if err != nil {
				t.Errorf("ParseHandle(%q) = %v, not nil", tt.arg, err)
			}

			if ret != tt.exp {
				t.Errorf("created handle not expected")
			}

		})
	}
}

func TestParseClassID(t *testing.T) {
	for _, tt := range []struct {
		name string
		arg  string
		exp  uint32
	}{
		{
			name: "ClassID_1:1",
			arg:  "1:1",
			exp:  (1 << 16) + 1,
		},
		{
			name: "ClassID_FFFF:FFFF",
			arg:  "FFFF:FFFF",
			exp:  0xFFFFFFFF,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ret, err := trafficctl.ParseClassID(tt.arg)
			if err != nil {
				t.Errorf("ParseHandle(%q) = %v, not nil", tt.arg, err)
			}

			if ret != tt.exp {
				t.Errorf("created handle not expected")
			}

		})
	}
}

func TestGetSize(t *testing.T) {
	for _, tt := range []struct {
		arg string
		val uint64
		err error
	}{
		{
			arg: "10k",
			val: 1024 * 10,
		},
		{
			arg: "10m",
			val: 1024 * 1024 * 10,
		},
		{
			arg: "10g",
			val: 1024 * 1024 * 1024 * 10,
		},
		{
			arg: "10kbit",
			val: 1024 * 10 / 8,
		},
		{
			arg: "10mbit",
			val: 1024 * 1024 * 10 / 8,
		},
		{
			arg: "10gbit",
			val: 1024 * 1024 * 1024 * 10 / 8,
		},
		{
			arg: "10a;sipdjghfilahbjsdfg",
			val: 0,
			err: strconv.ErrSyntax,
		},
	} {
		t.Run(tt.arg, func(t *testing.T) {
			tt := tt
			sz, err := trafficctl.ParseSize(tt.arg)
			if !errors.Is(err, tt.err) {
				t.Errorf("GetSize(%q) = %v, not %v", tt.arg, err, tt.err)
			}

			if sz != tt.val {
				t.Errorf("got %d, but want: %d", sz, tt.val)
			}
		})
	}
}

func TestParseRate(t *testing.T) {
	for _, tt := range []struct {
		arg string
		exp uint64
		err error
	}{
		{arg: "5mbit", exp: 625000},
	} {
		t.Run(tt.arg, func(t *testing.T) {
			ret, err := trafficctl.ParseRate(tt.arg)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseRate(%q) = %v, not %v", tt.arg, err, tt.err)
			}

			if ret != tt.exp {
				t.Errorf("got :%d, not %d", ret, tt.exp)
			}
		})
	}
}

func TestParseLinkLayer(t *testing.T) {
	for _, tt := range []struct {
		arg string
		exp uint8
		err error
	}{
		{arg: "ethernet", exp: 1},
		{arg: "atm", exp: 2},
		{arg: "ads1", exp: 2},
		{arg: "random", exp: 0xFF, err: trafficctl.ErrUnknownLinkLayer},
	} {
		t.Run(tt.arg, func(t *testing.T) {
			ret, err := trafficctl.ParseLinkLayer(tt.arg)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseLinkLayer(%q) = %v, not %v", tt.arg, err, tt.err)
			}

			if ret != tt.exp {
				t.Errorf("ParseLinkLayer(%q) = %v, not %v", tt.arg, ret, tt.exp)
			}
		})
	}
}
