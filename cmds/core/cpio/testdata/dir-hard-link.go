// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program was used to generate dir-hard-link.cpio, which is
// an archive containing two directories with the same inode (0).
// Depending on how the cpio package generates inodes in the future,
// it may not reproduce the file.
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

func main() {
	archiver, err := cpio.Format("newc")
	if err != nil {
		log.Fatal(err)
	}
	rw := archiver.Writer(os.Stdout)
	for _, rec := range []cpio.Record{
		cpio.Directory("directory1", 0o755),
		cpio.Directory("directory2", 0o755),
	} {
		rec.UID = uint64(os.Getuid())
		rec.GID = uint64(os.Getgid())
		if err := rw.WriteRecord(rec); err != nil {
			log.Fatal(err)
		}
	}
	if err := cpio.WriteTrailer(rw); err != nil {
		log.Fatal(err)
	}
}
