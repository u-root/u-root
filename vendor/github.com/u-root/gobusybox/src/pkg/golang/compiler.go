// Copyright 2015-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package golang is an API to the Go compiler.
package golang

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type CompilerType int

const (
	CompilerGo CompilerType = iota
	CompilerTinygo
	CompilerUnkown
)

// Cached information about the compiler used.
type Compiler struct {
	Path          string
	Identifier    string // e.g. 'tinygo' or 'go'
	Type          CompilerType
	Version       string // compiler-tool version, e.g. '0.32.0' for tinygo, go1.22.2 for
	VersionGo     string // version of go: same as 'Version' for standard go
	VersionOutput string // output of calling 'tool version'
	IsInit        bool   // CompilerInit() succeeded
}

// Map the compiler's identifier ("tinygo" or "go") to enum.
func CompilerTypeFromString(name string) CompilerType {
	val, ok := map[string]CompilerType{
		"go":     CompilerGo,
		"tinygo": CompilerTinygo,
	}[name]
	if ok {
		return val
	}
	return CompilerUnkown
}

// Sets the compiler for Build() / BuildDir() functions.
func WithCompiler(p string) Opt {
	return func(c *Environ) {
		c.Compiler.Path = p
		c.Compiler.IsInit = false
	}
}

// GoCmd runs a go command. It is used by, among other things, u-root testing
// for such things as go tool.
func (c Environ) GoCmd(gocmd string, args ...string) *exec.Cmd {
	return c.compilerCmd(gocmd, args...)
}

// Returns a compiler command to be run in the environment.
func (c Environ) compilerCmd(gocmd string, args ...string) *exec.Cmd {
	goBin := c.Compiler.Path
	if "" == goBin {
		goBin = filepath.Join(c.GOROOT, "bin", "go")
	}
	args = append([]string{gocmd}, args...)
	cmd := exec.Command(goBin, args...)
	if c.GBBDEBUG {
		log.Printf("GBB Go invocation: %s %s %#v", c, goBin, args)
	}
	cmd.Dir = c.Dir
	cmd.Env = append(os.Environ(), c.Env()...)
	return cmd
}

// compilerAbs does nothing and returns nil if c.Compiler.Path is not set.
// Otherwise, c.Compiler.Path is set to the absolute path of the compiler
// found by exec.LookPath.
//
// LookPath alone is not enough. It only promises an absolute result for a
// bare name it resolved against PATH; a path containing a slash is tried
// directly and returned verbatim, so "./tinygo" comes back relative with no
// error. compilerCmd runs the compiler with cmd.Dir set, which would resolve
// such a path against the build directory rather than the one the user named
// it in, so make it absolute here.
func (c *Environ) compilerAbs() error {
	if len(c.Compiler.Path) == 0 {
		return nil
	}

	fname, err := exec.LookPath(c.Compiler.Path)
	if err != nil {
		return fmt.Errorf("go-compiler: %w", err)
	}

	fname, err = filepath.Abs(fname)
	if err != nil {
		return fmt.Errorf("go-compiler: %w", err)
	}

	c.Compiler.Path = fname
	return nil
}

var (
	// ErrNoVersionString reports output that held no version line at all.
	ErrNoVersionString = errors.New("no compiler version string in output")

	// ErrVersionSyntax reports a version line that was found but could not
	// be parsed into an identifier and a version.
	ErrVersionSyntax = errors.New("unrecognized compiler version string")
)

// versionLine extracts the version line from `go version` output.
//
// The compiler may print progress to the same stream before the version
// itself, e.g. "go: downloading go1.22.4 (linux/amd64)" when GOTOOLCHAIN
// selects a toolchain that is not installed yet. The version is always the
// last non-empty line.
func versionLine(out string) (string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, raw := range slices.Backward(lines) {
		if line := strings.TrimSpace(raw); line != "" {
			return line, nil
		}
	}
	return "", ErrNoVersionString
}

// versionSyntaxError reports a compiler version line that was found but could
// not be parsed into an identifier and a version.
//
// It wraps ErrVersionSyntax so callers can tell malformed version output from
// output that held no version line at all (ErrNoVersionString).
func versionSyntaxError(v string) error {
	return fmt.Errorf("go-compiler version %q: %w", v, ErrVersionSyntax)
}

// Runs compilerCmd("version") and parse/caches output to c.Compiler.
func (c *Environ) CompilerInit() error {
	if c.Compiler.IsInit {
		return nil
	}

	if err := c.compilerAbs(); err != nil {
		return err
	}

	cmd := c.compilerCmd("version")
	vb, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	v, err := versionLine(string(vb))
	if err != nil {
		return err
	}

	s := strings.Fields(v)

	compiler := c.Compiler
	compiler.VersionOutput = strings.TrimSpace(v)
	compiler.Identifier = s[0]
	compiler.Type = CompilerTypeFromString(compiler.Identifier)
	compiler.IsInit = true

	switch compiler.Type {

	case CompilerGo:
		if len(s) < 3 {
			return versionSyntaxError(v)
		}
		compiler.Version = s[2]
		compiler.VersionGo = s[2]

	case CompilerTinygo:
		// e.g. "tinygo version 0.33.0 darwin/arm64 (using go version go1.22.2 and LLVM version 18.1.2)"
		if len(s) < 8 {
			return versionSyntaxError(v)
		}
		compiler.Version = s[2]
		compiler.VersionGo = s[7]

		// Fetch additional go-build-tags from tinygo
		// package fetch needs correct tags to prune
		cmd := c.compilerCmd("info", "-json")
		infov, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		var info map[string]interface{}
		err = json.Unmarshal(infov, &info)
		if err != nil {
			return err
		}

		// extract unique build tags
		tags := make(map[string]struct{})
		for _, tag := range c.BuildTags {
			tags[tag] = struct{}{}
		}
		for _, tag := range info["build_tags"].([]interface{}) {
			tags[tag.(string)] = struct{}{}
		}
		for tag := range tags {
			c.BuildTags = append(c.BuildTags, tag)
		}

	case CompilerUnkown:
		return versionSyntaxError(v)
	}
	c.Compiler = compiler
	return nil
}

// Returns the Go version string that runtime.Version would return for the Go
// compiler in this environ.
func (c *Environ) Version() (string, error) {
	if err := c.CompilerInit(); err != nil {
		return "", err
	}
	return c.Compiler.VersionGo, nil
}

func (c Environ) build(dirPath string, binaryPath string, pattern []string, opts *BuildOpts) error {
	if err := c.CompilerInit(); err != nil {
		return err
	}

	args := []string{
		"-o", binaryPath,
	}

	if c.GO111MODULE != "off" && len(c.Mod) > 0 {
		args = append(args, "-mod", string(c.Mod))
	}
	if c.InstallSuffix != "" {
		args = append(args, "-installsuffix", c.InstallSuffix)
	}

	switch c.Compiler.Type {
	case CompilerGo:

		// Force rebuilding of packages.
		args = append(args, "-a")

		if opts == nil || !opts.EnableInlining {
			// Disable "function inlining" to get a (likely) smaller binary.
			args = append(args, "-gcflags=all=-l")
		}

		if opts == nil || !opts.NoStrip {
			// Strip all symbols, and don't embed a Go build ID to be reproducible.
			args = append(args, "-ldflags", "-s -w -buildid=")
		}

		if opts == nil || !opts.NoTrimPath {
			// Reproducible builds: Trim any GOPATHs out of the executable's
			// debugging information.
			//
			// E.g. Trim /tmp/bb-*/ from /tmp/bb-12345567/src/github.com/...
			args = append(args, "-trimpath")
		}

	case CompilerTinygo:

		// TODO: handle force-rebuild of packages (-a to standard go)
		// TODO: handle EnableInlining

		// Strip all symbols. TODO: not sure about buildid
		if opts == nil || !opts.NoStrip {
			// Strip all symbols
			args = append(args, "-no-debug")
		}

		// TODO: handle NoTrimpPath

	}

	if len(c.BuildTags) > 0 {
		args = append(args, fmt.Sprintf("-tags=%s", strings.Join(c.BuildTags, ",")))
	}

	if opts != nil {
		args = append(args, opts.ExtraArgs...)
	}

	args = append(args, pattern...)

	cmd := c.compilerCmd("build", args...)
	if dirPath != "" {
		cmd.Dir = dirPath
	}

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package in %q: %v, %v", dirPath, string(o), err)
	}

	return nil
}
