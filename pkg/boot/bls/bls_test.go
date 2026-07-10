// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bls

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

var blsEntries = []struct {
	entry string
	err   string
}{
	{
		entry: "entry-1.conf",
	},
	{
		entry: "entry-2.conf",
		err:   "neither linux, efi, nor multiboot present in BootLoaderSpec config",
	},
}

// Enable this temporarily to generate new configs. Double-check them by hand.
func DISABLEDTestGenerateConfigs(t *testing.T) {
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			imgs, err := ScanBLSEntries(ulogtest.Logger{TB: t}, configPath, nil, "")
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.ToJSONFile(imgs, test); err != nil {
				t.Errorf("failed to generate file: %v", err)
			}
		})
	}
}

func TestParseBLSEntries(t *testing.T) {
	fsRoot := "./testdata/madeup"
	dir := filepath.Join(fsRoot, "loader/entries")

	for _, tt := range blsEntries {
		t.Run(tt.entry, func(t *testing.T) {
			image, err := parseBLSEntry(filepath.Join(dir, tt.entry), fsRoot, nil, false)
			if err != nil {
				if tt.err == "" {
					t.Fatalf("Got error %v", err)
				}
				if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("Got error %v, expected error to contain %s", err, tt.err)
				}
				return
			}
			if tt.err != "" {
				t.Fatalf("Expected error %s, got no error", tt.err)
			}
			t.Logf("Got image: %s", image.String())
		})
	}
}

func TestScanBLSEntries(t *testing.T) {
	// find all saved configs
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			want, err := os.ReadFile(test)
			if err != nil {
				t.Errorf("Failed to read test json '%v':%v", test, err)
			}

			imgs, err := ScanBLSEntries(ulogtest.Logger{TB: t}, configPath, nil, "")
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.CompareImagesToJSON(imgs, want); err != nil {
				t.Errorf("ParseLocalConfig(): %v", err)
			}
		})
	}
}

func TestSetBLSRank(t *testing.T) {
	fsRoot := "./testdata/madeup"
	dir := filepath.Join(fsRoot, "loader/entries")
	testRank := 2
	originRank := os.Getenv("BLS_BOOT_RANK")
	os.Setenv("BLS_BOOT_RANK", strconv.Itoa(testRank))

	for _, tt := range blsEntries {
		t.Run(tt.entry, func(t *testing.T) {
			image, err := parseBLSEntry(filepath.Join(dir, tt.entry), fsRoot, nil, false)
			if err != nil {
				if tt.err == "" {
					t.Fatalf("Got error %v", err)
				}
				if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("Got error %v, expected error to contain %s", err, tt.err)
				}
				return
			}
			if tt.err != "" {
				t.Fatalf("Expected error %s, got no error", tt.err)
			}

			if image.Rank() != testRank {
				t.Errorf("Expected rank %d, got %d", testRank, image.Rank())
			}
		})
	}

	os.Setenv("BLS_BOOT_RANK", originRank)
}

func TestParseConf(t *testing.T) {
	m, err := parseConf("testdata/typeone/fedora.conf")
	if err != nil {
		t.Fatal(err)
	}

	if m["title"] != "Fedora 19 (Rawhide)" {
		t.Errorf("expected title 'Fedora 19 (Rawhide)', got %v", m["title"])
	}
	if m["sort-key"] != "fedora" {
		t.Errorf("expected sort-key 'fedora', got %v", m["sort-key"])
	}
	if m["machine-id"] != "6a9857a393724b7a981ebb5b8495b9ea" {
		t.Errorf("expected machine-id '6a9857a393724b7a981ebb5b8495b9ea', got %v", m["sort-key"])
	}
	if m["version"] != "3.8.0-2.fc19.x86_64" {
		t.Errorf("expected version '3.8.0-2.fc19.x86_64', got %v", m["version"])
	}
	if m["options"] != "root=UUID=6d3376e4-fc93-4509-95ec-a21d68011da2 quiet" {
		t.Errorf("expected options 'root=UUID=6d3376e4-fc93-4509-95ec-a21d68011da2 quiet', got %v", m["options"])
	}
	if m["architecture"] != "x64" {
		t.Errorf("expected architecture 'fedora', got %v", m["architecture"])
	}
	if m["linux"] != "/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/linux" {
		t.Errorf("expected linux '/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/linux', got %v", m["linux"])
	}
	if m["initrd"] != "/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/initrd" {
		t.Errorf("expected initrd '/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/initrd', got %v", m["initrd"])
	}
	// TODO: extra entry may appear multiple times
	if m["extra"] != "/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/somedata.cred" {
		t.Errorf("expected extra '6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/somedata.cred', got %v", m["extra"])
	}
	if _, ok := m["key-without-value1"]; ok {
		t.Errorf("key-without-value should not be in map")
	}
	if _, ok := m["key-without-value2"]; ok {
		t.Errorf("key-without-value should not be in map")
	}
}
