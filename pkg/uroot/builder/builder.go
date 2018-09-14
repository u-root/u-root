// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

var (
	BusyBox = BBBuilder{}
	Source  = SourceBuilder{}
	Binary  = BinaryBuilder{}
)

var Builders = map[string]Builder{
	"bb":     BusyBox,
	"source": Source,
	"binary": Binary,
}

// Opts are options passed to the Builder.Build function.
type Opts struct {
	// Env is the Go compiler environment.
	Env golang.Environ

	// Packages are the Go packages to compile.
	//
	// Only an explicit list of Go import paths is accepted.
	//
	// E.g. cmd/go or github.com/u-root/u-root/cmds/ls.
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
	Build(*initramfs.Files, Opts) error

	// DefaultBinaryDir is the initramfs' default directory for binaries
	// built using this builder.
	DefaultBinaryDir() string
}

// GetBuilder returns the Build function for the named build mode.
func GetBuilder(name string) (Builder, error) {
	build, ok := Builders[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find builder %q", name)
	}
	return build, nil
}
