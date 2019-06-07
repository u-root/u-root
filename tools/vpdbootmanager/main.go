// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

const usage = `Usage:
%s add [netboot [dhcpv6|dhcpv4] [MAC] | localboot [grub|path [Device GUID] [Kernel Path]]]

Ex.
add localboot grub
add netboot dhcpv6 AA:BB:CC:DD:EE:FF

Flags for netboot:

-override-url - an optional URL used to override the boot file URL used
-retries - the number of times a DHCP request should be retried if failed

Flags for localboot:

-kernel-args - additional kernel args
-ramfs - path of ramfs to be used for kexec'ing into the target kernel

Global flags:

-vpd-dir - VPD dir to use

`

func main() {
	if err := cli(os.Args[1:]); err != nil {
		fmt.Printf(usage, os.Args[0])
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
}

func cli(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("You need to provide action")
	}
	switch args[0] {
	case "add":
		return add(args[1], args[2:])
	}
	return fmt.Errorf("Unrecognized action")
}
