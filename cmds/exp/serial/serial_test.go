// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"reflect"
	"testing"

	"go.bug.st/serial"
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		device     string
		parity     string
		baud       uint
		databits   int
		wantError  error
		wantParams params
	}{
		{
			device:    "",
			wantError: errUsage,
		},
		{
			device:     "dev",
			parity:     "odd",
			baud:       300,
			databits:   8,
			wantParams: params{device: "dev", parity: serial.OddParity, baud: 300, databits: 8},
		},
		{
			device:     "dev",
			parity:     "even",
			baud:       300,
			databits:   8,
			wantParams: params{device: "dev", parity: serial.EvenParity, baud: 300, databits: 8},
		},
		{
			device:     "dev",
			parity:     "no",
			baud:       300,
			databits:   8,
			wantParams: params{device: "dev", parity: serial.NoParity, baud: 300, databits: 8},
		},
		{
			device:    "dev",
			parity:    "other",
			baud:      300,
			databits:  8,
			wantError: errUsage,
		},
		{
			device:    "dev",
			parity:    "no",
			baud:      300,
			databits:  16,
			wantError: errUsage,
		},
	}

	for _, test := range tests {
		p, err := parseParams(test.device, test.parity, test.baud, test.databits)
		if !errors.Is(err, test.wantError) {
			t.Fatalf("want: %v got: %v", test.wantError, err)
		}

		if err == nil {
			if !reflect.DeepEqual(p, test.wantParams) {
				t.Errorf("want: %+v, got: %+v", test.wantParams, p)
			}
		}
	}
}
