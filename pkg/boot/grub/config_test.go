// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
	"context"
	"io"
	"log"
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

func FuzzParseGrubConfig(f *testing.F) {
	baseDir := f.TempDir()

	dirPath := filepath.Join(baseDir, "EFI", "uefi")
	err := os.MkdirAll(dirPath, 0o777)
	if err != nil {
		f.Fatalf("failed %v: %v", dirPath, err)
	}

	path := filepath.Join(dirPath, "grub.cfg")
	devices := block.BlockDevices{&block.BlockDev{
		Name:   dirPath,
		FSType: "test",
		FsUUID: strings.TrimSpace("07338180-4a96-4611-aa6a-a452600e4cfe"),
	}}
	mountPool := &mount.Pool{}
	mountPool.Add(&mount.MountPoint{
		Path:   dirPath,
		Device: filepath.Join("/dev", dirPath),
		FSType: "test",
	})

	// no log output
	log.SetOutput(io.Discard)
	log.SetFlags(0)

	// get seed corpora from testdata_new files
	seeds, err := filepath.Glob("testdata_new/*/*/*/grub.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	seeds2, err := filepath.Glob("testdata_new/*/*/grub.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	seeds = append(seeds, seeds2...)
	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}

	f.Add([]byte("multiBoot 0\nmodule --nounzip"))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			return
		}

		// do not allow arbitrary files reads
		if bytes.Contains(data, []byte("include")) {
			return
		}

		err = os.WriteFile(path, data, 0o777)
		if err != nil {
			t.Fatalf("Failed to create configfile '%v':%v", path, err)
		}

		ParseLocalConfig(context.Background(), baseDir, devices, mountPool)
	})
}
