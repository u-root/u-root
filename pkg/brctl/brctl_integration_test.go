// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package brctl

// Sometimes manual testing might be necessary or just more straight forward.
// To setup a local test environment similar to the integration test, run the following commands.
// Since the tests issue raw ioctl calls, they have to be run as root.
//
// ```
// ip link add eth10 type dummy
// ip link add eth10 type dummy
// brctl addbr br0
// brctl addbr br1
// brctl addif br0 eth0
// brctl addif br1 eth1
// ````

import (
	"fmt"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

var (
	BRCTL_TEST_IFACE_0 = "eth0"
	BRCTL_TEST_IFACE_1 = "eth1"
	BRCTL_TEST_IFACES  = []string{BRCTL_TEST_IFACE_0, BRCTL_TEST_IFACE_1}

	BRCTL_TEST_BR_0    = "br0"
	BRCTL_TEST_BR_1    = "br1"
	BRCTL_TEST_BRIDGES = []string{BRCTL_TEST_BR_0, BRCTL_TEST_BR_1}
)

// TODO: Since ioctl needs root privileges, we need to run the tests in a VM with root privileges.
func TestVM(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	// https://elixir.bootlin.com/linux/v6.0/source/net/socket.c#L1136
	// kernel needs to have the bridge built in
	// CONFIG_BRIDGE=y
	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/brctl"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", BRCTL_TEST_IFACE_0)),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", BRCTL_TEST_IFACE_1)),
		),
	)
}
