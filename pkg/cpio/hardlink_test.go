// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package cpio

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// Create a file, with one hard link, and verify that we
// create the records and then unpack it correctly.
func TestHardLink(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestHardLink")
	if err != nil {
		t.Fatalf("Tempdir: got %v, want nil", err)
	}
	t.Logf("Testing in %v", dir)
	defer os.RemoveAll(dir)
	f1 := filepath.Join(dir, "a")
	if err := ioutil.WriteFile(f1, nil, 0666); err != nil {
		t.Fatalf("create %s: got %v, want nil", f1, err)
	}
	f2 := filepath.Join(dir, "b")
	if err := unix.Link(f1, f2); err != nil {
		t.Fatalf("link %s to %s: got %v, want nil", f1, f2, err)
	}
	buf := &bytes.Buffer{}
	w := Newc.Writer(buf)
	cr := NewRecorder()
	names := []string{f1, f2}
	var rec Record
	for _, n := range names {
		rec, err = cr.GetRecord(n)

		if err != nil {
			t.Fatalf("Getting record of %q: got %v, want nil", n, err)
		}
		if err := w.WriteRecord(rec); err != nil {
			t.Fatalf("Writing record %q: got %v, want nil", n, err)
		}
	}

	// Programatically create a Hardlink record.
	f3 := filepath.Join(dir, "c")
	hl := Hardlink(f3, rec.Info)

	if err := w.WriteRecord(hl); err != nil {
		t.Fatalf("Writing record %s: got %v, want nil", rec, err)
	}
	for _, n := range names {
		if err := os.Remove(n); err != nil {
			t.Fatalf("remove: got %v, want nil", err)
		}
	}
	if err := WriteTrailer(w); err != nil {
		t.Fatalf("Writing trailer record: got %v, want nil", err)
	}

	rr := Newc.Reader(bytes.NewReader(buf.Bytes()))

	ww := NewUnixFiler(func(f *UnixFiler) {
		f.Root = "/"
	})
	var nread int
	for {
		rec, err := rr.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Reading records: got %v, want nil", err)
		}
		t.Logf("Creating %s\n", rec)
		// "/" seems scary but just recreates the tempdir
		// we had before.
		if err := ww.Create(rec); err != nil {
			log.Printf("Creating %q failed: %v", rec.Name, err)
		}
		nread++
	}

	if nread != 3 {
		t.Errorf("reading records: got %d, want 3", nread)
	}
	fi, err := os.Stat(names[0])
	if err != nil {
		t.Fatalf("Stat %q: got %v, want nil", names[0], err)
	}
	ino := fi.Sys().(*syscall.Stat_t).Ino

	names = append(names, f3)
	for _, n := range names {
		fi, err := os.Stat(n)
		if err != nil {
			t.Fatalf("Stat %q: got %v, want nil", n, err)
		}
		t.Logf("Stat %q gets %v", n, fi)
		if ino != fi.Sys().(*syscall.Stat_t).Ino {
			t.Errorf("%q: inode got %v, want %v", n, fi.Sys().(*syscall.Stat_t).Ino, ino)
		}
	}
}
