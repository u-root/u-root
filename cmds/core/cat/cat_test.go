// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// by Rafael Campos Nunes <rafaelnunes@engineer.com>

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// setup writes a set of files, putting 1 byte in each file.
func setup(t *testing.T, data []byte) (string, error) {
	t.Logf(":: Creating simulation data. ")
	dir, err := ioutil.TempDir("", "cat.dir")
	if err != nil {
		return "", err
	}

	for i, d := range data {
		n := fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i)
		if err := ioutil.WriteFile(n, []byte{d}, 0666); err != nil {
			return "", err
		}
	}

	return dir, nil
}

// TestCat test cat function against 4 files, in each file it is written a bit of someData
// array and the test expect the cat to return the exact same bit from someData array with
// the corresponding file.
func TestCat(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir, err := setup(t, someData)
	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}
	defer os.RemoveAll(dir)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}

	var b bytes.Buffer
	if err := cat(&b, files); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", b.Bytes(), someData)
	}
}
