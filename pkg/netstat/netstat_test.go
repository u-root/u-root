// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat_test

import (
	"errors"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/netstat"
)

func TestOutputNewOutput(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fmt     netstat.FmtFlags
		rootreq bool
		experr  error
	}{
		{
			name:    "SuccessNoProgNames",
			fmt:     netstat.FmtFlags{},
			rootreq: false,
			experr:  nil,
		},
		{
			name:    "SuccessProgNamesRoot",
			fmt:     netstat.FmtFlags{ProgNames: true},
			rootreq: true,
			experr:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.rootreq {
				guest.SkipIfNotInVM(t)
			}
			_, err := netstat.NewOutput(tt.fmt)
			if !errors.Is(err, tt.experr) {
				t.Error(err)
			}
		})
	}
}

func TestOutputInitIPSocketTitel(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fmt     netstat.FmtFlags
		rootreq bool
		experr  error
	}{
		{
			name:    "SuccessAllFlagsFalse",
			fmt:     netstat.FmtFlags{},
			rootreq: false,
			experr:  nil,
		},
		{
			name:    "SuccessProgNameSet",
			fmt:     netstat.FmtFlags{ProgNames: true},
			rootreq: true,
			experr:  nil,
		},
		{
			name:    "SuccessExtendSet",
			fmt:     netstat.FmtFlags{Extend: true},
			rootreq: false,
			experr:  nil,
		},
		{
			name:    "SuccessTimerSet",
			fmt:     netstat.FmtFlags{Timer: true},
			rootreq: false,
			experr:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.rootreq {
				guest.SkipIfNotInVM(t)
			}

			out, err := netstat.NewOutput(tt.fmt)
			if !errors.Is(err, tt.experr) {
				t.Error(err)
			}

			out.InitIPSocketTitels()

			if !strings.Contains(out.String(), "Proto") {
				t.Error("formatted output string does not contain 'Proto'")
			}

			if !strings.Contains(out.String(), "Recv-Q") {
				t.Error("formatted output string does not contain 'Recv-Q'")
			}

			if !strings.Contains(out.String(), "Send-Q") {
				t.Error("formatted output string does not contain 'Send-Q'")
			}

			if !strings.Contains(out.String(), "Local Address") {
				t.Error("formatted output string does not contain 'Local Address'")
			}

			if !strings.Contains(out.String(), "Foreign Address") {
				t.Error("formatted output string does not contain 'Foreign Address'")
			}

			if !strings.Contains(out.String(), "State") {
				t.Error("formatted output string does not contain 'State'")
			}

			if tt.fmt.Extend {
				if !strings.Contains(out.String(), "User") {
					t.Error("formatted output string does not contain 'User'")
				}

				if !strings.Contains(out.String(), "Inode") {
					t.Error("formatted output string does not contain 'Inode'")
				}
			}

			if tt.fmt.Timer && !strings.Contains(out.String(), "Timer") {
				t.Error("formatted output string does not contain 'Timer'")
			}

			if tt.fmt.ProgNames && !strings.Contains(out.String(), "PID/Program name") {
				t.Error("formatted output string does not contain 'PID/Program name'")
			}

			out.Builder.Reset()
		})
	}
}

func TestOutputInitUnixSocketTitels(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fmt     netstat.FmtFlags
		rootreq bool
		expErr  error
	}{
		{
			name:    "SuccessNoFlagsSet",
			fmt:     netstat.FmtFlags{},
			rootreq: false,
			expErr:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.rootreq {
				guest.SkipIfNotInVM(t)
			}
			out, err := netstat.NewOutput(tt.fmt)
			if err != nil {
				t.Error(err)
			}

			out.InitUnixSocketTitels()

			if !strings.Contains(out.String(), "Proto") {
				t.Error("formatted output string does not contain 'Proto'")
			}

			if !strings.Contains(out.String(), "RefCnt") {
				t.Error("formatted output string does not contain 'RefCnt'")
			}

			if !strings.Contains(out.String(), "Flags") {
				t.Error("formatted output string does not contain 'Flags'")
			}

			if !strings.Contains(out.String(), "Type") {
				t.Error("formatted output string does not contain 'Type'")
			}

			if !strings.Contains(out.String(), "State") {
				t.Error("formatted output string does not contain 'State'")
			}

			if !strings.Contains(out.String(), "I-Node") {
				t.Error("formatted output string does not contain 'I-Node'")
			}

			if tt.fmt.ProgNames && !strings.Contains(out.String(), "PID/Program name") {
				t.Error("formatted output string does not contain 'PID/Program name'")
			}
		})
	}
}

func TestOutputAddIPSocket(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fmt     netstat.FmtFlags
		rootreq bool
		experr  error
	}{
		{
			name:    "SuccessAllFlagsFalse",
			fmt:     netstat.FmtFlags{},
			rootreq: false,
			experr:  nil,
		},
		{
			name:    "SuccessProgNameSet",
			fmt:     netstat.FmtFlags{ProgNames: true},
			rootreq: true,
			experr:  nil,
		},
		{
			name:    "SuccessExtendSet",
			fmt:     netstat.FmtFlags{Extend: true},
			rootreq: false,
			experr:  nil,
		},
		{
			name: "SuccessExtendNumUser",
			fmt: netstat.FmtFlags{
				Extend:   true,
				NumUsers: true,
			},
			rootreq: false,
			experr:  nil,
		},
		{
			name:    "SuccessTimerSet",
			fmt:     netstat.FmtFlags{Timer: true},
			rootreq: false,
			experr:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.rootreq {
				guest.SkipIfNotInVM(t)
			}

			out, err := netstat.NewOutput(tt.fmt)
			if !errors.Is(err, tt.experr) {
				t.Error(err)
			}

			out.InitIPSocketTitels()

			sock, err := netstat.NewSocket(netstat.PROT_TCP)
			if err != nil {
				t.Error(err)
			}

			_, err = sock.SocketsString(false, false, out)
			if !errors.Is(err, tt.experr) {
				t.Error(err)
			}
		})
	}
}

func TestOutputAddUnixSocket(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fmt     netstat.FmtFlags
		rootreq bool
		expErr  error
	}{
		{
			name:    "SuccessNoFlagsSet",
			fmt:     netstat.FmtFlags{},
			rootreq: false,
			expErr:  nil,
		},
		{
			name:    "SuccessProgNameFlagsSet",
			fmt:     netstat.FmtFlags{ProgNames: true},
			rootreq: true,
			expErr:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if tt.rootreq {
				guest.SkipIfNotInVM(t)
			}
			out, err := netstat.NewOutput(tt.fmt)
			if err != nil {
				t.Error(err)
			}

			out.InitUnixSocketTitels()

			sock, err := netstat.NewSocket(netstat.PROT_UNIX)
			if err != nil {
				t.Error(err)
			}

			_, err = sock.SocketsString(false, false, out)
			if !errors.Is(err, tt.expErr) {
				t.Error(err)
			}
		})
	}
}

func TestConstructTimer(t *testing.T) {
	for _, tt := range []struct {
		name  string
		state uint8
		tl    uint64
		retr  uint64
		to    uint64
	}{
		{
			name:  "TimerOff",
			state: 0,
			tl:    0,
			retr:  0,
			to:    0,
		},
		{
			name:  "TimerOn",
			state: 1,
			tl:    0,
			retr:  0,
			to:    0,
		},
		{
			name:  "TimerKeepAlive",
			state: 2,
			tl:    0,
			retr:  0,
			to:    0,
		},
		{
			name:  "TimerTimeWait",
			state: 3,
			tl:    0,
			retr:  0,
			to:    0,
		},
		{
			name:  "TimerProbe",
			state: 4,
			tl:    0,
			retr:  0,
			to:    0,
		},
		{
			name:  "TimerUnknown",
			state: 5,
			tl:    0,
			retr:  0,
			to:    0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			out, err := netstat.NewOutput(netstat.FmtFlags{})
			if err != nil {
				t.Error(err)
			}

			out.ConstructTimer(tt.state, tt.tl, tt.retr, tt.to)

			out.Builder.Reset()
		})
	}
}

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
	grp := netstat.Groups{
		{
			IFace: "eth0",
			Grp:   net.IPv4(224, 0, 0, 1),
			Users: 5,
		},
		{
			IFace: "eth0",
			Grp:   net.ParseIP("ff02::1"),
			Users: 2,
		},
	}

	str := grp.String()

	if !(len(str) > 0) {
		t.Error("length of groups string is zero")
	}

	if !strings.Contains(str, "224.0.0.1") || !strings.Contains(str, "ff02::1") {
		t.Error("formatted output string does not contain expected group addresses")
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
		NumUsers: true,
	})
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

	_, err = ipv4.RoutesFormatString(false)
	if err != nil {
		t.Error(err)
	}

	var out strings.Builder

	if err := ipv4.PrintStatistics(&out); err != nil {
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

	_, err = ipv6.RoutesFormatString(false)
	if err != nil {
		t.Error(err)
	}

	var out strings.Builder

	if err := ipv6.PrintStatistics(&out); err != nil {
		t.Error(err)
	}
}

func TestPrintInterfaceTable(t *testing.T) {
	if _, err := os.Open(netstat.ProcNetDevPath); err != nil {
		t.Log("file dependencies not satisfied")
		t.Skip()
	}

	var out strings.Builder

	if err := netstat.PrintInterfaceTable("", false, &out); err != nil {
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

	var out strings.Builder

	if err := netstat.PrintMulticastGroups(true, true, &out); err != nil {
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
		NumUsers: true,
	})
	if err != nil {
		t.Error(err)
	}

	for _, s := range sock {
		_, err := s.SocketsString(false, false, output)
		if err != nil {
			t.Error(err)
		}
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
