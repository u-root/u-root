// Copyright 2013-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"strings"
	"syscall"
	"testing"
)

func TestArchiveMethods(t *testing.T) {
	r1 := StaticFile("/bin/r1", "content1", 0o644)
	r2 := StaticFile("/bin/r2", "content2", 0o755)
	records := []Record{r1, r2}

	ar1 := ArchiveFromRecords(records)
	r := ar1.Reader()
	ar2, err := ArchiveFromReader(r)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	if ar2.Empty() {
		t.Fatalf("archive is empty, expected to have two records")
	}

	for _, r := range records {
		if !ar2.Contains(r) {
			t.Errorf("records %v missing", r)
		}
	}

	_, ok := ar2.Get("/bin/r1")
	if !ok {
		t.Errorf("record %v not found", r1)
	}

	_, ok = ar2.Get("/bin/r2")
	if !ok {
		t.Errorf("record %v not found", r2)
	}

	s := ar2.String()
	if !strings.Contains(s, "bin/r1") {
		t.Errorf("missing 'bin/r1' from ar2.String()")
	}
	if !strings.Contains(s, "bin/r2") {
		t.Errorf("missing 'bin/r2' from ar2.String()")
	}
}

func FuzzWriteReadInMemArchive(f *testing.F) {
	var fileCount uint64 = 4
	content := []byte("Content")
	name := "fileName"
	var ino, mode, uid, gid, nlink, mtime, major, minor, rmajor, rminor uint64 = 1, S_IFREG | 2, 3, 4, 5, 6, 7, 8, 9, 10
	f.Add(fileCount, content, name, ino, mode, uid, gid, nlink, mtime, major, minor, rmajor, rminor)
	f.Fuzz(func(t *testing.T, fileCount uint64, content []byte, name string, ino uint64, mode uint64, uid uint64, gid uint64, nlink uint64, mtime uint64, major uint64, minor uint64, rmajor uint64, rminor uint64) {
		if len(name) > 64 || len(content) > 200 || fileCount > 8 {
			return
		}
		recs := []Record{}
		var i uint64
		for i = range fileCount {
			recs = append(recs, StaticRecord(content, Info{
				Ino:      ino | i,
				Mode:     syscall.S_IFREG | mode | i,
				UID:      uid | i,
				GID:      gid | i,
				NLink:    nlink | i,
				MTime:    mtime | i,
				FileSize: uint64(len(content)),
				Major:    major | i,
				Minor:    minor | i,
				Rmajor:   rmajor | i,
				Rminor:   rminor | i,
				Name:     Normalize(name) + fmt.Sprintf("%d", i),
			}))
		}

		arch := ArchiveFromRecords(recs)
		archReader := arch.Reader()

		for _, rec := range recs {
			readRec, err := archReader.ReadRecord()
			if err != nil {
				t.Fatalf("failed to read record from archive")
			}

			if !Equal(rec, readRec) {
				t.Fatalf("records not equal: %v %v", rec, readRec)
			}

			if !arch.Contains(rec) {
				t.Fatalf("record not in archive %v %#v", rec, arch)
			}
		}
	})
}
