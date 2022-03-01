// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		want       string
	}{
		{
			name: "len(args) < 1",
		},
		{
			name: "len(args) < 2 && *reference",
			args: []string{"arg"},
		},
		{
			name: "file does not exist",
			args: []string{"g-rx", "filedoesnotexist"},
			want: "stat filedoesnotexist: no such file or directory",
		},
		{
			name: "Value should be less than or equal to 0777",
			args: []string{"7777", f.Name()},
			want: fmt.Sprintf("invalid octal value %0o. Value should be less than or equal to 0777", 0o7777),
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
			want: fmt.Sprintf("unable to decode mode %q. Please use an octal value or a valid mode string", "a=9rwx"),
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
			want:      "bad reference file: stat filedoesnotexist: no such file or directory",
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
			want:      "chmod pathdoes not exist: no such file or directory",
		},
		{
			name:       "correct path filepath",
			args:       []string{"0777", f.Name()},
			recursive:  true,
			modeBefore: 0o777,
			modeAfter:  0o777,
			want:       "chmod pathdoes not exist: no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*recursive = tt.recursive
			*reference = tt.reference
			os.Chmod(f.Name(), tt.modeBefore)
			mode, got := chmod(tt.args...)
			if got != nil {
				if got.Error() != tt.want {
					t.Errorf("chmod() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if mode != tt.modeAfter {
					t.Errorf("chmod() = %v, want: %v", mode, tt.modeAfter)
				}
			}
		})
	}
}
