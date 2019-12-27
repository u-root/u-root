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
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestDhclientQEMU4 uses QEMU's DHCP server to test dhclient.
func TestDhclientQEMU4(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
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
	if vmtest.TestArch() != "amd64" {
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
