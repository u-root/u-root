// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/diskboot"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/mount"
)

var (
	v       = flag.Bool("v", false, "Print debug messages")
	verbose = func(string, ...interface{}) {}
	dryrun  = flag.Bool("dryrun", false, "Only print out kexec commands")

	devGlob           = flag.String("dev", "/sys/class/block/*", "Device glob")
	sDeviceIndex      = flag.String("d", "", "Device index")
	sConfigIndex      = flag.String("c", "", "Config index")
	sEntryIndex       = flag.String("n", "", "Entry index")
	removeCmdlineItem = flag.String("remove", "console", "comma separated list of kernel params value to remove from parsed kernel configuration (default to console)")
	reuseCmdlineItem  = flag.String("reuse", "console", "comma separated list of kernel params value to reuse from current kernel (default to console)")
	appendCmdline     = flag.String("append", "", "Additional kernel params")

	devices []*diskboot.Device
)

func getDevice() (*diskboot.Device, error) {
	devices = diskboot.FindDevices(*devGlob)
	if len(devices) == 0 {
		return nil, errors.New("No devices found")
	}

	verbose("Got devices: %#v", devices)
	var err error
	deviceIndex := 0
	if len(devices) > 1 {
		if *sDeviceIndex == "" {
			for i, device := range devices {
				log.Printf("Device #%v: path: %v type: %v",
					i, device.DevPath, device.Fstype)
			}
			return nil, errors.New("Multiple devices found - must specify a device index")
		}
		if deviceIndex, err = strconv.Atoi(*sDeviceIndex); err != nil ||
			deviceIndex < 0 || deviceIndex >= len(devices) {
			return nil, fmt.Errorf("invalid device index %q", *sDeviceIndex)
		}
	}
	return devices[deviceIndex], nil
}

func getConfig(device *diskboot.Device) (*diskboot.Config, error) {
	configs := device.Configs
	if len(configs) == 0 {
		return nil, errors.New("No config found")
	}

	verbose("Got configs: %#v", configs)
	var err error
	configIndex := 0
	if len(configs) > 1 {
		if *sConfigIndex == "" {
			for i, config := range configs {
				log.Printf("Config #%v: path: %v", i, config.ConfigPath)
			}
			return nil, errors.New("Multiple configs found - must specify a config index")
		}
		if configIndex, err = strconv.Atoi(*sConfigIndex); err != nil ||
			configIndex < 0 || configIndex >= len(configs) {
			return nil, fmt.Errorf("invalid config index %q", *sConfigIndex)
		}
	}
	return configs[configIndex], nil
}

func getEntry(config *diskboot.Config) (*diskboot.Entry, error) {
	verbose("Got entries: %#v", config.Entries)
	var err error
	entryIndex := 0
	if *sEntryIndex != "" {
		if entryIndex, err = strconv.Atoi(*sEntryIndex); err != nil ||
			entryIndex < 0 || entryIndex >= len(config.Entries) {
			return nil, fmt.Errorf("invalid entry index %q", *sEntryIndex)
		}
	} else if config.DefaultEntry >= 0 {
		entryIndex = config.DefaultEntry
	} else {
		for i, entry := range config.Entries {
			log.Printf("Entry #%v: %#v", i, entry)
		}
		return nil, errors.New("No entry specified")
	}
	return &config.Entries[entryIndex], nil
}

func bootEntry(config *diskboot.Config, entry *diskboot.Entry) error {
	verbose("Booting entry: %v", entry)
	filter := cmdline.NewUpdateFilter(*appendCmdline, strings.Split(*removeCmdlineItem, ","), strings.Split(*reuseCmdlineItem, ","))
	err := entry.KexecLoad(config.MountPath, filter, *dryrun)
	if err != nil {
		return fmt.Errorf("wrror doing kexec load: %v", err)
	}

	if *dryrun {
		return nil
	}

	err = kexec.Reboot()
	if err != nil {
		return fmt.Errorf("error doing kexec reboot: %v", err)
	}
	return nil
}

func cleanDevices() {
	for _, device := range devices {
		if err := mount.Unmount(device.MountPath, true, false); err != nil {
			log.Printf("Error unmounting device %v: %v", device.DevPath, err)
		}
	}
}

func main() {
	flag.Parse()
	if *v {
		verbose = log.Printf
	}
	defer cleanDevices()

	device, err := getDevice()
	if err != nil {
		log.Panic(err)
	}
	config, err := getConfig(device)
	if err != nil {
		log.Panic(err)
	}
	entry, err := getEntry(config)
	if err != nil {
		log.Panic(err)
	}
	if err := bootEntry(config, entry); err != nil {
		log.Panic(err)
	}
}
