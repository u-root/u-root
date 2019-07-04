// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
)

func TestDhclient(t *testing.T) {
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	network := qemu.NewNetwork()
	var sb, cb wc
	_, scleanup := QEMUTest(t, &Options{
		Name:         "TestDhclient_Server",
		SerialOutput: &sb,
		Cmds: []string{
			"github.com/u-root/u-root/cmds/core/echo",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/sleep",
			"github.com/u-root/u-root/cmds/core/shutdown",
			"github.com/u-root/u-root/integration/testcmd/pxeserver",
		},
		Uinit: []string{
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip link set eth0 up",
			"ip route add 255.255.255.255/32 dev eth0",
			"echo RUN THAT PXE SERVER",
			"pxeserver",
			"echo ALL DONE",
			"sleep 15",
			"shutdown -h",
		},
		Network: network,
	})
	defer scleanup()

	dhcpClient, ccleanup := QEMUTest(t, &Options{
		Name:         "TestDhclient_Client",
		SerialOutput: &cb,
		Cmds: []string{
			"github.com/u-root/u-root/cmds/core/echo",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/dhclient",
			"github.com/u-root/u-root/cmds/core/shutdown",
		},
		Uinit: []string{
			"echo DO THAT DHCLIENT",
			"dhclient -ipv6=false -v",
			"echo BACK, WHAT IP",
			"ip a",
			"echo OK, ALL DONE",
			"shutdown -h",
		},
		Network: network,
		Timeout: 30 * time.Second,
	})
	defer ccleanup()

	t.Logf("Now we wait!")

	if err := dhcpClient.Expect("Configured eth0 with IPv4 DHCP Lease"); err != nil {
		t.Error(err)
	}

	if err := dhcpClient.Expect("inet 192.168.1.0"); err != nil {
		t.Error(err)
	}
	t.Logf("Server: %s\nClient: %s", sb.String(), cb.String())
}

func TestPxeboot(t *testing.T) {
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	network := qemu.NewNetwork()
	dhcpServer, scleanup := QEMUTest(t, &Options{
		Name: "TestPxeboot_Server",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/integration/testcmd/pxeserver",
		},
		Uinit: []string{
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip link set eth0 up",
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeserver -dir=/pxeroot",
		},
		Files: []string{
			"./testdata/pxe:pxeroot",
		},
		Network: network,
	})
	defer scleanup()

	dhcpClient, ccleanup := QEMUTest(t, &Options{
		Name: "TestPxeboot_Client",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/shutdown",
			"github.com/u-root/u-root/cmds/boot/pxeboot",
		},
		Uinit: []string{
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeboot --dry-run",
			"echo PXE SUCCESSFUL",
			"shutdown -h",
		},
		Network: network,
	})
	defer ccleanup()

	if err := dhcpClient.Expect("PXE SUCCESSFUL"); err != nil {
		t.Errorf("Expected PXE SUCCESSFUL: %v", err)
	}
	dhcpServer.Expect("")
}
