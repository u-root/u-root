// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"unicode"

	"github.com/bobuhiro11/gokvm/vmm"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

var (
	fwFilename = "CLOUDHV.fd"
	kvmPath    = "/dev/kvm"

	blockList = flag.String("block", "", "comma separated list of pci vendor and device ids to ignore (format vendor:device). E.g. 0x8086:0x1234,0x8086:0xabcd")
)

func main() {
	flag.Parse()

	cfg := vmm.Config{}
	cfg.Dev = kvmPath
	cfg.MemSize = (1 << 30)
	cfg.NCPUs = 1 // We dont want all CPUs for the VM so we leave 2 for the host.

	blockDevs, err := block.GetBlockDevices()
	if err != nil {
		log.Fatal("No available block devices to boot from")
	}

	// Try to only boot from "good" block devices.
	blockDevs = blockDevs.FilterZeroSize()

	// Parse and filter blocklist
	if *blockList != "" {
		blockDevs, err = blockDevs.FilterBlockPCIString(*blockList)
		if err != nil {
			log.Fatal(err)
		}
	}

	blockDevs = filterNonPartitions(blockDevs)

	if len(blockDevs) < 1 {
		log.Fatal("no block devices to mount")
	}

	var mountPool mount.Pool
	if err := mountBlockDevs(blockDevs, &mountPool); err != nil {
		log.Fatal(err)
	}
	defer mountPool.UnmountAll(0)

	mp, err := findFirmwareInBlockMounts(&mountPool)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Kernel = path.Join(mp.Path, fwFilename)

	fmt.Println(cfg.Kernel)

	vmm := vmm.New(cfg)

	if err := vmm.Init(); err != nil {
		log.Fatalf("vmm.Init failed: %v", err)
	}

	if err := vmm.Setup(); err != nil {
		log.Fatalf("vmm.Setup failed: %v", err)
	}

	if err := vmm.Boot(); err != nil {
		log.Fatalf("vmm.Boot failed: %v", err)
	}
}

func mountBlockDevs(devs block.BlockDevices, mp *mount.Pool) error {
	for _, dev := range devs {
		mpDir, err := os.MkdirTemp("", "vmboot-")
		if err != nil {
			return err
		}
		m, err := mount.TryMount(dev.DevicePath(), mpDir, "", 0)
		if err != nil {
			return err
		}
		mp.Add(m)
	}

	return nil
}

func filterNonPartitions(devs block.BlockDevices) block.BlockDevices {
	var ret block.BlockDevices
	for _, dev := range devs {
		if unicode.IsDigit(rune(dev.DevName()[len(dev.DevName())-1])) {
			ret = append(ret, dev)
		}
	}
	return ret
}

func findFirmwareInBlockMounts(mPool *mount.Pool) (*mount.MountPoint, error) {
	for _, mp := range mPool.MountPoints {
		fwFilePathCandidate := path.Join(mp.Path, fwFilename)
		_, err := os.Stat(fwFilePathCandidate)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		} else if errors.Is(err, fs.ErrNotExist) {
			continue
		} else {
			return mp, nil
		}
	}
	return nil, fmt.Errorf("no mount point with firmware image found")
}
