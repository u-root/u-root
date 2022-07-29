// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package findpkg finds packages from user-input strings that are either file
// paths or Go package paths.
package findpkg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
	"golang.org/x/tools/go/packages"

	"github.com/u-root/gobusybox/src/pkg/bb/bbinternal"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/ulog"
)

// modules returns a list of module directories => directories of packages
// inside that module as well as packages that have no discernible module.
//
// The module for a package is determined by the **first** parent directory
// that contains a go.mod.
func modules(filesystemPaths []string) (map[string][]string, []string) {
	// list of module directory => directories of packages it likely contains
	moduledPackages := make(map[string][]string)
	var noModulePkgs []string
	for _, fullPath := range filesystemPaths {
		components := strings.Split(fullPath, "/")

		inModule := false
		for i := len(components); i >= 1; i-- {
			prefixPath := "/" + filepath.Join(components[:i]...)
			if _, err := os.Stat(filepath.Join(prefixPath, "go.mod")); err == nil {
				moduledPackages[prefixPath] = append(moduledPackages[prefixPath], fullPath)
				inModule = true
				break
			}
		}
		if !inModule {
			noModulePkgs = append(noModulePkgs, fullPath)
		}
	}
	return moduledPackages, noModulePkgs
}

// We load file system paths differently, because there is a big difference between
//
//    go list -json ../../foobar
//
// and
//
//    (cd ../../foobar && go list -json .)
//
// Namely, PWD determines which go.mod to use. We want each
// package to use its own go.mod, if it has one.
func loadFSPackages(l ulog.Logger, env golang.Environ, filesystemPaths []string) ([]*packages.Package, error) {
	var absPaths []string
	for _, fsPath := range filesystemPaths {
		absPath, err := filepath.Abs(fsPath)
		if err != nil {
			return nil, fmt.Errorf("could not find package at %q", fsPath)
		}
		absPaths = append(absPaths, absPath)
	}

	var allps []*packages.Package

	// Find each packages' module, and batch package queries together by module.
	//
	// Query all packages that don't have a module at all together, as well.
	//
	// Batching these queries saves a *lot* of time; on the order of
	// several minutes for 30+ commands.
	mods, noModulePkgDirs := modules(absPaths)

	for moduleDir, pkgDirs := range mods {
		pkgs, err := loadFSPkgs(env, moduleDir, pkgDirs...)
		if err != nil {
			return nil, fmt.Errorf("could not find packages in module %s: %v", moduleDir, err)
		}
		for _, pkg := range pkgs {
			allps, err = addPkg(l, allps, pkg)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(noModulePkgDirs) > 0 {
		// The directory we choose can be any dir that does not have a
		// go.mod anywhere in its parent tree.
		vendoredPkgs, err := loadFSPkgs(env, noModulePkgDirs[0], noModulePkgDirs...)
		if err != nil {
			return nil, fmt.Errorf("could not find packages: %v", err)
		}
		for _, p := range vendoredPkgs {
			allps, err = addPkg(l, allps, p)
			if err != nil {
				return nil, err
			}
		}
	}
	return allps, nil
}

func addPkg(l ulog.Logger, plist []*packages.Package, p *packages.Package) ([]*packages.Package, error) {
	if len(p.Errors) > 0 {
		packages.PrintErrors([]*packages.Package{p})
		return plist, fmt.Errorf("failed to add package %v for errors:", p)
	} else if len(p.GoFiles) == 0 {
		l.Printf("Skipping package %v because it has no Go files", p)
	} else if p.Name != "main" {
		l.Printf("Skipping package %v because it is not a command (must be `package main`)", p)
	} else {
		plist = append(plist, p)
	}
	return plist, nil
}

// NewPackages collects package metadata about all named packages.
//
// names can either be directory paths or Go import paths.
func NewPackages(l ulog.Logger, env golang.Environ, names ...string) ([]*bbinternal.Package, error) {
	var goImportPaths []string
	var filesystemPaths []string

	for _, name := range names {
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
			filesystemPaths = append(filesystemPaths, name)
		} else if _, err := os.Stat(name); err == nil {
			filesystemPaths = append(filesystemPaths, name)
		} else {
			goImportPaths = append(goImportPaths, name)
		}
	}

	var ps []*packages.Package
	if len(goImportPaths) > 0 {
		importPkgs, err := loadPkgs(env, "", goImportPaths...)
		if err != nil {
			return nil, fmt.Errorf("failed to load package %v: %v", goImportPaths, err)
		}
		for _, p := range importPkgs {
			ps, err = addPkg(l, ps, p)
		}
	}

	pkgs, err := loadFSPackages(l, env, filesystemPaths)
	if err != nil {
		return nil, fmt.Errorf("could not load packages from file system: %v", err)
	}
	ps = append(ps, pkgs...)

	var ips []*bbinternal.Package
	for _, p := range ps {
		ips = append(ips, bbinternal.NewPackage(path.Base(p.PkgPath), p))
	}
	return ips, nil
}

// loadFSPkgs looks up importDirs packages, making the import path relative to
// `dir`. `go list -json` requires the import path to be relative to the dir
// when the package is outside of a $GOPATH and there is no go.mod in any parent directory.
func loadFSPkgs(env golang.Environ, dir string, importDirs ...string) ([]*packages.Package, error) {
	// Eligibility check: does each directory contain files that are
	// compilable under the current GOROOT/GOPATH/GOOS/GOARCH and build
	// tags?
	//
	// In Go 1.14 and Go 1.15, this is done by packages.Load on a
	// per-package basis, which is why batching queries works out well. In
	// Go 1.13, the entire query fails with no indication of which package
	// made it fail, so we need to filter out commands that do not have
	// compilable files first.
	var compilableImportDirs []string
	for _, importDir := range importDirs {
		f, err := os.Open(importDir)
		if err != nil {
			return nil, err
		}
		names, err := f.Readdirnames(0)
		if errors.Is(err, unix.ENOTDIR) {
			return nil, fmt.Errorf("Go busybox requires a list of directories; failed to read directory %s: %v", importDir, err)
		} else if err != nil {
			return nil, fmt.Errorf("could not determine file names for %s: %v", importDir, err)
		}
		foundOne := false
		for _, name := range names {
			if match, err := env.Context.MatchFile(importDir, name); err != nil {
				// This pretty much only returns an error if
				// the file cannot be opened or read.
				return nil, fmt.Errorf("could not determine Go build constraints of %s: %v", importDir, err)
			} else if match {
				foundOne = true
				break
			}
		}
		if foundOne {
			compilableImportDirs = append(compilableImportDirs, importDir)
		} else {
			log.Printf("Skipping directory %s because build constraints exclude all Go files", importDir)
		}
	}

	if len(compilableImportDirs) == 0 {
		return nil, fmt.Errorf("build constraints excluded all requested commands")
	}

	// Make all paths relative, because packages.Load/`go list -json` does
	// not like absolute paths sometimes.
	var relImportDirs []string
	for _, importDir := range compilableImportDirs {
		relImportDir, err := filepath.Rel(dir, importDir)
		if err != nil {
			return nil, fmt.Errorf("Go package path %s is not relative to %s: %v", importDir, dir, err)
		}

		// N.B. `go list -json cmd/foo` is not the same as `go list -json ./cmd/foo`.
		//
		// The former looks for cmd/foo in $GOROOT or $GOPATH, while
		// the latter looks in the relative directory ./cmd/foo.
		relImportDirs = append(relImportDirs, "./"+relImportDir)
	}
	return loadPkgs(env, dir, relImportDirs...)
}

func loadPkgs(env golang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedFiles | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedCompiledGoFiles | packages.NeedModule | packages.NeedEmbedFiles,
		Env:  append(os.Environ(), env.Env()...),
		Dir:  dir,
	}
	return packages.Load(cfg, patterns...)
}
