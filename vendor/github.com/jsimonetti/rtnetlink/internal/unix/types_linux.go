// +build linux

package unix

import (
	linux "golang.org/x/sys/unix"
)

const (
	AF_INET              = linux.AF_INET
	AF_INET6             = linux.AF_INET6
	AF_UNSPEC            = linux.AF_UNSPEC
	NETLINK_ROUTE        = linux.NETLINK_ROUTE
	SizeofIfAddrmsg      = linux.SizeofIfAddrmsg
	SizeofIfInfomsg      = linux.SizeofIfInfomsg
	SizeofNdMsg          = linux.SizeofNdMsg
	SizeofRtMsg          = linux.SizeofRtMsg
	RTM_NEWADDR          = linux.RTM_NEWADDR
	RTM_DELADDR          = linux.RTM_DELADDR
	RTM_GETADDR          = linux.RTM_GETADDR
	RTM_NEWLINK          = linux.RTM_NEWLINK
	RTM_DELLINK          = linux.RTM_DELLINK
	RTM_GETLINK          = linux.RTM_GETLINK
	RTM_SETLINK          = linux.RTM_SETLINK
	RTM_NEWROUTE         = linux.RTM_NEWROUTE
	RTM_DELROUTE         = linux.RTM_DELROUTE
	RTM_GETROUTE         = linux.RTM_GETROUTE
	RTM_NEWNEIGH         = linux.RTM_NEWNEIGH
	RTM_DELNEIGH         = linux.RTM_DELNEIGH
	RTM_GETNEIGH         = linux.RTM_GETNEIGH
	IFA_UNSPEC           = linux.IFA_UNSPEC
	IFA_ADDRESS          = linux.IFA_ADDRESS
	IFA_LOCAL            = linux.IFA_LOCAL
	IFA_LABEL            = linux.IFA_LABEL
	IFA_BROADCAST        = linux.IFA_BROADCAST
	IFA_ANYCAST          = linux.IFA_ANYCAST
	IFA_CACHEINFO        = linux.IFA_CACHEINFO
	IFA_MULTICAST        = linux.IFA_MULTICAST
	IFA_FLAGS            = linux.IFA_FLAGS
	IFF_UP               = linux.IFF_UP
	IFF_BROADCAST        = linux.IFF_BROADCAST
	IFF_LOOPBACK         = linux.IFF_LOOPBACK
	IFF_POINTOPOINT      = linux.IFF_POINTOPOINT
	IFF_MULTICAST        = linux.IFF_MULTICAST
	IFLA_UNSPEC          = linux.IFLA_UNSPEC
	IFLA_ADDRESS         = linux.IFLA_ADDRESS
	IFLA_BROADCAST       = linux.IFLA_BROADCAST
	IFLA_IFNAME          = linux.IFLA_IFNAME
	IFLA_MTU             = linux.IFLA_MTU
	IFLA_LINK            = linux.IFLA_LINK
	IFLA_QDISC           = linux.IFLA_QDISC
	IFLA_OPERSTATE       = linux.IFLA_OPERSTATE
	IFLA_STATS           = linux.IFLA_STATS
	IFLA_STATS64         = linux.IFLA_STATS64
	IFLA_LINKINFO        = linux.IFLA_LINKINFO
	IFLA_MASTER          = linux.IFLA_MASTER
	IFLA_INFO_KIND       = linux.IFLA_INFO_KIND
	IFLA_INFO_SLAVE_KIND = linux.IFLA_INFO_SLAVE_KIND
	IFLA_INFO_DATA       = linux.IFLA_INFO_DATA
	IFLA_INFO_SLAVE_DATA = linux.IFLA_INFO_SLAVE_DATA
	NDA_UNSPEC           = linux.NDA_UNSPEC
	NDA_DST              = linux.NDA_DST
	NDA_LLADDR           = linux.NDA_LLADDR
	NDA_CACHEINFO        = linux.NDA_CACHEINFO
	NDA_IFINDEX          = linux.NDA_IFINDEX
	RTA_UNSPEC           = linux.RTA_UNSPEC
	RTA_DST              = linux.RTA_DST
	RTA_PREFSRC          = linux.RTA_PREFSRC
	RTA_GATEWAY          = linux.RTA_GATEWAY
	RTA_OIF              = linux.RTA_OIF
	RTA_PRIORITY         = linux.RTA_PRIORITY
	RTA_TABLE            = linux.RTA_TABLE
	RTA_EXPIRES          = linux.RTA_EXPIRES
	NTF_PROXY            = linux.NTF_PROXY
	RTN_UNICAST          = linux.RTN_UNICAST
	RT_TABLE_MAIN        = linux.RT_TABLE_MAIN
	RTPROT_BOOT          = linux.RTPROT_BOOT
	RTPROT_STATIC        = linux.RTPROT_STATIC
	RT_SCOPE_UNIVERSE    = linux.RT_SCOPE_UNIVERSE
	RT_SCOPE_HOST        = linux.RT_SCOPE_HOST
	RT_SCOPE_LINK        = linux.RT_SCOPE_LINK
)
