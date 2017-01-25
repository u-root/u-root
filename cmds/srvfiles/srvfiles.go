// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Serve files on the network.
//
// Synopsis:
//     srvfiles [--h=HOST] [--p=PORT] [--d=DIR]
//
// Options:
//     --h: hostname (default: 127.0.0.1)
//     --p: port number (default: 8080)
//     --d: directory to serve (default: .)
package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	host = flag.String("h", "127.0.0.1", "hostname")
	port = flag.String("p", "8080", "port number")
	dir  = flag.String("d", ".", "directory to serve")
)

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(*host+":"+*port, http.FileServer(http.Dir(*dir))))
}
