// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blkid prints information about blocks.
package main

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/mount/block"
)

func main() {
	devices, err := block.GetBlockDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Print(device.DevicePath())
		if device.FsUUID != "" {
			fmt.Printf(` UUID="%s"`, device.FsUUID)
		}
		if device.FSType != "" {
			fmt.Printf(` TYPE="%s"`, device.FSType)
		}
		fmt.Println()
	}
}
