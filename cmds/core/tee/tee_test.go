// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTee(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		input  string
		args   []string
		append bool
	}{
		{
			"hello",
			[]string{filepath.Join(tempDir, "a1"), filepath.Join(tempDir, "a2")},
			false,
		},
		{
			"a\nb\n",
			[]string{filepath.Join(tempDir, "b1")},
			true,
		},
	}

	for _, test := range tests {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		*cat = test.append
		err := run(strings.NewReader(test.input), &stdout, &stderr, test.args)
		if err != nil {
			t.Error(err)
		}

		for _, name := range test.args {

			b, err := os.ReadFile(name)
			if err != nil {
				t.Error(err)
			}
			res := string(b)

			if res != test.input {
				t.Errorf("want: %q, got %q", test.input, res)
			}

			if res != stdout.String() {
				t.Errorf("want: %q, got %q", test.input, res)
			}
		}
	}
}
