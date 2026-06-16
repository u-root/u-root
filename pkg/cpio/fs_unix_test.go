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
