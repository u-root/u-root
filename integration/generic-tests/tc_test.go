// Copyright 2021-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func TCVM(t *testing.T, name, script string, net *qnetwork.InterVM) *qemu.VM {
	var classWant string

	if os.Getenv("VMTEST_ARCH") == "arm" {
		classWant = "testdata/tc/class.want.arm:class.want"
	} else {
		classWant = "testdata/tc/class.want:class.want"
	}

	return scriptvm.Start(t, name, script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cmp",
				"github.com/u-root/u-root/cmds/core/tee",
			),
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/exp/tc",
			),
			uimage.WithFiles(
				"testdata/tc/qdisc.want:qdisc.want",
				classWant,
				"testdata/tc/filter.want:filter.want",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			net.NewVM(),
		),
	)
}

func TestTC(t *testing.T) {
	net := qnetwork.NewInterVM()

	// Based on <https://lartc.org/howto/lartc.qdisc.classful.html#AEN1079> and
	// <https://wiki.archlinux.org/title/Advanced_traffic_control>.
	script := `
		tc qdisc add dev eth0 root handle 1: htb default 30

		tc class add dev eth0 parent 1:  classid 1:1  htb rate 6mbit            burst 15k

		tc class add dev eth0 parent 1:1 classid 1:10 htb rate 5mbit            burst 15k
		tc class add dev eth0 parent 1:1 classid 1:20 htb rate 3mbit ceil 6mbit burst 15k
		tc class add dev eth0 parent 1:1 classid 1:30 htb rate 1kbit ceil 6mbit burst 15k

		tc qdisc add dev eth0 parent 1:10 handle 10: qfq
		tc qdisc add dev eth0 parent 1:20 handle 20: qfq
		tc qdisc add dev eth0 parent 1:30 handle 30: qfq

		# - Internet Header Length: 5*4 = 20 octets
		# - protocol: TCP (6)
		# - TCP destination port 80
		tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
			match u32 0x05000000 0x0f000000 at  0 \
			match u32 0x00060000 0x00ff0000 at  8 \
			match u32 0x00000050 0x0000ffff at 20 \
			flowid 1:10

		# - Internet Header Length: 5*4 = 20 octets
		# - protocol: TCP (6)
		# - TCP source port 25
		tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
			match u32 0x05000000 0x0f000000 at  0 \
			match u32 0x00060000 0x00ff0000 at  8 \
			match u32 0x00190000 0xffff0000 at 20 \
			flowid 1:20

		for obj in qdisc class filter; do
			tc $obj show dev eth0 | tee $obj.got
			cmp $obj.want $obj.got
		done
	`
	vm := TCVM(t, "tc_test", script, net)

	if _, err := vm.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("tc_test: %v", err)
	}

	vm.Wait()
}
