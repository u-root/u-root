// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestHexdump(t *testing.T) {
	for _, tt := range []struct {
		in  []byte
		out []byte
	}{
		{
			in: []byte("abcdefghijklmnopqrstuvwxyz"),
			out: []byte(
				`00000000  61 62 63 64 65 66 67 68  69 6a 6b 6c 6d 6e 6f 70  |abcdefghijklmnop|
00000010  71 72 73 74 75 76 77 78  79 7a                    |qrstuvwxyz|
`),
		},
	} {
		cmd := testutil.Command(t)
		cmd.Stdin = bytes.NewReader(tt.in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(out, tt.out) {
			t.Errorf("want=%#v; got=%#v", tt.out, tt)
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
