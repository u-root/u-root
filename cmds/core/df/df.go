// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// df reports details of mounted filesystems.
//
// Synopsis
//  df [-k] [-m]
//
// Description
//  read mount information from /proc/mounts and
//  statfs syscall and display summary information for all
//  mount points that have a non-zero block count.
//  Users can choose to see the diplay in KB or MB.
//
// Options
//  -k: display values in KB (default)
//  -m: dispaly values in MB
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"syscall"
)

var (
	inKB  = flag.Bool("k", false, "Express the values in kilobytes (default)")
	inMB  = flag.Bool("m", false, "Express the values in megabytes")
	units uint64
)

const procmountsFile = "/proc/mounts"

const (
	// B is Bytes
	B = 1
	// KB is kilobytes
	KB = 1024 * B
	// MB is megabytes
	MB = 1024 * KB
)

// Mount is a structure used to contain mount point data
type Mount struct {
	Device         string
	MountPoint     string
	FileSystemType string
	Flags          string
	Bsize          int64
	Blocks         uint64
	Total          uint64
	Used           uint64
	Avail          uint64
	PCT            uint8
}

type mountinfomap map[string]Mount

// mountinfo returns a map of mounts representing
// the data in /proc/mounts
func mountinfo() (mountinfomap, error) {
	buf, err := ioutil.ReadFile(procmountsFile)
	if err != nil {
		return nil, err
	}
	return mountinfoFromBytes(buf)
}

// returns a map generated from the bytestream returned
// from /proc/mounts
// for tidiness, we decide to ignore filesystems of size 0
// to exclude cgroup, procfs and sysfs types
func mountinfoFromBytes(buf []byte) (mountinfomap, error) {
	ret := make(mountinfomap)
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		kv := bytes.SplitN(line, []byte{' '}, 6)
		if len(kv) != 6 {
			// can't interpret this
			continue
		}
		key := string(kv[1])
		var mnt Mount
		mnt.Device = string(kv[0])
		mnt.MountPoint = string(kv[1])
		mnt.FileSystemType = string(kv[2])
		mnt.Flags = string(kv[3])
		DiskUsage(&mnt)
		if mnt.Blocks == 0 {
			continue
		} else {
			ret[key] = mnt
		}
	}
	return ret, nil
}

// DiskUsage calculates the usage statistics of a mount point
// note: arm7 Bsize is int32; all others are int64
func DiskUsage(mnt *Mount) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(mnt.MountPoint, &fs)
	if err != nil {
		return
	}
	mnt.Blocks = fs.Blocks * uint64(fs.Bsize) / units
	mnt.Bsize = int64(fs.Bsize)
	mnt.Total = fs.Blocks * uint64(fs.Bsize) / units
	mnt.Avail = fs.Bavail * uint64(fs.Bsize) / units
	mnt.Used = (fs.Blocks - fs.Bfree) * uint64(fs.Bsize) / units
	pct := float64((fs.Blocks - fs.Bfree)) * 100 / float64(fs.Blocks)
	mnt.PCT = uint8(math.Ceil(pct))
}

// SetUnits takes the command line flags and configures
// the correct units used to calculate display values
func SetUnits() {
	if *inKB && *inMB {
		log.Fatal("options -k and -m are mutually exclusive")
	}
	if *inMB {
		units = MB
	} else {
		units = KB
	}
}

func df() {
	SetUnits()
	mounts, _ := mountinfo()
	var blocksize = "1K"
	if *inMB {
		blocksize = "1M"
	}
	fmt.Printf("Filesystem           Type         %v-blocks       Used    Available  Use%% Mounted on\n", blocksize)
	for _, mnt := range mounts {
		fmt.Printf("%-20v %-9v %12v %10v %12v %4v%% %-13v\n",
			mnt.Device,
			mnt.FileSystemType,
			mnt.Blocks,
			mnt.Used,
			mnt.Avail,
			mnt.PCT,
			mnt.MountPoint)
	}
}

func main() {
	flag.Parse()
	df()
}
