// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package builder has methods for building many Go commands into an initramfs
// archive.
package builder

import (
	"errors"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage/initramfs"
	"github.com/u-root/uio/llog"
)

var (
	// Busybox is a shared GBBBuilder instance.
	Busybox = &GBBBuilder{}

	// Binary is a shared BinaryBuilder instance.
	Binary = &BinaryBuilder{}
)

// Possible build errors.
var (
	ErrBusyboxFailed  = errors.New("gobusybox build failed")
	ErrBinaryFailed   = errors.New("binary build failed")
	ErrEnvMissing     = errors.New("must specify Go build environment")
	ErrTempDirMissing = errors.New("must supply temporary directory for build")
)

// Opts are options passed to the Builder.Build function.
type Opts struct {
	// Env is the Go compiler environment.
	Env *golang.Environ

	// Build options for building go binaries. Ultimate this holds all the
	// args that end up being passed to `go build`.
	BuildOpts *golang.BuildOpts

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
	Build(*llog.Logger, *initramfs.Files, Opts) error

	// DefaultBinaryDir is the initramfs' default directory for binaries
	// built using this builder.
	DefaultBinaryDir() string
}
