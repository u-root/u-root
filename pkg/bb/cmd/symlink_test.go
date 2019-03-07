// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSymlink(t *testing.T) {
	d, err := ioutil.TempDir("", "u-root-test-")
	if err != nil {
		t.Fatal(err)
	}

	foo, err := os.Create(filepath.Join(d, "foo"))
	if err != nil {
		t.Fatal(err)
	}
	foo.Close()

	// /abs/parent/baz -> ../parent/bat -> bar -> foo
	os.Symlink("foo", filepath.Join(d, "bar"))
	os.Symlink("bar", filepath.Join(d, "bat"))
	os.Symlink(fmt.Sprintf("../%s/bat", filepath.Base(d)), filepath.Join(d, "baz"))

	want := filepath.Join(d, "bar")
	if got := resolveUntilLastSymlink(filepath.Join(d, "baz")); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
