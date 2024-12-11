// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bb builds one busybox-like binary out of many Go command sources.
//
// This allows you to take two Go commands, such as Go implementations of `sl`
// and `cowsay` and compile them into one binary, callable like `./bb sl` and
// `./bb cowsay`. Which command is invoked is determined by `argv[0]` or
// `argv[1]` if `argv[0]` is not recognized.
//
// Under the hood, bb implements a Go source-to-source transformation on pure
// Go code. This AST transformation does the following:
//
//   - Takes a Go command's source files and rewrites them into Go package files
//     without global side effects.
//   - Writes a `main.go` file with a `main()` that calls into the appropriate Go
//     command package based on `argv[0]`.
//
// Principally, the AST transformation moves all global side-effects into
// callable package functions. E.g. `main` becomes `registeredMain`, each
// `init` becomes `initN`, and global variable assignments are moved into their
// own `initN`.
package bb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/u-root/gobusybox/src/pkg/bb/bbinternal"
	"github.com/u-root/gobusybox/src/pkg/bb/findpkg"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/ulog"

	_ "embed"
)

//go:embed bbmain/cmd/main.go
var bbMainSource []byte

//go:embed bbmain/register.go
var bbRegisterSource []byte

func checkDuplicate(cmds []*bbinternal.Package) error {
	seen := make(map[string]string)
	for _, cmd := range cmds {
		if path, ok := seen[cmd.Name]; ok {
			return fmt.Errorf("failed to build with bb: found duplicate command %s (%s and %s)", cmd.Name, path, cmd.Pkg.PkgPath)
		}
		seen[cmd.Name] = cmd.Pkg.PkgPath
	}
	return nil
}

// Opts are the arguments to BuildBusybox.
type Opts struct {
	// Env are the environment variables used in Go compilation and package
	// discovery.
	Env *golang.Environ

	// LookupEnv is the environment for looking up and resolving command
	// paths.
	//
	// If left unset, DefaultEnv will be used.
	LookupEnv *findpkg.Env

	// GenSrcDir is an empty directory to generate the busybox source code
	// in.
	//
	// If GenSrcDir has children, BuildBusybox will return an error. If
	// GenSrcDir does not exist, it will be created. If no GenSrcDir is
	// given, a temporary directory will be generated. The generated
	// directory will be deleted if compilation succeeds.
	//
	// In GOPATH mode, GOPATH=GenSrcDir for compilation.
	GenSrcDir string

	// CommandPaths is a list of file system directories containing Go
	// commands, or Go import paths.
	CommandPaths []string

	// BinaryPath is the file to write the binary to.
	BinaryPath string

	// GoBuildOpts is configuration for the `go build` command that
	// compiles the busybox binary.
	GoBuildOpts *golang.BuildOpts

	// Generate the tree but don't build it. This is useful for systems
	// like Tamago which have their own way of building.
	GenerateOnly bool
}

// BuildBusybox builds a busybox of many Go commands. opts contains both the
// commands to build and other options.
//
// For documentation on how this works, please refer to the README at the top
// of the repository.
func BuildBusybox(l ulog.Logger, opts *Opts) (nerr error) {
	if opts == nil {
		return fmt.Errorf("no options given for busybox build")
	} else if opts.Env == nil {
		return fmt.Errorf("Go build environment unspecified for busybox build")
	} else if err := opts.Env.Valid(); err != nil {
		return err
	}

	var tmpDir string
	if opts.GenSrcDir != "" {
		var relTmpDir string
		dirents, err := ioutil.ReadDir(opts.GenSrcDir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(opts.GenSrcDir, 0700); err != nil {
				return fmt.Errorf("could not create directory for busybox generated source: %w", err)
			}
			relTmpDir = opts.GenSrcDir
		} else if err != nil {
			return fmt.Errorf("could not read directory supplied for busybox generated source: %w", err)
		} else if len(dirents) > 0 {
			return fmt.Errorf("directory supplied for busybox generated source is not an empty directory")
		} else {
			relTmpDir = opts.GenSrcDir
		}
		absDir, err := filepath.Abs(relTmpDir)
		if err != nil {
			return fmt.Errorf("busybox gen src dir %s could not be made absolute: %v", relTmpDir, err)
		}
		tmpDir = absDir
	} else {
		if opts.GenerateOnly {
			return fmt.Errorf("GenerateOnly switch requires that the GenSrcDir directory be supplied")
		}
		var err error
		tmpDir, err = ioutil.TempDir("", "bb-")
		if err != nil {
			return err
		}
		defer func() {
			if nerr != nil {
				l.Printf("Preserving bb generated source directory at %s due to error", tmpDir)
			} else {
				os.RemoveAll(tmpDir)
			}
		}()
	}

	bbDir := filepath.Join(tmpDir, "src/bb.u-root.com/bb")
	if err := os.MkdirAll(bbDir, 0700); err != nil {
		return err
	}
	pkgDir := filepath.Join(tmpDir, "src")

	var lookupEnv findpkg.Env
	if opts.LookupEnv != nil {
		lookupEnv = *opts.LookupEnv
	} else {
		lookupEnv = findpkg.DefaultEnv()
	}

	// Ask go about all the commands in one batch for dependency caching.
	cmds, err := findpkg.NewPackages(l, opts.Env, lookupEnv, opts.CommandPaths...)
	if err != nil {
		return fmt.Errorf("finding packages failed: %v", err)
	}
	if len(cmds) == 0 {
		return fmt.Errorf("no valid commands given")
	}

	// Collect all packages that we need to actually re-write.
	if err := checkDuplicate(cmds); err != nil {
		return err
	}

	modules := make(map[string]struct{})
	var numNoModule int
	for _, cmd := range cmds {
		if cmd.Pkg.Module != nil {
			modules[cmd.Pkg.Module.Path] = struct{}{}
		} else {
			numNoModule++
		}
	}
	if len(modules) > 0 && numNoModule > 0 {
		return fmt.Errorf("gobusybox does not support mixed module/non-module compilation -- commands contain main modules %v", strings.Join(maps.Keys(modules), ", "))
	}

	// List of packages to import in the real main file.
	var bbImports []string
	// Rewrite commands to packages.
	for _, cmd := range cmds {
		destination := filepath.Join(pkgDir, cmd.Pkg.PkgPath)

		if err := cmd.Rewrite(destination, "bb.u-root.com/bb/pkg/bbmain"); err != nil {
			return fmt.Errorf("rewriting command %q failed: %v", cmd.Pkg.PkgPath, err)
		}
		bbImports = append(bbImports, cmd.Pkg.PkgPath)
	}

	// Collect and write dependencies into pkgDir.
	if err := copyAllDeps(l, opts.Env, bbDir, tmpDir, pkgDir, cmds); err != nil {
		return fmt.Errorf("collecting and putting dependencies in place failed: %v", err)
	}

	if err := writeBBMain(bbDir, tmpDir, bbImports); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}

	if opts.GenerateOnly {
		return nil
	}

	// Get ready to compile bb.
	buildEnv := opts.Env.Copy(golang.WithGO111MODULE("off"), golang.WithGOPATH(tmpDir), golang.WithMod(""))
	if err := buildEnv.BuildDir(bbDir, opts.BinaryPath, opts.GoBuildOpts); err != nil {
		return &ErrBuild{
			CmdDir: bbDir,
			GOPATH: tmpDir,
			Err:    err,
		}
	}
	return nil
}

// ErrBuild is returned for a go build failure when modules were disabled.
type ErrBuild struct {
	CmdDir string
	GOPATH string
	Err    error
}

// Unwrap implements error.Unwrap.
func (e *ErrBuild) Unwrap() error {
	return e.Err
}

// Error implements error.Error.
func (e *ErrBuild) Error() string {
	return fmt.Sprintf("`(cd %s && GOPATH=%s GO111MODULE=off go build)` failed: %v", e.CmdDir, e.GOPATH, e.Err)
}

// writeBBMain writes $TMPDIR/src/bb.u-root.com/bb/pkg/bbmain/register.go and
// $TMPDIR/src/bb.u-root.com/bb/main.go.
//
// They are taken from ./bbmain/register.go and ./bbmain/cmd/main.go, but they
// do not retain their original import paths because the main command must be
// in a module that doesn't conflict with any bb commands. If one were to
// compile github.com/u-root/gobusybox/src/cmd/* into a busybox, we'd have
// problems -- the src/go.mod would conflict with our generated go.mod, and
// it'd be complicated to merge them. So they are transplanted into the
// bb.u-root.com/bb module.
func writeBBMain(bbDir, tmpDir string, bbImports []string) error {
	if err := os.MkdirAll(filepath.Join(bbDir, "pkg/bbmain"), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(bbDir, "pkg/bbmain/register.go"), bbRegisterSource, 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(bbDir, "main.go"), bbMainSource, 0755); err != nil {
		return err
	}

	bbFset, bbFiles, _, err := bbinternal.ParseAST("main", []string{filepath.Join(bbDir, "main.go")})
	if err != nil {
		return err
	}
	if len(bbFiles) == 0 {
		return fmt.Errorf("bb package not found")
	}

	// Fix the import path for bbmain, since we wrote bbmain/register.go into bbDir above.
	if !astutil.RewriteImport(bbFset, bbFiles[0], "github.com/u-root/gobusybox/src/pkg/bb/bbmain", "bb.u-root.com/bb/pkg/bbmain") {
		return fmt.Errorf("could not rewrite import")
	}

	// Create bb main.go.
	if err := bbinternal.CreateBBMainSource(bbFset, bbFiles, bbImports, bbDir); err != nil {
		return fmt.Errorf("creating bb main.go file failed: %v", err)
	}
	return nil
}

func copyAllDeps(l ulog.Logger, env *golang.Environ, bbDir, tmpDir, pkgDir string, mainPkgs []*bbinternal.Package) error {
	var deps []*packages.Package
	for _, p := range mainPkgs {
		deps = append(deps, collectDeps(p.Pkg)...)
	}

	// Copy local dependency packages into module directories at
	// tmpDir/src.
	seenIDs := make(map[string]struct{})
	for _, p := range deps {
		if _, ok := seenIDs[p.ID]; !ok {
			if err := bbinternal.WritePkg(p, filepath.Join(pkgDir, p.PkgPath)); err != nil {
				return fmt.Errorf("writing package %s failed: %v", p, err)
			}
			seenIDs[p.ID] = struct{}{}
		}
	}
	return nil
}

// deps recursively iterates through imports and returns the set of packages
// for which filter returns true.
func deps(p *packages.Package, filter func(p *packages.Package) bool) []*packages.Package {
	var pkgs []*packages.Package
	packages.Visit([]*packages.Package{p}, nil, func(pkg *packages.Package) {
		if filter(pkg) {
			pkgs = append(pkgs, pkg)
		}
	})
	return pkgs
}

func collectDeps(p *packages.Package) []*packages.Package {
	// If modules are not enabled, we need a copy of *ALL*
	// non-standard-library dependencies in the temporary directory.
	return deps(p, func(pkg *packages.Package) bool {
		// First component of package path contains a "."?
		//
		// Poor man's standard library test.
		firstComp := strings.SplitN(pkg.PkgPath, "/", 2)
		return strings.Contains(firstComp[0], ".")
	})
}
