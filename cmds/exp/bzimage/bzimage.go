// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bzImage is used to modify bzImage files.
// It reads the image in, applies an operator, and writes a new one out.
//
// Synopsis:
//
//	bzImage [copy <in> <out> ] | [diff <image> <image> ] | [dump <file>] | [initramfs input-bzimage initramfs output-bzimage]
//
// Description:
//
//	Read a bzImage in, change it, write it out, or print info.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/bzimage"
	"github.com/u-root/u-root/pkg/uroot/util"
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

const usage = `bzimage:
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

flags`

var (
	debug   = flag.Bool("d", false, "enable debug printing")
	jsonOut = flag.Bool("j", false, "json output ('ver' subcommand only)")
)

func run(w io.Writer, args ...string) error {
	if *debug {
		bzimage.Debug = log.Printf
	}
	if len(args) < 2 {
		flag.Usage()
		return nil
	}
	n, ok := argcounts[args[0]]
	if !ok || len(args) != n {
		flag.Usage()
		return nil
	}

	br := &bzimage.BzImage{}
	var image []byte
	switch args[0] {
	case "diff", "dump", "ver":
		br.NoDecompress = true
		fallthrough
	case "copy", "initramfs", "extract", "cfg":
		var err error
		image, err = os.ReadFile(args[1])
		if err != nil {
			return err
		}
		if err = br.UnmarshalBinary(image); err != nil {
			return err
		}
	}

	switch args[0] {
	case "copy":
		o, err := br.MarshalBinary()
		if err != nil {
			return err
		}
		if len(image) != len(o) {
			log.Printf("copy: input len is %d, output len is %d, they have to match", len(image), len(o))
			var br2 bzimage.BzImage
			if err = br2.UnmarshalBinary(o); err != nil {
				return err
			}
			fmt.Fprintf(w, "Input: %s\n", strings.Join(br.Header.Show(), "\n\t"))
			fmt.Fprintf(w, "Output: %s\n", strings.Join(br2.Header.Show(), "\n\t"))
			log.Printf("%s", br.Header.Diff(&br2.Header))
			return fmt.Errorf("there is no hope")
		}
		if err := os.WriteFile(args[2], o, 0o666); err != nil {
			return fmt.Errorf("writing %v: %w", args[2], err)
		}
	case "diff":
		b2, err := os.ReadFile(args[2])
		if err != nil {
			return err
		}
		br2 := &bzimage.BzImage{}
		if err = br2.UnmarshalBinary(b2); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s", br.Header.Diff(&br2.Header))
	case "dump":
		fmt.Fprintf(w, "%s\n", strings.Join(br.Header.Show(), "\n"))
	case "extract":
		bzimage.Debug = log.Printf
		var i []byte
		s, e, err := br.InitRAMFS()
		if err != nil {
			fmt.Fprintf(w, "Warning: could not extract initramfs: %v", err)
		} else {
			i = br.KernelCode[s:e]
		}
		// Need to add a trailer record to i
		fmt.Fprintf(w, "ramfs is %d bytes", len(i))

		for _, v := range []struct {
			n string
			b []byte
		}{
			{args[2] + ".boot", br.BootCode},
			{args[2] + ".head", br.HeadCode},
			{args[2] + ".kern", br.KernelCode},
			{args[2] + ".tail", br.TailCode},
			{args[2] + ".ramfs", i},
		} {
			if v.b == nil {
				fmt.Fprintf(w, "Warning: %s is nil", v.n)
				continue
			}
			if err := os.WriteFile(v.n, v.b, 0o666); err != nil {
				return fmt.Errorf("writing %v: %w", v, err)
			}
		}
	case "initramfs":
		if err := br.AddInitRAMFS(args[2]); err != nil {
			return err
		}

		b, err := br.MarshalBinary()
		if err != nil {
			return err
		}

		if err := os.WriteFile(args[3], b, 0o644); err != nil {
			return err
		}
	case "ver":
		v, err := br.KVer()
		if err != nil {
			return err
		}
		if *jsonOut {
			info, err := bzimage.ParseDesc(v)
			if err != nil {
				return err
			}
			j, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				return err
			}
			fmt.Fprintln(w, string(j))
		} else {
			fmt.Fprintln(w, v)
		}
	case "cfg":
		cfg, err := br.ReadConfig()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\n", cfg)
	}
	return nil
}

func main() {
	flag.Usage = util.Usage(flag.Usage, usage)
	flag.Parse()
	if err := run(os.Stdout, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
