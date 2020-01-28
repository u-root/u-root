// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uzip

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFromZip(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ziptest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	f := filepath.Join(tmpDir, "test.zip")
	err = ToZip("testdata/testFolder", f)
	if err != nil {
		t.Fatal(err)
	}

	z, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, z)

	out := filepath.Join(tmpDir, "unziped")
	err = os.MkdirAll(out, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = FromZip(f, out)
	if err != nil {
		t.Fatal(err)
	}

	f1 := filepath.Join(out, "file1")
	f2 := filepath.Join(out, "file2")
	f3 := filepath.Join(out, "subFolder", "file3")
	f4 := filepath.Join(out, "subFolder", "file4")

	f1Expected, err := ioutil.ReadFile("testdata/testFolder/file1")
	if err != nil {
		t.Fatal(err)
	}
	f2Expected, err := ioutil.ReadFile("testdata/testFolder/file2")
	if err != nil {
		t.Fatal(err)
	}
	f3Expected, err := ioutil.ReadFile("testdata/testFolder/subFolder/file3")
	if err != nil {
		t.Fatal(err)
	}
	f4Expected, err := ioutil.ReadFile("testdata/testFolder/subFolder/file4")
	if err != nil {
		t.Fatal(err)
	}

	require.FileExists(t, f1)
	require.FileExists(t, f2)
	require.FileExists(t, f3)
	require.FileExists(t, f4)

	var x []byte

	x, err = ioutil.ReadFile(f1)
	require.NoError(t, err)
	require.Equal(t, x, f1Expected)

	x, err = ioutil.ReadFile(f2)
	require.NoError(t, err)
	require.Equal(t, x, f2Expected)

	x, err = ioutil.ReadFile(f3)
	require.NoError(t, err)
	require.Equal(t, x, f3Expected)

	x, err = ioutil.ReadFile(f4)
	require.NoError(t, err)
	require.Equal(t, x, f4Expected)
}
