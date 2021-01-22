// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/spi"
)

// fakeSpidev toggles the ASCII case of each byte in the transfer.
type fakeSpidev struct {
	fakeTransferErr error
	transfers       []spi.Transfer
}

func (f *fakeSpidev) Transfer(t []spi.Transfer) error {
	if f.fakeTransferErr != nil {
		return f.fakeTransferErr
	}
	n := copy(t[0].Rx, t[0].Tx)
	for i := 0; i < n; i++ {
		// This toggles the ASCII case. Don't think too hard about it
		// -- we need to apply a transformation to fake some sort of
		// SPI device.
		t[0].Rx[i] ^= 1 << 5
	}
	f.transfers = append(f.transfers, t...)
	return nil
}

func (f *fakeSpidev) Close() error {
	return nil
}

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name            string
		args            []string
		input           []byte
		fakeOpenErr     error
		fakeTransferErr error
		wantTransfers   []spi.Transfer
		wantOutput      string
		wantErr         error
	}{
		{
			name:    "invalid arguments",
			args:    []string{"--invalid", "raw"},
			wantErr: errors.New("unknown flag: --invalid"),
		},
		{
			name:    "invalid subcommand",
			args:    []string{"potato"},
			wantErr: errors.New("expected 'raw' subcommand"),
		},
		{
			name:    "too many arguments",
			args:    []string{"raw", "potato"},
			wantErr: errors.New("expected 'raw' subcommand"),
		},
		{
			name:        "open error",
			args:        []string{"raw"},
			fakeOpenErr: errors.New("fake open error"),
			wantErr:     errors.New("fake open error"),
		},
		{
			name:            "transfer error",
			args:            []string{"raw"},
			input:           []byte("abcd"),
			fakeTransferErr: errors.New("fake transfer error"),
			wantErr:         errors.New("fake transfer error"),
		},
		{
			name: "empty transfer",
			args: []string{"raw"},
			// Note wantTransfers is an empty slice. There is no
			// need to even perform an ioctl.
		},
		{
			name:  "single transfer",
			args:  []string{"raw"},
			input: []byte("abcdefg"),
			wantTransfers: []spi.Transfer{
				{
					Tx:       []byte("abcdefg"),
					Rx:       []byte("ABCDEFG"),
					CSChange: true,
					SpeedHz:  500000,
				},
			},
			wantOutput: "ABCDEFG",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeSpidev{
				fakeTransferErr: tt.fakeTransferErr,
			}
			openFakeSpi := func(dev string) (spidev, error) {
				if tt.fakeOpenErr != nil {
					return nil, tt.fakeOpenErr
				}
				return fake, nil
			}

			output := &bytes.Buffer{}
			gotErr := run(tt.args, openFakeSpi, bytes.NewBuffer(tt.input), output)

			if gotErrString, wantErrString := fmt.Sprint(gotErr), fmt.Sprint(tt.wantErr); gotErrString != wantErrString {
				t.Errorf("run() got err %q; want err %q", gotErrString, wantErrString)
			}

			if gotOutputString := output.String(); gotOutputString != tt.wantOutput {
				t.Errorf("run() got output %q; want output %q", gotOutputString, tt.wantOutput)
			}

			if !reflect.DeepEqual(fake.transfers, tt.wantTransfers) {
				t.Errorf("run() got transfers %#v; want transfers %#v", fake.transfers, tt.wantTransfers)
			}
		})
	}
}
