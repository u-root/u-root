// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import "errors"

var ErrRouteCacheIPv6only = errors.New("route cache printing for IPv6 only")

type Protocol string

func (p *Protocol) String() string {
	return string(*p)
}

const (
	PROT_IGMP  Protocol = "igmp"
	PROT_IGMP6 Protocol = "igmp6"
	PROT_TCP   Protocol = "tcp"
	PROT_TCP6  Protocol = "tcp6"
	PROT_UDP   Protocol = "udp"
	PROT_UDP6  Protocol = "udp6"
	PROT_UDPL  Protocol = "udplite"
	PROT_UDPL6 Protocol = "udplite6"
	PROT_RAW   Protocol = "raw"
	PROT_RAW6  Protocol = "raw6"
	PROT_UNIX  Protocol = "unix"
)

var ProcnetPath = "/proc/net"

type NetState uint8

const (
	TCP_ESTABLISHED NetState = 1 + iota
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
)

func (n *NetState) String() string {
	switch *n {
	case TCP_ESTABLISHED:
		return "ESTABLISHED"
	case TCP_SYN_SENT:
		return "SYN_SENT"
	case TCP_SYN_RECV:
		return "SYN_RECV"
	case TCP_FIN_WAIT1:
		return "FIN_WAIT1"
	case TCP_FIN_WAIT2:
		return "FIN_WAIT2"
	case TCP_TIME_WAIT:
		return "TIME_WAIT"
	case TCP_CLOSE:
		return "CLOSE"
	case TCP_CLOSE_WAIT:
		return "CLOSE_WAIT"
	case TCP_LAST_ACK:
		return "LAST_ACK"
	case TCP_LISTEN:
		return "LISTEN"
	case TCP_CLOSING:
		return "CLOSING"
	default:
		return "unknown state"
	}
}

type SockState uint8

const (
	SSFREE = 0 + iota
	SSUNCONNECTED
	SSCONNECTING
	SSCONNECTED
	SSDISCONNECTING
)

func (s *SockState) parseState(flag uint32) string {
	switch *s {
	case SSFREE:
		return "FREE"
	case SSUNCONNECTED:
		if flag&SSACCEPTCON > 0 {
			return "LISTENING"
		}
		return "UNCONNECTED"
	case SSCONNECTING:
		return "CONNECTING"
	case SSCONNECTED:
		return "CONNECTED"
	case SSDISCONNECTING:
		return "DISCONNECTING"
	default:
		return ""
	}
}
