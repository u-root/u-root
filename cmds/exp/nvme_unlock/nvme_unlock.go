// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The nvme_unlock command is used to unlock a NVMe drive with a
// HSS-derived password, and rescan the drive to enumerate the
// unlocked partitions.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/finddrive"
	"github.com/u-root/u-root/pkg/hsskey"
	"github.com/u-root/u-root/pkg/mount/block"
)

const (
	opalLockUnlockIoctl = 1092120797
)

type opalKey struct {
	lr     byte
	keyLen byte
	align  [6]byte
	key    [256]byte
}

type opalSessionInfo struct {
	sum     uint32
	who     uint32
	opalKey opalKey
}

type opalLockUnlock struct {
	session opalSessionInfo
	lState  uint32
	align   [4]byte
}

var (
	disk               = flag.String("disk", "", "The disk to be unlocked.  If left blank, will search for a boot disk.")
	verbose            = flag.Bool("d", false, "print debug output")
	verboseNoSanitize  = flag.Bool("dangerously-disable-sanitize", false, "Print sensitive information - this should only be used for testing!")
	noRereadPartitions = flag.Bool("no-reread-partitions", false, "Only attempt to unlock the disk, don't re-read the partition table.")
	lock               = flag.Bool("lock", false, "Lock instead of unlocking")
	salt               = flag.String("salt", hsskey.DefaultPasswordSalt, "Salt for password generation")
	eepromSysfwPath    = flag.String("eeprom-sysfs-path", "", "Additional path (relative to /sys/bus) used with eeprom-pattern to locate the Host Secret Seeds")
	eepromPattern      = flag.String("eeprom-pattern", "", "The pattern used to match EEPROM sysfs paths where the Host Secret Seeds are located")
	hssFiles           = flag.String("hss-files", "", "Comma deliminated list of files or directories containing additional Host Secret Seed (HSS)")
)

func verboseLog(msg string) {
	if *verbose {
		log.Print(msg)
	}
}

func getSysfsInfo(index string, field string) (string, error) {
	path := fmt.Sprintf("/sys/class/nvme/nvme%s/%s", index, field)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading sysfs info at path %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func run(disk string, verbose bool, verboseNoSanitize bool, noRereadPartitions bool, lock bool) error {
	if disk == "" {
		disks, err := finddrive.FindSlotType(finddrive.M2MKeySlotType)
		if err != nil {
			return fmt.Errorf("error finding boot disk: %w", err)
		}
		if len(disks) == 0 {
			return fmt.Errorf("no boot disk found")
		}
		disk = disks[0]
		if len(disks) > 1 {
			log.Printf("Multiple boot disk candidates found, using the first from the following list: %v", disks)
		} else if verbose {
			log.Printf("Found boot disk %s", disk)
		}
	}

	commandName := "unlock"
	if lock {
		commandName = "lock"
	}

	diskIDRegexp := regexp.MustCompile(`/dev/nvme(\d+)n.*`)
	diskIDMatches := diskIDRegexp.FindStringSubmatch(disk)
	if diskIDMatches == nil {
		return fmt.Errorf("unable to parse device path %s", disk)
	}
	diskID := diskIDMatches[1]

	serial, err := getSysfsInfo(diskID, "serial")
	if err != nil {
		return err
	}
	model, err := getSysfsInfo(diskID, "model")
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Serial %s", serial)
		log.Printf("Model %s", model)
	}

	diskFd, err := os.Open(disk)
	if err != nil {
		return fmt.Errorf("error opening disk: %w", err)
	}
	defer diskFd.Close()

	sysfsPaths := []string{hsskey.BaseSysfsPattern}
	if eepromSysfwPath != nil {
		sysfsPaths = append(sysfsPaths, fmt.Sprintf("/sys/bus/%s", *eepromSysfwPath))
	}

	hssList, err := hsskey.GetAllHssWithPaths(os.Stdout, verboseNoSanitize,
		sysfsPaths,
		*eepromPattern, *hssFiles)
	if err != nil {
		return fmt.Errorf("error getting HSS: %w", err)
	}

	if len(hssList) == 0 {
		return fmt.Errorf("no HSS found - can't unlock disk")
	}

	if verbose {
		log.Printf("Found %d Host Secret Seeds.", len(hssList))
	}

	succeeded := false
	for i, hss := range hssList {
		password, err := hsskey.GenPassword(hss, *salt, serial, model)
		if err != nil {
			log.Printf("Couldn't generate password with HSS %d: %v", i, err)
			continue
		}

		var state uint32 = 0x02 // OPAL_RW
		if lock {
			state = 0x04 // OPAL_LK
		}
		arg := opalLockUnlock{
			session: opalSessionInfo{
				sum: 0,
				who: 0,
				opalKey: opalKey{
					keyLen: 32,
				},
			},
			lState: state,
		}
		copy(arg.session.opalKey.key[0:32], password)

		r1, _, errNo := syscall.Syscall(syscall.SYS_IOCTL, diskFd.Fd(),
			uintptr(opalLockUnlockIoctl), uintptr(unsafe.Pointer(&arg)))
		if errNo != 0 {
			log.Printf("%s failed with errno: %v", commandName, errNo)
		} else if r1 != 0 {
			log.Printf("%s returned nonzero value %v, password may be incorrect", commandName, r1)
		} else {
			succeeded = true
			break
		}
	}

	if succeeded {
		log.Printf("Successfully %sed disk %s.", commandName, disk)
	} else {
		log.Printf("Failed to %s disk %s with any HSS.", commandName, disk)
		return fmt.Errorf("all HSS failed")
	}

	if noRereadPartitions {
		return nil
	}

	// Update partitions on the on the disk.
	if verbose {
		log.Print("Reloading disk partitions...")
	}
	diskdev, err := block.Device(disk)
	if err != nil {
		return fmt.Errorf("could not find %s: %w", disk, err)
	}

	if err := diskdev.ReadPartitionTable(); err != nil && !lock {
		return fmt.Errorf("could not re-read partition table: %w", err)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(*disk, *verbose, *verboseNoSanitize, *noRereadPartitions, *lock); err != nil {
		log.Fatalf("nvme_unlock: %v", err)
	}
}
