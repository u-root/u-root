// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkinitramfs creates a u-root initramfs given the list of files on the
// command line.
package mkinitramfs

import (
	"flag"
	"log"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uroot"
)

var (
	outputFile = flag.String("o", "initramfs.cpio", "Initramfs output file")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("must specify at least one file to include in initramfs")
	}
	archiver := uroot.CPIOArchiver{
		RecordFormat: cpio.Newc,
	}

	// Open the target initramfs file.
	w, err := archiver.OpenWriter(*outputFile, "", "")
	if err != nil {
		log.Fatalf("failed to open cpio archive %q: %v", *outputFile, err)
	}

	files := uroot.NewArchiveFiles()
	archive := uroot.ArchiveOpts{
		ArchiveFiles:   files,
		OutputFile:     w,
		DefaultRecords: uroot.DefaultRamfs,
	}
	if err := uroot.ParseExtraFiles(archive.ArchiveFiles, flag.Args(), false); err != nil {
		log.Fatalf("failed to parse file names %v: %v", flag.Args(), err)
	}

	if err := archive.Write(); err != nil {
		log.Fatalf("failed to write archive %q: %v", *outputFile, err)
	}
}
