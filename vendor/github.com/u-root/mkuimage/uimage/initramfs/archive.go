// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package initramfs can write archives out to CPIO or directories.
package initramfs

import (
	"errors"
	"io"

	"github.com/u-root/mkuimage/cpio"
)

// Possible errors.
var (
	ErrNoPath = errors.New("invalid argument: must specify path")
)

// ReadOpener opens a cpio.RecordReader.
type ReadOpener interface {
	OpenReader() (cpio.RecordReader, error)
}

// WriteOpener opens a Writer.
type WriteOpener interface {
	OpenWriter() (Writer, error)
}

// Writer is an initramfs archive that files can be written to.
type Writer interface {
	cpio.RecordWriter

	// Finish finishes the archive.
	Finish() error
}

// Opts are options for building an initramfs archive.
type Opts struct {
	// Files are the files to be included.
	//
	// Files here generally have priority over files in DefaultRecords or
	// BaseArchive.
	*Files

	// OutputFile is the file to write to.
	OutputFile WriteOpener

	// BaseArchive is an existing archive to add files to.
	//
	// BaseArchive may be nil.
	BaseArchive ReadOpener

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
		base, err := opts.BaseArchive.OpenReader()
		if err != nil {
			return err
		}
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
			f, err := base.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			// TODO: ignore only the error where it already exists
			// in archive.
			_ = opts.Files.AddRecord(transform(f))
		}
	}

	out, err := opts.OutputFile.OpenWriter()
	if err != nil {
		return err
	}
	if err := opts.Files.WriteTo(out); err != nil {
		return err
	}
	return out.Finish()
}
