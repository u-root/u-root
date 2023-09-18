// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

// Test echo cmd
func TestEcho(t *testing.T) {
	// Creating test table for test
	for _, tt := range []struct {
		name                      string
		input                     string
		want                      string
		noNewline                 bool
		interpretEscapes          bool
		interpretBackslashEscapes bool
	}{
		{
			name:      "test1 noNewline",
			input:     "simple test1",
			want:      "simple test1",
			noNewline: true,
		},
		{
			name:  "test2",
			input: "simple test2",
			want:  "simple test2\n",
		},
		{
			name:             "test3 interpretEscapes",
			input:            "simple\\ttest3",
			want:             "simple\ttest3\n",
			interpretEscapes: true,
		},
		{
			name:             "test4 interpretEscapes",
			input:            "simple\\ttest4",
			want:             "simple\ttest4\n",
			interpretEscapes: true,
		},
		{
			name:             "test5 interpretEscapes",
			input:            "simple\\tte\\cst5",
			want:             "simple\tte\n",
			interpretEscapes: true,
		},
		{
			name:             "test6 interpretEscapes and noNewline",
			input:            "simple\\tte\\cst6",
			want:             "simple\tte",
			noNewline:        true,
			interpretEscapes: true,
		},
		{
			name:             "test7 interpretEscapes and noNewline",
			input:            "simple\\x56 test7",
			want:             "simpleV test7",
			noNewline:        true,
			interpretEscapes: true,
		},
		{
			name:             "test8 interpretEscapes and noNewline",
			input:            "simple\\x56 \\0113test8",
			want:             "simpleV Ktest8",
			noNewline:        true,
			interpretEscapes: true,
		},
		{
			name:             "test9 interpretEscapes and noNewline",
			input:            "\\\\9",
			want:             "\\9",
			noNewline:        true,
			interpretEscapes: true,
		},
		{
			name:             "test10 empty String",
			input:            "",
			want:             "",
			noNewline:        true,
			interpretEscapes: true,
		},
		{
			name:                      "test11 interpretBackslashEscapes true",
			input:                     "simple test11",
			want:                      "simple test11\n",
			interpretBackslashEscapes: true,
		},
	} {
		// Run tests
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			if err := echo(&got, tt.noNewline, tt.interpretEscapes, tt.interpretBackslashEscapes, tt.input); err != nil {
				t.Errorf("%q failed: %q", tt.name, err)
			}
			if got.String() != tt.want {
				t.Fatalf("%q failed. Got: %q, Want: %q", tt.name, got.String(), tt.want)
			}
		})
	}
}
