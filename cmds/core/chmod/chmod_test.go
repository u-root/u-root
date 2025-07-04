// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
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
			args:    []string{"chmod"},
			wantErr: true,
		},

		{
			name:    "len(args) < 2 && *reference",
			args:    []string{"chmod", "arg"},
			wantErr: true,
		},

		{
			name:    "file does not exist",
			args:    []string{"chmod", "g-rx", "filedoesnotexist"},
			wantErr: true,
		},
		{
			name:    "Value should be less than or equal to 0777",
			args:    []string{"chmod", "7777", f.Name()},
			wantErr: true,
		},
		{
			name:       "mode 0777 correct",
			args:       []string{"chmod", "0777", f.Name()},
			modeBefore: 0x000,
			modeAfter:  0o777,
		},
		{
			name:       "mode 0644 correct",
			args:       []string{"chmod", "0644", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o644,
		},
		{
			name:    "unable to decode mode",
			args:    []string{"chmod", "a=9rwx", f.Name()},
			wantErr: true,
		},
		{
			name:       "mode u-rwx correct",
			args:       []string{"chmod", "u-rwx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o077,
		},
		{
			name:       "mode g-rx correct",
			args:       []string{"chmod", "g-rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o727,
		},
		{
			name:       "mode a-xr correct",
			args:       []string{"chmod", "a-xr", f.Name()},
			modeBefore: 0o222,
			modeAfter:  0o222,
		},
		{
			name:       "mode a-xw correct",
			args:       []string{"chmod", "a-xw", f.Name()},
			modeBefore: 0o666,
			modeAfter:  0o444,
		},
		{
			name:       "mode u-xw correct",
			args:       []string{"chmod", "u-xw", f.Name()},
			modeBefore: 0o666,
			modeAfter:  0o466,
		},
		{
			name:       "mode a= correct",
			args:       []string{"chmod", "a=", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o000,
		},
		{
			name:       "mode u= correct",
			args:       []string{"chmod", "u=", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o077,
		},
		{
			name:       "mode u- correct",
			args:       []string{"chmod", "u-", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o777,
		},
		{
			name:       "mode o+ correct",
			args:       []string{"chmod", "o+", f.Name()},
			modeBefore: 0o700,
			modeAfter:  0o700,
		},
		{
			name:       "mode g=rx correct",
			args:       []string{"chmod", "g=rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o757,
		},
		{
			name:       "mode u=rx correct",
			args:       []string{"chmod", "u=rx", f.Name()},
			modeBefore: 0o077,
			modeAfter:  0o577,
		},
		{
			name:       "mode o=rx correct",
			args:       []string{"chmod", "o=rx", f.Name()},
			modeBefore: 0o077,
			modeAfter:  0o075,
		},
		{
			name:       "mode u=xw correct",
			args:       []string{"chmod", "u=xw", f.Name()},
			modeBefore: 0o742,
			modeAfter:  0o342,
		},
		{
			name:       "mode a-rwx correct",
			args:       []string{"chmod", "a-rwx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o000,
		},
		{
			name:       "mode a-rx correct",
			args:       []string{"chmod", "a-rx", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o222,
		},
		{
			name:       "mode a-x correct",
			args:       []string{"chmod", "a-x", f.Name()},
			modeBefore: 0o777,
			modeAfter:  0o666,
		},
		{
			name:       "mode o+rwx correct",
			args:       []string{"chmod", "o+rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o007,
		},
		{
			name:       "mode a+rwx correct",
			args:       []string{"chmod", "a+rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a+xrw correct",
			args:       []string{"chmod", "a+xrw", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a+xxxxxxxx correct",
			args:       []string{"chmod", "a+xxxxxxxx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o111,
		},
		{
			name:       "mode o+xxxxx correct",
			args:       []string{"chmod", "o+xxxxx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o001,
		},
		{
			name:       "mode a+rx correct",
			args:       []string{"chmod", "a+rx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o555,
		},
		{
			name:       "mode a+r correct",
			args:       []string{"chmod", "a+r", f.Name()},
			modeBefore: 0o111,
			modeAfter:  0o555,
		},
		{
			name:       "mode a=rwx correct",
			args:       []string{"chmod", "a=rwx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o777,
		},
		{
			name:       "mode a=rx correct",
			args:       []string{"chmod", "a=rx", f.Name()},
			modeBefore: 0o000,
			modeAfter:  0o555,
		},
		{
			name:      "bad reference file",
			args:      []string{"chmod", "-reference", "filedoesnotexist", f.Name()},
			reference: "filedoesnotexist",
			wantErr:   true,
		},
		{
			name:       "correct reference file",
			args:       []string{"chmod", "-reference", f.Name(), f.Name()},
			modeBefore: 0o222,
			modeAfter:  0o222,
			reference:  f.Name(),
		},
		{
			name:      "bad filepath",
			args:      []string{"chmod", "-recursive", "a=rx", "pathdoes not exist"},
			recursive: true,
			wantErr:   true,
		},
		{
			name:       "correct path filepath",
			args:       []string{"chmod", "-recursive", "0777", f.Name()},
			recursive:  true,
			modeBefore: 0o777,
			modeAfter:  0o777,
		},
		{
			name:       "mode +x correct",
			args:       []string{"chmod", "+x", f.Name()},
			modeBefore: 0o644,
			modeAfter:  0o755,
		},
		{
			name:       "mode -x correct",
			args:       []string{"chmod", "-x", f.Name()},
			modeBefore: 0o755,
			modeAfter:  0o644,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			os.Chmod(f.Name(), tt.modeBefore)

			cmd := chmod.New()
			cmd.SetIO(nil, io.Discard, io.Discard)

			exitCode, err := cmd.Run(context.Background(), tt.args...)

			if tt.wantErr {
				if err == nil && exitCode == 0 {
					t.Errorf("chmod(%q) expected error, got none", tt.args)
				}
				return
			}

			if err != nil {
				t.Errorf("chmod(%q) = %v, want nil", tt.args, err)
				return
			}

			if exitCode != 0 {
				t.Errorf("chmod(%q) = exit code %d, want 0", tt.args, exitCode)
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

	exitCode, err := cmd.Run(context.Background(), "chmod", "0777", f1.Name(), "filenotexists", f2.Name())
	if err == nil && exitCode == 0 {
		t.Errorf("expected error for non-existent file")
	}

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
