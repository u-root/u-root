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
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/goterm/term"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/u-root/gobusybox/src/pkg/bb/bbinternal"
	"github.com/u-root/gobusybox/src/pkg/bb/findpkg"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/cp"
	"github.com/u-root/uio/ulog"
)

func listStrings(m map[string]struct{}) []string {
	var l []string
	for k := range m {
		l = append(l, k)
	}
	return l
}

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

	// AllowMixedMode allows mixed mode (module / non-module) compilation.
	//
	// If this is done with GO111MODULE=on,
	AllowMixedMode bool

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
	if !opts.AllowMixedMode && len(modules) > 0 && numNoModule > 0 {
		return fmt.Errorf("gobusybox does not support mixed module/non-module compilation -- commands contain main modules %v", strings.Join(listStrings(modules), ", "))
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
	if err := copyLocalDeps(l, opts.Env, bbDir, tmpDir, pkgDir, cmds); err != nil {
		return fmt.Errorf("collecting and putting dependencies in place failed: %v", err)
	}

	if err := writeBBMain(bbDir, tmpDir, bbImports); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}

	// Get ready to compile bb.
	if opts.Env.GO111MODULE == "off" || numNoModule > 0 {
		opts.Env.GOPATH = tmpDir
	} else {
		// Run go mod tidy in order to get the go.sum file
		// sorted. This likely requires people to be online.
		//
		// Sorting a go.sum file seems to be required now in order to
		// get builds to work. Sorting is only necessary when we merge
		// more than one go.sum file (i.e. we are compiling commands
		// from more than one module, e.g. u-root and u-bmc).
		//
		// This can only be done after writeBBMain, as it reads what
		// main.go depends on and prunes everything that isn't needed.
		cmd := opts.Env.GoCmd("mod", "tidy")
		cmd.Dir = bbDir
		if o, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("running `go mod tidy` on the generated busybox main package failed (%v): %s", err, o)
		}
	}

	if opts.GenerateOnly {
		return nil
	}

	// Compile bb.
	if err := opts.Env.BuildDir(bbDir, opts.BinaryPath, opts.GoBuildOpts); err != nil {
		if opts.Env.GO111MODULE == "off" || numNoModule > 0 {
			return &ErrGopathBuild{
				CmdDir: bbDir,
				GOPATH: tmpDir,
				Err:    err,
			}
		} else {
			return &ErrModuleBuild{
				CmdDir: bbDir,
				Err:    err,
			}
		}
	}
	return nil
}

// ErrModuleBuild is returned for a go build failure when modules were enabled.
type ErrModuleBuild struct {
	CmdDir string
	Err    error
}

// Unwrap implements error.Unwrap.
func (e *ErrModuleBuild) Unwrap() error {
	return e.Err
}

// Error implements error.Error.
func (e *ErrModuleBuild) Error() string {
	return fmt.Sprintf("go build with modules failed: %v", e.Err)
}

// ErrGopathBuild is returned for a go build failure when modules were disabled.
type ErrGopathBuild struct {
	CmdDir string
	GOPATH string
	Err    error
}

// Unwrap implements error.Unwrap.
func (e *ErrGopathBuild) Unwrap() error {
	return e.Err
}

// Error implements error.Error.
func (e *ErrGopathBuild) Error() string {
	return fmt.Sprintf("non-module go build failed: %v", e.Err)
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

func isReplacedModuleLocal(m *packages.Module) bool {
	// From "replace directive": https://golang.org/ref/mod#go
	//
	//   If the path on the right side of the arrow is an absolute or
	//   relative path (beginning with ./ or ../), it is interpreted as the
	//   local file path to the replacement module root directory, which
	//   must contain a go.mod file. The replacement version must be
	//   omitted in this case.
	return strings.HasPrefix(m.Path, "./") || strings.HasPrefix(m.Path, "../") || strings.HasPrefix(m.Path, "/")
}

// copyLocalGoMods copies the go.sum and go.mod of all modules that are locally
// present in the file system into the right place in the generated tree.
func copyLocalGoMods(pkgDir, bbDir string, modules map[string]*packages.Module) error {
	copyGoMod := func(mod *packages.Module) error {
		if mod == nil {
			return nil
		}

		if err := os.MkdirAll(filepath.Join(pkgDir, mod.Path), 0755); os.IsExist(err) {
			return nil
		} else if err != nil {
			return err
		}

		// Use the module file for all outside dependencies.
		if err := cp.Copy(mod.GoMod, filepath.Join(pkgDir, mod.Path, "go.mod")); err != nil {
			return err
		}

		// As of Go 1.16, the Go build system expects an accurate
		// go.sum in the main module directory. We build it by
		// concatenating all constituent go.sums.
		//
		// If it doesn't exist, that's okay!
		gosum := filepath.Join(filepath.Dir(mod.GoMod), "go.sum")
		if err := cp.Copy(gosum, filepath.Join(pkgDir, mod.Path, "go.sum")); os.IsNotExist(err) {
			// Modules without dependencies don't have or need a go.sum.
			return nil
		} else if err != nil {
			return err
		}

		gosumf, err := os.Open(gosum)
		if err != nil {
			return err
		}
		defer gosumf.Close()
		f, err := os.OpenFile(filepath.Join(bbDir, "go.sum"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, gosumf)
		return err
	}

	for modPath, mod := range modules {
		if err := copyGoMod(mod); err != nil {
			return fmt.Errorf("failed to copy go.mod for %s: %v", modPath, err)
		}
	}
	return nil
}

// findLocalModules finds all modules that are locally present on the file
// system and raises an error if any modules dependencies conflict with each
// other (e.g. because one module requests a different `replace` dependency
// than another). Conflicts generally only arise when a module is replaced by a
// local directory, as Go takes care of other conflicts using
// minimum-version-selection (MVS).
func findLocalModules(l ulog.Logger, mainPkgs []*bbinternal.Package) (map[string]*packages.Module, error) {
	type localModule struct {
		m          *packages.Module
		provenance string
	}

	localModules := make(map[string]*localModule)

	// These are all modules that the user requested to compile commands
	// from. They are assumed to be local directories.
	for _, p := range mainPkgs {
		if p.Pkg.Module != nil {
			if _, ok := localModules[p.Pkg.Module.Path]; !ok {
				localModules[p.Pkg.Module.Path] = &localModule{
					m:          p.Pkg.Module,
					provenance: p.Pkg.PkgPath,
				}
			}
		}
	}

	// This finds all modules that are `replace`d in go.mod files of the
	// commands.
	for _, p := range mainPkgs {
		replacedModules := locallyReplacedModules(p.Pkg)
		for modPath, module := range replacedModules {
			if original, ok := localModules[modPath]; ok {
				// Is this module different from one that a
				// previous definition provided?
				//
				// This can happen if:
				// - there are 2 modules that have conflicting
				//   replace directives:
				//     replace u-root => ../foo
				//     replace u-root => ../bar
				//
				// - there is a module that has a replace
				//   directive that doesn't match the directory
				//   the user requested commands from, e.g.
				//     ./makebb ~/u-root/cmds/core/ip ~/cpu/cmds/cpu
				//   but cpu/go.mod has
				//     replace u-root => ../foobar (which is not ~/u-root!)
				//
				// TODO: write a pretty log message for the
				// user with a suggestion of what to do.
				if original.m.Dir != module.Dir {
					return nil, fmt.Errorf("two conflicting versions of module %s have been requested; one from %s, the other from %s's go.mod",
						modPath, original.provenance, p.Pkg.Module.Path)
				}
			} else {
				localModules[modPath] = &localModule{
					m:          module,
					provenance: fmt.Sprintf("%s's go.mod (%s)", p.Pkg.Module.Path, p.Pkg.Module.GoMod),
				}
			}
		}
	}

	// Look for versioning conflicts between all modules.
	//
	// Go through the entire dependency graph of every command the user
	// requested. Every dependency has a version through go.mod files, and
	// that version number may conflict with either a replace directive or
	// the fact that the user requested to compile a command from the
	// dependency module.
	//
	// E.g. if u-bmc depends on u-root @ v0.8.0, but we are also compiling
	// u-root from a local directory, those are conflicting requirements.
	var conflict bool

	// seen is a map of user-requested-command-module =>
	// local-dependency-module combination that a warning has already been
	// printed about.
	seen := map[string]map[string]struct{}{}

	for _, mainPkg := range mainPkgs {
		// Initialize the inner seen map.
		if mainPkg.Pkg.Module != nil {
			if _, ok := seen[mainPkg.Pkg.Module.Path]; !ok {
				seen[mainPkg.Pkg.Module.Path] = map[string]struct{}{}
			}
		}

		// Visit visits all packages in the dependency graph of the named package.
		packages.Visit([]*packages.Package{mainPkg.Pkg}, nil, func(p *packages.Package) {
			if p.Module == nil {
				return
			}
			if _, ok := seen[mainPkg.Pkg.Module.Path][p.Module.Path]; ok {
				return
			}

			if lm, ok := localModules[p.Module.Path]; ok && lm.m.Dir != p.Module.Dir {
				gbbstrict, set := os.LookupEnv("GBB_STRICT")
				if set == false {
					l.Printf("GBB_STRICT is not set.")
				}
				if gbbstrict != "1" {
					if p.Module.Version == "" || lm.m.Version == "" || (p.Module.Replace != nil && isReplacedModuleLocal(p.Module.Replace)) || (lm.m.Replace != nil && isReplacedModuleLocal(lm.m.Replace)) {
						seen[mainPkg.Pkg.Module.Path][p.Module.Path] = struct{}{}
						l.Printf("[WARNING] %s depends on %s @ %s\n", mainPkg.Pkg.PkgPath, p.Module.Path, moduleVersionIdentifier(p.Module))
						l.Printf("\tUsing %s @ %s to build it.", lm.m.Path, moduleVersionIdentifier(lm.m))
						return
					}
				}
				fmt.Fprintln(os.Stderr, "")
				l.Printf("Conflicting module dependencies on %s:", p.Module.Path)
				l.Printf("  %s depends on %s @ %s", mainPkg.Pkg.PkgPath, p.Module.Path, moduleVersionIdentifier(p.Module))
				l.Printf("  %s depends on %s @ %s", lm.provenance, lm.m.Path, moduleVersionIdentifier(lm.m))
				replacePath, err := filepath.Rel(mainPkg.Pkg.Module.Dir, lm.m.Dir)
				if err != nil {
					replacePath = lm.m.Dir
				}
				fmt.Fprintln(os.Stderr, "")
				l.Printf("%s: add `replace %s => %s` to %s", term.Bold("Suggestion to resolve"), p.Module.Path, replacePath, mainPkg.Pkg.Module.GoMod)
				fmt.Fprintln(os.Stderr, "")
				conflict = true

				// Don't print this particular warning combo again.
				seen[mainPkg.Pkg.Module.Path][p.Module.Path] = struct{}{}
			}
		})
	}
	if conflict {
		return nil, fmt.Errorf("conflicting module dependencies found")
	}

	modules := make(map[string]*packages.Module)
	for modPath, mod := range localModules {
		modules[modPath] = mod.m
	}
	return modules, nil
}

func moduleVersionIdentifier(m *packages.Module) string {
	// This module was one of the commands requested by the user, and hence was a local directory.
	if m.Version == "" {
		return fmt.Sprintf("directory %s", m.Dir)
	}

	// This module was replaced by some dependency to a local directory.
	if m.Replace != nil && isReplacedModuleLocal(m.Replace) {
		return fmt.Sprintf("directory %s", m.Replace.Path)
	}

	return fmt.Sprintf("version %s", m.Version)
}

// copyLocalDeps tries to suss out local files that need to be in the generated tree.
//
// It copies files from all dependency packages and modules that need to be in
// the generated tree, but NOT for the main commands that are going to be
// rewritten.
//
// It helps to have read https://golang.org/ref/mod when editing this function.
//
// Module-enabled Go programs resolve their dependencies in one of two ways:
//
//   - versioned dependencies: via a version control system at a specific
//     version, potentially remotely downloaded
//
//   - locally: a module that is either `replace`d with a local file system
//     directory, or a command that is being built from a module that is on the
//     local file system (e.g. ./makebb ../u-root/cmds/core/ip -- here, ../u-root
//     will be a local directory module)
//
// Go minimum version selection (MVS) will take care of all versioned
// dependencies on its own.
//
// *We* have to take care of all files that are in local modules: everything
// required for compilation of the gobusybox within local modules has to be
// copied to the generated tree.
//
// For local dependencies, we copy all dependency packages' files over, as well
// as the local modules' go.sum and go.mod files.
//
// Then, in the generated tree's main module, we create a go.mod file with
// replace directives for all the local modules we just copied over.
func copyLocalDeps(l ulog.Logger, env *golang.Environ, bbDir, tmpDir, pkgDir string, mainPkgs []*bbinternal.Package) error {
	localModules, err := findLocalModules(l, mainPkgs)
	if err != nil {
		return err
	}
	// Copy go.sum and go.mod files for all local modules to the generated tree.
	if err := copyLocalGoMods(pkgDir, bbDir, localModules); err != nil {
		return err
	}

	// Find all packages that need to be copied over to the generated tree.
	//
	// This is going to be all the source code that
	//   (a) a requested command depends on (one of mainPkgs), AND
	//   (b) that is in one of the localModules.
	var localDepPkgs []*packages.Package
	for _, p := range mainPkgs {
		// Find all dependency packages that are *within* module boundaries for this package.
		localDeps := collectDeps(p.Pkg, localModules)
		localDepPkgs = append(localDepPkgs, localDeps...)
	}

	// TODO(chrisko): We need to go through mainPkgs Module definitions to
	// find exclude and replace directives, which only have an effect in
	// the main module's go.mod, which will be the top-level go.mod we
	// write.
	//
	// mainPkgs module files expect to be "the main module", since those
	// are where Go compilation would normally occur.
	//
	// The top-level go.mod must have copies of the mainPkgs' modules'
	// replace and exclude directives. If they conflict, we need to have a
	// legible error message for the user.

	// Copy local dependency packages into module directories at
	// tmpDir/src.
	seenIDs := make(map[string]struct{})
	for _, p := range localDepPkgs {
		if _, ok := seenIDs[p.ID]; !ok {
			if err := bbinternal.WritePkg(p, filepath.Join(pkgDir, p.PkgPath)); err != nil {
				return fmt.Errorf("writing package %s failed: %v", p, err)
			}
			seenIDs[p.ID] = struct{}{}
		}
	}

	// Avoid go.mod in the case of GO111MODULE=(auto|off) if there are no modules.
	if env.GO111MODULE == "on" || len(localModules) > 0 {
		// go.mod for the bb binary.
		//
		// Add local replace rules for all modules we're compiling.
		//
		// This is the only way to locally reference another modules'
		// repository. Otherwise, go'll try to go online to get the source.
		//
		// The module name is something that'll never be online, lest Go
		// decides to go on the internet.
		var mod modfile.File

		mod.AddModuleStmt("bb.u-root.com/bb")
		for mpath, module := range localModules {
			v := module.Version
			if len(v) == 0 {
				// When we don't know the version, we can plug
				// in a "Go-generated" version number to get
				// past the validation in the compiler.
				//
				// We don't always do this because if the
				// module path has a /v2 or /v3, Go expects the
				// version number to match. So we use the
				// module.Version when available, because it's
				// the most accurate thing.
				v = "v0.0.0"
			}
			if err := mod.AddRequire(mpath, v); err != nil {
				return fmt.Errorf("could not add requiring %v to go.mod: %v", mpath, err)
			}
			if err := mod.AddReplace(mpath, "", path.Join("..", "..", mpath), ""); err != nil {
				return fmt.Errorf("could not add replace rule for %v to go.mod: %v", mpath, err)
			}
		}

		gomod, err := mod.Format()
		if err != nil {
			return fmt.Errorf("could not generated go.mod: %v", err)
		}

		// TODO(chrisko): add other go.mod files' replace and exclude
		// directives.
		//
		// Warn the user if they are potentially incompatible.
		if err := ioutil.WriteFile(filepath.Join(bbDir, "go.mod"), gomod, 0755); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func versionNum(mpath string) string {
	last := path.Base(mpath)
	if len(last) == 0 {
		return "v0"
	}
	if matched, _ := regexp.Match("v[0-9]+", []byte(last)); matched {
		return last
	}
	return "v0"
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

func locallyReplacedModules(p *packages.Package) map[string]*packages.Module {
	if p.Module == nil {
		return nil
	}

	m := make(map[string]*packages.Module)
	// Collect all "local" dependency packages that are in `replace`
	// directives somewhere, to be copied into the temporary directory
	// structure later.
	packages.Visit([]*packages.Package{p}, nil, func(pkg *packages.Package) {
		if pkg.Module != nil && pkg.Module.Replace != nil && isReplacedModuleLocal(pkg.Module.Replace) {
			m[pkg.Module.Path] = pkg.Module
		}
	})
	return m
}

func collectDeps(p *packages.Package, localModules map[string]*packages.Module) []*packages.Package {
	if p.Module != nil {
		// Collect all "local" dependency packages, to be copied into
		// the temporary directory structure later.
		return deps(p, func(pkg *packages.Package) bool {
			// Replaced modules can be local on the file system.
			if pkg.Module != nil && pkg.Module.Replace != nil && isReplacedModuleLocal(pkg.Module.Replace) {
				return true
			}

			// Is this a dependency within a local module?
			for modulePath := range localModules {
				if strings.HasPrefix(pkg.PkgPath, modulePath) {
					return true
				}
			}
			return false
		})
	}

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
