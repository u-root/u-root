// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// Serve files on the network.
//
// Synopsis:
//
//	srvfiles [--h=HOST] [--p=PORT] [--d=DIR]
//
// Options:
//
//	--h: hostname (default: 127.0.0.1)
//	--p: port number (default: 8080)
//	--d: directory to serve (default: .)
package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

var (
	host = flag.String("h", "127.0.0.1", "hostname")
	port = flag.String("p", "8080", "port number")
	dir  = flag.String("d", ".", "directory to serve")
)

var cacheHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func maxAgeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range cacheHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	http.Handle("/", maxAgeHandler(http.FileServer(http.Dir(*dir))))
	log.Fatal(http.ListenAndServe(net.JoinHostPort(*host, *port), nil))
}
