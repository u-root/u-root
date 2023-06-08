// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hsskey

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
)

const (
	baseSysfsPattern = "/sys/bus/i2c/devices/%s/eeprom"
)

// toOctalEscapeSequence converts a byte slice into similar format as C++ protobuf DebugString()
// If the byte represents a printable ASCII character, add it as is to the output string
// If it is non-printable or non-ASCII, add it as octal notation
// This is used for comparing output from C++ tool for debugging purpose.
func toOctalEscapeSequence(b []byte) string {
	var shortHandNotions = map[string]string{
		"\011": "\\t",
		"\015": "\\r",
	}
	octalEscape := ""
	for _, v := range b {
		octalChar := fmt.Sprintf("\\%03o", v)
		if v >= 32 && v <= 126 {
			octalEscape += string(v)
		} else if notion, ok := shortHandNotions[string(v)]; ok {
			// DebugString prints these characters as shorthand notations
			octalEscape += notion
		} else {
			octalEscape += octalChar
		}
	}
	return octalEscape
}

// validateChecksum check if HSS data is valid based on the checksum value stored in the last 4
// bytes.
func validateChecksum(data []byte) bool {
	size := len(data)
	if size != hostSecretSeedStructSize {
		return false
	}
	// Split the data into the main part and the checksum. Last 4 bytes is used for storing
	// checksum value.
	mainPart := data[:size-hostSecretSeedChecksumBytes]
	checksumBytes := data[size-hostSecretSeedChecksumBytes:]

	// Convert the checksum bytes from big-endian.
	checksum := binary.BigEndian.Uint32(checksumBytes)

	// Compute the checksum of the main part.
	table := crc32.MakeTable(crc32.Castagnoli)
	computedChecksum := crc32.Checksum(mainPart, table)

	if computedChecksum != checksum {
		return false
	}
	return true
}

// deduplicate removes duplicate entries from a slice of byte slices.
func deduplicate(data [][]byte) [][]byte {
	seen := make(map[string]bool)
	result := [][]byte{}

	for _, entry := range data {
		strEntry := string(entry)
		if !seen[strEntry] {
			seen[strEntry] = true
			result = append(result, entry)
		}
	}
	return result
}

// getHssEepromPaths takes a glob pattern for EEPROM sysfs paths and return all the matching paths.
// This function will return error if the provided glob pattern doesn't match any existing path.
//
// For example, if busDevicePattern is "0-007[01]", this function will look for all the paths
// matching "/sys/bus/i2c/devices/0-007[01]/eeprom"
func getHssEepromPaths(basePattern string, busDevicePattern string) ([]string, error) {
	fullPathGlob := fmt.Sprintf(basePattern, busDevicePattern)
	matches, err := filepath.Glob(fullPathGlob)
	if err != nil {
		return nil, fmt.Errorf("failed to search the file path %v", err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matching path found for glob pattern %s", fullPathGlob)
	}
	return matches, nil
}

// readFromFile reads HSS keys from the specified file. Expecting the file containing 4 consecutive
// HSS keys. Each HSS key is 64 bytes long and has a checksum for validation.
func readHssFromFile(filePath string) ([][]byte, error) {
	validLen := hostSecretSeedStructSize * hostSecretSeedCount
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v", err)
	}
	if len(data) < validLen {
		return nil, fmt.Errorf("file size %d is less than expected", len(data))
	}
	if len(data) > validLen {
		log.Printf("Expecting %d bytes but got more %d bytes", validLen, len(data))
	}

	var hssList [][]byte
	for i := 0; i < hostSecretSeedCount; i++ {
		start := i * hostSecretSeedStructSize
		end := start + hostSecretSeedStructSize
		hssKey := data[start:end]

		// Verify the checksums.
		if !validateChecksum(hssKey) {
			log.Printf("Checksum validation failed for HSS key at offset %d", start)
			continue
		}
		hssList = append(hssList, hssKey[:hostSecretSeedLen])
	}
	return hssList, nil
}

func getHssFromFile(verbose bool, verboseDangerous bool, filePaths []string) ([][]byte, error) {
	allHss := [][]byte{}
	for _, f := range filePaths {
		hssKeys, err := readHssFromFile(f)
		if err != nil {
			log.Printf("Failed to read HSS keys from file %s. err: %v", f, err)
			continue
		}

		if verbose {
			msg := fmt.Sprintf("Reading HSS keys from file=%s", f)
			if verboseDangerous {
				for _, hss := range hssKeys {
					msg = msg + fmt.Sprintf("\nseed=%x, seed(octal escape sequence)=%s", hss,
						toOctalEscapeSequence(hss))
				}
			}
			log.Print(msg)
		}
		allHss = append(allHss, hssKeys...)
	}

	// Remove duplicate HSS key.
	allHss = deduplicate(allHss)

	return allHss, nil
}
