// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

// Command bootvars reads the current UEFI boot variables. Given -m arg, it
// mounts the filesystem pointed to by the current boot variable and prints the
// location.
//
// Note - does not check whether the fs is already mounted anywhere. User is
// responsible for unmounting and for removing the temp dir.
package main

import (
	"flag"
	"log"
	"os"
	fp "path/filepath"

	"github.com/u-root/u-root/pkg/uefivars/boot"
)

// must run as root, as efi vars are not accessible otherwise
func main() {
	m := flag.Bool("m", false, "Mount FS containing boot file, print path.")
	flag.Parse()

	bv, err := boot.ReadCurrentBootVar()
	if err != nil {
		log.Fatalf("Reading current boot var: %s", err)
	}
	if bv == nil {
		log.Fatalf("Unable to read var... are you root?")
	}
	if *m {
		// actually mount FS, locate file
		exe := fp.Base(os.Args[0])
		tmp, err := os.MkdirTemp("", exe)
		if err != nil {
			log.Fatalf("creating temp dir: %s", err)
		}
		path := tmp
		for _, element := range bv.FilePathList {
			res, err := element.Resolver()
			if err != nil {
				log.Fatalf("%s", err)
			}
			path, err = res.Resolve(path)
			if err != nil {
				log.Printf("Resolving element %s: %s", element, err)
			}
		}
		log.Printf("file corresponding to CurrentBoot var can be found at %s\n", path)
		log.Printf("you will need to unmount the filesystem, remove temp dir, etc when done.\n")
	} else {
		// just print out elements in variable
		log.Printf("%s", bv)
		for _, element := range bv.FilePathList {
			res, err := element.Resolver()
			if err != nil {
				log.Fatalf("%s", err)
			}
			log.Printf("%s", res.String())
		}
	}
}
