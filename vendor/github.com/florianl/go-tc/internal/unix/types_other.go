//go:build !linux
// +build !linux

package unix

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

const (
	AF_UNSPEC     = 0x0
	NETLINK_ROUTE = 0x0
	IFLA_EXT_MASK = 0x1d
	RTM_GETLINK   = 0x12
	RTNLGRP_TC    = 0x4
)

const (
	RTM_NEWTFILTER = 44
	RTM_DELTFILTER = 45
	RTM_GETTFILTER = 46
)

const (
	RTM_NEWTCLASS = 40
	RTM_DELTCLASS = 41
	RTM_GETTCLASS = 42
)

const (
	RTM_NEWQDISC = 36
	RTM_DELQDISC = 37
	RTM_GETQDISC = 38
)

const (
	RTM_NEWCHAIN = 100
	RTM_DELCHAIN = 101
	RTM_GETCHAIN = 102
)

const (
	RTM_NEWACTION = 48
	RTM_DELACTION = 49
	RTM_GETACTION = 50
)
