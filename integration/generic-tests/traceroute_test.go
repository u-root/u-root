// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func TestTraceroute(t *testing.T) {
	// To deterministically test traceroute, we have to set up a controlled network
	// environment.
	//
	// ┌───────────────┐         ┌───────────────┐         ┌───────────────┐
	// │       A       │◄───────►│       B       │◄───────►│       C       |
	// │               │         │ 192.168.10.1  │         │               |
	// │ 192.168.10.10 │         │ 192.168.20.1  │         │ 192.168.20.10 │
	// │               │         │               │         │               │
	// │ listen:PortAB │         │listen:PortBC  │         │               │
	// │               │         │con:PortAB     │         │con:PortBC     │
	// └───────────────┘         └───────────────┘         └───────────────┘
	//
	// Host C is the target of the traceroute. Host A is the source of the traceroute.
	// Host B is the router. Host B will be running a netcat server that will respond
	// to the traceroute packets. Host A will be running the traceroute command.
	// To enforce this network typology and also keep ICMP compatible, we use
	// QEMU socket networking.

	const (
		hostAIPv4       = "192.168.10.10"
		hostBGatewayAv4 = "192.168.10.1"
		hostBGatewayCv4 = "192.168.20.1"
		hostCIPv4       = "192.168.20.10"
		hostsv4         = hostAIPv4 + "vmA\n" + hostCIPv4 + "vmC\n"

		hostAIPv6         = "2001:db8:ab::2"
		hostBGatewayAIPv6 = "2001:db8:ab::1"
		hostBGatewayCIPv6 = "2001:db8:bc::1"
		hostCIPv6         = "2001:db8:bc::2"

		qemuSocketPortAB = "1234"
		qemuSocketPortBC = "1235"
		qemuSocketHost   = "127.0.0.1"

		tcpTestPorts          = "80 443 8080 33434 31337"
		tcpTestRegexHostCIPv4 = "TTL: 2+[[:space:]]+" + hostCIPv4 + "+[[:space:]]+\\([0-9]+(\\.[0-9]+)?[[:space:]]*(ms|us|μs|s)\\)"
		tcpTestRegexHostCIPv6 = "TTL: 2+[[:space:]]+" + hostCIPv6 + "+[[:space:]]+\\([0-9]+(\\.[0-9]+)?[[:space:]]*(ms|us|μs|s)\\)"
	)

	var (
		scriptHostA, scriptHostB, scriptHostC strings.Builder
	)

	fmt.Fprint(&scriptHostA, `
		sleep 30

		# IPv4 setup
		ip addr add `+hostAIPv4+`/24 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip route add default via `+hostBGatewayAv4+` dev eth0 || exit 1

		# IPv6 setup
		ip -6 addr add `+hostAIPv6+`/64 dev eth0 || exit 1
		ip -6 route add default via `+hostBGatewayAIPv6+` dev eth0 || exit 1

		sleep 30
		# Test IPv4
		ping -c 1 `+hostCIPv4+` || exit 1

		# udp4
		# default traceroute with parameters explicitly set
		traceroute -4 -m udp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`"  || exit 1

		# traceroute udp custom port
		traceroute -4 -m udp -p 33434 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# traceroute udp must fail
		traceroute -m udp -p 111111 `+hostCIPv4+` || return 0 && exit 1

		# tcp4
		traceroute -m tcp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# traceroute tcp custom port
		traceroute -m tcp -p 8080 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# icmp4
		traceroute -m icmp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		echo "icmp ipv4 custom port / sequence number"
		traceroute -m icmp -p 2 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		exit 0

		# udp6
		traceroute -6 -m udp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# traceroute udp6 custom port
		traceroute -6 -m udp -p 33434 `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# tcp6
		traceroute -6 -m tcp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# traceroute tcp6 custom port
		traceroute -6 -m tcp -p 8080 `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# icmp6
		traceroute -6 -m icmp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1
		`)

	fmt.Fprint(&scriptHostB, `
		sleep 30
		# IPv4 setup
		ip addr add `+hostBGatewayAv4+`/24 dev eth0 || exit 1
		ip addr add `+hostBGatewayCv4+`/24 dev eth1 || exit 1
		ip link set eth0 up || exit 1
		ip link set eth1 up || exit 1
		echo "1" > /proc/sys/net/ipv4/ip_forward || exit 1

		# Ipv6 setup
		ip -6 addr add `+hostBGatewayAIPv6+`/64 dev eth0 || exit 1
		ip -6 addr add `+hostBGatewayCIPv6+`/64 dev eth1 || exit 1
		echo "1" > /proc/sys/net/ipv6/conf/all/forwarding || exit 1

		# configure nat
		ip link set eth2 up || exit 1
		ip addr add 10.0.2.15/24 dev eth2 || exit 1
		ip route add default via 10.0.2.2 dev eth2 || exit 1

		sleep 30

		# debugging
		ip a
		ip r

		# router does nothing
		sleep 200
	`)

	fmt.Fprint(&scriptHostC, `
		sleep 30
		# IPv4 setup
		ip addr add `+hostCIPv4+`/24 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip route add default via `+hostBGatewayCv4+` dev eth0 || exit 1

		# IPv6 setup
		ip -6 addr add `+hostCIPv6+`/64 dev eth0 || exit 1
		ip -6 route add default via `+hostBGatewayCIPv6+` dev eth0 || exit 1

		sleep 10
		# check connectivity
		ping -c 1 `+hostBGatewayCv4+` || exit 1

		# tcp4
		# keep the tcp servers open and running
		for port in `+tcpTestPorts+`; do
			echo $port
			netcat -l 192.168.20.10 $port &
		done

		sleep 300
	`)

	vmA := scriptvm.Start(t, "vmA", scriptHostA.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/core/ping",
				"github.com/u-root/u-root/cmds/core/grep",
				"github.com/u-root/u-root/cmds/core/tail",
			),
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,listen=:"+qemuSocketPortAB),
		),
	)

	vmB := scriptvm.Start(t, "vmB", scriptHostB.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/core/ping",
				// "github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect="+qemuSocketHost+":"+qemuSocketPortAB),
			qemu.ArbitraryArgs("-nic", "socket,listen=:"+qemuSocketPortBC),
			// NAT to the outside world
			qemu.ArbitraryArgs("-netdev", "user,id=net2"),
			qemu.ArbitraryArgs("-device", "e1000,netdev=net2"),
		),
	)

	vmC := scriptvm.Start(t, "vmC", scriptHostC.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/core/ping",
				"github.com/u-root/u-root/cmds/core/netcat",
				// "github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect="+qemuSocketHost+":"+qemuSocketPortBC),
		),
	)

	// go vmA.Wait()
	vmB.Wait()
	vmC.Wait()

	if _, err := vmA.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("tc_test: %v", err)
	}
	if _, err := vmB.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("tc_test: %v", err)
	}
	if _, err := vmC.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("tc_test: %v", err)
	}

	vmA.Wait()
}
