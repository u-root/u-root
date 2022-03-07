// Copyright 2021-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package main

import (
	"os"
	"path/filepath"
	"testing"

	"src.elv.sh/pkg/prog"
)

func TestElvish(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		status int
	}{
		{
			name:   "success",
			args:   []string{""},
			status: 0,
		},
		{
			name:   "failure",
			args:   []string{"foo", "bar"},
			status: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			inFile, err := os.Create(filepath.Join(tmp, "inFile"))
			if err != nil {
				t.Errorf("Failed to create file: %v", err)
			}
			outFile, err := os.Create(filepath.Join(tmp, "outFile"))
			if err != nil {
				t.Errorf("Failed to create file: %v", err)
			}

			if result := run(inFile, outFile, outFile, tt.args); result != tt.status {
				t.Errorf("Want: %d, Got: %d", tt.status, result)
			}
		})
	}
}

func TestDaemonStub(t *testing.T) {
	for _, tt := range []struct {
		name    string
		daemon  bool
		wantErr error
	}{
		{
			name:    "nodaemon",
			daemon:  false,
			wantErr: prog.ErrNotSuitable,
		},
		{
			name:    "daemon",
			daemon:  true,
			wantErr: ErrNotSupported,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := daemonStub{}
			if err := d.Run([3]*os.File{}, &prog.Flags{Daemon: tt.daemon}, []string{}); err != nil {
				if err != tt.wantErr {
					t.Errorf("Want: %q, Got: %v", tt.wantErr, err)
				}
			}
		})
	}
}
