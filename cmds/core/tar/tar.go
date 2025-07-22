// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Create and extract tar archives.
//
// Synopsis:
//
//	tar [OPTION...] [FILE]...
//
// Description:
//
//	This command line can be used only in the following ways:
//	   tar -cvf x.tar directory/         # create
//	   tar -cvf x.tar file1 file2 ...    # create
//	   tar -tvf x.tar                    # list
//	   tar -xvf x.tar directory/         # extract
//
// Options:
//
//	-c: create a new tar archive from the given directory
//	-x: extract a tar archive to the given directory
//	-v: verbose, print each filename (optional)
//	-f: tar filename (required)
//	-t: list the contents of an archive
//
// TODO: The arguments deviates slightly from gnu tar.
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/tar"
)

func main() {
	cmd := tar.New()
	if err := cmd.Run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
