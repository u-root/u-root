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
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/util"
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

func handleFlags(u *util.Utsname) string {
	info := make([]string, 0, 6)

	if *all || flag.NFlag() == 0 {
		info = append(info, u.Sysname, u.Nodename, u.Release, u.Version, u.Machine, u.Domainname)
		goto end
	}
	if *kernel {
		info = append(info, u.Sysname)
	}
	if *node {
		info = append(info, u.Nodename)
	}
	if *release {
		info = append(info, u.Release)
	}
	if *version {
		info = append(info, u.Version)
	}
	if *machine {
		info = append(info, u.Machine)
	}
	if *domain {
		info = append(info, u.Domainname)
	}

end:
	return strings.Join(info, " ")
}

func main() {
	flag.Parse()

	if u, err := util.Uname(); err != nil {
		log.Fatalf("%v", err)
	} else {
		info := handleFlags(u)
		fmt.Println(info)
	}
}
