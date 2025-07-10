// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/core/mkdir"
	"golang.org/x/sys/unix"
)

type flags struct {
	mode    string
	mkall   bool
	verbose bool
}

func TestMkdir(t *testing.T) {
	d := t.TempDir()
	for _, tt := range []struct {
		name      string
		flags     flags
		args      []string
		wantMode  string
		wantPrint string
		want      error
	}{
		{
			name:     "Create 1 directory",
			flags:    flags{mode: "755"},
			args:     []string{filepath.Join(d, "stub0")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:      "Directory already exists",
			flags:     flags{mode: "755"},
			args:      []string{filepath.Join(d, "stub0")},
			wantMode:  "drwxr-xr-x",
			wantPrint: fmt.Sprintf("%s: %s file exists", filepath.Join(d, "stub0"), filepath.Join(d, "stub0")),
		},
		{
			name: "Create 1 directory verbose",
			flags: flags{
				mode:    "755",
				verbose: true,
			},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Create 2 directories",
			flags:    flags{mode: "755"},
			args:     []string{filepath.Join(d, "stub2"), filepath.Join(d, "stub3")},
			wantMode: "drwxr-xr-x",
		},
		{
			name: "Create a sub directory directly",
			flags: flags{
				mode:  "755",
				mkall: true,
			},
			args:     []string{filepath.Join(d, "stub4"), filepath.Join(d, "stub4/subdir")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:  "Perm Mode Bits over 7 Error",
			flags: flags{mode: "7778"},
			args:  []string{filepath.Join(d, "stub1")},
			want:  fmt.Errorf(`invalid mode "7778"`),
		},
		{
			name:     "More than 4 Perm Mode Bits Error",
			flags:    flags{mode: "11111"},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxrwxr-x",
			want:     fmt.Errorf(`invalid mode "11111"`),
		},
		{
			name:     "Custom Perm in Octal Form",
			flags:    flags{mode: "0777"},
			args:     []string{filepath.Join(d, "stub6")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm not in Octal Form",
			flags:    flags{mode: "777"},
			args:     []string{filepath.Join(d, "stub7")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit",
			flags:    flags{mode: "1777"},
			args:     []string{filepath.Join(d, "stub8")},
			wantMode: "dtrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SGID Bit",
			flags:    flags{mode: "2777"},
			args:     []string{filepath.Join(d, "stub9")},
			wantMode: "dgrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SUID Bit",
			flags:    flags{mode: "4777"},
			args:     []string{filepath.Join(d, "stub10")},
			wantMode: "durwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit and SUID Bit",
			flags:    flags{mode: "5777"},
			args:     []string{filepath.Join(d, "stub11")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name:     "Custom Perm for 2 Directories",
			flags:    flags{mode: "5777"},
			args:     []string{filepath.Join(d, "stub12"), filepath.Join(d, "stub13")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name:     "Default creation mode",
			args:     []string{filepath.Join(d, "stub14")},
			wantMode: "drwxr-xr-x",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := mkdir.New()
			var stdout, stderr bytes.Buffer
			cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

			// don't depend on system umask value, if mode is not specified
			if tt.flags.mode == "" {
				m := unix.Umask(unix.S_IWGRP | unix.S_IWOTH)
				defer func() {
					unix.Umask(m)
				}()
			}

			// Build command arguments
			var args []string
			if tt.flags.mode != "" {
				args = append(args, "-m", tt.flags.mode)
			}
			if tt.flags.mkall {
				args = append(args, "-p")
			}
			if tt.flags.verbose {
				args = append(args, "-v")
			}
			args = append(args, tt.args...)

			got := cmd.Run(args...)
			if got != nil {
				if tt.want == nil || got.Error() != tt.want.Error() {
					t.Errorf("Run() = '%v', want: '%v'", got, tt.want)
				}
			} else {
				if stderr.String() != "" {
					if !strings.Contains(stderr.String(), "file exist") {
						t.Errorf("Stderr = '%v', want to contain 'file exist'", stderr.String())
					}
				}
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
