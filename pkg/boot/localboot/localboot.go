// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package localboot contains helper functions for booting off local disks.
package localboot

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bls"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

func parse(device *block.BlockDev, mountDir string) []boot.OSImage {
	imgs, err := bls.ScanBLSEntries(ulog.Log, mountDir)
	if err != nil {
		log.Printf("Failed to parse systemd-boot BootLoaderSpec configs, trying another format...: %v", err)
	}

	grubImgs, err := grub.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		log.Printf("Failed to parse GRUB configs from %s, trying another format...: %v", device, err)
	}
	imgs = append(imgs, grubImgs...)

	syslinuxImgs, err := syslinux.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		log.Printf("Failed to parse syslinux configs from %s: %v", device, err)
	}
	imgs = append(imgs, syslinuxImgs...)

	return imgs
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot() ([]boot.OSImage, []*mount.MountPoint, error) {
	blockDevs, err := block.GetBlockDevices()
	if err != nil {
		return nil, nil, errors.New("no available block devices to boot from")
	}

	// Try to only boot from "good" block devices.
	blockDevs = blockDevs.FilterZeroSize()
	log.Printf("Booting from the following block devices: %v", blockDevs)

	mountPoints, err := ioutil.TempDir("", "u-root-boot")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create tmpdir: %v", err)
	}

	var images []boot.OSImage
	var mps []*mount.MountPoint
	for _, device := range blockDevs {
		dir := filepath.Join(mountPoints, device.Name)

		os.MkdirAll(dir, 0777)
		mp, err := device.Mount(dir, mount.ReadOnly)
		if err != nil {
			continue
		}

		imgs := parse(device, dir)
		images = append(images, imgs...)
		mps = append(mps, mp)
	}
	return images, mps, nil
}
