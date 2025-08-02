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
	"reflect"
	"strconv"
	"strings"
)

type IPv6 struct {
	*Output
}

var ProcNetRoutePath6 = "/proc/net/ipv6_route"

func (i *IPv6) RoutesFormatString(cache bool) (string, error) {
	file, err := os.Open(ProcNetRoutePath6)
	if err != nil {
		return "", err
	}

	i.Output.InitRoute6Titel()
	s := bufio.NewScanner(file)

	for s.Scan() {
		r, err := parseRoutev6(s.Text())
		if err != nil {
			return "", fmt.Errorf("failed to parse route: %w", err)
		}
		if cache {
			if strings.Contains(r.Flags, "C") {
				i.Output.AddRoute6(*r)
			}
		} else {
			i.Output.AddRoute6(*r)
		}
	}
	return i.Output.String(), nil
}

type routev6 struct {
	Dest         net.IP
	DestPrefix   uint32
	Source       net.IP
	SourcePrefix uint32
	NextHop      net.IP
	Metric       uint32
	RefCnt       uint32
	Use          uint32
	Flags        string
	IFace        string
}

func parseRoutev6(line string) (*routev6, error) {
	r := &routev6{}
	var dest, destprefix, src, srcprefix, nHop string
	var flag uint32

	_, err := fmt.Sscanf(line, "%s %s %s %s %s %x %x %d %x %s",
		&dest,
		&destprefix,
		&src,
		&srcprefix,
		&nHop,
		&r.Metric,
		&r.RefCnt,
		&r.Use,
		&flag,
		&r.IFace,
	)
	if err != nil {
		fmt.Printf("%v\n", line)
		return nil, err
	}

	destip, err := newIPAddress(dest)
	if err != nil {
		return nil, err
	}

	r.Dest = destip.Address

	d, err := strconv.ParseUint(destprefix, 16, 32)
	if err != nil {
		return nil, err
	}
	r.DestPrefix = uint32(d)

	srcip, err := newIPAddress(src)
	if err != nil {
		return nil, err
	}

	r.Source = srcip.Address

	d, err = strconv.ParseUint(srcprefix, 16, 32)
	if err != nil {
		return nil, err
	}
	r.SourcePrefix = uint32(d)

	nHopip, err := newIPAddress(nHop)
	if err != nil {
		return nil, err
	}

	r.NextHop = nHopip.Address

	r.Flags = parseIPv6Flags(flag)

	return r, nil
}

func parseIPv6Flags(flags uint32) string {
	var s strings.Builder

	f := []struct {
		flag uint32
		n    string
	}{
		{RTFUP, "U"},
		{RTFREJECT, "!"},
		{RTFGATEWAY, "G"},
		{RTFHOST, "H"},
		{RTFDEFAULT, "D"},
		{RTFADDRCONF, "A"},
		{RTFCACHE, "C"},
		{RTFALLONLINK, "a"},
		{RTFEXPIRES, "e"},
		{RTFMODIFIED, "m"},
		{RTFNONEXTHOP, "n"},
		{RTFFLOW, "f"},
	}

	for _, f := range f {
		if (f.flag & flags) > 0 {
			s.WriteString(f.n)
		}
	}

	return s.String()
}

func (i *IPv6) PrintStatistics(out io.Writer) error {
	snmp6, err := newSNMP6()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", snmp6.String())

	return nil
}

type SNMP6 struct {
	IP6   ip
	ICMP6 icmp
	UDP6  udp
	UDPL6 udp
	TCP6  tcp
}

func (s *SNMP6) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "%s\n", s.IP6.String())
	fmt.Fprintf(&str, "%s\n", s.ICMP6.String())
	fmt.Fprintf(&str, "%s\n", s.UDP6.String())
	fmt.Fprintf(&str, "%s\n", s.TCP6.String())

	return str.String()
}

var (
	prefixes  = []string{"Ip6", "Icmp6", "Udp6", "UdpLite6"}
	SNMP6file = "/proc/net/snmp6"
)

func newSNMP6() (*SNMP6, error) {
	i := &ip{}
	ic := &icmp{}
	u := &udp{}
	ul := &udp{}

	file, err := os.Open(SNMP6file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	s := bufio.NewScanner(file)

	for s.Scan() {
		// Each line contains the Prefix (Ip6, Icmp6, Udp6, UdpLite6) to each key-value pair

		for _, prefix := range prefixes {
			cutstring, ok := strings.CutPrefix(s.Text(), prefix)
			if !ok {
				continue
			}

			var name, val string
			fmt.Sscanf(cutstring, "%s %s", &name, &val)

			var refVal reflect.Value
			switch prefix {
			case "Ip6":
				refVal = reflect.ValueOf(i)
			case "Icmp6":
				refVal = reflect.ValueOf(ic)
			case "Udp6":
				refVal = reflect.ValueOf(u)
			case "UdpLite6":
				refVal = reflect.ValueOf(ul)
			default:
				// We skip UdpLite for now
				continue
			}

			field := refVal.Elem().FieldByName(name)
			if !field.CanSet() {
				continue
			}
			field.SetString(val)
		}
	}

	nsfile, err := os.Open(SNMP4file)
	if err != nil {
		return nil, err
	}

	s = bufio.NewScanner(nsfile)
	tcp := &tcp{}
	for s.Scan() {

		titleLine := s.Bytes()
		titleSlice := strings.Split(string(titleLine), ": ")

		s.Scan()
		dataLine := s.Bytes()
		// First element (dataSlice[0]) will hold the Prefix
		dataSlice := strings.Split(string(dataLine), ": ")

		switch dataSlice[0] {
		case "Tcp":
			parseData(titleSlice, dataSlice, reflect.ValueOf(tcp))
		}
	}
	return &SNMP6{IP6: *i, ICMP6: *ic, UDP6: *u, UDPL6: *ul, TCP6: *tcp}, nil
}

func (i *IPv6) ClearOutput() {
	i.Output.Builder.Reset()
}
