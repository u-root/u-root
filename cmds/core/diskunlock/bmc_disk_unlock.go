// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The bmc_disk_unlock command is used to unlock a disk drive as follows:
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
	"fmt"
	"io"
	"log"
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
)

var (
	diskName = "/dev/sda"
)

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
	skmPrefix := "/skm/hss/"

	// Read from all /skm/hss/* blobs.
	for j := 0; j < blobCount; j++ {
		id, err := h.BlobEnumerate(j)
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(id, skmPrefix) {
			continue
		}

		hss, err := readHssBlob(id, h)
		if err != nil {
			log.Printf("failed to read HSS of id %s: %v", id, err)
		} else {
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

func main() {
	// Obtain 32 byte Host Secret Seed (HSS) from IPMI.
	hssList, err := getAllHss()
	if err != nil {
		log.Fatalf("error getting HSS: %v", err)
	}

	if len(hssList) == 0 {
		log.Fatalf("no HSS found - can't unlock disk.")
	}

	log.Printf("Found %d Host Secret Seeds.", len(hssList))

	// Open the disk. Read its identity, and use it to unlock the disk.
	disk, err := scuzz.NewSGDisk(diskName)
	if err != nil {
		log.Fatalf("failed to open disk %v: %v", diskName, err)
	}

	info, err := disk.Identify()
	if err != nil {
		log.Fatalf("failed to read disk %v identity: %v", diskName, err)
	}

	log.Printf("Disk info for %s: %s", diskName, info.String())

	// Try using each HSS to unlock the disk - only 1 should work.
	unlocked := false
	for i, hss := range hssList {
		key, err := genPassword(hss, info)
		if err != nil {
			log.Printf("Couldn't generate password with HSS %d: %v", i, err)
			continue
		}

		if err := disk.Unlock((string)(key), false); err != nil {
			log.Printf("Couldn't unlock disk with HSS %d: %v", i, err)
		} else {
			unlocked = true
			break
		}
	}

	if unlocked {
		log.Printf("Successfully unlocked disk %s.", diskName)
	} else {
		log.Fatalf("Failed to unlock disk %s with any HSS.", diskName)
	}

	// Update partitions on the on the disk.
	diskdev, err := block.Device("/dev/sda")
	if err != nil {
		log.Fatalf("Could not find /dev/sda: %v", err)
	}

	if err := diskdev.ReadPartitionTable(); err != nil {
		log.Fatalf("Could not re-read partition table: %v", err)
	}

	parts, err := filepath.Glob("/sys/class/block/sda*")
	if err != nil {
		log.Fatalf("Could not find /sys/class/block/sda* files: %v", err)
	}

	log.Printf("Found these sda partitions: %v", parts)

}
