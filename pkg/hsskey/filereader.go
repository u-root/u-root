// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hsskey

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	BaseSysfsPattern = "/sys/bus/i2c/devices/%s/eeprom"
)

// toOctalEscapeSequence converts a byte slice into similar format as C++ protobuf DebugString()
// If the byte represents a printable ASCII character, add it as is to the output string
// If it is non-printable or non-ASCII, add it as octal notation
// This is used for comparing output from C++ tool for debugging purpose.
func toOctalEscapeSequence(b []byte) string {
	shortHandNotions := map[string]string{
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

	return computedChecksum == checksum
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

// getHssEepromPaths takes glob patterns for EEPROM sysfs paths and returns all the matching paths.
// This function will return error if the provided glob patterns don't match any existing path.
//
// For example, if busDevicePattern is "0-007[01]" and basePatterns is ["/sys/bus/i2c/devices/*/eeprom"]
// this function will look for all the paths matching "/sys/bus/i2c/devices/0-007[01]/eeprom"
func getHssEepromPaths(basePatterns []string, busDevicePattern string) ([]string, error) {
	matches := []string{}
	globs := []string{}
	for _, pattern := range basePatterns {
		fullPathGlob := fmt.Sprintf(pattern, busDevicePattern)
		m, _ := filepath.Glob(fullPathGlob)
		globs = append(globs, fullPathGlob)
		matches = append(matches, m...)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matching path found for glob pattern %s", strings.Join(globs, ","))
	}
	return matches, nil
}

// ReadHssFromFile reads HSS keys from the specified file.
// Each HSS key is 64 bytes long and has a checksum for validation.
func ReadHssFromFile(filePath string, minHssPerFile int) ([][]byte, error) {
	minValidLen := hostSecretSeedStructSize * minHssPerFile
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %w", err)
	}

	if len(data) < minValidLen {
		return nil, fmt.Errorf("file size %d is less than expected", len(data))
	}
	if len(data) > minValidLen {
		log.Printf("Expecting %d bytes but got more %d bytes", minValidLen, len(data))
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

// GetHssFromFile reads HSS keys from the specified files.
// Each HSS key is 64 bytes long and has a checksum for validation. Duplicate HSS are removed.
func GetHssFromFile(warnings io.Writer, verboseDangerous bool, filePaths []string, minHssPerFile int) ([][]byte, error) {
	allHss := [][]byte{}
	for _, f := range filePaths {
		hssKeys, err := ReadHssFromFile(f, minHssPerFile)
		if err != nil {
			log.Printf("Failed to read HSS keys from file %s. err: %v", f, err)
			continue
		}

		if verboseDangerous && warnings != nil {
			msg := fmt.Sprintf("Reading HSS keys from file=%s", f)
			for _, hss := range hssKeys {
				msg = msg + fmt.Sprintf("\nseed=%x, seed(octal escape sequence)=%s", hss,
					toOctalEscapeSequence(hss))
			}
			io.WriteString(warnings, msg+"\n")
		}
		allHss = append(allHss, hssKeys...)
	}

	// Remove duplicate HSS key.
	allHss = deduplicate(allHss)

	return allHss, nil
}

// WriteHssToTempFile writes a list of HSS to a tmpfs file where the filepath is returned.
// See WriteHssToFile for HSS details.
func WriteHssToTempFile(warnings io.Writer, verboseDangerous bool, hss [][]byte) (string, error) {
	file, err := os.CreateTemp("", "hss_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file %w", err)
	}

	err = WriteHssToFile(warnings, verboseDangerous, file, hss)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return "", err
	}

	file.Close()
	return file.Name(), nil
}

// WriteHssToFile writes a list of HSS to an open file.
// Each HSS key is expected to be 64 bytes long and has a checksum for validation.
// HSS that fail to validate are not written to the file.
// Function returns error if no HSS are written or writing can no longer continue.
func WriteHssToFile(warnings io.Writer, verboseDangerous bool, file *os.File, hss [][]byte) error {
	if len(hss) == 0 {
		return fmt.Errorf("InvalidArgument: No HSS to write")
	}
	if file == nil {
		return fmt.Errorf("InvalidArgument: file is nil")
	}
	// Count of successful writes
	written := 0

	for i := 0; i < len(hss); i++ {
		hssKey := hss[i]

		// Verify the checksums.
		if !validateChecksum(hssKey) {
			if warnings != nil {
				msg := fmt.Sprintf("Checksum validation failed for HSS key index %d", i)
				if verboseDangerous {
					msg = msg + fmt.Sprintf("\nseed=%x, seed(octal escape sequence)=%s", hssKey,
						toOctalEscapeSequence(hssKey))
				}
				io.WriteString(warnings, msg+"\n")
			}
			continue
		}

		if verboseDangerous && warnings != nil {
			msg := fmt.Sprintf("Writing HSS key to file=%s\nseed=%x, seed(octal escape sequence)=%s",
				file.Name(), hssKey, toOctalEscapeSequence(hssKey))
			io.WriteString(warnings, msg+"\n")
		}

		n, err := file.Write(hssKey)
		if n != len(hssKey) {
			return fmt.Errorf("failed to write entire HSS key to file %s", file.Name())
		}
		if err != nil {
			return err
		}
		written++
	}

	if written == 0 {
		return fmt.Errorf("no HSS keys were written to file %s", file.Name())
	}
	return nil
}
