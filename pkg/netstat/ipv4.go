// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
)

type AddressFamily interface {
	PrintStatistics() error
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

func (i *IPv4) PrintStatistics() error {
	snmp, err := newSNMP()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", snmp.String())

	netstat, err := newNetstat()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", netstat.String())

	return nil
}

var (
	errUnknownPrefix = errors.New("unknown prefix found")
)

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
