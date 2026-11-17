// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/tools/go/packages"

	"github.com/u-root/gobusybox/src/pkg/golang"
)

// bbModulePath is the module path of the generated busybox module.
//
// The generated main command must live in a module that cannot collide with
// any of the commands being compiled into the busybox. If one were to compile
// github.com/u-root/gobusybox/src/cmd/* into a busybox, a module named after
// this repository would conflict with src/go.mod, and merging the two would be
// complicated. So the main command and bbmain are transplanted here instead.
//
// Only an exact match is a problem. A vendored package is reached through
// vendor/modules.txt, which maps module paths to vendored directories, so a
// module merely sitting underneath this path -- bb.u-root.com/bb/extra, say,
// or a package bb.u-root.com/bb/cmd/foo belonging to module bb.u-root.com --
// resolves to vendor/ and builds correctly. A module whose path is exactly
// this one is different: it would be required by a go.mod that already
// declares it as the module being built. See checkModulePathCollision.
const bbModulePath = "bb.u-root.com/bb"

// bbMainImportPath is the import path of the generated bbmain package. Every
// rewritten command imports it in order to register itself.
const bbMainImportPath = bbModulePath + "/pkg/bbmain"

// syntheticVersion is the module version recorded for packages whose real
// version is inapplicable: those in a main module, which has no version.
//
// A vendored build never resolves these versions against a proxy or the module
// cache, so any syntactically valid version will do. It only has to agree
// between go.mod and vendor/modules.txt.
const syntheticVersion = "v0.0.0"

// noGoDirectiveVersion is the language version the go tool assumes for a
// module whose go.mod has no go directive.
const noGoDirectiveVersion = "1.16"

// vendoredModule is one module stanza in a generated vendor/modules.txt.
type vendoredModule struct {
	Path      string
	Version   string
	GoVersion string
	Packages  []string
}

// errModulePathCollision is returned when an input module's path is exactly
// bbModulePath, so that the generated go.mod would have to require the very
// module it declares. The go tool rejects that with a message about the
// vendored module needing to be required explicitly, which says nothing about
// the actual cause.
//
// Unexported: reaching this needs a module named after the sentinel path the
// generated tree uses, which in practice means re-running makebb inside a
// previously generated tree. No caller has reason to branch on it. It exists
// to make the message legible, not to be matched.
var errModulePathCollision = errors.New("input module collides with the generated busybox module path")

// hasPathPrefix reports whether modPath is a path-element prefix of pkgPath.
func hasPathPrefix(modPath, pkgPath string) bool {
	return pkgPath == modPath || strings.HasPrefix(pkgPath, modPath+"/")
}

// syntheticVersionFor returns a version that is legal for modPath.
//
// A module path carrying a major version suffix may only be required at that
// major version, so a flat v0.0.0 is not always acceptable to the go tool.
//
// The suffix takes two forms: "/v2" for most paths, and gopkg.in's ".v2".
// Parsing only the first produced "require gopkg.in/yaml.v2 v0.0.0", which the
// go tool rejects with `version "v0.0.0" invalid: should be v2, not v0`, so
// defer to x/mod rather than matching the string here. Note that gopkg.in's
// ".v1" demands v1, where a bare path is free to use v0.
func syntheticVersionFor(modPath string) string {
	_, pathMajor, ok := module.SplitPathVersion(modPath)
	if !ok || pathMajor == "" {
		return syntheticVersion
	}
	// pathMajor is "/vN" or ".vN"; both carry the "v".
	return strings.TrimLeft(pathMajor, "/.") + ".0.0"
}

// vendorManifest groups packages by the module they will be attributed to in
// vendor/modules.txt.
//
// Every package must belong to a module whose path is a prefix of its import
// path. Callers guarantee this by rejecting GOPATH-mode builds up front; see
// checkModules.
func vendorManifest(pkgs map[string]*packages.Package) ([]vendoredModule, error) {
	type modKey struct{ path, version string }
	mods := make(map[modKey]*vendoredModule)

	importPaths := make([]string, 0, len(pkgs))
	for ip := range pkgs {
		importPaths = append(importPaths, ip)
	}
	sort.Strings(importPaths)

	for _, ip := range importPaths {
		p := pkgs[ip]

		mod := p.Module
		if mod == nil {
			return nil, fmt.Errorf("%w: %s", ErrNoModule, ip)
		}
		if mod.Path == bbModulePath {
			return nil, fmt.Errorf("%w: package %s is in module %s, which is the path of the generated busybox module itself -- the generated go.mod would have to require the module it declares; rename the module", errModulePathCollision, ip, mod.Path)
		}
		if !hasPathPrefix(mod.Path, ip) {
			return nil, fmt.Errorf("package %s is attributed to module %s, which is not a prefix of its import path", ip, mod.Path)
		}

		key := modKey{path: mod.Path, version: mod.Version}
		goVer := mod.GoVersion
		if goVer == "" {
			goVer = noGoDirectiveVersion
		}
		// A main module carries no version of its own.
		if key.version == "" {
			key.version = syntheticVersionFor(key.path)
		}

		m, ok := mods[key]
		if !ok {
			m = &vendoredModule{Path: key.path, Version: key.version, GoVersion: goVer}
			mods[key] = m
		}
		m.Packages = append(m.Packages, ip)
	}

	out := make([]vendoredModule, 0, len(mods))
	for _, m := range mods {
		sort.Strings(m.Packages)
		out = append(out, *m)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Version < out[j].Version
	})
	return out, nil
}

// writeModule writes the go.mod and vendor/modules.txt that make the generated
// tree a self-contained, vendored Go module.
//
// The two files are generated together from the same package list because the
// go tool rejects a vendor directory whose modules.txt disagrees with go.mod.
func writeModule(dir string, pkgs map[string]*packages.Package, goVersion string) error {
	mods, err := vendorManifest(pkgs)
	if err != nil {
		return err
	}

	// The generated module must declare a go version at least as high as
	// every module it vendors, the way `go mod tidy` would. Otherwise the
	// go tool refuses the build with "X in vendor/modules.txt requires go
	// >= Y", and GOTOOLCHAIN never gets the chance to select a toolchain
	// new enough to satisfy the dependency.
	modGoVersion := goVersion
	for _, m := range mods {
		if golang.CompareGoVersions(m.GoVersion, modGoVersion) > 0 {
			modGoVersion = m.GoVersion
		}
	}

	var gomod strings.Builder
	fmt.Fprintf(&gomod, "module %s\n\ngo %s\n", bbModulePath, modGoVersion)
	if len(mods) > 0 {
		gomod.WriteString("\nrequire (\n")
		for _, m := range mods {
			fmt.Fprintf(&gomod, "\t%s %s\n", m.Path, m.Version)
		}
		gomod.WriteString(")\n")
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod.String()), 0o644); err != nil {
		return fmt.Errorf("writing generated go.mod: %w", err)
	}

	var modulesTxt strings.Builder
	for _, m := range mods {
		fmt.Fprintf(&modulesTxt, "# %s %s\n", m.Path, m.Version)
		fmt.Fprintf(&modulesTxt, "## explicit; go %s\n", m.GoVersion)
		for _, p := range m.Packages {
			fmt.Fprintf(&modulesTxt, "%s\n", p)
		}
	}
	path := filepath.Join(dir, "vendor", "modules.txt")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(modulesTxt.String()), 0o644); err != nil {
		return fmt.Errorf("writing generated vendor/modules.txt: %w", err)
	}
	return nil
}

// goLangVersion returns the major.minor language version of a Go compiler
// version, e.g. "1.22" for go1.22.4.
//
// This is the floor for the generated module's go directive. Using the
// compiler's own version keeps the directive from ever exceeding the toolchain
// that is about to build the tree.
func goLangVersion(version string) string {
	return golang.MajorMinorGoVersion(version, noGoDirectiveVersion)
}
