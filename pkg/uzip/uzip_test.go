// Copyright 2018-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uzip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFromZip(t *testing.T) {
	tmpDir := t.TempDir()

	f := filepath.Join(tmpDir, "test.zip")
	if err := ToZip("testdata/testFolder", f, ""); err != nil {
		t.Fatalf(`ToZip("testdata/testFolder", %q, "") = %v, want nil`, f, err)
	}

	z, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf(`os.ReadFile(%q) = %v, want nil`, f, err)
	}
	if len(z) == 0 {
		t.Errorf("len(%v) == %d, want not 0", z, len(z))
	}
	if len(z) < 1 {
		t.Errorf("no content read from file %q", f)
	}

	out := filepath.Join(tmpDir, "unziped")
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		t.Fatalf(`os.MkdirAll(%q, %v)  = %v, want nil`, out, os.ModePerm, err)
	}

	if err := FromZip(f, out); err != nil {
		t.Fatalf(`FromZip(%q, %q) = %v, want nil`, f, out, err)
	}

	f1 := filepath.Join(out, "file1")
	f2 := filepath.Join(out, "file2")
	f3 := filepath.Join(out, "subFolder", "file3")
	f4 := filepath.Join(out, "subFolder", "file4")

	f1Expected, err := os.ReadFile("testdata/testFolder/file1")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/file1") = _, %v, want nil`, err)
	}
	f2Expected, err := os.ReadFile("testdata/testFolder/file2")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/file2") = _, %v, want nil`, err)
	}
	f3Expected, err := os.ReadFile("testdata/testFolder/subFolder/file3")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/subFolder/file3") = _, %v, want nil`, err)
	}
	f4Expected, err := os.ReadFile("testdata/testFolder/subFolder/file4")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/subFolder/file4") = _, %v, want nil`, err)
	}

	var x []byte

	x, err = ioutil.ReadFile(f1)
	if err != nil {
		t.Errorf("open file: %q failed with: %q", f1, err)
	}
	if !bytes.Equal(x, f1Expected) {
		t.Errorf("file %q and file %q are not equal", f1, "testdata/testFolder/file1")
	}
	x, err = ioutil.ReadFile(f2)
	if err != nil {
		t.Errorf("open file: %q failed with: %q", f2, err)
	}
	if !bytes.Equal(x, f2Expected) {
		t.Errorf("file %q and file %q are not equal", f2, "testdata/testFolder/file2")
	}

	x, err = ioutil.ReadFile(f3)
	if err != nil {
		t.Errorf("open file: %q failed with: %q", f3, err)
	}
	if !bytes.Equal(x, f3Expected) {
		t.Errorf("file %q and file %q are not equal", f3, "testdata/testFolder/file3")
	}

	x, err = ioutil.ReadFile(f4)
	if err != nil {
		t.Errorf("open file: %q failed with: %q", f4, err)
	}
	if !bytes.Equal(x, f4Expected) {
		t.Errorf("file %q and file %q are not equal", f4, "testdata/testFolder/file4")
	}
}

func TestFromZipNoValidFile(t *testing.T) {
	f, err := ioutil.TempFile("", "testfile-")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	if err := FromZip(f.Name(), "someDir"); err == nil {
		t.Errorf("FromZip succeeded but shouldn't")
	}
}

func TestAppendZip(t *testing.T) {
	_, err := os.Create("appendTest.zip")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("appendTest.zip")

	if err := AppendZip("testdata/testFolder", "appendTest.zip", "Test append zip"); err != nil {
		t.Error(err)
	}
}

func TestAppendZipNoDir1(t *testing.T) {
	if err := AppendZip("doesNotExist", "alsoNotExist", "Whythough"); err == nil {
		t.Error("AppendZip succeeded but shouldn't")
	}
}

func TestAppendZipNoDir2(t *testing.T) {
	f, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Errorf("creating testfile failed: %v", err)
	}
	defer f.Close()
	if err := AppendZip(f.Name(), f.Name(), "no comment"); err == nil {
		t.Error("AppendZip succeeded but shouldn't")
	}
}

func TestComment(t *testing.T) {
	comment, err := Comment("test.zip")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(comment)
}

func TestToZip(t *testing.T) {
	if err := ToZip(".", "testfile.zip", "test comment"); err != nil {
		t.Error(err)
	}
	defer os.Remove("testfile.zip")
}

func TestToZipInvalidDir(t *testing.T) {
	f, err := ioutil.TempFile("", "testfile-")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	if err := ToZip(f.Name(), "invalid", "no need"); err == nil {
		t.Errorf("ToZip succeeded but shouldn't")
	}
}
