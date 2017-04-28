// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"os"
)

// Config contains various pieces of configuration created by main and passed
// around.
type Config struct {
	ArchiveFormat string
	BuildFormats  []string
	OutputPath    string
	Packages      []string
	Run           bool
	Verbose       bool
}

// Intermediate files are stored in memory so you do not need sudo to create
// cpio files. Some fields such as atime and mtime are not included here
// because they are useless for our purposes.
type file struct {
	path string
	data []byte
	mode os.FileMode
	uid  uint32
	gid  uint32
	rdev uint64
}

// Generates files for inclusion into the archive.
type buildGenerator interface {
	generate(Config) ([]file, error)
}

// Create an archive given a slice of files.
type archiveGenerator interface {
	generate(Config, []file) error
	run(Config) error
}
