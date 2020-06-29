// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/mount"
)

const (
	dataPartitionFSType = "ext4"
	dataPartitionLabel  = "STDATA"
	dataMountPoint      = "data"
)

func findDataPartition() error {
	debug("Search data partition with label %s ...", dataPartitionLabel)
	fs, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return err
	}
	if !strings.Contains(string(fs), dataPartitionFSType) {
		return fmt.Errorf("filesystem unknown: %s", dataPartitionFSType)
	}

	devices, err := getBlockDevs()
	if err != nil {
		return fmt.Errorf("block devices: %v", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("no non-loopback block devices found")
	}

	device, err := deviceByPartLabel(devices, dataPartitionLabel)
	if err != nil {
		return err
	}

	mp, err := mount.Mount(device, dataMountPoint, dataPartitionFSType, "", 0)
	if err != nil {
		return err
	}

	debug("data partition %s mounted at %s", mp.Device, mp.Path)
	return nil
}

func getBlockDevs() ([]string, error) {
	devnames := make([]string, 0)
	root := "/sys/class/block"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if strings.Contains(rel, "loop") {
			return nil
		}
		dev := filepath.Join("/dev", rel)
		devnames = append(devnames, dev)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return devnames, nil
}

func deviceByPartLabel(devices []string, label string) (string, error) {
	var d string
	var p string
	for _, device := range devices {
		fd, err := os.Open(device)
		if err != nil {
			debug("Skip %s: %v", device, err)
			continue
		}
		defer fd.Close()
		if _, err = fd.Seek(512, io.SeekStart); err != nil {
			debug("Skip %s: %v", device, err)
			continue
		}
		table, err := gpt.ReadTable(fd, 512)
		if err != nil {
			debug("Skip %s: %v", device, err)
			continue
		}
		for n, part := range table.Partitions {
			if part.IsEmpty() {
				debug("Skip %s: no partitions found", device)
				continue
			}
			l, err := decodeLabel(part.PartNameUTF16[:])
			if err != nil {
				debug("Skip %s partition %d: %v", device, n+1, err)
				continue
			}
			if l == label {
				d = device
				p = strconv.Itoa(n + 1)
				info("Found data partition on %s , partition %s", device, p)
				break
			}
			debug("Skip %s partition %d: label does not match %s", device, n+1, label)
		}
		if d != "" && p != "" {
			break
		}
	}
	if d != "" && p != "" {
		for _, device := range devices {
			if !strings.HasPrefix(device, d) {
				continue
			}
			part := strings.TrimPrefix(device, d)
			if strings.Contains(part, p) {
				return device, nil
			}
		}
		return "", fmt.Errorf("Cannot find partition %s of %s in %v", p, d, devices)
	}
	return "", fmt.Errorf("No device with partition labeled %s found", label)
}

func decodeLabel(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("label has odd number of bytes")
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return strings.Trim(ret.String(), "\x00"), nil
}
