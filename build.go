// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Synopsis:
//     u-root [OPTIONS] [GLOBS...]
//
// Options:
//     -build_format: one of src or bb (default "src")
//     -format:       one of chroot, cpio, docker or list (default "chroot")
//     -no_def_glob:  disable the default glob ("github.com/u-root/u-root/cmds/*")
//     -o             output file or directory
//     -run:          run the generated ramfs
//     -v:            verbose
//
// Bugs:
//     installcommand does not recognize main packages outside of
//     "github.com/u-root/u-root/cmds".
//
//     initialCpio and existingInit are not yet implemented.
//
//     docker and bb are not yet implemented.
package main

import (
	"flag"
	gobuild "go/build"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/u-root/u-root/scripts/build"
)

var (
	archiveFormat = flag.String("format", "chroot", "one of chroot, cpio, docker or list")
	buildFormat   = flag.String("build_format", "src", "one of src or bb")
	noDefGlob     = flag.Bool("no_def_glob", false, "disable the default glob (\"github.com/u-root/u-root/cmds/*\")")
	output        = flag.String("o", "", "output file or directory")
	run           = flag.Bool("run", false, "run the generated ramfs")
	verbose       = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()

	// Perform a sanity check on the environment.
	if *verbose {
		log.Printf("env: %#v\n", gobuild.Default)
	}
	if gobuild.Default.GOPATH == "" {
		log.Println("warning: GOPATH is unset, upgrade to Go 1.8")
	}
	if gobuild.Default.GOROOT == "" {
		log.Println("warning: GOROOT is unset")
	}
	if gobuild.Default.CgoEnabled {
		log.Println("warning: CGO is enabled")
		if err := os.Setenv("CGO_ENABLED", "0"); err == nil {
			gobuild.Default.CgoEnabled = false
			log.Println("warning: disabling CGO on your behalf")
		}
	}

	// De-glob package names.
	globs := flag.Args()
	if !*noDefGlob {
		globs = append(globs, "github.com/u-root/u-root/cmds/*") // default
	}
	packages := []string{}
	for _, g := range globs {
		// Expand glob in GOROOT and all GOPATHs.
		for _, srcDir := range gobuild.Default.SrcDirs() {
			pkgs, err := filepath.Glob(path.Join(srcDir, g))
			if err != nil {
				log.Fatalf("%q is an invalid glob: %v", g, err)
			}
			for _, absPkg := range pkgs {
				// Exclude non-directories.
				if fi, err := os.Stat(absPkg); err == nil && fi.IsDir() {
					// Convert path to be relative to GOPATH or GOROOT and canonicalize.
					relPkg, err := filepath.Rel(srcDir, absPkg)
					if err != nil {
						log.Fatalf("Package %q not in GOPATH or GOROOT: %v", absPkg, err)
					}
					packages = append(packages, relPkg)
				}
			}
		}
	}
	packages = build.Uniq(packages)

	// Create a temporary output path if one does not exist.
	outputPath := *output
	if outputPath == "" {
		var err error
		outputPath, err = ioutil.TempDir("", "uroot")
		if err != nil {
			log.Fatal(err)
		}
		outputPath = filepath.Join(outputPath, "output")
	}

	// Setup config.
	config := build.Config{
		ArchiveFormat: *archiveFormat,
		BuildFormats:  []string{*buildFormat, "dev"}, // always include dev generator
		OutputPath:    outputPath,
		Packages:      packages,
		Run:           *run,
		Verbose:       *verbose,
	}
	if *verbose {
		log.Printf("config: %+v\n", config)
	}

	// Build and optionally run.
	if err := build.Build(config); err != nil {
		log.Fatalln("fatal:", err)
	}
}
