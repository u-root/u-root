// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"strconv"
	"strings"
)

type ProtocolOptions struct {
	IPType     IPType
	SocketType SocketType
}

type IPType int

const (
	IP_NONE IPType = iota
	IP_V4
	IP_V6
	IP_V4_V6
	IP_V4_STRICT
	IP_V6_STRICT
)

type SocketType int

// UDP can be combined with one of {UNIX, VSOCK} or stand alone
const (
	SOCKET_TYPE_TCP SocketType = iota
	SOCKET_TYPE_UDP
	SOCKET_TYPE_UNIX
	SOCKET_TYPE_VSOCK
	SOCKET_TYPE_SCTP
	SOCKET_TYPE_UDP_VSOCK
	SOCKET_TYPE_UDP_UNIX
	SOCKET_TYPE_NONE
)

func (s SocketType) String() string {
	return [...]string{
		"tcp",
		"udp",
		"unix",
		"vsock",
		"sctp",
		"udp-vsock",
		"unixgram",
		"none",
	}[s]
}

func (p *ProtocolOptions) Network() (string, error) {
	switch p.SocketType {
	case SOCKET_TYPE_TCP:
		switch p.IPType {
		case IP_V4, IP_V4_STRICT:
			return "tcp4", nil
		case IP_V6, IP_V6_STRICT:
			return "tcp6", nil
		default:
			return "tcp", nil
		}
	case SOCKET_TYPE_SCTP:
		switch p.IPType {
		case IP_V4, IP_V4_STRICT:
			return "sctp4", nil
		case IP_V6, IP_V6_STRICT:
			return "sctp6", nil
		default:
			return "sctp", nil
		}
	case SOCKET_TYPE_UDP:
		switch p.IPType {
		case IP_V4, IP_V4_STRICT:
			return "udp4", nil
		case IP_V6, IP_V6_STRICT:
			return "udp6", nil
		default:
			return "udp", nil
		}

	case SOCKET_TYPE_UNIX:
		return "unix", nil

	case SOCKET_TYPE_UDP_UNIX:
		return "unixgram", nil

	// VSOCK connections don't require a network specification
	case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK:
		return "", nil
	}

	return "", fmt.Errorf("invalid/unimplemented combination of socket and ip type (%v - %v)", p.SocketType, p.IPType)
}

func ParseSocketType(udp, unix, vsock, sctp bool) (SocketType, error) {
	// tcp ^ (udp || udp && unix || udp && vsock) ^ unix ^ vsock ^ sctp
	if !(udp || unix || vsock || sctp) {
		return SOCKET_TYPE_TCP, nil
	}

	if udp && !(unix || vsock || sctp) {
		return SOCKET_TYPE_UDP, nil
	}

	if udp && (unix != vsock) && !sctp {
		if unix {
			return SOCKET_TYPE_UDP_UNIX, nil
		} else if vsock {
			return SOCKET_TYPE_UDP_VSOCK, nil
		}
	}

	if unix && !(udp || vsock || sctp) {
		return SOCKET_TYPE_UNIX, nil
	}

	if vsock && !(udp || unix || sctp) {
		return SOCKET_TYPE_VSOCK, nil
	}

	if sctp && !(udp || unix || vsock) {
		return SOCKET_TYPE_SCTP, nil
	}

	return SOCKET_TYPE_NONE, fmt.Errorf("invalid socket type combination")
}

func SplitVSockAddr(addr string) (uint32, uint32, error) {
	splitAddr := strings.Split(addr, ":")

	if len(splitAddr) != 2 {
		return 0, 0, fmt.Errorf("invalid vsock address %q", addr)
	}

	cid, err := strconv.ParseUint(splitAddr[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse vsock CID %q: %w", splitAddr[0], err)
	}

	port, err := strconv.ParseUint(splitAddr[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse vsock port %q: %w", splitAddr[1], err)
	}

	return uint32(cid), uint32(port), nil
}
