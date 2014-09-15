// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"syscall"
	"unsafe"
)

var l = log.New(os.Stdout, "ip: ", 0)

func adddelip(op, ip, dev string) error {
	addr, network, err := net.ParseCIDR(ip)
	if err != nil {
		addr = net.ParseIP(ip)
	}

	iface, err := net.InterfaceByName(dev)
	if err != nil {
		l.Fatalf("%v not found", dev)
		return err
	}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		l.Fatalf("socket %v", err)
		return err
	}
	// let's face it. The whole bsd interface sucks. They force you
	// to think about endianness, and what byte goes where. It's never been right.
	// It's a 30 year old botch. Let's not play this stupid game.
	// How I miss Plan 9 at times. You're welcome to fix this, but to do it right
	// you need to fix the netlink support in the net package, and I don't have the
	// time to do that just now.
	newaddr := &[128]byte{}
	copy(newaddr[0:], dev)
	newaddr[16+0] = syscall.AF_INET
	// that's how bad this all is.
	newaddr[20] = addr[12]
	newaddr[21] = addr[13]
	newaddr[22] = addr[14]
	newaddr[23] = addr[15]

	rv1, rv2, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFADDR, uintptr(unsafe.Pointer(&newaddr[0])))
	l.Printf("addr %v network %v iface %v fd %v rv1 %v rv2 %v",
		addr, network, iface, fd, rv1, rv2)
	if errno != 0 {
		l.Fatalf("ioctl SIOCSIFADDR BAD %v", error(errno))
		return err
	}

	// now bring it up.
	// this is a short cut. I have other things to get right first.
	flags := uint16(syscall.IFF_UP | syscall.IFF_BROADCAST | syscall.IFF_RUNNING)
	newaddr[16] = byte(flags >> 8)
	newaddr[17] = byte(flags)
	rv1, rv2, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&newaddr[0])))
	if errno != 0 {
		l.Fatalf("ioctl SIOCSIFFLAGS BAD %v", error(errno))
		return err
	}
	return nil

}
func addroute(ip, dev string) error {
	addr, network, err := net.ParseCIDR(ip)
	if err != nil {
		addr = net.ParseIP(ip)
	}

	iface, err := net.InterfaceByName(dev)
	if err != nil {
		l.Fatalf("%v not found", dev)
		return err
	}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		l.Fatalf("socket %v", err)
		return err
	}
	// It's a 30 year old botch. Let's not play this stupid game.
	// How I miss Plan 9 at times. You're welcome to fix this, but to do it right
	// you need to fix the netlink support in the net package, and I don't have the
	// time to do that just now.
	newaddr := &[128]byte{}
	copy(newaddr[0:], dev)
	newaddr[16+0] = syscall.AF_INET
	// that's how bad this all is.
	newaddr[20] = addr[12]
	newaddr[21] = addr[13]
	newaddr[22] = addr[14]
	newaddr[23] = addr[15]

	rv1, rv2, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFADDR, uintptr(unsafe.Pointer(&newaddr[0])))
	l.Printf("addr %v network %v iface %v fd %v rv1 %v rv2 %v",
		addr, network, iface, fd, rv1, rv2)
	if errno != 0 {
		l.Fatalf("ioctl SIOCSIFADDR BAD %v", error(errno))
		return err
	}

	// now bring it up.
	// this is a short cut. I have other things to get right first.
	flags := uint16(syscall.IFF_UP | syscall.IFF_BROADCAST | syscall.IFF_RUNNING)
	newaddr[16] = byte(flags >> 8)
	newaddr[17] = byte(flags)
	rv1, rv2, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&newaddr[0])))
	if errno != 0 {
		l.Fatalf("ioctl SIOCSIFFLAGS BAD %v", error(errno))
		return err
	}
	return nil

}
func main() {
	flag.Parse()
	arg := flag.Args()
	if len(arg) < 1 {
		l.Fatalf("arg count")
	}
	switch {
	case len(arg) == 5 && arg[0] == "addr" && arg[1] == "add" && arg[3] == "dev":
		adddelip(arg[1], arg[2], arg[4])

	case len(arg) == 1 && arg[0] == "link":
		fallthrough
	case len(arg) == 2 && arg[0] == "link" && arg[1] == "show":
		ifaces, err := net.Interfaces()
		if err != nil {
			l.Fatalf("Can't enumerate interfaces? %v", err)
		}
		for _, v := range ifaces {
			addrs, err := v.Addrs()
			if err != nil {
				l.Printf("Can't enumerate addresses")
			}
			l.Printf("%v: %v", v, addrs)
		}
	case len(arg) == 1 && arg[0] == "route":
		if b, err := ioutil.ReadFile("/proc/net/route"); err == nil {
			l.Printf("%s", string(b))
		} else {
			l.Fatalf("Route failed: %v", err)
		}
	// oh, barf.
	case len(arg) == 8 && arg[0] == "route" && arg[1] == "add" && arg[2] == "default" && arg[3] == "via" && arg[5] == "dev":
	     AddDefaultGw(arg[4], arg[6])
	default:
		l.Fatalf("We don't do this: %v; try addr or link or route", arg)
	}
}
