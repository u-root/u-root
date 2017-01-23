// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func TestSimple(t *testing.T) {
	var dat = bytes.NewReader(testCPIO)

	r, err := NewcReader(dat)

	if err != nil {
		t.Error(err)
	}
	var f *File
	var i int
	for f, err = r.RecRead(); err == nil; f, err = r.RecRead() {
		if f.String() != testResult[i] {
			t.Errorf("Value %d: got \n%s, want \n%s", i, f.String(), testResult[i])
		}
		t.Logf("Value %d: got \n%s, want \n%s", i, f.String(), testResult[i])
		i++
	}
}
