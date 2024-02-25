// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command vmmount mounts 9P directories as defined by env vars, runs a
// command, and unmounts them.
//
// The 9P directories are mounted via virtio; their tags are derived from any
// env var that matches VMTEST_MOUNT9P_*=$tag. The mount location is
// /mount/9p/$tag.
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hugelgupf/vmtest/guest"
)

func run() error {
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "VMTEST_MOUNT9P_") {
			continue
		}

		e := strings.SplitN(v, "=", 2)
		mp, err := guest.Mount9PDir(filepath.Join("/mount/9p", e[1]), e[1])
		if err != nil {
			log.Printf("Tried to mount 9P tag %s at /mount/9p/%s: %v", e[1], e[1], err)
		}
		defer func() {
			if err := mp.Unmount(0); err != nil {
				log.Printf("Failed to unmount: %v", err)
			}
		}()
	}

	args := flag.Args()
	if len(args) == 0 {
		return nil
	}
	c := exec.Command(args[0], args[1:]...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Printf("Failed: %v", err)
	}
}
