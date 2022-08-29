// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// md5sum prints an md5 hash generated from file contents.
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/uroot/util"
)

var usage = "md5sum: md5sum <File Name>"

func init() {
	flag.Usage = util.Usage(flag.Usage, usage)
}

func calculateMd5Sum(r io.Reader) ([]byte, error) {
	md5Generator := md5.New()
	if _, err := io.Copy(md5Generator, r); err != nil {
		return nil, err
	}
	return md5Generator.Sum(nil), nil
}

func md5Sum(w io.Writer, r io.Reader, args ...string) error {
	var err error

	if len(args) == 0 {
		h, err := calculateMd5Sum(r)
		if err != nil {
			fmt.Println("Error getting input.")
			return err
		}
		_, err = fmt.Fprintf(w, "%x\n", h)
		if err != nil {
			return err
		}
	} else {
		fileDesc, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer fileDesc.Close()
		h, err := calculateMd5Sum(fileDesc)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%x %s\n", h, args[0])
		if err != nil {
			return err
		}
	}
	return err
}

func main() {
	flag.Parse()
	if err := md5Sum(os.Stdout, os.Stdin, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
