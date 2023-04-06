// Copyright 2018-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uzip

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFromZip(t *testing.T) {
	tmpDir := t.TempDir()

	f := filepath.Join(tmpDir, "test.zip")
	if err := ToZip("testdata/testFolder", f, ""); err != nil {
		t.Fatalf(`ToZip("testdata/testFolder", %q, "") = %q, not nil`, f, err)
	}

	z, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf(`os.ReadFile(%q) = %q, not nil`, f, err)
	}
	if len(z) == 0 {
		t.Errorf("len(z) = 0 , want > 0")
	}
	if len(z) < 1 {
		t.Errorf(`len(z) < 1, not > 0`)
	}

	out := filepath.Join(tmpDir, "unziped")
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		t.Fatalf(`os.MkdirAll(%q, %v)  = %q, not nil`, out, os.ModePerm, err)
	}

	if err := FromZip(f, out); err != nil {
		t.Fatalf(`FromZip(%q, %q) = %q, not nil`, f, out, err)
	}

	f1 := filepath.Join(out, "file1")
	f2 := filepath.Join(out, "file2")
	f3 := filepath.Join(out, "subFolder", "file3")
	f4 := filepath.Join(out, "subFolder", "file4")

	f1Expected, err := os.ReadFile("testdata/testFolder/file1")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/file1") = _, %q, not _, nil`, err)
	}
	f2Expected, err := os.ReadFile("testdata/testFolder/file2")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/file2") = _, %q, not _, nil`, err)
	}
	f3Expected, err := os.ReadFile("testdata/testFolder/subFolder/file3")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/subFolder/file3") = _, %q, not _, nil`, err)
	}
	f4Expected, err := os.ReadFile("testdata/testFolder/subFolder/file4")
	if err != nil {
		t.Fatalf(`os.ReadFile("testdata/testFolder/subFolder/file4") = _, %q, not _, nil`, err)
	}

	var x []byte

	x, err = os.ReadFile(f1)
	if err != nil {
		t.Errorf("ioutil.ReadFile(%q) = %q, not nil", f1, err)
	}
	if !bytes.Equal(x, f1Expected) {
		t.Logf("\nGot:\t %v\nWant:\t %v", x[:30], f1Expected[:30])
		t.Errorf("file %q and file %q are not equal", f1, "testdata/testFolder/file1")
	}
	x, err = os.ReadFile(f2)
	if err != nil {
		t.Errorf("ioutil.ReadFile(%q) = %q, not nil", f2, err)
	}
	if !bytes.Equal(x, f2Expected) {
		t.Logf("\nGot:\t %v\nWant:\t %v", x[:30], f2Expected[:30])
		t.Errorf("file %q and file %q are not equal", f2, "testdata/testFolder/file2")
	}

	x, err = os.ReadFile(f3)
	if err != nil {
		t.Errorf("ioutil.ReadFile(%q) = %q, not nil", f3, err)
	}
	if !bytes.Equal(x, f3Expected) {
		t.Logf("\nGot:\t %v\nWant:\t %v", x[:30], f3Expected[:30])
		t.Errorf("file %q and file %q are not equal", f3, "testdata/testFolder/file3")
	}

	x, err = os.ReadFile(f4)
	if err != nil {
		t.Errorf("ioutil.ReadFile(%q) = %q, not nil", f4, err)
	}
	if !bytes.Equal(x, f4Expected) {
		t.Logf("\nGot:\t %v\nWant:\t %v", x[:30], f4Expected[:30])
		t.Errorf("file %q and file %q are not equal\n", f4, "testdata/testFolder/file4")
	}
}

func TestFromZipNoValidFile(t *testing.T) {
	f, err := os.CreateTemp("", "testfile-")
	if err != nil {
		t.Errorf(`ioutil.TempFile("", "testfile-") = %q, not nil`, err)
	}
	defer f.Close()
	if err := FromZip(f.Name(), "someDir"); err == nil {
		t.Errorf(`FromZip(f.Name(), "someDir") = %q, not nil`, err)
	}
}

func TestAppendZip(t *testing.T) {
	_, err := os.Create("appendTest.zip")
	if err != nil {
		t.Errorf(`os.Create("appendTest.zip") = %q, not nil`, err)
	}
	defer os.Remove("appendTest.zip")

	if err := AppendZip("testdata/testFolder", "appendTest.zip", "Test append zip"); err != nil {
		t.Errorf(`AppendZip("testdata/testFolder", "appendTest.zip", "Test append zip") = %q, not nil`, err)
	}
}

func TestAppendZipNoDir1(t *testing.T) {
	if err := AppendZip("doesNotExist", "alsoNotExist", "Whythough"); err == nil {
		t.Errorf(`AppendZip("doesNotExist", "alsoNotExist", "Whythough") = %q, not nil`, err)
	}
}

func TestAppendZipNoDir2(t *testing.T) {
	f, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Errorf(`ioutil.TempFile("", "testfile") = _, %q, not _, nil`, err)
	}
	defer f.Close()
	if err := AppendZip(f.Name(), f.Name(), "no comment"); err == nil {
		t.Errorf(`AppendZip(f.Name(), f.Name(), "no comment") = %q, not nil`, err)
	}
}

func TestComment(t *testing.T) {
	comment, err := Comment("test.zip")
	if err != nil {
		t.Errorf(`Comment("test.zip") = %q, not nil`, err)
	}
	fmt.Println(comment)
}

func TestToZip(t *testing.T) {
	if err := ToZip(".", "testfile.zip", "test comment"); err != nil {
		t.Errorf(`ToZip(".", "testfile.zip", "test comment") = %q, not nil`, err)
	}
	defer os.Remove("testfile.zip")
}

func TestToZipInvalidDir(t *testing.T) {
	f, err := os.CreateTemp("", "testfile-")
	if err != nil {
		t.Errorf(`ioutil.TempFile("", "testfile-") = %q, not nil`, err)
	}
	defer f.Close()
	if err := ToZip(f.Name(), "invalid", "no need"); err == nil {
		t.Errorf(`ToZip(f.Name(), "invalid", "no need") = %q, not nil`, err)
	}
}
