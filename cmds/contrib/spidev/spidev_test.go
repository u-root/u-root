// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/flash/spimock"
	"github.com/u-root/u-root/pkg/spidev"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name               string
		args               []string
		input              []byte
		ForceOpenErr       error
		ForceTransferErr   error
		ForceSetSpeedHzErr error
		wantTransfers      []spidev.Transfer
		wantSpeed          uint32
		wantOutput         string
		wantOutputRegex    *regexp.Regexp
		wantErr            error
	}{
		{
			name:    "invalid arguments",
			args:    []string{"--invalid", "raw"},
			wantErr: errors.New("unknown flag: --invalid"),
		},
		{
			name:    "invalid subcommand",
			args:    []string{"potato"},
			wantErr: errors.New("unknown subcommand"),
		},
		{
			name:    "too many arguments",
			args:    []string{"raw", "potato"},
			wantErr: errors.New("expected one subcommand"),
		},
		{
			name:         "open error",
			args:         []string{"raw"},
			ForceOpenErr: errors.New("fake open error"),
			wantErr:      errors.New("fake open error"),
		},
		{
			name:             "transfer error",
			args:             []string{"raw"},
			input:            []byte("abcd"),
			ForceTransferErr: errors.New("fake transfer error"),
			wantErr:          errors.New("fake transfer error"),
		},
		{
			name:               "setspeedhz error",
			args:               []string{"raw"},
			input:              []byte("abcd"),
			ForceSetSpeedHzErr: errors.New("fake setspeedhz error"),
			wantErr:            errors.New("fake setspeedhz error"),
		},
		{
			name: "empty transfer",
			args: []string{"raw"},
			// Note wantTransfers is an empty slice. There is no
			// need to even perform an ioctl.
			wantSpeed: 5000000,
		},
		{
			name: "single transfer",
			args: []string{"raw"},
			// This test sends a raw sfdp read command.
			input: []byte{0x5a, 0, 0, 0, 0xff, 0, 0, 0, 0},
			wantTransfers: []spidev.Transfer{
				{
					Tx: []byte{0x5a, 0, 0, 0, 0xff, 0, 0, 0, 0},
					Rx: []byte{0, 0, 0, 0, 0, 'S', 'F', 'D', 'P'},
				},
			},
			wantSpeed:  5000000,
			wantOutput: "\x00\x00\x00\x00\x00SFDP",
		},
		{
			name:            "sfdp",
			args:            []string{"sfdp"},
			wantSpeed:       5000000,
			wantOutputRegex: regexp.MustCompile("FlashMemoryDensity *0x1fffffff"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := spimock.New()
			s.ForceTransferErr = tt.ForceTransferErr
			s.ForceSetSpeedHzErr = tt.ForceSetSpeedHzErr
			openFakeSpi := func(dev string) (spi, error) {
				if tt.ForceOpenErr != nil {
					return nil, tt.ForceOpenErr
				}
				return s, nil
			}

			output := &bytes.Buffer{}
			gotErr := run(tt.args, openFakeSpi, bytes.NewBuffer(tt.input), output)

			if gotErrString, wantErrString := fmt.Sprint(gotErr), fmt.Sprint(tt.wantErr); gotErrString != wantErrString {
				t.Errorf("run() got err %q; want err %q", gotErrString, wantErrString)
			}

			gotOutputString := output.String()
			if tt.wantOutputRegex != nil {
				if !tt.wantOutputRegex.MatchString(gotOutputString) {
					t.Errorf("run() got output %q; want output regex %q", gotOutputString, tt.wantOutputRegex)
				}
			} else if gotOutputString != tt.wantOutput {
				t.Errorf("run() got output %q; want output %q", gotOutputString, tt.wantOutput)
			}

			if tt.wantTransfers != nil && !reflect.DeepEqual(s.Transfers, tt.wantTransfers) {
				t.Errorf("run() got transfers %#v; want transfers %#v", s.Transfers, tt.wantTransfers)
			}
		})
	}
}
