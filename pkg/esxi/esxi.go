// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package esxi contains an ESXi boot config parser.
//
// This package can read ESXi boot configurations from disks or CDROMs.
package esxi

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/gpt"
	"github.com/u-root/u-root/pkg/mount"
)

// LoadOS loads an ESXi multiboot kernel from disk device's partition.
//
// If device is /dev/sda, and partition is 5, /dev/sda5 will be mounted at
// mountPoint.
func LoadOS(mountPoint string, device string, partition int) (*boot.MultibootImage, error) {
	partitionDev := fmt.Sprintf("%s%d", device, partition)
	if err := mount.Mount(partitionDev, mountPoint, "vfat", "", unix.MS_RDONLY|unix.MS_NOATIME); err != nil {
		return nil, err
	}
	return loadConfig(filepath.Join(mountPoint, "boot.cfg"), device, partition)
}

// LoadCDROM loads an ESXi multiboot kernel from a CDROM at device.
//
// device will be mounted at mountPoint.
func LoadCDROM(mountPoint string, device string) (*boot.MultibootImage, error) {
	if err := mount.Mount(device, mountPoint, "iso9660", "", unix.MS_RDONLY|unix.MS_NOATIME); err != nil {
		return nil, err
	}
	// Don't pass the device to ESXi. It doesn't need it.
	return loadConfig(filepath.Join(mountPoint, "boot.cfg"), "", 0)
}

// LoadConfig loads an ESXi configuration from configFile.
func LoadConfig(configFile string) (*boot.MultibootImage, error) {
	return loadConfig(configFile, "", 0)
}

func loadConfig(configFile, device string, partition int) (*boot.MultibootImage, error) {
	opts, err := parse(configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config at %s: %v", configFile, err)
	}
	if len(device) > 0 {
		if err := opts.addUUID(device, partition); err != nil {
			return nil, fmt.Errorf("cannot add boot uuid of %s: %v", device, err)
		}
	}

	return &boot.MultibootImage{
		Path:    opts.kernel,
		Cmdline: opts.args,
		Modules: opts.modules,
	}, nil
}

const (
	kernel  = "kernel"
	args    = "kernelopt"
	modules = "modules"

	comment = '#'
	sep     = "---"

	uuidMagic = "VMWARE FAT16    "
	uuidSize  = 32
)

type options struct {
	kernel  string
	args    string
	modules []string
}

func getUUID(device string, partition int) (string, error) {
	device = strings.TrimRight(device, "/")
	blockSize, err := gpt.GetBlockSize(device)
	if err != nil {
		return "", err
	}

	f, err := os.Open(fmt.Sprintf("%s%d", device, partition))
	if err != nil {
		return "", err
	}

	// Boot uuid is stored in the second block of the disk
	// in the following format:
	//
	// VMWARE FAT16    <uuid>
	// <---128 bit----><128 bit>
	data := make([]byte, uuidSize)
	n, err := f.ReadAt(data, int64(blockSize))
	if err != nil {
		return "", err
	}
	if n != uuidSize {
		return "", io.ErrUnexpectedEOF
	}

	if magic := string(data[:len(uuidMagic)]); magic != uuidMagic {
		return "", fmt.Errorf("bad uuid magic %q, want %q", magic, uuidMagic)
	}

	uuid := hex.EncodeToString(data[len(uuidMagic):])
	return fmt.Sprintf("bootUUID=%s", uuid), nil
}

func (o *options) addUUID(device string, partition int) error {
	uuid, err := getUUID(device, partition)
	if err != nil {
		return err
	}
	o.args += " " + uuid
	return nil
}

func parse(configFile string) (options, error) {
	dir := filepath.Dir(configFile)

	f, err := os.Open(configFile)
	if err != nil {
		return options{}, err
	}
	defer f.Close()

	var opt options

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == comment {
			continue
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return opt, fmt.Errorf("bad line %q", line)
		}
		key := strings.TrimSpace(tokens[0])
		val := strings.TrimSpace(tokens[1])
		switch key {
		case kernel:
			opt.kernel = filepath.Join(dir, val)
		case args:
			opt.args = val
		case modules:
			for _, tok := range strings.Split(val, sep) {
				// Each module is "filename arg0 arg1 arg2" and
				// the filename is relative to the directory
				// the module is in.
				tok = strings.TrimSpace(tok)
				if len(tok) > 0 {
					entry := strings.Fields(tok)
					entry[0] = filepath.Join(dir, entry[0])
					opt.modules = append(opt.modules, strings.Join(entry, " "))
				}
			}
		}
	}

	err = scanner.Err()
	return opt, err
}
