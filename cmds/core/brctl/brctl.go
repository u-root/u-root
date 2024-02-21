// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// brctl - ethernet bridge administration
//
// Synopsis:
//
// brctl addbr <name> creates a new instance of the ethernet bridge
// brctl delbr <name> deletes the instance <name> of the ethernet bridge
// brctl show shows all current instances of the ethernet bridge
//
// brctl addif <brname> <ifname> will make the interface <ifname> a port of the bridge <brname>
// brctl delif <brname> <ifname> will detach the interface <ifname> from the bridge <brname>
// brctl show <brname> will show some information on the bridge and its attached ports
//
// brctl showmacs <brname> shows a list of learned MAC addresses for this bridge
// brctl setageingtime <brname> <time> sets the ethernet (MAC) address ageing time, in seconds [OPT]
// brctl setgcint <brname> <time> sets the garbage collection interval for the bridge <brname> to <time> seconds [OPT]
//
// TODO: Spanning Tree Protocol
// See: https://elixir.bootlin.com/busybox/latest/source/networking/brctl.c
// Author:
//
//	Leon Gross (leon.gross@9elements.com)
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func ioctl_str(fd int, req uint, raw string) (int, error) {
	local_bytes := append([]byte(raw), 0)
	err_int, _, err_str := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&local_bytes[0])))
	return int(err_int), fmt.Errorf("%s", err_str)
}

const usage = "brctl [commands]"

type subcommand struct {
	name  string
	nargs int
}

// subcommands
// https://elixir.bootlin.com/busybox/latest/source/networking/brctl.c#L583
func addbr(name string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := ioctl_str(brctl_socket, unix.SIOCBRADDBR, name); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// TODO: merge add and del?
func delbr(name string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := ioctl_str(brctl_socket, unix.SIOCBRDELBR, name); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func addif(name string, iface string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	ifreg, err := unix.NewIfreq(name)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ioctl_str(brctl_socket, unix.SIOCBRADDIF, iface)
}

// func show(names ...string) error {
// 	if len(names) == 0 {
// 		// show all bridges

// 	} else {
// 		for _, name := range names {
// 		}
// 	}
// 	return nil
// }

func run(out io.Writer, argv []string) error {
	command := argv[0]
	args := argv[1:]

	fmt.Printf("argv = %v\n", args)

	var err error

	switch command {
	case "addbr":
		if len(args) != 1 {
			return fmt.Errorf("too few args")
		}
		err = addbr(args[0])

	case "delbr":
		if len(args) != 1 {
			return fmt.Errorf("too few args")
		}
		err = delbr(args[0])
	// case "show":
	// 	err = show()
	case "addif":
		if len(args) != 2 {
			return fmt.Errorf("too few args")
		}
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return err
}

func main() {
	argv := os.Args

	if len(argv) < 2 {
		log.Fatal(usage)
		os.Exit(1)
	}

	fmt.Printf("argv = %v, argc = %d\n", argv, len(argv))
	if err := run(os.Stdout, argv[1:]); err != nil {
		log.Fatalf("brctl: %v", err)
	}
}
