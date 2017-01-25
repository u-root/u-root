// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestEmit(t *testing.T) {

	var buf = []byte("hello\nthis is a test\n")
	var buf2 = []byte("hello\nthiz is a text")

	c1 := make(chan byte, 8192)
	c2 := make(chan byte, 8192)
	c3 := make(chan byte, 8192)

	r := bytes.NewReader(buf)
	err := emit(r, c1, 0)
	if err != io.EOF {
		t.Errorf("%v\n", err)
	}

	r = bytes.NewReader(buf2)
	err = emit(r, c2, 0)
	if err != io.EOF {
		t.Errorf("%v", err)
	}

	err = emit(os.Stdin, c3, 0)
	if err != io.EOF {
		t.Errorf("%v", err)
	}

}
