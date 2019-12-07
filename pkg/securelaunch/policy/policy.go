// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package policy locates and parses a JSON policy file.
package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot/diskboot"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
)

/*
 * Policy describes the TPM measurements to take and the OS to boot.
 *
 * The policy is stored as a JSON file.
 */
type Policy struct {
	DefaultAction string
	Collectors    []measurement.Collector
}

/*
 * scanKernelCmdLine scans the kernel cmdline
 * for 'sl_policy' flag. when set, this flag provides location of
 * of policy file on disk enabling the function to return policy file as
 * a byte slice.
 *
 * format of sl_policy flag is as follows
 * sl_policy=<block device identifier>:<path>
 * e.g sda:/boot/securelaunch.policy
 * e.g 4qccd342-12zr-4e99-9ze7-1234cb1234c4:/foo/securelaunch.policy
 */
func scanKernelCmdLine() []byte {

	slaunch.Debug("scanKernelCmdLine: scanning kernel cmd line for *sl_policy* flag")
	val, ok := cmdline.Flag("sl_policy")
	if !ok {
		log.Printf("scanKernelCmdLine: sl_policy cmdline flag is not set")
		return nil
	}

	// val is of type sda:path/to/file or UUID:path/to/file
	mntFilePath, mountPath, e := slaunch.GetMountedFilePath(val, mount.MS_RDONLY) // false means readonly mount
	if e != nil {
		log.Printf("scanKernelCmdLine: GetMountedFilePath err=%v", e)
		return nil
	}
	slaunch.Debug("scanKernelCmdLine: Reading file=%s", mntFilePath)

	d, err := ioutil.ReadFile(mntFilePath)
	if e := mount.Unmount(mountPath, true, false); e != nil {
		log.Printf("Unmount failed. PANIC")
		panic(e)
	}

	if err != nil {
		log.Printf("Error reading policy file:mountPath=%s, passed=%s", mntFilePath, val)
		return nil
	}
	return d
}

/*
 *  scanBlockDevice scans an already mounted block device inside directories
 *	"/", "/efi" and "/boot" for policy file and if found, returns the policy byte as a byte slice.
 *
 *	e.g: if you mount /dev/sda1 on /tmp/sda1,
 *	then mountPath would be /tmp/sda1
 *	and searchPath would be /tmp/sda1/securelaunch.policy,
 * /tmp/sda1/efi/securelaunch.policy and /tmp/sda1/boot/securelaunch.policy
 *	respectively for each iteration of loop over SearchRoots slice.
 */
func scanBlockDevice(mountPath string) []byte {

	log.Printf("scanBlockDevice")
	// scan for securelaunch.policy under /, /efi, or /boot
	var SearchRoots = []string{"/", "/efi", "/boot"}
	for _, c := range SearchRoots {

		searchPath := filepath.Join(mountPath, c, "securelaunch.policy")
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		d, err := ioutil.ReadFile(searchPath)
		if err != nil {
			// Policy File not found. Moving on to next search root...
			log.Printf("Error reading policy file %s, continuing", searchPath)
			continue
		}
		log.Printf("policy file found on mountPath=%s, directory =%s", mountPath, c)
		return d // return when first policy file found
	}

	return nil
}

/*
 * locate searches for policy file on the kernel cmdline.
 * if not found on cmdline, it looks for policy file on each block device
 * under "/", "efi" and "/boot" directories.
 *
 * Steps:
 * 1. Check if kernel param sl_policy is set,
 * 		parse the string
 * 2. Iterate through each local block device,
 *	- mount the block device
 *	- scan for securelaunch.policy under /, /efi, or /boot
 * 3  Read in policy file
 */
func locate() ([]byte, error) {

	d := scanKernelCmdLine()
	if d != nil {
		return d, nil
	}

	slaunch.Debug("Searching and mounting block devices with bootable configs")
	blkDevices := diskboot.FindDevices("/sys/class/block/*") // FindDevices find and *mounts* the devices.
	if len(blkDevices) == 0 {
		return nil, errors.New("no block devices found")
	}

	for _, device := range blkDevices {
		devicePath, mountPath := device.MountPoint.Device, device.MountPoint.Path
		slaunch.Debug("scanning for policy file under devicePath=%s, mountPath=%s", devicePath, mountPath)
		raw := scanBlockDevice(mountPath)
		if e := mount.Unmount(mountPath, true, false); e != nil {
			log.Printf("Unmount failed. PANIC")
			panic(e)
		}

		if raw == nil {
			log.Printf("no policy file found under this device")
			continue
		}

		slaunch.Debug("policy file found at devicePath=%s", devicePath)
		return raw, nil
	}

	return nil, errors.New("policy file not found anywhere")
}

/*
 * parse accepts a JSON file as input, parses
 * it into a well defined Policy structure (parse) and
 * returns a pointer to Policy structure.
 */
func parse(pf []byte) (*Policy, error) {
	p := &Policy{}
	var parse struct {
		DefaultAction string            `json:"default_action"`
		Collectors    []json.RawMessage `json:"collectors"`
		Attestor      json.RawMessage   `json:"attestor"`
		Launcher      json.RawMessage   `json:"launcher"`
		EventLog      json.RawMessage   `json:"eventlog"`
	}

	if err := json.Unmarshal(pf, &parse); err != nil {
		log.Printf("parse SL Policy: Unmarshall error for entire policy file!! err=%v", err)
		return nil, err
	}

	p.DefaultAction = parse.DefaultAction

	for _, c := range parse.Collectors {
		collector, err := measurement.GetCollector(c)
		if err != nil {
			log.Printf("GetCollector err:c=%s, collector=%v", c, collector)
			return nil, err
		}
		p.Collectors = append(p.Collectors, collector)
	}

	return p, nil
}

/*
 * Get locates and parses the policy file.
 *
 * The file is located by the following priority:
 *
 *  (1) kernel cmdline "sl_policy" argument.
 *  (2) a file on any partition on any disk called "securelaunch.policy"
 */
func Get() (*Policy, error) {
	b, err := locate()
	if err != nil {
		return nil, err
	}
	policy, err := parse(b)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("no policy found")
	}
	return policy, nil
}
