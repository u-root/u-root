// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRunGPGV(t *testing.T) {
	for _, tt := range []struct {
		name     string
		keyfile  string
		sigfile  string
		datafile string
		wantErr  error
	}{
		{
			name:     "SuccessSignature",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
		},
		{
			name:     "UsageErrorMissingDatafile",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.sig",
			datafile: "",
			wantErr:  errUsage,
		},
		{
			name:     "UsageErrorMissingSigFile",
			keyfile:  "testdata/key.pub",
			sigfile:  "",
			datafile: "testdata/datafile.txt",
			wantErr:  errUsage,
		},
		{
			name:     "UsageErrorMissingKeyFile",
			keyfile:  "",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			wantErr:  errUsage,
		},
		{
			name:     "UsageError",
			keyfile:  "",
			sigfile:  "",
			datafile: "",
			wantErr:  errUsage,
		},
		{
			name:     "InvalidKeyFile",
			keyfile:  "testdata/datafile.txt",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			wantErr:  fmt.Errorf("tag byte does not have MSB set"),
		},
		{
			name:     "InvalidKeyFilePrivateKey",
			keyfile:  "testdata/private.key",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			wantErr:  errExpectedPacket,
		},
		{
			name:     "InvalidSignatureFile",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.txt",
			datafile: "testdata/datafile.txt",
			wantErr:  fmt.Errorf("tag byte does not have MSB set"),
		},
		{
			name:     "KeyFileIsNotExists",
			keyfile:  "testdata/missing",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			wantErr:  os.ErrNotExist,
		},
		{
			name:     "SigFileIsNotExists",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/missing",
			datafile: "testdata/datafile.txt",
			wantErr:  os.ErrNotExist,
		},
		{
			name:     "DataFileIsNotExists",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/missing",
			wantErr:  os.ErrNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := runGPGV(&buf, tt.keyfile, tt.sigfile, tt.datafile); !errors.Is(err, tt.wantErr) {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("runGPGV(&buf, verbose, args):= %q, want %q", err, tt.wantErr)
				}
			}
		})
	}
}
