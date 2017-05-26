// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bbramfs builds a simple initramfs given an existing built bb; see bb.go
// You have to run bb first, which creates cmds/bb/bbsh. cd to that directory,
// and run bbramfs, and you have a single binary which does all u-root commands.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
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
	if err := w.WriteRecords(devCPIO[:]); err != nil {
		log.Fatalf("%v\n", err)
	}

	bbdir := filepath.Join(config.Gopath, "src/github.com/u-root/u-root/bb/bbsh")

	if err := filepath.Walk(filepath.Join(bbdir, "init"), func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cn, err := filepath.Rel(bbdir, name)
		if err != nil {
			log.Fatalf("filepath.Rel(%v, %v): %v", bbdir, name, err)
		}
		debug("%v\n", cn)
		rec, err := cpio.GetRecord(name)
		if err != nil {
			log.Fatalf("Getting record of %q failed: %v", cn, err)
		}
		// the name in the cpio is relative to our starting point.
		rec.Name = cn
		if err := w.WriteRecord(rec); err != nil {
			log.Fatalf("%v\n", err)
		}
		return nil
	}); err != nil {
		log.Fatalf("bbsh walk failed: %v", err)
	}

	if err := w.WriteTrailer(); err != nil {
		log.Fatalf("Error writing trailer record: %v", err)
	}
	fmt.Printf("Output file is in %v\n", oname)
}
