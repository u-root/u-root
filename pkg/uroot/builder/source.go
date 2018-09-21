// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// SourceBuilder includes full source for Go commands in the initramfs.
//
// SourceBuilder is an implementation of Builder.
//
// It also includes the Go toolchain in the initramfs, and a tool called
// installcommand that can compile the other commands using symlinks.
//
// E.g. if "ls" is an included command, "ls" will be a symlink to
// "installcommand" in the initramfs, which uses argv[0] to figure out which
// command to compile.
type SourceBuilder struct{}

// DefaultBinaryDir implements Builder.DefaultBinaryDir.
//
// The initramfs default binary dir is buildbin.
func (SourceBuilder) DefaultBinaryDir() string {
	return "buildbin"
}

// Build is an implementation of Builder.Build.
func (SourceBuilder) Build(af *initramfs.Files, opts Opts) error {
	// TODO: this is a failure to collect the correct dependencies.
	if err := af.AddFile(filepath.Join(opts.Env.GOROOT, "pkg/include"), "go/pkg/include"); err != nil {
		return err
	}

	var installcommand string
	log.Printf("Collecting package files and dependencies...")
	deps := make(map[string]struct{})
	for _, pkg := range opts.Packages {
		name := path.Base(pkg)
		if name == "installcommand" {
			installcommand = pkg
		}

		// Add high-level packages' src files to archive.
		p := goListPkg(opts, pkg, af)
		if p == nil {
			continue
		}
		for _, d := range p.Deps {
			deps[d] = struct{}{}
		}

		if name != "installcommand" {
			// Add a symlink to installcommand. This means source mode can
			// work with any init.
			if err := af.AddRecord(cpio.Symlink(
				path.Join(opts.BinaryDir, name),
				path.Join("/", opts.BinaryDir, "installcommand"))); err != nil {
				return err
			}
		}
	}
	if len(installcommand) == 0 {
		return fmt.Errorf("must include a version of installcommand in source mode")
	}

	// Add src files of dependencies to archive.
	for dep := range deps {
		goListPkg(opts, dep, af)
	}

	// Add Go toolchain.
	log.Printf("Building go toolchain...")
	if err := buildToolchain(opts); err != nil {
		return err
	}
	if err := opts.Env.Build(installcommand, filepath.Join(opts.TempDir, opts.BinaryDir, "installcommand"), golang.BuildOpts{}); err != nil {
		return err
	}

	// Add Go toolchain and installcommand to archive.
	return af.AddFile(opts.TempDir, "")
}

// buildToolchain builds the needed Go toolchain binaries: go, compile, link,
// asm.
func buildToolchain(opts Opts) error {
	goBin := filepath.Join(opts.TempDir, "go/bin/go")
	tcbo := golang.BuildOpts{
		ExtraArgs: []string{"-tags", "cmd_go_bootstrap"},
	}
	if err := opts.Env.Build("cmd/go", goBin, tcbo); err != nil {
		return err
	}

	toolDir := filepath.Join(opts.TempDir, fmt.Sprintf("go/pkg/tool/%v_%v", opts.Env.GOOS, opts.Env.GOARCH))
	for _, pkg := range []string{"compile", "link", "asm"} {
		c := filepath.Join(toolDir, pkg)
		if err := opts.Env.Build(fmt.Sprintf("cmd/%s", pkg), c, golang.BuildOpts{}); err != nil {
			return err
		}
	}
	return nil
}

func goListPkg(opts Opts, importPath string, out *initramfs.Files) *golang.ListPackage {
	p, err := opts.Env.Deps(importPath)
	if err != nil {
		log.Printf("Can't list Go dependencies for %v; ignoring.", importPath)
		return nil
	}

	// Add Go files in this package to archive.
	for _, file := range append(append(p.GoFiles, p.SFiles...), p.HFiles...) {
		relPath := filepath.Join("src", p.ImportPath, file)
		srcFile := filepath.Join(p.Root, relPath)
		if p.Goroot {
			out.AddFile(srcFile, filepath.Join("go", relPath))
		} else {
			out.AddFile(srcFile, relPath)
		}
	}
	return p
}
