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
	network := qemu.NewNetwork()
	_, scleanup := QEMUTest(t, &Options{
		Name: "TestDhclient_Server",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/integration/testcmd/pxeserver",
		},
		Uinit: []string{
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip link set eth0 up",
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeserver",
			"shutdown -h",
		},
		Network: network,
	})
	defer scleanup()

	dhcpClient, ccleanup := QEMUTest(t, &Options{
		Name: "TestDhclient_Client",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/dhclient",
			"github.com/u-root/u-root/cmds/shutdown",
		},
		Uinit: []string{
			"dhclient -ipv6=false -verbose",
			"ip a",
			"shutdown -h",
		},
		Network: network,
		Timeout: 30 * time.Second,
	})
	defer ccleanup()

	if err := dhcpClient.Expect("err from done <nil>"); err != nil {
		t.Fatal(err)
	}

	if err := dhcpClient.Expect("inet 192.168.1.0"); err != nil {
		t.Fatal(err)
	}
}

func TestPxeboot(t *testing.T) {
	network := qemu.NewNetwork()
	dhcpServer, scleanup := QEMUTest(t, &Options{
		Name: "TestPxeboot_Server",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/ip",
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
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/cmds/pxeboot",
		},
		Uinit: []string{
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeboot --dry-run",
			"shutdown -h",
		},
		Network: network,
	})
	defer ccleanup()

	dhcpClient.Expect("")
	dhcpServer.Expect("")
}
