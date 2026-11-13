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
// `init` becomes a numbered `busyboxInitN`, and global variable assignments
// are moved into their own `busyboxInitN`.
//
// Opts.LazyInit extends that to named library packages, which keep their own
// name and import path and gain an exported hook that the commands importing
// them call. See docs/adr/0004 and the README's Library Rewrite
// section.
package bb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

// ErrNoModule is returned when a command does not belong to a Go module.
var ErrNoModule = errors.New("command is not in a Go module")

// gopathRemedy is appended to every ErrNoModule message. Refusing a build that
// used to work is only defensible if the refusal says what to do instead.
const gopathRemedy = "GOPATH mode is no longer supported -- put each command in a module, and use `goanywhere` if they span several"

// checkModules rejects commands that do not belong to a Go module.
//
// GOPATH mode is not merely unfashionable here, it is unrepresentable. It lets
// one import path exist at two versions at once, each reached through a
// different nested vendor directory, and the generated module vendors every
// dependency into a single flat directory that can hold only one of them.
// Collapsing the two would silently hand some command a version it was never
// compiled against.
//
// In module mode the go tool has already chosen exactly one version per import
// path -- by minimal version selection within a module, or across a workspace
// -- before gobusybox sees anything, so the flat directory can always hold the
// whole build.
func checkModules(cmds []*bbinternal.Package) error {
	var noModule []string
	for _, cmd := range cmds {
		if cmd.Pkg.Module == nil {
			noModule = append(noModule, cmd.Pkg.PkgPath)
		}
	}
	if len(noModule) == 0 {
		return nil
	}
	sort.Strings(noModule)
	return fmt.Errorf("%w: %s; %s", ErrNoModule, strings.Join(noModule, ", "), gopathRemedy)
}

// checkDuplicate rejects commands whose names collide.
//
// A command is reached in the busybox by the base name of its directory, so
// two packages sharing a base name cannot both be present however far apart
// their import paths are.
//
// Every collision is reported rather than the first, because the patterns that
// cause this match hundreds of packages and collide on a handful. Stopping at
// the first turns one fix into a sequence of rebuilds, each of which pays for
// a full load to reveal the next.
func checkDuplicate(cmds []*bbinternal.Package) error {
	paths := make(map[string][]string)
	for _, cmd := range cmds {
		paths[cmd.Name] = append(paths[cmd.Name], cmd.Pkg.PkgPath)
	}

	var names []string
	for name, p := range paths {
		if len(p) > 1 {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)

	var b strings.Builder
	fmt.Fprintf(&b, "failed to build with bb: %d duplicate command name", len(names))
	if len(names) > 1 {
		b.WriteString("s")
	}
	for _, name := range names {
		p := paths[name]
		sort.Strings(p)
		fmt.Fprintf(&b, "\n\t%s: %s", name, strings.Join(p, ", "))
	}
	b.WriteString("\n\nName the ones to leave out as exclusions, e.g. -")
	b.WriteString(paths[names[0]][0])
	return errors.New(b.String())
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
	// GenSrcDir becomes the root of a self-contained, fully vendored Go
	// module. It can be built on its own with
	// `cd $GenSrcDir && GOWORK=off go build -mod=vendor`.
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

	// LazyInit selects the library packages whose initialization is
	// deferred to the commands that import them. Each entry is a whole
	// import path, or a path holding "..." to reach a subtree.
	//
	// Empty means every package initializes at program start, as it always
	// has, and nothing below the command level is rewritten. A selected
	// package gives up its package-level variables and its init together,
	// and so does every package that imports it: an eager package
	// initializes before any deferred one and so may not import one.
	//
	// An exact path that cannot be deferred fails the build. A pattern
	// defers what it can and reports the rest, failing only when it moved
	// nothing. Either way, asking for laziness and silently not getting it
	// is the failure this is hardest to notice.
	LazyInit []string
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
	// Reject GOPATH mode before resolving anything. Package resolution
	// fails in its own confusing ways under GO111MODULE=off, and the
	// reason gobusybox cannot proceed has nothing to do with whichever
	// path happens to fail first.
	if opts.Env.GO111MODULE == "off" {
		return fmt.Errorf("%w: GO111MODULE=off; %s", ErrNoModule, gopathRemedy)
	}

	var tmpDir string
	if opts.GenSrcDir != "" {
		var relTmpDir string
		dirents, err := os.ReadDir(opts.GenSrcDir)
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
		tmpDir, err = os.MkdirTemp("", "bb-")
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

	// The generated tree is a single, fully vendored Go module.
	//
	// The main command sits at the module root; every rewritten command
	// and every non-standard-library dependency is written into vendor/ at
	// its original import path, so that the rewritten sources continue to
	// compile without touching any of their import statements.
	//
	// Vendoring is what keeps a traditionally offline compilation offline:
	// -mod=vendor consults neither the network nor the module cache.
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0700); err != nil {
		return err
	}

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

	if err := checkModules(cmds); err != nil {
		return err
	}

	// Which packages wait for a command that imports them, and which
	// initialize at program start as they always have.
	plan, err := planInit(l, cmds, opts.LazyInit)
	if err != nil {
		return err
	}

	// Every package written into vendor/, keyed by the import path it is
	// written at. go.mod and vendor/modules.txt are derived from this same
	// map, so the manifest cannot drift from what is on disk.
	vendored := make(map[string]*packages.Package)

	// The hook of every deferred package, by import path. A hook's name is
	// chosen against the package's own scope, so it is not predictable
	// from the path and has to be recorded as each package is rewritten.
	hooks := make(map[string]bbinternal.Hook)

	// Deferred libraries before the commands that call them, so that every
	// hook a command names has a name by the time the command is written.
	if err := rewriteLibraries(vendorDir, cmds, plan, vendored, hooks); err != nil {
		return err
	}

	// List of packages to import in the real main file.
	var bbImports []string
	// Rewrite commands to packages.
	for _, cmd := range cmds {
		ipath := cmd.Pkg.PkgPath
		if err := cmd.Rewrite(filepath.Join(vendorDir, ipath), bbMainImportPath, hookCalls(plan, hooks, ipath)); err != nil {
			return fmt.Errorf("rewriting command %q failed: %v", cmd.Pkg.PkgPath, err)
		}
		vendored[ipath] = cmd.Pkg
		bbImports = append(bbImports, ipath)
	}

	// Collect and write dependencies into vendorDir.
	if err := copyAllDeps(vendorDir, cmds, vendored); err != nil {
		return fmt.Errorf("collecting and putting dependencies in place failed: %v", err)
	}

	if err := writeBBMain(tmpDir, bbImports, len(hooks) > 0); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}

	reportLazyInit(l, plan, hooks)

	version, err := opts.Env.Version()
	if err != nil {
		return fmt.Errorf("could not determine Go version: %w", err)
	}
	if err := writeModule(tmpDir, vendored, goLangVersion(version)); err != nil {
		return err
	}

	if opts.GenerateOnly {
		return nil
	}

	// Get ready to compile bb.
	//
	// GOWORK=off because the generated module must be built on its own
	// terms even if it happens to be written underneath a directory
	// containing a go.work.
	buildEnv := opts.Env.Copy(
		golang.WithGO111MODULE("on"),
		golang.WithGOWORK("off"),
		golang.WithMod(golang.ModVendor),
	)
	if err := buildEnv.BuildDir(tmpDir, opts.BinaryPath, opts.GoBuildOpts); err != nil {
		return &ErrBuild{
			CmdDir: tmpDir,
			Err:    err,
		}
	}
	return nil
}

// ErrBuild is returned when compiling the generated busybox module fails.
type ErrBuild struct {
	CmdDir string
	Err    error
}

// Unwrap implements error.Unwrap.
func (e *ErrBuild) Unwrap() error {
	return e.Err
}

// Error implements error.Error.
func (e *ErrBuild) Error() string {
	return fmt.Sprintf("`(cd %s && GO111MODULE=on GOWORK=off go build -mod=vendor)` failed: %v", e.CmdDir, e.Err)
}

// writeBBMain writes $GENDIR/pkg/bbmain/register.go and $GENDIR/main.go.
//
// They are taken from ./bbmain/register.go and ./bbmain/cmd/main.go, but they
// do not retain their original import paths, because the main command must be
// in a module that doesn't conflict with any bb commands. See bbModulePath.
//
// lazy arms the flag register.go declares but never sets, which is what makes
// it refuse a second dispatched command. That refusal is only correct once
// initialization has been deferred: while every package initializes before
// main, a second command runs exactly the init Go would have run for it.
func writeBBMain(bbDir string, bbImports []string, lazy bool) error {
	if err := os.MkdirAll(filepath.Join(bbDir, "pkg/bbmain"), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(bbDir, "pkg/bbmain/register.go"), bbRegisterSource, 0644); err != nil {
		return err
	}
	if lazy {
		if err := os.WriteFile(filepath.Join(bbDir, "pkg/bbmain/lazyinit.go"), []byte(lazyInitArmSource), 0644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(bbDir, "main.go"), bbMainSource, 0644); err != nil {
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
	if !astutil.RewriteImport(bbFset, bbFiles[0], "github.com/u-root/gobusybox/src/pkg/bb/bbmain", bbMainImportPath) {
		return fmt.Errorf("could not rewrite import")
	}

	// Create bb main.go.
	if err := bbinternal.CreateBBMainSource(bbFset, bbFiles, bbImports, bbDir); err != nil {
		return fmt.Errorf("creating bb main.go file failed: %v", err)
	}
	return nil
}

// copyAllDeps writes every non-standard-library dependency of mainPkgs into
// vendorDir, recording each one in vendored.
//
// Packages already present in vendored -- the rewritten commands -- are left
// alone, so that a command is never overwritten by its own unrewritten source.
func copyAllDeps(vendorDir string, mainPkgs []*bbinternal.Package, vendored map[string]*packages.Package) error {
	var deps []*packages.Package
	for _, p := range mainPkgs {
		deps = append(deps, collectDeps(p.Pkg)...)
	}

	for _, p := range deps {
		ipath := p.PkgPath
		if _, ok := vendored[ipath]; ok {
			continue
		}
		if err := bbinternal.WritePkg(p, filepath.Join(vendorDir, ipath)); err != nil {
			return fmt.Errorf("writing package %s failed: %v", p, err)
		}
		vendored[ipath] = p
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
	// The generated module vendors everything, so we need a copy of *ALL*
	// non-standard-library dependencies in the temporary directory.
	return deps(p, func(pkg *packages.Package) bool {
		// Standard library packages belong to no module, and every
		// other package belongs to one: its own main module, a
		// dependency module, or a workspace member. go/packages fills
		// this in because findpkg asks for packages.NeedModule, so it
		// is the go tool's own answer rather than a guess at one.
		//
		// The guess this replaced was whether the first path component
		// contained a dot. Module paths are not required to contain
		// one -- `module nodot` is legal for anything not fetched by
		// path -- and such a package was taken for standard library,
		// left unvendored, and failed the build it was part of:
		//
		//	vendor/nodot/cmd/hi/main.go:6:2:
		//	package nodot/pkg/hello is not in std
		return pkg.Module != nil
	})
}
