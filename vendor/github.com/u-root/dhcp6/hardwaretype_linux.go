// +build linux

// Copyright (c) 2012 The Go Authors. All rights reserved.
// Source code in this file is based on src/net/interface_linux.go,
// from the Go standard library.  The Go license can be found here:
// https://golang.org/LICENSE.

package dhcp6

import (
	"net"
	"os"
	"syscall"
	"unsafe"
)

// HardwareType returns the IANA ARP parameter Hardware Type value
// for an input network interface.  A table of known values can
// be found here:
// http://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml.
//
// If an error occurs, a syscall error is returned.  If no hardware type
// is found, ErrParseHardwareType is returned.
func HardwareType(ifi *net.Interface) (uint16, error) {
	// Get link information from netlink
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	if err != nil {
		return 0, os.NewSyscallError("netlink rib", err)
	}

	// Parse information into messages
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return 0, os.NewSyscallError("netlink message", err)
	}

	// Check messages for information
	for _, m := range msgs {
		switch m.Header.Type {
		// No more messages
		case syscall.NLMSG_DONE:
			break
		// Network interface message
		case syscall.RTM_NEWLINK:
			ifim := (*syscall.IfInfomsg)(unsafe.Pointer(&m.Data[0]))
			if ifi.Index == int(ifim.Index) {
				return ifim.Type, nil
			}
		}
	}

	// Did not find hardware type
	return 0, ErrParseHardwareType
}
