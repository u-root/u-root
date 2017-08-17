// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print process information.
//
// Synopsis:
//     id
//
// Description:
//     id displays the uid, guid and groups of the calling process
//
// Options:
package main

import (
	"fmt"
	"log"
	"syscall"
)

var ()

func main() {
	uid := syscall.Getuid()
	gid := syscall.Getgid()
	groups, err := syscall.Getgroups()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uid: %d\n", uid)
	fmt.Printf("gid: %d\n", gid)

	fmt.Print("groups: ")
	for _, group := range groups {
		fmt.Printf("%d ", group)
	}
	fmt.Println()

}
