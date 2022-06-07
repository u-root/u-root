// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

// fakeDevices returns a list of fake block devices and a pool of mount points.
// The pool is pre-populated so that Mount is never called.
func fakeDevices() (block.BlockDevices, *mount.Pool, error) {
	// For some reason, Glob("testdata_new/*/") does not work here.
	files, err := os.ReadDir("testdata_new")
	if err != nil {
		return nil, nil, err
	}
	dirs := []string{}
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, filepath.Join("testdata_new", f.Name()))
		}
	}

	devices := block.BlockDevices{}
	mountPool := &mount.Pool{}
	for _, dir := range dirs {
		// TODO: Also add LABEL to BlockDev
		fsUUID, _ := os.ReadFile(filepath.Join(dir, "UUID"))
		devices = append(devices, &block.BlockDev{
			Name:   dir,
			FSType: "test",
			FsUUID: strings.TrimSpace(string(fsUUID)),
		})
		mountPool.Add(&mount.MountPoint{
			Path:   dir,
			Device: filepath.Join("/dev", dir),
			FSType: "test",
		})
	}
	return devices, mountPool, nil
}

// Enable this to generate new configs.
func DISABLEDTestGenerateConfigs(t *testing.T) {
	tests, err := filepath.Glob("testdata_new/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			devices, mountPool, err := fakeDevices()
			if err != nil {
				t.Fatal(err)
			}
			imgs, err := ParseLocalConfig(context.Background(), configPath, devices, mountPool)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.ToJSONFile(imgs, test); err != nil {
				t.Errorf("failed to generate file: %v", err)
			}
		})
	}
}

func TestConfigs(t *testing.T) {
	// find all saved configs
	tests, err := filepath.Glob("testdata_new/*.json")
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

			devices, mountPool, err := fakeDevices()
			if err != nil {
				t.Fatal(err)
			}
			imgs, err := ParseLocalConfig(context.Background(), configPath, devices, mountPool)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.CompareImagesToJSON(imgs, want); err != nil {
				t.Errorf("ParseLocalConfig(): %v", err)
			}
		})
	}
}
