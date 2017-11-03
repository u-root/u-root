// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/u-root/u-root/pkg/golang"
)

type SourceBuilder struct{}

// buildToolchain builds the needed Go toolchain binaries: go, compile, link,
// asm.
func (s SourceBuilder) buildToolchain(opts BuilderOpts, out *ArchiveFiles) error {
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

	// Add Go toolchain to archive.
	return out.AddFile(opts.TempDir, "")
}

func goListPkg(opts BuilderOpts, pkg string, out *ArchiveFiles) *golang.ListPackage {
	p, err := opts.Env.ListDeps(pkg)
	if err != nil {
		log.Printf("Can't list Go dependencies for %v; ignoring.", pkg)
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

func (s SourceBuilder) Build(opts BuilderOpts) (ArchiveFiles, error) {
	af := NewArchiveFiles()

	if err := af.AddFile(filepath.Join(opts.Env.GOROOT, "pkg/include"), "go/pkg/include"); err != nil {
		return ArchiveFiles{}, err
	}

	log.Printf("Collecting package files and dependencies...")
	deps := make(map[string]struct{})
	for _, pkg := range opts.Packages {
		// Add high-level packages' src files to archive.
		p := goListPkg(opts, pkg, &af)
		if p == nil {
			continue
		}
		for _, d := range p.Deps {
			deps[d] = struct{}{}
		}
	}
	// Add src files of dependencies to archive.
	for dep := range deps {
		goListPkg(opts, dep, &af)
	}

	// Add Go toolchain.
	log.Printf("Building go toolchain...")
	if err := s.buildToolchain(opts, &af); err != nil {
		return ArchiveFiles{}, err
	}
	return af, nil
}
