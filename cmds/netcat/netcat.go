// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
netcat connects to a place and sends data to it and from it.
*/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	var c net.Conn
	var err error
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	if c, err = net.Dial("tcp", os.Args[1]); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
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
