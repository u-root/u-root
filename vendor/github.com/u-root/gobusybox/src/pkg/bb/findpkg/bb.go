// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package findpkg finds packages from user-input strings that are either file
// paths or Go package paths.
//
// findpkg supports globs and exclusions in addition to the normal `go list`
// syntax, as described in the [NewPackage] documentation.
package findpkg

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/bb/bbinternal"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/ulog"
	"golang.org/x/tools/go/packages"
	"mvdan.cc/sh/v3/shell"
)

func addPkg(plist []*packages.Package, p *packages.Package) ([]*packages.Package, error) {
	if len(p.Errors) > 0 {
		var merr error
		for _, e := range p.Errors {
			merr = errors.Join(merr, e)
		}
		return plist, fmt.Errorf("failed to add package %v for errors: %w", p, merr)
	} else if len(p.GoFiles) > 0 {
		plist = append(plist, p)
	}
	return plist, nil
}

func newPackages(l ulog.Logger, genv *golang.Environ, env Env, patterns ...string) ([]*packages.Package, error) {
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
	importPkgs, err := loadPkgs(genv, paths...)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %v: %v", paths, err)
	}
	var pkgs []*packages.Package
	for _, p := range importPkgs {
		pkgs, err = addPkg(pkgs, p)
		if err != nil {
			return nil, err
		}
	}
	return pkgs, nil
}

// NewPackages collects package metadata about all named packages.
//
// names can either be directory paths or Go import paths, with globs.
//
// It skips directories that do not have Go files subject to the build
// constraints in env and logs a "Skipping package {}" statement about such
// directories/packages.
//
// All given names have to be resolvable by GOPATH or by Go modules. Generally,
// if `go list <name>` works, it should work here. (Except go list does not
// support globs or exclusions.)
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

func loadPkgs(env *golang.Environ, patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedName | packages.NeedImports | packages.NeedFiles | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedCompiledGoFiles | packages.NeedModule | packages.NeedEmbedFiles
	return env.Lookup(mode, patterns...)
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
				merr = errors.Join(merr, fmt.Errorf("package %s: %w", p.PkgPath, e))
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
func lookupPkgNameAndFiles(env *golang.Environ, patterns ...string) ([]*packages.Package, error) {
	return env.Lookup(packages.NeedName|packages.NeedFiles, patterns...)
}

func couldBeGlob(s string) bool {
	return strings.ContainsAny(s, "*?[") || strings.Contains(s, `\\`)
}

// lookupPkgsWithGlob resolves globs in Go package paths to a realized list of
// Go command paths. It may return a list that contains errors.
//
// Precondition: couldBeGlob(pattern) is true
func lookupPkgsWithGlob(env *golang.Environ, pattern string) ([]*packages.Package, error) {
	elems := strings.Split(pattern, "/")

	globIndex := 0
	for i, e := range elems {
		if couldBeGlob(e) {
			globIndex = i
			break
		}
	}

	nonGlobPath := strings.Join(append(elems[:globIndex], "..."), "/")

	pkgs, err := lookupPkgNameAndFiles(env, nonGlobPath)
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
func lookupCompilablePkgsWithGlob(l ulog.Logger, env *golang.Environ, patterns ...string) ([]string, error) {
	var pkgs []*packages.Package
	// Batching saves time. Patterns with globs cannot be batched.
	//
	// When you batch requests you cannot attribute which result came from
	// which individual request. For globs, we need to be able to do
	// path.Match-ing on the results. So no batching of globs.
	var batchedPatterns []string
	for _, pattern := range patterns {
		if couldBeGlob(pattern) {
			ps, err := lookupPkgsWithGlob(env, pattern)
			if err != nil {
				return nil, err
			}
			pkgs = append(pkgs, ps...)
		} else {
			batchedPatterns = append(batchedPatterns, pattern)
		}
	}
	if len(batchedPatterns) > 0 {
		ps, err := lookupPkgNameAndFiles(env, batchedPatterns...)
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

func filterGoPaths(l ulog.Logger, env *golang.Environ, gopathIncludes, gopathExcludes []string) ([]string, error) {
	goInc, err := lookupCompilablePkgsWithGlob(l, env, gopathIncludes...)
	if err != nil {
		return nil, err
	}

	goExc, err := lookupCompilablePkgsWithGlob(l, env, gopathExcludes...)
	if err != nil {
		return nil, err
	}
	return excludePaths(goInc, goExc), nil
}

var errNoMatch = fmt.Errorf("no Go commands match the given patterns")

// Glob evaluates the given pattern.
func (e Env) Glob(l ulog.Logger, pattern string) (isPath bool, absPaths []string) {
	prefixes := append([]string{""}, e.GBBPath...)

	// Special sauce for backwards compatibility with old u-root behavior:
	// if urootSource is set, try to catch Go package paths and convert
	// them to file system lookups.
	if len(e.URootSource) > 0 {
		// Prefer urootSource to gbbPaths in this case.
		prefixes = append([]string{"", e.URootSource}, e.GBBPath...)
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
}

func (e Env) String() string {
	return fmt.Sprintf("GBB_PATH=%s UROOT_SOURCE=%s", strings.Join(e.GBBPath, ":"), e.URootSource)
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
		GBBPath:     gbbPaths,
		URootSource: os.Getenv("UROOT_SOURCE"),
	}
}

func splitExclusions(patterns []string) (include []string, exclude []string) {
	for _, pattern := range patterns {
		isExclude := strings.HasPrefix(pattern, "-")
		if isExclude {
			exclude = append(exclude, pattern[1:])
		} else {
			include = append(include, pattern)
		}
	}
	return
}

func expand(patterns []string) []string {
	var p []string
	for _, pattern := range patterns {
		fields, err := shell.Fields(pattern, func(_ string) string {
			return ""
		})
		if err != nil {
			p = append(p, pattern)
		} else {
			p = append(p, fields...)
		}
	}
	return p
}

func glob(l ulog.Logger, env Env, patterns []string) []string {
	var ret []string
	for _, pattern := range patterns {
		if match, directories := env.Glob(l, pattern); len(directories) > 0 {
			ret = append(ret, directories...)
		} else if !match {
			ret = append(ret, pattern)
		}
	}
	return ret
}

// ResolveGlobs takes a list of Go paths and directories that may
// include globs and returns a valid list of Go commands (either addressed by
// Go package path or directory path).
//
// It returns only packages that have Go files subject to
// the build constraints in env and logs a "Skipping package {}" statement
// about packages that are excluded due to build constraints.
//
// ResolveGlobs always returns normalized Go package paths.
//
// ResolveGlobs should work in all cases that `go list` works.
//
// See NewPackages for allowed formats.
func ResolveGlobs(l ulog.Logger, genv *golang.Environ, env Env, patterns []string) ([]string, error) {
	if genv == nil {
		return nil, fmt.Errorf("Go build environment must be specified")
	}
	includes, excludes := splitExclusions(patterns)
	includes = glob(l, env, expand(includes))
	excludes = glob(l, env, expand(excludes))

	paths, err := filterGoPaths(l, genv, includes, excludes)
	if err != nil {
		if strings.Contains(err.Error(), "go.mod file not found") {
			return nil, fmt.Errorf("%w: gobusybox has removed previous multi-module functionality in favor of Go workspaces -- read https://github.com/u-root/gobusybox#path-resolution--multi-module-builds for more", err)
		}
		return nil, err
	}

	if len(paths) == 0 {
		return nil, errNoMatch
	}
	sort.Strings(paths)
	return paths, nil
}

// Modules returns a list of module directories => directories of packages
// inside that module as well as packages that have no discernible module.
//
// The module for a package is determined by the **first** parent directory
// that contains a go.mod.
func Modules(paths []string) (map[string][]string, []string) {
	// list of module directory => directories of packages it likely contains
	moduledPackages := make(map[string][]string)
	var noModulePkgs []string
	for _, fullPath := range paths {
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

func globPaths(l ulog.Logger, env Env, patterns []string) []string {
	var ret []string
	for _, pattern := range patterns {
		if match, directories := env.Glob(l, pattern); len(directories) > 0 {
			ret = append(ret, directories...)
		} else if !match {
			l.Printf("No match found for %v", pattern)
		}
	}
	return ret
}

// GlobPaths resolves file path globs in env with exclusions and shell
// expansions.
func GlobPaths(l ulog.Logger, env Env, patterns ...string) []string {
	includes, excludes := splitExclusions(patterns)
	includes = globPaths(l, env, expand(includes))
	excludes = globPaths(l, env, expand(excludes))
	paths := excludePaths(includes, excludes)
	return paths
}
