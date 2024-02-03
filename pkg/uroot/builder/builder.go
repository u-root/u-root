// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"github.com/u-root/uio/ulog"
)

var (
	// BusyBox is a shared GBBBuilder instance.
	BusyBox = GBBBuilder{}
	// Binary is a shared BinaryBuilder instance.
	Binary = BinaryBuilder{}
)

// Opts are options passed to the Builder.Build function.
type Opts struct {
	// Env is the Go compiler environment.
	Env *gbbgolang.Environ

	// Build options for building go binaries. Ultimate this holds all the
	// args that end up being passed to `go build`.
	BuildOpts *gbbgolang.BuildOpts

	// Packages are the Go packages to compile.
	//
	// Only an explicit list of absolute directory paths is accepted.
	Packages []string

	// TempDir is a temporary directory where the compilation mode compiled
	// binaries can be placed.
	//
	// TempDir should contain no files.
	TempDir string

	// BinaryDir is the initramfs directory for built binaries.
	//
	// BinaryDir must be specified.
	BinaryDir string
}

// Builder builds Go packages and adds the binaries to an initramfs.
//
// The resulting files need not be binaries per se, but exec'ing the resulting
// file should result in the Go program being executed.
type Builder interface {
	// Build uses the given options to build Go packages and adds its files
	// to be included in the initramfs to the given ArchiveFiles.
	Build(ulog.Logger, *initramfs.Files, Opts) error

	// DefaultBinaryDir is the initramfs' default directory for binaries
	// built using this builder.
	DefaultBinaryDir() string
}
