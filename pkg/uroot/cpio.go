// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
	"github.com/u-root/u-root/pkg/ramfs"
)

type CPIOArchiver struct{}

func (ca CPIOArchiver) DefaultExtension() string {
	return "cpio"
}

func (ca CPIOArchiver) Archive(opts ArchiveOpts) error {
	archiver, err := cpio.Format("newc")
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

	for _, record := range opts.ArchiveFiles.Records {
		if err := init.WriteRecord(record); err != nil {
			return err
		}
	}
	for dest, src := range opts.ArchiveFiles.Files {
		if err := init.WriteFile(src, dest); err != nil {
			return err
		}
	}
	return init.WriteTrailer()
}
