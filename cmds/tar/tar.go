// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Create and extract tar archives.
//
// Synopsis:
//     tar [OPTION...] [FILE]...
//
// Description:
//     This command line can be used only in the following ways:
//        tar -cvf x.tar directory/         # create
//        tar -cvf x.tar file1 file2 ...    # create
//        tar -tvf x.tar                    # list
//        tar -xvf x.tar directory/         # extract
//
// Options:
//     -c: create a new tar archive from the given directory
//     -x: extract a tar archive to the given directory
//     -v: verbose, print each filename (optional)
//     -f: tar filename (required)
//     -t: list the contents of an archive
//
// TODO: The arguments deviates slightly from gnu tar.
package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/tarutil"
)

var (
	create  = flag.BoolP("create", "c", false, "create a new tar archive from the given directory")
	extract = flag.BoolP("extract", "x", false, "extract a tar archive from the given directory")
	verbose = flag.BoolP("verbose", "v", false, "print each filename")
	file    = flag.StringP("file", "f", "", "tar file")
	list    = flag.BoolP("list", "t", false, "list the contents of an archive")
)

func main() {
	flag.Parse()

	if *create && *extract {
		log.Fatal("cannot supply both -c and -x")
	} else if *create && *list {
		log.Fatal("cannot supply both -c and -t")
	} else if *extract && *list {
		log.Fatal("cannot supply both -x and -t")
	}

	if *file == "" {
		log.Fatal("tar filename is required")
	}

	var filters []tarutil.Filter
	if *verbose {
		filters = []tarutil.Filter{tarutil.VerboseFilter}
	}

	switch {
	case *create:
		f, err := os.Create(*file)
		if err != nil {
			log.Fatal(err)
		}
		if err := tarutil.CreateTarFilter(f, flag.Args(), filters); err != nil {
			f.Close()
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	case *extract:
		if flag.NArg() != 1 {
			flag.Usage()
			os.Exit(1)
		}
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := tarutil.ExtractDirFilter(f, flag.Arg(0), filters); err != nil {
			log.Fatal(err)
		}
	case *list:
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := tarutil.ListArchive(f); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("must supply at least one of: -c, -x, -t")
	}
}
