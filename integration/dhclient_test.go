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
	var sb, cb wc
	_, scleanup := QEMUTest(t, &Options{
		Name:         "TestDhclient_Server",
		SerialOutput: &sb,
		Cmds: []string{
			"github.com/u-root/u-root/cmds/echo",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/sleep",
			"github.com/u-root/u-root/cmds/shutdown",
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
			"github.com/u-root/u-root/cmds/echo",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/dhclient",
			"github.com/u-root/u-root/cmds/shutdown",
		},
		Uinit: []string{
			"echo DO THAT DHCLIENT",
			"dhclient -ipv6=false -verbose",
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

	if err := dhcpClient.Expect("err from done <nil>"); err != nil {
		t.Error(err)
	}

	if err := dhcpClient.Expect("inet 192.168.1.0"); err != nil {
		t.Error(err)
	}
	t.Logf("Server: %s\nClient: %s", sb.String(), cb.String())
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
