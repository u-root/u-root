// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat_test

import (
	"net"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/netstat"
)

func TestProtocolString(t *testing.T) {
	tcpProt := netstat.PROT_TCP
	tcpString := "tcp"

	if tcpProt.String() != tcpString {
		t.Errorf("%s not equal %s", tcpProt.String(), tcpString)
	}
}

func TestStateString(t *testing.T) {
	for _, tt := range []struct {
		input uint8
		exp   string
	}{
		{
			input: 0,
			exp:   "unknown state",
		},
		{
			input: 1,
			exp:   "ESTABLISHED",
		},
		{
			input: 2,
			exp:   "SYN_SENT",
		},
		{
			input: 3,
			exp:   "SYN_RECV",
		},
		{
			input: 4,
			exp:   "FIN_WAIT1",
		},
		{
			input: 5,
			exp:   "FIN_WAIT2",
		},
		{
			input: 6,
			exp:   "TIME_WAIT",
		},
		{
			input: 7,
			exp:   "CLOSE",
		},
		{
			input: 8,
			exp:   "CLOSE_WAIT",
		},
		{
			input: 9,
			exp:   "LAST_ACK",
		},
		{
			input: 10,
			exp:   "LISTEN",
		},
		{
			input: 11,
			exp:   "CLOSING",
		},
	} {
		t.Run(tt.exp, func(t *testing.T) {
			tt := tt
			s := netstat.NetState(tt.input)

			str := s.String()

			if !(len(str) > 0) {
				t.Error("len of state is zero")
			}

			if str != tt.exp {
				t.Errorf("s.String() = %s, not %s", s.String(), tt.exp)
			}
		})
	}
}

func TestIPAddressString(t *testing.T) {
	ipv4 := netstat.IPAddress{
		Address: net.IPv4(192, 168, 0, 1),
		Port:    12345,
	}

	ipv6 := netstat.IPAddress{
		Address: net.IPv6loopback,
		Port:    0,
	}

	ipv6str := ipv6.String()

	if !(len(ipv6str) > 0) {
		t.Error("string length of IPv6 address is zero")
	}

	ipstr := ipv4.String()

	if !(len(ipstr) > 0) {
		t.Error("string length of IPv4 address is zero")
	}
}

func TestSMNPString(t *testing.T) {
	snmp := netstat.SNMP{}

	str := snmp.String()

	if !(len(str) > 0) {
		t.Error("string of snmp structure is zero")
	}
}

func TestSNMP6String(t *testing.T) {
	snmp := netstat.SNMP6{}

	str := snmp.String()

	if !(len(str) > 0) {
		t.Error("string of snmp structure is zero")
	}
}

func TestAddressFamiliy(t *testing.T) {
	fmt, err := netstat.NewOutput(netstat.FmtFlags{})
	if err != nil {
		t.Error(err)
	}
	_ = netstat.NewAddressFamily(false, fmt)
	_ = netstat.NewAddressFamily(true, fmt)

}

func TestGroupString(t *testing.T) {
	grp := netstat.Groups{}

	str := grp.String()

	if !(len(str) > 0) {
		t.Error("length of groups string is zero")
	}
}

func TestNetstatString(t *testing.T) {
	nstat := netstat.NetStat{}

	str := nstat.String()

	if !(len(str) > 0) {
		t.Error("length of netstat string is zero")
	}
}

func TestIPv4(t *testing.T) {
	fmt, err := netstat.NewOutput(netstat.FmtFlags{
		NumHosts: true,
		NumPorts: true,
		NumUsers: true})
	if err != nil {
		t.Error(err)
	}
	ipv4 := netstat.NewAddressFamily(false, fmt)

	deps := []string{
		netstat.ProcNetRoutePath4,
		netstat.SNMP4file,
	}

	if !checkFileDependencies(deps...) {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	if err := ipv4.PrintRoutes(false, false); err != nil {
		t.Error(err)
	}

	if err := ipv4.PrintStatistics(); err != nil {
		t.Error(err)
	}
}

func TestIPv6(t *testing.T) {
	fmt, err := netstat.NewOutput(netstat.FmtFlags{
		NumHosts: true,
		NumPorts: true,
		NumUsers: true,
	})
	if err != nil {
		t.Error(err)
	}
	ipv6 := netstat.NewAddressFamily(true, fmt)

	deps := []string{
		netstat.ProcNetRoutePath6,
		netstat.SNMP4file,
		netstat.SNMP6file,
	}

	if !checkFileDependencies(deps...) {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	if err := ipv6.PrintRoutes(false, false); err != nil {
		t.Error(err)
	}

	if err := ipv6.PrintStatistics(); err != nil {
		t.Error(err)
	}
}

func TestPrintInterfaceTable(t *testing.T) {
	if _, err := os.Open(netstat.ProcNetDevPath); err != nil {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	if err := netstat.PrintInterfaceTable("", false); err != nil {
		t.Error(err)
	}
}

func TestPrintMulticastGroups(t *testing.T) {
	deps := []string{
		netstat.ProcNetigmpv4path,
		netstat.ProcNetRoutePath6,
	}

	if !checkFileDependencies(deps...) {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	if err := netstat.PrintMulticastGroups(true, true); err != nil {
		t.Error(err)
	}
}

func TestPrintNetFiles(t *testing.T) {
	deps := []string{
		netstat.ProcnetPath + "/tcp",
		netstat.ProcnetPath + "/tcp6",
		netstat.ProcnetPath + "/udp",
		netstat.ProcnetPath + "/udp6",
	}

	if !checkFileDependencies(deps...) {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	sock := make([]netstat.Socket, 0)

	tcpsock, err := netstat.NewSocket(netstat.PROT_TCP)
	if err != nil {
		t.Error(err)
	}

	sock = append(sock, tcpsock)

	output, err := netstat.NewOutput(netstat.FmtFlags{
		NumHosts: true,
		NumPorts: true,
		NumUsers: true})
	if err != nil {
		t.Error(err)
	}

	for _, s := range sock {
		s.PrintSockets(false, false, output)
	}
}

func checkFileDependencies(files ...string) bool {
	for _, file := range files {
		if _, err := os.Open(file); err != nil {
			return false
		}
	}
	return true
}
