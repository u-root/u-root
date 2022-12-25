// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
	"github.com/u-root/u-root/pkg/uroot/builder"
	itest "github.com/u-root/u-root/pkg/uroot/initramfs/test"
)

type inMemArchive struct {
	*cpio.Archive
}

// Finish implements initramfs.Writer.Finish.
func (inMemArchive) Finish() error { return nil }

func TestResolvePackagePathsSpecialCases(t *testing.T) {
	gopath1, err := filepath.Abs("test/gopath1")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath2, err := filepath.Abs("test/gopath2")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	everythingEnv := gbbgolang.Default()
	everythingEnv.GO111MODULE = "off"
	everythingEnv.GOPATH = gopath1 + ":" + gopath2

	l := &ulogtest.Logger{TB: t}

	for _, tc := range []struct {
		env      gbbgolang.Environ
		in       []string
		expected []string
		wantErr  bool
	}{
		{
			env: everythingEnv,
			in:  []string{"foo", "mypkga"},
			expected: []string{
				"foo",    // from gopath1
				"mypkga", // from gopath2
			},
			wantErr: false,
		},
	} {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := ResolvePackagePaths(l, tc.env, tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolvePackagePaths(%v, %v) = %v, want err is %t", tc.env, tc.in, err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("ResolvePackagePaths(%v, %v) = %v; want %v", tc.env, tc.in, out, tc.expected)
			}
		})
	}
}

func TestResolvePackagePathsUrootGOPATH(t *testing.T) {
	urootpath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	foopath, err := filepath.Abs("test/gopath1/src/foo")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	moduleOffEnv := gbbgolang.Default()
	moduleOffEnv.GO111MODULE = "off"

	moduleOnEnv := gbbgolang.Default()
	moduleOnEnv.GO111MODULE = "on"

	l := &ulogtest.Logger{TB: t}

	for _, env := range []gbbgolang.Environ{moduleOnEnv, moduleOffEnv} {
		for _, tc := range []struct {
			in       []string
			expected []string
			wantErr  bool
		}{
			// Nonexistent Package
			{
				in:       []string{"fakepackagename"},
				expected: nil,
				wantErr:  true,
			},
			// Single go package import
			{
				in:       []string{"github.com/u-root/u-root/cmds/core/ls"},
				expected: []string{"github.com/u-root/u-root/cmds/core/ls"},
				wantErr:  false,
			},
			// Single package directory relative to working dir
			{
				in:       []string{"test/gopath1/src/foo"},
				expected: []string{filepath.Join(urootpath, "/pkg/uroot/test/gopath1/src/foo")},
				wantErr:  false,
			},
			// Single package directory with absolute path
			{
				in:       []string{foopath},
				expected: []string{filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo")},
				wantErr:  false,
			},
			// Package directory glob
			{
				in: []string{"test/gopath2/src/mypkg*"},
				expected: []string{
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
				},
				wantErr: false,
			},
			// Same package specified twice
			{
				in: []string{"test/gopath2/src/mypkga", "test/gopath2/src/mypkga"},
				// TODO: This returns the package twice. Is this preferred?
				expected: []string{
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
				},
				wantErr: false,
			},
			// Excludes
			{
				in: []string{"test/gopath2/src/*", "-test/gopath2/src/mypkga"},
				expected: []string{
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
				},
				wantErr: false,
			},
		} {
			t.Run(fmt.Sprintf("GO111MODULE=%s-%q", env.GO111MODULE, tc.in), func(t *testing.T) {
				out, err := ResolvePackagePaths(l, env, tc.in)
				if (err != nil) != tc.wantErr {
					t.Fatalf("ResolvePackagePaths(%v, %v) = %v, want err is %t", env, tc.in, err, tc.wantErr)
				}
				if !reflect.DeepEqual(out, tc.expected) {
					t.Errorf("ResolvePackagePaths(%v, %v) = %v; want %v", env, tc.in, out, tc.expected)
				}
			})
		}
	}
}

func TestCreateInitramfs(t *testing.T) {
	dir := t.TempDir()
	syscall.Umask(0)

	tmp777 := filepath.Join(dir, "tmp777")
	if err := os.MkdirAll(tmp777, 0o777); err != nil {
		t.Error(err)
	}

	l := &ulogtest.Logger{TB: t}

	for i, tt := range []struct {
		name       string
		opts       Opts
		want       string
		validators []itest.ArchiveValidator
	}{
		{
			name: "BB archive with ls and init",
			opts: Opts{
				TempDir:         dir,
				ExtraFiles:      nil,
				UseExistingInit: false,
				InitCmd:         "init",
				DefaultShell:    "ls",
				Commands: []Commands{
					{
						Builder: builder.BusyBox,
						Packages: []string{
							"github.com/u-root/u-root/cmds/core/init",
							"github.com/u-root/u-root/cmds/core/ls",
						},
					},
				},
			},
			want: "",
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bbin/bb"},
				itest.HasRecord{cpio.Symlink("bbin/init", "bb")},
				itest.HasRecord{cpio.Symlink("bbin/ls", "bb")},
				itest.HasRecord{cpio.Symlink("bin/defaultsh", "../bbin/ls")},
				itest.HasRecord{cpio.Symlink("bin/sh", "../bbin/ls")},
			},
		},
		{
			name: "no temp dir",
			opts: Opts{
				InitCmd:      "init",
				DefaultShell: "",
			},
			want: "temp dir \"\" must exist: stat : no such file or directory",
			validators: []itest.ArchiveValidator{
				itest.IsEmpty{},
			},
		},
		{
			name: "no commands",
			opts: Opts{
				TempDir: dir,
			},
			want: "",
			validators: []itest.ArchiveValidator{
				itest.MissingFile{"bbin/bb"},
			},
		},
		{
			name: "init specified, but not in commands",
			opts: Opts{
				TempDir:      dir,
				DefaultShell: "zoocar",
				InitCmd:      "foobar",
				Commands: []Commands{
					{
						Builder: builder.Binary,
						Packages: []string{
							"github.com/u-root/u-root/cmds/core/ls",
						},
					},
				},
			},
			want: "could not create symlink from \"init\" to \"foobar\": command or path \"foobar\" not included in u-root build: specify -initcmd=\"\" to ignore this error and build without an init (or, did you specify a list, and are you missing github.com/u-root/u-root/cmds/core/init?)",
			validators: []itest.ArchiveValidator{
				itest.IsEmpty{},
			},
		},
		{
			name: "init symlinked to absolute path",
			opts: Opts{
				TempDir: dir,
				InitCmd: "/bin/systemd",
			},
			want: "",
			validators: []itest.ArchiveValidator{
				itest.HasRecord{cpio.Symlink("init", "bin/systemd")},
			},
		},
		{
			name: "multi-mode archive",
			opts: Opts{
				TempDir:         dir,
				ExtraFiles:      nil,
				UseExistingInit: false,
				InitCmd:         "init",
				DefaultShell:    "ls",
				Commands: []Commands{
					{
						Builder: builder.BusyBox,
						Packages: []string{
							"github.com/u-root/u-root/cmds/core/init",
							"github.com/u-root/u-root/cmds/core/ls",
						},
					},
					{
						Builder: builder.Binary,
						Packages: []string{
							"github.com/u-root/u-root/cmds/core/cp",
							"github.com/u-root/u-root/cmds/core/dd",
						},
					},
				},
			},
			want: "",
			validators: []itest.ArchiveValidator{
				itest.HasRecord{cpio.Symlink("init", "bbin/init")},

				// bb mode.
				itest.HasFile{"bbin/bb"},
				itest.HasRecord{cpio.Symlink("bbin/init", "bb")},
				itest.HasRecord{cpio.Symlink("bbin/ls", "bb")},
				itest.HasRecord{cpio.Symlink("bin/defaultsh", "../bbin/ls")},
				itest.HasRecord{cpio.Symlink("bin/sh", "../bbin/ls")},

				// binary mode.
				itest.HasFile{"bin/cp"},
				itest.HasFile{"bin/dd"},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %d [%s]", i, tt.name), func(t *testing.T) {
			archive := inMemArchive{cpio.InMemArchive()}
			tt.opts.OutputFile = archive
			// Compare error type or error string.
			if err := CreateInitramfs(l, tt.opts); (err != nil && err.Error() != tt.want) || (len(tt.want) > 0 && err == nil) {
				t.Errorf("CreateInitramfs(%v) = %v, want %v", tt.opts, err, tt.want)
			}

			for _, v := range tt.validators {
				if err := v.Validate(archive.Archive); err != nil {
					t.Errorf("validator failed: %v / archive:\n%s", err, archive)
				}
			}
		})
	}
}

func TestResolveGlobsSpecialCases(t *testing.T) {
	gopath1, err := filepath.Abs("test/gopath1")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	urootpath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath2, err := filepath.Abs("test/gopath2")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	moduleOffEnv := gbbgolang.Default()
	moduleOffEnv.GO111MODULE = "off"

	gopath1Env := moduleOffEnv
	gopath1Env.GOPATH = gopath1

	gopath2Env := moduleOffEnv
	gopath2Env.GOPATH = gopath2

	l := &ulogtest.Logger{TB: t}

	for _, tc := range []struct {
		env      gbbgolang.Environ
		in       string
		expected []string
		wantErr  bool
	}{
		// Single package directory relative to GOPATH
		{
			env:      gopath1Env,
			in:       filepath.Join(gopath1, "src/foo"),
			expected: []string{filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo")},
			wantErr:  false,
		},
		// Go import path glob
		{
			env: gopath2Env,
			in:  "mypkg*",
			expected: []string{
				"mypkga",
				"mypkgb",
			},
			wantErr: false,
		},
	} {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := resolveGlobs(l, tc.env, tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("resolveGlobs(%v, %v) = (%v, %v), wantErr is %t", tc.env, tc.in, out, err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("resolveGlobs(%v, %v) = %#v; want %#v", tc.env, tc.in, out, tc.expected)
			}
		})
	}
}

func TestResolveGlobsUrootGOPATH(t *testing.T) {
	urootpath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	foopath, err := filepath.Abs("test/gopath1/src/foo")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	moduleOffEnv := gbbgolang.Default()
	moduleOffEnv.GO111MODULE = "off"

	moduleOnEnv := gbbgolang.Default()
	moduleOnEnv.GO111MODULE = "on"

	l := &ulogtest.Logger{TB: t}

	for _, env := range []gbbgolang.Environ{moduleOffEnv, moduleOnEnv} {
		for _, tc := range []struct {
			in       string
			expected []string
			wantErr  bool
		}{
			// Nonexistent Package
			{
				in:       "fakepackagename",
				expected: nil,
				wantErr:  true,
			},
			// Single go package import
			{
				in:       "github.com/u-root/u-root/cmds/core/ls",
				expected: []string{"github.com/u-root/u-root/cmds/core/ls"},
				wantErr:  false,
			},
			// Single package directory relative to working dir
			{
				in:       "test/gopath1/src/foo",
				expected: []string{filepath.Join(urootpath, "/pkg/uroot/test/gopath1/src/foo")},
				wantErr:  false,
			},
			// Single package directory with absolute path
			{
				in:       foopath,
				expected: []string{filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo")},
				wantErr:  false,
			},
			// Single package with Plan 9 only files.
			{
				in:      "github.com/u-root/u-root/cmds/core/bind",
				wantErr: true,
			},
			// Single package with Plan 9 only files.
			{
				in:      filepath.Join(urootpath, "cmds/core/bind"),
				wantErr: true,
			},
			// Package directory glob
			{
				in: "test/gopath2/src/mypkg*",
				expected: []string{
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
					filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
				},
				wantErr: false,
			},
			// Package directory glob does not match anything
			{
				in:      "test/gopath2/src/foo*",
				wantErr: true,
			},
			// Go import path glob
			{
				in: "github.com/u-root/u-root/pkg/uroot/test/gopath2/src/my*",
				expected: []string{
					"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
					"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkgb",
				},
				wantErr: false,
			},
			// Glob does not match anything
			{
				in:      "github.com/u-root/u-root/pkg/uroot/test/gopath2/src/foo*",
				wantErr: true,
			},
			// Go glob support
			{
				in: "github.com/u-root/u-root/pkg/uroot/test/gopath2/src/...",
				expected: []string{
					"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
					"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkgb",
				},
				wantErr: false,
			},
			// Go import path glob
			{
				in: "github.com/u-root/u-root/cmds/core/i*",
				expected: []string{
					"github.com/u-root/u-root/cmds/core/id",
					"github.com/u-root/u-root/cmds/core/init",
					"github.com/u-root/u-root/cmds/core/insmod",
					"github.com/u-root/u-root/cmds/core/io",
					"github.com/u-root/u-root/cmds/core/ip",
				},
				wantErr: false,
			},
		} {
			t.Run(fmt.Sprintf("GO111MODULE=%s-%q", env.GO111MODULE, tc.in), func(t *testing.T) {
				out, err := resolveGlobs(l, env, tc.in)
				if (err != nil) != tc.wantErr {
					t.Fatalf("resolveGlobs(%v, %v) = (%v, %v), wantErr is %t", env, tc.in, out, err, tc.wantErr)
				}
				if !reflect.DeepEqual(out, tc.expected) {
					t.Errorf("resolveGlobs(%v, %v) = %#v; want %#v", env, tc.in, out, tc.expected)
				}
			})
		}
	}
}
