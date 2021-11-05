// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lacp

import (
	"net"
	"reflect"
	"runtime"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/vmtest"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type tearDown func()

func setupNetlinkTest(t *testing.T) tearDown {
	// Creating a temp namespace and editing interfaces requires root
	testutil.SkipIfNotRoot(t)

	// Skip test architectures where netlink ops fail with unsupported.
	if vmtest.TestArch() == "amd64" || vmtest.TestArch() == "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	runtime.LockOSThread()
	// Save the current network namespace
	origns, _ := netns.Get()

	// new temporary namespace so we don't pollute the host
	// lock thread since the namespace is thread local
	ns, err := netns.New()
	if err != nil {
		t.Fatal("Failed to create newns", err)
	}
	err = netns.Set(ns)
	if err != nil {
		t.Fatal("Failed to set newns", err)
	}

	return func() {
		netns.Set(origns)
		ns.Close()
		origns.Close()
		runtime.UnlockOSThread()
	}
}

func mac(t *testing.T, s string) net.HardwareAddr {
	mac, err := net.ParseMAC(s)
	if err != nil {
		t.Fatal("Test setup failed with bad mac address", s)
	}
	return mac
}

func setupDefaultLinks(t *testing.T) []netlink.Link {
	links := []netlink.Link{
		&netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name:         "uroot0",
				HardwareAddr: mac(t, "AA:AA:AA:AA:AA:AA"),
			},
		},
		&netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name:         "uroot1",
				HardwareAddr: mac(t, "AA:AA:AA:AA:AA:BB"),
			},
		},
		&netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name:         "uroot2",
				HardwareAddr: mac(t, "AA:AA:AA:AA:AA:CC"),
			},
		},
	}

	for _, link := range links {
		err := netlink.LinkAdd(link)
		if err != nil {
			t.Fatalf("Failed adding link %v during setup: %v", link, err)
		}
	}
	return links
}

func checkLinkMac(t *testing.T, linkName string, wantMac net.HardwareAddr) {
	if got, err := netlink.LinkByName(linkName); err != nil {
		t.Errorf("Failed to find link %s in interface list: %v", linkName, err)
	} else {
		if !reflect.DeepEqual(got.Attrs().HardwareAddr, wantMac) {
			t.Errorf("Bad MAC on interface %s. Got: %s, Want: %s", linkName, got.Attrs().HardwareAddr, wantMac)
		}
	}
}

func checkLinkFlags(t *testing.T, linkName string, want net.Flags) {
	if link, err := netlink.LinkByName(linkName); err != nil {
		t.Errorf("Failed getting interface for %s: %v", linkName, err)
	} else if link.Attrs().Flags&want != want {
		t.Errorf("Failed interface flags (got & want != want): Want %v; got %v", want, link.Attrs().Flags)
	}
}

func checkLinks(t *testing.T, wants []netlink.Link) {
	for _, want := range wants {
		checkLinkMac(t, want.Attrs().Name, want.Attrs().HardwareAddr)
	}
}

func checkInterfaceCount(t *testing.T, want int) {
	if ifs, err := net.Interfaces(); err != nil {
		t.Errorf("Failed getting interface list for count: %v", err)
	} else if len(ifs) != want {
		t.Errorf("Failed interface count: Want (links+bond) %d; got %d: %v", want, len(ifs), ifs)
	}
}

func TestRemoveExistingBonds(t *testing.T) {
	tearDown := setupNetlinkTest(t)
	defer tearDown()
	links := setupDefaultLinks(t)

	bondName := "urBond"
	bond := netlink.NewLinkBond(netlink.LinkAttrs{
		Name:         bondName,
		HardwareAddr: mac(t, "AA:BB:CC:DD:EE:FF"),
	})
	err := netlink.LinkAdd(bond)
	if err != nil {
		t.Fatalf("Failed adding bond %v: %v", bond, err)
	}

	err = RemoveExistingBonds()
	if err != nil {
		t.Fatalf("Failed RemoveExistingBonds. Got %v, want: nil", err)
	}

	got, err := netlink.LinkByName(bondName)
	if err == nil {
		t.Errorf("Found bond with name %s: Got: %v; Want nil", got.Attrs().Name, got)
	}

	// Verify all other interfaces still present
	checkLinks(t, links)
	checkInterfaceCount(t, len(links)+1) // Links + lo
}

func TestCreateLACPBond(t *testing.T) {
	tearDown := setupNetlinkTest(t)
	defer tearDown()
	links := setupDefaultLinks(t)

	bondName := "urBond"
	// Expected to use array[0] for MAC
	wantBondMac := links[0].Attrs().HardwareAddr
	_, err := CreateLACPBond([]netlink.Link{links[0], links[1]}, bondName)
	if err != nil {
		t.Fatalf("Failed CreateLACPBond. Got %v, want: nil", err)
	}

	// Verify all other interfaces still present. Skip links[1] who should report
	// the bond MAC.
	checkLinks(t, append(links[:0], links[2:]...))
	checkLinkMac(t, links[1].Attrs().Name, wantBondMac)
	checkLinkMac(t, bondName, wantBondMac)
	checkLinkFlags(t, bondName, net.FlagUp)
	checkInterfaceCount(t, len(links)+2) // Links + lo + new bond
}
