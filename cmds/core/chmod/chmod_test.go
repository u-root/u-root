// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
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
		err        error
	}{
		{
			name: "len(args) < 1",
			err:  errBadUsage,
		},

		{
			name: "len(args) < 2 && *reference",
			args: []string{"arg"},
			err:  errBadUsage,
		},

		{
			name: "file does not exist",
			args: []string{"g-rx", "filedoesnotexist"},
			err:  os.ErrNotExist,
		},
		{
			name: "Value should be less than or equal to 0777",
			args: []string{"7777", f.Name()},
			err:  strconv.ErrRange,
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
			name: "unable to decode mode",
			args: []string{"a=9rwx", f.Name()},
			err:  strconv.ErrSyntax,
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
			args:      []string{"a=rx", f.Name()},
			reference: "filedoesnotexist",
			err:       os.ErrNotExist,
		},
		{
			name:       "correct reference file",
			args:       []string{f.Name()},
			modeBefore: 0o222,
			modeAfter:  0o222,
			reference:  f.Name(),
		},
		{
			name:      "bad filepath",
			args:      []string{"a=rx", "pathdoes not exist"},
			recursive: true,
			err:       os.ErrNotExist,
		},
		{
			name:       "correct path filepath",
			args:       []string{"0777", f.Name()},
			recursive:  true,
			modeBefore: 0o777,
			modeAfter:  0o777,
			err:        nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			os.Chmod(f.Name(), tt.modeBefore)
			err := command(io.Discard, tt.recursive, tt.reference).run(tt.args...)
			if !errors.Is(err, tt.err) {
				t.Errorf("chmod(%v, %q, %q) = %v, want %v", tt.recursive, tt.reference, tt.args, err, tt.err)
				return
			}

			fi, err := os.Stat(f.Name())
			if err != nil {
				t.Fatalf("failed to stat file: %v", err)
			}

			if fi.Mode() != tt.modeAfter {
				t.Errorf("chmod(%v, %q, %q) = mode = %o, want %o", tt.recursive, tt.reference, tt.args, fi.Mode(), tt.modeAfter)
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

	err = command(stderr, false, "").run("0777", f1.Name(), "filenotexists", f2.Name())
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
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

	if stderr.String() != "chmod filenotexists: no such file or directory\n" {
		t.Errorf("expected stderr to be 'chmod filenotexists: no such file or directory', got %q", stderr.String())
	}
}
