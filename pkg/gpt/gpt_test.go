// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package gpt

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

var (
	disk = make([]byte, 0x100000000)
)

func InstallGPT() {
	for i, d := range block {
		copy(disk[i:], d)
	}
}

// TestGPTTtables tests whether we can match the primary and backup
// or, if they differ, we catch that error.
// We know from other tests that the tables read fine.
// This test verifies that they match and that therefore we
// are able to read the backup table and test that it is ok.
func TestGPTTables(t *testing.T) {
	var tests = []struct {
		mangle int
		what   string
	}{
		{-1, "No error test"},
		{0x10, "Should differ test"},
	}

	for _, test := range tests {
		InstallGPT()
		if test.mangle > -1 {
			disk[BlockSize+test.mangle] = 0xff
		}

		r := bytes.NewReader(disk)
		_, err := New(r)
		switch {
		case err != nil && test.mangle > -1:
			t.Logf("Got expected error %s", test.what)
		case err != nil && test.mangle == -1:
			t.Errorf("%s: got %s, want nil", test.what, err)
			continue
		case err == nil && test.mangle > -1:
			t.Errorf("%s: got nil, want err", test.what)
			continue
		}
		t.Logf("Passed %s", test.what)
	}
}

type iodisk []byte

func (d *iodisk) WriteAt(b []byte, off int64) (int, error) {
	copy([]byte(*d)[off:], b)
	return len(b), nil
}

func TestWrite(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	p, err := New(r)
	if err != nil {
		t.Fatalf("Reading partitions: got %v, want nil", err)
	}
	var targ = make(iodisk, len(disk))

	if err := Write(&targ, p); err != nil {
		t.Fatalf("Writing: got %v, want nil", err)
	}
	if n, err := New(bytes.NewReader([]byte(targ))); err != nil {
		t.Logf("Old GPT: %v", p.GPT)
		var b bytes.Buffer
		w := hex.Dumper(&b)
		io.Copy(w, bytes.NewReader(disk[:4096]))
		t.Logf("%s\n", b.String())
		t.Fatalf("Reading back new header: new:%s\n%v", n, err)
	}
}
