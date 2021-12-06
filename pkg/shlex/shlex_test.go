// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shlex

import (
	"fmt"
	"reflect"
	"testing"
)

func TestArgv(t *testing.T) {
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
		{
			desc: "quote forgot close",
			in:   "stuff var='more stuff",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "empty",
			in:   "",
			want: nil,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			got := Argv(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Argv = %#v, want %#v", got, tt.want)
			}
		})
	}
}
