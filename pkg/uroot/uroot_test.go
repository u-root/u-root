// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/builder"
)

type inMemArchive struct {
	*cpio.Archive
}

// Finish implements initramfs.Writer.Finish.
func (inMemArchive) Finish() error { return nil }

func TestResolvePackagePaths(t *testing.T) {
	defaultEnv := golang.Default()
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
			env: defaultEnv,
			in:  []string{"github.com/u-root/u-root/cmds/ls"},
			// We expect the full URL format because that's the path in our default GOPATH
			expected: []string{"github.com/u-root/u-root/cmds/ls"},
			wantErr:  false,
		},
		// Single package directory relative to working dir
		{
			env:      defaultEnv,
			in:       []string{"test/gopath1/src/foo"},
			expected: []string{"github.com/u-root/u-root/pkg/uroot/test/gopath1/src/foo"},
			wantErr:  false,
		},
		// Single package directory with absolute path
		{
			env:      defaultEnv,
			in:       []string{foopath},
			expected: []string{"github.com/u-root/u-root/pkg/uroot/test/gopath1/src/foo"},
			wantErr:  false,
		},
		// Single package directory relative to GOPATH
		{
			env: gopath1Env,
			in:  []string{"foo"},
			expected: []string{
				"foo",
			},
			wantErr: false,
		},
		// Package directory glob
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkg*"},
			expected: []string{
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkgb",
			},
			wantErr: false,
		},
		// GOPATH glob
		{
			env: gopath2Env,
			in:  []string{"mypkg*"},
			expected: []string{
				"mypkga",
				"mypkgb",
			},
			wantErr: false,
		},
		// Single ambiguous package - exists in both GOROOT and GOPATH
		{
			env: gopath1Env,
			in:  []string{"os"},
			expected: []string{
				"os",
			},
			wantErr: false,
		},
		// Packages from different gopaths
		{
			env: everythingEnv,
			in:  []string{"foo", "mypkga"},
			expected: []string{
				"foo",
				"mypkga",
			},
			wantErr: false,
		},
		// Same package specified twice
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkga", "test/gopath2/src/mypkga"},
			// TODO: This returns the package twice. Is this preferred?
			expected: []string{
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
			},
			wantErr: false,
		},
	} {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := ResolvePackagePaths(tc.env, tc.in)
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

type archiveValidator interface {
	validate(a *cpio.Archive) error
}

type hasRecord struct {
	r cpio.Record
}

func (hr hasRecord) validate(a *cpio.Archive) error {
	r, ok := a.Get(hr.r.Name)
	if !ok {
		return fmt.Errorf("archive does not contain %v", hr.r)
	}
	if !cpio.Equal(r, hr.r) {
		return fmt.Errorf("archive does not contain %v; instead has %v", hr.r, r)
	}
	return nil
}

type hasFile struct {
	path string
}

func (hf hasFile) validate(a *cpio.Archive) error {
	if _, ok := a.Get(hf.path); ok {
		return nil
	}
	return fmt.Errorf("archive does not contain %s, but should", hf.path)
}

type missingFile struct {
	path string
}

func (mf missingFile) validate(a *cpio.Archive) error {
	if _, ok := a.Get(mf.path); ok {
		return fmt.Errorf("archive contains %s, but shouldn't", mf.path)
	}
	return nil
}

type isEmpty struct{}

func (isEmpty) validate(a *cpio.Archive) error {
	if empty := a.Empty(); !empty {
		return fmt.Errorf("expected archive to be empty")
	}
	return nil
}

func TestCreateInitramfs(t *testing.T) {
	dir, err := ioutil.TempDir("", "foo")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)
	syscall.Umask(0)

	tmp777 := filepath.Join(dir, "tmp777")
	if err := os.MkdirAll(tmp777, 0777); err != nil {
		t.Error(err)
	}

	for i, tt := range []struct {
		name       string
		opts       Opts
		want       error
		validators []archiveValidator
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
				Commands: []Commands{
					{
						Builder: builder.BBBuilder{},
						Packages: []string{
							"github.com/u-root/u-root/cmds/init",
							"github.com/u-root/u-root/cmds/ls",
						},
					},
				},
			},
			want: nil,
			validators: []archiveValidator{
				hasFile{path: "bbin/bb"},
				hasRecord{cpio.Symlink("bbin/init", "bb")},
				hasRecord{cpio.Symlink("bbin/ls", "bb")},
				hasRecord{cpio.Symlink("bin/defaultsh", "/bbin/ls")},
			},
		},
		{
			name: "no temp dir",
			opts: Opts{
				Env:          golang.Default(),
				InitCmd:      "init",
				DefaultShell: "",
			},
			// TODO: Ew. Our error types suck.
			want: fmt.Errorf("temp dir \"\" must exist: stat : no such file or directory"),
			validators: []archiveValidator{
				isEmpty{},
			},
		},
		{
			name: "no commands",
			opts: Opts{
				Env:     golang.Default(),
				TempDir: dir,
			},
			want: nil,
			validators: []archiveValidator{
				missingFile{"bbin/bb"},
			},
		},
		{
			name: "init specified, but not in commands",
			opts: Opts{
				Env:          golang.Default(),
				TempDir:      dir,
				DefaultShell: "zoocar",
				InitCmd:      "foobar",
			},
			want: fmt.Errorf("could not find init: command or path \"foobar\" not included in u-root build"),
			validators: []archiveValidator{
				isEmpty{},
			},
		},
		{
			name: "init symlinked to absolute path",
			opts: Opts{
				Env:     golang.Default(),
				TempDir: dir,
				InitCmd: "/bin/systemd",
			},
			want: nil,
			validators: []archiveValidator{
				hasRecord{cpio.Symlink("init", "/bin/systemd")},
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
				Commands: []Commands{
					{
						Builder: builder.BBBuilder{},
						Packages: []string{
							"github.com/u-root/u-root/cmds/init",
							"github.com/u-root/u-root/cmds/ls",
						},
					},
					{
						Builder: builder.BinaryBuilder{},
						Packages: []string{
							"github.com/u-root/u-root/cmds/cp",
							"github.com/u-root/u-root/cmds/dd",
						},
					},
					{
						Builder: builder.SourceBuilder{},
						Packages: []string{
							"github.com/u-root/u-root/cmds/cat",
							"github.com/u-root/u-root/cmds/chroot",
							"github.com/u-root/u-root/cmds/installcommand",
						},
					},
				},
			},
			want: nil,
			validators: []archiveValidator{
				hasRecord{cpio.Symlink("init", "/bbin/init")},

				// bb mode.
				hasFile{path: "bbin/bb"},
				hasRecord{cpio.Symlink("bbin/init", "bb")},
				hasRecord{cpio.Symlink("bbin/ls", "bb")},
				hasRecord{cpio.Symlink("bin/defaultsh", "/bbin/ls")},

				// binary mode.
				hasFile{path: "bin/cp"},
				hasFile{path: "bin/dd"},

				// source mode.
				hasRecord{cpio.Symlink("buildbin/cat", "/buildbin/installcommand")},
				hasRecord{cpio.Symlink("buildbin/chroot", "/buildbin/installcommand")},
				hasFile{path: "buildbin/installcommand"},
				hasFile{path: "src/github.com/u-root/u-root/cmds/cat/cat.go"},
				hasFile{path: "src/github.com/u-root/u-root/cmds/chroot/chroot.go"},
				hasFile{path: "src/github.com/u-root/u-root/cmds/installcommand/installcommand.go"},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %d [%s]", i, tt.name), func(t *testing.T) {
			archive := inMemArchive{cpio.InMemArchive()}
			tt.opts.OutputFile = archive
			// Compare error type or error string.
			if err := CreateInitramfs(tt.opts); err != tt.want && (tt.want == nil || err.Error() != tt.want.Error()) {
				t.Errorf("CreateInitramfs(%v) = %v, want %v", tt.opts, err, tt.want)
			}

			for _, v := range tt.validators {
				if err := v.validate(archive.Archive); err != nil {
					t.Errorf("validator failed: %v / archive:\n%s", err, archive)
				}
			}
		})
	}
}
