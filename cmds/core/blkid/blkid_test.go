// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/mount/block"
)

func TestBlkid(t *testing.T) {
	for _, tt := range []struct {
		name         string
		BlockDevices []*block.BlockDev
		wantString   string
		want         error
	}{
		{
			name: "Mixed Block Devices",
			BlockDevices: []*block.BlockDev{
				{
					Name:   "nvme0n1p1",
					FsUUID: "51820b9c-d640-4c8c-8597-188689253e69",
				}, {
					Name:   "sda",
					FsUUID: "4c8c-8597",
				},
			},
			wantString: "/dev/nvme0n1p1 UUID=\"51820b9c-d640-4c8c-8597-188689253e69\"\n/dev/sda UUID=\"4c8c-8597\"\n",
			want:       nil,
		},
		{
			name: "Error Block Devices",
			BlockDevices: []*block.BlockDev{
				{
					Name:   "nvme0n1p1",
					FsUUID: "51820b9c-d640-4c8c-8597-188689253e69",
				}, {
					Name:   "sda",
					FsUUID: "4c8c-8597",
				},
			},
			want: fmt.Errorf("random error"),
		},
		{
			name: "Got FS Type",
			BlockDevices: []*block.BlockDev{
				{
					Name:   "nvme0n1p1",
					FSType: "Ext4",
					FsUUID: "51820b9c-d640-4c8c-8597-188689253e69",
				}, {
					Name:   "sda",
					FsUUID: "4c8c-8597",
				},
			},
			wantString: "/dev/nvme0n1p1 UUID=\"51820b9c-d640-4c8c-8597-188689253e69\" TYPE=\"Ext4\"\n/dev/sda UUID=\"4c8c-8597\"\n",
			want:       nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			blockGetBlockDevices := func() (block.BlockDevices, error) {
				return tt.BlockDevices, tt.want
			}
			var outBuf bytes.Buffer
			err := run(blockGetBlockDevices, &outBuf)
			if err != nil && !strings.Contains(err.Error(), tt.want.Error()) {
				t.Errorf("%q failed. Got '%v', want '%v'", tt.name, err, tt.want)
			}
			if outBuf.String() != tt.wantString {
				t.Errorf("Blkid.run() = '%s', want: '%s'", outBuf.String(), tt.wantString)
			}
		})
	}
}
