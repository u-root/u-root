// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The disk_unlock command is used to unlock a disk drive with a
// HSS-derived password, and rescan the drive to enumerate the
// unlocked partitions.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/hsskey"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/mount/scuzz"
)

const (
	// Master Password ID for SKM-based unlock.
	skmMPI = 0x0601
)

var (
	disk               = flag.String("disk", "/dev/sda", "The disk to be unlocked")
	verbose            = flag.Bool("d", false, "print debug output")
	verboseNoSanitize  = flag.Bool("dangerously-disable-sanitize", false, "Print sensitive information - this should only be used for testing!")
	noRereadPartitions = flag.Bool("no-reread-partitions", false, "Only attempt to unlock the disk, don't re-read the partition table.")
	retries            = flag.Int("num_retries", 1, "Number of times to retry password if unlocking fails for any reason other than the password being wrong.")
	salt               = flag.String("salt", hsskey.DefaultPasswordSalt, "Salt for password generation")
	eepromPattern      = flag.String("eeprom-pattern", "", "The pattern used to match EEPROM sysfs paths where the Host Secret Seeds are located")
	hssFiles           = flag.String("hss-files", "", "Comma deliminated list of files or directories containing additional Host Secret Seed (HSS)")
)

func verboseLog(msg string) {
	if *verbose {
		log.Print(msg)
	}
}

// writeFile is ioutil.WriteFile but disallows creating new file
func writeFile(filename string, contents string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	wlen, err := file.WriteString(contents)
	if err != nil && wlen < len(contents) {
		err = io.ErrShortWrite
	}
	// If Close() fails this likely indicates a write failure.
	if errClose := file.Close(); err == nil {
		err = errClose
	}
	return err
}

func main() {
	flag.Parse()

	// Open the disk. Read its identity, and use it to unlock the disk.
	sgdisk, err := scuzz.NewSGDisk(*disk)
	if err != nil {
		log.Fatalf("failed to open disk %v: %v", *disk, err)
	}

	info, err := sgdisk.Identify()
	if err != nil {
		log.Fatalf("failed to read disk %v identity: %v", *disk, err)
	}

	verboseLog(fmt.Sprintf("Disk info for %s: %s", *disk, info.String()))

	// Obtain 32 byte Host Secret Seed (HSS) from IPMI.
	hssList, err := hsskey.GetAllHss(os.Stdout, *verboseNoSanitize, *eepromPattern, *hssFiles)
	if err != nil {
		log.Fatalf("error getting HSS: %v", err)
	}

	if len(hssList) == 0 {
		log.Fatalf("no HSS found - can't unlock disk.")
	}

	verboseLog(fmt.Sprintf("Found %d Host Secret Seeds.", len(hssList)))

	switch {
	case !info.SecurityStatus.SecurityEnabled():
		log.Printf("Disk security is not enabled on %v.", *disk)
		return
	case info.SecurityStatus.SecurityFrozen():
		// If the disk is frozen, its security state cannot be changed until the next
		// power on or hardware reset. Disk unlock will fail anyways, so return.
		// This is unlikely to occur, since someone would need to freeze the drive's
		// security state between the last AC cycle and this code being run.
		log.Print("Disk security is frozen. Power cycle the machine to unfreeze the disk.")
		return
	case !info.SecurityStatus.SecurityLocked():
		log.Print("Disk is already unlocked.")
		return
	case info.SecurityStatus.SecurityCountExpired():
		// If the security count is expired, this means too many attempts have been
		// made to unlock the disk. Reset this with an AC cycle or hardware reset
		// on the disk.
		log.Fatalf("Security count expired on disk. Reset the password counter by power cycling the disk.")
	case info.MasterRevision != skmMPI:
		log.Fatalf("Disk is locked with unknown master password ID: %X (Do you have skm tools installed?)", info.MasterRevision)
	}

	// Try using each HSS to unlock the disk - only 1 should work.
	unlocked := false

TryAllHSS:
	for i, hss := range hssList {
		key, err := hsskey.GenPassword(hss, *salt, info.Serial, info.Model)
		if err != nil {
			log.Printf("Couldn't generate password with HSS %d: %v", i, err)
			continue
		}

		for r := 0; r < *retries; r++ {
			if err := sgdisk.Unlock(string(key), false); err != nil {
				log.Printf("Couldn't unlock disk with HSS %d: %v", i, err)
			} else {
				unlocked = true
				break TryAllHSS
			}
		}
	}

	if unlocked {
		log.Printf("Successfully unlocked disk %s.", *disk)
	} else {
		log.Fatalf("Failed to unlock disk %s with any HSS.", *disk)
	}

	if *noRereadPartitions {
		return
	}

	// Rescans all LUNs/Channels/Targets on the scsi_host. This ensures the kernel
	// updates the ATA driver to see the newly unlocked disk.
	verboseLog("Rescanning scsi...")
	if err := writeFile("/sys/class/scsi_host/host0/scan", "- - -"); err != nil {
		log.Fatalf("couldn't rescan SCSI to reload newly unlocked disk: %v", err)
	}

	// Update partitions on the on the disk.
	verboseLog("Reloading disk partitions...")
	diskdev, err := block.Device(*disk)
	if err != nil {
		log.Fatalf("Could not find %s: %v", *disk, err)
	}

	if err := diskdev.ReadPartitionTable(); err != nil {
		log.Fatalf("Could not re-read partition table: %v", err)
	}

	glob := filepath.Join("/sys/class/block", diskdev.Name+"*")
	parts, err := filepath.Glob(glob)
	if err != nil {
		log.Fatalf("Could not find disk partitions: %v", err)
	}

	verboseLog(fmt.Sprintf("Found these %s unlocked partitions: %v", *disk, parts))
}
