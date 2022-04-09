// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The disk_unlock command is used to unlock a disk drive as follows:
// 1. Via BMC, read a 32-byte secret seed known as the Host Secret Seed (HSS)
//    using the OpenBMC IPMI blob transfer protocol
// 2. Compute a password as follows:
//	We get the deterministically computed 32-byte HDKF-SHA256 using:
//	- salt: "SKM PROD_V2 ACCESS"
//	- hss: 32-byte HSS
//	- device identity: strings formed by concatenating the assembly serial
//	  number, the _ character, and the assembly part number.
// 3. Unlock the drive with the given password
// 4. Update the partition table for the disk
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/ipmi/blobs"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/mount/scuzz"
	"golang.org/x/crypto/hkdf"
)

const (
	hostSecretSeedLen = 32

	passwordSalt = "SKM PROD_V2 ACCESS"

	// Master Password ID for SKM-based unlock.
	skmMPI = 0x0601
)

var (
	disk               = flag.String("disk", "/dev/sda", "The disk to be unlocked")
	verbose            = flag.Bool("d", false, "print debug output")
	verboseNoSanitize  = flag.Bool("dangerously-disable-sanitize", false, "Print sensitive information - this should only be used for testing!")
	noRereadPartitions = flag.Bool("no-reread-partitions", false, "Only attempt to unlock the disk, don't re-read the partition table.")
	retries            = flag.Int("num_retries", 1, "Number of times to retry password if unlocking fails for any reason other than the password being wrong.")
)

func verboseLog(msg string) {
	if *verbose {
		log.Print(msg)
	}
}

// readHssBlob reads a host secret seed from the given blob id.
func readHssBlob(id string, h *blobs.BlobHandler) (data []uint8, rerr error) {
	sessionID, err := h.BlobOpen(id, blobs.BMC_BLOB_OPEN_FLAG_READ)
	if err != nil {
		return nil, fmt.Errorf("IPMI BlobOpen for %s failed: %v", id, err)
	}
	defer func() {
		// If the function returned successfully but failed to close the blob,
		// return an error.
		if err := h.BlobClose(sessionID); err != nil && rerr == nil {
			rerr = fmt.Errorf("IPMI BlobClose %s failed: %v", id, err)
		}
	}()

	data, err = h.BlobRead(sessionID, 0, hostSecretSeedLen)
	if err != nil {
		return nil, fmt.Errorf("IPMI BlobRead %s failed: %v", id, err)
	}

	if len(data) != hostSecretSeedLen {
		return nil, fmt.Errorf("HSS size incorrect: got %d for %s", len(data), id)
	}

	return data, nil
}

// getAllHss reads all host secret seeds over IPMI.
func getAllHss() ([][]uint8, error) {
	i, err := ipmi.Open(0)
	if err != nil {
		return nil, err
	}
	h := blobs.NewBlobHandler(i)

	blobCount, err := h.BlobGetCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get blob count: %v", err)
	}

	hssList := [][]uint8{}
	seen := make(map[string]bool)
	skmSubstr := "/skm/hss/"

	// Read from all */skm/hss/* blobs.
	for j := 0; j < blobCount; j++ {
		id, err := h.BlobEnumerate(j)
		if err != nil {
			return nil, fmt.Errorf("failed to enumerate blob %d: %v", j, err)
		}

		if !strings.Contains(id, skmSubstr) {
			continue
		}

		hss, err := readHssBlob(id, h)
		if err != nil {
			log.Printf("failed to read HSS of id %s: %v", id, err)
			continue
		}

		msg := fmt.Sprintf("HSS Entry: Id=%s", id)
		if *verboseNoSanitize {
			msg = msg + fmt.Sprintf(",Seed=%x", hss)
		}
		verboseLog(msg)

		hssStr := fmt.Sprint(hss)
		if _, ok := seen[hssStr]; !ok {
			seen[hssStr] = true
			hssList = append(hssList, hss)
		}
	}

	return hssList, nil
}

// Compute the password deterministically as the 32-byte HDKF-SHA256 of the
// HSS plus the device identity.
func genPassword(hss []byte, info *scuzz.Info) ([]byte, error) {
	hash := sha256.New
	devID := fmt.Sprintf("%s_%s", info.Serial, info.Model)

	r := hkdf.New(hash, hss, ([]byte)(passwordSalt), ([]byte)(devID))
	key := make([]byte, 32)

	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}
	return key, nil
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

	// Obtain 32 byte Host Secret Seed (HSS) from IPMI.
	hssList, err := getAllHss()
	if err != nil {
		log.Fatalf("error getting HSS: %v", err)
	}

	if len(hssList) == 0 {
		log.Fatalf("no HSS found - can't unlock disk.")
	}

	verboseLog(fmt.Sprintf("Found %d Host Secret Seeds.", len(hssList)))

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
	case info.MasterPasswordRev != skmMPI:
		log.Fatalf("Disk is locked with unknown master password ID: %X (Do you have skm tools installed?)", info.MasterPasswordRev)
	}

	// Try using each HSS to unlock the disk - only 1 should work.
	unlocked := false

TryAllHSS:
	for i, hss := range hssList {
		key, err := genPassword(hss, info)
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
