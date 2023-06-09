// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hsskey provides functionality for generating a key for unlocking
// drives based on the following procedure:
//  1. Via BMC, read a 32-byte secret seed known as the Host Secret Seed (HSS)
//     using the OpenBMC IPMI blob transfer protocol
//  2. Via EEPROM, read a 32-byte secret seed from EEPROM
//  3. Compute a password as follows:
//     We get the deterministically computed 32-byte HDKF-SHA256 using:
//     - salt: "SKM PROD_V2 ACCESS" (default)
//     - hss: 32-byte HSS
//     - device identity: strings formed by concatenating the assembly serial
//     number, the _ character, and the assembly part number.
package hsskey

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/ipmi/blobs"
	"golang.org/x/crypto/hkdf"
)

type blobReader interface {
	BlobOpen(id string, flags int16) (blobs.SessionID, error)
	BlobRead(sid blobs.SessionID, offset, size uint32) ([]uint8, error)
	BlobClose(sid blobs.SessionID) error
}

const (
	hostSecretSeedLen           = 32 // Size for HSS seed
	hostSecretSeedChecksumBytes = 4  // Size for HSS checksum
	hostSecretSeedStructSize    = 64 // Total size for serialized HSS data
	hostSecretSeedCount         = 4  // Number of HSS data stored per EEPROM

	DefaultPasswordSalt = "SKM PROD_V2 ACCESS"
)

// readHssBlob reads a host secret seed from the given blob id.
func readHssBlob(id string, h blobReader) (data []uint8, rerr error) {
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

// getHssFromIpmi reads all host secret seeds over IPMI.
func getHssFromIpmi(verbose bool, verboseDangerous bool) ([][]uint8, error) {
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
			log.Printf("Failed to read HSS of id %s: %v", id, err)
			continue
		}

		if verbose {
			msg := fmt.Sprintf("HSS Entry: Id=%s", id)
			if verboseDangerous {
				msg = msg + fmt.Sprintf(",Seed=%x", hss)
			}
			log.Print(msg)
		}

		hssStr := fmt.Sprint(hss)
		if !seen[hssStr] {
			seen[hssStr] = true
			hssList = append(hssList, hss)
		}
	}

	return hssList, nil
}

// GetAllHss reads all host secret seeds from IPMI or EEPROM.
//   - eepromPattern: A string pattern to find EEPROMs in sysfs paths. The glob string used for
//     searching will be in the format: "/sys/bus/i2c/devices/{eepromPattern}/eeprom".
//     For example, 0-005*
//     An empty string "" will skip the attempt to read from EEPROM.
func GetAllHss(verbose bool, verboseDangerous bool, eepromPattern string) ([][]uint8, error) {
	// Attempt to get HSS from IPMI.
	hssList, err := getHssFromIpmi(verbose, verboseDangerous)
	if err == nil && len(hssList) > 0 {
		// If it successfully reads from IPMI and the list is not empty, return the result.
		return hssList, nil
	}
	log.Printf("Failed to get HSS key from IPMI: %v", err)

	// Attempt to get HSS from EEPROM.
	if eepromPattern != "" {
		filePaths, err := getHssEepromPaths(baseSysfsPattern, eepromPattern)
		if len(filePaths) != 2 {
			log.Printf("Expecting 2 EEPROMs but found %d", len(filePaths))
		}
		if err != nil {
			return nil, fmt.Errorf("failed to find HSS EEPROM paths: %v", err)
		}
		hssList, err = getHssFromFile(verbose, verboseDangerous, filePaths)
		if err == nil && len(hssList) > 0 {
			return hssList, nil
		}
		log.Printf("Failed to get HSS key from file: %v", err)
	}

	return nil, fmt.Errorf("failed all HSS retrieval attempts")
}

// GenPassword computes the password deterministically as the 32-byte HDKF-SHA256 of the
// HSS plus the device identity.
func GenPassword(hss []byte, salt string, identifiers ...string) ([]byte, error) {
	hash := sha256.New
	devID := strings.Join(identifiers, "_")

	r := hkdf.New(hash, hss, ([]byte)(salt), ([]byte)(devID))
	key := make([]byte, 32)

	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}
	return key, nil
}
