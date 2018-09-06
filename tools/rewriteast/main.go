// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rewrite_ast takes a Go command's source and rewrites it to be a u-root
// busybox compatible library package.
package main

import (
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/u-root/u-root/pkg/monoimporter"
	"github.com/u-root/u-root/pkg/uflag"
	"github.com/u-root/u-root/pkg/uroot"
)

var (
	name         = flag.String("name", "", "Name of the command")
	pkg          = flag.String("package", "", "Go import package path")
	destDir      = flag.String("dest_dir", "", "Destination directory")
	bbImportPath = flag.String("bb_import_path", "", "BB import path")
	gorootDir    uflag.Strings
	archives     uflag.Strings
	sourceFiles  uflag.Strings
)

func init() {
	flag.Var(&gorootDir, "go_root_zip", "Go standard library zip archives containing stdlib object files")
	flag.Var(&archives, "archive", "Archives")
	flag.Var(&sourceFiles, "source", "Source files")
}

func main() {
	flag.Parse()

	if len(*name) == 0 {
		log.Fatal("rewrite_ast: no command name given")
	}

	c := build.Default

	bpkg := &build.Package{
		Name:       "main",
		Dir:        filepath.Dir(sourceFiles[0]),
		ImportPath: *pkg,
	}

	for _, path := range sourceFiles {
		basename := filepath.Base(path)
		// Check the file against build tags.
		//
		// TODO: build.Default may not be the actual build environment.
		// Fix it via flags from Skylark?
		ok, err := c.MatchFile(bpkg.Dir, basename)
		if ok {
			bpkg.GoFiles = append(bpkg.GoFiles, basename)
		} else if err != nil {
			log.Fatal(err)
		} else {
			b, err := ioutil.ReadFile(filepath.Join(bpkg.Dir, basename))
			if err != nil {
				log.Fatal(err)
			}
			// N.B. Hack: Blaze expects an output file for every
			// input file, even if we decide to do nothing with the
			// input file.  Just write a damn empty file. The
			// compiler will automagically exclude it based on the
			// same build tags.
			if err := ioutil.WriteFile(filepath.Join(*destDir, basename), b, 0644); err != nil {
				log.Fatal(err)
			}
		}
	}

	imp, err := monoimporter.NewFromZips(c, []string(archives), []string(gorootDir))
	if err != nil {
		log.Fatal(err)
	}

	bbPkg, err := uroot.NewPackage(*name, bpkg, imp)
	if err != nil {
		log.Fatal(err)
	}
	if err := bbPkg.Rewrite(*destDir, *bbImportPath); err != nil {
		log.Fatal(err)
	}
}
