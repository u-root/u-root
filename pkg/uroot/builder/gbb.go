// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/u-root/gobusybox/src/pkg/bb"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"github.com/u-root/uio/ulog"
)

// Commands to skip building in bb mode.
var skip = map[string]struct{}{
	"bb": {},
}

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
func (b GBBBuilder) Build(l ulog.Logger, af *initramfs.Files, opts Opts) error {
	// Build the busybox binary.
	if len(opts.TempDir) == 0 {
		return fmt.Errorf("opts.TempDir is empty")
	}
	bbPath := filepath.Join(opts.TempDir, "bb")

	if len(opts.BinaryDir) == 0 {
		return fmt.Errorf("must specify binary directory")
	}
	if opts.Env == nil {
		return fmt.Errorf("must specify Go build environment")
	}

	bopts := &bb.Opts{
		Env:          opts.Env,
		GenSrcDir:    opts.TempDir,
		CommandPaths: opts.Packages,
		BinaryPath:   bbPath,
		GoBuildOpts:  opts.BuildOpts,
	}

	if err := bb.BuildBusybox(l, bopts); err != nil {
		// Print the actual error. This may contain a suggestion for
		// what to do, actually.
		l.Printf("Gobusybox error: %v", err)

		// Return some instructions for the user; this is printed last in the u-root tool.
		//
		// TODO: yeah, this isn't a good way to do error handling. The
		// error should be the thing that's returned, I just wanted
		// that to be printed first, and the instructions for what to
		// do about it to be last.
		var errGopath *bb.ErrGopathBuild
		var errGomod *bb.ErrModuleBuild
		if errors.As(err, &errGopath) {
			return fmt.Errorf("preserving bb generated source directory at %s due to error. To reproduce build, `cd %s` and `GO111MODULE=off GOPATH=%s go build`: %v", opts.TempDir, errGopath.CmdDir, errGopath.GOPATH, err)
		} else if errors.As(err, &errGomod) {
			return fmt.Errorf("preserving bb generated source directory at %s due to error. To debug build, `cd %s` and use `go build` to build, or `go mod [why|tidy|graph]` to debug dependencies, or `go list -m all` to list all dependency versions:\n%v", opts.TempDir, errGomod.CmdDir, err)
		} else {
			return fmt.Errorf("preserving bb generated source directory at %s due to error:\n%v", opts.TempDir, err)
		}
	}

	if err := af.AddFile(bbPath, "bbin/bb"); err != nil {
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
