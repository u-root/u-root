// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestDhclientQEMU4(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "client"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{
					"-device", "e1000,netdev=host0",
					"-netdev", "user,id=host0,net=192.168.0.0/24,dhcpstart=192.168.0.10,ipv6=off",
				},
			},
		},
		TestCmds: []string{
			"dhclient -ipv6=false -v",
			"ip a",
			"sleep 5",
			"shutdown -h",
		},
	})
	defer ccleanup()

	if err := dhcpClient.Expect("Configured eth0 with IPv4 DHCP Lease"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
	if err := dhcpClient.Expect("inet 192.168.0.10"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
}

func TestDhclientTimesOut(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	network := qemu.NewNetwork()
	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestQEMUDHCPTimesOut",
		QEMUOpts: qemu.Options{
			Timeout: 50 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
		TestCmds: []string{
			"dhclient -v -retry 2 -timeout 10",
			"echo \"DHCP timed out\"",
			"shutdown -h",
		},
	})
	defer ccleanup()

	if err := dhcpClient.Expect("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if err := dhcpClient.Expect("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if err := dhcpClient.Expect("DHCP timed out"); err != nil {
		t.Error(err)
	}
}
