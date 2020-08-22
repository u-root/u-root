// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package block

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/pci"
	"github.com/u-root/u-root/pkg/testutil"
)

// TestFindMountPointNotExists checks that non existent
// entry is checked and nil returned
func TestFindMountPointNotExists(t *testing.T) {
	LinuxMountsPath = "testdata/mounts"
	_, err := GetMountpointByDevice("/dev/mapper/sys-oldxxxxxx")
	require.Error(t, err)
}

// TestFindMountPointValid check for valid output of
// test mountpoint.
func TestFindMountPointValid(t *testing.T) {
	LinuxMountsPath = "testdata/mounts"
	mountpoint, err := GetMountpointByDevice("/dev/mapper/sys-old")
	require.NoError(t, err)
	require.Equal(t, *mountpoint, "/media/usb")
}

func TestParsePCIBlockList(t *testing.T) {
	for _, tt := range []struct {
		name        string
		blockString string
		want        pci.Devices
		errStr      string
	}{
		{
			name:        "one device",
			blockString: "0x8086:0x1234",
			want:        pci.Devices{&pci.PCI{Vendor: "0x8086", Device: "0x1234"}},
			errStr:      "",
		},
		{
			name:        "two devices",
			blockString: "0x8086:0x1234,0x1234:0xabcd",
			want: pci.Devices{
				&pci.PCI{Vendor: "0x8086", Device: "0x1234"},
				&pci.PCI{Vendor: "0x1234", Device: "0xabcd"},
			},
			errStr: "",
		},
		{
			name:        "no 0x",
			blockString: "8086:1234,1234:abcd",
			want: pci.Devices{
				&pci.PCI{Vendor: "0x8086", Device: "0x1234"},
				&pci.PCI{Vendor: "0x1234", Device: "0xabcd"},
			},
			errStr: "",
		},
		{
			name:        "capitals",
			blockString: "0x8086:0x1234,0x1234:0xABCD",
			want: pci.Devices{
				&pci.PCI{Vendor: "0x8086", Device: "0x1234"},
				&pci.PCI{Vendor: "0x1234", Device: "0xabcd"},
			},
			errStr: "",
		},
		{
			name:        "not hex vendor",
			blockString: "0xghij:0x1234",
			want:        nil,
			errStr:      "BlockList needs to contain a hex vendor ID, got 0xghij, err strconv.ParseUint: parsing \"ghij\": invalid syntax",
		},
		{
			name:        "not hex vendor",
			blockString: "0x1234:0xghij",
			want:        nil,
			errStr:      "BlockList needs to contain a hex device ID, got 0xghij, err strconv.ParseUint: parsing \"ghij\": invalid syntax",
		},
		{
			name:        "bad format",
			blockString: "0xghij,0x1234",
			want:        nil,
			errStr:      "BlockList needs to be of format vendor1:device1,vendor2:device2...! got 0xghij,0x1234",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			devices, err := parsePCIBlockList(tt.blockString)
			if e := testutil.CheckError(err, tt.errStr); e != nil {
				t.Error(e)
			}
			if !reflect.DeepEqual(devices, tt.want) {
				// Need to do this because stringer does not print device and vendor
				s := "got:\n"
				for _, d := range devices {
					s = fmt.Sprintf("%s{Vendor: %v, Device %v}\n", s, d.Vendor, d.Device)
				}
				s = fmt.Sprintf("%swant:\n", s)
				for _, d := range tt.want {
					s = fmt.Sprintf("%s{Vendor: %v, Device %v}\n", s, d.Vendor, d.Device)
				}
				t.Error(s)
			}
		})
	}
}
