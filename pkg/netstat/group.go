// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var (
	ProcNetigmpv4path = "/proc/net/igmp"
	ProcNetigmpv6path = "/proc/net/igmp6"
)

type Groups []member

func (g *Groups) String() string {
	var s strings.Builder

	for _, mem := range *g {
		fmt.Fprintf(&s, "%-20s %-10d %s\n",
			mem.IFace,
			mem.Users,
			mem.Grp)
	}

	return s.String()
}

type member struct {
	IFace string
	Grp   net.IP
	Users uint32
}

func PrintMulticastGroups(ipv4, ipv6 bool, out io.Writer) error {
	g := Groups{}

	if ipv4 {
		members, err := parseigmp()
		if err != nil {
			return fmt.Errorf("failed to parse igmp: %w", err)
		}
		g = append(g, members...)
		fmt.Fprintf(out, "%s", "IPv4")
	}

	if ipv6 {
		members, err := parseigmp6()
		if err != nil {
			return fmt.Errorf("failed to parse igmp6: %w", err)
		}
		g = append(g, members...)
		if ipv4 {
			fmt.Fprintf(out, "%s", "/")
		}
		fmt.Fprintf(out, "%s", "IPv6")
	}

	fmt.Fprintf(out, "%s\n", " Group memberships")
	fmt.Fprintf(out, "%-20s %-10s %s\n", "Interface", "RefCnt", "Group")
	fmt.Fprintf(out, "%-20s %-10s %s\n", "------------------", "---------", "---------")

	fmt.Fprintf(out, "%s\n", g.String())

	return nil
}

func parseigmp() ([]member, error) {
	file, err := os.Open(ProcNetigmpv4path)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(file)

	// First line is worthless
	s.Scan()

	members := make([]member, 0)
	for s.Scan() {
		// Get rid of : in the line
		var idx, count int
		var iface, vstring string
		line := s.Text()
		fmt.Sscanf(line, "%d %s : %d %s",
			&idx,
			&iface,
			&count,
			&vstring,
		)

		for i := 0; i < count; i++ {
			m := member{
				IFace: iface,
			}
			s.Scan()
			entryline := s.Text()
			var mIP string
			var users, time, rep uint32
			fmt.Sscanf(entryline, "%s %d %s %d",
				&mIP,
				&users,
				&time,
				&rep,
			)

			ip, err := newIPAddress(mIP)
			if err != nil {
				return nil, err
			}

			m.Grp = ip.Address
			m.Users = users

			members = append(members, m)
		}
	}
	return members, nil
}

func parseigmp6() ([]member, error) {
	file, err := os.Open(ProcNetigmpv6path)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(file)

	retmem := make([]member, 0)
	for s.Scan() {
		m := member{}
		line := s.Text()
		var idx uint32
		var grpaddr string
		fmt.Sscanf(line, "%d %s %s %d",
			&idx,
			&m.IFace,
			&grpaddr,
			&m.Users,
		)

		ip, err := newIPAddress(grpaddr)
		if err != nil {
			return nil, err
		}

		m.Grp = ip.Address

		retmem = append(retmem, m)
	}
	return retmem, nil
}
