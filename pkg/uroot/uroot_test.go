// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/u-root/gobusybox/src/pkg/golang"
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

	l := ulogtest.Logger{TB: t}

	for i, tt := range []struct {
		name       string
		opts       Opts
		want       string
		validators []itest.ArchiveValidator
	}{
		{
			name: "BB archive with ls and init",
			opts: Opts{
				Env:             golang.Default(golang.DisableCGO()),
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
				Env:          golang.Default(golang.DisableCGO()),
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
				Env:             golang.Default(golang.DisableCGO()),
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
