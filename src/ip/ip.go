// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net"
	"os"
	"syscall"
)

var l = log.New(os.Stdout, "ip: ", 0)

/*
type NlMsghdr struct {
        Len   uint32
        Type  uint16
        Flags uint16
        Seq   uint32
        Pid   uint32
}
*/

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

	// do the deed.

/* FUCK NETLINK IN GO FOR NOW. It's not complete.
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, 0)
	if err != nil {
		l.Fatalf("Socket: %v", err)
	}
	defer syscall.Close(fd)
	sz := int(32768)
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, sz)
	if err != nil {
		l.Fatalf("setsocktop SNDBUF: %v", err)
	}
	sz = 1048576
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, sz)
	if err != nil {
		l.Fatalf("setsocktop RCVBUF: %v", err)
	}
    	lsa := &syscall.SockaddrNetlink{Family: syscall.AF_NETLINK}
	if err = syscall.Bind(fd, lsa); err != nil {
		l.Fatalf("Bind SockaddrNetline %v: %v", lsa, err)
	}
		rr := &syscall.NetlinkRouteRequest{}
		rr.Header.Len = uint32(syscall.NLMSG_HDRLEN + syscall.SizeofRtGenmsg)
		rr.Header.Type = uint16(syscall.IPPROTO_IP)
	rr.Header.Flags = syscall.NLM_F_REPLACE
		rr.Header.Seq = uint32(1)
		rr.Data.Family = uint8(syscall.AF_INET)
		tow := rr.toWireFormat()
	l.Printf("wire format: %v -> %v", rr, tow)
/*
type SockaddrNetlink struct {
    Family uint16
    Pad    uint16
    Pid    uint32
    Groups uint32
    raw    RawSockaddrNetlink
}
		defer Close(s)
    58		lsa := &SockaddrNetlink{Family: AF_NETLINK}
    59		if err := Bind(s, lsa); err != nil {
    60			return nil, err
    61		}
func newNetlinkRouteRequest(proto, seq, family int) []byte {
    41		rr := &NetlinkRouteRequest{}
    42		rr.Header.Len = uint32(NLMSG_HDRLEN + SizeofRtGenmsg)
    43		rr.Header.Type = uint16(proto)
    44		rr.Header.Flags = NLM_F_DUMP | NLM_F_REQUEST
    45		rr.Header.Seq = uint32(seq)
    46		rr.Data.Family = uint8(family)
    47		return rr.toWireFormat()
    48	}
*/
	// err = syscall.Bind(fd, sa_family, 12)
	l.Fatalf("sock: fd %d, addr %v, network %v, iface %v", fd, addr, network, iface)
	//1458  socket(PF_NETLINK, SOCK_RAW|SOCK_CLOEXEC, 0) = 3
	//1458  setsockopt(3, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
	//1458  setsockopt(3, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0
	//1458  bind(3, {sa_family=AF_NETLINK, pid=0, groups=00000000}, 12) = 0
	//1458  getsockname(3, {sa_family=AF_NETLINK, pid=1458, groups=00000000}, [12]) = 0
	//1458  sendmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"(\0\0\0\24\0\5\6* \16T\0\0\0\0\2 \0\0\f\0\0\0\10\0\2\0\300\250\0\1\10\0\1\//0\300\250\0\1", 40}], msg_controllen=0, msg_flags=0}, 0) = 40
	//1458  recvmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"$\0\0\0\2\0\0\0* \16T\262\5\0\0\0\0\0\0(\0\0\0\24\0\5\6* \16T\0\0\0\0", 16//384}], msg_controllen=0, msg_flags=0}, 0) = 36
	return nil
}
func main() {
	var err error
	flag.Parse()
	arg := flag.Args()
	if len(arg) < 1 {
		l.Fatalf("arg count")
	}
	switch {
	case len(arg) == 5 && arg[0] == "addr" && arg[1] == "add" && arg[3] == "dev":
		err = adddelip(arg[1], arg[2], arg[4])
	default:
		l.Fatalf("We don't do this: %v", arg)
	}
	if err != nil {
		l.Fatalf("%v: %v", arg, err)
	}
}
