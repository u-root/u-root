// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
	"unicode"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestUnescape(t *testing.T) {
	if _, err := unescape("a\\tb\\nc\\d"); err == nil {
		t.Errorf("unescape() expected error, got nil")
	}
	got, err := unescape("a\\tb\\nc\\\\")
	if err != nil {
		t.Fatalf("unescape() error: %v", err)
	}
	want := "a\tb\nc\\"
	if string(got) != want {
		t.Errorf("unescape() want %q, got %q", want, got)
	}
}

func TestTR(t *testing.T) {
	for _, test := range []struct {
		name   string
		input  string
		output string
		t      *transformer
	}{
		{
			name:   "alnum",
			input:  "0123456789!&?defgh",
			output: "zzzzzzzzzz!&?zzzzz",
			t:      setToRune(ALNUM, 'z'),
		},
		{
			name:   "alpha",
			input:  "0123456789abcdefgh",
			output: "0123456789zzzzzzzz",
			t:      setToRune(ALPHA, 'z'),
		},
		{
			name:   "digit",
			input:  "0123456789abcdefgh",
			output: "zzzzzzzzzzabcdefgh",
			t:      setToRune(DIGIT, 'z'),
		},
		{
			name:   "lower",
			input:  "0123456789abcdEFGH",
			output: "0123456789zzzzEFGH",
			t:      setToRune(LOWER, 'z'),
		},
		{
			name:   "upper",
			input:  "0123456789abcdEFGH",
			output: "0123456789abcdzzzz",
			t:      setToRune(UPPER, 'z'),
		},
		{
			name:   "punct",
			input:  "012345*{}[]!.?&()def",
			output: "012345zzzzzzzzzzzdef",
			t:      setToRune(PUNCT, 'z'),
		},
		{
			name:   "space",
			input:  "0123456789\t\ncdef",
			output: "0123456789zzcdef",
			t:      setToRune(SPACE, 'z'),
		},
		{
			name:   "graph",
			input:  "\f\t🔫123456789abcdEFG",
			output: "\f\tzzzzzzzzzzzzzzzzz",
			t:      setToRune(GRAPH, 'z'),
		},
		{
			name:   "lower_to_upper",
			input:  "0123456789abcdEFGH",
			output: "0123456789ABCDEFGH",
			t:      lowerToUpper(),
		},
		{
			name:   "upper_to_lower",
			input:  "0123456789abcdEFGH",
			output: "0123456789abcdefgh",
			t:      upperToLower(),
		},
		{
			name:   "runes_to_runes",
			input:  "0123456789abcdEFGH",
			output: "012x45678yabcdzFGH",
			t:      runesToRunes([]rune("39E"), 'x', 'y', 'z'),
		},
		{
			name:   "runes_to_runes_truncated",
			input:  "0123456789abcdEFGH",
			output: "012x45678yabcdyFGH",
			t:      runesToRunes([]rune("39E"), 'x', 'y'),
		},
		{
			name:   "delete_alnum",
			input:  "0123456789abcdEFGH",
			output: "",
			t:      setToRune(ALNUM, unicode.ReplacementChar),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			test.t.run(bytes.NewBufferString(test.input), out)
			res := out.String()
			if test.output != res {
				t.Errorf("run() want %q, got %q", test.output, res)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
