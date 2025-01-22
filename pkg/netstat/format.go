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
			return nil, fmt.Errorf("reading /proc filesystem: %w", err)
		}
		ret.ProgCache = cache
	}

	ret.Builder = strings.Builder{}
	ret.FmtFlags = flags

	return ret, nil
}

func (o *Output) InitIPSocketTitels() {
	fmt.Fprintf(&o.Builder, "%-6s %-6s %-6s %-35s %-35s %-17s",
		"Proto",
		"Recv-Q",
		"Send-Q",
		"Local Address",
		"Foreign Address",
		"State",
	)

	if o.Extend {
		fmt.Fprintf(&o.Builder, "%-9s %-11s", "User", "Inode")
	}

	if o.ProgNames {
		fmt.Fprintf(&o.Builder, "%-27s", "PID/Program name")
	}

	if o.Timer {
		fmt.Fprintf(&o.Builder, "%-10s", "Timer")
	}

	fmt.Fprint(&o.Builder, "\n")
}

func (o *Output) InitUnixSocketTitels() {
	fmt.Fprintf(&o.Builder, "%s\n", "Active sockets in the UNIX domain")
	fmt.Fprintf(&o.Builder, "%-8s%-8s%-8s%-10s%-16s%-10s",
		"Proto",
		"RefCnt",
		"Flags",
		"Type",
		"State",
		"I-Node",
	)

	if o.ProgNames {
		fmt.Fprintf(&o.Builder, "%-27s", "PID/Program name")
	}

	fmt.Fprintf(&o.Builder, "%s\n", "Path")
}

func (o *Output) AddIPSocket(ipsock netSocket) {
	fmt.Fprintf(&o.Builder, "%-6s %6d %6d %-35s %-35s %-17s",
		ipsock.Protocol,
		ipsock.RxQueue,
		ipsock.TxQueue,
		o.resolveAddress(ipsock.LocalAddr, net.CIDRMask(6, 8)),   // TODO: Transform in correct representation
		o.resolveAddress(ipsock.ForeignAddr, net.CIDRMask(6, 8)), // TODO: Transform in correct representation
		ipsock.State.String(),
	)

	if o.Extend {
		if o.NumUsers {
			fmt.Fprintf(&o.Builder, "%-10d", ipsock.UID)
		} else {
			username, err := convertUID(ipsock.UID)
			if err != nil {
				fmt.Fprintf(&o.Builder, "%-10d", ipsock.UID)
			}
			fmt.Fprintf(&o.Builder, "%-10s", username)
		}
		fmt.Fprintf(&o.Builder, "%-11d", ipsock.Inode)
	}

	if o.ProgNames {
		name, err := o.getNameFromInode(uint64(ipsock.Inode))
		if err != nil {
			fmt.Fprintf(&o.Builder, "%-27v", err)
		}
		fmt.Fprintf(&o.Builder, "%-27s", name)
	}

	if o.Timer {
		fmt.Fprintf(&o.Builder, "%-15s",
			o.ConstructTimer(ipsock.TimerRun,
				ipsock.TimerLen,
				ipsock.Retr,
				ipsock.Timeout))
	}

	fmt.Fprint(&o.Builder, "\n")
}

func (o *Output) AddUnixSocket(uSock unixSocket) {
	fmt.Fprintf(&o.Builder, "%-8s%-8d%-8s%-10s%-16s%-10d",
		"unix",
		uSock.RefCnt,
		parseUnixFlags(uSock.Flags),
		uSock.Type.String(),
		uSock.St.parseState(uSock.Flags),
		uSock.Inode,
	)

	if o.ProgNames {
		name, err := o.getNameFromInode(uint64(uSock.Inode))
		if err != nil {
			fmt.Fprintf(&o.Builder, "%-27v", err)
		}
		fmt.Fprintf(&o.Builder, "%-27s", name)
	}

	fmt.Fprintf(&o.Builder, "%s\n", uSock.Path)
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
		s.WriteString("/")
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

func (o *Output) ConstructTimer(tr uint8, tl, retr, to uint64) string {
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

func (o *Output) InitRoute4Titel() {
	fmt.Fprintf(&o.Builder, "%s\n", "Kernel IP routing table")
	fmt.Fprintf(&o.Builder, "%-16s %-16s %-16s %-8s %s %-7s %s %s\n",
		"Destination",
		"Gateway",
		"Genmask",
		"Flags",
		"MSS",
		"Window",
		"irrt",
		"Iface",
	)
}

func (o *Output) AddRoute4(r routev4) {
	var gw []string
	var dest string
	var err error

	if o.NumHosts {
		dest = r.Dest.String()
		gw = append(gw, r.Gateway.String())
	} else {
		gw, err = net.LookupAddr(r.Gateway.String())
		if err != nil {
			gw = append(gw, r.Gateway.String())
		}

		dest = o.resolveAddress(IPAddress{Address: r.Dest}, net.IPv4Mask(0, 0, 0, 0))
		if len(dest) < 1 {
			dest = r.Dest.String()
		}
	}

	fmt.Fprintf(&o.Builder, "%-16s %-16s %-16s %-8s %-3v %-7d %4d %s\n",
		dest,
		gw[0],
		r.Mask,
		r.Flags,
		r.MTU,
		r.Window,
		r.IRRT,
		r.IFace,
	)
}

func (o *Output) InitRoute6Titel() {
	fmt.Fprintf(&o.Builder, "%s\n", "Kernel IPv6 routing table")
	fmt.Fprintf(&o.Builder, "%-34s %-32s %-8s %-10s %-8s %-8s %s\n",
		"Destination",
		"Next Hop",
		"Flag",
		"Met",
		"Ref",
		"Use",
		"If",
	)
}

func (o *Output) AddRoute6(r routev6) {
	var dest, d, nextH []string
	var err error

	if o.NumHosts {
		if r.Dest.String() == "::" {
			d = append(d, "[::]")
		} else {
			d = append(d, r.Dest.String())
		}

		if r.NextHop.String() == "::" {
			nextH = append(nextH, "[::]")
		} else {
			nextH = append(nextH, r.NextHop.String())
		}
	} else {
		dest, err = net.LookupAddr(r.Dest.String())
		if err != nil {
			if r.Dest.String() == "::" {
				dest = append(dest, "[::]")
			} else {
				dest = append(dest, r.Dest.String())
			}
			d = append(d, fmt.Sprintf("%s/%d", dest[0], r.DestPrefix))
		} else {
			d = append(d, fmt.Sprintf("%s/%d", dest[0], r.DestPrefix))
		}

		nextH, err = net.LookupAddr(r.NextHop.String())
		if err != nil {
			if r.NextHop.String() == "::" {
				nextH = append(nextH, "[::]")
			} else {
				nextH = append(nextH, r.NextHop.String())
			}
		}
	}

	metstr := strconv.FormatUint(uint64(r.Metric), 10)

	if r.Metric == 0xFFFFFFFF {
		metstr = "-1"
	}

	fmt.Fprintf(&o.Builder, "%-34s %-32s %-8s %-10s %-8d %-8d %s\n",
		d[0],
		nextH[0],
		r.Flags,
		metstr,
		r.RefCnt,
		r.Use,
		r.IFace,
	)
}
