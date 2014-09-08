// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
       "flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

var l = log.New(os.Stdout, "ip: ", 0)

func adddelip(op, ip, dev string) error {
     i := clip(ip)
     iface = getiface(dev)
}
func main() {
     var err error
     flag.Parse()
     arg := flag.Args()
	if len(arg) < 2 {
		os.Exit(1)
	}
	switch {
	case len(arg) == 6 && arg[1] == "addr" && arg[2] == "add" && arg[4] == "dev":
	     err = adddelip(arg[1], argv[3], argv[5])
	default:
	l.Fatalf("We don't do this: %v", arg)
	}
	if err != nil {
	   l.Fatalf("%v: %v", arg,err)
	   }
}
