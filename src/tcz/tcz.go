// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Wget reads one file from the argument and writes it on the standard output.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

const tcz = "/tinycorelinux.net/5.x/x86/tcz"

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	cmdName := os.Args[1]

	if err := os.MkdirAll(tcz, 0600); err != nil {
		log.Fatal(err)
	}
	
	// path.Join doesn't quite work here. 
	filepath := path.Join(tcz, cmdName)
	cmd := "http:/" + filepath

	resp, err := http.Get(cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// we've to the whole tcz in resp.Body.
	// First, save it to /tcz/name
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Fatal(err)
	}
	/* OK, these are compressed tars ... */
}
