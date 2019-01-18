// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
)

func TestTczclient(t *testing.T) {
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	network := qemu.NewNetwork()
	// TODO: On the next iteration, this will serve and provide a missing tcz.
	var sb wc
	if true {
		q, scleanup := QEMUTest(t, &Options{
			Name:         "TestTczclient_Server",
			SerialOutput: &sb,
			Cmds: []string{
				"github.com/u-root/u-root/cmds/dmesg",
				"github.com/u-root/u-root/cmds/echo",
				"github.com/u-root/u-root/cmds/ip",
				"github.com/u-root/u-root/cmds/init",
				"github.com/u-root/u-root/cmds/shutdown",
				"github.com/u-root/u-root/cmds/sleep",
				"github.com/u-root/u-root/cmds/srvfiles",
			},
			Uinit: []string{
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
			Files: []string{
				"./testdata/tczserver:tcz",
			},
			Network: network,
		})
		if err := q.Expect("shutdown"); err != nil {
			t.Logf("got %v", err)
		}
		defer scleanup()

		t.Logf("Server SerialOutput: %s", sb.String())
	}

	var b wc
	tczClient, ccleanup := QEMUTest(t, &Options{
		Name:         "TestTczclient_Client",
		SerialOutput: &b,
		Cmds: []string{
			"github.com/u-root/u-root/cmds/cat",
			"github.com/u-root/u-root/cmds/echo",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/tcz",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/cmds/ls",
		},
		Uinit: []string{
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
		Files: []string{
			"./testdata/tczclient:tcz",
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
