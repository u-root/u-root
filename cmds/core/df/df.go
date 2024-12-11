// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// df reports details of mounted filesystems.
//
// Synopsis
//
//	df [-k] [-m]
//
// Description
//
//	read mount information from /proc/mounts and
//	statfs syscall and display summary information for all
//	mount points that have a non-zero block count.
//	Users can choose to see the diplay in KB or MB.
//
// Options
//
//	-k: display values in KB (default)
//	-m: dispaly values in MB
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"syscall"
)

type flags struct {
	k bool
	m bool
}

var (
	fargs = flags{}
	units uint64

	errKMExclusiv = errors.New("options -k and -m are mutually exclusive")
)

func init() {
	flag.BoolVar(&fargs.k, "k", false, "Express the values in kilobytes (default)")
	flag.BoolVar(&fargs.m, "m", false, "Express the values in megabytes")
}

const (
	// B is Bytes
	B = 1
	// KB is kilobytes
	KB = 1024 * B
	// MB is megabytes
	MB = 1024 * KB

	procmountsFile = "/proc/mounts"
)

// Mount is a structure used to contain mount point data
type mount struct {
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

type mountinfomap map[string]mount

// mountinfo returns a map of mounts representing
// the data in /proc/mounts
func mountinfo() (mountinfomap, error) {
	buf, err := os.ReadFile(procmountsFile)
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
		var mnt mount
		mnt.Device = string(kv[0])
		mnt.MountPoint = string(kv[1])
		mnt.FileSystemType = string(kv[2])
		mnt.Flags = string(kv[3])
		if err := diskUsage(&mnt); err != nil {
			return nil, err
		}
		if mnt.Blocks == 0 {
			continue
		}
		ret[key] = mnt
	}
	return ret, nil
}

// diskUsage calculates the usage statistics of a mount point
// note: arm7 Bsize is int32; all others are int64
func diskUsage(mnt *mount) error {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(mnt.MountPoint, &fs); err != nil {
		// skip mount point if df is running without root
		if os.IsPermission(err) {
			return nil
		}
		return err
	}
	mnt.Blocks = fs.Blocks * uint64(fs.Bsize) / units
	mnt.Bsize = int64(fs.Bsize)
	mnt.Total = fs.Blocks * uint64(fs.Bsize) / units
	mnt.Avail = uint64(fs.Bavail) * uint64(fs.Bsize) / units
	mnt.Used = (fs.Blocks - fs.Bfree) * uint64(fs.Bsize) / units
	pct := float64((fs.Blocks - fs.Bfree)) * 100 / float64(fs.Blocks)
	mnt.PCT = uint8(math.Ceil(pct))
	return nil
}

// setUnits takes the command line flags and configures
// the correct units used to calculate display values
func setUnits(inKB, inMB bool) error {
	if inKB && inMB {
		return errKMExclusiv
	}
	if inMB {
		units = MB
	} else {
		units = KB
	}
	return nil
}

func printHeader(w io.Writer, blockSize string) {
	fmt.Fprintf(w, "Filesystem           Type         %v-blocks       Used    Available  Use%% Mounted on\n", blockSize)
}

func printMount(w io.Writer, mnt mount) {
	fmt.Fprintf(w, "%-20v %-9v %12v %10v %12v %4v%% %-13v\n",
		mnt.Device,
		mnt.FileSystemType,
		mnt.Blocks,
		mnt.Used,
		mnt.Avail,
		mnt.PCT,
		mnt.MountPoint)
}

func df(w io.Writer, fargs flags, args []string) error {
	if err := setUnits(fargs.k, fargs.m); err != nil {
		return err
	}
	mounts, err := mountinfo()
	if err != nil {
		return fmt.Errorf("mountinfo()=_,%w, want: _,nil", err)
	}
	blocksize := "1K"
	if fargs.m {
		blocksize = "1M"
	}

	if len(args) == 0 {
		printHeader(w, blocksize)
		for _, mnt := range mounts {
			printMount(w, mnt)
		}

		return nil
	}

	var fileDevs []uint64
	for _, arg := range args {
		fileDev, err := deviceNumber(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "df: %v\n", err)
			continue
		}

		fileDevs = append(fileDevs, fileDev)
	}

	showHeader := true
	for _, mnt := range mounts {
		stDev, err := deviceNumber(mnt.MountPoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "df: %v\n", err)
			continue
		}

		for _, fDev := range fileDevs {
			if fDev == stDev {
				if showHeader {
					printHeader(w, blocksize)
					showHeader = false
				}
				printMount(w, mnt)
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if err := df(os.Stdout, fargs, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
