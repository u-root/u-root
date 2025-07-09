// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/core/chmod"
)

func TestChmod(t *testing.T) {
	f, err := os.Create(filepath.Join(t.TempDir(), "tmpfile"))
	if err != nil {
		t.Errorf("Failed to create tmp file, %v", err)
	}
	for _, tt := range []struct {
		name       string
		args       []string
		recursive  bool
		reference  string
		modeBefore os.FileMode
		modeAfter  os.FileMode
		wantErr    bool
	}{
		{
			name:    "len(args) < 1",
			args:    []string{},
			wantErr: true,
		},

		{
			name:    "len(args) < 2 && *reference",
			args:    []string{"arg"},
			wantErr: true,
		},

		{
			name:    "file does not exist",
			args:    []string{"g-rx", "filedoesnotexist"},
			wantErr: true,
		},
		{
			name:    "Value should be less than or equal to 0777",
			args:    []string{"7777", f.Name()},
			wantErr: true,
		},
		{
			name:       "mode 0777 correct",
			args:       []string{"0777", f.Name()},
			modeBefore: 0x000,
			modeAfter:  0o777,
		},
		{
			name:       "mode 0644 correct",
			args:       []string{"0644", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o644,
		},
		{
			name:    "unable to decode mode",
			args:    []string{"a=9rwx", f.Name()},
			wantErr: true,
		},
		{
			name:       "mode u-rwx correct",
			args:       []string{"u-rwx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o077,
		},
		{
			name:       "mode g-rx correct",
			args:       []string{"g-rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o727,
		},
		{
			name:       "mode a-xr correct",
			args:       []string{"a-xr", f.Name()},
			modeBefore: 0o222,
			modeAfter:  0o222,
		},
		{
			name:       "mode a-xw correct",
			args:       []string{"a-xw", f.Name()},
			modeBefore: 0o666,
			modeAfter:  0o444,
		},
		{
			name:       "mode u-xw correct",
			args:       []string{"u-xw", f.Name()},
			modeBefore: 0o666,
			modeAfter:  0o466,
		},
		{
			name:       "mode a= correct",
			args:       []string{"a=", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o000,
		},
		{
			name:       "mode u= correct",
			args:       []string{"u=", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o077,
		},
		{
			name:       "mode u- correct",
			args:       []string{"u-", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o777,
		},
		{
			name:       "mode o+ correct",
			args:       []string{"o+", f.Name()},
			modeBefore: 0o700,
			modeAfter:  0o700,
		},
		{
			name:       "mode g=rx correct",
			args:       []string{"g=rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o757,
		},
		{
			name:       "mode u=rx correct",
			args:       []string{"u=rx", f.Name()},
			modeBefore: 0o077,
			modeAfter:  0o577,
		},
		{
			name:       "mode o=rx correct",
			args:       []string{"o=rx", f.Name()},
			modeBefore: 0o077,
			modeAfter:  0o075,
		},
		{
			name:       "mode u=xw correct",
			args:       []string{"u=xw", f.Name()},
			modeBefore: 0o742,
			modeAfter:  0o342,
		},
		{
			name:       "mode a-rwx correct",
			args:       []string{"a-rwx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o000,
		},
		{
			name:       "mode a-rx correct",
			args:       []string{"a-rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o222,
		},
		{
			name:       "mode a-x correct",
			args:       []string{"a-x", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o666,
		},
		{
			name:       "mode o+rwx correct",
			args:       []string{"o+rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o007,
		},
		{
			name:       "mode a+rwx correct",
			args:       []string{"a+rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a+xrw correct",
			args:       []string{"a+xrw", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a+xxxxxxxx correct",
			args:       []string{"a+xxxxxxxx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o111,
		},
		{
			name:       "mode o+xxxxx correct",
			args:       []string{"o+xxxxx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o001,
		},
		{
			name:       "mode a+rx correct",
			args:       []string{"a+rx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o555,
		},
		{
			name:       "mode a+r correct",
			args:       []string{"a+r", f.Name()},
			modeBefore: 0o111,
			modeAfter:  0o555,
		},
		{
			name:       "mode a=rwx correct",
			args:       []string{"a=rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a=rx correct",
			args:       []string{"a=rx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o555,
		},
		{
			name:      "bad reference file",
			args:      []string{"-reference", "filedoesnotexist", f.Name()},
			reference: "filedoesnotexist",
			wantErr:   true,
		},
		{
			name:       "correct reference file",
			args:       []string{"-reference", f.Name(), f.Name()},
			modeBefore: 0o222,
			modeAfter:  0o222,
			reference:  f.Name(),
		},
		{
			name:      "bad filepath",
			args:      []string{"-recursive", "a=rx", "pathdoes not exist"},
			recursive: true,
			wantErr:   true,
		},
		{
			name:       "correct path filepath",
			args:       []string{"-recursive", "0777", f.Name()},
			recursive:  true,
			modeBefore: 0o777,
			modeAfter:  0o777,
		},
		{
			name:       "mode +x correct",
			args:       []string{"+x", f.Name()},
			modeBefore: 0o644,
			modeAfter:  0o755,
		},
		{
			name:       "mode -x correct",
			args:       []string{"-x", f.Name()},
			modeBefore: 0o755,
			modeAfter:  0o644,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			os.Chmod(f.Name(), tt.modeBefore)

			cmd := chmod.New()
			cmd.SetIO(nil, io.Discard, io.Discard)

			err := cmd.Run(tt.args...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("chmod(%q) expected error, got none", tt.args)
				}
				return
			}

			if err != nil {
				t.Errorf("chmod(%q) = %v, want nil", tt.args, err)
				return
			}

			fi, err := os.Stat(f.Name())
			if err != nil {
				t.Fatalf("failed to stat file: %v", err)
			}

			if fi.Mode() != tt.modeAfter {
				t.Errorf("chmod(%q) = mode = %o, want %o", tt.args, fi.Mode(), tt.modeAfter)
			}
		})
	}
}

func TestMultipleFiles(t *testing.T) {
	f1, err := os.Create(filepath.Join(t.TempDir(), "tmpfile1"))
	if err != nil {
		t.Fatalf("Failed to create tmp file, %v", err)
	}

	f2, err := os.Create(filepath.Join(t.TempDir(), "tmpfile2"))
	if err != nil {
		t.Fatalf("Failed to create tmp file, %v", err)
	}

	stderr := &bytes.Buffer{}

	cmd := chmod.New()
	cmd.SetIO(nil, io.Discard, stderr)

	_ = cmd.Run("0777", f1.Name(), "filenotexists", f2.Name())

	// but file1 and file2 should have been chmod'ed
	fi, err := os.Stat(f1.Name())
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if fi.Mode() != 0o777 {
		t.Errorf("chmod(%q) = %o, want %o", f1.Name(), fi.Mode(), 0o777)
	}

	fi, err = os.Stat(f2.Name())
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if fi.Mode() != 0o777 {
		t.Errorf("chmod(%q) = %o, want %o", f2.Name(), fi.Mode(), 0o777)
	}
}
