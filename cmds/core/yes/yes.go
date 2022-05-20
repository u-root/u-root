// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func runYes(w io.Writer, count uint64, args ...string) error {
	yes := "y"
	if len(args) > 0 {
		yes = strings.Join(args, " ")
	}
	// If count == 0 this loop runs to infinity, every other value of count
	// will iterate the loop for exactly "count" times before it breaks out.
	// The standard behavior is achieved with count = 0 as done in main()
	// Sadly this is required to make it testable with the recommended pattern.
	for {
		if _, err := fmt.Fprintf(w, "%s\n", yes); err != nil {
			return err
		}
		if count > 1 {
			count--
		} else if count == 1 {
			break
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := runYes(os.Stdout, 0, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
