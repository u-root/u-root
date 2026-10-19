// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp/packet"
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

// TestRunGPGVRejectsCertification feeds the key's user ID self-certification
// as a detached signature over the bytes it certifies.
func TestRunGPGVRejectsCertification(t *testing.T) {
	f, err := os.Open("testdata/key.pub")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var (
		pk  *packet.PublicKey
		uid *packet.UserId
		sig *packet.Signature
	)
	r := packet.NewReader(f)
	for sig == nil {
		p, err := r.Next()
		if err != nil {
			t.Fatalf("no user ID certification in key: %v", err)
		}
		switch p := p.(type) {
		case *packet.PublicKey:
			if pk == nil {
				pk = p
			}
		case *packet.UserId:
			uid = p
		case *packet.Signature:
			if uid != nil {
				sig = p
			}
		}
	}

	// The certified data: the primary key followed by the user ID.
	var content bytes.Buffer
	if err := pk.SerializeForHash(&content); err != nil {
		t.Fatal(err)
	}
	var hdr [5]byte
	hdr[0] = 0xb4
	binary.BigEndian.PutUint32(hdr[1:], uint32(len(uid.Id)))
	content.Write(hdr[:])
	content.WriteString(uid.Id)

	var sigBuf bytes.Buffer
	if err := sig.Serialize(&sigBuf); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	datafile := filepath.Join(dir, "data")
	sigfile := filepath.Join(dir, "data.sig")
	if err := os.WriteFile(datafile, content.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sigfile, sigBuf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := runGPGV(&buf, "testdata/key.pub", sigfile, datafile); !errors.Is(err, errNotDetachedSig) {
		t.Errorf("runGPGV with certification sig = %v (output %q), want %v", err, buf.String(), errNotDetachedSig)
	}
}
