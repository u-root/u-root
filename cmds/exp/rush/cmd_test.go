// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"testing"
)

// TestCmd tests the implementation of builtin commands.
func TestCmd(t *testing.T) {
	tests := []struct {
		name string
		line string
		out  string
		err  error
	}{
		{name: "cd no args", line: "cd", err: errCdUsage},
		{name: "cd bad dir", line: "cd ZARDOX", err: os.ErrNotExist},
		{name: "cd .", line: "cd ."},
		{name: "rushinfo", line: "rushinfo", err: nil, out: "ama"},
	}

	for _, tt := range tests {
		c, _, err := getCommand(bufio.NewReader(bytes.NewReader([]byte(tt.line))))
		if err != nil && !errors.Is(err, tt.err) {
			t.Errorf("%s: getCommand(%q): %v is not %v", tt.name, tt.line, err, tt.err)
			continue
		}
		// We don't test broken parsing here, just that we get some expected
		// arrays
		doArgs(c)
		if err := commands(c); err != nil {
			t.Errorf("commands: %v != nil", err)
			continue
		}
		t.Logf("cmd %q", c)
		// We don't do pipelines in this test.
		// We don't usually care about output.
		o, e := &bytes.Buffer{}, &bytes.Buffer{}
		c[0].Cmd.Stdout, c[0].Cmd.Stderr = o, e
		err = command(c[0])
		t.Logf("Command output: %q, %q", o.String(), e.String())
		if err == nil && tt.err == nil {
			continue
		}
		if err == nil && tt.err != nil {
			t.Errorf("%q: got nil, want %v", c, tt.err)
		}
		if !errors.Is(err, tt.err) {
			t.Errorf("%q: got %v, want %v", c, err, tt.err)
		}
		if len(tt.out) == 0 {
			continue
		}
		if o.Len() == 0 {
			t.Errorf("%q: stdout: got no data, want something", c)
		}
		t.Logf("Command stdout is %s", o.String())
	}
}
