// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
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
