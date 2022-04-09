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

func run(getBlock func() (block.BlockDevices, error), out io.ReadWriter) error {
	devices, err := getBlock()
	if err != nil {
		return fmt.Errorf("error getting Block devices: %v", err)
	}

	for _, device := range devices {
		fmt.Print(device.DevicePath())
		if device.FsUUID != "" {
			fmt.Fprintf(out, ` UUID="%s"`, device.FsUUID)
		}
		if device.FSType != "" {
			fmt.Fprintf(out, ` TYPE="%s"`, device.FSType)
		}
		fmt.Println()
	}
	return nil
}

func main() {
	if err := run(block.GetBlockDevices, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
