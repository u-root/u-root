// Copyright 2009-2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ldd prints the full path of dependencies.
//
// Description:
//     Unlike the standard one, you can use it in a script,
//     e.g. i=`ldd whatever` leaves you with a list of files you can usefully
//     copy. You can also feed it a long list of files (/bin/*) and get a
//     short list of libraries; further, it will read stdin.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/uroot"
)

func usage() {
	log.Fatalf("usage: ldd file [file...]")
}

func ldd(s ...string) ([]string, error) {
	var libs []string
	l, err := uroot.Ldd(s)
	if err != nil {
		return nil, err
	}
	for _, i := range l {
		if i.Mode().IsRegular() {
			libs = append(libs, i.FullName)
		}
	}
	return libs, nil
}

func main() {
	l, err := ldd(os.Args[1:]...)
	if err != nil {
		log.Fatalf("ldd: %v", err)
	}
	fmt.Printf("%v\n", l)
}
