// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
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

func helpPrinter() {
	fmt.Printf("Usage:\ncksum <File Name>\n")
	pflag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("cksum utility, URoot Version.")
	os.Exit(0)
}

func getInput(fileName string) (input []byte, err error) {
	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
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

func main() {
	var (
		help    bool
		version bool
	)
	cliArgs := ""
	pflag.BoolVarP(&help, "help", "h", false, "Show this help and exit")
	pflag.BoolVarP(&version, "version", "v", false, "Print Version")
	pflag.Parse()

	if help {
		helpPrinter()
	}

	if version {
		versionPrinter()
	}
	if len(os.Args) >= 2 {
		cliArgs = os.Args[1]
	}
	input, err := getInput(cliArgs)
	if err != nil {
		return
	}
	fmt.Println(calculateCksum(input), len(input), cliArgs)
}
