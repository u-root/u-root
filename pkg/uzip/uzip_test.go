// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uzip

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromZip(t *testing.T) {
	tmpDir := t.TempDir()

	f := filepath.Join(tmpDir, "test.zip")
	if err := ToZip("testdata/testFolder", f, ""); err != nil {
		t.Fatal(err)
	}

	z, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, z)

	out := filepath.Join(tmpDir, "unziped")
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := FromZip(f, out); err != nil {
		t.Fatal(err)
	}

	f1 := filepath.Join(out, "file1")
	f2 := filepath.Join(out, "file2")
	f3 := filepath.Join(out, "subFolder", "file3")
	f4 := filepath.Join(out, "subFolder", "file4")

	f1Expected, err := os.ReadFile("testdata/testFolder/file1")
	if err != nil {
		t.Fatal(err)
	}
	f2Expected, err := os.ReadFile("testdata/testFolder/file2")
	if err != nil {
		t.Fatal(err)
	}
	f3Expected, err := os.ReadFile("testdata/testFolder/subFolder/file3")
	if err != nil {
		t.Fatal(err)
	}
	f4Expected, err := os.ReadFile("testdata/testFolder/subFolder/file4")
	if err != nil {
		t.Fatal(err)
	}

	require.FileExists(t, f1)
	require.FileExists(t, f2)
	require.FileExists(t, f3)
	require.FileExists(t, f4)

	var x []byte

	x, err = os.ReadFile(f1)
	require.NoError(t, err)
	require.Equal(t, x, f1Expected)

	x, err = os.ReadFile(f2)
	require.NoError(t, err)
	require.Equal(t, x, f2Expected)

	x, err = os.ReadFile(f3)
	require.NoError(t, err)
	require.Equal(t, x, f3Expected)

	x, err = os.ReadFile(f4)
	require.NoError(t, err)
	require.Equal(t, x, f4Expected)
}
