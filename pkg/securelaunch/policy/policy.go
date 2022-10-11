// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package policy locates and parses a JSON policy file.
package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/securelaunch/launcher"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
)

// Policy describes the policy used to drive the security engine.
//
// The policy is stored as a JSON file.
type Policy struct {
	Collectors []measurement.Collector
	Launcher   launcher.Launcher
	EventLog   eventlog.EventLog
}

// policyBytes is a byte slice to hold a copy of the policy file in memory.
var policyBytes []byte

// scanKernelCmdLine scans the kernel cmdline for the 'sl_policy' flag.
// When set it provides location of the policy file on disk. It reads it and
// returns the policy file as a byte slice.
//
// The format of sl_policy flag is as follows:
// sl_policy=<block device identifier>:<path>
// e.g., sda:/boot/securelaunch.policy
// e.g., 4qccd342-12zr-4e99-9ze7-1234cb1234c4:/foo/securelaunch.policy
func scanKernelCmdLine() []byte {
	slaunch.Debug("scanKernelCmdLine: scanning kernel cmd line for *sl_policy* flag")
	val, ok := cmdline.Flag("sl_policy")
	if !ok {
		log.Printf("scanKernelCmdLine: sl_policy cmdline flag is not set")
		return nil
	}

	// val is of type sda:path/to/file or UUID:path/to/file
	mntFilePath, e := slaunch.GetMountedFilePath(val, mount.MS_RDONLY) // false means readonly mount
	if e != nil {
		log.Printf("scanKernelCmdLine: GetMountedFilePath err=%v", e)
		return nil
	}
	slaunch.Debug("scanKernelCmdLine: Reading file=%s", mntFilePath)

	d, err := os.ReadFile(mntFilePath)
	if err != nil {
		log.Printf("Error reading policy file:mountPath=%s, passed=%s", mntFilePath, val)
		return nil
	}
	return d
}

// scanBlockDevice scans for the policy file on already mounted block devices.
// It looks in the "/", "/efi", and "/boot" directories. If found, it returns
// the policy file as a byte slice.
//
//	e.g., if /dev/sda1 is mounted on /tmp/sda1, then mountPath would be
//
// /tmp/sda1 and searchPath would be /tmp/sda1/securelaunch.policy,
// /tmp/sda1/efi/securelaunch.policy, and /tmp/sda1/boot/securelaunch.policy
// respectively for each iteration of loop over SearchRoots slice.
func scanBlockDevice(mountPath string) []byte {
	log.Printf("scanBlockDevice")
	// scan for securelaunch.policy under /, /efi, or /boot
	SearchRoots := []string{"/", "/efi", "/boot"}
	for _, c := range SearchRoots {

		searchPath := filepath.Join(mountPath, c, "securelaunch.policy")
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		d, err := os.ReadFile(searchPath)
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

// locate searches for the policy file on the kernel cmdline.
// If it's not found on cmdline, it looks for the policy file on each block
// device in the "/", "efi", and "/boot" directories.
func locate() ([]byte, error) {
	d := scanKernelCmdLine()
	if d != nil {
		return d, nil
	}

	slaunch.Debug("Searching for block devices")
	if err := slaunch.GetBlkInfo(); err != nil {
		return nil, err
	}

	// devName = sda, mountPath = /tmp/sluinit-FOO/
	for _, device := range slaunch.StorageBlkDevices {

		devName := device.Name
		mountPath, err := slaunch.MountDevice(device, mount.MS_RDONLY)
		if err != nil {
			log.Printf("failed to mount %s, continuing to next block device", devName)
			continue
		}

		slaunch.Debug("scanning for policy file under devName=%s, mountPath=%s", devName, mountPath)
		raw := scanBlockDevice(mountPath)
		if raw == nil {
			log.Printf("no policy file found under this device")
			continue
		}

		slaunch.Debug("policy file found at devName=%s", devName)
		return raw, nil
	}

	return nil, errors.New("policy file not found anywhere")
}

// parse accepts a JSON file as input, parses it into a well defined Policy
// structure and returns a pointer to the Policy structure.
func parse(pf []byte) (*Policy, error) {
	p := &Policy{}
	var parse struct {
		Collectors []json.RawMessage `json:"collectors"`
		Attestor   json.RawMessage   `json:"attestor"`
		Launcher   json.RawMessage   `json:"launcher"`
		EventLog   json.RawMessage   `json:"eventlog"`
	}

	if err := json.Unmarshal(pf, &parse); err != nil {
		log.Printf("parse SL Policy: Unmarshall error for entire policy file!! err=%v", err)
		return nil, err
	}

	for _, c := range parse.Collectors {
		collector, err := measurement.GetCollector(c)
		if err != nil {
			log.Printf("GetCollector err:c=%s, collector=%v", c, collector)
			return nil, err
		}
		p.Collectors = append(p.Collectors, collector)
	}

	if len(parse.Launcher) > 0 {
		if err := json.Unmarshal(parse.Launcher, &p.Launcher); err != nil {
			log.Printf("parse policy: Launcher Unmarshall error=%v!!", err)
			return nil, err
		}
	}

	if len(parse.EventLog) > 0 {
		if err := json.Unmarshal(parse.EventLog, &p.EventLog); err != nil {
			log.Printf("parse policy: EventLog Unmarshall error=%v!!", err)
			return nil, err
		}
	}
	return p, nil
}

// Measure measures the policy file.
func Measure() error {
	if len(policyBytes) == 0 {
		return fmt.Errorf("policy file not yet loaded or empty")
	}

	eventDesc := "File Collector: measured securelaunch policy file"
	if err := measurement.HashBytes(policyBytes, eventDesc); err != nil {
		log.Printf("policy: ERR: could not measure policy file: %v", err)
		return err
	}

	return nil
}

// Get locates and parses the policy file.
//
// The file is located by the following priority:
//  1. the kernel cmdline `sl_policy` argument.
//  2. a file on any partition on any disk called "securelaunch.policy".
func Get() (*Policy, error) {
	policyBytes, err := locate()
	if err != nil {
		return nil, err
	}

	policy, err := parse(policyBytes)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("no policy found")
	}

	return policy, nil
}
