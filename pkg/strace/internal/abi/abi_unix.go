// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package abi

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// OpenMode represents the mode to open(2) a file.
var OpenMode = FlagSet{
	&Value{
		Value: syscall.O_RDWR,
		Name:  "O_RDWR",
	},
	&Value{
		Value: syscall.O_WRONLY,
		Name:  "O_WRONLY",
	},
	&Value{
		Value: syscall.O_RDONLY,
		Name:  "O_RDONLY",
	},
}

// OpenFlagSet is the set of open(2) flags.
var OpenFlagSet = FlagSet{
	&BitFlag{
		Value: syscall.O_APPEND,
		Name:  "O_APPEND",
	},
	&BitFlag{
		Value: syscall.O_ASYNC,
		Name:  "O_ASYNC",
	},
	&BitFlag{
		Value: syscall.O_CLOEXEC,
		Name:  "O_CLOEXEC",
	},
	&BitFlag{
		Value: syscall.O_CREAT,
		Name:  "O_CREAT",
	},
	&BitFlag{
		Value: syscall.O_DIRECT,
		Name:  "O_DIRECT",
	},
	&BitFlag{
		Value: syscall.O_DIRECTORY,
		Name:  "O_DIRECTORY",
	},
	&BitFlag{
		Value: syscall.O_EXCL,
		Name:  "O_EXCL",
	},
	&BitFlag{
		Value: syscall.O_NOATIME,
		Name:  "O_NOATIME",
	},
	&BitFlag{
		Value: syscall.O_NOCTTY,
		Name:  "O_NOCTTY",
	},
	&BitFlag{
		Value: syscall.O_NOFOLLOW,
		Name:  "O_NOFOLLOW",
	},
	&BitFlag{
		Value: syscall.O_NONBLOCK,
		Name:  "O_NONBLOCK",
	},
	&BitFlag{
		Value: 0x200000, // O_PATH
		Name:  "O_PATH",
	},
	&BitFlag{
		Value: syscall.O_SYNC,
		Name:  "O_SYNC",
	},
	&BitFlag{
		Value: syscall.O_TRUNC,
		Name:  "O_TRUNC",
	},
}

func Open(val uint64) string {
	s := OpenMode.Parse(val & syscall.O_ACCMODE)
	if flags := OpenFlagSet.Parse(val &^ syscall.O_ACCMODE); flags != "" {
		s += "|" + flags
	}
	return s
}

// socket

// SocketFamily are the possible socket(2) families.
var SocketFamily = FlagSet{
	&Value{
		Value: unix.AF_UNSPEC,
		Name:  "AF_UNSPEC",
	},
	&Value{
		Value: unix.AF_UNIX,
		Name:  "AF_UNIX",
	},
	&Value{
		Value: unix.AF_INET,
		Name:  "AF_INET",
	},
	&Value{
		Value: unix.AF_AX25,
		Name:  "AF_AX25",
	},
	&Value{
		Value: unix.AF_IPX,
		Name:  "AF_IPX",
	},
	&Value{
		Value: unix.AF_APPLETALK,
		Name:  "AF_APPLETALK",
	},
	&Value{
		Value: unix.AF_NETROM,
		Name:  "AF_NETROM",
	},
	&Value{
		Value: unix.AF_BRIDGE,
		Name:  "AF_BRIDGE",
	},
	&Value{
		Value: unix.AF_ATMPVC,
		Name:  "AF_ATMPVC",
	},
	&Value{
		Value: unix.AF_X25,
		Name:  "AF_X25",
	},
	&Value{
		Value: unix.AF_INET6,
		Name:  "AF_INET6",
	},
	&Value{
		Value: unix.AF_ROSE,
		Name:  "AF_ROSE",
	},
	&Value{
		Value: unix.AF_DECnet,
		Name:  "AF_DECnet",
	},
	&Value{
		Value: unix.AF_NETBEUI,
		Name:  "AF_NETBEUI",
	},
	&Value{
		Value: unix.AF_SECURITY,
		Name:  "AF_SECURITY",
	},
	&Value{
		Value: unix.AF_KEY,
		Name:  "AF_KEY",
	},
	&Value{
		Value: unix.AF_NETLINK,
		Name:  "AF_NETLINK",
	},
	&Value{
		Value: unix.AF_PACKET,
		Name:  "AF_PACKET",
	},
	&Value{
		Value: unix.AF_ASH,
		Name:  "AF_ASH",
	},
	&Value{
		Value: unix.AF_ECONET,
		Name:  "AF_ECONET",
	},
	&Value{
		Value: unix.AF_ATMSVC,
		Name:  "AF_ATMSVC",
	},
	&Value{
		Value: unix.AF_RDS,
		Name:  "AF_RDS",
	},
	&Value{
		Value: unix.AF_SNA,
		Name:  "AF_SNA",
	},
	&Value{
		Value: unix.AF_IRDA,
		Name:  "AF_IRDA",
	},
	&Value{
		Value: unix.AF_PPPOX,
		Name:  "AF_PPPOX",
	},
	&Value{
		Value: unix.AF_WANPIPE,
		Name:  "AF_WANPIPE",
	},
	&Value{
		Value: unix.AF_LLC,
		Name:  "AF_LLC",
	},
	&Value{
		Value: unix.AF_IB,
		Name:  "AF_IB",
	},
	&Value{
		Value: unix.AF_MPLS,
		Name:  "AF_MPLS",
	},
	&Value{
		Value: unix.AF_CAN,
		Name:  "AF_CAN",
	},
	&Value{
		Value: unix.AF_TIPC,
		Name:  "AF_TIPC",
	},
	&Value{
		Value: unix.AF_BLUETOOTH,
		Name:  "AF_BLUETOOTH",
	},
	&Value{
		Value: unix.AF_IUCV,
		Name:  "AF_IUCV",
	},
	&Value{
		Value: unix.AF_RXRPC,
		Name:  "AF_RXRPC",
	},
	&Value{
		Value: unix.AF_ISDN,
		Name:  "AF_ISDN",
	},
	&Value{
		Value: unix.AF_PHONET,
		Name:  "AF_PHONET",
	},
	&Value{
		Value: unix.AF_IEEE802154,
		Name:  "AF_IEEE802154",
	},
	&Value{
		Value: unix.AF_CAIF,
		Name:  "AF_CAIF",
	},
	&Value{
		Value: unix.AF_ALG,
		Name:  "AF_ALG",
	},
	&Value{
		Value: unix.AF_NFC,
		Name:  "AF_NFC",
	},
	&Value{
		Value: unix.AF_VSOCK,
		Name:  "AF_VSOCK",
	},
}

// SocketType are the possible socket(2) types.
var SocketType = FlagSet{
	&Value{
		Value: unix.SOCK_STREAM,
		Name:  "SOCK_STREAM",
	},
	&Value{
		Value: unix.SOCK_DGRAM,
		Name:  "SOCK_DGRAM",
	},
	&Value{
		Value: unix.SOCK_RAW,
		Name:  "SOCK_RAW",
	},
	&Value{
		Value: unix.SOCK_RDM,
		Name:  "SOCK_RDM",
	},
	&Value{
		Value: unix.SOCK_SEQPACKET,
		Name:  "SOCK_SEQPACKET",
	},
	&Value{
		Value: unix.SOCK_DCCP,
		Name:  "SOCK_DCCP",
	},
	&Value{
		Value: unix.SOCK_PACKET,
		Name:  "SOCK_PACKET",
	},
}

// SocketFlagSet are the possible socket(2) flags.
var SocketFlagSet = FlagSet{
	&BitFlag{
		Value: unix.SOCK_CLOEXEC,
		Name:  "SOCK_CLOEXEC",
	},
	&BitFlag{
		Value: unix.SOCK_NONBLOCK,
		Name:  "SOCK_NONBLOCK",
	},
}

// ipProtocol are the possible socket(2) types for INET and INET6 sockets.
var ipProtocol = FlagSet{
	&Value{
		Value: unix.IPPROTO_IP,
		Name:  "IPPROTO_IP",
	},
	&Value{
		Value: unix.IPPROTO_ICMP,
		Name:  "IPPROTO_ICMP",
	},
	&Value{
		Value: unix.IPPROTO_IGMP,
		Name:  "IPPROTO_IGMP",
	},
	&Value{
		Value: unix.IPPROTO_IPIP,
		Name:  "IPPROTO_IPIP",
	},
	&Value{
		Value: unix.IPPROTO_TCP,
		Name:  "IPPROTO_TCP",
	},
	&Value{
		Value: unix.IPPROTO_EGP,
		Name:  "IPPROTO_EGP",
	},
	&Value{
		Value: unix.IPPROTO_PUP,
		Name:  "IPPROTO_PUP",
	},
	&Value{
		Value: unix.IPPROTO_UDP,
		Name:  "IPPROTO_UDP",
	},
	&Value{
		Value: unix.IPPROTO_IDP,
		Name:  "IPPROTO_IDP",
	},
	&Value{
		Value: unix.IPPROTO_TP,
		Name:  "IPPROTO_TP",
	},
	&Value{
		Value: unix.IPPROTO_DCCP,
		Name:  "IPPROTO_DCCP",
	},
	&Value{
		Value: unix.IPPROTO_IPV6,
		Name:  "IPPROTO_IPV6",
	},
	&Value{
		Value: unix.IPPROTO_RSVP,
		Name:  "IPPROTO_RSVP",
	},
	&Value{
		Value: unix.IPPROTO_GRE,
		Name:  "IPPROTO_GRE",
	},
	&Value{
		Value: unix.IPPROTO_ESP,
		Name:  "IPPROTO_ESP",
	},
	&Value{
		Value: unix.IPPROTO_AH,
		Name:  "IPPROTO_AH",
	},
	&Value{
		Value: unix.IPPROTO_MTP,
		Name:  "IPPROTO_MTP",
	},
	&Value{
		Value: unix.IPPROTO_BEETPH,
		Name:  "IPPROTO_BEETPH",
	},
	&Value{
		Value: unix.IPPROTO_ENCAP,
		Name:  "IPPROTO_ENCAP",
	},
	&Value{
		Value: unix.IPPROTO_PIM,
		Name:  "IPPROTO_PIM",
	},
	&Value{
		Value: unix.IPPROTO_COMP,
		Name:  "IPPROTO_COMP",
	},
	&Value{
		Value: unix.IPPROTO_SCTP,
		Name:  "IPPROTO_SCTP",
	},
	&Value{
		Value: unix.IPPROTO_UDPLITE,
		Name:  "IPPROTO_UDPLITE",
	},
	&Value{
		Value: unix.IPPROTO_MPLS,
		Name:  "IPPROTO_MPLS",
	},
	&Value{
		Value: unix.IPPROTO_RAW,
		Name:  "IPPROTO_RAW",
	},
}

// SocketProtocol are the possible socket(2) protocols for each protocol family.
var SocketProtocol = map[int32]FlagSet{
	unix.AF_INET:  ipProtocol,
	unix.AF_INET6: ipProtocol,
	unix.AF_NETLINK: {
		&Value{
			Value: unix.NETLINK_ROUTE,
			Name:  "NETLINK_ROUTE",
		},
		&Value{
			Value: unix.NETLINK_UNUSED,
			Name:  "NETLINK_UNUSED",
		},
		&Value{
			Value: unix.NETLINK_USERSOCK,
			Name:  "NETLINK_USERSOCK",
		},
		&Value{
			Value: unix.NETLINK_FIREWALL,
			Name:  "NETLINK_FIREWALL",
		},
		&Value{
			Value: unix.NETLINK_SOCK_DIAG,
			Name:  "NETLINK_SOCK_DIAG",
		},
		&Value{
			Value: unix.NETLINK_NFLOG,
			Name:  "NETLINK_NFLOG",
		},
		&Value{
			Value: unix.NETLINK_XFRM,
			Name:  "NETLINK_XFRM",
		},
		&Value{
			Value: unix.NETLINK_SELINUX,
			Name:  "NETLINK_SELINUX",
		},
		&Value{
			Value: unix.NETLINK_ISCSI,
			Name:  "NETLINK_ISCSI",
		},
		&Value{
			Value: unix.NETLINK_AUDIT,
			Name:  "NETLINK_AUDIT",
		},
		&Value{
			Value: unix.NETLINK_FIB_LOOKUP,
			Name:  "NETLINK_FIB_LOOKUP",
		},
		&Value{
			Value: unix.NETLINK_CONNECTOR,
			Name:  "NETLINK_CONNECTOR",
		},
		&Value{
			Value: unix.NETLINK_NETFILTER,
			Name:  "NETLINK_NETFILTER",
		},
		&Value{
			Value: unix.NETLINK_IP6_FW,
			Name:  "NETLINK_IP6_FW",
		},
		&Value{
			Value: unix.NETLINK_DNRTMSG,
			Name:  "NETLINK_DNRTMSG",
		},
		&Value{
			Value: unix.NETLINK_KOBJECT_UEVENT,
			Name:  "NETLINK_KOBJECT_UEVENT",
		},
		&Value{
			Value: unix.NETLINK_GENERIC,
			Name:  "NETLINK_GENERIC",
		},
		&Value{
			Value: unix.NETLINK_SCSITRANSPORT,
			Name:  "NETLINK_SCSITRANSPORT",
		},
		&Value{
			Value: unix.NETLINK_ECRYPTFS,
			Name:  "NETLINK_ECRYPTFS",
		},
		&Value{
			Value: unix.NETLINK_RDMA,
			Name:  "NETLINK_RDMA",
		},
		&Value{
			Value: unix.NETLINK_CRYPTO,
			Name:  "NETLINK_CRYPTO",
		},
	},
}

var ControlMessageType = map[int32]string{
	unix.SCM_RIGHTS:      "SCM_RIGHTS",
	unix.SCM_CREDENTIALS: "SCM_CREDENTIALS",
	unix.SO_TIMESTAMP:    "SO_TIMESTAMP",
}

func SockType(stype int32) string {
	s := SocketType.Parse(uint64(stype & SOCK_TYPE_MASK))
	if flags := SocketFlagSet.Parse(uint64(stype &^ SOCK_TYPE_MASK)); flags != "" {
		s += "|" + flags
	}
	return s
}

func SockProtocol(family, protocol int32) string {
	protocols, ok := SocketProtocol[family]
	if !ok {
		return fmt.Sprintf("%#x", protocol)
	}
	return protocols.Parse(uint64(protocol))
}

func SockFlags(flags int32) string {
	if flags == 0 {
		return "0"
	}
	return SocketFlagSet.Parse(uint64(flags))
}

// MessageHeader64 is the 64-bit representation of the msghdr struct used in
// the recvmsg and sendmsg syscalls.
type MessageHeader64 struct {
	// Name is the optional pointer to a network address buffer.
	Name uint64

	// NameLen is the length of the buffer pointed to by Name.
	NameLen uint32
	_       uint32

	// Iov is a pointer to an array of io vectors that describe the memory
	// locations involved in the io operation.
	Iov uint64

	// IovLen is the length of the array pointed to by Iov.
	IovLen uint64

	// Control is the optional pointer to ancillary control data.
	Control uint64

	// ControlLen is the length of the data pointed to by Control.
	ControlLen uint64

	// Flags on the sent/received message.
	Flags int32
	_     int32
}
