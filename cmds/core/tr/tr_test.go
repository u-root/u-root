// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"
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
		name         string
		args         []string
		input        string
		wantOutput   string
		del          bool
		wantNewError bool
	}{
		{
			name:       "alnum",
			args:       []string{"[:alnum:]", "z"},
			input:      "0123456789!&?defgh",
			wantOutput: "zzzzzzzzzz!&?zzzzz",
		},
		{
			name:       "alpha",
			args:       []string{"[:alpha:]", "z"},
			input:      "0123456789abcdefgh",
			wantOutput: "0123456789zzzzzzzz",
		},
		{
			name:       "digit",
			args:       []string{"[:digit:]", "z"},
			input:      "0123456789abcdefgh",
			wantOutput: "zzzzzzzzzzabcdefgh",
		},
		{
			name:       "lower",
			args:       []string{"[:lower:]", "z"},
			input:      "0123456789abcdEFGH",
			wantOutput: "0123456789zzzzEFGH",
		},
		{
			name:       "upper",
			args:       []string{"[:upper:]", "z"},
			input:      "0123456789abcdEFGH",
			wantOutput: "0123456789abcdzzzz",
		},
		{
			name:       "punct",
			args:       []string{"[:punct:]", "z"},
			input:      "012345*{}[]!.?&()def",
			wantOutput: "012345zzzzzzzzzzzdef",
		},
		{
			name:       "space",
			args:       []string{"[:space:]", "z"},
			input:      "0123456789\t\ncdef",
			wantOutput: "0123456789zzcdef",
		},
		{
			name:       "graph",
			args:       []string{"[:graph:]", "z"},
			input:      "\f\t🔫123456789abcdEFG",
			wantOutput: "\f\tzzzzzzzzzzzzzzzzz",
		},
		{
			name:       "lower_to_upper",
			args:       []string{"[:lower:]", "[:upper:]"},
			input:      "0123456789abcdEFGH",
			wantOutput: "0123456789ABCDEFGH",
		},
		{
			name:       "upper_to_lower",
			args:       []string{"[:upper:]", "[:lower:]"},
			input:      "0123456789abcdEFGH",
			wantOutput: "0123456789abcdefgh",
		},
		{
			name:       "runes_to_runes",
			args:       []string{"39E", "xyz"},
			input:      "0123456789abcdEFGH",
			wantOutput: "012x45678yabcdzFGH",
		},
		{
			name:       "runes_to_runes_truncated",
			args:       []string{"39E", "xy"},
			input:      "0123456789abcdEFGH",
			wantOutput: "012x45678yabcdyFGH",
		},
		{
			name:       "delete_alnum",
			args:       []string{"[:alnum:]", "\uFFFD"},
			input:      "0123456789abcdEFGH",
			wantOutput: "",
		},
		{
			name:       "delete_all_flag_true",
			args:       []string{"[:alnum:]"},
			input:      "12345",
			del:        true,
			wantOutput: "",
		},
		{
			name:       "delete_lower_flag_true",
			args:       []string{"[:lower:]"},
			input:      "abcdeABCz",
			del:        true,
			wantOutput: "ABC",
		},
		{
			name:         "no args",
			args:         nil,
			wantNewError: true,
		},
		{
			name:         "del is true",
			args:         []string{"a", "b"},
			del:          true,
			wantNewError: true,
		},
		{
			name:         "more than two args",
			args:         []string{"a", "b", "c"},
			wantNewError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			in := strings.NewReader(test.input)
			out := &bytes.Buffer{}
			cmd, err := command(in, out, test.args, test.del)
			if test.wantNewError {
				if err == nil {
					t.Fatal("expected error got: nil")
				}
				t.Skip()
			}

			err = cmd.run()
			if err != nil {
				t.Fatal(err)
			}
			res := out.String()

			if test.wantOutput != res {
				t.Errorf("run() want %q, got %q", test.wantOutput, res)
			}
		})
	}
}
