// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/storage"
)

const (
	dataPartitionFSType = "ext4"
	dataPartitionLabel  = "STDATA"
)

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

	devices, err = storage.PartitionsByLable(devices, "STDATA")
	if err != nil || len(devices) == 0 {
		return nil, fmt.Errorf("no partitions with label %s", dataPartitionLabel)
	}
	if len(devices) > 1 {
		debug("WARNING: multiple data partitions found! Take %s", devices[0].Name)
	}

	devname := filepath.Join("/dev", devices[0].Name)
	path := filepath.Join("/mnt", devices[0].Name)
	mp, err := mount.Mount(devname, path, dataPartitionFSType, "", 0)
	if err != nil {
		return nil, err
	}

	debug("data partition %s mounted at %s", mp.Device, mp.Path)
	var p partition
	p.mountpoint = mp.Path
	return &p, nil
}
