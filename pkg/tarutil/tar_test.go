// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tarutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func extractAndCompare(t *testing.T, tarFile string, files []struct{ name, body string }) {
	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Extract tar to tmpDir.
	f, err := os.Open(tarFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := ExtractDir(f, tmpDir, &Opts{Filters: []Filter{SafeFilter}}); err != nil {
		t.Fatal(err)
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

func TestExtractDir(t *testing.T) {
	var files = []struct {
		name, body string
	}{
		{"a.txt", "hello\n"},
		{"dir/b.txt", "world\n"},
	}
	extractAndCompare(t, "test.tar", files)
}

func TestCreateTarSingleFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the tar file.
	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := CreateTar(f, []string{"test0"}, nil); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("tar", "-tf", filename).CombinedOutput()
	if err != nil {
		t.Fatalf("system tar could not parse the file: %v", err)
	}
	expected := `test0
test0/a.txt
test0/dir
test0/dir/b.txt
`
	if string(out) != expected {
		t.Fatalf("got %q, want %q", string(out), expected)
	}
}

func TestCreateTarMultFiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the tar file.
	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	files := []string{"test0", "test1", "test2.txt"}
	if err := CreateTar(f, files, nil); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("tar", "-tf", filename).CombinedOutput()
	if err != nil {
		t.Fatalf("system tar could not parse the file: %v", err)
	}
	expected := `test0
test0/a.txt
test0/dir
test0/dir/b.txt
test1
test1/a1.txt
test2.txt
`
	if string(out) != expected {
		t.Fatalf("got %q, want %q", string(out), expected)
	}
}

// TestCreateTarProcfsFile exercises the special case where stat on the file
// reports a size of 0, but the file actually has contents. For example, most
// of the files in /proc and /sys.
func TestCreateTarProcfsFile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("/proc/version is only on linux, but GOOS=%s", runtime.GOOS)
	}

	tmpDir, err := ioutil.TempDir("", "tartest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// /proc/version won't change during the test. The size according to
	// stat should also be 0.
	procfsFile := "/proc/version"
	contents, err := ioutil.ReadFile(procfsFile)
	if err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(procfsFile); err != nil {
		t.Fatal(err)
	} else if fi.Size() != 0 {
		t.Fatalf("Expected the size of %q to be 0, got %d", procfsFile, fi.Size())
	}

	// Create the tar file containing /proc/version.
	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := CreateTar(f, []string{"/proc", procfsFile}, &Opts{NoRecursion: true}); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	var expected = []struct {
		name, body string
	}{
		{procfsFile, string(contents)},
	}
	extractAndCompare(t, filename, expected)
}

func TestListArchive(t *testing.T) {
	f, err := os.Open("test.tar")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := ListArchive(f); err != nil {
		t.Fatal(err)
	}
}
