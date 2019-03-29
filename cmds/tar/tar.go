// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Create and extract tar archives.
//
// Synopsis:
//     tar [-cxvf] DIRECTORY
//
// Description:
//     There are only two ways to use this command line:
//        tar -cvf x.tar directory/  # create
//        tar -xvf x.tar directory/  # extract
//
// Options:
//     -c: create a new tar archive from the given directory
//     -x: extract a tar archive to the given directory
//     -v: verbose, print each filename (optional)
//     -f: tar filename (required)
//
// TODO: The arguments deviates slightly from gnu tar.
package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/tar"
)

var (
	create  = flag.BoolP("create", "c", false, "create a new tar archive from the given directory")
	extract = flag.BoolP("extract", "x", false, "extract a tar archive from the given directory")
	verbose = flag.BoolP("verbose", "v", false, "print each filename")
	file    = flag.StringP("file", "f", "", "tar file")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if !*create && !*extract {
		log.Fatal("must supply at least one of -c or -x")
	}
	if *create && *extract {
		log.Fatal("cannot supply both -c and -x")
	}

	if *file == "" {
		log.Fatal("tar filename is required")
	}

	filter := tar.NoFilter
	if *verbose {
		filter = tar.VerboseFilter
	}

	if *create {
		f, err := os.Create(*file)
		if err != nil {
			log.Fatal(err)
		}
		if err := tar.CreateDirFilter(f, flag.Arg(0), filter); err != nil {
			f.Close()
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := tar.ExtractDirFilter(f, flag.Arg(0), filter); err != nil {
			log.Fatal(err)
		}
	}
}
