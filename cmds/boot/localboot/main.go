// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/storage"
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
// The specified filesystems will be used in the mount attempts.
// If more than one partition is found with the given GUID, the first that is
// found is used.
// This function returns a storage.Mountpoint object, or an error if any.
func mountByGUID(devices []storage.BlockDev, filesystems []string, guid, baseMountpoint string) (*storage.Mountpoint, error) {
	log.Printf("Looking for partition with GUID %s", guid)
	partitions, err := storage.PartitionsByGUID(devices, guid)
	if err != nil || len(partitions) == 0 {
		return nil, fmt.Errorf("Error looking up for partition with GUID %s", guid)
	}
	log.Printf("Partitions with GUID %s: %+v", guid, partitions)
	if len(partitions) > 1 {
		log.Printf("Warning: more than one partition found with the given GUID. Using the first one")
	}
	dev := partitions[0]
	mountpath := path.Join(baseMountpoint, dev.Name)
	devname := path.Join("/dev", dev.Name)
	mountpoint, err := storage.Mount(devname, mountpath, filesystems)
	if err != nil {
		return nil, fmt.Errorf("mountByGUID: cannot mount %s (GUID %s) on %s: %v", devname, guid, mountpath, err)
	}
	return mountpoint, nil
}

// BootGrubMode tries to boot a kernel in GRUB mode. GRUB mode means:
// * look for the partition with the specified GUID, and mount it
// * if no GUID is specified, mount all of the specified devices
// * try to mount the device(s) using any of the kernel-supported filesystems
// * look for a GRUB configuration in various well-known locations
// * build a list of valid boot configurations from the found GRUB configuration files
// * try to boot every valid boot configuration until one succeeds
//
// The first parameter, `devices` is a list of storage.BlockDev . The function
// will look for bootable configurations on these devices
// The second parameter, `baseMountPoint`, is the directory where the mount
// points for each device will be created.
// The third parameter, `guid`, is the partition GUID to look for. If it is an
// empty string, will search boot configurations on all of the specified devices
// instead.
// The fourth parameter, `dryrun`, will not boot the found configurations if set
// to true.
func BootGrubMode(devices []storage.BlockDev, baseMountpoint string, guid string, dryrun bool, configIdx int) error {
	// get a list of supported file systems for real devices (i.e. skip nodev)
	debug("Getting list of supported filesystems")
	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	debug("Supported file systems: %v", filesystems)

	var mounted []storage.Mountpoint
	if guid == "" {
		// try mounting all the available devices, with all the supported file
		// systems
		debug("trying to mount all the available block devices with all the supported file system types")
		mounted = make([]storage.Mountpoint, 0)
		for _, dev := range devices {
			devname := path.Join("/dev", dev.Name)
			mountpath := path.Join(baseMountpoint, dev.Name)
			if mountpoint, err := storage.Mount(devname, mountpath, filesystems); err != nil {
				debug("Failed to mount %s on %s: %v", devname, mountpath, err)
			} else {
				mounted = append(mounted, *mountpoint)
			}
		}
		log.Printf("mounted: %+v", mounted)
		defer func() {
			// clean up
			for _, mountpoint := range mounted {
				syscall.Unmount(mountpoint.Path, syscall.MNT_DETACH)
			}
		}()
	} else {
		mount, err := mountByGUID(devices, filesystems, guid, baseMountpoint)
		if err != nil {
			return err
		}
		mounted = []storage.Mountpoint{*mount}
	}

	// search for a valid grub config and extracts the boot configuration
	bootconfigs := make([]bootconfig.BootConfig, 0)
	for _, mountpoint := range mounted {
		bootconfigs = append(bootconfigs, ScanGrubConfigs(devices, mountpoint.Path)...)
	}
	if len(bootconfigs) == 0 {
		return fmt.Errorf("No boot configuration found")
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
// The first parameter, `devices` is a list of storage.BlockDev . The function
// will look for bootable configurations on these devices
// The second parameter, `baseMountPoint`, is the directory where the mount
// points for each device will be created.
// The third parameter, `guid`, is the partition GUID to look for.
// The fourth parameter, `dryrun`, will not boot the found configurations if set
// to true.
func BootPathMode(devices []storage.BlockDev, baseMountpoint string, guid string, dryrun bool) error {
	debug("Getting list of supported filesystems")
	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	debug("Supported file systems: %v", filesystems)

	mount, err := mountByGUID(devices, filesystems, guid, baseMountpoint)
	if err != nil {
		return err
	}

	fullKernelPath := path.Join(mount.Path, *flagKernelPath)
	fullInitramfsPath := path.Join(mount.Path, *flagInitramfsPath)
	cfg := bootconfig.BootConfig{
		Kernel:     fullKernelPath,
		Initramfs:  fullInitramfsPath,
		KernelArgs: *flagKernelCmdline,
	}
	debug("Trying boot configuration %+v", cfg)
	if dryrun {
		log.Printf("Dry-run, will not actually boot")
	} else {
		if err := cfg.Boot(); err != nil {
			return fmt.Errorf("Failed to boot kernel %s: %v", cfg.Kernel, err)
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
	devices, err := storage.GetBlockStats()
	if err != nil {
		log.Fatal(err)
	}
	// print partition info
	if *flagDebug {
		for _, dev := range devices {
			log.Printf("Device: %+v", dev)
			table, err := storage.GetGPTTable(dev)
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

	// TODO boot from EFI system partitions. See storage.FilterEFISystemPartitions

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
