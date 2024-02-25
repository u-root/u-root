// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"log/slog"
	"path"
	"path/filepath"

	"github.com/u-root/gobusybox/src/pkg/bb"
	"github.com/u-root/mkuimage/cpio"
	"github.com/u-root/mkuimage/uimage/initramfs"
	"github.com/u-root/uio/llog"
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
func (b *GBBBuilder) Build(l *llog.Logger, af *initramfs.Files, opts Opts) error {
	// Build the busybox binary.
	if len(opts.TempDir) == 0 {
		return ErrTempDirMissing
	}
	if opts.Env == nil {
		return ErrEnvMissing
	}
	bbPath := filepath.Join(opts.TempDir, "bb")

	binaryDir := opts.BinaryDir
	if binaryDir == "" {
		binaryDir = b.DefaultBinaryDir()
	}

	bopts := &bb.Opts{
		Env:          opts.Env,
		GenSrcDir:    opts.TempDir,
		CommandPaths: opts.Packages,
		BinaryPath:   bbPath,
		GoBuildOpts:  opts.BuildOpts,
	}
	if err := bb.BuildBusybox(l.AtLevel(slog.LevelInfo), bopts); err != nil {
		return fmt.Errorf("%w: %w", ErrBusyboxFailed, err)
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
			if err := af.AddRecord(cpio.StaticFile(filepath.Join(binaryDir, b), "#!/bbin/bb #!"+b+"\n", 0o755)); err != nil {
				return err
			}
		} else if err := af.AddRecord(cpio.Symlink(filepath.Join(binaryDir, path.Base(pkg)), "bb")); err != nil {
			return err
		}
	}
	return nil
}
