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
)

var (
	// Paths contains the paths to put into the initramfs.
	// The index is a root directory, and the value is the place from which to
	// walk. The only required root is the bbsh dir itself, and the starting
	// walk is init -- i.e. we grab only one file. Should you wish to bring in,
	// e.g., /lib/modules/4.04, you would do add the root as / and the
	// starting point for the walk as lib/modules/4.04. That way we only preserve
	// as much of the path as we need, but we can preserve it all.
	paths      = map[string][]string{}
	extraPaths = flag.String("extra", "", "Extra paths to add in the form root:start, e.g. /:etc/hosts")
	extraCmds  = flag.String("cmds", "", "Extra commands to add (full path, comma-separated string)")
)

func sanity() {
	goBinGo := filepath.Join(config.Goroot, "bin/go")
	_, err := os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	// but does the one in go/bin/OS_ARCH exist too?
	goBinGo = filepath.Join(config.Goroot, fmt.Sprintf("bin/%s_%s/go", config.Goos, config.Arch))
	_, err = os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	if config.Go == "" {
		log.Fatalf("Can't find a go binary! Is GOROOT set correctly?")
	}
}

// copyCommands takes a list of commands, generates the list of libs,
// and creates cpio records, including directory records.
func copyCommands(w cpio.Writer, cmds []string) {
	debug("copyCommands: start with %v", cmds)
	var recs []cpio.Record
	libs, err := uroot.LddList(cmds)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cmds = append(cmds, libs...)
	for _, n := range cmds {
		debug("copyCommands: file %v", n)
		var dirlist []string
		for d := filepath.Dir(n); d != "/"; d = filepath.Dir(d) {
			dirlist = append([]string{d}, dirlist...)
		}
		dirlist = append(dirlist, n)
		for _, n := range dirlist {
			debug("copyCommands: %v", n)
			r, err := cpio.GetRecord(n)
			if err != nil {
				log.Fatalf("%v: %v", n, err)
			}
			recs = append(recs, r)
		}
	}
	cpio.MakeReproducible(recs)
	if err := w.WriteRecords(recs); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func ramfs() {
	archiver, err := cpio.Format("newc")
	if err != nil {
		log.Fatalf("Creating newc archiver: %v", err)
	}

	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", config.Goos, config.Arch)
	f, err := os.Create(oname)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	w := archiver.Writer(f)
	cpio.MakeReproducible(devCPIO[:])
	if err := w.WriteRecords(devCPIO[:]); err != nil {
		log.Fatalf("%v\n", err)
	}

	paths[filepath.Join(config.Gopath, "src/github.com/u-root/u-root/bb/bbsh")] = []string{"init", "ubin"}

	if *extraPaths != "" {
		extras := strings.Split(*extraPaths, " ")
		for _, x := range extras {
			paths["/"] = append(paths["/"], x)
		}
	}

	if *extraCmds != "" {
		copyCommands(w, strings.Split(*extraCmds, " "))
	}

	// For all the 'roots' in paths, start walking at the name.
	debug("PATHS: %v", paths)
	for r, list := range paths {
		for _, n := range list {
			if err := filepath.Walk(filepath.Join(r, n), func(name string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				cn, err := filepath.Rel(r, name)
				if err != nil {
					log.Fatalf("filepath.Rel(%v, %v): %v", r, name, err)
				}
				debug("%v\n", cn)
				rec, err := cpio.GetRecord(name)
				if err != nil {
					log.Fatalf("Getting record of %q failed: %v", cn, err)
				}
				// the name in the cpio is relative to our starting point.
				rec.Name = cn
				recs := []cpio.Record{rec}
				cpio.MakeReproducible(recs)
				if err := w.WriteRecords(recs); err != nil {
					log.Fatalf("%v\n", err)
				}
				return nil
			}); err != nil {
				log.Fatalf("bbsh walk failed: %v", err)
			}
		}
	}

	if err := w.WriteTrailer(); err != nil {
		log.Fatalf("Error writing trailer record: %v", err)
	}
	fmt.Printf("Output file is in %v\n", oname)
}
