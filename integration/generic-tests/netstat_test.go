// Copyright 2021-2025 the u-root Authors. All rights reserved
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
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func TestNetstat(t *testing.T) {
	tests := []struct {
		cmd string
		exp expect.ExpectOpt
	}{
		{
			cmd: "netstat -I lo",
			exp: expect.All(
				expect.String("Kernel Interface table"),
				expect.String("Iface            MTU      Rx-OK    Rx-ERR   Rx-DRP   Rx-OVR   TX-OK    TX-ERR   TX-DRP   TX-OVR   Flg"),
				expect.String("lo               65536    8        0        0        0        8        0        0        0        LUR"),
			),
		},
		{
			cmd: "netstat -r",
			exp: expect.All(
				expect.String("Kernel IP routing table"),
				expect.String("Destination      Gateway          Genmask          Flags    MSS Window  irrt Iface"),
				expect.String("default          0.0.0.0          0.0.0.0          U        0   0          0 eth0"),
				expect.String("192.168.0.0      0.0.0.0          255.255.255.0    U        0   0          0 eth0"),
			),
		},
		{
			cmd: "netstat -s",
			exp: expect.All(
				expect.String("ip:"),
				expect.String("Forwarding is 2"),
				expect.String("Default TTL is 64"),
				expect.String("13 total packets received"),
				expect.String("0 forwarded"),
				expect.String("0 incoming packets discarded"),
				expect.String("13 incoming packets delivered"),
				expect.String("13 requests sent out"),
				expect.String("icmp:"),
				expect.String("6 ICMP messages received"),
				expect.String("0 input ICMP message failed"),
				expect.String("6 ICMP messages sent"),
				expect.String("0 ICMP messages failed"),
				expect.String("Input historam:"),
				expect.String("destination unreachable: 4"),
				expect.String("echo requests: 1"),
				expect.String("echo replies: 1"),
				expect.String("Output historam:"),
				expect.String("IcmpMsg:"),
				expect.String("InType3: 4"),
				expect.String("OutType3: 4"),
				expect.String("tcp:"),
				expect.String("RtoAlgorithm: 1"),
				expect.String("RtoMin: 200"),
				expect.String("RtoMax: 120000"),
				expect.String("MaxConn: -1"),
				expect.String("2 active connection openings"),
				expect.String("2 passive connection openings"),
				expect.String("0 failed connection attempts"),
				expect.String("0 connection resets received"),
				expect.String("4 connections established"),
				expect.String("6 segments received"),
				expect.String("6 resets sent"),
				expect.String("0 segments retransmitted"),
				expect.String("0 bad segments received"),
				expect.String("0 segments sent out"),
				expect.String("udp:"),
				expect.String("0 packets received"),
				expect.String("4 packets to unknown port received"),
				expect.String("0 packet receive errors"),
				expect.String("4 packets sent"),
				expect.String("0 receive buffer errors"),
				expect.String("0 send buffer errors"),
				expect.String("tcpExt:"),
				expect.String("2 acknowledgments not containing data payload received"),
				expect.String("TCPDelivered: 2"),
				expect.String("ipExt:"),
				expect.String("InOctets: 1068"),
				expect.String("OutOctets: 1068"),
				expect.String("InNoECTPkts: 13"),
			),
		},
		{
			cmd: "netstat -4 --tcp --all --numeric",
			exp: expect.All(
				expect.String("Proto  Recv-Q Send-Q Local Address                       Foreign Address                     State"),
				expect.String("tcp         0      0 127.0.0.1:5005                      0.0.0.0:0                           LISTEN"),
				expect.RegexpPattern(`tcp\s+0\s+0\s+127\.0\.0\.1:5005\s+127\.0\.0\.1:\d+\s+ESTABLISHED`),
				expect.RegexpPattern(`tcp\s+0\s+0\s+127\.0\.0\.1:\d+\s+127\.0\.0\.1:5005\s+ESTABLISHED`),
			),
		},
		{
			cmd: "netstat -4 --udp --all --numeric",
			exp: expect.All(
				expect.String("Proto  Recv-Q Send-Q Local Address                       Foreign Address                     State"),
				expect.RegexpPattern(`udp\s+0\s+0\s+127\.0\.0\.1:\d+\s+127\.0\.0\.1:5005\s+ESTABLISHED`),
				expect.String("udp         0      0 127.0.0.1:5005                      0.0.0.0:0                           CLOSE"),
			),
		},
		{
			cmd: "netstat -6 --tcp --all --numeric",
			exp: expect.All(
				expect.String("Proto  Recv-Q Send-Q Local Address                       Foreign Address                     State"),
				expect.String("tcp6        0      0 ::1:5005                            :::0                                LISTEN"),
				expect.RegexpPattern(`tcp6\s+0\s+0\s+::1:5005\s+::1:\d+\s+ESTABLISHED`),
				expect.RegexpPattern(`tcp6\s+0\s+0\s+::1:\d+\s+::1:5005\s+ESTABLISHED`),
			),
		},
		{
			cmd: "netstat -6 --udp --all --numeric",
			exp: expect.All(
				expect.String("Proto  Recv-Q Send-Q Local Address                       Foreign Address                     State"),
				expect.RegexpPattern(`udp6\s+0\s+0\s+::1:\d+\s+::1:5005\s+ESTABLISHED`),
				expect.String("udp6        0      0 ::1:5005                            :::0                                CLOSE"),
			),
		},
		{
			cmd: "netstat --unix --all",
			exp: expect.All(
				expect.String("Active sockets in the UNIX domain"),
				expect.String("Proto   RefCnt  Flags   Type      State           I-Node    Path"),
				expect.RegexpPattern(`unix\s+3\s+\[\]\s+STREAM\s+CONNECTED\s+\d+\s+stream\.sock`),
				expect.RegexpPattern(`unix\s+2\s+\[\]\s+DGRAM\s+CONNECTED\s+\d+`),
				expect.RegexpPattern(`unix\s+3\s+\[\]\s+DGRAM\s+CONNECTED\s+\d+\s+datagram\.sock`),
				expect.RegexpPattern(`unix\s+2\s+\[ACC\]\s+STREAM\s+LISTENING\s+\d+\s+stream\.sock`),
				expect.RegexpPattern(`unix\s+3\s+\[\]\s+STREAM\s+CONNECTED\s+\d+`),
			),
		},
	}

	var script strings.Builder
	fmt.Fprint(&script, `
		ip addr add 192.168.0.1/24 dev eth0
		ip link set eth0 up
		ip route add 0.0.0.0/0 dev eth0
		ping -c 1 192.168.0.1

		{
			echo 127.0.0.1 localhost4
			echo ::1 localhost6
		} >>/etc/hosts

		mkfifo fifo
		sleep 3600 >fifo &

		netcat --listen -4		   127.0.0.1 5005 <fifo >/dev/null &
		netcat --listen -4	   --udp   127.0.0.1 5005 <fifo >/dev/null &
		netcat --listen -6		   ::1       5005 <fifo >/dev/null &
		netcat --listen -6	   --udp   ::1       5005 <fifo >/dev/null &
		netcat --listen --unixsock	   stream.sock    <fifo >/dev/null &
		netcat --listen --unixsock --udp   datagram.sock  <fifo >/dev/null &
		sleep 2

		netcat -4                 127.0.0.1 5005 <fifo >/dev/null &
		netcat -4         --udp   127.0.0.1 5005 <fifo >/dev/null &
		netcat -6                 ::1       5005 <fifo >/dev/null &
		netcat -6         --udp   ::1       5005 <fifo >/dev/null &
		netcat --unixsock         stream.sock    <fifo >/dev/null &
		netcat --unixsock --udp   datagram.sock  <fifo >/dev/null &
		sleep 2
	`)
	for _, test := range tests {
		fmt.Fprintln(&script, test.cmd)
	}

	vm := scriptvm.Start(t, "vm", script.String(),
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/mkfifo",
				"github.com/u-root/u-root/cmds/core/netcat",
				"github.com/u-root/u-root/cmds/core/ping",
				"github.com/u-root/u-root/cmds/core/sleep",
			),
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/core/netstat",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qnetwork.HostNetwork("192.168.0.0/24"),
		),
	)

	for _, test := range tests {
		t.Run(test.cmd, func(t *testing.T) {
			if _, err := vm.Console.Expect(test.exp); err != nil {
				t.Errorf("VM output did not match expectations: %v", err)
			}
		})
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
