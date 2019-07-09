// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/storage"
)

// Parser for BLS boot entries. See spec at https://systemd.io/BOOT_LOADER_SPECIFICATION
// This only implements non-UEFI ("Type 1") boot entries

const (
	entryBaseDir = "loader"
	entryPath    = "loader/entries"
)

// ScanDevices searches the given list of devices for valid BLS entries.
// This implementation only supports disks with a GPT partition table.
// baseMountPoint is used for temporary mounts to inspect the filesystems.
// This function keeps devices with valid boot configs mounted when it returns so the paths are valid.
func ScanDevices(devices []storage.BlockDev, baseMountpoint string) ([]BootConfig, error) {
	// get a list of supported file systems for real devices (i.e. skip nodev)
	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	result := []BootConfig{}

	for _, d := range devices {
		mount := probeBLSRoot(d, baseMountpoint, filesystems)
		if mount != nil {
			configs, err := ScanBLSConfigs(mount.Path)
			if err != nil {
				syscall.Unmount(mount.Path, syscall.MNT_DETACH)
				// quiet for now. TODO fix
				continue
			}
			result = append(result, configs...)
		}
	}
	return result, nil
}

// probeBLSRoot attempts to mount a potential BLS root fs.
// Returns mounts with valid BLS directories, immediately unmounts invalid ones.
func probeBLSRoot(dev storage.BlockDev, baseMountpoint string, filesystems []string) *storage.Mountpoint {
	m, err := mountBLSRoot(dev, baseMountpoint, filesystems)
	if err != nil {
		// quiet for now. TODO: Fix
		return nil
	}
	if !validBLSMount(m) {
		syscall.Unmount(m.Path, syscall.MNT_DETACH)
		return nil
	}
	return m
}

func mountBLSRoot(dev storage.BlockDev, baseMountpoint string, filesystems []string) (*storage.Mountpoint, error) {
	devname := path.Join("/dev", dev.Name)
	mountpath := path.Join(baseMountpoint, dev.Name)
	mountpoint, err := storage.Mount(devname, mountpath, filesystems)
	if err != nil {
		return nil, fmt.Errorf("failed to mount %s on %s: %v", devname, mountpath, err)
	}
	return mountpoint, nil
}

// validates if the given mountpoint contains a valid BLS directory structure
// (that is, a "/loader/entries/" directory)
func validBLSMount(mount *storage.Mountpoint) bool {
	if _, err := os.Stat(path.Join(mount.Path, entryPath)); err != nil {
		// not a valid BLS or we can't read it
		return false
	}
	return true
}

// ScanBLSConfigs scans the filesystem root for valid BLS entries.
// This function skips over invalid or unreadable entries in an effort
// to return everything that is bootable
func ScanBLSConfigs(fsBase string) ([]BootConfig, error) {
	dir := path.Join(fsBase, entryPath)
	files, err := filepath.Glob(path.Join(dir, "*.conf"))
	if err != nil {
		return nil, err
	}
	result := []BootConfig{}
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println("Skipping BLS entry", f, " err: ", err)
			continue
		}
		conf, err := parseBLSEntry(string(data))
		if err != nil {
			fmt.Println("Skipping over invalid BLS entry", f, " err: ", err)
			continue
		}
		patchDirs(conf, fsBase)
		result = append(result, *conf)
	}
	return result, nil
}

// patchDirs converts the relative paths from the BLS config
// into absolute paths that can be used by a generic loader
func patchDirs(bc *BootConfig, fsBase string) {
	// no multiboot in BLS
	bc.Kernel = path.Join(fsBase, entryBaseDir, bc.Kernel)
	bc.Initramfs = path.Join(fsBase, entryBaseDir, bc.Initramfs)
}

// Parse takes the content of a Type #1 BLS entry and returns a BootConfig
// An error is returned if the syntax is wrong or required keys are missing
func parseBLSEntry(content string) (*BootConfig, error) {
	bc := &BootConfig{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		sline := strings.SplitN(line, " ", 2)
		if len(sline) != 2 {
			continue
		}
		key, val := sline[0], strings.TrimSpace(sline[1])

		switch key {
		case "title":
			bc.Name = val
		case "linux":
			bc.Kernel = val
		case "initrd":
			bc.Initramfs = val
		case "options":
			bc.KernelArgs = val
		}
	}
	// validate - spec says kernel and initrd are required
	if bc.Kernel == "" || bc.Initramfs == "" {
		return nil, fmt.Errorf("malformed BLS config, kernel or initrd missing")
	}
	return bc, nil
}
