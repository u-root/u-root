// Copyright 2018-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func netcatVM(t *testing.T, name, script string, net *qnetwork.InterVM, mods ...uimage.Modifier) *qemu.VM {
	fixedMods := []uimage.Modifier{
		uimage.WithBusyboxCommands(
			"github.com/u-root/u-root/cmds/core/basename",
			"github.com/u-root/u-root/cmds/core/cat",
			"github.com/u-root/u-root/cmds/core/dirname",
			"github.com/u-root/u-root/cmds/core/echo",
			"github.com/u-root/u-root/cmds/core/grep",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/kill",
			"github.com/u-root/u-root/cmds/core/mkfifo",
			"github.com/u-root/u-root/cmds/core/rm",
			"github.com/u-root/u-root/cmds/core/seq",
			"github.com/u-root/u-root/cmds/core/shasum",
			"github.com/u-root/u-root/cmds/core/sleep",
		),
		uimage.WithCoveredCommands(
			"github.com/u-root/u-root/cmds/core/netcat",
		),
	}

	return scriptvm.Start(t, name, script,
		scriptvm.WithUimage(append(fixedMods, mods...)...),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			net.NewVM(),
		),
	)
}

func TestNetcatStream(t *testing.T) {
	net := qnetwork.NewInterVM()

	serverScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.2/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::2/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		seq -w 0 99999 >input.txt

		# loopback tests disabled due to https://github.com/mvdan/sh/issues/1142
		#
		# mkfifo fifo fifo6
		#
		# # TCPv4 server: loopback
		# : >fifo &
		# netcat -l 192.168.0.2 5005 <fifo >fifo &
		#
		# # TCPv4 server: checksum
		# netcat -l 192.168.0.2 5006 <input.txt | shasum >5006.out &
		#
		# # TCPv6 server: loopback
		# : >fifo6 &
		# netcat -l fd51:3681:1eb4::2 5005 <fifo6 >fifo6 &
		#
		# # TCPv6 server: checksum
		# netcat -l fd51:3681:1eb4::2 5006 <input.txt | shasum >5006-6.out &

		# accept file from TCPv4 client
		netcat -l 192.168.0.2 5007 </dev/null | shasum >5007.out &

		# send file to TCPv4 client
		netcat -l 192.168.0.2 5008 <input.txt &

		# exchange files with TCPv4 client
		netcat -l 192.168.0.2 5009 <input.txt | shasum >5009.out &

		# accept file from TCPv6 client
		netcat -l fd51:3681:1eb4::2 5007 </dev/null | shasum >5007-6.out &

		# send file to TCPv6 client
		netcat -l fd51:3681:1eb4::2 5008 <input.txt &

		# exchange files with TCPv6 client
		netcat -l fd51:3681:1eb4::2 5009 <input.txt | shasum >5009-6.out &

		wait

		# loopback tests disabled due to https://github.com/mvdan/sh/issues/1142
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5006.out
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5006-6.out

		# verify files from TCPv4 client
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5007.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5009.out

		# verify files from TCPv6 client
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5007-6.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5009-6.out

		# run TCPv4/v6 servers in keep-open mode for about 20 seconds
		# run TCPv4/v6 servers in broker (chat) mode for about 20 seconds
		netcat -l -k     192.168.0.2       5010 </dev/null | shasum >5010.out   &
		netcat -l -k     fd51:3681:1eb4::2 5010 </dev/null | shasum >5010-6.out &
		netcat -l --chat 192.168.0.2       5011 >5011.out   &
		netcat -l --chat fd51:3681:1eb4::2 5011 >5011-6.out &
		sleep 20
		grep -l netcat /proc/*/comm |
			while read P; do
				kill $(basename $(dirname $P))
			done
		wait


		# verify output from keep-open mode servers
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5010.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5010-6.out

		# verify output from chat servers
		expected=$(
			echo 'user<1>: hello-1'
			echo 'user<2>: hello-2'
			echo 'user<3>: hello-3'
		)
		got=$(cat 5011.out)
		test "$expected" = "$got"
		got=$(cat 5011-6.out)
		test "$expected" = "$got"
	`
	clientScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.1/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::1/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		seq -w     0 49999 >input-1.txt
		seq -w 50000 99999 >input-2.txt
		cat input-1.txt input-2.txt >input.txt

		# wait a bit for the server to come up
		sleep 3

		# loopback tests disabled due to https://github.com/mvdan/sh/issues/1142
		#
		# mkfifo fifo
		#
		# # TCPv4 client: checksum
		# netcat 192.168.0.2 5005 <input.txt | shasum >5005.out
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5005.out
		#
		# # TCPv4 client: loopback
		# : >fifo &
		# netcat 192.168.0.2 5006 <fifo >fifo
		#
		# # TCPv6 client: checksum
		# netcat fd51:3681:1eb4::2 5005 <input.txt | shasum >5005-6.out
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5005-6.out
		#
		# # TCPv6 client: loopback
		# : >fifo &
		# netcat fd51:3681:1eb4::2 5006 <fifo >fifo
		#
		# # unix server: loopback
		# : >fifo &
		# netcat -l -U stream.sock <fifo >fifo &
		# sleep 1
		#
		# # unix client: checksum
		# netcat -U stream.sock <input.txt | shasum >stream.client.out
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream.client.out
		# wait
		# rm stream.sock
		#
		# # unix server: checksum
		# netcat -l -U stream.sock <input.txt | shasum >stream.server.out &
		# sleep 1
		#
		# # unix client: loopback
		# : >fifo &
		# netcat -U stream.sock <fifo >fifo
		#
		# wait
		# rm stream.sock
		# grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream.server.out

		# upload file to TCPv4 server
		netcat 192.168.0.2 5007 <input.txt

		# download file from TCPv4 server
		netcat 192.168.0.2 5008 </dev/null | shasum >5008.out

		# exchange files with TCPv4 server
		netcat 192.168.0.2 5009 <input.txt | shasum >5009.out

		# verify files from TCPv4 server
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5008.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5009.out

		# upload file to TCPv6 server
		netcat fd51:3681:1eb4::2 5007 <input.txt

		# download file from TCPv6 server
		netcat fd51:3681:1eb4::2 5008 </dev/null | shasum >5008-6.out

		# exchange files with TCPv6 server
		netcat fd51:3681:1eb4::2 5009 <input.txt | shasum >5009-6.out

		# verify files from TCPv6 server
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5008-6.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa 5009-6.out


		# wait a bit until the keep-open and chat servers start up
		sleep 3

		# upload file in two parts to each keep-open server
		netcat 192.168.0.2	 5010 <input-1.txt
		netcat 192.168.0.2	 5010 <input-2.txt
		netcat fd51:3681:1eb4::2 5010 <input-1.txt
		netcat fd51:3681:1eb4::2 5010 <input-2.txt

		# Connect with three clients to each chat server in a predefined order (at 0,
		# 2, and 4 seconds), and once they're all connected (which happens slightly
		# after the 4 second mark), make them send strings in a predefined order (at
		# 6, 8, and 10 seconds from the start). Each client lingers until the 12
		# second mark (so that everyone can hear everyone).

		(sleep 6; echo hello-1; sleep 6) | netcat 192.168.0.2       5011 >5011-1.out   &
		(sleep 6; echo hello-1; sleep 6) | netcat fd51:3681:1eb4::2 5011 >5011-6-1.out &
		sleep 2
		(sleep 6; echo hello-2; sleep 4) | netcat 192.168.0.2       5011 >5011-2.out   &
		(sleep 6; echo hello-2; sleep 4) | netcat fd51:3681:1eb4::2 5011 >5011-6-2.out &
		sleep 2
		(sleep 6; echo hello-3; sleep 2) | netcat 192.168.0.2       5011 >5011-3.out   &
		(sleep 6; echo hello-3; sleep 2) | netcat fd51:3681:1eb4::2 5011 >5011-6-3.out &
		wait

		# verify output from each chat client
		expected1=$(
			echo 'user<2>: hello-2'
			echo 'user<3>: hello-3'
		)
		expected2=$(
			echo 'user<1>: hello-1'
			echo 'user<3>: hello-3'
		)
		expected3=$(
			echo 'user<1>: hello-1'
			echo 'user<2>: hello-2'
		)

		got1=$(cat 5011-1.out)
		got2=$(cat 5011-2.out)
		got3=$(cat 5011-3.out)
		test "$expected1" = "$got1"
		test "$expected2" = "$got2"
		test "$expected3" = "$got3"

		got1=$(cat 5011-6-1.out)
		got2=$(cat 5011-6-2.out)
		got3=$(cat 5011-6-3.out)
		test "$expected1" = "$got1"
		test "$expected2" = "$got2"
		test "$expected3" = "$got3"


		# accept file from unix client
		netcat -l -U stream-1.sock </dev/null | shasum >stream-1.server.out &

		# send file to unix client
		netcat -l -U stream-2.sock <input.txt &

		# exchange files with unix client
		netcat -l -U stream-3.sock <input.txt | shasum >stream-3.server.out &

		sleep 1

		# upload file to unix server
		netcat -U stream-1.sock <input.txt

		# download file from unix server
		netcat -U stream-2.sock </dev/null | shasum >stream-2.client.out

		# exchange files with unix server
		netcat -U stream-3.sock <input.txt | shasum >stream-3.client.out

		# verify files from unix client
		wait
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream-1.server.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream-3.server.out

		# verify files from unix server
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream-2.client.out
		grep -q a7ffaef825af40e08daef5a1e0804d851904b5aa stream-3.client.out
	`

	serverVM := netcatVM(t, "netcat_server", serverScript, net)
	clientVM := netcatVM(t, "netcat_client", clientScript, net)

	if _, err := serverVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("serverVM: %v", err)
	}
	if _, err := clientVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("clientVM: %v", err)
	}

	clientVM.Wait()
	serverVM.Wait()
}

func TestNetcatDatagram(t *testing.T) {
	net := qnetwork.NewInterVM()

	serverScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.2/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::2/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		# Start four simple "echo servers", and let them run for 30 seconds. (There
		# is no EOF propagation over datagram sockets, and so these netcat servers
		# would run forever; we need to kill them manually.)
		function reply
		(
			while read CLIENT_MSG; do
				echo "back from server: $CLIENT_MSG"
			done <netcat.$1.out.fifo >netcat.$1.in.fifo
		)
		for ((K=0; K<4; K++)); do
			mkfifo netcat.$K.{in,out}.fifo
		done

		# gosh bug: when starting a background command in a loop using variable
		# substitution, gosh does not seem to perform the substitution *first*.
		# Instead, each command running in the background sees the loop variable
		# continue changing. O_o ... fixed in upstream gosh commit 87e88a4ca0ba
		# ("interp: make a full copy of the environment for background subshells",
		# 2025-03-29)
		reply 0 &
		reply 1 &
		reply 2 &
		reply 3 &

		# Listen on port 5005.
		netcat --listen --udp 192.168.0.2       5005 >netcat.0.out.fifo <netcat.0.in.fifo &
		netcat --listen --udp fd51:3681:1eb4::2 5005 >netcat.1.out.fifo <netcat.1.in.fifo &

		# Listen on port 5006.
		netcat --listen --udp 192.168.0.2       5006 >netcat.2.out.fifo <netcat.2.in.fifo &
		netcat --listen --udp fd51:3681:1eb4::2 5006 >netcat.3.out.fifo <netcat.3.in.fifo &

		sleep 30
		grep -l netcat /proc/*/comm |
			while read P; do
				kill $(basename $(dirname $P))
			done
		wait
	`
	clientScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.1/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::1/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		# wait a bit for the server to come up
		sleep 3

		# Produce three lines of output, sleeping 1 second aftrer each, for two
		# purposes: (a) each netcat client below should read each line separately,
		# and send it to the corresponding netcat server in a separate datagram, (b)
		# we should give enough time to the server for responding to each datagram.
		# (The netcat datagram client does exit upon EOF from stdin.)
		function hello_single
		{
			echo "hello-$1"
			sleep 1
		}
		function hello
		{
			hello_single 1
			hello_single 2
			hello_single 3
		}
		expected=$(
			echo 'back from server: hello-1'
			echo 'back from server: hello-2'
			echo 'back from server: hello-3'
		)

		# Trigger echoes from the first two servers using fixed source ports. Each
		# server locks on to the source address:port of the first datagram that it
		# receives. Because we use fixed source ports, we can use distinct netcat
		# client processes for sending the datagrams, and still satisfy the servers.
		(
			hello_single 1 | netcat --udp --source 192.168.0.1       --source-port 12345 192.168.0.2       5005 >>netcat.0.out
			hello_single 2 | netcat --udp --source 192.168.0.1       --source-port 12345 192.168.0.2       5005 >>netcat.0.out
			hello_single 3 | netcat --udp --source 192.168.0.1       --source-port 12345 192.168.0.2       5005 >>netcat.0.out
		) &
		(
			hello_single 1 | netcat --udp --source fd51:3681:1eb4::1 --source-port 12345 fd51:3681:1eb4::2 5005 >>netcat.1.out
			hello_single 2 | netcat --udp --source fd51:3681:1eb4::1 --source-port 12345 fd51:3681:1eb4::2 5005 >>netcat.1.out
			hello_single 3 | netcat --udp --source fd51:3681:1eb4::1 --source-port 12345 fd51:3681:1eb4::2 5005 >>netcat.1.out
		) &
		wait

		# Trigger echoes from the last two servers using OS-assigned source ports.
		hello | netcat --udp 192.168.0.2       5006 >netcat.2.out &
		hello | netcat --udp fd51:3681:1eb4::2 5006 >netcat.3.out &
		wait

		# Repeat the same tests locally, with unix domain datagram sockets.
		function reply
		(
			while read CLIENT_MSG; do
				echo "back from server: $CLIENT_MSG"
			done <netcat.$1.out.fifo >netcat.$1.in.fifo
		)

		mkfifo netcat.4.{in,out}.fifo
		mkfifo netcat.5.{in,out}.fifo
		reply 4 &
		reply 5 &
		netcat --listen --udp --unixsock dgram.4.sock >netcat.4.out.fifo <netcat.4.in.fifo &
		netcat --listen --udp --unixsock dgram.5.sock >netcat.5.out.fifo <netcat.5.in.fifo &

		sleep 1

		hello_single 1 | netcat --udp --unixsock --source source.dgram.sock dgram.4.sock >>netcat.4.out
		rm source.dgram.sock
		hello_single 2 | netcat --udp --unixsock --source source.dgram.sock dgram.4.sock >>netcat.4.out
		rm source.dgram.sock
		hello_single 3 | netcat --udp --unixsock --source source.dgram.sock dgram.4.sock >>netcat.4.out
		rm source.dgram.sock
		hello          | netcat --udp --unixsock                            dgram.5.sock >>netcat.5.out

		# Kill the local servers.
		grep -l netcat /proc/*/comm |
			while read P; do
				kill $(basename $(dirname $P))
			done
		wait
		rm dgram.4.sock dgram.5.sock

		# Verify replies.
		for ((K=0; K<6; K++)); do
			got=$(cat netcat.$K.out)
			test "$expected" = "$got"
		done
	`

	serverVM := netcatVM(t, "netcat_server", serverScript, net)
	clientVM := netcatVM(t, "netcat_client", clientScript, net)

	if _, err := serverVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("serverVM: %v", err)
	}
	if _, err := clientVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("clientVM: %v", err)
	}

	clientVM.Wait()
	serverVM.Wait()
}

func TestNetcatExec(t *testing.T) {
	net := qnetwork.NewInterVM()

	serverScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.2/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::2/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		# Launch four netcat servers; two for a single reception each, and two for
		# multiple receptions each. Kill the latter after 30 seconds.
		netcat -l    192.168.0.2       5005 </dev/null >5005.out   &
		netcat -l    fd51:3681:1eb4::2 5005 </dev/null >5005-6.out &
		netcat -l -k 192.168.0.2       5006 </dev/null >5006.out   &
		netcat -l -k fd51:3681:1eb4::2 5006 </dev/null >5006-6.out &

		sleep 30
		grep -l netcat /proc/*/comm |
			while read P; do
				kill $(basename $(dirname $P))
			done
		wait

		# Check outputs.
		expected_double=$(
			echo 'hello world'
			echo 'hello world'
		)

		test "$(<5005.out)"   = "hello world"
		test "$(<5005-6.out)" = "hello world"
		test "$(<5006.out)"   = "$expected_double"
		test "$(<5006-6.out)" = "$expected_double"
	`
	clientScript := `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this virtual
		# network, and it will only prevent netcat from binding our unique local
		# address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad

		ip    addr add 192.168.0.1/24        dev eth0
		ip -6 addr add fd51:3681:1eb4::1/126 dev eth0
		ip link set eth0 up
		ip    route add 0.0.0.0/0 dev eth0
		ip -6 route add ::/0      dev eth0
		echo "192.168.0.1       netcat_client" >>/etc/hosts
		echo "fd51:3681:1eb4::1 netcat_client" >>/etc/hosts
		echo "192.168.0.2       netcat_server" >>/etc/hosts
		echo "fd51:3681:1eb4::2 netcat_server" >>/etc/hosts

		# wait a bit for the server to come up
		sleep 3

		# Single sends.
		netcat --exec 'echo hello world' 192.168.0.2       5005
		netcat --exec 'echo hello world' fd51:3681:1eb4::2 5005

		# Repeated sends.
		netcat --exec 'echo hello world' 192.168.0.2       5006
		netcat --exec 'echo hello world' 192.168.0.2       5006
		netcat --exec 'echo hello world' fd51:3681:1eb4::2 5006
		netcat --exec 'echo hello world' fd51:3681:1eb4::2 5006

		# Repeat the tests locally.
		netcat -l    -U netcat-1.sock </dev/null >netcat-1.out &
		netcat -l -k -U netcat-2.sock </dev/null >netcat-2.out &
		sleep 1
		netcat --exec 'echo hello world' -U netcat-1.sock
		netcat --exec 'echo hello world' -U netcat-2.sock
		netcat --exec 'echo hello world' -U netcat-2.sock
		sleep 2
		grep -l netcat /proc/*/comm |
			while read P; do
				kill $(basename $(dirname $P))
			done
		wait
		rm netcat-1.sock netcat-2.sock

		expected_double=$(
			echo 'hello world'
			echo 'hello world'
		)
		test "$(<netcat-1.out)" = "hello world"
		test "$(<netcat-2.out)" = "$expected_double"
	`

	serverVM := netcatVM(t, "netcat_server", serverScript, net)
	clientVM := netcatVM(t, "netcat_client", clientScript, net)

	if _, err := serverVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("serverVM: %v", err)
	}
	if _, err := clientVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("clientVM: %v", err)
	}

	clientVM.Wait()
	serverVM.Wait()
}
