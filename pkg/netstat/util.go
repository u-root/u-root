// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type IPAddress struct {
	Address net.IP
	Port    uint16
}

func (s *IPAddress) String() string {
	return fmt.Sprintf("%v:%d", s.Address, s.Port)
}

const (
	ipv4Len    = 8
	ipv6Len    = 32
	ipv6GrpLen = 4
)

func newIPAddress(addr string) (IPAddress, error) {
	retAddr := IPAddress{}
	var ip net.IP

	splitAddr := strings.Split(addr, ":")

	switch len(splitAddr[0]) {
	case ipv4Len:
		v, err := strconv.ParseUint(splitAddr[0], 16, 32)
		if err != nil {
			return retAddr, err
		}
		ip = make(net.IP, net.IPv4len)
		binary.LittleEndian.PutUint32(ip, uint32(v))
	case ipv6Len:
		ip = make(net.IP, net.IPv6len)
		addr := splitAddr[0]
		i, j := 0, 4
		for len(addr) != 0 {
			grpStr := addr[0:8]
			grp, err := strconv.ParseUint(grpStr, 16, 32)
			if err != nil {
				return retAddr, err
			}
			binary.LittleEndian.PutUint32(ip[i:j], uint32(grp))
			i, j = i+ipv6GrpLen, j+ipv6GrpLen
			addr = addr[8:]
		}
	default:
		return retAddr, errors.New("unknown ip address length")
	}

	var v uint64
	var err error

	if len(splitAddr) > 1 {
		v, err = strconv.ParseUint(splitAddr[1], 16, 16)
		if err != nil {
			return retAddr, err
		}
	}

	retAddr.Address = ip
	retAddr.Port = uint16(v)

	return retAddr, nil
}
