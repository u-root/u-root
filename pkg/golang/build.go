// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Package is the basic information a build system should provide about a Go
// package.
type Package struct {
	Name       string
	Dir        string
	ImportPath string
	GoFiles    []string
	IsCommand  bool
}

// Environ is the shared toolchain interface for all build systems building Go
// packages.
//
// This is to have a common denominator between the standard Go toolchain and
// build systems such as bazel, buck, or blaze.
type Environ interface {
	// Package retrieves information about a package by its Go import path.
	Package(importPath string) (*Package, error)
}

// PuppetEnviron is an Environ that can be used with the u-root busybox build
// system by build systems that don't use the standard Go toolchain, such as
// blaze, bazel, or buck.
type PuppetEnviron struct {
	packages map[string]*Package
}

var _ Environ = &PuppetEnviron{}

// NewPuppetEnviron returns a new puppet toolchain environment.
func NewPuppetEnviron(pkgs map[string]*Package) *PuppetEnviron {
	return &PuppetEnviron{
		packages: pkgs,
	}
}

// Package implements Environ.Package.
func (p *PuppetEnviron) Package(importPath string) (*Package, error) {
	pkg, ok := p.packages[importPath]
	if ok {
		return pkg, nil
	}
	return nil, fmt.Errorf("build system did not pass enough information to find package %v", importPath)
}

// StandardGoEnviron gives a Go API to interact with the standard Go toolchain
// and implements Environ.
type StandardGoEnviron struct {
	build.Context
}

var _ Environ = &StandardGoEnviron{}

// Default is the default build environment comprised of the default GOPATH,
// GOROOT, GOOS, GOARCH, and CGO_ENABLED values.
func Default() *StandardGoEnviron {
	return &StandardGoEnviron{Context: build.Default}
}

// PackageByPath retrieves information about a package by its file system path.
//
// `path` is assumed to be the directory containing the package.
func (c *StandardGoEnviron) PackageByPath(path string) (*build.Package, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return c.Context.ImportDir(abs, 0)
}

// GoPackage retrieves information about a package by its Go import path.
func (c *StandardGoEnviron) GoPackage(importPath string) (*build.Package, error) {
	return c.Context.Import(importPath, "", 0)
}

// Package implements Environ.Package.
func (c *StandardGoEnviron) Package(importPath string) (*Package, error) {
	p, err := c.GoPackage(importPath)
	if err != nil {
		return nil, err
	}
	return &Package{
		Name:       p.Name,
		Dir:        p.Dir,
		ImportPath: p.ImportPath,
		GoFiles:    p.GoFiles,
		IsCommand:  p.IsCommand(),
	}, nil
}

// ListPackage matches a subset of the JSON output of the `go list -json`
// command.
//
// See `go help list` for the full structure.
//
// This currently contains an incomplete list of dependencies.
type ListPackage struct {
	Dir        string
	Deps       []string
	GoFiles    []string
	SFiles     []string
	HFiles     []string
	Goroot     bool
	Root       string
	ImportPath string
}

func (c StandardGoEnviron) goCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(c.GOROOT, "bin", "go"), args...)
	cmd.Env = append(os.Environ(), c.Env()...)
	return cmd
}

// Deps lists all dependencies of the package given by `importPath`.
func (c StandardGoEnviron) Deps(importPath string) (*ListPackage, error) {
	// The output of this is almost the same as build.Import, except for
	// the dependencies.
	cmd := c.goCmd("list", "-json", importPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var p ListPackage
	if err := json.Unmarshal(out, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c StandardGoEnviron) Env() []string {
	var env []string
	if c.GOARCH != "" {
		env = append(env, fmt.Sprintf("GOARCH=%s", c.GOARCH))
	}
	if c.GOOS != "" {
		env = append(env, fmt.Sprintf("GOOS=%s", c.GOOS))
	}
	if c.GOROOT != "" {
		env = append(env, fmt.Sprintf("GOROOT=%s", c.GOROOT))
	}
	if c.GOPATH != "" {
		env = append(env, fmt.Sprintf("GOPATH=%s", c.GOPATH))
	}
	var cgo int8
	if c.CgoEnabled {
		cgo = 1
	}
	env = append(env, fmt.Sprintf("CGO_ENABLED=%d", cgo))
	return env
}

// String implements fmt.Stringer.
func (c StandardGoEnviron) String() string {
	return strings.Join(c.Env(), " ")
}

// Optional arguments to Environ.Build.
type BuildOpts struct {
	// ExtraArgs are extra arguments to pass to `go build`.
	ExtraArgs []string
}

// Build compiles the package given by `importPath`, writing the build object
// to `binaryPath`.
func (c StandardGoEnviron) Build(importPath string, binaryPath string, opts BuildOpts) error {
	p, err := c.Package(importPath)
	if err != nil {
		return err
	}

	return c.BuildDir(p.Dir, binaryPath, opts)
}

// BuildDir compiles the package in the directory `dirPath`, writing the build
// object to `binaryPath`.
func (c StandardGoEnviron) BuildDir(dirPath string, binaryPath string, opts BuildOpts) error {
	args := []string{
		"build",
		"-a", // Force rebuilding of packages.
		"-o", binaryPath,
		"-installsuffix", "uroot",
		"-ldflags", "-s -w", // Strip all symbols.
	}
	if opts.ExtraArgs != nil {
		args = append(args, opts.ExtraArgs...)
	}
	// We always set the working directory, so this is always '.'.
	args = append(args, ".")

	cmd := c.goCmd(args...)
	cmd.Dir = dirPath

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package in %q: %v, %v", dirPath, string(o), err)
	}
	return nil
}
