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

type Environ struct {
	build.Context
}

// Default is the default build environment comprised of the default GOPATH,
// GOROOT, GOOS, GOARCH, and CGO_ENABLED values.
func Default() Environ {
	return Environ{Context: build.Default}
}

// PackageByPath retrieves information about a package by its file system path.
//
// `path` is assumed to be the directory containing the package.
func (c Environ) PackageByPath(path string) (*build.Package, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return c.Context.ImportDir(abs, 0)
}

// Package retrieves information about a package by its Go import path.
func (c Environ) Package(importPath string) (*build.Package, error) {
	return c.Context.Import(importPath, "", 0)
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

func (c Environ) goCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(c.GOROOT, "bin", "go"), args...)
	cmd.Env = append(os.Environ(), c.Env()...)
	return cmd
}

// Deps lists all dependencies of the package given by `importPath`.
func (c Environ) Deps(importPath string) (*ListPackage, error) {
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

func (c Environ) Env() []string {
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

func (c Environ) String() string {
	return strings.Join(c.Env(), " ")
}

// Optional arguments to Environ.Build.
type BuildOpts struct {
	// ExtraArgs to `go build`.
	ExtraArgs []string
}

// Build compiles the package given by `importPath`, writing the build object
// to `binaryPath`.
func (c Environ) Build(importPath string, binaryPath string, opts BuildOpts) error {
	p, err := c.Package(importPath)
	if err != nil {
		return err
	}

	return c.BuildDir(p.Dir, binaryPath, opts)
}

// BuildDir compiles the package in the directory `dirPath`, writing the build
// object to `binaryPath`.
func (c Environ) BuildDir(dirPath string, binaryPath string, opts BuildOpts) error {
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
