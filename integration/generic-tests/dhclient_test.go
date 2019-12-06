// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

func RunTestDhclient(t *testing.T, initramfs string) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	network := qemu.NewNetwork()
	_, scleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name:      "TestDhclient_Server",
		Initramfs: initramfs,
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "server"),
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
		TestCmds: []string{
			"ip link set eth0 up",
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip route add 0.0.0.0/0 dev eth0",
			"pxeserver",
		},
	})
	defer scleanup()

	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name:      "TestDhclient_Client",
		Initramfs: initramfs,
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "client"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
		TestCmds: []string{
			"dhclient -ipv6=false -v",
			"ip a",
			// Sleep so serial console output gets flushed. The expect library is racy.
			"sleep 5",
			"shutdown -h",
		},
	})
	defer ccleanup()

	if err := dhcpClient.Expect("Configured eth0 with IPv4 DHCP Lease"); err != nil {
		t.Error(err)
	}
	if err := dhcpClient.Expect("inet 192.168.0.2"); err != nil {
		t.Error(err)
	}
}

// TestPxeboot runs a server and client to test pxebooting a node.
// TODO: FIX THIS TEST!
// Change the t.Logf below back to t.Errorf
func RunTestPxeboot(t *testing.T, initramfs string) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	var serverInitramfs string
	var clientInitramfs string
	if len(initramfs) == 0 {
		serverOpts := uroot.Opts{
			ExtraFiles: []string{"../testdata/pxe:pxeroot"},
		}
		serverf, err := vmtest.CreateTestInitramfs(
			serverOpts, "github.com/u-root/u-root/integration/testcmd/generic/uinit", "")
		if err != nil {
			t.Errorf("failed to create server test initramfs: %v", err)
		}
		defer os.Remove(serverf)
		serverInitramfs = serverf

		clientOpts := uroot.Opts{
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/elvish",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/shutdown",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/boot/pxeboot",
			),
		}
		clientf, err := vmtest.CreateTestInitramfs(
			clientOpts, "github.com/u-root/u-root/integration/testcmd/generic/uinit", "")
		if err != nil {
			t.Errorf("failed to create client test initramfs: %v", err)
		}
		defer os.Remove(clientf)
		clientInitramfs = clientf
	} else {
		serverInitramfs = initramfs
		clientInitramfs = initramfs
	}

	network := qemu.NewNetwork()
	dhcpServer, scleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name:      "TestPxeboot_Server",
		Initramfs: serverInitramfs,
		TestCmds: []string{
			"ip addr add 192.168.0.1/24 dev eth0",
			"ip link set eth0 up",
			"ip route add 0.0.0.0/0 dev eth0",
			"ls -l /pxeroot",
			"pxeserver -tftp-dir=/pxeroot",
		},
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "server"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
	})
	defer scleanup()

	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name:      "TestPxeboot_Client",
		Initramfs: clientInitramfs,
		TestCmds: []string{
			"pxeboot --dry-run --no-load -v",
			// Sleep so serial console output gets flushed. The expect library is racy.
			"sleep 5",
			"shutdown -h",
		},
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "client"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
	})
	defer ccleanup()

	if err := dhcpServer.Expect("starting file server"); err != nil {
		t.Logf("File server: %v", err)
	}
	if err := dhcpClient.Expect("Got DHCPv4 lease on eth0:"); err != nil {
		t.Logf("Lease %v:", err)
	}
	if err := dhcpClient.Expect("Boot URI: tftp://192.168.0.1/pxelinux.0"); err != nil {
		t.Logf("Boot: %v", err)
	}
}

func RunTestQEMUDHCPTimesOut(t *testing.T, initramfs string) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name:      "TestQEMUDHCPTimesOut",
		Initramfs: initramfs,
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "client"),
			Timeout:      40 * time.Second,
		},
		TestCmds: []string{
			// loopback should time out and it can't have configured anything.
			"dhclient -v -retry 1 -timeout 10 lo",
			"echo \"DHCP timed out\"",
			// Sleep so serial console output gets flushed. The expect library is racy.
			"sleep 5",
			"shutdown -h",
		},
	})
	defer ccleanup()

	// Make sure that dhclient does not hang forever.
	if err := dhcpClient.Expect("DHCP timed out"); err != nil {
		t.Error(err)
	}
}
