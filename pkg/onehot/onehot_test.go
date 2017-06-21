// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package onehot

import (
	"fmt"
	"io"
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
	if err := f.Close(); err != nil {
		t.Fatalf("Closing %v: want nil, got %v", name, err)
	}
}

func TestOneHotFail(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestOneHot")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)
	bad := filepath.Join(tmpDir, "bad")
	if err := ioutil.WriteFile(bad, []byte{}, 0644); err != nil {
		t.Fatalf("Writing test data: got %v, want nil", err)
	}
	f, err := Open(bad)
	if err != nil {
		t.Fatalf("Opening %v: got %v, want nil", bad, err)
	}
	var b [1]byte
	if _, err := f.Read(b[:]); err == nil {
		t.Fatalf("Writing test data: got nil, want err")
	}

}

func TestOneHotTooMany(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestOneHot")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)
	// This is kinda weird, but basically we make files until we get E2MANY.
	// Then we close the last one that worked, meaning we can only have one open.
	// Then we use onehot to iterate through them all, which should not fail.
	var closeme *os.File
	var names []string
	var nfiles int
	for ; ; nfiles++ {
		n := filepath.Join(tmpDir, fmt.Sprintf("%d", nfiles))
		f, err := os.Create(n)
		t.Logf("%v: %v, %v", n, f, err)
		// I can't believe there's not a better way to do this but ...
		want := fmt.Sprintf("open %v: too many open files", n)
		if err != nil && err.Error() != want {
			t.Fatalf("Creating the %d'th file: want %v, got %v", nfiles, want, err.Error())
		}
		if err != nil {
			break
		}
		// write one byte into it to verify we can read it
		if n, err := f.Write([]byte{1}); err != nil || n != 1 {
			t.Fatalf("Trying to write one byte to %v: want [1, nil] , got [%v, %v]", f, n, err)
		}
		t.Logf("Opened %v", n)
		closeme = f
		names = append(names, n)
	}

	t.Logf("Opened %d files", nfiles)

	if err := closeme.Close(); err != nil {
		t.Fatalf("Closing the closeme file, %v: want nil, got %v", closeme, err)
	}

	var onehots []io.ReadCloser

	for _, n := range names {
		f, err := Open(n)
		t.Logf("onehot open %v", n)
		if err != nil {
			t.Fatalf("Opening %v: got %v, want nil", n, err)
		}
		onehots = append(onehots, f)
	}

	t.Logf("opened them all")
	for _, f := range onehots {
		var b [1]byte
		if n, err := f.Read(b[:]); err != nil || n != 1 {
			t.Fatalf("Reading test data: want [1, nil], got [%v, %v]", n, err)
		}
		t.Logf("Read 1 byte from %v", f)
	}

	// Go around again. This tests interleaved reads.
	for _, f := range onehots {
		var b [1]byte
		if n, err := f.Read(b[:]); err != nil || n != 1 {
			t.Fatalf("Reading test data: want [1, nil], got [%v, %v]", n, err)
		}
		t.Logf("Read 1 byte from %v", f)
	}

	for _, f := range onehots {
		if err := f.Close(); err != nil {
			t.Fatalf("Closing onehot , %v: want nil, got %v", f, err)
		}
	}

	// Go around again. This tests read after close.
	for _, f := range onehots {
		var b [1]byte
		if n, err := f.Read(b[:]); err != io.EOF || n != -1 {
			t.Fatalf("Reading test data: want [-1, io.EOF], got [%v, %v]", n, err)
		}
		t.Logf("Read 1 byte from %v", f)
	}

}
