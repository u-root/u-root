// Copyright 2018-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uzip

import (
	"bytes"
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

	if _, err := os.Stat(f1); err != nil {
		t.Errorf(`os.Stat(%q) =_, %v, want nil`, f1, err)
	}
	if _, err = os.Stat(f2); err != nil {
		t.Errorf(`os.Stat(%q) = _, %v, want nil`, f2, err)
	}
	if _, err := os.Stat(f3); err != nil {
		t.Errorf(`os.Stat(%q) = _, %v, want nil`, f3, err)
	}
	if _, err := os.Stat(f4); err != nil {
		t.Errorf(`os.Stat(%q) = _, %v, want nil`, f4, err)
	}

	if x, err := os.ReadFile(f1); err != nil || !bytes.Equal(x, f1Expected) {
		t.Errorf(`os.ReadFile(%q) = %v, %v, want %v, nil`, f1, x, err, f1Expected)
	}

	if x, err := os.ReadFile(f2); err != nil || !bytes.Equal(x, f2Expected) {
		t.Errorf(`os.ReadFile(%q) = %v, %v, want %v, nil`, f2, x, err, f2Expected)
	}

	if x, err := os.ReadFile(f3); err != nil || !bytes.Equal(x, f3Expected) {
		t.Errorf(`os.ReadFile(%q) = %v, %v, want %v, nil`, f3, x, err, f3Expected)
	}

	if x, err := os.ReadFile(f4); err != nil || !bytes.Equal(x, f4Expected) {
		t.Errorf(`os.ReadFile(%q) = %v, %v, want %v, nil`, f4, x, err, f4Expected)
	}
}
