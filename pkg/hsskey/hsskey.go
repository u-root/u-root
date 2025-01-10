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
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("IPMI BlobOpen for %s failed: %w", id, err)
	}
	defer func() {
		// If the function returned successfully but failed to close the blob,
		// return an error.
		if err := h.BlobClose(sessionID); err != nil && rerr == nil {
			rerr = fmt.Errorf("IPMI BlobClose %s failed: %w", id, err)
		}
	}()

	data, err = h.BlobRead(sessionID, 0, hostSecretSeedLen)
	if err != nil {
		return nil, fmt.Errorf("IPMI BlobRead %s failed: %w", id, err)
	}

	if len(data) != hostSecretSeedLen {
		return nil, fmt.Errorf("HSS size incorrect: got %d for %s", len(data), id)
	}

	return data, nil
}

// getHssFromIpmi reads all host secret seeds over IPMI.
func getHssFromIpmi(warnings io.Writer, verboseDangerous bool) ([][]uint8, error) {
	i, err := ipmi.Open(0)
	if err != nil {
		return nil, err
	}
	h := blobs.NewBlobHandler(i)

	blobCount, err := h.BlobGetCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get blob count: %w", err)
	}

	hssLists := [][][]uint8{}
	hssPrefixes := []string{"/skm/hss/", "/skm/hss-backup/"}
	for range hssPrefixes {
		hssLists = append(hssLists, [][]uint8{})
	}

	// Gather all HSS entries, expected to look like /skm/hss/0, /skm/hss-backup/3, etc.
	for j := 0; j < blobCount; j++ {
		id, err := h.BlobEnumerate(j)
		if err != nil {
			return nil, fmt.Errorf("failed to enumerate blob %d: %w", j, err)
		}

		// Ignore entries with a trailing /, which don't actually represent a HSS
		if strings.HasSuffix(id, "/") {
			continue
		}

		prefixIdx := -1
		for k, prefix := range hssPrefixes {
			if strings.HasPrefix(id, prefix) {
				prefixIdx = k
				break
			}
		}
		if prefixIdx == -1 {
			continue
		}

		hss, err := readHssBlob(id, h)
		if err != nil {
			log.Printf("Failed to read HSS of id %s: %v", id, err)
			continue
		}

		if warnings != nil && verboseDangerous {
			fmt.Fprintf(warnings, "HSS Entry: Id=%s, Seed=%x\n", id, hss)
		}

		hssLists[prefixIdx] = append(hssLists[prefixIdx], hss)
	}

	// Deduplicate repeated HSS entries, and order by hssPrefixes
	seen := make(map[string]bool)
	hssList := [][]uint8{}
	for _, list := range hssLists {
		for _, hss := range list {
			hssStr := fmt.Sprint(hss)
			if !seen[hssStr] {
				seen[hssStr] = true
				hssList = append(hssList, hss)
			}
		}
	}

	return hssList, nil
}

// GetAllHss reads all host secret seeds from IPMI or EEPROM.
//   - eepromPattern: A string pattern to find EEPROMs in sysfs paths. The glob string used for
//     searching will be in the format: "/sys/bus/i2c/devices/{eepromPattern}/eeprom".
//     For example, 0-005*
//     An empty string "" will skip the attempt to read from EEPROM.
func GetAllHss(warnings io.Writer, verboseDangerous bool, eepromPattern string, hssFiles string) ([][]uint8, error) {
	return GetAllHssWithPaths(warnings, verboseDangerous, []string{BaseSysfsPattern}, eepromPattern, hssFiles)
}

// GetAllHssWithPaths reads all host secret seeds from IMPI or any given sysfs paths.
//   - sysFsPatterns: A list of sysfs path to search in.
//   - eepromPattern: A string pattern to find EEPROMs in the sysfs paths.
func GetAllHssWithPaths(warnings io.Writer, verboseDangerous bool, sysFsPatterns []string, eepromPattern string, hssFiles string) ([][]uint8, error) {
	// Attempt to get HSS from IPMI.
	hssList, err := getHssFromIpmi(warnings, verboseDangerous)
	if err != nil || len(hssList) == 0 {
		fmt.Fprintf(warnings, "Failed to get HSS key from IPMI: %v\n", err)
	}

	// Attempt to get HSS from EEPROM.
	if eepromPattern != "" {
		filePaths, err := getHssEepromPaths(sysFsPatterns, eepromPattern)
		if err != nil {
			fmt.Fprintf(warnings, "Failed to find HSS EEPROM paths: %v\n", err)
		}
		hssEeprom, err := GetHssFromFile(warnings, verboseDangerous, filePaths, hostSecretSeedCount)
		if err == nil && len(hssEeprom) > 0 {
			hssList = append(hssList, hssEeprom...)
		} else {
			fmt.Fprintf(warnings, "Failed to get HSS key from file: %v\n", err)
		}
	}

	// Attempt to get HSS from files.
	if hssFiles != "" {
		hssFileArr := []string{}
		// Parse if list of paths were given.
		for _, hssFile := range strings.Split(strings.TrimSpace(hssFiles), ",") {
			hssFile, err := evaluateHssPath(hssFile)
			if err != nil {
				fmt.Fprintf(warnings, "Error parsing path %s for file: %v\n", hssFile, err)
			}
			hssFileArr = append(hssFileArr, hssFile...)
		}

		if len(hssFileArr) > 0 {
			hss, err := GetHssFromFile(warnings, verboseDangerous, hssFileArr, 0)
			if err != nil {
				fmt.Fprintf(warnings, "Error parsing files %s for HSS: %v\n", hss, err)
			}
			hssList = append(hssList, hss...)
		}
	}

	if len(hssList) > 0 {
		return hssList, nil
	}
	return nil, fmt.Errorf("failed all HSS retrieval attempts")
}

// evaluateHssPath evaluates a filepath to return the file or contents if a directory.
// This function evaluates the base pointer for symlinks.
func evaluateHssPath(hssFiles string) ([]string, error) {
	// Follow symlink if needed.
	hssFiles, err := filepath.EvalSymlinks(hssFiles)
	if err != nil {
		return nil, fmt.Errorf("failed evaluating HSS path: %w", err)
	}

	hssFileArr := []string{}

	fi, err := os.Stat(hssFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to stat HSS files: %w", err)
	}

	// Check if hssFiles is a directory or a file. Add all regular files in the base directory.
	switch mode := fi.Mode(); {
	case mode.IsDir():
		filepath.Walk(hssFiles, func(fpath string, info os.FileInfo, _ error) error {
			if info.Mode().IsRegular() {
				hssFileArr = append(hssFileArr, fpath)
			}
			return nil
		})
	case mode.IsRegular():
		hssFileArr = append(hssFileArr, hssFiles)
	}

	return hssFileArr, nil
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
