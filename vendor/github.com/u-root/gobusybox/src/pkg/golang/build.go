// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package golang is an API to the Go compiler.
package golang

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/uflag"
)

// Environ are the environment variables for the Go compiler.
type Environ struct {
	build.Context

	GO111MODULE string
	GBBDEBUG    bool
}

// RegisterFlags registers flags for Environ.
func (c *Environ) RegisterFlags(f *flag.FlagSet) {
	arg := (*uflag.Strings)(&c.BuildTags)
	f.Var(arg, "go-build-tags", "Go build tags")
}

// Valid returns an error if GOARCH, GOROOT, or GOOS are unset.
func (c Environ) Valid() error {
	if c.GOARCH == "" && c.GOROOT == "" && c.GOOS == "" {
		return fmt.Errorf("golang.Environ should use golang.Default(), not empty value")
	}
	if c.GOARCH == "" {
		return fmt.Errorf("empty GOARCH")
	}
	if c.GOROOT == "" {
		return fmt.Errorf("empty GOROOT")
	}
	if c.GOOS == "" {
		return fmt.Errorf("empty GOOS")
	}
	return nil
}

func parseBool(s string) bool {
	ok, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return ok
}

// Opt is an option function applied to Environ.
type Opt func(*Environ)

// DisableCGO is an option that disables cgo.
func DisableCGO() Opt {
	return func(c *Environ) {
		c.CgoEnabled = false
	}
}

// WithGOARCH is an option that overrides GOARCH.
func WithGOARCH(goarch string) Opt {
	return func(c *Environ) {
		c.GOARCH = goarch
	}
}

// WithGOPATH is an option that overrides GOPATH.
func WithGOPATH(gopath string) Opt {
	return func(c *Environ) {
		c.GOPATH = gopath
	}
}

// WithGOROOT is an option that overrides GOROOT.
func WithGOROOT(goroot string) Opt {
	return func(c *Environ) {
		c.GOROOT = goroot
	}
}

// WithGO111MODULE is an option that overrides GO111MODULE.
func WithGO111MODULE(go111module string) Opt {
	return func(c *Environ) {
		c.GO111MODULE = go111module
	}
}

// Default is the default build environment comprised of the default GOPATH,
// GOROOT, GOOS, GOARCH, and CGO_ENABLED values.
func Default(opt ...Opt) *Environ {
	env := &Environ{
		Context:     build.Default,
		GO111MODULE: os.Getenv("GO111MODULE"),
		GBBDEBUG:    parseBool(os.Getenv("GBBDEBUG")),
	}
	for _, o := range opt {
		o(env)
	}
	return env
}

// GoCmd runs a go command in the environment.
func (c Environ) GoCmd(args ...string) *exec.Cmd {
	goBin := filepath.Join(c.GOROOT, "bin", "go")
	cmd := exec.Command(goBin, args...)
	if c.GBBDEBUG {
		log.Printf("GBB Go invocation: %s %s %#v", c, goBin, args)
	}
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

func (c Environ) envCommon() []string {
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
	}
	return env
}

func (c Environ) EnvHuman() []string {
	env := c.envCommon()
	if c.GOROOT != "" {
		env = append(env, fmt.Sprintf("PATH=%s:$PATH", filepath.Join(c.GOROOT, "bin")))
	}
	return env
}

// Env returns all environment variables for invoking a Go command.
func (c Environ) Env() []string {
	env := c.envCommon()
	if c.GOROOT != "" {
		// If GOROOT is set to a different version of Go, we must
		// ensure that $GOROOT/bin is also in path to make the "go"
		// binary available to golang.org/x/tools/packages.
		env = append(env, fmt.Sprintf("PATH=%s:%s", filepath.Join(c.GOROOT, "bin"), os.Getenv("PATH")))
	}
	return env
}

// String returns all environment variables for Go invocations.
func (c Environ) String() string {
	return strings.Join(c.EnvHuman(), " ")
}

// Optional arguments to Environ.BuildDir.
type BuildOpts struct {
	// NoStrip builds an unstripped binary.
	//
	// Symbols and Build ID will be left in the binary.
	//
	// If NoTrimPath and NoStrip are false, the binary produced will be
	// reproducible.
	NoStrip bool

	// EnableInlining enables function inlining.
	EnableInlining bool

	// NoTrimPath produces a binary whose stack traces contain the module
	// root dirs, GOPATHs, and GOROOTs.
	//
	// If NoTrimPath and NoStrip are false, the binary produced will be
	// reproducible.
	NoTrimPath bool

	// ExtraArgs to `go build`.
	ExtraArgs []string
}

// RegisterFlags registers flags for BuildOpts.
func (b *BuildOpts) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&b.NoStrip, "go-no-strip", false, "Do not strip symbols & Build ID from the binary (will not produce a reproducible binary)")
	f.BoolVar(&b.EnableInlining, "go-enable-inlining", false, "Enable inlining (will likely produce a larger binary)")
	f.BoolVar(&b.NoTrimPath, "go-no-trimpath", false, "Disable -trimpath (will not produce a reproducible binary)")
	arg := (*uflag.Strings)(&b.ExtraArgs)
	f.Var(arg, "go-extra-args", "Extra args to 'go build'")
}

// BuildDir compiles the package in the directory `dirPath`, writing the build
// object to `binaryPath`.
func (c Environ) BuildDir(dirPath string, binaryPath string, opts *BuildOpts) error {
	args := []string{
		"build",

		// Force rebuilding of packages.
		"-a",

		"-o", binaryPath,
	}
	if c.InstallSuffix != "" {
		args = append(args, "-installsuffix", c.Context.InstallSuffix)
	}
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
	if opts != nil {
		args = append(args, opts.ExtraArgs...)
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
