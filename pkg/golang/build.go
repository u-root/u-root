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

// Deps lists all dependencies of the package given by `importPath`.
func (c Environ) Deps(importPath string) (*ListPackage, error) {
	// The output of this is almost the same as build.Import, except for
	// the dependencies.
	cmd := exec.Command("go", "list", "-json", importPath)
	env := os.Environ()
	env = append(env, c.Env()...)
	cmd.Env = env
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

	cmd := exec.Command("go", args...)
	cmd.Dir = p.Dir
	cmd.Env = append(os.Environ(), c.Env()...)

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package %v: %v, %v", importPath, string(o), err)
	}
	return nil
}
