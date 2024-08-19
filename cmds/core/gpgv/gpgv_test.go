// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestRunGPGV(t *testing.T) {
	for _, tt := range []struct {
		name     string
		keyfile  string
		sigfile  string
		datafile string
		verbose  bool
		wantErr  error
	}{
		{
			name:     "SuccessSignature",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			verbose:  true,
		},
		{
			name:     "UsageErrorMissingDatafile",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.sig",
			datafile: "",
			verbose:  true,
			wantErr:  errUsage,
		},
		{
			name:     "UsageErrorMissingSigFile",
			keyfile:  "testdata/key.pub",
			sigfile:  "",
			datafile: "testdata/datafile.txt",
			verbose:  true,
			wantErr:  errUsage,
		},
		{
			name:     "UsageErrorMissingKeyFile",
			keyfile:  "",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			verbose:  true,
			wantErr:  errUsage,
		},
		{
			name:     "UsageError",
			keyfile:  "",
			sigfile:  "",
			datafile: "",
			verbose:  true,
			wantErr:  errUsage,
		},
		{
			name:     "InvalidKeyFile",
			keyfile:  "testdata/datafile.txt",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			verbose:  true,
			wantErr:  fmt.Errorf("tag byte does not have MSB set"),
		},
		{
			name:     "InvalidKeyFilePrivateKey",
			keyfile:  "testdata/private.key",
			sigfile:  "testdata/datafile.sig",
			datafile: "testdata/datafile.txt",
			verbose:  true,
			wantErr:  fmt.Errorf("openpgp: invalid data: expected first packet to be PublicKey"),
		},
		{
			name:     "InvalidSignatureFile",
			keyfile:  "testdata/key.pub",
			sigfile:  "testdata/datafile.txt",
			datafile: "testdata/datafile.txt",
			verbose:  true,
			wantErr:  fmt.Errorf("tag byte does not have MSB set"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := runGPGV(&buf, tt.verbose, tt.keyfile, tt.sigfile, tt.datafile); !errors.Is(err, tt.wantErr) {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("runGPGV(&buf, verbose, args):= %q, want %q", err, tt.wantErr)
				}
			}
		})
	}
}
