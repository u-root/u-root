// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package onehot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const tval = "testing"

func TestOneHotFileClose(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestOneHot")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	name := filepath.Join(tmpDir, "test")
	if err := ioutil.WriteFile(name, []byte(tval), 0644); err != nil {
		t.Fatalf("Writing test data: got %v, want nil", err)
	}

	f, err := Open(name)
	if err != nil {
		t.Fatalf("Opening %v after writing it: got %v, want nil", name, err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Closing %v after opening it: got %v, want nil", name, err)
	}
	err = f.Close()
	if err == nil {
		t.Fatalf("Closing %v after closing it: got nil, want err", name)
	}
}

func TestOneHotFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestOneHot")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	name := filepath.Join(tmpDir, "test")
	if err := ioutil.WriteFile(name, []byte(tval), 0644); err != nil {
		t.Fatalf("Writing test data: got %v, want nil", err)
	}

	f, err := Open(name)
	if err != nil {
		t.Fatalf("Opening %v after writing it: got %v, want nil", name, err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Writing test data: got %v, want nil", err)
	}
	if len(b) != 7 {
		t.Fatalf("Reading back file: got len %d, want %d", len(b), len(tval))
	}
	if string(b) != tval {
		t.Fatalf("Reading back file: got %s, want %s", string(b), tval)
	}
}

func TestOneHotFail(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestOneHot")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)
	bad := filepath.Join(tmpDir, "bad")
	f, err := Open(bad)
	if err != nil {
		t.Fatalf("Opening %v: got %v, want nil", err)
	}
	var b [1]byte
	if _, err := f.Read(b[:]); err == nil {
		t.Fatalf("Writing test data: got %v, want nil", err)
	}

}
