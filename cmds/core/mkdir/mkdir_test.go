// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestMkdir(t *testing.T) {
	d := t.TempDir()
	for _, tt := range []struct {
		name     string
		flags    []string
		args     []string
		wantMode string
		err      error
	}{
		{
			name:     "Create 1 directory",
			flags:    []string{"-m", "755"},
			args:     []string{filepath.Join(d, "stub0")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Directory already exists",
			flags:    []string{"-m", "755"},
			args:     []string{filepath.Join(d, "stub0")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Create 1 directory verbose",
			flags:    []string{"-m", "755", "-v"},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Create 2 directories",
			flags:    []string{"-m", "755"},
			args:     []string{filepath.Join(d, "stub2"), filepath.Join(d, "stub3")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Create a sub directory directly",
			flags:    []string{"-m", "755", "-p"},
			args:     []string{filepath.Join(d, "stub4"), filepath.Join(d, "stub4/subdir")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:  "Perm Mode Bits over 7 Error",
			flags: []string{"-m", "7778"},
			args:  []string{filepath.Join(d, "stub1")},
			err:   errInvalidMode,
		},
		{
			name:     "More than 4 Perm Mode Bits Error",
			flags:    []string{"-m", "11111"},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxrwxr-x",
			err:      errInvalidMode,
		},
		{
			name:     "Custom Perm in Octal Form",
			flags:    []string{"-m", "0777"},
			args:     []string{filepath.Join(d, "stub6")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm not in Octal Form",
			flags:    []string{"-m", "777"},
			args:     []string{filepath.Join(d, "stub7")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit",
			flags:    []string{"-m", "1777"},
			args:     []string{filepath.Join(d, "stub8")},
			wantMode: "dtrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SGID Bit",
			flags:    []string{"-m", "2777"},
			args:     []string{filepath.Join(d, "stub9")},
			wantMode: "dgrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SUID Bit",
			flags:    []string{"-m", "4777"},
			args:     []string{filepath.Join(d, "stub10")},
			wantMode: "durwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit and SUID Bit",
			flags:    []string{"-m", "5777"},
			args:     []string{filepath.Join(d, "stub11")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name:     "Custom Perm for 2 Directories",
			flags:    []string{"-m", "5777"},
			args:     []string{filepath.Join(d, "stub12"), filepath.Join(d, "stub13")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name: "No dirs",
			err:  errNoDirs,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var ca []string
			ca = append(ca, tt.flags...)
			ca = append(ca, tt.args...)
			var stdout bytes.Buffer

			err := mkdir(&stdout, io.Discard, ca)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}

			if slices.Contains(tt.flags, "-v") {
				for _, arg := range tt.args {
					if !strings.Contains(stdout.String(), arg) {
						t.Errorf("expected to contain %q with verbose flag, got: %s", arg, stdout.String())
					}
				}
			}

			if err == nil {
				for _, name := range tt.args {
					if stat, err := os.Stat(name); err == nil {
						if stat.Mode().String() != tt.wantMode {
							t.Errorf("Mode = '%v', want: '%v'", stat.Mode().String(), tt.wantMode)
						}
					}
				}
			}
		})
	}
}
