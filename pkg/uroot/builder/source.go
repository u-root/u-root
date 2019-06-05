// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

var (
	goCommandFile      = "zzzzinit.go"
	addInitToGoCommand = []byte(`// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

)

func init() {
	if os.Args[0] != "/init" {
		return
	}

	c := exec.Command("/go/bin/go", "build", "-o", "/buildbin/installcommand", "github.com/u-root/u-root/cmds/core/installcommand")
	c.Env = append(c.Env,  []string{"GOROOT=/go", "GOPATH=/",}...)
	o, err := c.CombinedOutput()
	if err != nil {
		log.Printf("building installcommand: %s, %v", string(o), err)
		return
	}
	if err := syscall.Exec("/buildbin/init", []string{"init"}, []string{}); err != nil {
		log.Printf("Exec of /buildbin/init failed. %v", err)
	}
}
`)
)

// SourceBuilder includes full source for Go commands in the initramfs.
//
// SourceBuilder is an implementation of Builder.
//
// It also includes the Go toolchain in the initramfs, and a tool called
// installcommand that can compile the other commands using symlinks.
//
// E.g. if "ls" is an included command, "ls" will be a symlink to
// "installcommand" in the initramfs, which uses argv[0] to figure out which
// command to compile.
type SourceBuilder struct {
	// FourBins, if true, will cause us to not build
	// an installcommand. This only makes sense if you are using the
	// fourbins command in the u-root command, but that's your call.
	// In operation, the default behavior is the one most people will want,
	// i.e. the installcommand will be built.
	FourBins bool
}

// DefaultBinaryDir implements Builder.DefaultBinaryDir.
//
// The initramfs default binary dir is buildbin.
func (SourceBuilder) DefaultBinaryDir() string {
	return "buildbin"
}

// Build is an implementation of Builder.Build.
func (sb SourceBuilder) Build(af *initramfs.Files, opts Opts) error {
	// TODO: this is a failure to collect the correct dependencies.
	if err := af.AddFile(filepath.Join(opts.Env.GOROOT, "pkg/include"), "go/pkg/include"); err != nil {
		return err
	}

	var installcommand string
	log.Printf("Collecting package files and dependencies...")
	deps := make(map[string]struct{})
	for _, pkg := range opts.Packages {
		name := path.Base(pkg)
		if name == "installcommand" {
			installcommand = pkg
		}

		// Add high-level packages' src files to archive.
		p := goListPkg(opts, pkg, af)
		if p == nil {
			continue
		}
		for _, d := range p.Deps {
			deps[d] = struct{}{}
		}

		if name != "installcommand" {
			// Add a symlink to installcommand. This means source mode can
			// work with any init.
			if err := af.AddRecord(cpio.Symlink(path.Join(opts.BinaryDir, name), "installcommand")); err != nil {
				return err
			}
		}
	}
	if len(installcommand) == 0 {
		return fmt.Errorf("must include a version of installcommand in source mode")
	}

	// Add src files of dependencies to archive.
	for dep := range deps {
		goListPkg(opts, dep, af)
	}

	// If we are doing "four bins" mode, or maybe I should call it Go of
	// Four, then we need to drop a file into the Go command source
	// directory before we build, and we need to remove it after.  And we
	// need to verify that we're not supplanting something.
	if sb.FourBins {
		goCmd := filepath.Join(opts.Env.GOROOT, "src/cmd/go")
		if _, err := os.Stat(goCmd); err != nil {
			return fmt.Errorf("stat(%q): %v", goCmd, err)
		}

		z := filepath.Join(goCmd, goCommandFile)
		if _, err := os.Stat(z); err == nil {
			return fmt.Errorf("%q exists, and we will not overwrite it", z)
		}

		if err := ioutil.WriteFile(z, addInitToGoCommand, 0444); err != nil {
			return err
		}
		defer os.Remove(z)
	}

	// Add Go toolchain.
	log.Printf("Building go toolchain...")
	if err := buildToolchain(opts); err != nil {
		return err
	}
	if !sb.FourBins {
		if err := opts.Env.Build(installcommand, filepath.Join(opts.TempDir, opts.BinaryDir, "installcommand"), golang.BuildOpts{}); err != nil {
			return err
		}
	}

	// Add Go toolchain and installcommand to archive.
	return af.AddFile(opts.TempDir, "")
}

// buildToolchain builds the needed Go toolchain binaries: go, compile, link,
// asm.
func buildToolchain(opts Opts) error {
	goBin := filepath.Join(opts.TempDir, "go/bin/go")
	tcbo := golang.BuildOpts{
		ExtraArgs: []string{"-tags", "cmd_go_bootstrap"},
	}
	if err := opts.Env.Build("cmd/go", goBin, tcbo); err != nil {
		return err
	}

	toolDir := filepath.Join(opts.TempDir, fmt.Sprintf("go/pkg/tool/%v_%v", opts.Env.GOOS, opts.Env.GOARCH))
	for _, pkg := range []string{"compile", "link", "asm"} {
		c := filepath.Join(toolDir, pkg)
		if err := opts.Env.Build(fmt.Sprintf("cmd/%s", pkg), c, golang.BuildOpts{}); err != nil {
			return err
		}
	}
	return nil
}

func goListPkg(opts Opts, importPath string, out *initramfs.Files) *golang.ListPackage {
	p, err := opts.Env.Deps(importPath)
	if err != nil {
		log.Printf("Can't list Go dependencies for %v; ignoring.", importPath)
		return nil
	}

	// Add Go files in this package to archive.
	for _, file := range append(append(p.GoFiles, p.SFiles...), p.HFiles...) {
		relPath := filepath.Join("src", p.ImportPath, file)
		srcFile := filepath.Join(p.Root, relPath)
		if p.Goroot {
			out.AddFile(srcFile, filepath.Join("go", relPath))
		} else {
			out.AddFile(srcFile, relPath)
		}
	}
	return p
}
