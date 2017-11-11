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
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/ldd"
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

	env := golang.Default()
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)

	if err := Build(env, flag.Args(), *build, *format, *tmpDir, *base, *outputPath, strings.Fields(*extraFiles), *useExistingInit); err != nil {
		log.Fatalf("u-root: %v", err)
	}
}

func Build(env golang.Environ, pkgs []string, build, format, tempDir, base, outputPath string, extraFiles []string, useExistingInit bool) error {
	builder, err := uroot.GetBuilder(build)
	if err != nil {
		return err
	}
	archiver, err := uroot.GetArchiver(format)
	if err != nil {
		return err
	}

	if tempDir == "" {
		var err error
		tempDir, err = ioutil.TempDir("", "u-root")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
	}

	// Resolve globs into package imports.
	//
	// Currently allowed formats:
	//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/*
	if len(pkgs) == 0 {
		var err error
		pkgs, err = uroot.DefaultPackageImports(env)
		if err != nil {
			return err
		}
	}

	var importPaths []string
	// Resolve file system paths to package import paths.
	for _, pkg := range pkgs {
		importPath, err := env.FindPackageByPath(pkg)
		if err != nil {
			if _, perr := env.FindPackageDir(pkg); perr != nil {
				return fmt.Errorf("%q is neither package or path: %v / %v", pkg, err, perr)
			}
			importPath = pkg
		}
		importPaths = append(importPaths, importPath)
	}

	// Build the packages.
	bOpts := uroot.BuildOpts{
		Env:      env,
		Packages: importPaths,
		TempDir:  tempDir,
	}
	files, err := builder(bOpts)
	if err != nil {
		return fmt.Errorf("error building %#v: %v", bOpts, err)
	}

	// Open the target initramfs file.
	if outputPath == "" {
		outputPath = fmt.Sprintf("/tmp/initramfs.%s_%s.%s", env.GOOS, env.GOARCH, archiver.DefaultExtension())
	}
	f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	var baseFile *os.File
	if base != "" {
		var err error
		baseFile, err = os.Open(base)
		if err != nil {
			return err
		}
		defer baseFile.Close()
	}

	archive := uroot.ArchiveOpts{
		ArchiveFiles:    files,
		OutputFile:      f,
		BaseArchive:     baseFile,
		UseExistingInit: useExistingInit,
	}

	// Add files from command line.
	for _, file := range extraFiles {
		path, err := filepath.Abs(file)
		if err != nil {
			return fmt.Errorf("couldn't find absolute path for %q: %v", file, err)
		}
		if err := archive.AddFile(path, path[1:]); err != nil {
			return fmt.Errorf("couldn't add %q to archive: %v", file, err)
		}

		// Pull dependencies in the case of binaries. If `path` is not
		// a binary, `libs` will just be empty.
		libs, err := ldd.List([]string{path})
		if err != nil {
			return fmt.Errorf("couldn't list ldd dependencies for %q: %v", file, err)
		}
		for _, lib := range libs {
			if err := archive.AddFile(lib, lib[1:]); err != nil {
				return fmt.Errorf("couldn't add %q to archive: %v", lib, err)
			}
		}
	}

	// Finally, write the archive.
	if err := archiver.Archive(archive); err != nil {
		return fmt.Errorf("error archiving: %v", err)
	}

	log.Printf("Filename is %s", outputPath)
	return nil
}
