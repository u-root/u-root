// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bbramfs builds a simple initramfs given an existing built bb; see bb.go
// You have to run bb first, which creates cmds/bb/bbsh. cd to that directory,
// and run bbramfs, and you have a single binary which does all u-root commands.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/ramfs"
)

var (
	extraAdd  = flag.String("add", "", "Extra commands or directories to add (full path, space-separated string)")
	extraCpio = flag.String("cpio", "", "A list of cpio archives to include in the output")
)

func initramfs(goos string, arch string) error {
	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", goos, arch)
	f, err := os.Create(oname)
	if err != nil {
		return err
	}
	defer f.Close()

	archiver, err := cpio.Format("newc")
	if err != nil {
		return err
	}

	init, err := ramfs.NewInitramfs(archiver.Writer(f))
	if err != nil {
		return err
	}

	// Paths contains the paths to put into the initramfs. The index is a
	// root directory, and the value is the place from which to walk.
	//
	// The only required root is the bbsh dir itself, and the starting walk
	// is init -- i.e. we grab only one file. Should you wish to bring in,
	// e.g., /lib/modules/4.04, you would do add the root as / and the
	// starting point for the walk as lib/modules/4.04. That way we only
	// preserve as much of the path as we need, but we can preserve it all.
	paths := make(map[string][]string)
	paths[filepath.Join(config.Gopath, "src/github.com/u-root/u-root/bb/bbsh")] = []string{"init", "ubin"}

	if *extraAdd != "" {
		copyc := strings.Fields(*extraAdd)
		for _, path := range copyc {
			s, err := os.Stat(path)
			if err != nil {
				return err
			}

			paths["/"] = append(paths["/"], path)
			if s.Mode().IsRegular() && (s.Mode()&0111 != 0) {
				libs, err := ldd.List([]string{path})
				if err != nil {
					return err
				}
				paths["/"] = append(paths["/"], libs...)
			}
		}
	}

	if *extraCpio != "" {
		extras := strings.Fields(*extraCpio)
		for _, x := range extras {
			f, err := os.Open(x)
			if err != nil {
				log.Fatalf("%v: %v", x, err)
			}
			defer f.Close()

			if err := init.Concat(archiver.Reader(f), cpio.MakeReproducible); err != nil {
				return err
			}
		}
	}

	for root, list := range paths {
		if err := init.WriteFiles(root, "", list); err != nil {
			return err
		}
	}

	if err := init.WriteTrailer(); err != nil {
		return err
	}

	fmt.Printf("Output file is in %v\n", oname)
	return nil
}
