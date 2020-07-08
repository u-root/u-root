// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package localboot contains helper functions for booting off local disks.
package localboot

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bls"
	"github.com/u-root/u-root/pkg/boot/esxi"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

// parse treats device as a block device with a file system.
func parse(l ulog.Logger, device *block.BlockDev, mountDir string) []boot.OSImage {
	imgs, err := bls.ScanBLSEntries(l, mountDir)
	if err != nil {
		l.Printf("No systemd-boot BootLoaderSpec configs found on %s, trying another format...: %v", device, err)
	}

	grubImgs, err := grub.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		l.Printf("No GRUB configs found on %s, trying another format...: %v", device, err)
	}
	imgs = append(imgs, grubImgs...)

	syslinuxImgs, err := syslinux.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		l.Printf("No syslinux configs found on %s: %v", device, err)
	}
	imgs = append(imgs, syslinuxImgs...)

	return imgs
}

// parseUnmounted treats device as unmounted, with or without partitions.
func parseUnmounted(l ulog.Logger, device *block.BlockDev) ([]boot.OSImage, []*mount.MountPoint) {
	// This will try to mount device partition 5 and 6.
	imgs, mps, err := esxi.LoadDisk(device.DevicePath())
	if err != nil {
		l.Printf("No ESXi disk configs found on %s: %v", device, err)
	}

	// This tries to mount the device itself, in case it's an installer CD.
	img, mp, err := esxi.LoadCDROM(device.DevicePath())
	if err != nil {
		l.Printf("No ESXi CDROM configs found on %s: %v", device, err)
	}
	if img != nil {
		imgs = append(imgs, img)
		mps = append(mps, mp)
	}
	// Convert from *MultibootImage to OSImage.
	var images []boot.OSImage
	for _, i := range imgs {
		images = append(images, i)
	}
	return images, mps
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot(l ulog.Logger, blockDevs block.BlockDevices) ([]boot.OSImage, []*mount.MountPoint, error) {
	mountPoints, err := ioutil.TempDir("", "u-root-boot")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create tmpdir: %v", err)
	}

	var images []boot.OSImage
	var mps []*mount.MountPoint
	for _, device := range blockDevs {
		imgs, mmps := parseUnmounted(l, device)
		if len(imgs) > 0 {
			images = append(images, imgs...)
			mps = append(mps, mmps...)
		} else {
			dir := filepath.Join(mountPoints, device.Name)

			os.MkdirAll(dir, 0777)
			mp, err := device.Mount(dir, mount.ReadOnly)
			if err != nil {
				os.RemoveAll(dir)
				continue
			}

			imgs = parse(l, device, dir)
			images = append(images, imgs...)
			mps = append(mps, mp)
		}
	}
	return images, mps, nil
}
