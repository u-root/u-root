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

var (
	netType = flag.String("net", "tcp", "What net type to use, e.g. tcp, unix, etc.")
	listen  = flag.Bool("l", false, "Listen for connections.")
	verbose = flag.Bool("v", false, "Verbose output.")
)

func init() {
	util.Usage(usage)
}

func main() {
	var c net.Conn
	var err error
	if flag.Parse(); len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	addr := flag.Args()[0]

	if *listen {
		ln, err := net.Listen(*netType, addr)
		if err != nil {
			log.Fatalln(err)
		}
		if *verbose {
			fmt.Fprintln(os.Stderr, "Listening on", ln.Addr())
		}

		c, err = ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if c, err = net.Dial(*netType, addr); err != nil {
			log.Fatalln(err)
		}
	}
	if *verbose {
		fmt.Fprintln(os.Stderr, "Connected to", c.RemoteAddr())
	}

	go func() {
		if _, err := io.Copy(c, os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	if _, err = io.Copy(os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if *verbose {
		fmt.Fprintln(os.Stderr, "Disconnected")
	}
}
