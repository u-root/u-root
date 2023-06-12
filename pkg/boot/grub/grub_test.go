// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"errors"
	"fmt"
	"testing"
)

func TestCmdlineQuote(t *testing.T) {
	for i, tt := range []struct {
		desc string
		in   []string
		want string
	}{
		{
			desc: "nothing to do",
			in:   []string{"stuff"},
			want: "stuff",
		},
		{
			desc: "split",
			in:   []string{"stuff", "more", "stuff"},
			want: "stuff more stuff",
		},
		{
			desc: "escape",
			in:   []string{`escape\quote'double"`},
			want: `escape\\quote\'double\"`,
		},
		{
			desc: "quote spaced",
			in:   []string{`some stuff`},
			want: `"some stuff"`,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			got := cmdlineQuote(tt.in)
			if got != tt.want {
				t.Errorf("cmdlineQuote = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFindkeywordGrubEnv(t *testing.T) {
	file := "EFI/centos/grubenv"
	fsRoot := "testdata_new/CentOS_8_Stream_x86_64_blscfg_sda1"

	tests := []struct {
		name string
		key  string
		want string
		err  error
	}{
		{
			name: "Return correct Grubenv value",
			key:  "saved_entry",
			want: "9af7b02ac08149d985841c07c8ff366e-5.18.0",
			err:  nil,
		},
		{
			name: "Return error key",
			key:  "saved_en",
			want: "",
			err:  errMissingKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			got, err := findkeywordGrubEnv(file, fsRoot, tt.key)
			if got != tt.want || !errors.Is(err, tt.err) {
				t.Errorf("findkeywordGrubEnv() = %v, want %v : %v", got, tt.want, err)
			}
		})
	}
}
