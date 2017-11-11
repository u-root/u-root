// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

// Flags for u-root builder.
var (
	build  = flag.String("build", "source", "u-root build format (e.g. bb or source)")
	format = flag.String("format", "cpio", "Archival format (e.g. cpio)")

	tmpDir = flag.String("tmpdir", "", "Temporary directory to put binaries in.")

	base            = flag.String("base", "", "Base archive to add files to")
	useExistingInit = flag.Bool("useinit", false, "Use existing init from base archive (only if --base was specified).")

	extraFiles = flag.String("files", "", "Additional files, directories, and binaries (with their ldd dependencies) to add to archive.")

	outputPath = flag.String("o", "", "Path to output initramfs file.")
)

func main() {
	flag.Parse()

	// Main is in a separate functions so defer's run on return.
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	env := golang.Default()
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)
	if env.GOOS != "linux" {
		log.Printf("GOOS is not linux. Did you mean to set GOOS=linux?")
	}

	builder, err := uroot.GetBuilder(*build)
	if err != nil {
		return err
	}
	archiver, err := uroot.GetArchiver(*format)
	if err != nil {
		return err
	}

	tempDir := *tmpDir
	if tempDir == "" {
		var err error
		tempDir, err = ioutil.TempDir("", "u-root")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
	} else if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return fmt.Errorf("temporary directory %q did not exist; tried to mkdir but failed: %v", tempDir, err)
		}
	}

	// Resolve globs into package imports.
	//
	// Currently allowed formats:
	//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/*
	pkgs := flag.Args()
	if len(pkgs) == 0 {
		var err error
		pkgs, err = uroot.DefaultPackageImports(env)
		if err != nil {
			return err
		}
	}

	// Open the target initramfs file.
	filename := *outputPath
	if filename == "" {
		filename = fmt.Sprintf("/tmp/initramfs.%s_%s.%s", env.GOOS, env.GOARCH, archiver.DefaultExtension())
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	var baseFile *os.File
	if *base != "" {
		var err error
		baseFile, err = os.Open(*base)
		if err != nil {
			return err
		}
		defer baseFile.Close()
	}

	opts := uroot.Opts{
		Env:             env,
		Builder:         builder,
		Archiver:        archiver,
		TempDir:         tempDir,
		Packages:        pkgs,
		ExtraFiles:      strings.Fields(*extraFiles),
		OutputFile:      f,
		BaseArchive:     baseFile,
		UseExistingInit: *useExistingInit,
	}
	if err := uroot.CreateInitramfs(opts); err != nil {
		return err
	}
	log.Printf("Filename is %s", filename)
	return nil
}
