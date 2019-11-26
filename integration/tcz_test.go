// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestTczclient(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	t.Skip("This test is flaky, and must be fixed")

	network := qemu.NewNetwork()
	// TODO: On the next iteration, this will serve and provide a missing tcz.
	var sb wc
	if true {
		q, scleanup := vmtest.QEMUTest(t, &vmtest.Options{
			Name: "TestTczclient_Server",
			BuildOpts: uroot.Opts{
				Commands: uroot.BusyBoxCmds(
					"github.com/u-root/u-root/cmds/core/dmesg",
					"github.com/u-root/u-root/cmds/core/echo",
					"github.com/u-root/u-root/cmds/core/ip",
					"github.com/u-root/u-root/cmds/core/init",
					"github.com/u-root/u-root/cmds/core/shutdown",
					"github.com/u-root/u-root/cmds/core/sleep",
					"github.com/u-root/u-root/cmds/exp/srvfiles",
				),
				ExtraFiles: []string{
					"./testdata/tczserver:tcz",
				},
			},
			TestCmds: []string{
				"dmesg",
				"ip l",
				"echo NOW DO IT",
				"ip addr add 192.168.0.1/24 dev eth0",
				"ip link set eth0 up",
				"ip route add 255.255.255.255/32 dev eth0",
				"ip l",
				"ip a",
				"echo NOW SERVER IT",
				"srvfiles -h 192.168.0.1 -d /",
				"echo The Server Completes",
				"shutdown -h",
			},
			QEMUOpts: qemu.Options{
				SerialOutput: &sb,
				Devices: []qemu.Device{
					network.NewVM(),
				},
			},
		})
		if err := q.Expect("shutdown"); err != nil {
			t.Logf("got %v", err)
		}
		defer scleanup()

		t.Logf("Server SerialOutput: %s", sb.String())
	}

	var b wc
	tczClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestTczclient_Client",
		BuildOpts: uroot.Opts{
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/echo",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/exp/tcz",
				"github.com/u-root/u-root/cmds/core/shutdown",
				"github.com/u-root/u-root/cmds/core/ls",
			),
			ExtraFiles: []string{
				"./testdata/tczclient:tcz",
			},
		},
		TestCmds: []string{
			"ip addr add 192.168.0.2/24 dev eth0",
			"ip link set eth0 up",
			//"ip route add 255.255.255.255/32 dev eth0",
			"ip a",
			"ls -l /",
			"ls -l /dev",
			"cat /proc/devices",
			"cat /proc/filesystems",
			"ip l",
			"echo let us do this now",
			"tcz -d -h 192.168.0.1 -p 8080 libXcomposite libXdamage libXinerama libxshmfence",
			"tcz -d -h 192.168.0.1 -p 8080 libXdmcp",
			"ls -l /proc/mounts",
			"cat /proc/mounts",
			"echo HI THERE",
			"ls /TinyCorePackages/tcloop",
			"shutdown -h",
		},
		QEMUOpts: qemu.Options{
			SerialOutput: &b,
		},
	})
	defer ccleanup()

	// The directory list is the last thing we get. At that point,
	// b will have the output we care about and the VM will have shut
	// down. We can do the rest of the RE matching on b.String()
	// This is a bit of a hack but it frees us from worrying
	// about the order in which things appear.
	tczs := []string{"libXcomposite", "libXdamage", "libXinerama", "libxshmfence"}
	for _, s := range tczs {
		if err := tczClient.Expect(s); err != nil {
			t.Logf("Client SerialOutput: %s", b.String())
			t.Errorf("got %v, want nil", err)
		}
		t.Logf("Matched %s", s)
	}

	if false {
		for _, s := range tczs {
			re, err := regexp.Compile(".*loop.*" + s)
			if err != nil {
				t.Errorf("Check loop device re %s: got %v, want nil", s, err)
				continue
			}
			if ok := re.MatchString(b.String()); !ok {
				t.Errorf("Check loop device %s: got no match, want match", s)
				continue
			}
		}
	}

}
