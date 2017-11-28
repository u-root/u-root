// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print build information about the kernel and machine.
//
// Synopsis:
//     uname [-asnrvmd]
//
// Options:
//     -a: print everything
//     -s: print the kernel name
//     -n: print the network node name
//     -r: print the kernel release
//     -v: print the kernel version
//     -m: print the machine hardware name
//     -d: print your domain name
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	all     = flag.Bool("a", false, "print everything")
	kernel  = flag.Bool("s", false, "print the kernel name")
	node    = flag.Bool("n", false, "print the network node name")
	release = flag.Bool("r", false, "print the kernel release")
	version = flag.Bool("v", false, "print the kernel version")
	machine = flag.Bool("m", false, "print the machine hardware name")
	domain  = flag.Bool("d", false, "print your domain name")
)

func toString(d []byte) string {
	return string(d[:bytes.IndexByte(d[:], 0)])
}

func handleFlags(u *unix.Utsname) string {
	Sysname, Nodename := toString(u.Sysname[:]), toString(u.Nodename[:])
	Release, Version := toString(u.Release[:]), toString(u.Version[:])
	Machine, Domainname := toString(u.Machine[:]), toString(u.Domainname[:])
	info := make([]string, 0, 6)

	if *all || flag.NFlag() == 0 {
		info = append(info, Sysname, Nodename, Release, Version, Machine, Domainname)
		goto end
	}
	if *kernel {
		info = append(info, Sysname)
	}
	if *node {
		info = append(info, Nodename)
	}
	if *release {
		info = append(info, Release)
	}
	if *version {
		info = append(info, Version)
	}
	if *machine {
		info = append(info, Machine)
	}
	if *domain {
		info = append(info, Domainname)
	}

end:
	return strings.Join(info, " ")
}

func main() {
	flag.Parse()

	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		log.Fatalf("%v", err)
	}

	info := handleFlags(&u)
	fmt.Println(info)
}
