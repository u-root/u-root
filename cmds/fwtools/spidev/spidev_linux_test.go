// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
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
		err                error
	}{
		{
			name: "no args",
			args: []string{""},
			err:  ErrCommand,
		},
		{
			name:            "id",
			args:            []string{"id"},
			wantOutputRegex: regexp.MustCompile("[0-9a-fA-F]*\n"),
		},
		{
			name:             "id failing IO",
			args:             []string{"id"},
			ForceTransferErr: os.ErrInvalid,
			err:              os.ErrInvalid,
		},
		{
			name: "invalid arguments",
			args: []string{"--invalid", "raw"},
			err:  ErrCommand,
		},
		{
			name: "bad hex to raw command",
			args: []string{"raw", "zrg"},
			err:  ErrConvert,
		},
		{
			name: "invalid subcommand",
			args: []string{"potato"},
			err:  ErrCommand,
		},
		{
			name: "too many arguments",
			args: []string{"id", "potato"},
			err:  ErrCommand,
		},
		{
			name:         "open error",
			args:         []string{"raw"},
			ForceOpenErr: os.ErrPermission,
			err:          os.ErrPermission,
		},
		{
			name:             "transfer error",
			args:             []string{"raw"},
			input:            []byte("abcd"),
			ForceTransferErr: os.ErrInvalid,
			err:              os.ErrInvalid,
		},
		{
			name:               "setspeedhz error",
			args:               []string{"-s", "1", "raw"},
			input:              []byte("abcd"),
			ForceSetSpeedHzErr: os.ErrInvalid,
			err:                os.ErrInvalid,
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
			err := run(tt.args, openFakeSpi, bytes.NewBuffer(tt.input), output)

			if !errors.Is(err, tt.err) {
				t.Errorf("run(): %v != %v", err, tt.err)
			}

			if err != nil {
				return
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
