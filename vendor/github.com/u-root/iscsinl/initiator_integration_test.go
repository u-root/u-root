// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package iscsinl

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	// This iSCSI target is a disaster, but it works without any kernel
	// modules.
	iscsiconfig "github.com/gostor/gotgt/pkg/config"
	_ "github.com/gostor/gotgt/pkg/port/iscsit"
	"github.com/gostor/gotgt/pkg/scsi"
	_ "github.com/gostor/gotgt/pkg/scsi/backingstore"

	"github.com/hugelgupf/p9/fsimpl/test/vmdriver"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

func ISCSIConfig(filename, volume, portal string) *iscsiconfig.Config {
	return &iscsiconfig.Config{
		Storages: []iscsiconfig.BackendStorage{
			{
				DeviceID:         1000,
				Path:             fmt.Sprintf("file:%s", filename),
				Online:           true,
				ThinProvisioning: true,
				BlockShift:       0,
			},
		},
		ISCSIPortals: []iscsiconfig.ISCSIPortalInfo{
			{
				ID:     0,
				Portal: portal,
			},
		},
		ISCSITargets: map[string]iscsiconfig.ISCSITarget{
			volume: {
				TPGTs: map[string][]uint64{
					// ???
					"1": {0},
				},
				LUNs: map[string]uint64{
					// Map LUN ID to DeviceID (see above).
					"1": 1000,
				},
			},
		},
	}
}

func StartISCSIService(config *iscsiconfig.Config) error {
	if err := scsi.InitSCSILUMap(config); err != nil {
		return err
	}

	scsiTarget := scsi.NewSCSITargetService()
	targetDriver, err := scsi.NewTargetDriver("iscsi", scsiTarget)
	if err != nil {
		return err
	}

	for tgtname := range config.ISCSITargets {
		if err := targetDriver.NewTarget(tgtname, config); err != nil {
			return err
		}
	}

	// comment this to avoid concurrent issue
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// run a service
	go targetDriver.Run() //nolint:errcheck
	return nil
}

func TestIntegration(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skipf("test")
	}
	volumeName := "iqn.2016-09.com.example:foobar"

	// 1MB.ext4_vfat has two partitions:
	//
	//   part 1 ext4
	//   part 2 vfat
	if err := StartISCSIService(ISCSIConfig("./testdata/1MB.ext4_vfat", volumeName, "127.0.0.1:3260")); err != nil {
		t.Fatalf("iscsi target failed to start: %v", err)
	}

	vmtest.GolangTest(t,
		[]string{"github.com/u-root/iscsinl"},
		&vmtest.Options{
			BuildOpts: uroot.Opts{
				Commands: uroot.BusyBoxCmds(
					"github.com/u-root/u-root/cmds/core/dhclient",
					"github.com/u-root/u-root/cmds/core/ls",
				),
			},
			QEMUOpts: qemu.Options{
				Devices: []qemu.Device{
					vmdriver.HostNetwork{
						Net: &net.IPNet{
							// 192.168.0.0/24
							IP:   net.IP{192, 168, 0, 0},
							Mask: net.CIDRMask(24, 32),
						},
					},
				},
				KernelArgs: fmt.Sprintf("TGT_PORT=3260 TGT_SERVER=192.168.0.2 TGT_VOLUME=%s -v", volumeName),
				Timeout:    30 * time.Second,
			},
		},
	)
}
