// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package findpkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/uio/ulog/ulogtest"
)

type testCase struct {
	// name of the test case
	name string
	// envs to try it in (if unset, default will be GO111MODULE=on and off)
	envs []*golang.Environ
	// wd sets the findpkg.Env.WorkingDirectory
	// WorkingDirectory is the directory used for module-enabled
	// `go list` lookups. The go.mod in this directory (or one of
	// its parents) is used to resolve package paths.
	wd string
	// GBB_PATH
	gbbPath []string
	// Input patterns
	in []string
	// Expected result from ResolveGlobs
	want []string
	// or, how many results you want at least.
	atLeast int
	// Expected result from NewPackages each packages' PkgPath
	wantPkgPath []string
	// Error expected?
	wantErr bool
	// If set, expected error.
	err error
}

func TestResolve(t *testing.T) {
	gbbmod, err := filepath.Abs("../../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gbbroot := filepath.Dir(gbbmod)
	cmdRoot := filepath.Join(gbbroot, "cmds/exp/cmd2pkg")
	moduleOffEnv := golang.Default(golang.WithGO111MODULE("on"))
	moduleOnEnv := golang.Default(golang.WithGO111MODULE("on"))
	// TODO: re-enable when https://github.com/golang/go/issues/62114 is resolved.
	// noGoToolEnv := golang.Default(golang.WithGOROOT(t.TempDir()))

	if err := os.Mkdir("./test/resolvebroken", 0777); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll("./test/resolvebroken") })
	if err := os.WriteFile("./test/resolvebroken/main.go", []byte("broken"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir("./test/parsebroken", 0777); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll("./test/parsebroken") })
	if err := os.WriteFile("./test/parsebroken/main.go", []byte("package main\n\nimport \"fmt\""), 0777); err != nil {
		t.Fatal(err)
	}

	l := &ulogtest.Logger{TB: t}

	sharedTestCases := []testCase{
		// Nonexistent Package
		{
			name:    "fakepackage",
			in:      []string{"fakepackagename"},
			wantErr: true,
		},
		// Single package, file system path.
		{
			name: "fspath-single",
			in:   []string{filepath.Join(gbbmod, "exp/cmd2pkg")},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
		},
		// Single package, file system path, GBB_PATHS.
		{
			name:    "fspath-gbbpath-single",
			gbbPath: []string{gbbmod},
			in:      []string{"exp/cmd2pkg"},
			want:    []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
		},
		// Single package, Go package path.
		{
			name: "pkgpath-single",
			in:   []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
		},
		// globbed file system path.
		{
			name: "fspath-glob",
			in:   []string{filepath.Join(gbbmod, "exp/cmd2*")},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
		},
		// globbed Go package path.
		{
			name: "pkgpath-glob",
			in:   []string{"github.com/u-root/u-root/cmds/exp/cmd2*"},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg"},
		},
		// Globbed file system path of non-existent packages.
		{
			name:    "fspath-glob-doesnotexist",
			in:      []string{filepath.Join(gbbmod, "cmd/makeq*")},
			wantErr: true,
			err:     errNoMatch,
		},
		// Globbed package path of non-existent packages.
		{
			name:    "pkgpath-glob-doesnotexist",
			in:      []string{"github.com/u-root/gobusybox/src/cmd/makeq*"},
			wantErr: true,
			err:     errNoMatch,
		},
		// Two packages (file system paths), one excluded by build constraints.
		// The issue is that depending on whether you run on OSX, Plan 9, or Linux, or ...
		// you will get a different number of matches
		{
			name:    "fspath-log-buildconstrained",
			in:      []string{"./test/buildconstraint", filepath.Join(gbbmod, "core/b*")},
			atLeast: 1,
		},
		// Two packages (Go package paths), some excluded by build constraints.
		{
			name:    "pkgpath-log-buildconstrained",
			in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg", "github.com/u-root/u-root/cmds/core/bind", "github.com/u-root/u-root/cmds/core/ip"},
			atLeast: 1,
		},
		{
			name:    "fspath-log-buildconstrained-onlyone",
			in:      []string{"./test/buildconstraint"},
			err:     errNoMatch,
			wantErr: true,
		},
		// Package excluded by build constraints (Go package paths).
		{
			name:    "pkgpath-log-buildconstrained-onlyone",
			in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/buildconstraint"},
			err:     errNoMatch,
			wantErr: true,
		},
		// Go glob support (Go package path).
		{
			name: "pkgpath-go-glob",
			in:   []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/..."},
			want: []string{
				"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/echo",
				"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/foo",
			},
		},
		// Go glob support (relative Go package path).
		{
			name: "pkgpath-relative-go-glob",
			in:   []string{"./test/goglob/..."},
			want: []string{
				"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/echo",
				"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/foo",
			},
		},
		// Go glob support ("relative" Go package path, without ./ -- follows Go semantics).
		//
		// This is actually just a Go package path, and not interpreted
		// as a file system path by the Go lookup (because they must be
		// able to distinguish between "cmd/compile" and
		// "./cmd/compile").
		{
			name:    "pkgpath-relative-go-glob-broken",
			in:      []string{"test/goglob/..."},
			wantErr: true,
			err:     errNoMatch,
		},
		{
			name:    "fspath-empty-directory",
			in:      []string{"./test/empty"},
			wantErr: true,
		},
		{
			name:    "pkgpath-empty-directory",
			in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/empty"},
			wantErr: true,
		},
		// resolvebroken is not compilable.
		{
			name:    "fspath-broken-go",
			in:      []string{"./test/resolvebroken"},
			wantErr: true,
		},
		{
			name:    "pkgpath-broken-go",
			in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/resolvebroken"},
			wantErr: true,
		},
		// Contains test/resolvebroken which is not compilable.
		{
			name:    "fspath-glob-with-errors",
			in:      []string{"./test/*"},
			wantErr: true,
		},
		{
			name:    "pkgpath-glob-with-errors",
			in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/*"},
			wantErr: true,
		},
		// Multi module resolution, package path. (GO111MODULE=on only)
		//
		// Unless we put u-root and p9 in GOPATH in the local version
		// of this test, this is an ON only test.
		{
			name: "pkgpath-multi-module",
			envs: []*golang.Environ{moduleOnEnv},
			wd:   filepath.Join(cmdRoot, "test/resolve-modules"),
			in: []string{
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/dhclient",
				"github.com/hugelgupf/p9/cmd/p9ufs",
			},
			want: []string{
				"github.com/hugelgupf/p9/cmd/p9ufs",
				"github.com/u-root/u-root/cmds/core/dhclient",
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/ip",
			},
		},

		// Shell expansions.
		{
			name: "pkgpath-shell-expansion",
			envs: []*golang.Environ{moduleOnEnv},
			wd:   filepath.Join(cmdRoot, "test/resolve-modules"),
			in: []string{
				"github.com/u-root/u-root/cmds/core/{init,ip,dhclient}",
			},
			want: []string{
				"github.com/u-root/u-root/cmds/core/dhclient",
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/ip",
			},
		},
		// Exclusion, single package, file system path.
		{
			name: "fspath-exclusion",
			in:   []string{"./test/goglob/*", "-test/goglob/echo"},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/foo"},
		},
		// Exclusion, single package, Go package path.
		{
			name: "pkgpath-exclusion",
			in:   []string{"./test/goglob/...", "-github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/echo"},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/foo"},
		},
		// Exclusion, single package, mixed.
		{
			name: "path-exclusion",
			in:   []string{"./test/goglob/...", "-test/goglob/echo"},
			want: []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/foo"},
		},
		// Globs in exclusions should work.
		//
		// Unless we put u-root and p9 in GOPATH in the local version
		// of this test, this is an ON only test.
		{
			name: "pkgpath-multi-module-exclusion-glob",
			envs: []*golang.Environ{moduleOnEnv},
			wd:   filepath.Join(cmdRoot, "test/resolve-modules"),
			in: []string{
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/yes",
				"github.com/hugelgupf/p9/cmd/p9ufs",
				"-github.com/u-root/u-root/cmds/core/y*",
			},
			want: []string{
				"github.com/hugelgupf/p9/cmd/p9ufs",
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/ip",
			},
		},
		// File system path. Not a directory.
		{
			name:    "fspath-not-a-directory",
			in:      []string{"./bb_test.go"},
			wantErr: true,
			err:     errNoMatch,
		},
		// Some error cases where $GOROOT/bin/go is unavailable, so packages.Load fails.
		/*
			TODO: re-enable when https://github.com/golang/go/issues/62114 is resolved.
			{
				name:    "fspath-load-fails",
				envs:    []*golang.Environ{noGoToolEnv},
				in:      []string{"./test/goglob/*"},
				wantErr: true,
			},
			{
				name:    "pkgpath-batched-load-fails",
				envs:    []*golang.Environ{noGoToolEnv},
				in:      []string{"./test/goglob/..."},
				wantErr: true,
			},
			{
				name:    "pkgpath-glob-load-fails",
				envs:    []*golang.Environ{noGoToolEnv},
				in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/goglob/*"},
				wantErr: true,
			},
		*/
	}

	for _, tc := range sharedTestCases {
		envs := []*golang.Environ{moduleOffEnv, moduleOnEnv}
		if tc.envs != nil {
			envs = tc.envs
		}
		for _, env := range envs {
			env = env.Copy(golang.WithWorkingDir(tc.wd))
			t.Run(fmt.Sprintf("ResolveGlobs-GO111MODULE=%s-%s", env.GO111MODULE, tc.name), func(t *testing.T) {
				e := Env{
					GBBPath: tc.gbbPath,
				}
				out, err := ResolveGlobs(l, env, e, tc.in)
				if tc.err != nil && !errors.Is(err, tc.err) {
					t.Errorf("ResolveGlobs(%v, %v) = %v, want %v", e, tc.in, err, tc.err)
				}
				if (err != nil) != tc.wantErr {
					t.Errorf("ResolveGlobs(%v, %v) = (%v, %v), wantErr is %t", e, tc.in, out, err, tc.wantErr)
				}
				if len(out) < tc.atLeast {
					t.Errorf("ResolveGlobs(%v, %v) = %d items; want at least %d", e, tc.in, len(out), tc.atLeast)
				}
				if tc.atLeast > 0 {
					return
				}
				if !reflect.DeepEqual(out, tc.want) {
					t.Errorf("ResolveGlobs(%v, %v) = %#v; want %#v", e, tc.in, out, tc.want)
				}
			})
		}
	}

	//noGopathModuleOffEnv := golang.Default(golang.WithGO111MODULE("off"), golang.WithGOPATH(t.TempDir()))

	newPkgTests := append(sharedTestCases, testCase{
		name:    "fspath-parse-broken",
		in:      []string{"./test/parsebroken"},
		wantErr: true,
	}, testCase{
		name:    "pkgpath-parse-broken",
		in:      []string{"github.com/u-root/u-root/cmds/exp/cmd2pkg/findpkg/test/parsebroken"},
		wantErr: true,
	})
	for _, tc := range newPkgTests {
		envs := []*golang.Environ{moduleOffEnv, moduleOnEnv}
		if tc.envs != nil {
			envs = tc.envs
		}
		for _, env := range envs {
			env = env.Copy(golang.WithWorkingDir(tc.wd))
			t.Run(fmt.Sprintf("NewPackage-GO111MODULE=%s-%s", env.GO111MODULE, tc.name), func(t *testing.T) {
				e := Env{
					GBBPath: tc.gbbPath,
				}
				out, err := NewPackages(l, env, e, tc.in...)
				if tc.err != nil && !errors.Is(err, tc.err) {
					t.Errorf("NewPackages(%v, %v) = %v, want %v", e, tc.in, err, tc.err)
				}
				if (err != nil) != tc.wantErr {
					t.Errorf("NewPackages(%v, %v) = (%v, %v), wantErr is %t", e, tc.in, out, err, tc.wantErr)
				}

				if len(out) < tc.atLeast {
					t.Errorf("ResolveGlobs(%v, %v) = %d items; want at least %d", e, tc.in, len(out), tc.atLeast)
				}
				if tc.atLeast > 0 {
					return
				}
				var pkgPaths []string
				for _, p := range out {
					pkgPaths = append(pkgPaths, p.Pkg.PkgPath)
				}
				sort.Strings(pkgPaths)
				if !reflect.DeepEqual(pkgPaths, tc.want) {
					t.Errorf("NewPackages(%v, %v) = %v; want %v", e, tc.in, out, tc.want)
				}
			})

		}
	}
}

func TestDefaultEnv(t *testing.T) {
	for _, tc := range []struct {
		GBB_PATH     string
		UROOT_SOURCE string
		s            string
		want         Env
	}{
		{
			GBB_PATH:     "foo:bar",
			UROOT_SOURCE: "./foo",
			s:            "GBB_PATH=foo:bar UROOT_SOURCE=./foo",
			want: Env{
				GBBPath:     []string{"foo", "bar"},
				URootSource: "./foo",
			},
		},
		{
			GBB_PATH: "foo",
			s:        "GBB_PATH=foo UROOT_SOURCE=",
			want: Env{
				GBBPath: []string{"foo"},
			},
		},
		{
			s:    "GBB_PATH= UROOT_SOURCE=",
			want: Env{},
		},
	} {
		t.Run(tc.s, func(t *testing.T) {
			os.Setenv("GBB_PATH", tc.GBB_PATH)
			os.Setenv("UROOT_SOURCE", tc.UROOT_SOURCE)
			e := DefaultEnv()
			if !reflect.DeepEqual(e, tc.want) {
				t.Errorf("Env = %#v, want %#v", e, tc.want)
			}
			if e.String() != tc.s {
				t.Errorf("Env.String() = %v, want %v", e, tc.s)
			}
		})
	}
}

func TestModules(t *testing.T) {
	dir := t.TempDir()

	_ = os.MkdirAll(filepath.Join(dir, "mod1/cmd/cmd1"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "mod1/cmd/cmd2"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "mod1/nestedmod1/cmd/cmd5"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "mod1/nestedmod2/cmd/cmd6"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "mod2/cmd/foo3"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "mod2/cmd/foo4"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, "nomod/cmd/cmd7"), 0755)
	_ = os.WriteFile(filepath.Join(dir, "mod1/go.mod"), nil, 0644)
	_ = os.WriteFile(filepath.Join(dir, "mod1/nestedmod1/go.mod"), nil, 0644)
	_ = os.WriteFile(filepath.Join(dir, "mod1/nestedmod2/go.mod"), nil, 0644)
	_ = os.WriteFile(filepath.Join(dir, "mod2/go.mod"), nil, 0644)

	paths := []string{
		filepath.Join(dir, "mod1/cmd/cmd1"),
		filepath.Join(dir, "mod1/cmd/cmd2"),
		filepath.Join(dir, "mod1/nestedmod1/cmd/cmd5"),
		filepath.Join(dir, "mod1/nestedmod2/cmd/cmd6"),
		filepath.Join(dir, "mod2/cmd/foo3"),
		filepath.Join(dir, "mod2/cmd/foo4"),
		filepath.Join(dir, "nomod/cmd/cmd7"),
	}
	mods, noModulePkgs := Modules(paths)

	want := map[string][]string{
		filepath.Join(dir, "mod1"): {
			filepath.Join(dir, "mod1/cmd/cmd1"),
			filepath.Join(dir, "mod1/cmd/cmd2"),
		},
		filepath.Join(dir, "mod1/nestedmod1"): {
			filepath.Join(dir, "mod1/nestedmod1/cmd/cmd5"),
		},
		filepath.Join(dir, "mod1/nestedmod2"): {
			filepath.Join(dir, "mod1/nestedmod2/cmd/cmd6"),
		},
		filepath.Join(dir, "mod2"): {
			filepath.Join(dir, "mod2/cmd/foo3"),
			filepath.Join(dir, "mod2/cmd/foo4"),
		},
	}
	if !reflect.DeepEqual(mods, want) {
		t.Errorf("modules() = %v, want %v", mods, want)
	}
	wantNoModule := []string{
		filepath.Join(dir, "nomod/cmd/cmd7"),
	}
	if !reflect.DeepEqual(noModulePkgs, wantNoModule) {
		t.Errorf("modules() no module pkgs = %v, want %v", noModulePkgs, wantNoModule)
	}

	wantG := []string{
		filepath.Join(dir, "mod1/cmd/cmd1"),
		filepath.Join(dir, "mod2/cmd/foo3"),
	}
	e := Env{
		GBBPath: []string{filepath.Join(dir, "mod1"), filepath.Join(dir, "mod2")},
	}
	if got := GlobPaths(&ulogtest.Logger{TB: t}, e, `cmd/cmd*`, "cmd/foo3", "-cmd/cmd2"); !reflect.DeepEqual(got, wantG) {
		t.Errorf("GlobPaths = %v, want %v", got, wantG)
	}
}
