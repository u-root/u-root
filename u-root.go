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

	extraFiles = flag.String("files", "", "Additional files and directories to add to archive.")
	binaries   = flag.String("binaries", "", "Additional binaries and their ldd dependencies to add to archive.")

	outputPath = flag.String("o", "", "Path to output initramfs file.")
)

func main() {
	flag.Parse()

	env := golang.Default()
	if env.CgoEnabled {
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)

	builder, err := uroot.GetBuilder(*build)
	if err != nil {
		log.Fatalf("%v", err)
	}
	archiver, err := uroot.GetArchiver(*format)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tempDir := *tmpDir
	if tempDir == "" {
		var err error
		tempDir, err = ioutil.TempDir("", "u-root")
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer os.RemoveAll(tempDir)
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
			log.Fatalf("%v", err)
		}
	}

	var importPaths []string
	// Resolve file system paths to package import paths.
	for _, pkg := range pkgs {
		importPath, err := env.FindPackageByPath(pkg)
		if err != nil {
			if _, perr := env.FindPackageDir(pkg); perr != nil {
				log.Fatalf("%q is neither package or path: %v / %v", pkg, err, perr)
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
		log.Fatalf("Error building %#v: %v", bOpts, err)
	}

	// Open the target initramfs file.
	filename := *outputPath
	if filename == "" {
		filename = fmt.Sprintf("/tmp/initramfs.%s_%s.%s", env.GOOS, env.GOARCH, archiver.DefaultExtension())
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalf("Couldn't open file %q: %v", filename, err)
	}
	defer f.Close()

	var baseFile *os.File
	if *base != "" {
		var err error
		baseFile, err = os.Open(*base)
		if err != nil {
			log.Fatalf("Couldn't open %q: %v", *base, err)
		}
		defer baseFile.Close()
	}

	archive := uroot.ArchiveOpts{
		ArchiveFiles:    files,
		OutputFile:      f,
		BaseArchive:     baseFile,
		UseExistingInit: *useExistingInit,
	}

	// Add files from command line.
	for _, file := range strings.Fields(*extraFiles) {
		path, err := filepath.Abs(file)
		if err != nil {
			log.Fatalf("Couldn't find absolute path for %q: %v", file, err)
		}
		if err := archive.AddFile(path, path[1:]); err != nil {
			log.Fatalf("Couldn't add %q to archive: %v", file, err)
		}
	}

	// Add binaries from command line.
	for _, binary := range strings.Fields(*binaries) {
		path, err := filepath.Abs(binary)
		if err != nil {
			log.Fatalf("Couldn't find absolute path for %q: %v", binary, err)
		}
		if err := archive.AddFile(path, path[1:]); err != nil {
			log.Fatalf("Couldn't add %q to archive: %v", binary, err)
		}

		libs, err := ldd.List([]string{path})
		if err != nil {
			log.Fatalf("Couldn't list dependencies for %q: %v", binary, err)
		}
		for _, lib := range libs {
			if err := archive.AddFile(lib, lib[1:]); err != nil {
				log.Fatalf("Couldn't add %q to archive: %v", lib, err)
			}
		}
	}

	// Finally, write the archive.
	if err := archiver.Archive(archive); err != nil {
		log.Fatalf("Error archiving: %v", err)
	}

	log.Printf("Filename is %s", filename)
}
