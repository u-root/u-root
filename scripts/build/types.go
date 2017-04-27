// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"io"
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
	data io.Reader // Only for special files, may this be nil.
	mode os.FileMode
	uid  uint32
	gid  uint32
	rdev uint64
}

// From Linux header: /include/uapi/linux/kdev_t.h
const (
	minorBits = 8
	minorMask = (1 << minorBits) - 1
)

// dev returns the device number given the major and minor numbers.
func dev(major, minor uint64) uint64 {
	return major<<minorBits + minor
}

// major returns the device number's major number.
func major(dev uint64) uint64 {
	return dev >> minorBits
}

// minor returns the device number's minor number.
func minor(dev uint64) uint64 {
	return dev & minorMask
}

// Generate files for inclusion into the archive.
type builder interface {
	generate(Config) ([]file, error)
}

// Create an archive given a slice of files.
type archiver interface {
	generate(Config, []file) error
	run(Config) error
}
