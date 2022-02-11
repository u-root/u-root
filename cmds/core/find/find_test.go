// Copyright 2013-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunFind(t *testing.T) {
	for _, tt := range []struct {
		name       string
		arg        []string
		fargs      flags
		wantErr    error
		wantErrOut string
	}{
		{
			name: "Usage",
			arg:  []string{"", ""},
		},
		{
			name: "Passwd_NoFlags",
			arg:  []string{"/etc/passwd"},
		},
		{
			name: "Passwd_Long",
			arg:  []string{"/etc/passwd"},
			fargs: flags{
				filetype: "file",
			},
		},
		{
			name: "Passwd_File-Debug",
			arg:  []string{"/etc/passwd"},
			fargs: flags{
				filetype: "file",
				debug:    true,
			},
		},
		{
			name: "Passwd_Dir-Debug",
			arg:  []string{"/etc/passwd"},
			fargs: flags{
				filetype: "d",
				debug:    true,
			},
		},
		{
			name: "Invalid Filetype",
			arg:  []string{"/etc/passwd"},
			fargs: flags{
				filetype: "l",
				debug:    true,
			},
			wantErr: errors.New("l is not a valid file type\n valid types are f,file,d,directory"),
		},
		{
			name:       "Find_FileNotExist",
			arg:        []string{"/etc/rancid"},
			wantErrOut: "/etc/rancid: lstat /etc/rancid: no such file or directory\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var inBuf, errBuf bytes.Buffer
			if err := runFind(&inBuf, &errBuf, tt.fargs, tt.arg); !errors.Is(err, tt.wantErr) {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("runFind(out, errOut, tt.args) = %q, want %q", err, tt.wantErr)
				}
				return
			}
			if tt.wantErr != nil {
				return
			}
			if tt.wantErrOut != errBuf.String() {
				t.Errorf("Expected errBuf == %q, not %q", tt.wantErrOut, errBuf.String())
			}
		})
	}
}
