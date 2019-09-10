// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package golang is an API to the Go compiler.
package golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Environ struct {
	build.Context

	GO111MODULE string
}

// Default is the default build environment comprised of the default GOPATH,
// GOROOT, GOOS, GOARCH, and CGO_ENABLED values.
func Default() Environ {
	return Environ{
		Context:     build.Default,
		GO111MODULE: os.Getenv("GO111MODULE"),
	}
}

// Package matches a subset of the JSON output of the `go list -json`
// command.
//
// See `go help list` for the full structure.
//
// This currently contains an incomplete list of dependencies.
type Package struct {
	Name       string
	Dir        string
	Deps       []string
	GoFiles    []string
	SFiles     []string
	HFiles     []string
	Goroot     bool
	Root       string
	Module     *Module
	ImportPath string
}

type Module struct {
	Path      string       // module path
	Version   string       // module version
	Versions  []string     // available module versions (with -versions)
	Replace   *Module      // replaced by this module
	Time      *time.Time   // time version was created
	Update    *Module      // available update, if any (with -u)
	Main      bool         // is this the main module?
	Indirect  bool         // is this module only an indirect dependency of main module?
	Dir       string       // directory holding files for this module, if any
	GoMod     string       // path to go.mod file for this module, if any
	GoVersion string       // go version used in module
	Error     *ModuleError // error loading module
}

type ModuleError struct {
	Err string
}

// GoCmd runs a go command in the environment.
func (c Environ) GoCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(c.GOROOT, "bin", "go"), args...)
	cmd.Env = append(os.Environ(), c.Env()...)
	return cmd
}

// Version returns the Go version string that runtime.Version would return for
// the Go compiler in this environ.
func (c Environ) Version() (string, error) {
	cmd := c.GoCmd("version")
	v, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	s := strings.Fields(string(v))
	if len(s) < 3 {
		return "", fmt.Errorf("unknown go version, tool returned weird output for 'go version': %v", string(v))
	}
	return s[2], nil
}

// Find lists all dependencies of the package given by `importPath`.
func (c Environ) Find(pattern string) ([]*Package, error) {
	// The output of this is almost the same as build.Import, except for
	// the dependencies.
	cmd := c.GoCmd("list", "-json", pattern)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list -json %q: %v", pattern, stderr.String())
	}

	var ps []*Package
	for dec := json.NewDecoder(stdout); dec.More(); {
		var p Package
		if err := dec.Decode(&p); err != nil {
			return nil, fmt.Errorf("json unmarshal of go list -json %q: %v", pattern, err)
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (c Environ) FindCmds(pattern string) ([]*Package, error) {
	ps, err := c.Find(pattern)
	if err != nil {
		return nil, err
	}
	var cmds []*Package
	for _, p := range ps {
		if p.Name == "main" {
			cmds = append(cmds, p)
		}
	}
	if len(cmds) == 0 {
		return nil, fmt.Errorf("pattern %q did not find commands, only packages", pattern)
	}
	return cmds, nil
}

func (c Environ) FindOne(pattern string) (*Package, error) {
	ps, err := c.Find(pattern)
	if err != nil {
		return nil, err
	}
	if len(ps) != 1 {
		return nil, fmt.Errorf("pattern %q returned %d packages, wanted one", pattern, len(ps))
	}
	return ps[0], nil
}

func (c Environ) FindOneCmd(pattern string) (*Package, error) {
	ps, err := c.FindCmds(pattern)
	if err != nil {
		return nil, err
	}
	if len(ps) != 1 {
		return nil, fmt.Errorf("pattern %q returned %d packages, wanted one", pattern, len(ps))
	}
	return ps[0], nil
}

// Env returns all environment variables for invoking a Go command.
func (c Environ) Env() []string {
	var env []string
	if c.GOARCH != "" {
		env = append(env, fmt.Sprintf("GOARCH=%s", c.GOARCH))
	}
	if c.GOOS != "" {
		env = append(env, fmt.Sprintf("GOOS=%s", c.GOOS))
	}
	if c.GOPATH != "" {
		env = append(env, fmt.Sprintf("GOPATH=%s", c.GOPATH))
	}
	var cgo int8
	if c.CgoEnabled {
		cgo = 1
	}
	env = append(env, fmt.Sprintf("CGO_ENABLED=%d", cgo))
	env = append(env, fmt.Sprintf("GO111MODULE=%s", c.GO111MODULE))

	if c.GOROOT != "" {
		env = append(env, fmt.Sprintf("GOROOT=%s", c.GOROOT))

		// If GOROOT is set to a different version of Go, we must
		// ensure that $GOROOT/bin is also in path to make the "go"
		// binary available to golang.org/x/tools/packages.
		env = append(env, fmt.Sprintf("PATH=%s:%s", filepath.Join(c.GOROOT, "bin"), os.Getenv("PATH")))
	}
	return env
}

// String returns all environment variables for Go invocations.
func (c Environ) String() string {
	return strings.Join(c.Env(), " ")
}

// Optional arguments to Environ.Build.
type BuildOpts struct {
	// NoStrip builds an unstripped binary.
	NoStrip bool
	// ExtraArgs to `go build`.
	ExtraArgs []string
}

// Build compiles the package given by `importPath`, writing the build object
// to `binaryPath`.
func (c Environ) Build(importPath string, binaryPath string, opts BuildOpts) error {
	p, err := c.FindOneCmd(importPath)
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

		// Force rebuilding of packages.
		"-a",

		// Strip all symbols, and don't embed a Go build ID to be reproducible.
		"-ldflags", "-s -w -buildid=",

		"-o", binaryPath,
		"-installsuffix", "uroot",

		"-gcflags=all=-l", // Disable "function inlining" to get a smaller binary
	}
	if !opts.NoStrip {
		args = append(args, `-ldflags=-s -w`) // Strip all symbols.
	}

	v, err := c.Version()
	if err != nil {
		return err
	}

	// Reproducible builds: Trim any GOPATHs out of the executable's
	// debugging information.
	//
	// E.g. Trim /tmp/bb-*/ from /tmp/bb-12345567/src/github.com/...
	if strings.Contains(v, "go1.13") || strings.Contains(v, "go1.14") || strings.Contains(v, "gotip") {
		args = append(args, "-trimpath")
	} else {
		args = append(args, "-gcflags", fmt.Sprintf("-trimpath=%s", c.GOPATH))
		args = append(args, "-asmflags", fmt.Sprintf("-trimpath=%s", c.GOPATH))
	}

	if len(c.BuildTags) > 0 {
		args = append(args, []string{"-tags", strings.Join(c.BuildTags, " ")}...)
	}
	// We always set the working directory, so this is always '.'.
	args = append(args, ".")

	cmd := c.GoCmd(args...)
	cmd.Dir = dirPath

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package in %q: %v, %v", dirPath, string(o), err)
	}
	return nil
}
