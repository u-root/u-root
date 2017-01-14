// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gopxe is a pxe program written in Go.
//
// Description:
//     It includes the functionality of the the PXE loader and pxelinux. It is
//     incomplete but anyone with interest in filling it out should find this
//     pretty easy. You can extend it to fetch a file following the rules of
//     pxe file name formation, then fetch the files, then use the kexec system
//     call (see kexec.go) to exec it.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	tftp "github.com/vcabbage/trivialt"
)

func main() {
	// in.tftpd does not support ClientTransferSize but the trivialt package sets it by default.
	opts := []tftp.ClientOpt{tftp.ClientTransferSize(false)}
	flag.Parse()

	c, err := tftp.NewClient(opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client is %v\n", c)

	f, err := c.Get("tftp://192.168.28.128//pxelinux.0")
	fmt.Printf("Get: r %v err %v\n", f, err)

	data, err := ioutil.ReadAll(f)
	fmt.Printf("len(data) %v, err %v", len(data), err)
}
