// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hsskey

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

func TestToOctalEscapeSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "case 1: printable ASCII characters",
			input:    []byte("hello world!"),
			expected: "hello world!",
		},
		{
			name:     "case 2: non-printable ASCII characters",
			input:    []byte{0, 12, 200, 255},
			expected: "\\000\\014\\310\\377",
		},
		{
			name:     "case 3: mix printable and non-printable ASCII characters",
			input:    []byte{53, 193, 224, 43, 147, 134, 239, 118, 50, 115},
			expected: "5\\301\\340+\\223\\206\\357v2s",
		},
		{
			name:     "case 4: multiple characters with shorthand notations",
			input:    []byte{251, 5, 232, 226, 13, 9},
			expected: "\\373\\005\\350\\342\\r\\t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toOctalEscapeSequence(tt.input)
			if result != tt.expected {
				t.Fatalf("toOctalEscapeSequence(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Create dummy hss data, last 4 bytes are reserved for checksum
func createDummyHssData() []byte {
	data := make([]byte, hostSecretSeedStructSize-hostSecretSeedChecksumBytes)
	rand.Read(data)
	checksum := crc32.Checksum(data, castagnoliTable)

	checksumBytes := make([]byte, hostSecretSeedChecksumBytes)
	binary.BigEndian.PutUint32(checksumBytes, checksum)
	data = append(data, checksumBytes...)

	return data
}

func TestValidateChecksum(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		expect bool
	}{
		{
			name:   "valid checksum",
			expect: true,
			data:   createDummyHssData(),
		},
		{
			name:   "invalid checksum",
			expect: false,
			data: func() []byte {
				data := createDummyHssData()
				// Change the last byte of checksum to make it invalid
				data[len(data)-1]++
				return data
			}(),
		},
		{
			name:   "invalid size",
			expect: false,
			data:   []byte{1, 2, 3},
		},
		{
			name:   "zero bytes",
			expect: false,
			data:   make([]byte, 64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateChecksum(tt.data)
			if result != tt.expect {
				t.Fatalf("validateChecksum(%v) = %v, want %v", tt.data, result, tt.expect)
			}
		})
	}
}

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]uint8
		expected [][]uint8
	}{
		{
			name: "No duplicates",
			input: [][]uint8{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			expected: [][]uint8{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name: "With duplicates",
			input: [][]uint8{
				{1, 2, 3},
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: [][]uint8{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
		{
			name: "All duplicates",
			input: [][]uint8{
				{1, 2, 3},
				{1, 2, 3},
				{1, 2, 3},
			},
			expected: [][]uint8{
				{1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicate(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("deduplicate(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetHssEepromPaths(t *testing.T) {
	tests := []struct {
		name          string
		pattern       string
		allFiles      []string
		expectedFiles []string
		isErr         bool
	}{
		{
			name:    "Matching all EEPROMs",
			pattern: "*",
			allFiles: []string{
				"/sys/bus/i2c/devices/0-0050/eeprom",
				"/sys/bus/i2c/devices/1-0050/eeprom",
				"/sys/bus/i2c/devices/1-0050/device",
			},
			expectedFiles: []string{
				"/sys/bus/i2c/devices/0-0050/eeprom",
				"/sys/bus/i2c/devices/1-0050/eeprom",
			},
			isErr: false,
		},
		{
			name:    "Matching multiple EEPROMs",
			pattern: "0-007[01]",
			allFiles: []string{
				"/sys/bus/i2c/devices/0-0070/eeprom",
				"/sys/bus/i2c/devices/0-0071/eeprom",
				"/sys/bus/i2c/devices/1-0070/device",
			},
			expectedFiles: []string{
				"/sys/bus/i2c/devices/0-0070/eeprom",
				"/sys/bus/i2c/devices/0-0071/eeprom",
			},
			isErr: false,
		},
		{
			name:    "Matching multiple EEPROMs behind mux",
			pattern: "0-0050/channel-[01]/*",
			allFiles: []string{
				"/sys/bus/i2c/devices/0-0050/channel-0/0-0070/eeprom",
				"/sys/bus/i2c/devices/0-0050/channel-1/0-0075/eeprom",
				"/sys/bus/i2c/devices/0-0055/eeprom",
			},
			expectedFiles: []string{
				"/sys/bus/i2c/devices/0-0050/channel-0/0-0070/eeprom",
				"/sys/bus/i2c/devices/0-0050/channel-1/0-0075/eeprom",
			},
			isErr: false,
		},
		{
			name:    "No matching found",
			pattern: "0-0050",
			allFiles: []string{
				"/sys/bus/i2c/devices/0-0070/eeprom",
				"/sys/bus/foo/devices/skm_eeprom0/foo",
			},
			expectedFiles: []string{},
			isErr:         true,
		},
		{
			name:    "Matching other sysfs",
			pattern: "*",
			allFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom0/foo",
				"/sys/bus/foo/devices/skm_eeprom1/foo",
				"/sys/bus/foo/devices/skm_eeprom1/power",
				"/sys/bus/foo/devices/hoth_mailbox0/foo",
			},
			expectedFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom0/foo",
				"/sys/bus/foo/devices/skm_eeprom1/foo",
				"/sys/bus/foo/devices/hoth_mailbox0/foo",
			},
			isErr: false,
		},
		{
			name:    "Matching SKM FOOs",
			pattern: "skm_eeprom*",
			allFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom0/foo",
				"/sys/bus/foo/devices/skm_eeprom1/foo",
				"/sys/bus/foo/devices/skm_eeprom1/power",
				"/sys/bus/foo/devices/hoth_mailbox0/foo",
			},
			expectedFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom0/foo",
				"/sys/bus/foo/devices/skm_eeprom1/foo",
			},
			isErr: false,
		},
		{
			name:    "Matching One SKM FOO",
			pattern: "skm_eeprom1*",
			allFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom0/foo",
				"/sys/bus/foo/devices/skm_eeprom1/foo",
				"/sys/bus/foo/devices/skm_eeprom1/power",
				"/sys/bus/foo/devices/hoth_mailbox0/foo",
			},
			expectedFiles: []string{
				"/sys/bus/foo/devices/skm_eeprom1/foo",
			},
			isErr: false,
		},
		{
			name:    "Matching both i2c and foo",
			pattern: "*0",
			allFiles: []string{
				"/sys/bus/i2c/devices/0-0070/eeprom",
				"/sys/bus/foo/devices/skm_eeprom0/foo",
			},
			expectedFiles: []string{
				"/sys/bus/i2c/devices/0-0070/eeprom",
				"/sys/bus/foo/devices/skm_eeprom0/foo",
			},
			isErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp dir
			tempDir := t.TempDir()
			// Create all files
			for _, path := range tt.allFiles {
				fullPath := filepath.Join(tempDir, path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
				if err != nil {
					t.Fatal(err)
				}
				_, err = os.Create(fullPath)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Use temp dir to create base sysfs pattern
			testBaseI2cSysfsPattern := filepath.Join(tempDir, BaseSysfsPattern)
			testBaseNvmemSysfsPattern := filepath.Join(tempDir, "/sys/bus/foo/devices/%s/foo")
			result, err := getHssEepromPaths([]string{testBaseI2cSysfsPattern, testBaseNvmemSysfsPattern}, tt.pattern)

			// Check if error matches expectation
			if err != nil && !tt.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err == nil && tt.isErr {
				t.Fatalf("Expected error but got none")
			}

			// Check if result matches expectation
			expectedPaths := make(map[string]struct{})
			for _, path := range tt.expectedFiles {
				expectedPaths[filepath.Join(tempDir, path)] = struct{}{}
			}

			for _, path := range result {
				if _, ok := expectedPaths[path]; !ok {
					t.Fatalf("Unexpected path: %v", path)
				}
				delete(expectedPaths, path)
			}
			if len(expectedPaths) > 0 {
				t.Fatalf("Missing path: %v", expectedPaths)
			}
		})
	}
}

func TestReadHssFromFile(t *testing.T) {
	validHssBlock1 := createDummyHssData()
	validHssBlock2 := createDummyHssData()
	invalidChecksumHssBlock1 := createDummyHssData()
	invalidChecksumHssBlock2 := createDummyHssData()
	// make checksum invalid
	invalidChecksumHssBlock1[len(invalidChecksumHssBlock1)-1]++
	invalidChecksumHssBlock2[len(invalidChecksumHssBlock2)-1]++

	tests := []struct {
		name         string
		fileData     [][]byte
		expectedData [][]byte
		expectErr    bool
	}{
		{
			name:         "Four valid hss blocks",
			fileData:     [][]byte{validHssBlock1, validHssBlock1, validHssBlock2, validHssBlock2},
			expectedData: [][]byte{validHssBlock1[:32], validHssBlock1[:32], validHssBlock2[:32], validHssBlock2[:32]},
			expectErr:    false,
		},
		{
			name:         "Two valid and two invalid",
			fileData:     [][]byte{validHssBlock1, invalidChecksumHssBlock1, validHssBlock2, invalidChecksumHssBlock2},
			expectedData: [][]byte{validHssBlock1[:32], validHssBlock2[:32]},
			expectErr:    false,
		},
		{
			name:         "Mix invalid block",
			fileData:     [][]byte{make([]byte, 5), validHssBlock1, validHssBlock1, validHssBlock1, validHssBlock1},
			expectedData: nil,
			expectErr:    true,
		},
		{
			name:         "File length is longer than 4 blocks",
			fileData:     [][]byte{validHssBlock1, invalidChecksumHssBlock1, validHssBlock2, validHssBlock2, validHssBlock2},
			expectedData: [][]byte{validHssBlock1[:32], validHssBlock2[:32], validHssBlock2[:32]},
			expectErr:    false,
		},
		{
			name:         "File too short",
			fileData:     [][]byte{validHssBlock1},
			expectedData: nil,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file with the test HSS data.
			tempFile, err := os.CreateTemp("", "testData")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			var fileContent []byte
			for _, bytes := range tt.fileData {
				fileContent = append(fileContent, bytes...)
			}
			if err := os.WriteFile(tempFile.Name(), fileContent, 0o644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			hssList, err := ReadHssFromFile(tempFile.Name(), hostSecretSeedCount)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else if !reflect.DeepEqual(hssList, tt.expectedData) {
				t.Fatalf("ReadHssFromFile(%v) = %v, want %v", tempFile.Name(), hssList, tt.expectedData)
			}
		})
	}
}

func TestGetHssFromFile(t *testing.T) {
	// Create temp files.
	tempFile1, err := os.CreateTemp("", "testData")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile1.Name())

	tempFile2, err := os.CreateTemp("", "testData")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile2.Name())

	tempConfigFile, err := os.CreateTemp("", "testConfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempConfigFile.Name())

	// Write test data into temp file 1
	// 2 same key + another 2 same key
	var fileContent1 []byte
	f1td1 := createDummyHssData()
	fileContent1 = append(fileContent1, f1td1...)
	fileContent1 = append(fileContent1, f1td1...)
	f1td2 := createDummyHssData()
	fileContent1 = append(fileContent1, f1td2...)
	fileContent1 = append(fileContent1, f1td2...)
	if err := os.WriteFile(tempFile1.Name(), fileContent1, 0o644); err != nil {
		t.Fatal(err)
	}

	// Write test data into temp file 2
	// 2 same key + 1 invalid key + 1 another key
	var fileContent2 []byte
	f2td1 := createDummyHssData()
	fileContent2 = append(fileContent2, f2td1...)
	fileContent2 = append(fileContent2, f2td1...)
	f2td2 := createDummyHssData()
	f2td2[len(f2td2)-1]++
	fileContent2 = append(fileContent2, f2td2...)
	f2td3 := createDummyHssData()
	fileContent2 = append(fileContent2, f2td3...)
	if err := os.WriteFile(tempFile2.Name(), fileContent2, 0o644); err != nil {
		t.Fatal(err)
	}

	hssList, err := GetHssFromFile(os.Stdout, true, []string{tempFile1.Name(), tempFile2.Name()}, hostSecretSeedCount)
	// Expected to read two unique HSS key from td1 and td4
	expectedResult := [][]byte{f1td1[:32], f1td2[:32], f2td1[:32], f2td3[:32]}
	if err != nil {
		t.Fatalf("getHssFromFile() error = %v", err)
	}
	if !reflect.DeepEqual(hssList, expectedResult) {
		t.Fatalf("getHssFromFile() = %v, want %v", hssList, expectedResult)
	}
}

func TestWriteHssToFile(t *testing.T) {
	validHssBlock1 := createDummyHssData()
	validHssBlock2 := createDummyHssData()
	invalidChecksumHssBlock1 := createDummyHssData()
	invalidChecksumHssBlock2 := createDummyHssData()
	// make checksum invalid
	invalidChecksumHssBlock1[len(invalidChecksumHssBlock1)-1]++
	invalidChecksumHssBlock2[len(invalidChecksumHssBlock2)-1]++

	tests := []struct {
		name         string
		writeData    [][]byte
		expectedData [][]byte
		expectErr    bool
	}{
		{
			name:         "Four valid hss blocks",
			writeData:    [][]byte{validHssBlock1, validHssBlock1, validHssBlock2, validHssBlock2},
			expectedData: [][]byte{validHssBlock1, validHssBlock1, validHssBlock2, validHssBlock2},
			expectErr:    false,
		},
		{
			name:         "Two valid and two invalid",
			writeData:    [][]byte{validHssBlock1, invalidChecksumHssBlock1, validHssBlock2, invalidChecksumHssBlock2},
			expectedData: [][]byte{validHssBlock1, validHssBlock2},
			expectErr:    false,
		},
		{
			name:         "Mix invalid block",
			writeData:    [][]byte{make([]byte, 5), validHssBlock1, validHssBlock1, validHssBlock1},
			expectedData: [][]byte{validHssBlock1, validHssBlock1, validHssBlock1},
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file with the test HSS data.
			tempFile, err := os.CreateTemp("", "testData")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			err = WriteHssToFile(os.Stdout, true, tempFile, tt.writeData)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			var wantFileContent []byte
			for _, bytes := range tt.expectedData {
				wantFileContent = append(wantFileContent, bytes...)
			}

			gotFileContent, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			if !reflect.DeepEqual(gotFileContent, wantFileContent) {
				t.Fatalf("WriteHssToFile(%v, %v) =\ngot: %v\nwant: %v", tempFile.Name(), tt.writeData, gotFileContent, wantFileContent)
			}
		})
	}
}

func TestWriteHssToTempFile(t *testing.T) {
	validHssBlock1 := createDummyHssData()
	validHssBlock2 := createDummyHssData()
	invalidChecksumHssBlock1 := createDummyHssData()
	invalidChecksumHssBlock2 := createDummyHssData()
	// make checksum invalid
	invalidChecksumHssBlock1[len(invalidChecksumHssBlock1)-1]++
	invalidChecksumHssBlock2[len(invalidChecksumHssBlock2)-1]++

	tests := []struct {
		name          string
		writeData     [][]byte
		expectedData  [][]byte
		expectErr     bool
		warnings      io.Writer
		dangerVerbose bool
	}{
		{
			name:          "Four valid hss blocks",
			writeData:     [][]byte{validHssBlock1, validHssBlock1, validHssBlock2, validHssBlock2},
			expectedData:  [][]byte{validHssBlock1, validHssBlock1, validHssBlock2, validHssBlock2},
			expectErr:     false,
			warnings:      os.Stdout,
			dangerVerbose: true,
		},
		{
			name:          "Two valid and two invalid",
			writeData:     [][]byte{validHssBlock1, invalidChecksumHssBlock1, validHssBlock2, invalidChecksumHssBlock2},
			expectedData:  [][]byte{validHssBlock1, validHssBlock2},
			expectErr:     false,
			warnings:      os.Stdout,
			dangerVerbose: true,
		},
		{
			name:          "Mix invalid block",
			writeData:     [][]byte{make([]byte, 5), validHssBlock1, validHssBlock1, validHssBlock1},
			expectedData:  [][]byte{validHssBlock1, validHssBlock1, validHssBlock1},
			expectErr:     true,
			warnings:      os.Stdout,
			dangerVerbose: true,
		},
		{
			name:          "Two valid and two invalid skipping warnings",
			writeData:     [][]byte{validHssBlock1, invalidChecksumHssBlock1, validHssBlock2, invalidChecksumHssBlock2},
			expectedData:  [][]byte{validHssBlock1, validHssBlock2},
			expectErr:     false,
			warnings:      nil,
			dangerVerbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file with the test HSS data.
			tempFilePath, err := WriteHssToTempFile(os.Stdout, true, tt.writeData)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			var tempFile *os.File
			if tempFilePath != "" {
				defer os.Remove(tempFilePath)
				tempFile, err = os.Open(tempFilePath)
				if err != nil {
					t.Fatalf("Failed to open temporary file: %v", err)
				}
				defer tempFile.Close()
			}

			var wantFileContent []byte
			for _, bytes := range tt.expectedData {
				wantFileContent = append(wantFileContent, bytes...)
			}

			gotFileContent, err := os.ReadFile(tempFilePath)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			if !reflect.DeepEqual(gotFileContent, wantFileContent) {
				t.Fatalf("WriteHssToFile(%v, %v) =\ngot: %v\nwant: %v", tempFilePath, tt.writeData, gotFileContent, wantFileContent)
			}
		})
	}
}
