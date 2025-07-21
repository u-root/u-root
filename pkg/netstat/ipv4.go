// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strings"
)

type AddressFamily interface {
	RoutesFormatString(bool) (string, error)
	PrintStatistics(io.Writer) error
	ClearOutput()
}

func NewAddressFamily(ipv6Flag bool, output *Output) AddressFamily {
	if !ipv6Flag {
		return &IPv4{output}
	}
	return &IPv6{output}
}

type IPv4 struct {
	*Output
}

var ProcNetRoutePath4 = "/proc/net/route"

func (i *IPv4) RoutesFormatString(_ bool) (string, error) {
	i.Output.InitRoute4Titel()

	file, err := os.Open(ProcNetRoutePath4)
	if err != nil {
		return "", err
	}

	s := bufio.NewScanner(file)

	// Scan the first line, which contains only title texts
	s.Scan()

	// All subsequent lines contain routes
	for s.Scan() {
		r, err := parseRoutev4(s.Text())
		if err != nil {
			return "", fmt.Errorf("failed to parse route: %w", err)
		}

		i.Output.AddRoute4(*r)
	}

	return i.Output.String(), nil
}

type routev4 struct {
	IFace   string
	Dest    net.IP
	Gateway net.IP
	Flags   string
	RefCnt  uint32
	Use     uint32
	Metric  uint32
	Mask    net.IP
	MTU     uint32
	Window  uint32
	IRRT    uint32
}

func parseRoutev4(line string) (*routev4, error) {
	retr := &routev4{}
	var dest, gate, mask string
	var flag uint32

	_, err := fmt.Sscanf(line, "%s %s %s %d %d %d %d %s %d %d %d",
		&retr.IFace,
		&dest,
		&gate,
		&flag,
		&retr.RefCnt,
		&retr.Use,
		&retr.Metric,
		&mask,
		&retr.MTU,
		&retr.Window,
		&retr.IRRT,
	)
	if err != nil {
		return nil, nil
	}

	ip, err := newIPAddress(dest)
	if err != nil {
		return nil, err
	}

	retr.Dest = ip.Address

	ip, err = newIPAddress(gate)
	if err != nil {
		return nil, err
	}

	retr.Gateway = ip.Address

	genmask, err := newIPAddress(mask)
	if err != nil {
		return nil, err
	}

	retr.Mask = genmask.Address

	retr.Flags = convertFlagData(flag)

	return retr, nil
}

func (i *IPv4) PrintStatistics(out io.Writer) error {
	snmp, err := newSNMP()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", snmp.String())

	netstat, err := newNetstat()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", netstat.String())

	return nil
}

var errUnknownPrefix = errors.New("unknown prefix found")

func newSNMP() (*SNMP, error) {
	ret := &SNMP{}
	file, err := os.Open(SNMP4file)
	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		titleLine := scan.Bytes()
		titleSlice := strings.Split(string(titleLine), ": ")

		scan.Scan()
		dataLine := scan.Bytes()
		// First element (dataSlice[0]) will hold the Prefix
		dataSlice := strings.Split(string(dataLine), ": ")

		switch dataSlice[0] {
		case "Ip":
			i := &ip{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(i))
			ret.IP = *i
		case "Icmp":
			ic := &icmp{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(ic))
			ret.ICMP = *ic
		case "IcmpMsg":
			icmsg := &icmpmsg{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(icmsg))
			ret.ICMPMsg = *icmsg
		case "Tcp":
			t := &tcp{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(t))
			ret.TCP = *t
		case "Udp":
			u := &udp{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(u))
			ret.UDP = *u
		case "UdpLite":
			udpl := &udp{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(udpl))
			ret.UDPL = *udpl
		default:
			return nil, errUnknownPrefix
		}
	}

	return ret, nil
}

func newNetstat() (*NetStat, error) {
	ret := &NetStat{}
	file, err := os.Open(netstatfile)
	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		titleLine := string(scan.Bytes())
		titleSlice := strings.Split(titleLine, ": ")

		scan.Scan()
		dataLine := string(scan.Bytes())
		dataSlice := strings.Split(dataLine, ": ")

		switch dataSlice[0] {
		case "TcpExt":
			te := &tcpExt{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(te))
			ret.TCPExt = *te
		case "IpExt":
			ie := &ipExt{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(ie))
			ret.IPExt = *ie
		case "MPTcpExt":
			me := &mptcpext{}
			parseData(titleSlice, dataSlice, reflect.ValueOf(me))
			ret.MPTCPExt = *me
		default:
			return nil, errUnknownPrefix
		}
	}
	return ret, nil
}

func parseData(ts, ds []string, stc reflect.Value) {
	titleSlice := strings.Split(ts[1], " ")
	dataSlice := strings.Split(ds[1], " ")

	for i, title := range titleSlice {
		field := stc.Elem().FieldByName(title)
		if !field.CanSet() {
			continue
		}
		field.SetString(dataSlice[i])
	}
}

func (i *IPv4) ClearOutput() {
	i.Output.Builder.Reset()
}
