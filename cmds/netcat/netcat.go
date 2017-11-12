// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Netcat pipes over the network.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "netcat [go-style network address]"

var netType = flag.String("net", "tcp", "What net type to use, e.g. tcp, unix, etc.")

func init() {
	util.Usage(usage)
}

func main() {
	flag.Parse()
	var c net.Conn
	var err error
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if c, err = net.Dial(*netType, flag.Args()[0]); err != nil {
		log.Fatalf("%v\n", err)
	}
	go func() {
		if _, err := io.Copy(c, os.Stdin); err != nil {
			fmt.Printf("%v", err)
		}
	}()
	if _, err = io.Copy(os.Stdout, c); err != nil {
		fmt.Printf("%v", err)
	}
}
