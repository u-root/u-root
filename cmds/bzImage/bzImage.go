// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bzImage is used to modify bzImage files.
// It reads the image in, applies an operator, and writes a new one out.
//
// Synopsis:
//     bzImage [dump <file>] | [initramfs input-bzimage initramfs output-bzimage]
//
// Description:
//	Read a bzImage in, change it, write it out, or print info.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/bzimage"
)

var argcounts = map[string]int{
	"dump":      2,
	"initramfs": 4,
}

var cmdUsage = "Usage: bzImage  [dump <file>] | [initramfs input-bzimage initramfs output-bzimage]"

func usage() {
	log.Fatalf(cmdUsage)
}

func main() {
	flag.Parse()

	a := flag.Args()
	if len(a) < 2 {
		usage()
	}
	n, ok := argcounts[a[0]]
	if !ok || len(a) != n {
		usage()
	}

	var br = &bzimage.BzImage{}
	switch a[0] {
	case "dump", "initramfs":
		b, err := ioutil.ReadFile(a[1])
		if err != nil {
			log.Fatal(err)
		}
		if err = br.UnmarshalBinary(b); err != nil {
			log.Fatal(err)
		}
	}

	switch a[0] {
	case "dump":
		fmt.Printf("%s\n", strings.Join(br.Header.Show(), "\n"))
	case "initramfs":
		if err := br.AddInitRAMFS(a[2]); err != nil {
			log.Fatal(err)
		}

		b, err := br.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(a[3], b, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
