// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && !race
// +build !tinygo,!race

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestIntegrationPox(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/cmds/exp/pox"),
		govmtest.WithQEMUFn(qemu.WithVMTimeout(time.Minute)),
	)
}

func TestPox(t *testing.T) {
	guest.SkipIfNotInVM(t)

	f := filepath.Join(t.TempDir(), "x.tcz")
	tmpFile, err := os.Create("/bin/bash")
	if err != nil {
		t.Errorf("Couldn't create /bin/bash: %v", err)
	}
	defer tmpFile.Close()
	for _, tt := range []struct {
		name    string
		args    []string
		verbose bool
		run     bool
		create  bool
		file    string
		extra   string
		wantErr string
	}{
		{
			name:    "err in usage",
			args:    []string{"/bin/bash"},
			file:    f,
			wantErr: "pox [-[-verbose]|v] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]",
		},
		{
			name:    "err in pox.Create",
			args:    []string{},
			create:  true,
			file:    f,
			wantErr: "pox [-[-verbose]|v] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]",
		},
		{
			name:    "err in extraMounts(*extra)",
			args:    []string{"/bin/bash"},
			extra:   "::",
			file:    f,
			wantErr: "[\"\" \"\" \"\"] is not in the form src:target",
		},
		{
			name:    "err in extraMounts(os.Getenv('POX_EXTRA'))",
			args:    []string{"/bin/bash"},
			verbose: true,
			file:    f,
			wantErr: "[\"\" \"\" \"\"] is not in the form src:target",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*verbose = tt.verbose
			*run = tt.run
			*create = tt.create
			*file = tt.file
			*extra = tt.extra
			if tt.name == "err in extraMounts(os.Getenv('POX_EXTRA'))" {
				os.Setenv("POX_EXTRA", "::")
			}
			if got := pox(tt.args...); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("pox() = %q, want: %q", got.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestPoxCreate(t *testing.T) {
	guest.SkipIfNotInVM(t)

	f := filepath.Join(t.TempDir(), "x.tcz")
	tmpFile, err := os.Create("/bin/bash")
	if err != nil {
		t.Errorf("Couldn't create /bin/bash: %v", err)
	}
	defer tmpFile.Close()
	for _, tt := range []struct {
		name    string
		args    []string
		zip     bool
		self    bool
		file    string
		wantErr string
	}{
		{
			name:    "len(bin) == 0",
			args:    []string{},
			wantErr: "pox [-[-verbose]|v] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]",
		},
		{
			name:    "error in ldd.Ldd",
			args:    []string{""},
			wantErr: "running ldd on []: open : no such file or directory ",
		},
		{
			name:    "self = false, zip = false, err in c.CombinedOutput()",
			args:    []string{"/bin/bash"},
			wantErr: "executable file not found in $PATH",
		},
		{
			name: "self = true, zip = false, no err",
			args: []string{"/bin/bash"},
			self: true,
			file: f,
		},
		{
			name:    "self = true, zip = false, err in cp.Copy",
			args:    []string{"/bin/bash"},
			self:    true,
			file:    "",
			wantErr: "open : no such file or directory",
		},
		{
			name:    "self = false, zip = true, no err",
			args:    []string{"/bin/bash"},
			zip:     true,
			file:    f,
			wantErr: "open : no such file or directory",
		},
		{
			name:    "self = false, zip = true, err in uzip.ToZip",
			args:    []string{"/bin/bash"},
			zip:     true,
			file:    "",
			wantErr: "open : no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*zip = tt.zip
			*self = tt.self
			*file = tt.file
			if got := poxCreate(tt.args...); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("poxCreate() = %q, want: %q", got.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestPoxRun(t *testing.T) {
	guest.SkipIfNotInVM(t)

	f := filepath.Join(t.TempDir(), "x.tcz")
	if _, err := os.Create(f); err != nil {
		t.Errorf("Couldn't create file: %v", err)
	}
	for _, tt := range []struct {
		name    string
		args    []string
		zip     bool
		file    string
		wantErr string
	}{
		{
			name:    "len(args) == 0",
			args:    []string{},
			wantErr: "pox [-[-verbose]|v] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]",
		},
		{
			name:    "zip = true with error",
			args:    []string{"/bin/bash"},
			zip:     true,
			file:    f,
			wantErr: "zip: not a valid zip file",
		},
		{
			name:    "zip = false with error",
			args:    []string{"/bin/bash"},
			file:    "/mnt",
			wantErr: "open /mnt: no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*zip = tt.zip
			*file = tt.file
			if got := poxRun(tt.args...); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("poxRun() = %q, want: %q", got.Error(), tt.wantErr)
				}
			}
		})
	}
}
