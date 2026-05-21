// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iso

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

func fakeDev(t *testing.T) (*block.BlockDev, string) {
	tmp := t.TempDir()

	old := blockPath
	blockPath = filepath.Join(tmp, "block")
	t.Cleanup(func() { blockPath = old })

	sda := filepath.Join(blockPath, "sda")
	if err := os.MkdirAll(sda, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sda, "removable"), []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dev := &block.BlockDev{
		Name:   "sda",
		FsUUID: "1234-abcd",
	}

	mountDir := filepath.Join(tmp, "mnt")
	if err := os.MkdirAll(mountDir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dev, mountDir
}

// Read json configs, test against testdata dirs.
func TestConfigs(t *testing.T) {
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		name := filepath.Base(configPath)
		t.Run(configPath, func(t *testing.T) {
			want, err := os.ReadFile(test)
			if err != nil {
				t.Fatalf("Failed to read test json '%v': %v", test, err)
			}

			dev, mountDir := fakeDev(t)

			// Create a dummy iso file so WalkDir finds it
			isoPath := filepath.Join(mountDir, name+".iso")
			if err := os.WriteFile(isoPath, []byte("dummy"), 0o644); err != nil {
				t.Fatal(err)
			}

			// Mock mountISO to point to testdata/<name>
			old := mountISO
			mountISO = func(path string, pool *mount.Pool) (string, error) {
				absTestData, err := filepath.Abs("testdata")
				if err != nil {
					return "", err
				}
				return filepath.Join(absTestData, name), nil
			}
			t.Cleanup(func() { mountISO = old })

			l := ulogtest.Logger{TB: t}
			imgs, err := ParseISOFiles(l, mountDir, dev, &mount.Pool{})
			if err != nil {
				t.Fatalf("ParseISOFiles failed: %v", err)
			}

			if err := boottest.CompareImagesToJSON(imgs, want); err != nil {
				t.Error(err)
			}
		})
	}
}

// Read testdata dirs, generate json configs.
// Only run this if adding new testdata, via:
// GENERATE_CONFIGS=1 go test ./pkg/boot/iso -run TestGenerateConfigs
func TestGenerateConfigs(t *testing.T) {
	if os.Getenv("GENERATE_CONFIGS") != "1" {
		t.Skip("Skipping config generation. Set GENERATE_CONFIGS=1 to enable.")
	}

	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		name := f.Name()

		t.Run(name, func(t *testing.T) {
			dev, mountDir := fakeDev(t)

			isoPath := filepath.Join(mountDir, name+".iso")
			os.WriteFile(isoPath, []byte("dummy"), 0o644)

			old := mountISO
			mountISO = func(path string, pool *mount.Pool) (string, error) {
				absTestData, _ := filepath.Abs("testdata")
				return filepath.Join(absTestData, name), nil
			}
			t.Cleanup(func() { mountISO = old })

			l := ulogtest.Logger{TB: t}
			imgs, err := ParseISOFiles(l, mountDir, dev, &mount.Pool{})
			if err != nil {
				t.Fatalf("failed to parse %s: %v", f, err)
			}

			jsonPath := filepath.Join("testdata", name+".json")
			if err := boottest.ToJSONFile(imgs, jsonPath); err != nil {
				t.Errorf("failed to generate json: %v", err)
			}
		})
	}
}
