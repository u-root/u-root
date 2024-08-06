//go:build linux
// +build linux

package unix

import linux "golang.org/x/sys/unix"

// IfInfomsg makes linter happy with this comment.
type IfInfomsg = linux.IfInfomsg

// Make linter happy with this comment.
const (
	AF_UNSPEC     = linux.AF_UNSPEC
	NETLINK_ROUTE = linux.NETLINK_ROUTE
	IFLA_EXT_MASK = linux.IFLA_EXT_MASK
	RTM_GETLINK   = linux.RTM_GETLINK
	RTNLGRP_TC    = linux.RTNLGRP_TC
)

// Make linter happy with this comment.
const (
	RTM_NEWTFILTER = linux.RTM_NEWTFILTER
	RTM_DELTFILTER = linux.RTM_DELTFILTER
	RTM_GETTFILTER = linux.RTM_GETTFILTER
)

// Make linter happy with this comment.
const (
	RTM_NEWTCLASS = linux.RTM_NEWTCLASS
	RTM_DELTCLASS = linux.RTM_DELTCLASS
	RTM_GETTCLASS = linux.RTM_GETTCLASS
)

// Make linter happy with this comment.
const (
	RTM_NEWQDISC = linux.RTM_NEWQDISC
	RTM_DELQDISC = linux.RTM_DELQDISC
	RTM_GETQDISC = linux.RTM_GETQDISC
)

// Make linter happy with this comment.
const (
	RTM_NEWCHAIN = linux.RTM_NEWCHAIN
	RTM_DELCHAIN = linux.RTM_DELCHAIN
	RTM_GETCHAIN = linux.RTM_GETCHAIN
)

// Make linter happy with this comment.
const (
	RTM_NEWACTION = linux.RTM_NEWACTION
	RTM_DELACTION = linux.RTM_DELACTION
	RTM_GETACTION = linux.RTM_GETACTION
)
