// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// losetup sets up and controls loop devices.
//
// Synopsis:
//
//	losetup [-Ad] FILE
//	losetup [-Ad] DEV FILE
//
// Options:
//
//	-A: pick any device
//	-d: detach the device
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/mount/loop"
)

var detach = flag.Bool("d", false, "Detach the device")

func main() {
	var (
		filename, devicename string
		err                  error
	)

	flag.Parse()
	args := flag.Args()
	if *detach {
		if len(args) == 1 {
			if err := loop.ClearFile(args[0]); err != nil {
				log.Fatal("Error clearing device: ", err)
			}
			log.Println("Detached", args[0])
			os.Exit(0)
		}
		flag.Usage()
		log.Fatal("Syntax Error")
	}

	if len(args) == 1 {
		devicename, err = loop.FindDevice()
		if err != nil {
			log.Fatalf("can't find a loop: %v", err)
		}
		filename = args[0]
	} else if len(args) == 2 {
		devicename = args[0]
		filename = args[1]
	} else {
		flag.Usage()
		log.Fatal("Syntax Error")
	}

	if err := loop.SetFile(devicename, filename); err != nil {
		log.Fatal("Could not set loop device:", err)
	}

	log.Printf("Attached %s to %s", devicename, filename)
}
