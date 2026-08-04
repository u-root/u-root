// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func TestPing(t *testing.T) {
	const serverScript = `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this
		# virtual network, and it will only prevent the client from binding our
		# unique local address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad
		ip addr add 192.168.0.2/24 dev eth0 || exit 1
		ip -6 addr add fd51:3681:1eb4::2/64 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		sleep 30
	`
	serverVM := scriptvm.Start(t, "ping_server", serverScript,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,listen=:1236"),
		),
	)
	t.Cleanup(func() {
		_ = serverVM.Kill()
		_ = serverVM.Wait()
	})

	const clientScript = `
		# Disable IPv6 Duplicate Address Discovery. We don't need it on this
		# virtual network, and it will only prevent us from binding our unique
		# local address (ULA) for several seconds.
		echo 0 >/proc/sys/net/ipv6/conf/eth0/accept_dad
		ip addr add 192.168.0.1/24 dev eth0 || exit 1
		ip -6 addr add fd51:3681:1eb4::1/64 dev eth0 || exit 1
		ip link set eth0 up || exit 1
		ip link set lo up || exit 1
		sleep 2
		ping -c 1 -w 5000 192.168.0.2 || exit 1
		ping -6 -c 1 -w 5000 fd51:3681:1eb4::2 || exit 1

		# Regression test for #2974: pinging loopback with -c > 1 used to
		# fail on the second ping with "got seq 1; want 2".
		ping -c 5 -i 200 -w 5000 127.0.0.1 || exit 1
		ping -6 -c 5 -i 200 -w 5000 ::1 || exit 1
		echo "TESTS PASSED MARKER"
	`
	clientVM := scriptvm.Start(t, "ping_client", clientScript,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
			),
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/ping"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-nic", "socket,connect=127.0.0.1:1236"),
		),
	)
	if _, err := clientVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("clientVM: %v", err)
	}
	if err := clientVM.Wait(); err != nil {
		t.Errorf("clientVM wait: %v", err)
	}
}
