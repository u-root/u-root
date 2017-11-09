// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"sort"

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

	// Reproducible builds: Files should be added to the archive in the
	// same order.
	r := opts.ArchiveFiles.Records
	f := opts.ArchiveFiles.Files
	paths := make([]string, 0, len(r)+len(f))
	for path := range r {
		paths = append(paths, path)
	}
	for path := range f {
		paths = append(paths, path)
	}
	sort.Sort(sort.StringSlice(paths))

	for _, path := range paths {
		if record, ok := r[path]; ok {
			if err := init.WriteRecord(record); err != nil {
				return err
			}
		}
		if src, ok := f[path]; ok {
			if err := init.WriteFile(src, path); err != nil {
				return err
			}
		}
	}

	return init.WriteTrailer()
}
