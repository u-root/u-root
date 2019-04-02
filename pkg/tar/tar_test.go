// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tar

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExtractDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	f, err := os.Open("test.tar")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := ExtractDir(f, tmpDir); err != nil {
		t.Fatal(err)
	}

	var files = []struct {
		name, body string
	}{
		{"a.txt", "hello\n"},
		{"dir/b.txt", "world\n"},
	}
	for _, f := range files {
		body, err := ioutil.ReadFile(filepath.Join(tmpDir, f.name))
		if err != nil {
			t.Errorf("could not read %s: %v", f.name, err)
			continue
		}
		if string(body) != f.body {
			t.Errorf("for file %s, got %q, want %q",
				f.name, string(body), f.body)
		}
	}
}

func TestCreateDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := CreateDir(f, "test"); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("tar", "-tf", filename).CombinedOutput()
	if err != nil {
		t.Fatalf("system tar could not parse the file: %v", err)
	}
	expected := `test
test/a.txt
test/dir
test/dir/b.txt
`
	if string(out) != expected {
		t.Fatalf("got %q, want %q", string(out), expected)
	}
}
