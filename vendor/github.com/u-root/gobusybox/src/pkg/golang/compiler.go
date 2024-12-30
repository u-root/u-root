// Copyright 2015-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package golang is an API to the Go compiler.
package golang

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

// If go-compiler specified, return its absolute-path, otherwise return 'nil'.
func (c *Environ) compilerAbs() error {
	if c.Compiler.Path != "" {
		fname, err := exec.LookPath(string(c.Compiler.Path))
		if err == nil {
			fname, err = filepath.Abs(fname)
		}
		if err != nil {
			return fmt.Errorf("build: %v", err)
		}
		c.Compiler.Path = fname
	}
	return nil
}

// Runs compilerCmd("version") and parse/caches output to c.Compiler.
func (c *Environ) CompilerInit() error {
	if c.Compiler.IsInit {
		return nil
	}

	c.compilerAbs()

	cmd := c.compilerCmd("version")
	vb, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	v := string(vb)

	efmt := "go-compiler 'version' output unrecognized: %v"
	s := strings.Fields(v)
	if len(s) < 1 {
		return fmt.Errorf(efmt, v)
	}

	compiler := c.Compiler
	compiler.VersionOutput = strings.TrimSpace(v)
	compiler.Identifier = s[0]
	compiler.Type = CompilerTypeFromString(compiler.Identifier)
	compiler.IsInit = true

	switch compiler.Type {

	case CompilerGo:
		if len(s) < 3 {
			return fmt.Errorf(efmt, v)
		}
		compiler.Version = s[2]
		compiler.VersionGo = s[2]

	case CompilerTinygo:
		// e.g. "tinygo version 0.33.0 darwin/arm64 (using go version go1.22.2 and LLVM version 18.1.2)"
		if len(s) < 8 {
			return fmt.Errorf(efmt, v)
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
		return fmt.Errorf(efmt, v)
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
		args = append(args, "-installsuffix", c.Context.InstallSuffix)
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
