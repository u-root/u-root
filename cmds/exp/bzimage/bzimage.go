// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bzImage is used to modify bzImage files.
// It reads the image in, applies an operator, and writes a new one out.
//
// Synopsis:
//     bzImage [copy <in> <out> ] | [diff <image> <image> ] | [dump <file>] | [initramfs input-bzimage initramfs output-bzimage]
//
// Description:
//	Read a bzImage in, change it, write it out, or print info.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/boot/bzimage"
)

var argcounts = map[string]int{
	"copy":      3,
	"diff":      3,
	"dump":      2,
	"initramfs": 4,
	"extract":   3,
	"ver":       2,
	"cfg":       2,
}

const cmdUsage = `Performs various operations on kernel images. Usage:
bzimage copy <in> <out>
	Create a copy of <in> at <out>, parsing structures.
bzimage diff <image> <image>
	Compare headers of two kernel images.
bzimage extract <file> <elf-file>
	extract parts of the kernel into separate files with self-
	explainatory extensions .boot, .head, .kern, .tail, .ramfs
bzimage dump <file>
    Dumps header.
bzimage initramfs <input-bzimage> <new-initramfs> <output-bzimage>
	Replaces initramfs in input-bzimage, creating output-bzimage.
bzimage ver <image>
	Dump version info similar to 'file <image>'.
bzimage cfg <image>
	Dump embedded config.

flags:`

var debug = flag.BoolP("debug", "d", false, "enable debug printing")
var jsonOut = flag.BoolP("json", "j", false, "json output ('ver' subcommand only)")

func usage() {
	fmt.Fprintln(os.Stderr, cmdUsage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *debug {
		bzimage.Debug = log.Printf
	}
	a := flag.Args()
	if len(a) < 2 {
		usage()
	}
	n, ok := argcounts[a[0]]
	if !ok || len(a) != n {
		usage()
	}

	var br = &bzimage.BzImage{}
	var image []byte
	switch a[0] {
	case "diff", "dump", "ver":
		br.NoDecompress = true
		fallthrough
	case "copy", "initramfs", "extract", "cfg":
		var err error
		image, err = ioutil.ReadFile(a[1])
		if err != nil {
			log.Fatal(err)
		}
		if err = br.UnmarshalBinary(image); err != nil {
			log.Fatal(err)
		}
	}

	switch a[0] {
	case "copy":
		o, err := br.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}
		if len(image) != len(o) {
			log.Printf("copy: input len is %d, output len is %d, they have to match", len(image), len(o))
			var br2 bzimage.BzImage
			if err = br2.UnmarshalBinary(o); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Input: %s\n", strings.Join(br.Header.Show(), "\n\t"))
			fmt.Printf("Output: %s\n", strings.Join(br2.Header.Show(), "\n\t"))
			log.Printf("%s", br.Header.Diff(&br2.Header))
			log.Fatalf("there is no hope")
		}
		if err := ioutil.WriteFile(a[2], o, 0666); err != nil {
			log.Fatalf("Writing %v: %v", a[2], err)
		}
	case "diff":
		b2, err := ioutil.ReadFile(a[2])
		if err != nil {
			log.Fatal(err)
		}
		var br2 = &bzimage.BzImage{}
		if err = br2.UnmarshalBinary(b2); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", br.Header.Diff(&br2.Header))
	case "dump":
		fmt.Printf("%s\n", strings.Join(br.Header.Show(), "\n"))
	case "extract":
		bzimage.Debug = log.Printf
		var i []byte
		s, e, err := br.InitRAMFS()
		if err != nil {
			fmt.Printf("Warning: could not extract initramfs: %v", err)
		} else {
			i = br.KernelCode[s:e]
		}
		// Need to add a trailer record to i
		fmt.Printf("ramfs is %d bytes", len(i))

		for _, v := range []struct {
			n string
			b []byte
		}{
			{a[2] + ".boot", br.BootCode},
			{a[2] + ".head", br.HeadCode},
			{a[2] + ".kern", br.KernelCode},
			{a[2] + ".tail", br.TailCode},
			{a[2] + ".ramfs", i},
		} {
			if v.b == nil {
				fmt.Printf("Warning: %s is nil", v.n)
				continue
			}
			if err := ioutil.WriteFile(v.n, v.b, 0666); err != nil {
				log.Fatalf("Writing %v: %v", v, err)
			}
		}
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
	case "ver":
		v, err := br.KVer()
		if err != nil {
			log.Fatal(err)
		}
		if *jsonOut {
			info, err := bzimage.ParseDesc(v)
			if err != nil {
				log.Fatal(err)
			}
			j, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(j))
		} else {
			fmt.Println(v)
		}
	case "cfg":
		cfg, err := br.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", cfg)
	}
}
