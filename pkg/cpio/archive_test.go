// Copyright 2013-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"syscall"
	"testing"
)

func FuzzWriteReadInMemArchive(f *testing.F) {
	var fileCount uint64 = 4
	var content = []byte("Content")
	var name = "fileName"
	var ino, mode, uid, gid, nlink, mtime, major, minor, rmajor, rminor uint64 = 1, S_IFREG | 2, 3, 4, 5, 6, 7, 8, 9, 10
	f.Add(fileCount, content, name, ino, mode, uid, gid, nlink, mtime, major, minor, rmajor, rminor)
	f.Fuzz(func(t *testing.T, fileCount uint64, content []byte, name string, ino uint64, mode uint64, uid uint64, gid uint64, nlink uint64, mtime uint64, major uint64, minor uint64, rmajor uint64, rminor uint64) {
		if len(name) > 64 || len(content) > 200 || fileCount > 8 {
			return
		}
		recs := []Record{}
		var i uint64
		for i = 0; i < fileCount; i++ {

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
