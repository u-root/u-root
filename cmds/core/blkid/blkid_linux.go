// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blkid prints information about blocks.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/mount/block"
)

func run(getBlock func() (block.BlockDevices, error), out io.Writer) error {
	devices, err := getBlock()
	if err != nil {
		return fmt.Errorf("error getting Block devices: %w", err)
	}

	for _, device := range devices {
		fmt.Fprint(out, device.DevicePath())
		if device.FsUUID != "" {
			fmt.Fprintf(out, " UUID=%q", device.FsUUID)
		}
		if device.FSType != "" {
			fmt.Fprintf(out, " TYPE=%q", device.FSType)
		}
		fmt.Fprintln(out)
	}
	return nil
}

func main() {
	if err := run(block.GetBlockDevices, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
