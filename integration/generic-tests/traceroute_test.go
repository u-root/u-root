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

		hostAIPv6         = "fd51:3681:1eb4::2"
		hostBGatewayAIPv6 = "fd51:3681:1eb4::1"
		hostBGatewayCIPv6 = "fd51:3681:1eb5::1"
		hostCIPv6         = "fd51:3681:1eb5::2"

		qemuSocketPortAB = "1234"
		qemuSocketPortBC = "1235"
		qemuSocketHost   = "127.0.0.1"

		tcpTestPorts          = "80 443 8080 33434 31337"
		tcpTestRegexHostCIPv4 = "TTL: 2+[[:space:]]+" + hostCIPv4 + "+[[:space:]]+\\([0-9]+(\\.[0-9]+)?[[:space:]]*(ms|us|μs|s)\\)"
		tcpTestRegexHostCIPv6 = "TTL: 2+[[:space:]]+" + hostCIPv6 + "+[[:space:]]+\\([0-9]+(\\.[0-9]+)?[[:space:]]*(ms|us|μs|s)\\)"
	)

	var scriptHostA, scriptHostB, scriptHostC strings.Builder

	fmt.Fprint(&scriptHostA, `
		sleep 30

		# IPv4 setup
		ip addr add `+hostAIPv4+`/24 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip route add default via `+hostBGatewayAv4+` dev eth0 || exit 1

		# IPv6 setup
		ip -6 addr add `+hostAIPv6+`/64 dev eth0 || exit 1
		ip -6 route add default via `+hostBGatewayAIPv6+` dev eth0 || exit 1

		echo "`+hostAIPv4+`	vmA" > /etc/hosts || exit 1
		echo "`+hostCIPv4+`	vmC" >> /etc/hosts || exit 1
		echo "`+hostAIPv6+`	vmA" >> /etc/hosts || exit 1
		echo "`+hostCIPv6+`	vmC" >> /etc/hosts || exit 1

		sleep 30
		# udp4
		# default traceroute with parameters explicitly set
		traceroute -4 -m udp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`"  || exit 1

		# traceroute udp custom port
		traceroute -4 -m udp -p 33434 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# tcp4
		traceroute -m tcp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# traceroute tcp custom port
		traceroute -m tcp -p 8080 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# icmp4
		traceroute -m icmp `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# icmp ipv4 custom port / sequence number
		traceroute -m icmp -p 2 `+hostCIPv4+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv4+`" || exit 1

		# finish up tcp4 tests
		echo "TCP4DONE" | netcat `+hostCIPv4+` 22222 || exit 1

		# wait for tcp6 servers to start
		echo "Waiting for TCP6 servers to start"
		sleep 10

		# udp6
		traceroute -6 -m udp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# traceroute udp6 custom port
		traceroute -6 -m udp -p 33434 `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# tcp6 default port
		traceroute -6 -m tcp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# traceroute tcp6 custom port
		traceroute -6 -m tcp -p 8080 `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# icmp6
		traceroute -6 -m icmp `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		# icmp6 custom port / sequence number
		traceroute -6 -m icmp -p 2 `+hostCIPv6+` | tail -n 1 | grep -q --regexp "`+tcpTestRegexHostCIPv6+`" || exit 1

		echo "ALL TESTS PASSED MARKER"

		# signal the other VMs that the tests are done
		echo "TESTDONE" | netcat `+hostBGatewayAv4+` 22222 || exit 1
		echo "TESTDONE" | netcat `+hostCIPv4+` 22222 || exit 1

		echo "ALL TESTS PASSED MARKER"
		`)

	fmt.Fprint(&scriptHostB, `
		sleep 30
		# IPv4 setup
		ip addr add `+hostBGatewayAv4+`/24 dev eth0 || exit 1
		ip addr add `+hostBGatewayCv4+`/24 dev eth1 || exit 1
		echo "1" > /proc/sys/net/ipv4/ip_forward || exit 1

		# Ipv6 setup
		ip -6 addr add `+hostBGatewayAIPv6+`/64 dev eth0 || exit 1
		ip -6 addr add `+hostBGatewayCIPv6+`/64 dev eth1 || exit 1

		echo "1" > /proc/sys/net/ipv6/conf/all/forwarding || exit 1
		echo "1" > /proc/sys/net/ipv6/conf/all/accept_ra || exit 1
		echo "2" > /proc/sys/net/ipv6/conf/all/router_solicitations || exit 1
		echo "1" > /proc/sys/net/ipv6/conf/all/accept_ra_defrtr || exit 1

		ip link set eth0 up || exit 1
		ip link set eth1 up || exit 1

		# configure nat
		ip link set eth2 up || exit 1
		ip addr add 10.0.2.15/24 dev eth2 || exit 1
		ip route add default via 10.0.2.2 dev eth2 || exit 1

		echo "`+hostAIPv4+`	vmA" > /etc/hosts
		echo "`+hostCIPv4+`	vmC" >> /etc/hosts
		echo "`+hostAIPv6+`	vmA" >> /etc/hosts
		echo "`+hostCIPv6+`	vmC" >> /etc/hosts

		# wait for done signal from vmA
		netcat -l `+hostBGatewayAv4+` 22222 | grep -q "TESTDONE" && echo "ALL TESTS PASSED MARKER" || exit 1 

		# idle
		#sleep 300
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

		# export GODEBUG=netdns=go+1

		echo "`+hostAIPv4+`	vmA" > /etc/hosts
		echo "`+hostCIPv4+`	vmC" >> /etc/hosts
		echo "`+hostAIPv6+`	vmA" >> /etc/hosts
		echo "`+hostCIPv6+`	vmC" >> /etc/hosts

		# tcp4
		# keep the tcp servers open and running
		# gosh bugs...
		#for port in `+tcpTestPorts+`; do
		#	netcat -l `+hostCIPv4+` $port &
		#done

		netcat -l -k `+hostCIPv4+` 80 &
		netcat -l -k `+hostCIPv4+` 443 &
		netcat -l -k `+hostCIPv4+` 8080 &
		netcat -l -k `+hostCIPv4+` 8443 &
		netcat -l -k `+hostCIPv4+` 33434 &
		netcat -l -k `+hostCIPv4+` 31337 &

		# blocking for ipv4 tcp tests to finish
		netcat -l `+hostCIPv4+` 22222 | grep -q "TCP4DONE" || exit 1
		kill $(pidof netcat) || exit 1

		# tcp6
		# for that we need to close down the other netcat TCP servers.
		# we open another netcat connection and wait for the FIN signal.
		# once received, close down the other netcat tcp4 servers and open up
		# the tcp6 servers.

		#for port in `+tcpTestPorts+`; do
		#	netcat -l `+hostCIPv6+` $port &
		#done

		netcat -l `+hostCIPv6+` 80 &
		netcat -l `+hostCIPv6+` 443 &
		netcat -l `+hostCIPv6+` 8080 &
		netcat -l `+hostCIPv6+` 33434 &
		netcat -l `+hostCIPv6+` 31337 &

		netcat -l `+hostCIPv4+` 22222 | grep -q "TESTDONE" && echo "ALL TESTS PASSED MARKER" || exit 1
		#sleep 300
	`)

	vmA := scriptvm.Start(t, "vmA", scriptHostA.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/core/ping",
				"github.com/u-root/u-root/cmds/core/grep",
				"github.com/u-root/u-root/cmds/core/tail",
				"github.com/u-root/u-root/cmds/core/netcat",
			),
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(3*time.Minute),
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
				"github.com/u-root/u-root/cmds/core/netcat",
				"github.com/u-root/u-root/cmds/core/grep",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(3*time.Minute),
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
				"github.com/u-root/u-root/cmds/core/netcat",
				"github.com/u-root/u-root/cmds/core/grep",
				"github.com/u-root/u-root/cmds/core/tail",
				"github.com/u-root/u-root/cmds/core/kill",
				"github.com/u-root/u-root/cmds/core/pidof",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(3*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect="+qemuSocketHost+":"+qemuSocketPortBC),
		),
	)

	if _, err := vmA.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("traceroute vmA: %v", err)
	}
	if _, err := vmB.Console.ExpectEOF(); err != nil {
		t.Errorf("traceroute vmB: %v", err)
	}
	if _, err := vmC.Console.ExpectEOF(); err != nil {
		t.Errorf("traceroute vmC: %v", err)
	}

	vmA.Wait()
	vmB.Wait()
	vmC.Wait()
}
