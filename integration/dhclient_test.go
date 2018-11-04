// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
)

func TestDhclient(t *testing.T) {
	QEMUTestSetup(t)

	network := qemu.NewNetwork()
	dhcpServer, err := QEMU(&Options{
		Cmds: []string{"github.com/u-root/u-root/integration/testcmd/pxeserver"},
		Uinit: []string{
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip link set eth0 up",
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeserver",
		},
		Logger: &TestLogger{t},
	})
	if err != nil {
		t.Fatal(err)
	}
	network.NewVM(dhcpServer)

	dhcpClient, err := QEMU(&Options{
		Uinit: []string{
			"dhclient -ipv6=false -verbose",
			"ip a",
		},
		Logger: &TestLogger{t},
	})
	if err != nil {
		t.Fatal(err)
	}
	network.NewVM(dhcpClient)

	t.Logf("server cmdline:\n%s", dhcpServer.CmdlineQuoted())
	t.Logf("client cmdline:\n%s", dhcpClient.CmdlineQuoted())
	if err := dhcpServer.Start(); err != nil {
		t.Fatal(err)
	}
	defer dhcpServer.Close()

	if err := dhcpClient.Start(); err != nil {
		t.Fatal(err)
	}
	defer dhcpClient.Close()

	if err := dhcpClient.Expect("err from done <nil>"); err != nil {
		//t.Logf("Client out: %v", dhcpClient.Output())
		t.Fatal(err)
	}

	if err := dhcpClient.Expect("inet 192.168.1.0"); err != nil {
		//t.Logf("Client out: %v", dhcpClient.Output())
		t.Fatal(err)
	}
}

func TestPxeboot(t *testing.T) {
	QEMUTestSetup(t)

	network := qemu.NewNetwork()
	dhcpServer, err := QEMU(&Options{
		Cmds: []string{
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
	})
	if err != nil {
		t.Fatal(err)
	}
	network.NewVM(dhcpServer)

	dhcpClient, err := QEMU(&Options{
		Uinit: []string{
			"ip route add 255.255.255.255/32 dev eth0",
			"pxeboot --dry-run",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	network.NewVM(dhcpClient)

	t.Logf("server cmdline:\n%s", dhcpServer.CmdlineQuoted())
	t.Logf("client cmdline:\n%s", dhcpClient.CmdlineQuoted())

	if err := dhcpServer.Start(); err != nil {
		t.Fatal(err)
	}
	defer dhcpServer.Close()

	if err := dhcpClient.Start(); err != nil {
		t.Fatal(err)
	}
	defer dhcpClient.Close()

	dhcpClient.Expect("")
	dhcpServer.Expect("")

	//t.Logf("Client out: %v", dhcpClient.Output())
	//t.Logf("Server out: %v", dhcpServer.Output())
}
