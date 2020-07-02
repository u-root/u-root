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
			desc: "single quote specials",
			in:   `stuff var='more stuff $ \ \$ " '`,
			want: []string{"stuff", `var=more stuff $ \ \$ " `},
		},
		{
			desc: "double quote specials",
			in:   `stuff var="more stuff $ \ \$ \" \n "`,
			want: []string{"stuff", `var=more stuff $ \ $ " \n `},
		},
		{
			desc: "quote forgot close",
			in:   "stuff var='more stuff",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "empty",
			in:   "",
			want: []string{},
		},
		{
			in: `This string has an embedded apostrophe, doesn't it?`,
			want: []string{
				"This",
				"string",
				"has",
				"an",
				"embedded",
				"apostrophe,",
				"doesnt it?",
			},
		},
		{
			in: "This string has embedded \"double quotes\" and 'single quotes' in it,\nand even \"a 'nested example'\".\n",
			want: []string{
				"This",
				"string",
				"has",
				"embedded",
				`double quotes`,
				"and",
				`single quotes`,
				"in",
				"it,",
				"and",
				"even",
				`a 'nested example'.`,
			},
		},
		{
			in: `Hello world!, こんにちは　世界！`,
			want: []string{
				"Hello",
				"world!,",
				"こんにちは",
				"世界！",
			},
		},
		{
			in:   `Do"Not"Separate`,
			want: []string{`DoNotSeparate`},
		},
		{
			in: `Escaped \e Character not in quotes`,
			want: []string{
				"Escaped",
				"e",
				"Character",
				"not",
				"in",
				"quotes",
			},
		},
		{
			in: `Escaped "\e" Character in double quotes`,
			want: []string{
				"Escaped",
				`\e`,
				"Character",
				"in",
				"double",
				"quotes",
			},
		},
		{
			in: `Escaped '\e' Character in single quotes`,
			want: []string{
				"Escaped",
				`\e`,
				"Character",
				"in",
				"single",
				"quotes",
			},
		},
		{
			in: `Escaped '\'' \"\'\" single quote`,
			want: []string{
				"Escaped",
				`\ \"\"`,
				"single",
				"quote",
			},
		},
		{
			in: `Escaped "\"" \'\"\' double quote`,
			want: []string{
				"Escaped",
				`"`,
				`'"'`,
				"double",
				"quote",
			},
		},
		{
			in:   `"'Strip extra layer of quotes'"`,
			want: []string{`'Strip extra layer of quotes'`},
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
