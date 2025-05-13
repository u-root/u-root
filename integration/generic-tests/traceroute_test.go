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

	"github.com/Netflix/go-expect"
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
		hostAIP       = "192.168.10.10"
		hostBGatewayA = "192.168.10.1"
		hostBGatewayC = "192.168.20.1"
		hostCIP       = "192.168.20.10"
		hosts         = hostAIP + "vmA\n" + hostCIP + "vmC\n"

		qemuSocketPortAB = "1234"
		qemuSocketPortBC = "1235"
		qemuSocketHost   = "127.0.0.1"
	)

	tests := []struct {
		cmd string
		exp expect.ExpectOpt
	}{
		{
			//traceroute to vmB (192.168.20.10), 20 hops max, 60 byte packets
			//TTL: 1    192.168.10.1         (15.990 ms) 192.168.10.1         (24.019 ms)
			//TTL: 2    192.168.20.10        (24.049 ms)
			cmd: "traceroute " + hostCIP,
			exp: expect.All(
				expect.String("traceroute to vmB ("+hostCIP+"), 20 hops max, 60 byte packets "),
				expect.RegexpPattern(`TTL: 1\s+`+hostBGatewayA+`.*\(\d+\.\d+ ms\)\s+`+hostBGatewayA+`.*\(\d+\.\d+ ms\)`),
				expect.RegexpPattern(`TTL: 2\s+`+hostCIP+`.*\(\d+\.\d+ ms\)`),
			),
		},
	}

	var (
		scriptHostA strings.Builder
		scriptHostB strings.Builder
		scriptHostC strings.Builder
	)

	// fmt.Fprint(&scriptHostA, `
	// 	ip addr add 192.168.10.10/24 dev eth0
	// 	ip link set eth0 up
	// 	ip route add default via 192.168.10.1 eht0
	// `)

	fmt.Fprint(&scriptHostA, `
		sleep 20
		ip addr add `+hostAIP+`/24 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip route add default via `+hostBGatewayA+` eth0 || exit 1
		echo "`+hosts+`" > /etc/hosts || exit 1

		traceroute `+hostCIP+` || exit 1
		`)

	// fmt.Fprint(&scriptHostB, `
	// 	ip addr add 192.168.10.1/24 dev eth0
	// 	ip addr add 192.168.20.1/24 dev eth1
	// 	ip link set eth0 up
	// 	ip link set eth1 up
	// 	echo "1" > /proc/sys/net/ipv4/ip_forward # make this our router
	// `)

	fmt.Fprint(&scriptHostB, `
		ip addr add `+hostBGatewayA+`/24 dev eth0 || exit 1
		ip addr add `+hostBGatewayC+`/24 dev eth1 || exit 1
		ip link set eth0 up || exit 1
		ip link set eth1 up || exit 1
		echo "`+hosts+`" > /etc/hosts || exit 1
		echo "1" > /proc/sys/net/ipv4/ip_forward || exit 1

		ip a

		sleep 200
	`)

	// fmt.Fprint(&scriptHostC, `
	// 	ip addr add 192.168.20.10/24 dev eth0
	// 	ip link set eth0 up
	// 	ip route add default via 192.168.20.1 eth0
	// `)

	fmt.Fprint(&scriptHostC, `
		ip addr add `+hostCIP+`/24 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip route add default via `+hostBGatewayC+` eth0 || exit 1
		echo "`+hosts+`" > /etc/hosts || exit 1

		# idle and wait for traceroute to finish
		sleep 200
	`)

	for _, test := range tests {
		fmt.Fprintln(&scriptHostA, test.cmd)
	}

	vmA := scriptvm.Start(t, "vmA", scriptHostA.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/exp/traceroute",
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

	vmC := scriptvm.Start(t, "vmC", scriptHostC.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect="+qemuSocketHost+":"+qemuSocketPortBC),
		),
	)

	vmB := scriptvm.Start(t, "vmB", scriptHostB.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/exp/traceroute",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect="+qemuSocketHost+":"+qemuSocketPortAB),
			qemu.ArbitraryArgs("-nic", "socket,listen=:"+qemuSocketPortBC),
		),
	)

	// Run these tests on vmA
	for _, test := range tests {
		t.Run(test.cmd, func(t *testing.T) {
			if _, err := vmA.Console.Expect(test.exp); err != nil {
				t.Errorf("VM output did not match expectations: %v", err)
			}
		})
	}

	go vmB.Wait()
	go vmC.Wait()

	if err := vmA.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
