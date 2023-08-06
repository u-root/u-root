// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package findpkg finds packages from user-input strings that are either file
// paths or Go package paths.
package findpkg

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/u-root/gobusybox/src/pkg/bb/bbinternal"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/ulog"
	"golang.org/x/tools/go/packages"
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

func loadRelative(moduleDir string, pkgDirs []string, loadFunc func(moduleDir string, dirs []string) error) error {
	// Make all paths relative, because packages.Load/`go list -json` does
	// not like absolute paths when what's being looked up is outside of
	// GOPATH. It's fine though when those paths are relative to the PWD.
	// Don't ask me why...
	//
	// E.g.
	// * package $HOME/u-root/cmds/core/ip: -: import "$HOME/u-root/cmds/core/ip": cannot import absolute path
	var relPkgDirs []string
	for _, pkgDir := range pkgDirs {
		relPkgDir, err := filepath.Rel(moduleDir, pkgDir)
		if err != nil {
			return fmt.Errorf("Go package path %s is not relative to directory %s: %v", pkgDir, moduleDir, err)
		}

		// N.B. `go list -json cmd/foo` is not the same as `go list -json ./cmd/foo`.
		//
		// The former looks for cmd/foo in $GOROOT or $GOPATH, while
		// the latter looks in the relative directory ./cmd/foo.
		relPkgDirs = append(relPkgDirs, "./"+relPkgDir)
	}
	return loadFunc(moduleDir, relPkgDirs)
}

// Find each packages' module, and batch package queries together by module.
//
// Query all packages that don't have a module at all together, as well.
//
// Batching these queries saves a *lot* of time; on the order of
// several minutes for 30+ commands.
func batchFSPackages(l ulog.Logger, absPaths []string, loadFunc func(moduleDir string, dirs []string) error) error {
	mods, noModulePkgDirs := modules(absPaths)

	for moduleDir, pkgDirs := range mods {
		if err := loadRelative(moduleDir, pkgDirs, loadFunc); err != nil {
			return err
		}
	}

	if len(noModulePkgDirs) > 0 {
		if err := loadRelative(noModulePkgDirs[0], noModulePkgDirs, loadFunc); err != nil {
			return err
		}
	}
	return nil
}

// We look up file system paths differently, because there is a big difference between
//
//	go list -json ../../foobar
//
// and
//
//	(cd ../../foobar && go list -json .)
//
// Namely, PWD determines which go.mod to use. We want each
// package to use its own go.mod, if it has one.
//
// The easiest implementation would be to do (cd $packageDir && go list -json
// .), however doing that N times is very expensive -- takes several minutes
// for 30 packages. So here, we figure out every module involved and do one
// query per module and one query for everything that isn't in a module.
func batchLoadFSPackages(l ulog.Logger, env *golang.Environ, absPaths []string) ([]*packages.Package, error) {
	var allps []*packages.Package

	err := batchFSPackages(l, absPaths, func(moduleDir string, packageDirs []string) error {
		pkgs, err := loadPkgs(env, moduleDir, packageDirs...)
		if err != nil {
			return fmt.Errorf("could not find packages in module %s: %v", moduleDir, err)
		}
		for _, pkg := range pkgs {
			allps, err = addPkg(l, allps, pkg)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allps, nil
}

func addPkg(l ulog.Logger, plist []*packages.Package, p *packages.Package) ([]*packages.Package, error) {
	if len(p.Errors) > 0 {
		var merr error
		for _, e := range p.Errors {
			merr = multierror.Append(merr, e)
		}
		return plist, fmt.Errorf("failed to add package %v for errors: %v", p, merr)
	} else if len(p.GoFiles) > 0 {
		plist = append(plist, p)
	}
	return plist, nil
}

func newPackages(l ulog.Logger, genv *golang.Environ, env Env, patterns ...string) ([]*packages.Package, error) {
	var goImportPaths []string
	var filesystemPaths []string

	// Two steps:
	//
	// 1. Resolve globs, filter packages with build constraints.
	//    Produce an explicit list of packages.
	//
	// 2. Look up every piece of information necessary from those packages.
	//    (Includes optimizations to reduce the amount of time it takes to
	//    do type-checking, etc.)

	// Step 1.
	paths, err := ResolveGlobs(l, genv, env, patterns)
	if err != nil {
		return nil, err
	}

	// Step 2.
	for _, name := range paths {
		// ResolveGlobs returns either an absolute file system path or
		// a Go import path.
		if strings.HasPrefix(name, "/") {
			filesystemPaths = append(filesystemPaths, name)
		} else {
			goImportPaths = append(goImportPaths, name)
		}
	}

	var ps []*packages.Package
	if len(goImportPaths) > 0 {
		importPkgs, err := loadPkgs(genv, env.WorkingDirectory, goImportPaths...)
		if err != nil {
			return nil, fmt.Errorf("failed to load package %v: %v", goImportPaths, err)
		}
		for _, p := range importPkgs {
			ps, err = addPkg(l, ps, p)
			if err != nil {
				return nil, err
			}
		}
	}

	pkgs, err := batchLoadFSPackages(l, genv, filesystemPaths)
	if err != nil {
		return nil, fmt.Errorf("could not load packages from file system: %v", err)
	}
	ps = append(ps, pkgs...)
	return ps, nil
}

// NewPackages collects package metadata about all named packages.
//
// names can either be directory paths or Go import paths, with globs.
//
// It skips directories that do not have Go files subject to the build
// constraints in env and logs a "Skipping package {}" statement about such
// directories/packages.
//
// Allowed formats for names:
//
//   - relative and absolute paths including globs following Go's
//     filepath.Match format.
//
//   - Go package paths; e.g. github.com/u-root/u-root/cmds/core/ls
//
//   - Globs of Go package paths, e.g github.com/u-root/u-root/cmds/i* (using
//     path.Match format).
//
//   - Go package path expansions with ..., e.g.
//     github.com/u-root/u-root/cmds/core/...
//
//   - file system paths (with globs in filepath.Match format) relative to
//     GBB_PATH, e.g. cmds/core/ls if GBB_PATH contains $HOME/u-root.
//
//   - backwards compatibility: UROOT_SOURCE is a GBB_PATH, and patterns that
//     begin with github.com/u-root/u-root/ will attempt to use UROOT_SOURCE
//     first to find Go commands within.
//
// If a pattern starts with "-", it excludes the matching package(s).
//
// Globs of Go package paths must be within module boundaries to give accurate
// results, i.e. a glob that spans 2 Go modules may give unpredictable results.
//
// Examples of valid inputs:
//
//   - ./foobar
//
//   - ./foobar/glob*
//
//   - github.com/u-root/u-root/cmds/core/...
//
//   - github.com/u-root/u-root/cmds/core/ip
//
//   - github.com/u-root/u-root/cmds/core/g*lob
//
//   - GBB_PATH=$HOME/u-root:$HOME/yourproject cmds/core/* cmd/foobar
//
//   - UROOT_SOURCE=$HOME/u-root github.com/u-root/u-root/cmds/core/ip
func NewPackages(l ulog.Logger, genv *golang.Environ, env Env, names ...string) ([]*bbinternal.Package, error) {
	if genv == nil {
		return nil, fmt.Errorf("Go build environment must be specified")
	}
	ps, err := newPackages(l, genv, env, names...)
	if err != nil {
		return nil, err
	}

	var ips []*bbinternal.Package
	for _, p := range ps {
		ips = append(ips, bbinternal.NewPackage(path.Base(p.PkgPath), p))
	}
	return ips, nil
}

func loadPkgs(env *golang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedFiles | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedCompiledGoFiles | packages.NeedModule | packages.NeedEmbedFiles,
		Env:  append(os.Environ(), env.Env()...),
		Dir:  dir,
	}
	if len(env.Context.BuildTags) > 0 {
		tags := fmt.Sprintf("-tags=%s", strings.Join(env.Context.BuildTags, ","))
		cfg.BuildFlags = []string{tags}
	}
	return packages.Load(cfg, patterns...)
}

func filterDirectoryPaths(l ulog.Logger, env *golang.Environ, directories []string, excludes []string) ([]string, error) {
	// Eligibility check: does each directory contain files that are
	// compilable under the current GOROOT/GOPATH/GOOS/GOARCH and build
	// tags?
	//
	// We filter this out first, because while packages.Load will give us
	// an error for this, it is not distinguishable from other errors. We
	// would like to give only a warning for these.
	//
	// This eligibility check requires Go 1.15, as before Go 1.15 the
	// package loader would return an error "cannot find package" for
	// packages not meeting build constraints.
	var allps []*packages.Package
	err := batchFSPackages(l, directories, func(moduleDir string, packageDirs []string) error {
		pkgs, err := lookupPkgNameAndFiles(env, moduleDir, packageDirs...)
		if err != nil {
			return fmt.Errorf("could not look up packages %q: %v", packageDirs, err)
		}
		allps = append(allps, pkgs...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	eligiblePkgs, err := checkEligibility(l, allps)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, p := range eligiblePkgs {
		paths = append(paths, filepath.Dir(p.GoFiles[0]))
	}
	return excludePaths(paths, excludes), nil
}

func checkEligibility(l ulog.Logger, pkgs []*packages.Package) ([]*packages.Package, error) {
	var goodPkgs []*packages.Package
	var merr error
	for _, p := range pkgs {
		// If there's a build constraint issue, short out early and
		// neither add the package nor add an error -- just log a skip
		// note.
		if len(p.GoFiles) == 0 && len(p.IgnoredFiles) > 0 {
			l.Printf("Skipping package %s because build constraints exclude all Go files", p.PkgPath)
		} else if len(p.Errors) == 0 {
			if p.Name != "main" {
				l.Printf("Skipping package %s because it is not a command (must be `package main`)", p.PkgPath)
			} else {
				goodPkgs = append(goodPkgs, p)
			}
		} else {
			// We'll definitely return an error in the end, but
			// we're not returning early because we want to give
			// the user as much information as possible.
			for _, e := range p.Errors {
				merr = multierror.Append(merr, fmt.Errorf("package %s: %w", p.PkgPath, e))
			}
		}
	}
	if merr != nil {
		return nil, merr
	}
	return goodPkgs, nil
}

func excludePaths(paths []string, exclusions []string) []string {
	excludes := map[string]struct{}{}
	for _, p := range exclusions {
		excludes[p] = struct{}{}
	}

	var result []string
	for _, p := range paths {
		if _, ok := excludes[p]; !ok {
			result = append(result, p)
		}
	}
	return result
}

// Just looking up the stuff that doesn't take forever to parse.
func lookupPkgNameAndFiles(env *golang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
		Env:  append(os.Environ(), env.Env()...),
		Dir:  dir,
	}
	return packages.Load(cfg, patterns...)
}

func couldBeGlob(s string) bool {
	return strings.ContainsAny(s, "*?[") || strings.Contains(s, `\\`)
}

// lookupPkgsWithGlob resolves globs in Go package paths to a realized list of
// Go command paths. It may return a list that contains errors.
//
// Precondition: couldBeGlob(pattern) is true
func lookupPkgsWithGlob(env *golang.Environ, wd string, pattern string) ([]*packages.Package, error) {
	elems := strings.Split(pattern, "/")

	globIndex := 0
	for i, e := range elems {
		if couldBeGlob(e) {
			globIndex = i
			break
		}
	}

	nonGlobPath := strings.Join(append(elems[:globIndex], "..."), "/")

	pkgs, err := lookupPkgNameAndFiles(env, wd, nonGlobPath)
	if err != nil {
		return nil, fmt.Errorf("%q is neither package or path/glob -- could not lookup %q (import path globs have to be within modules): %v", pattern, nonGlobPath, err)
	}

	// Apply the glob.
	var filteredPkgs []*packages.Package
	for _, p := range pkgs {
		if matched, err := path.Match(pattern, p.PkgPath); err != nil {
			return nil, fmt.Errorf("could not match %q to %q: %v", pattern, p.PkgPath, err)
		} else if matched {
			filteredPkgs = append(filteredPkgs, p)
		}
	}
	return filteredPkgs, nil
}

// lookupCompilablePkgsWithGlob resolves Go package path globs to a realized
// list of Go command paths. It filters out packages that have no files
// matching our build constraints and other errors.
func lookupCompilablePkgsWithGlob(l ulog.Logger, env *golang.Environ, wd string, patterns ...string) ([]string, error) {
	var pkgs []*packages.Package
	// Batching saves time. Patterns with globs cannot be batched.
	//
	// When you batch requests you cannot attribute which result came from
	// which individual request. For globs, we need to be able to do
	// path.Match-ing on the results. So no batching of globs.
	var batchedPatterns []string
	for _, pattern := range patterns {
		if couldBeGlob(pattern) {
			ps, err := lookupPkgsWithGlob(env, wd, pattern)
			if err != nil {
				return nil, err
			}
			pkgs = append(pkgs, ps...)
		} else {
			batchedPatterns = append(batchedPatterns, pattern)
		}
	}
	if len(batchedPatterns) > 0 {
		ps, err := lookupPkgNameAndFiles(env, wd, batchedPatterns...)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, ps...)
	}

	eligiblePkgs, err := checkEligibility(l, pkgs)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, p := range eligiblePkgs {
		paths = append(paths, p.PkgPath)
	}
	return paths, nil
}

func filterGoPaths(l ulog.Logger, env *golang.Environ, wd string, gopathIncludes, gopathExcludes []string) ([]string, error) {
	goInc, err := lookupCompilablePkgsWithGlob(l, env, wd, gopathIncludes...)
	if err != nil {
		return nil, err
	}

	goExc, err := lookupCompilablePkgsWithGlob(l, env, wd, gopathExcludes...)
	if err != nil {
		return nil, err
	}
	return excludePaths(goInc, goExc), nil
}

var errNoMatch = fmt.Errorf("no Go commands match the given patterns")

func findDirectoryMatches(l ulog.Logger, env Env, pattern string) (bool, []string) {
	prefixes := append([]string{""}, env.GBBPath...)

	// Special sauce for backwards compatibility with old u-root behavior:
	// if urootSource is set, try to catch Go package paths and convert
	// them to file system lookups.
	if len(env.URootSource) > 0 {
		// Prefer urootSource to gbbPaths in this case.
		prefixes = append([]string{"", env.URootSource}, env.GBBPath...)
		pattern = strings.TrimPrefix(pattern, "github.com/u-root/u-root/")
	}

	// We track matches because we want to ignore individual files.
	//
	// Go can look up and compile individual files, which show up as a
	// package called "command-line-arguments". That does not make sense in
	// a busybox.
	//
	// We don't return an error for convenience of globbing directories
	// that may have files & directories with commands in them.
	//
	// We track the bool so that even if no matches are found for this
	// pattern, it is not looked up as a Go package path later.
	foundMatch := false
	for _, prefix := range prefixes {
		matches, _ := filepath.Glob(filepath.Join(prefix, pattern))
		if len(matches) == 0 {
			continue
		}

		foundMatch = true
		var dirs []string
		for _, match := range matches {
			fileInfo, _ := os.Stat(match)
			if !fileInfo.IsDir() {
				l.Printf("Skipping %s because it is not a directory", match)
				continue
			}
			absPath, _ := filepath.Abs(match)
			dirs = append(dirs, absPath)
		}
		if len(dirs) > 0 {
			return true, dirs
		}
	}
	return foundMatch, nil
}

// Env is configuration for package lookups.
type Env struct {
	// GBBPath provides directories in which to look for Go commands.
	//
	// The default is to use a colon-separated list from the env var
	// GBB_PATH.
	GBBPath []string

	// URootSource is a special GBBPath. It's a directory that will be used
	// to look for u-root commands. If a u-root command is given as a
	// pattern with the "github.com/u-root/u-root/" Go package path prefix,
	// URootSource will be used to find the command source.
	//
	// The default is to use UROOT_SOURCE env var.
	URootSource string

	// WorkingDirectory is the directory used for module-enabled `go list`
	// lookups. The go.mod in this directory (or one of its parents) is
	// used to resolve Go package paths.
	WorkingDirectory string
}

func (e Env) String() string {
	return fmt.Sprintf("GBB_PATH=%s UROOT_SOURCE=%s PWD=%s", strings.Join(e.GBBPath, ":"), e.URootSource, e.WorkingDirectory)
}

// DefaultEnv is the default environment derived from environment variables and
// the current working directory.
func DefaultEnv() Env {
	gbbPath := os.Getenv("GBB_PATH")
	// strings.Split("", ":") is []string{""}, but we want nil
	var gbbPaths []string
	if gbbPath != "" {
		gbbPaths = strings.Split(gbbPath, ":")
	}
	return Env{
		GBBPath:          gbbPaths,
		URootSource:      os.Getenv("UROOT_SOURCE"),
		WorkingDirectory: "",
	}
}

// ResolveGlobs takes a list of Go paths and directories that may
// include globs and returns a valid list of Go commands (either addressed by
// Go package path or directory path).
//
// It returns only directories that have Go files subject to
// the build constraints in env and logs a "Skipping package {}" statement
// about packages that are excluded due to build constraints.
//
// ResolveGlobs always returns either an absolute file system path and
// normalized Go package paths. The return list may be mixed.
//
// See NewPackages for allowed formats.
func ResolveGlobs(logger ulog.Logger, genv *golang.Environ, env Env, patterns []string) ([]string, error) {
	if genv == nil {
		return nil, fmt.Errorf("Go build environment must be specified")
	}
	var dirIncludes []string
	var dirExcludes []string
	var gopathIncludes []string
	var gopathExcludes []string
	for _, pattern := range patterns {
		isExclude := strings.HasPrefix(pattern, "-")
		if isExclude {
			pattern = pattern[1:]
		}
		if hasFileMatch, directories := findDirectoryMatches(logger, env, pattern); len(directories) > 0 {
			if !isExclude {
				dirIncludes = append(dirIncludes, directories...)
			} else {
				dirExcludes = append(dirExcludes, directories...)
			}
		} else if !hasFileMatch {
			if !isExclude {
				gopathIncludes = append(gopathIncludes, pattern)
			} else {
				gopathExcludes = append(gopathExcludes, pattern)
			}
		}
	}

	directories, err := filterDirectoryPaths(logger, genv, dirIncludes, dirExcludes)
	if err != nil {
		return nil, err
	}

	gopaths, err := filterGoPaths(logger, genv, env.WorkingDirectory, gopathIncludes, gopathExcludes)
	if err != nil {
		return nil, err
	}

	result := append(directories, gopaths...)
	if len(result) == 0 {
		return nil, errNoMatch
	}
	sort.Strings(result)
	return result, nil
}
