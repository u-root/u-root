// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/term"
)

func TestRun(t *testing.T) {
	tests := []struct {
		err    error
		name   string
		output string
		args   []string
	}{
		{
			name: "no arguments",
			args: []string{"nohup"},
			err:  errUsage,
		},
		{
			name: "invalid command",
			args: []string{"nohup", "invalidcommand"},
			err:  errStart,
		},
		{
			name: "false command",
			args: []string{"nohup", "false"},
			err:  errFinish,
		},
		{
			name:   "valid command",
			args:   []string{"nohup", "echo", "hello"},
			output: "hello\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			err := os.Chdir(dir)
			if err != nil {
				t.Fatalf("can't chdir into %q", dir)
			}

			err = run(test.args)

			if !errors.Is(err, test.err) {
				t.Errorf("expected %v, got %v", test.err, err)
			}

			if test.output != "" && term.IsTerminal(int(os.Stdout.Fd())) {
				b, err := os.ReadFile(filepath.Join(dir, "nohup.out"))
				if err != nil {
					t.Fatalf("can't open nohup.out: %v", err)
				}

				if string(b) != test.output {
					t.Errorf("expected %q, got %q", test.output, string(b))
				}
			}
		})
	}
}
