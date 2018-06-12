package main

import (
	"flag"
	"log"
	"path"
	"syscall"

	"github.com/insomniacslk/systemboot/pkg/storage"
)

// TODO backward compatibility for BIOS mode with partition type 0xee
// TODO read and write non-volatile variables on the flash ROM via VPD
//      https://chromium.googlesource.com/chromiumos/platform/vpd/
//      via sysfs:
//      https://github.com/torvalds/linux/blob/master/drivers/firmware/google/vpd.c
// TODO use a proper parser for grub config (see grub.go)

var (
	baseMountPoint = flag.String("m", "/mnt", "Base mount point where to mount partiions")
	doDebug        = flag.Bool("d", false, "Print debug output")
)

var debug = func(string, ...interface{}) {}

func main() {
	flag.Parse()

	if *doDebug {
		debug = log.Printf
	}

	// Get all the available block devices
	devices, err := storage.GetBlockStats()
	if err != nil {
		log.Fatal(err)
	}
	// print partition info
	for _, dev := range devices {
		log.Printf("Device: %+v", dev)
		table, err := storage.GetGPTTable(dev)
		if err != nil {
			continue
		}
		log.Printf("  Table: %+v", table)
		for _, part := range table.Partitions {
			log.Printf("    Partition: %+v", part)
			if !part.IsEmpty() {
				log.Printf("      UUID: %s", part.Type.String())
			}
		}
	}

	// get a list of supported file systems for real devices (i.e. skip nodev)
	debug("Getting list of supported filesystems")
	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	debug("Supported file systems: %v", filesystems)

	// detect EFI system partitions
	// TODO currently, this is not necessary, but will be once we have VPD.
	debug("Searching for EFI system partitions")
	esps, err := storage.FilterEFISystemPartitions(devices)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d system partitions: %+v", len(esps), esps)

	// try mounting all the available devices, with all the supported file
	// systems
	debug("trying to mount all the available block devices with all the supported file system types")
	mounted := make([]storage.Mountpoint, 0)
	for _, dev := range devices {
		devname := path.Join("/dev", dev.Name)
		mountpath := path.Join(*baseMountPoint, dev.Name)
		if mountpoint, err := storage.Mount(devname, mountpath, filesystems); err != nil {
			debug("Failed to mount %s on %s: %v", devname, mountpath, err)
		} else {
			mounted = append(mounted, *mountpoint)
		}
	}
	debug("mounted: %+v", mounted)

	// search for a valid grub config and extracts the boot configuration
	bootconfigs := make([]BootConfig, 0)
	for _, mountpoint := range mounted {
		bootconfigs = append(bootconfigs, ScanGrubConfigs(mountpoint.Path)...)
	}
	log.Printf("Found %d boot configs", len(bootconfigs))
	for _, cfg := range bootconfigs {
		debug("%+v", cfg)
	}

	// try to kexec into every boot config kernel until one succeeds
	for _, cfg := range bootconfigs {
		log.Printf("trying to boot %s", cfg.KernelName)
		if err := cfg.Boot(); err != nil {
			log.Printf("Failed to boot kernel %s: %v", cfg.KernelName, err)
			cfg.Close()
		}
	}
	log.Print("No boot configuration succeeded")

	// if we are here, booting failed, so let's clean things up by closing all
	// open file descriptors and unmounting the devices that we mounted
	for _, cfg := range bootconfigs {
		cfg.Close()
	}
	for _, mountpoint := range mounted {
		syscall.Unmount(mountpoint.Path, syscall.MNT_DETACH)
	}
}
