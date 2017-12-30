package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/diskboot"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/mount"
)

var (
	v       = flag.Bool("v", false, "Print debug messages")
	verbose = func(string, ...interface{}) {}
	dryrun  = flag.Bool("dryrun", false, "Only print out kexec commands")

	devGlob       = flag.String("dev", "/sys/class/block/*", "Device glob")
	sDeviceIndex  = flag.String("d", "", "Device index")
	sConfigIndex  = flag.String("c", "", "Config index")
	sEntryIndex   = flag.String("n", "", "Entry index")
	appendCmdline = flag.String("append", "", "Additional kernel params")

	devices []*diskboot.Device
)

func getDevice() *diskboot.Device {
	devices = diskboot.FindDevices(*devGlob)
	if len(devices) == 0 {
		log.Fatal("No devices found")
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
			log.Fatal("Multiple devices found - must specify a device index")
		}
		if deviceIndex, err = strconv.Atoi(*sDeviceIndex); err != nil ||
			deviceIndex < 0 || deviceIndex >= len(devices) {
			log.Fatal("Invalid device index:", *sDeviceIndex)
		}
	}
	return devices[deviceIndex]
}

func getConfig(device *diskboot.Device) *diskboot.Config {
	configs := device.Configs
	if len(configs) == 0 {
		log.Fatal("No config found")
	}

	verbose("Got configs: %#v", configs)
	var err error
	configIndex := 0
	if len(configs) > 1 {
		if *sConfigIndex == "" {
			for i, config := range configs {
				log.Printf("Config #%v: path: %v", i, config.ConfigPath)
			}
			log.Fatal("Multiple configs found - must specify a config index")
		}
		if configIndex, err = strconv.Atoi(*sConfigIndex); err != nil ||
			configIndex < 0 || configIndex >= len(configs) {
			log.Fatal("Invalid config index:", *sConfigIndex)
		}
	}
	return configs[configIndex]
}

func getEntry(config *diskboot.Config) *diskboot.Entry {
	verbose("Got entries: %#v", config.Entries)
	var err error
	entryIndex := 0
	if *sEntryIndex != "" {
		if entryIndex, err = strconv.Atoi(*sEntryIndex); err != nil ||
			entryIndex < 0 || entryIndex >= len(config.Entries) {
			log.Fatal("Invalid entry index:", *sEntryIndex)
		}
	} else if config.DefaultEntry >= 0 {
		entryIndex = config.DefaultEntry
	} else {
		for i, entry := range config.Entries {
			log.Printf("Entry #%v: %#v", i, entry)
		}
		log.Fatal("No entry specified")
	}
	return &config.Entries[entryIndex]
}

func bootEntry(config *diskboot.Config, entry *diskboot.Entry) {
	verbose("Booting entry: %v", entry)
	err := entry.KexecLoad(config.MountPath, *appendCmdline, *dryrun)
	if err != nil {
		log.Fatal("Error doing kexec load:", err)
	}

	if *dryrun {
		return
	}

	err = kexec.Reboot()
	if err != nil {
		log.Fatal("Error doing kexec reboot:", err)
	}
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

	device := getDevice()
	config := getConfig(device)
	entry := getEntry(config)
	bootEntry(config, entry)
}
