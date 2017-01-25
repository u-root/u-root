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
		s         string
		nonewline bool
	}
	tests := []test{{s: "simple\ttest", nonewline: false}, {s: "simple\ttest\t2", nonewline: true}}
	bufs := make([]bytes.Buffer, len(tests))

	for i, v := range tests {
		if err := echo(&bufs[i], v.s); err != nil {
			t.Errorf("%s", err)
		}
		if !*nonewline {
			v.s = v.s + "\n"
		}
		if string(bufs[i].Bytes()) != v.s {
			t.Fatalf("Want %v, got %v", v.s, string(bufs[i].Bytes()))
		}
	}
}
