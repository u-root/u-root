// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uroot/logger"
)

var (
	CPIO = CPIOArchiver{
		RecordFormat: cpio.Newc,
	}

	Dir = DirArchiver{}

	// Archivers are the supported initramfs archivers at the moment.
	//
	// - cpio: writes the initramfs to a cpio.
	// - dir:  writes the initramfs relative to a specified directory.
	Archivers = map[string]Archiver{
		"cpio": CPIO,
		"dir":  Dir,
	}
)

// Archiver is an archive format that builds an archive using a given set of
// files.
type Archiver interface {
	// OpenWriter opens an archive writer at `path`.
	//
	// If `path` is unspecified, implementations may choose an arbitrary
	// default location, potentially based on `goos` and `goarch`.
	OpenWriter(l logger.Logger, path, goos, goarch string) (Writer, error)

	// Reader returns a Reader that allows reading files from a file.
	Reader(file io.ReaderAt) Reader
}

// GetArchiver finds a registered initramfs archiver by name.
//
// Good to use with command-line arguments.
func GetArchiver(name string) (Archiver, error) {
	archiver, ok := Archivers[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find archival format %q", name)
	}
	return archiver, nil
}

// Writer is an initramfs archive that files can be written to.
type Writer interface {
	cpio.RecordWriter

	// Finish finishes the archive.
	Finish() error
}

// Reader is an object that files can be read from.
type Reader cpio.RecordReader

// Opts are options for building an initramfs archive.
type Opts struct {
	// Files are the files to be included.
	//
	// Files here generally have priority over files in DefaultRecords or
	// BaseArchive.
	*Files

	// OutputFile is the file to write to.
	OutputFile Writer

	// BaseArchive is an existing archive to add files to.
	//
	// BaseArchive may be nil.
	BaseArchive Reader

	// UseExistingInit determines whether the init from BaseArchive is used
	// or not, if BaseArchive is specified.
	//
	// If this is false, the "init" file in BaseArchive will be renamed
	// "inito" (for init-original) in the output archive.
	UseExistingInit bool
}

// Write uses the given options to determine which files to write to the output
// initramfs.
func Write(opts *Opts) error {
	// Write base archive.
	if opts.BaseArchive != nil {
		transform := cpio.MakeReproducible

		// Rename init to inito if user doesn't want the existing init.
		if !opts.UseExistingInit && opts.Contains("init") {
			transform = func(r cpio.Record) cpio.Record {
				if r.Name == "init" {
					r.Name = "inito"
				}
				return cpio.MakeReproducible(r)
			}
		}
		// If user wants the base archive init, but specified another
		// init, make the other one inito.
		if opts.UseExistingInit && opts.Contains("init") {
			opts.Rename("init", "inito")
		}

		for {
			f, err := opts.BaseArchive.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			// TODO: ignore only the error where it already exists
			// in archive.
			opts.Files.AddRecord(transform(f))
		}
	}

	if err := opts.Files.WriteTo(opts.OutputFile); err != nil {
		return err
	}
	return opts.OutputFile.Finish()
}
