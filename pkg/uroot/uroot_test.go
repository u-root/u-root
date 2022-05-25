// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/builder"
	itest "github.com/u-root/u-root/pkg/uroot/initramfs/test"
)

type inMemArchive struct {
	*cpio.Archive
}

// Finish implements initramfs.Writer.Finish.
func (inMemArchive) Finish() error { return nil }

func TestResolvePackagePaths(t *testing.T) {
	defaultEnv := golang.Default()
	urootpath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath1, err := filepath.Abs("test/gopath1")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath2, err := filepath.Abs("test/gopath2")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath1Env := defaultEnv
	gopath1Env.GOPATH = gopath1
	gopath2Env := defaultEnv
	gopath2Env.GOPATH = gopath2
	everythingEnv := defaultEnv
	everythingEnv.GOPATH = gopath1 + ":" + gopath2
	foopath, err := filepath.Abs("test/gopath1/src/foo")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	// Why doesn't the log package export this as a default?
	l := log.New(os.Stdout, "", log.LstdFlags)

	for _, tc := range []struct {
		env      golang.Environ
		in       []string
		expected []string
		wantErr  bool
	}{
		// Nonexistent Package
		{
			env:      defaultEnv,
			in:       []string{"fakepackagename"},
			expected: nil,
			wantErr:  true,
		},
		// Single go package import
		{
			env:      defaultEnv,
			in:       []string{"github.com/u-root/u-root/cmds/core/ls"},
			expected: []string{filepath.Join(urootpath, "cmds/core/ls")},
			wantErr:  false,
		},
		// Single package directory relative to working dir
		{
			env:      defaultEnv,
			in:       []string{"test/gopath1/src/foo"},
			expected: []string{filepath.Join(urootpath, "/pkg/uroot/test/gopath1/src/foo")},
			wantErr:  false,
		},
		// Single package directory with absolute path
		{
			env:      defaultEnv,
			in:       []string{foopath},
			expected: []string{filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo")},
			wantErr:  false,
		},
		// Single package directory relative to GOPATH
		{
			env:      gopath1Env,
			in:       []string{filepath.Join(gopath1, "src/foo")},
			expected: []string{filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo")},
			wantErr:  false,
		},
		// Package directory glob
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkg*"},
			expected: []string{
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
			},
			wantErr: false,
		},
		// GOPATH glob
		{
			env: gopath2Env,
			in:  []string{filepath.Join(gopath2, "src/mypkg*")},
			expected: []string{
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
			},
			wantErr: false,
		},
		// Single ambiguous package - exists in both GOROOT and GOPATH
		// This test doesn't work anymore with file paths
		// {
		// 	env: gopath1Env,
		// 	in:  []string{"os"},
		// 	expected: []string{
		// 		"os",
		// 	},
		// 	wantErr: false,
		// },
		// Packages from different gopaths
		{
			env: everythingEnv,
			in:  []string{filepath.Join(gopath1, "src/foo"), filepath.Join(gopath2, "src/mypkga")},
			expected: []string{
				filepath.Join(urootpath, "pkg/uroot/test/gopath1/src/foo"),
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
			},
			wantErr: false,
		},
		// Same package specified twice
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkga", "test/gopath2/src/mypkga"},
			// TODO: This returns the package twice. Is this preferred?
			expected: []string{
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkga"),
			},
			wantErr: false,
		},
		// Excludes
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/*", "-test/gopath2/src/mypkga"},
			expected: []string{
				filepath.Join(urootpath, "pkg/uroot/test/gopath2/src/mypkgb"),
			},
			wantErr: false,
		},
	} {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := ResolvePackagePaths(l, tc.env, urootpath, tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolvePackagePaths(%#v, %v) err != nil is %v, want %v\nerr is %v",
					tc.env, tc.in, err != nil, tc.wantErr, err)
			}
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("ResolvePackagePaths(%#v, %v) = %v; want %v",
					tc.env, tc.in, out, tc.expected)
			}
		})
	}
}

func TestCreateInitramfs(t *testing.T) {
	dir := t.TempDir()
	syscall.Umask(0)

	urootpath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	tmp777 := filepath.Join(dir, "tmp777")
	if err := os.MkdirAll(tmp777, 0o777); err != nil {
		t.Error(err)
	}

	// Why doesn't the log package export this as a default?
	l := log.New(os.Stdout, "", log.LstdFlags)

	for i, tt := range []struct {
		name       string
		opts       Opts
		want       string
		validators []itest.ArchiveValidator
	}{
		{
			name: "BB archive with ls and init",
			opts: Opts{
				Env:             golang.Default(),
				TempDir:         dir,
				ExtraFiles:      nil,
				UseExistingInit: false,
				InitCmd:         "init",
				DefaultShell:    "ls",
				UrootSource:     urootpath,
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
				Env:          golang.Default(),
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
				Env:     golang.Default(),
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
				Env:          golang.Default(),
				TempDir:      dir,
				DefaultShell: "zoocar",
				InitCmd:      "foobar",
				UrootSource:  urootpath,
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
				Env:     golang.Default(),
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
				Env:             golang.Default(),
				TempDir:         dir,
				ExtraFiles:      nil,
				UseExistingInit: false,
				InitCmd:         "init",
				DefaultShell:    "ls",
				UrootSource:     urootpath,
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
