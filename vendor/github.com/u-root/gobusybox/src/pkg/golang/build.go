// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package golang is an API to the Go compiler.
package golang

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/uflag"
	"golang.org/x/tools/go/packages"
)

type ModBehavior string

const (
	ModNone     ModBehavior = ""
	ModOnly     ModBehavior = "mod"
	ModReadonly ModBehavior = "readonly"
	ModVendor   ModBehavior = "vendor"
)

// Environ are the environment variables for the Go compiler.
type Environ struct {
	build.Context

	GO111MODULE string
	Mod         ModBehavior
	GBBDEBUG    bool

	Compiler Compiler
}

// Copy makes a copy of Environ with the given changes.
func (c *Environ) Copy(opts ...Opt) *Environ {
	e := &Environ{
		Context:     c.Context,
		GO111MODULE: c.GO111MODULE,
		Mod:         c.Mod,
		GBBDEBUG:    c.GBBDEBUG,
		Compiler:    c.Compiler,
	}
	e.Apply(opts...)
	return e
}

// RegisterFlags registers flags for Environ.
func (c *Environ) RegisterFlags(f *flag.FlagSet) {
	f.Var((*uflag.Strings)(&c.BuildTags), "go-build-tags", "Go build tags")
	f.StringVar((*string)(&c.Mod), "go-mod", string(c.Mod), "Value of -mod to go commands (allowed: (empty), vendor, mod, readonly)")
	f.StringVar((*string)(&c.Compiler.Path), "compiler", "",
		"override go compiler used (e.g. \"/path/to/tinygo\")")
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
	if goarch == "" {
		return nil
	}
	return func(c *Environ) {
		c.GOARCH = goarch
	}
}

// WithGOOS is an option that overrides GOOS.
func WithGOOS(goos string) Opt {
	if goos == "" {
		return nil
	}
	return func(c *Environ) {
		c.GOOS = goos
	}
}

// WithBuildTag is an option that appends build tags.
func WithBuildTag(tag ...string) Opt {
	return func(c *Environ) {
		c.BuildTags = append(c.BuildTags, tag...)
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

// WithMod is an option that overrides module behavior.
func WithMod(mod ModBehavior) Opt {
	return func(c *Environ) {
		c.Mod = mod
	}
}

// WithWorkingDir sets the working directory for calls to `go`.
func WithWorkingDir(wd string) Opt {
	return func(c *Environ) {
		c.Dir = wd
	}
}

// Default is the default build environment comprised of the default GOPATH,
// GOROOT, GOOS, GOARCH, and CGO_ENABLED values.
func Default(opts ...Opt) *Environ {
	env := &Environ{
		Context:     build.Default,
		GO111MODULE: os.Getenv("GO111MODULE"),
		GBBDEBUG:    parseBool(os.Getenv("GBBDEBUG")),
	}
	env.Apply(opts...)
	return env
}

// Apply applies additional opts to the environment.
func (c *Environ) Apply(opts ...Opt) {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
}

// Lookup looks up packages by patterns relative to dir, using the Go environment from c.
func (c *Environ) Lookup(mode packages.LoadMode, patterns ...string) ([]*packages.Package, error) {

	// required to compute compiler build tags
	if err := c.CompilerInit(); err != nil {
		return nil, err
	}

	cfg := &packages.Config{
		Mode: mode,
		Env:  append(os.Environ(), c.Env()...),
		Dir:  c.Dir,
	}
	if len(c.Context.BuildTags) > 0 {
		tags := fmt.Sprintf("-tags=%s", strings.Join(c.Context.BuildTags, ","))
		cfg.BuildFlags = []string{tags}
	}
	if c.GO111MODULE != "off" && len(c.Mod) > 0 {
		cfg.BuildFlags = append(cfg.BuildFlags, "-mod", string(c.Mod))
	}
	return packages.Load(cfg, patterns...)
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
		env = append(env, fmt.Sprintf("PATH=%s", filepath.Join(c.GOROOT, "bin")))
	}
	return env
}

// Env returns all environment variables for invoking a Go command.
func (c Environ) Env() []string {
	env := c.envCommon()
	if c.GOROOT != "" {
		// If GOROOT is set to a different version of Go, we must
		// ensure that $GOROOT/bin is the path to make the "go" binary
		// available to golang.org/x/tools/packages.
		env = append(env, fmt.Sprintf("PATH=%s", filepath.Join(c.GOROOT, "bin")))
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
	f.Var((*uflag.Strings)(&b.ExtraArgs), "go-extra-args", "Extra args to 'go build'")
}

// BuildDir compiles the package in the directory `dirPath`, writing the build
// object to `binaryPath`.
func (c Environ) BuildDir(dirPath string, binaryPath string, opts *BuildOpts) error {
	return c.build(dirPath, binaryPath, []string{"."}, opts)
}

// Build compiles the pattern.
func (c Environ) Build(binaryPath string, pattern []string, opts *BuildOpts) error {
	return c.build("", binaryPath, pattern, opts)
}
