// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package main

import (
	"os"
	"os/exec"
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
			wantErr: ErrUsage.Error(),
		},
		{
			name:    "err in pox.Create",
			args:    []string{},
			create:  true,
			file:    f,
			wantErr: ErrUsage.Error(),
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
			c := cmd{
				debug:   func(s string, i ...interface{}) {},
				verbose: tt.verbose,
				run:     tt.run,
				create:  tt.create,
				file:    tt.file,
				extra:   tt.extra,
				args:    tt.args,
			}
			if tt.name == "err in extraMounts(os.Getenv('POX_EXTRA'))" {
				os.Setenv("POX_EXTRA", "::")
			}
			if got := c.start(); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("cmd.start() = %q, want: %q", got.Error(), tt.wantErr)
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
			wantErr: ErrUsage.Error(),
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
			c := cmd{
				debug: func(s string, i ...interface{}) {},
				zip:   tt.zip,
				self:  tt.self,
				file:  tt.file,
				args:  tt.args,
			}
			if got := c.poxCreate(); got != nil {
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
			c := cmd{
				debug: func(s string, i ...interface{}) {},
				zip:   tt.zip,
				file:  tt.file,
				args:  tt.args,
			}
			if got := c.poxRun(); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("poxRun() = %q, want: %q", got.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestSelfEmbedding(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skipf("skipping; must be root")
	}
	pwd, err := exec.LookPath("pwd")
	if err != nil {
		t.Skip("no pwd on this system")
	}

	// The above covers many code paths. The code below tests
	// self-embedding.
	if err := exec.Command("go", "build", "-o", "pox").Run(); err != nil {
		t.Fatalf("building pox: got %v, want nil", err)
	}

	if err := exec.Command("./pox", "-cvsf", "pwd.pox", pwd).Run(); err != nil {
		t.Fatalf("building pwd.pox: got %v, want nil", err)
	}

	out, err := exec.Command("./pwd.pox").CombinedOutput()
	if err != nil {
		t.Fatalf("running pwd.pox: got (%s,%v), want nil", out, err)
	}
	if string(out) != "/\n" {
		t.Fatalf("running pwd.pox: got %s, want %s", string(out), "/")
	}
}
