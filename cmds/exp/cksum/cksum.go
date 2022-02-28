// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
)

var (
	help    = pflag.BoolP("help", "h", false, "Show this help and exit")
	version = pflag.BoolP("version", "v", false, "Print Version")
)

func reverse(x uint32) uint32 {
	x = (((x & 0xaaaaaaaa) >> 1) | ((x & 0x55555555) << 1))
	x = (((x & 0xcccccccc) >> 2) | ((x & 0x33333333) << 2))
	x = (((x & 0xf0f0f0f0) >> 4) | ((x & 0x0f0f0f0f) << 4))
	x = (((x & 0xff00ff00) >> 8) | ((x & 0x00ff00ff) << 8))
	return ((x >> 16) | (x << 16))
}

func reverseByte(b byte) byte {
	b = (b&0xF0)>>4 | (b&0x0F)<<4
	b = (b&0xCC)>>2 | (b&0x33)<<2
	b = (b&0xAA)>>1 | (b&0x55)<<1
	return b
}

func reverseBytes(b []byte) []byte {
	for i := range b {
		b[i] = reverseByte(b[i])
	}
	return b
}

func helpPrinter(w io.Writer) error {
	fmt.Fprintf(w, "Usage:\ncksum <File Name>\n")
	pflag.PrintDefaults()
	return nil
}

func versionPrinter(w io.Writer) error {
	fmt.Fprintln(w, "cksum utility, URoot Version.")
	return nil
}

func getInput(r io.Reader, fileName string) (input []byte, err error) {
	if fileName != "" {
		return os.ReadFile(fileName)
	}
	return io.ReadAll(r)
}

func appendLengthToData(data []byte) []byte {
	length := len(data)
	for length > 0 {
		data = append(data, byte(length))
		length = (length >> 8)
	}
	return data
}

func calculateCksum(data []byte) uint32 {
	return reverse(crc32.Update(0xffffffff, crc32.MakeTable(0xEDB88320), reverseBytes(appendLengthToData(data))))
}

func cksum(w io.Writer, r io.Reader, args ...string) error {
	cliArgs := ""
	if *help {
		return helpPrinter(w)
	}
	if *version {
		return versionPrinter(w)
	}
	if len(args) >= 2 {
		cliArgs = args[1]
	}
	input, err := getInput(r, cliArgs)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, calculateCksum(input), len(input), cliArgs)
	return nil
}

func main() {
	pflag.Parse()
	if err := cksum(os.Stdout, os.Stdin, os.Args...); err != nil {
		log.Fatal(err)
	}
}
