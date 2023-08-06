// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"golang.org/x/tools/go/packages"
)

func lookupPkgs(env *gbbgolang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
		Env:  append(os.Environ(), env.Env()...),
		Dir:  dir,
	}
	return packages.Load(cfg, patterns...)
}

func dirFor(env *gbbgolang.Environ, pkg string) (string, error) {
	pkgs, err := lookupPkgs(env, "", pkg)
	if err != nil {
		return "", fmt.Errorf("failed to look up package %q: %v", pkg, err)
	}

	// One directory = one package in standard Go, so
	// finding the first file's parent directory should
	// find us the package directory.
	var dir string
	for _, p := range pkgs {
		if len(p.GoFiles) > 0 {
			dir = filepath.Dir(p.GoFiles[0])
		}
	}
	if dir == "" {
		return "", fmt.Errorf("could not find package directory for %q", pkg)
	}
	return dir, nil
}

// BinaryBuilder builds each Go command as a separate binary.
//
// BinaryBuilder is an implementation of Builder.
type BinaryBuilder struct{}

// DefaultBinaryDir implements Builder.DefaultBinaryDir.
//
// "bin" is the default initramfs binary directory for these binaries.
func (BinaryBuilder) DefaultBinaryDir() string {
	return "bin"
}

// Build implements Builder.Build.
func (BinaryBuilder) Build(l ulog.Logger, af *initramfs.Files, opts Opts) error {
	if opts.Env == nil {
		return fmt.Errorf("must specify Go build environment")
	}
	result := make(chan error, len(opts.Packages))

	var wg sync.WaitGroup
	for _, pkg := range opts.Packages {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			dir, err := dirFor(opts.Env, p)
			if err != nil {
				result <- err
				return
			}
			result <- opts.Env.BuildDir(
				dir,
				filepath.Join(opts.TempDir, opts.BinaryDir, filepath.Base(p)),
				opts.BuildOpts)
		}(pkg)
	}

	wg.Wait()
	close(result)

	for err := range result {
		if err != nil {
			return err
		}
	}

	// Add bin directory to archive.
	return af.AddFile(opts.TempDir, "")
}
