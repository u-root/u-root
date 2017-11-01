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
	// Paths contains the paths to put into the initramfs.
	// The index is a root directory, and the value is the place from which to
	// walk. The only required root is the bbsh dir itself, and the starting
	// walk is init -- i.e. we grab only one file. Should you wish to bring in,
	// e.g., /lib/modules/4.04, you would do add the root as / and the
	// starting point for the walk as lib/modules/4.04. That way we only preserve
	// as much of the path as we need, but we can preserve it all.
	paths     = map[string][]string{}
	extraAdd  = flag.String("add", "", "Extra commands or directories to add (full path, space-separated string)")
	extraCpio = flag.String("cpio", "", "A list of cpio archives to include in the output")
)

func sanity() {
	goBinGo := filepath.Join(config.Goroot, "bin/go")
	_, err := os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	// But does the one in go/bin/OS_ARCH exist too?
	goBinGo = filepath.Join(config.Goroot, fmt.Sprintf("bin/%s_%s/go", config.Goos, config.Arch))
	_, err = os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	if config.Go == "" {
		log.Fatalf("Can't find a go binary! Is GOROOT set correctly?")
	}
}

// dirComponents takes a string and returns an array of strings,
// such that we can create the directory records for intermediate
// directories.
func dirComponents(dir string) []string {
	var dirlist []string
	if filepath.Dir(dir) == "." {
		return []string{}
	}
	for d := filepath.Dir(dir); d != "/"; d = filepath.Dir(d) {
		dirlist = append([]string{d}, dirlist...)
	}
	dirlist = append(dirlist, dir)
	return dirlist
}

func initramfs(goos string, arch string) {
	archiver, err := cpio.Format("newc")
	if err != nil {
		log.Fatalf("Creating newc archiver: %v", err)
	}

	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", goos, arch)
	f, err := os.Create(oname)
	if err != nil {
		log.Fatalf("%v", err)
	}

	w := archiver.Writer(f)
	if err := w.WriteRecords(ramfs.DevCPIO); err != nil {
		log.Fatalf("%v", err)
	}

	paths[filepath.Join(config.Gopath, "src/github.com/u-root/u-root/bb/bbsh")] = []string{"init", "ubin"}

	if *extraAdd != "" {
		copyc := strings.Fields(*extraAdd)
		for _, eachPath := range copyc {
			debug("each path is %s\n", eachPath)
			// Must not include ~ in path
			if strings.Count(eachPath, ":") > 1 {
				log.Fatalf(" Input has more than one :")
			}

			modPath := strings.Replace(eachPath, ":", "/", 1)
			debug("modpath is %s\n", modPath)
			statval, err := os.Stat(modPath)
			if err != nil {
				log.Fatalf("%v", err)
			}

			p := strings.Split(eachPath, ":")
			if len(p) != 2 {
				p = append([]string{"/"}, p...)
			}
			debug("P is %v\n", p)
			debug("Paths currently is %v\n", paths)
			paths[p[0]] = append(paths[p[0]], p[1])
			debug("putps")
			debug("Paths currently is %v\n", paths)

			if !statval.IsDir() {
				tmpSlice := []string{modPath}
				libs, err := ldd.List(tmpSlice)
				if err != nil {
					log.Fatalf("%v", err)
				}
				paths["/"] = append(paths["/"], libs...)
				debug("Paths currently is (because command) %v\n", paths)
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

			recs, err := archiver.Reader(f).ReadRecords()
			if err != nil {
				log.Fatalf("read records: %v", err)
			}

			cpio.MakeAllReproducible(recs)
			if err := w.WriteRecords(recs); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
	//append the ldd'd files to /
	// For all the 'roots' in paths, start walking at the name.
	debug("PATHS: %v", paths)
	for r, list := range paths {
		debug("PATHS: root %v", r)
		// We need to make all the path prefix directories.
		for _, n := range list {
			debug("\troot %v, name %v", r, n)
			for _, d := range dirComponents(n) {
				debug("\t\troot %v, name %v, component %v", r, n, d)
				rec, err := cpio.GetRecord(d)
				if err != nil {
					log.Fatalf("Getting record of %q failed: %v", d, err)
				}

				if err := w.WriteRecord(cpio.MakeReproducible(rec)); err != nil {
					log.Fatalf("%v", err)
				}
			}
		}

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
				if err := w.WriteRecord(cpio.MakeReproducible(rec)); err != nil {
					log.Fatalf("%v", err)
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
