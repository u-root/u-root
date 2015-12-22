// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Wget reads one file from the argument and writes it on the standard output.
package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func wget(arg string, w io.Writer) error {
	resp, err := http.Get(arg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	if err := wget(os.Args[1], os.Stdout); err != nil {
		log.Fatalf("%v\n", err)
	}
}
