// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"fmt"
	"net"
	"os/user"
	"strconv"
	"strings"
)

type FmtFlags struct {
	Extend    bool // Adds fields User and Inode for ip sockets
	Wide      bool // Unknown
	NumHosts  bool // dont resolve host ip to names - applies to ip socks
	NumPorts  bool // dont resolve port numbers to names - applies to ip socks
	NumUsers  bool // dont resolve user id to usernames - applies to ip and unix socks
	ProgNames bool // Adds fields PID/Program name for sockets (for ip and unix sockets)
	Timer     bool // Adds field timer for sockets (not unix socket)
	Symbolic  bool // Rounting table/cache -> no ip to name conversion, same as -n
}

type Output struct {
	strings.Builder
	FmtFlags

	ProgCache map[int]ProcNode
}

func NewOutput(
	flags FmtFlags,
) (*Output, error) {
	ret := &Output{}

	if flags.ProgNames {
		cache, err := readProgFS()
		if err != nil {
			return nil, err
		}
		ret.ProgCache = cache
	}

	ret.Builder = strings.Builder{}
	ret.FmtFlags = flags

	return ret, nil
}


func convertUID(uid uint32) (string, error) {
	var s strings.Builder
	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return "", err
	}
	s.WriteString(user.Username)
	return s.String(), nil
}

func (o *Output) getNameFromInode(inode uint64) (string, error) {
	var s strings.Builder
	pnote := o.ProgCache[int(inode)]
	if pnote.PID == 0 {
		s.WriteString("-")
	} else {
		s.WriteString(strconv.Itoa(pnote.PID))
	}

	s.WriteString(pnote.Name)
	return s.String(), nil
}

// CLK_TCK is a constant on Linux for all architectures except alpha and ia64.
// See e.g.
// https://git.musl-libc.org/cgit/musl/tree/src/conf/sysconf.c#n30
// https://github.com/containerd/cgroups/pull/12
// https://lore.kernel.org/lkml/agtlq6$iht$1@penguin.transmeta.com/
const SystemClkTck = 100

func (o *Output) constructTimer(tr uint8, tl, retr, to uint64) string {
	var s strings.Builder
	clktick := SystemClkTck
	switch tr {
	case 0:
		fmt.Fprintf(&s, "off: (0.00/%d/%d)", retr, to)
	case 1:
		fmt.Fprintf(&s, "on: (%2.2f/%d/%d)", float64(tl)/float64(clktick), retr, to)
	case 2:
		fmt.Fprintf(&s, "keepalive: (%2.2f/%d/%d)", float64(tl)/float64(clktick), retr, to)
	case 3:
		fmt.Fprintf(&s, "timewait: (%2.2f/%d/%d)", float64(tl)/float64(clktick), retr, to)
	case 4:
		fmt.Fprintf(&s, "probe: (%2.2f/%d/%d)", float64(tl)/float64(clktick), retr, to)
	default:
		fmt.Fprintf(&s, "unknown: %d (%2.2f/%d/%d)", tr, float64(tl)/float64(clktick), retr, to)
	}
	return s.String()
}

var (
	ANYADDR = "0.0.0.0"
	ANYMASK = "00000000"
)

func (o *Output) resolveAddress(addr IPAddress, netmask net.IPMask) string {
	var s strings.Builder

	if o.NumHosts && o.NumPorts {
		s.WriteString(addr.String())
		return s.String()
	}

	switch len(addr.Address) {
	case 4:
		// IPv4 case
		// Always assume 0xFFFFFF00 as netmask?
		if addr.Address.String() == ANYADDR {
			if netmask.String() == ANYMASK {
				s.WriteString("default")
			} else {
				s.WriteString("*")
			}
			return s.String()
		}

		if addr.Address.Mask(netmask).String() != ANYADDR {
			hn, err := net.LookupAddr(addr.Address.String())
			if len(hn) > 0 {
				if err != nil {
					s.WriteString("unable to resolve")
				}
				if o.Wide {
					hncut, ok := strings.CutSuffix(hn[0], ".")
					if !ok {
						s.WriteString(hn[0])
					} else {
						s.WriteString(hncut)
					}
					s.WriteString(":" + strconv.Itoa(int(addr.Port)))
				} else {
					hnsplit := strings.Split(hn[0], ".")
					s.WriteString(hnsplit[0])
					s.WriteString(":" + strconv.Itoa(int(addr.Port)))
				}
			} else {
				s.WriteString(addr.String())
			}
		}
	case 16:
		loc, err := net.LookupAddr(addr.Address.String())
		if len(loc) > 0 {
			s.WriteString(loc[0])
			s.WriteString(":" + strconv.Itoa(int(addr.Port)))
		}
		if err != nil {
			s.WriteString(addr.String())
		}
	}

	return s.String()
}
