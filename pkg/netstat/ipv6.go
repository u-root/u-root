// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type IPv6 struct {
	*Output
}

func (i *IPv6) PrintStatistics() error {
	snmp6, err := newSNMP6()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", snmp6.String())

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
				//We skip UdpLite for now
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
