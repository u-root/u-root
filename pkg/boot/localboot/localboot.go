// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package localboot contains helper functions for booting off local disks.
package localboot

import (
	"context"
	"sort"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bls"
	"github.com/u-root/u-root/pkg/boot/esxi"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

// Sort the image in descending order by rank
type byRank []boot.OSImage

func (a byRank) Less(i, j int) bool { return a[i].Rank() > a[j].Rank() }
func (a byRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRank) Len() int           { return len(a) }

// parse treats device as a block device with a file system.
func parse(l ulog.Logger, device *block.BlockDev, devices block.BlockDevices, mountDir string, mountPool *mount.Pool) []boot.OSImage {
	imgs, err := bls.ScanBLSEntries(l, mountDir, nil, "")
	if err != nil {
		l.Printf("No systemd-boot BootLoaderSpec configs found on %s, trying another format...: %v", device, err)
	}

	// Grub parser may want to load files (kernel, initramfs, modules, ...)
	// from another partition, thus it is given devices and mountPool in
	// order to reuse mounts and mount more file systems.
	grubImgs, err := grub.ParseLocalConfig(context.Background(), mountDir, devices, mountPool)
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
func parseUnmounted(l ulog.Logger, device *block.BlockDev, mountPool *mount.Pool) []boot.OSImage {
	// This will try to mount device partition 5 and 6.
	imgs, mps, err := esxi.LoadDisk(device.DevicePath())
	if mps != nil {
		mountPool.Add(mps...)
	}
	if err != nil {
		l.Printf("No ESXi disk configs found on %s: %v", device, err)
	}

	// This tries to mount the device itself, in case it's an installer CD.
	img, mp, err := esxi.LoadCDROM(device.DevicePath())
	if mp != nil {
		mountPool.Add(mp)
	}
	if err != nil {
		l.Printf("No ESXi CDROM configs found on %s: %v", device, err)
	}
	if img != nil {
		imgs = append(imgs, img)
	}
	// Convert from *MultibootImage to OSImage.
	var images []boot.OSImage
	for _, i := range imgs {
		images = append(images, i)
	}
	return images
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot(l ulog.Logger, blockDevs block.BlockDevices, mp *mount.Pool) ([]boot.OSImage, error) {
	var images []boot.OSImage
	for _, device := range blockDevs {
		imgs := parseUnmounted(l, device, mp)
		if len(imgs) > 0 {
			images = append(images, imgs...)
		} else {
			m, err := mp.Mount(device, mount.ReadOnly)
			if err != nil {
				continue
			}
			imgs = parse(l, device, blockDevs, m.Path, mp)
			images = append(images, imgs...)
		}
	}

	sort.Sort(byRank(images))
	return images, nil
}
