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
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

var (
	outputFile = flag.String("o", "initramfs.cpio", "Initramfs output file")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("must specify at least one file to include in initramfs")
	}
	archiver := initramfs.CPIOArchiver{
		RecordFormat: cpio.Newc,
	}

	// Open the target initramfs file.
	w, err := archiver.OpenWriter(*outputFile, "", "")
	if err != nil {
		log.Fatalf("failed to open cpio archive %q: %v", *outputFile, err)
	}

	files := initramfs.NewFiles()
	archive := &initramfs.Opts{
		Files:       files,
		OutputFile:  w,
		BaseArchive: uroot.DefaultRamfs.Reader(),
	}
	if err := uroot.ParseExtraFiles(archive.Files, flag.Args(), false); err != nil {
		log.Fatalf("failed to parse file names %v: %v", flag.Args(), err)
	}

	if err := initramfs.Write(archive); err != nil {
		log.Fatalf("failed to write archive %q: %v", *outputFile, err)
	}
}
