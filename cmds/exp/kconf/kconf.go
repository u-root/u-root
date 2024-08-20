// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/bzimage"
)

const (
	cfgfile = "/proc/config.gz"
	notset  = "is not set"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Reads kernel config from /proc/config.gz or bzimage, optionally filtering\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-k /path/to/bzimage] [-y|-m|-n] [-f filterstr]\n", os.Args[0])
		flag.PrintDefaults()
	}
	y := flag.Bool("y", false, "show only built-in")
	m := flag.Bool("m", false, "show only modules")
	n := flag.Bool("n", false, "show only not configured")
	f := flag.String("f", "", "filter on config symbol, case insensitive")
	p := flag.Bool("p", true, "pretty: trim prefix; also suffix if -y/-m/-n")
	k := flag.String("k", "", "use kernel image rather than /proc/config.gz")
	flag.Parse()

	var cfgIn io.Reader
	if len(*k) > 0 {
		image, err := os.ReadFile(*k)
		if err != nil {
			log.Fatal(err)
		}
		br := &bzimage.BzImage{}
		if err = br.UnmarshalBinary(image); err != nil {
			log.Fatal(err)
		}
		cfg, err := br.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}
		cfgIn = strings.NewReader(cfg)
	} else {
		configgz, err := os.Open(cfgfile)
		if err != nil {
			log.Fatalf("cannot open %s: %s", cfgfile, err)
		}
		defer configgz.Close()
		gz, err := gzip.NewReader(configgz)
		if err != nil {
			log.Fatalf("decompress %s: %s", cfgfile, err)
		}
		defer gz.Close()
		cfgIn = gz
	}
	filter := strings.ToUpper(*f)
	scanner := bufio.NewScanner(cfgIn)
	for scanner.Scan() {
		line := scanner.Text()
		if *y && !strings.HasSuffix(line, "=y") {
			continue
		}
		if *m && !strings.HasSuffix(line, "=m") {
			continue
		}
		if *n && !strings.HasSuffix(line, notset) {
			continue
		}
		if len(filter) > 0 && !strings.Contains(line, filter) {
			continue
		}
		if *p {
			if *n {
				line = strings.TrimPrefix(line, "# ")
			}
			line = strings.TrimPrefix(line, "CONFIG_")
			if *y || *m {
				line = strings.Split(line, "=")[0]
			}
			if *n {
				line = strings.TrimSuffix(line, notset)
			}
		}
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		log.Fatalf("reading %s: %s", cfgfile, err)
	}
}
