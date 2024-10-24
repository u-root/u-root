// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot/jsonboot"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

// TODO backward compatibility for BIOS mode with partition type 0xee
// TODO use a proper parser for grub config (see grub.go)

var (
	flagBaseMountPoint = flag.String("m", "/mnt", "Base mount point where to mount partitions")
	flagDryRun         = flag.Bool("dryrun", false, "Do not actually kexec into the boot config")
	flagDebug          = flag.Bool("d", false, "Print debug output")
	flagConfigIdx      = flag.Int("config", -1, "Specify the index of the configuration to boot. The order is determined by the menu entries in the Grub config")
	flagGrubMode       = flag.Bool("grub", false, "Use GRUB mode, i.e. look for valid Grub/Grub2 configuration in default locations to boot a kernel. GRUB mode ignores -kernel/-initramfs/-cmdline")
	flagKernelPath     = flag.String("kernel", "", "Specify the path of the kernel to execute. If using -grub, this argument is ignored")
	flagInitramfsPath  = flag.String("initramfs", "", "Specify the path of the initramfs to load. If using -grub, this argument is ignored")
	flagKernelCmdline  = flag.String("cmdline", "", "Specify the kernel command line. If using -grub, this argument is ignored")
	flagDeviceGUID     = flag.String("guid", "", "GUID of the device where the kernel (and optionally initramfs) are located. Ignored if -grub is set or if -kernel is not specified")
)

var debug = func(string, ...interface{}) {}

// mountByGUID looks for a partition with the given GUID, and tries to mount it
// in a subdirectory under the specified mount point. The subdirectory has the
// same name of the device (e.g. /your/base/mountpoint/sda1).
// If more than one partition is found with the given GUID, the first that is
// found is used.
// This function returns a mount.Mountpoint object, or an error if any.
func mountByGUID(devices block.BlockDevices, guid, baseMountpoint string) (*mount.MountPoint, error) {
	log.Printf("Looking for partition with GUID %s", guid)
	partitions := devices.FilterPartType(guid)
	if len(partitions) == 0 {
		return nil, fmt.Errorf("no partitions with GUID %s", guid)
	}
	log.Printf("Partitions with GUID %s: %+v", guid, partitions)
	if len(partitions) > 1 {
		log.Printf("Warning: more than one partition found with the given GUID. Using the first one")
	}

	mountpath := filepath.Join(baseMountpoint, partitions[0].Name)
	return partitions[0].Mount(mountpath, mount.MS_RDONLY, func() error { return os.MkdirAll(mountpath, 0o666) })
}

// BootGrubMode tries to boot a kernel in GRUB mode. GRUB mode means:
// * look for the partition with the specified GUID, and mount it
// * if no GUID is specified, mount all of the specified devices
// * try to mount the device(s) using any of the kernel-supported filesystems
// * look for a GRUB configuration in various well-known locations
// * build a list of valid boot configurations from the found GRUB configuration files
// * try to boot every valid boot configuration until one succeeds
//
// The first parameter, `devices` is a list of block.BlockDev . The function
// will look for bootable configurations on these devices
// The second parameter, `baseMountPoint`, is the directory where the mount
// points for each device will be created.
// The third parameter, `guid`, is the partition GUID to look for. If it is an
// empty string, will search boot configurations on all of the specified devices
// instead.
// The fourth parameter, `dryrun`, will not boot the found configurations if set
// to true.
func BootGrubMode(devices block.BlockDevices, baseMountpoint string, guid string, dryrun bool, configIdx int) error {
	var mounted []*mount.MountPoint
	if guid == "" {
		// try mounting all the available devices, with all the supported file
		// systems
		debug("trying to mount all the available block devices with all the supported file system types")
		for _, dev := range devices {
			mountpath := filepath.Join(baseMountpoint, dev.Name)
			if mountpoint, err := dev.Mount(mountpath, mount.MS_RDONLY, func() error { return os.MkdirAll(mountpath, 0o666) }); err != nil {
				debug("Failed to mount %s on %s: %v", dev, mountpath, err)
			} else {
				mounted = append(mounted, mountpoint)
			}
		}
	} else {
		mount, err := mountByGUID(devices, guid, baseMountpoint)
		if err != nil {
			return err
		}
		mounted = append(mounted, mount)
	}

	log.Printf("mounted: %+v", mounted)
	defer func() {
		// clean up
		for _, mountpoint := range mounted {
			if err := mountpoint.Unmount(mount.MNT_DETACH); err != nil {
				debug("Failed to unmount %v: %v", mountpoint, err)
			}
		}
	}()

	// search for a valid grub config and extracts the boot configuration
	bootconfigs := make([]jsonboot.BootConfig, 0)
	for _, mountpoint := range mounted {
		bootconfigs = append(bootconfigs, ScanGrubConfigs(devices, mountpoint.Path)...)
	}
	if len(bootconfigs) == 0 {
		return fmt.Errorf("no boot configuration found")
	}
	log.Printf("Found %d boot configs", len(bootconfigs))
	for _, cfg := range bootconfigs {
		debug("%+v", cfg)
	}
	for n, cfg := range bootconfigs {
		log.Printf("  %d: %s\n", n, cfg.Name)
	}
	if configIdx > -1 {
		for n, cfg := range bootconfigs {
			if configIdx == n {
				if dryrun {
					debug("Dry-run mode: will not boot the found configuration")
					debug("Boot configuration: %+v", cfg)
					return nil
				}
				if err := cfg.Boot(); err != nil {
					log.Printf("Failed to boot kernel %s: %v", cfg.Kernel, err)
				}
			}
		}
		log.Printf("Invalid arg -config %d: there are only %d bootconfigs available\n", configIdx, len(bootconfigs))
		return nil
	}
	if dryrun {
		cfg := bootconfigs[0]
		debug("Dry-run mode: will not boot the found configuration")
		debug("Boot configuration: %+v", cfg)
		return nil
	}

	// try to kexec into every boot config kernel until one succeeds
	for _, cfg := range bootconfigs {
		debug("Trying boot configuration %+v", cfg)
		if err := cfg.Boot(); err != nil {
			log.Printf("Failed to boot kernel %s: %v", cfg.Kernel, err)
		}
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return nil
}

// BootPathMode tries to boot a kernel in PATH mode. This means:
// * look for a partition with the given GUID and mount it
// * look for the kernel and initramfs in the provided locations
// * boot the kernel with the provided command line
//
// The first parameter, `devices` is a list of block.BlockDev . The function
// will look for bootable configurations on these devices
// The second parameter, `baseMountPoint`, is the directory where the mount
// points for each device will be created.
// The third parameter, `guid`, is the partition GUID to look for.
// The fourth parameter, `dryrun`, will not boot the found configurations if set
// to true.
func BootPathMode(devices block.BlockDevices, baseMountpoint string, guid string, dryrun bool) error {
	mount, err := mountByGUID(devices, guid, baseMountpoint)
	if err != nil {
		return err
	}

	fullKernelPath := filepath.Join(mount.Path, *flagKernelPath)

	var fullInitramfsPath string
	if len(*flagInitramfsPath) != 0 {
		fullInitramfsPath = filepath.Join(mount.Path, *flagInitramfsPath)
	}

	cfg := jsonboot.BootConfig{
		Kernel:     fullKernelPath,
		Initramfs:  fullInitramfsPath,
		KernelArgs: *flagKernelCmdline,
	}
	debug("Trying boot configuration %+v", cfg)
	if dryrun {
		log.Printf("Dry-run, will not actually boot")
	} else {
		if err := cfg.Boot(); err != nil {
			return fmt.Errorf("failed to boot kernel %s: %w", cfg.Kernel, err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *flagGrubMode && *flagKernelPath != "" {
		log.Fatal("Options -grub and -kernel are mutually exclusive")
	}
	if *flagDebug {
		debug = log.Printf
	}

	// Get all the available block devices
	devices, err := block.GetBlockDevices()
	if err != nil {
		log.Fatal(err)
	}
	// print partition info
	if *flagDebug {
		for _, dev := range devices {
			log.Printf("Device: %+v", dev)
			table, err := dev.GPTTable()
			if err != nil {
				continue
			}
			log.Printf("  Table: %+v", table)
			for _, part := range table.Partitions {
				log.Printf("    Partition: %+v\n", part)
				if !part.IsEmpty() {
					log.Printf("      UUID: %s\n", part.Type.String())
				}
			}
		}
	}

	// TODO boot from EFI system partitions.

	if *flagGrubMode {
		if err := BootGrubMode(devices, *flagBaseMountPoint, *flagDeviceGUID, *flagDryRun, *flagConfigIdx); err != nil {
			log.Fatal(err)
		}
	} else if *flagKernelPath != "" {
		if err := BootPathMode(devices, *flagBaseMountPoint, *flagDeviceGUID, *flagDryRun); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("You must specify either -grub or -kernel")
	}
	os.Exit(1)
}
