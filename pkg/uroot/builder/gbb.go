// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/u-root/gobusybox/src/pkg/bb"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// GBBBuilder is an implementation of Builder that compiles many Go commands
// into one busybox-style binary.
//
// GBBBuilder will also include symlinks for each command to the busybox binary.
//
// GBBBuilder does all this by rewriting the source files of the packages given
// to create one busybox-like binary containing all commands.
//
// The compiled binary uses argv[0] to decide which Go command to run.
//
// See bb/README.md for a detailed explanation of the implementation of busybox
// mode.
type GBBBuilder struct {
	// ShellBang means generate #! files instead of symlinks.
	// ShellBang are more portable and just as efficient.
	ShellBang bool
}

// DefaultBinaryDir implements Builder.DefaultBinaryDir.
//
// The default initramfs binary dir is bbin for busybox binaries.
func (GBBBuilder) DefaultBinaryDir() string {
	return "bbin"
}

// Build is an implementation of Builder.Build for a busybox-like initramfs.
func (b GBBBuilder) Build(af *initramfs.Files, opts Opts) error {
	// Build the busybox binary.
	if len(opts.TempDir) == 0 {
		return fmt.Errorf("opts.TempDir is empty")
	}
	bbPath := filepath.Join(opts.TempDir, "bb")

	// gobusybox has its own copy of the golang package, but Environ stayed
	// (mostly) the same.
	env := golang.Environ{
		Context:     opts.Env.Context,
		GO111MODULE: os.Getenv("GO111MODULE"),
	}

	if len(opts.BinaryDir) == 0 {
		return fmt.Errorf("must specify binary directory")
	}

	bopts := &bb.Opts{
		Env:          env,
		GenSrcDir:    opts.TempDir,
		CommandPaths: opts.Packages,
		BinaryPath:   bbPath,
		GoBuildOpts:  &golang.BuildOpts{},
	}
	bopts.GoBuildOpts.RegisterFlags(flag.CommandLine)

	if err := bb.BuildBusybox(bopts); err != nil {
		var errGopath *bb.ErrGopathBuild
		var errGomod *bb.ErrModuleBuild
		if errors.As(err, &errGopath) {
			return fmt.Errorf("preserving bb generated source directory at %s due to error. To reproduce build, `cd %s` and `GO111MODULE=off GOPATH=%s go build`", opts.TempDir, errGopath.CmdDir, errGopath.GOPATH)
		} else if errors.As(err, &errGomod) {
			return fmt.Errorf("preserving bb generated source directory at %s due to error. To debug build, `cd %s` and use `go build` to build, or `go mod [why|tidy|graph]` to debug dependencies, or `go list -m all` to list all dependency versions", opts.TempDir, errGomod.CmdDir)
		} else {
			return fmt.Errorf("preserving bb generated source directory at %s due to error", opts.TempDir)
		}
	}

	if err := af.AddFile(bbPath, "/bbin/bb"); err != nil {
		return err
	}

	// Add symlinks for included commands to initramfs.
	for _, pkg := range opts.Packages {
		if _, ok := skip[path.Base(pkg)]; ok {
			continue
		}

		// Add a symlink /bbin/{cmd} -> /bbin/bb to our initramfs.
		// Or add a #! file if b.ShellBang is set ...
		if b.ShellBang {
			b := path.Base(pkg)
			if err := af.AddRecord(cpio.StaticFile(filepath.Join(opts.BinaryDir, b), "#!/bbin/bb #!"+b+"\n", 0o755)); err != nil {
				return err
			}
		} else if err := af.AddRecord(cpio.Symlink(filepath.Join(opts.BinaryDir, path.Base(pkg)), "bb")); err != nil {
			return err
		}
	}
	return nil
}
