// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fdtdump prints a readable version of Flattened Device Tree or dtb.
//
// Synopsis:
//
//	fdtdump [-json] FILE
//
// Options:
//
//	-json: Print json with base64 encoded values.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/dt"
)

var asJSON = flag.Bool("json", false, "Print json with base64 encoded values.")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("usage: %s [-json] FILE", os.Args[0])
	}

	// Open and parse Device Tree.
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fdt, err := dt.ReadFDT(f)
	if err != nil {
		log.Fatal(err)
	}

	if *asJSON {
		out, err := json.MarshalIndent(fdt, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out))
	} else {
		//nolint:staticcheck
		if err := fdt.PrintDTS(os.Stdout); err != nil {
			log.Fatalf("error printing dts: %v", err)
		}
	}
}
