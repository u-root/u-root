// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintenv(t *testing.T) {
	var buf bytes.Buffer

	want := os.Environ()

	printenv(&buf)

	found := strings.Split(buf.String(), "\n")

	for i, v := range want {
		if v != found[i] {
			t.Fatalf("want %s, got %s", v, found[i])
		}
	}
}
