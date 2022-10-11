// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// partprobe prompts the OS to re-read partition tables.
//
// Synopsis:
//
//	partprobe [device]...
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/mount/block"
)

func main() {
	flag.Parse()

	devs := flag.Args()

	if len(devs) == 0 {
		log.Printf("Usage: partprobe [device]...")
		os.Exit(0)
	}

	for _, dev := range devs {
		d, err := block.Device(dev)
		if err != nil {
			log.Printf("Failed to find device %s: %v", dev, err)
			continue
		}

		if err := d.ReadPartitionTable(); err != nil {
			log.Printf("Failed to read partition table for %s: %v", dev, err)
		}
	}
}
