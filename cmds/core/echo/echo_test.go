// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func TestEcho(t *testing.T) {
	type test struct {
		s string
		r string
		f flags
	}
	var buf bytes.Buffer
	tests := []test{
		{s: "simple test1", r: "simple test1", f: flags{noNewline: true}},
		{s: "simple test2", r: "simple test2\n", f: flags{}},
		{s: "simple\\ttest3", r: "simple\ttest3\n", f: flags{interpretEscapes: true}},
		{s: "simple\\ttest4", r: "simple\ttest4\n", f: flags{interpretEscapes: true}},
		{s: "simple\\tte\\cst5", r: "simple\tte\n", f: flags{interpretEscapes: true}},
		{s: "simple\\tte\\cst6", r: "simple\tte", f: flags{true, true}},
		{s: "simple\\x56 test7", r: "simpleV test7", f: flags{true, true}},
		{s: "simple\\x56 \\0113test7", r: "simpleV Ktest7", f: flags{true, true}},
		{s: "\\\\8", r: "\\8", f: flags{true, true}},
	}

	for _, v := range tests {
		if err := echo(v.f, &buf, v.s); err != nil {
			t.Errorf("%s", err)
		}
		if string(buf.Bytes()) != v.r {
			t.Fatalf("Want \"%v\", got \"%v\"", v.r, string(buf.Bytes()))
		}
		buf.Reset()
	}
}
