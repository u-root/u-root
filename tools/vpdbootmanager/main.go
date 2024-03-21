// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

func getUsage(progname string) string {
	return fmt.Sprintf(`Usage:
%s add [netboot [dhcpv6|dhcpv4] [MAC] | localboot [grub|path [Device GUID] [Kernel Path]]]
%s get [variable name]
%s set [variable name] [variable value]
%s delete [variable name]
%s dump

Ex.
add localboot grub
add netboot dhcpv6 AA:BB:CC:DD:EE:FF
get
get firmware_version
set systemboot_log_level 6
delete systemboot_log_level

Flags for netboot:

-override-url - an optional URL used to override the boot file URL used
-retries - the number of times a DHCP request should be retried if failed

Flags for localboot:

-kernel-args - additional kernel args
-ramfs - path of ramfs to be used for kexec'ing into the target kernel

Global flags:

-vpd-dir - VPD dir to use

`, progname, progname, progname, progname, progname)
}

func main() {
	if err := cli(os.Args[1:]); err != nil {
		fmt.Println(getUsage(os.Args[0]))
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
}

func cli(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("you need to provide action")
	}
	switch args[0] {
	case "add":
		return add(args[1], args[2:])
	case "get":
		varname := ""
		if len(args) > 1 {
			varname = args[1]
		}
		getter := NewGetter()
		return getter.Print(varname)
	case "set":
		if len(args) == 3 {
			err := set(args[1], args[2])
			if err == nil {
				fmt.Println("Successfully set, it will take effect after reboot")
			}
			return err
		}
	case "delete":
		if len(args) == 2 {
			err := remove(args[1])
			if err == nil {
				fmt.Println("Successfully deleted, it will take effect after reboot")
			}
			return err
		}
	case "dump":
		return dump()
	}
	return fmt.Errorf("unrecognized action")
}
