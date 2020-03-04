// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/storage"
)

const dataPartitionFSType = "ext4"

type dataPartition interface {
	get(filename string) ([]byte, error)
}

type partition struct {
	mountpoint string
}

func (p *partition) get(filename string) ([]byte, error) {
	f := filepath.Join(p.mountpoint, filename)
	return ioutil.ReadFile(f)
}

func findDataPartition() (dataPartition, error) {
	fs, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return nil, err
	}
	if !strings.Contains(string(fs), dataPartitionFSType) {
		return nil, fmt.Errorf("filesystem unknown: %s", dataPartitionFSType)
	}

	devices, err := storage.GetBlockStats()
	if err != nil {
		return nil, fmt.Errorf("no block devices: %v", err)
	}

	var mounted []*mount.MountPoint
	for _, dev := range devices {
		if strings.Contains(dev.Name, "loop") {
			continue
		}
		devname := filepath.Join("/dev", dev.Name)
		path := filepath.Join("/mnt", dev.Name)
		mp, err := mount.Mount(devname, path, dataPartitionFSType, "", 0)
		if err != nil {
			debug("Skip %s: %v", devname, err)
			continue
		}
		mounted = append(mounted, mp)
	}
	var p partition
	for _, mp := range mounted {
		f := filepath.Join(mp.Path, provisioningServerFile)
		if _, err := os.Stat(f); err != nil {
			debug("Skip %s : %v", mp.Device, err)
			continue
		}
		debug("data partition %s mounted at %s", mp.Device, mp.Path)
		p.mountpoint = mp.Path
		return &p, nil
	}
	return nil, errors.New("No stboot data partition found")
}
