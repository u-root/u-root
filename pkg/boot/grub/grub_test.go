// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	for i, tt := range []struct {
		desc string
		in   string
		want []string
	}{
		{
			desc: "nothing to do",
			in:   "stuff",
			want: []string{"stuff"},
		},
		{
			desc: "split",
			in:   "stuff more stuff",
			want: []string{"stuff", "more", "stuff"},
		},
		{
			desc: "escape",
			in:   "stuff\\ more stuff",
			want: []string{"stuff more", "stuff"},
		},
		{
			desc: "quote",
			in:   "stuff var='more stuff'",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "double quote",
			in:   "stuff var=\"more stuff\"",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "quote specials",
			in:   "stuff var='more stuff $ \\ \\$ \" '",
			want: []string{"stuff", "var=more stuff $ \\ \\$ \" "},
		},
		{
			desc: "double quote",
			in:   "stuff var=\"more stuff $ \\ \\$ \\\" \"",
			want: []string{"stuff", "var=more stuff $ \\ $ \" "},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			got := fields(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fields = %#v, want %#v", got, tt.want)
			}

		})
	}
}

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
