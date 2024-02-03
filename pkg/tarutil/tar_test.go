// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tarutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func extractAndCompare(t *testing.T, tarFile string, files []struct{ name, body string }) {
	tmpDir := t.TempDir()

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
		body, err := os.ReadFile(filepath.Join(tmpDir, f.name))
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
	files := []struct {
		name, body string
	}{
		{"a.txt", "hello\n"},
		{"dir/b.txt", "world\n"},
	}
	extractAndCompare(t, "testdata/test.tar", files)
}

func TestCreateTarSingleFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the tar file.
	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := CreateTar(f, []string{"testdata/test0"}, &Opts{
		NoRecursion: true,
		Filters:     []Filter{VerboseFilter, VerboseLogFilter},
	}); err != nil {
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
	expected := "testdata/test0\n"
	if string(out) != expected {
		t.Fatalf("got %q, want %q", string(out), expected)
	}
}

func TestCreateTarMultFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the tar file.
	filename := filepath.Join(tmpDir, "test.tar")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	files := []string{"testdata/test0", "testdata/test1", "testdata/test2.txt"}
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
	expected := `testdata/test0
testdata/test0/a.txt
testdata/test0/dir
testdata/test0/dir/b.txt
testdata/test1
testdata/test1/a1.txt
testdata/test2.txt
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

	tmpDir := t.TempDir()

	// /proc/version won't change during the test. The size according to
	// stat should also be 0.
	procfsFile := "/proc/version"
	contents, err := os.ReadFile(procfsFile)
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

	opts := &Opts{
		NoRecursion: true,
		// We want the path names in the tar archive to be relative to
		// root. We cannot store absolute paths in the tar file because
		// it will fail the zipslip check on extraction.
		ChangeDirectory: "/",
	}
	if err := CreateTar(f, []string{"proc", "proc/version"}, opts); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		name, body string
	}{
		{procfsFile, string(contents)},
	}
	extractAndCompare(t, filename, expected)
}

func TestListArchive(t *testing.T) {
	f, err := os.Open("testdata/test.tar")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := ListArchive(f); err != nil {
		t.Fatal(err)
	}
}
