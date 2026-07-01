// Copyright 2013-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package cpio

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
)

func TestCreateFileInRoot(t *testing.T) {
	tmp := t.TempDir()
	fileName := "file"
	content := "content"
	r := StaticFile(fileName, content, 0o644)
	err := CreateFileInRoot(r, tmp, false)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	b, err := os.ReadFile(filepath.Join(tmp, "file"))
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	if string(b) != content {
		t.Errorf("expected %q got %q", content, string(b))
	}
}

// modeSpyReaderAt records the on-disk mode of statPath the first time the
// record content is read, i.e. while CreateFileInRoot is still copying into
// the freshly created file and before it applies the final mode.
type modeSpyReaderAt struct {
	data     []byte
	statPath string
	once     sync.Once
	mode     os.FileMode
}

func (s *modeSpyReaderAt) ReadAt(p []byte, off int64) (int, error) {
	s.once.Do(func() {
		if fi, err := os.Stat(s.statPath); err == nil {
			s.mode = fi.Mode()
		}
	})
	if off >= int64(len(s.data)) {
		return 0, io.EOF
	}
	n := copy(p, s.data[off:])
	if off+int64(n) >= int64(len(s.data)) {
		return n, io.EOF
	}
	return n, nil
}

func TestCreateFileInRootSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	// First record drops a symlink pointing outside root. Creating it is fine,
	// it is only a link.
	if err := CreateFileInRoot(Symlink("escape", outside), root, false); err != nil {
		t.Fatalf("creating symlink: %v", err)
	}
	// The second record tries to write a file through that symlink. The write
	// must be refused, not followed out of the root.
	if err := CreateFileInRoot(StaticFile("escape/pwned", "owned", 0o644), root, false); err == nil {
		t.Errorf("writing through escaping symlink succeeded, want error")
	}

	if _, err := os.Stat(filepath.Join(outside, "pwned")); err == nil {
		t.Fatalf("file written outside root at %s", filepath.Join(outside, "pwned"))
	}
}

func TestCreateFileInRootSymlinkWithinRoot(t *testing.T) {
	root := t.TempDir()

	// A symlink that stays inside root must keep working: real/foo is reachable
	// through the link/foo path.
	recs := []Record{
		Directory("real", 0o755),
		Symlink("link", "real"),
		StaticFile("link/foo", "content", 0o644),
	}
	for _, r := range recs {
		if err := CreateFileInRoot(r, root, false); err != nil {
			t.Fatalf("CreateFileInRoot(%q): %v", r.Name, err)
		}
	}

	b, err := os.ReadFile(filepath.Join(root, "real", "foo"))
	if err != nil {
		t.Fatalf("reading real/foo: %v", err)
	}
	if string(b) != "content" {
		t.Errorf("real/foo = %q, want %q", b, "content")
	}
}

func TestCreateFileInRootMode(t *testing.T) {
	// Pin umask so the observed create mode is deterministic.
	old := syscall.Umask(0)
	defer syscall.Umask(old)

	tmp := t.TempDir()
	name := "secret"
	spy := &modeSpyReaderAt{
		data:     []byte("topsecret"),
		statPath: filepath.Join(tmp, name),
	}
	rec := Record{
		ReaderAt: spy,
		Info: Info{
			Name:     name,
			Mode:     S_IFREG | 0o600,
			FileSize: uint64(len(spy.data)),
		},
	}

	if err := CreateFileInRoot(rec, tmp, false); err != nil {
		t.Fatalf("CreateFileInRoot: %v", err)
	}

	if spy.mode&0o077 != 0 {
		t.Errorf("file group/world accessible while being written: mode %#o", spy.mode.Perm())
	}
	if got := spy.mode.Perm(); got != 0o600 {
		t.Errorf("create mode = %#o, want 0o600", got)
	}
}
