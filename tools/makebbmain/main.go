// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// makebbmain adds u-root command package imports to an existing main()
// template file.
package main

import (
	"flag"
	"go/build"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/uflag"
	"github.com/u-root/u-root/pkg/uroot"
)

var (
	pkg      = flag.String("template_pkg", "", "Go import package path")
	destDir  = flag.String("dest_dir", "", "Destination directory")
	pkgFiles uflag.Strings
	commands uflag.Strings
)

func init() {
	flag.Var(&pkgFiles, "package_file", "package files")
	flag.Var(&commands, "command", "Go package path for command to import")
}

func main() {
	flag.Parse()

	var dir string
	var gofiles []string
	for _, file := range pkgFiles {
		if len(dir) == 0 {
			dir = filepath.Dir(file)
		} else if dir != filepath.Dir(file) {
			log.Fatal("all package source files must be in the same directory")
		}
		gofiles = append(gofiles, filepath.Base(file))
	}

	bpkg := &build.Package{
		Name:       "main",
		Dir:        dir,
		ImportPath: *pkg,
		GoFiles:    gofiles,
	}

	fset, astp, err := uroot.ParseAST(bpkg)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*destDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := uroot.CreateBBMainSource(fset, astp, commands, *destDir); err != nil {
		log.Fatal(err)
	}
}
