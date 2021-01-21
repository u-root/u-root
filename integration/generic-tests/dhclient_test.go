// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestDhclientQEMU4 uses QEMU's DHCP server to test dhclient.
func TestDhclientQEMU4(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	// Create the file to download
	dir, err := ioutil.TempDir("", "dhclient-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	want := "conteeent"
	foobarFile := filepath.Join(dir, "foobar")
	if err := ioutil.WriteFile(foobarFile, []byte(want), 0644); err != nil {
		t.Fatal(err)
	}

	// Serve HTTP on the host on a random port.
	http.Handle("/", http.FileServer(http.Dir(dir)))
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	s := &http.Server{}
	wg.Add(1)
	go func() {
		_ = s.Serve(ln)
		wg.Done()
	}()
	defer wg.Wait()
	defer s.Close()

	port := ln.Addr().(*net.TCPAddr).Port

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
			// Download a file to make sure dhclient configures kernel networking correctly.
			fmt.Sprintf("wget http://192.168.0.2:%d/foobar", port),
			"cat ./foobar",
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
	// "cat ./foobar" should be outputting this.
	if err := dhcpClient.Expect(want); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
}

func TestDhclientTimesOut(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	network := qemu.NewNetwork()
	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestQEMUDHCPTimesOut",
		QEMUOpts: qemu.Options{
			Timeout: 50 * time.Second,
			Devices: []qemu.Device{
				// An empty new network is easier than
				// configuring QEMU not to expose any
				// networking. At the moment.
				network.NewVM(),
			},
		},
		TestCmds: []string{
			"dhclient -v -retry 2 -timeout 10",
			"echo \"DHCP timed out\"",
			"sleep 5",
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

func TestDhclient6(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	// QEMU doesn't support DHCPv6 for getting IP configuration, so we have
	// to supply our own server.
	//
	// We don't currently have a radvd server we can use, so we also cannot
	// try to download a file using the DHCP configuration.
	network := qemu.NewNetwork()
	dhcpServer, scleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestDhclient6_Server",
		TestCmds: []string{
			"ip link set eth0 up",
			"pxeserver -6 -your-ip6=fec0::3 -4=false",
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
		Name: "TestDhclient6_Client",
		TestCmds: []string{
			"dhclient -ipv4=false -vv",
			"ip a",
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

	if err := dhcpServer.Expect("starting dhcpv6 server"); err != nil {
		t.Errorf("%s dhcpv6 server: %v", testutil.NowLog(), err)
	}
	if err := dhcpClient.Expect("Configured eth0 with IPv6 DHCP Lease IP fec0::3"); err != nil {
		t.Errorf("%s configure: %v", testutil.NowLog(), err)
	}
	if err := dhcpClient.Expect("inet6 fec0::3"); err != nil {
		t.Errorf("%s ip: %v", testutil.NowLog(), err)
	}
}
