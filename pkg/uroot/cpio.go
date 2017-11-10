// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
	"github.com/u-root/u-root/pkg/ramfs"
)

// CPIOArchiver is an implementation of Archiver for the cpio format.
type CPIOArchiver struct {
	// Format is the name of the cpio format to use.
	Format string
}

// DefaultExtension implements Archiver.DefaultExtension.
func (ca CPIOArchiver) DefaultExtension() string {
	return "cpio"
}

// Archive implements Archiver.Archive.
func (ca CPIOArchiver) Archive(opts ArchiveOpts) error {
	archiver, err := cpio.Format(ca.Format)
	if err != nil {
		return err
	}

	init, err := ramfs.NewInitramfs(archiver.Writer(opts.OutputFile))
	if err != nil {
		return err
	}

	if opts.BaseArchive != nil {
		transform := cpio.MakeReproducible

		// Rename init to inito if there is another init.
		if !opts.UseExistingInit && opts.Contains("init") {
			transform = func(r cpio.Record) cpio.Record {
				if r.Name == "init" {
					r.Name = "inito"
				}
				return cpio.MakeReproducible(r)
			}
		}

		if err := init.Concat(archiver.Reader(opts.BaseArchive), transform); err != nil {
			return err
		}
	}

	// Reproducible builds: Files should be added to the archive in the
	// same order.
	for _, path := range opts.ArchiveFiles.SortedKeys() {
		if record, ok := opts.ArchiveFiles.Records[path]; ok {
			if err := init.WriteRecord(record); err != nil {
				return err
			}
		}
		if src, ok := opts.ArchiveFiles.Files[path]; ok {
			if err := init.WriteFile(src, path); err != nil {
				return err
			}
		}
	}
	return init.WriteTrailer()
}
