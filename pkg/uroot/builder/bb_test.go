// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"debug/elf"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gbb "github.com/u-root/gobusybox/src/pkg/bb"
	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

func TestBBBuildPreModules(t *testing.T) {
	if os.Getenv("GO111MODULE") != "off" {
		t.Skip("Skipping non-modular test")
	}
	dir := t.TempDir()

	opts := Opts{
		Env: golang.Default(),
		Packages: []string{
			"github.com/u-root/u-root/pkg/uroot/test/foo",
			"github.com/u-root/u-root/cmds/core/elvish",
		},
		TempDir:   dir,
		BinaryDir: "bbin",
	}
	af := initramfs.NewFiles()
	var bbb BBBuilder
	if err := bbb.Build(af, opts); err != nil {
		t.Error(err)
	}

	mustContain := []string{
		"bbin/elvish",
		"bbin/foo",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q; archive: %v", name, af)
		}
	}
}

// This is a simple test for use of the go busy box.
// It technically need not be here, but it's not bad to make sure
// things are still working.
func TestBBBuildModules(t *testing.T) {
	if os.Getenv("GO111MODULE") == "off" {
		t.Skip("Skipping modular test")
	}
	dir := t.TempDir()

	bopts := &gbbgolang.BuildOpts{}
	bopts.RegisterFlags(flag.CommandLine)

	o := filepath.Join(dir, "initramfs.cpio")

	env := gbbgolang.Default()
	if env.CgoEnabled {
		t.Logf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	t.Logf("Build environment: %s", env)

	tmpDir, err := ioutil.TempDir("", "bb-")
	if err != nil {
		t.Fatalf("Could not create busybox source directory: %v", err)
	}

	opts := &gbb.Opts{
		Env:       env,
		GenSrcDir: tmpDir,
		CommandPaths: []string{
			"../test/foo",
			"../../..//cmds/core/elvish",
		},
		BinaryPath:  o,
		GoBuildOpts: bopts,
	}
	if err := gbb.BuildBusybox(opts); err != nil {
		var errGopath *gbb.ErrGopathBuild
		var errGomod *gbb.ErrModuleBuild
		if errors.As(err, &errGopath) {
			t.Fatalf("preserving bb generated source directory at %s due to error. To reproduce build, `cd %s` and `GO111MODULE=off GOPATH=%s go build`", tmpDir, errGopath.CmdDir, errGopath.GOPATH)
		} else if errors.As(err, &errGomod) {
			t.Fatalf("preserving bb generated source directory at %s due to error. To debug build, `cd %s` and use `go build` to build, or `go mod [why|tidy|graph]` to debug dependencies, or `go list -m all` to list all dependency versions", tmpDir, errGomod.CmdDir)
		} else {
			t.Fatalf("preserving bb generated source directory at %s due to error", tmpDir)
		}
	}

	// The gobusybox code we use does not build an initramfs. It builds an elf program.
	// The only test, therefore, is if the output file can be read as an ELF.
	f, err := elf.Open(o)
	if err != nil {
		t.Fatalf("Opening initramfs(%v): got %v, want nil", o, err)
	}

	// A valid ELF file has at least one loadable Program.
	// More than this, we can not say.
	var foundLoadable bool
	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}
		foundLoadable = true
	}
	if !foundLoadable {
		t.Errorf("ELF check: %q has no segments with PT_LOAD flag, want at least one", o)
	}
}
