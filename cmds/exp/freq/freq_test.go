// Copyright 2013 the u-root Authors. All rights reserved
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

func TestFreq(t *testing.T) {
	t.Run("test stdin", func(t *testing.T) {
		stdin := strings.NewReader("hello\n")
		stdout := &bytes.Buffer{}

		p := params{chr: true}
		c := command(stdin, stdout, nil, p)
		err := c.run()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		expectedOutput := `-        1
e        1
h        1
l        2
o        1
`
		if stdout.String() != expectedOutput {
			t.Errorf("expected %q, got %q", expectedOutput, stdout.String())
		}
	})

	t.Run("test file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "input.txt")
		err := os.WriteFile(path, []byte("hello\n"), 0o644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		stdout := &bytes.Buffer{}
		c := command(nil, stdout, nil, params{utf: true}, path)
		err = c.run()
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expectedOutput := ` 10 012 0a -        1
101 145 65 e        1
104 150 68 h        1
108 154 6c l        2
111 157 6f o        1
`

		if stdout.String() != expectedOutput {
			t.Errorf("expected %q, got %q", expectedOutput, stdout.String())
		}
	})
}
