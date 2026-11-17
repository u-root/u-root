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

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/uio/uio"
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

func TestSearchSeparateFATBootPartition(t *testing.T) {
	for _, device := range []struct {
		name string
		uuid string
	}{
		{name: "nvme0n1", uuid: "a1e8-94da"},
		{name: "mmcblk0", uuid: "6cfc-948e"},
	} {
		t.Run(device.name, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				search string
			}{
				{name: "lowercase UUID", search: "search --fs-uuid --set=root " + device.uuid},
				{name: "uppercase UUID", search: "search --fs-uuid --set=root " + strings.ToUpper(device.uuid)},
				{name: "generated search", search: "search --no-floppy --fs-uuid --set=root " + strings.ToUpper(device.uuid)},
				{name: "generated search with hints", search: "search --no-floppy --fs-uuid --set=root --hint-bios=hd0,gpt2 --hint-efi=hd0,gpt2 --hint-baremetal=ahci0,gpt2 " + strings.ToUpper(device.uuid)},
				{name: "lowercase UUID with hints", search: "search --no-floppy --fs-uuid --set=root --hint-bios=hd0,gpt2 --hint-efi=hd0,gpt2 --hint-baremetal=ahci0,gpt2 " + device.uuid},
				{name: "file search with no-floppy", search: "search --no-floppy --file --set=root /vmlinuz-linux"},
			} {
				t.Run(tt.name, func(t *testing.T) {
					rootDir := filepath.Join(t.TempDir(), "p3")
					bootDir := filepath.Join(t.TempDir(), "p2")
					configDir := filepath.Join(rootDir, "boot", "grub")
					for _, dir := range []string{configDir, bootDir} {
						if err := os.MkdirAll(dir, 0o755); err != nil {
							t.Fatal(err)
						}
					}
					for name, contents := range map[string]string{
						"vmlinuz-linux":       "test kernel on p2",
						"initramfs-linux.img": "test initramfs on p2",
					} {
						if err := os.WriteFile(filepath.Join(bootDir, name), []byte(contents), 0o644); err != nil {
							t.Fatal(err)
						}
					}
					config := "menuentry 'test' {\nset root='hd0,gpt2'\n" + tt.search + "\nlinux /vmlinuz-linux\ninitrd /initramfs-linux.img\n}\n"
					if err := os.WriteFile(filepath.Join(configDir, "grub.cfg"), []byte(config), 0o644); err != nil {
						t.Fatal(err)
					}

					devices := block.BlockDevices{
						&block.BlockDev{Name: device.name + "p2", FsUUID: device.uuid, FSType: "vfat"},
						&block.BlockDev{Name: device.name + "p3", FsUUID: "test-root-uuid", FSType: "ext4"},
					}
					mountPool := &mount.Pool{}
					mountPool.Add(
						&mount.MountPoint{Path: bootDir, Device: devices[0].DevicePath(), FSType: "vfat"},
						&mount.MountPoint{Path: rootDir, Device: devices[1].DevicePath(), FSType: "ext4"},
					)
					images, err := ParseLocalConfig(context.Background(), rootDir, devices, mountPool)
					if err != nil {
						t.Fatal(err)
					}
					if len(images) != 1 {
						t.Fatalf("Got %d images, want 1", len(images))
					}
					image, ok := images[0].(*boot.LinuxImage)
					if !ok {
						t.Fatalf("Got image type %T, want *boot.LinuxImage", images[0])
					}
					for _, file := range []struct {
						name   string
						reader io.ReaderAt
						want   string
					}{
						{name: "kernel", reader: image.Kernel, want: "test kernel on p2"},
						{name: "initramfs", reader: image.Initrd, want: "test initramfs on p2"},
					} {
						got, err := uio.ReadAll(file.reader)
						if err != nil {
							t.Errorf("Failed to read %s from p2: %v", file.name, err)
							continue
						}
						if string(got) != file.want {
							t.Errorf("Read %s = %q, want %q", file.name, got, file.want)
						}
					}
				})
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
