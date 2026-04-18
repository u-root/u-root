// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iso

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

func fakeDevices(t *testing.T, tmp string) {
	// Mock blockPath for isRemovable
	old := blockPath
	blockPath = filepath.Join(tmp, "block")
	t.Cleanup(func() { blockPath = old })

	// Simulate /sys/class/block/sda (removable)
	sda := filepath.Join(blockPath, "sda")
	if err := os.MkdirAll(sda, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sda, "removable"), []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Simulate /sys/class/block/sda1 (partition on removable)
	realSda1 := filepath.Join(sda, "sda1")
	if err := os.MkdirAll(realSda1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(realSda1, "partition"), []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Add the symlink
	if err := os.Symlink(realSda1, filepath.Join(blockPath, "sda1")); err != nil {
		t.Fatal(err)
	}

	// Simulate /sys/class/block/nvme0n1 (not removable)
	nvme := filepath.Join(blockPath, "nvme0n1")
	if err := os.MkdirAll(nvme, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nvme, "removable"), []byte("0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Simulate /sys/class/block/broken (no "removable" file)
	broken := filepath.Join(blockPath, "broken")
	if err := os.MkdirAll(broken, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestIsRemovable(t *testing.T) {
	tmp := t.TempDir()
	fakeDevices(t, tmp)

	for _, tt := range []struct {
		name      string
		devName   string
		want      bool
		wantError bool
	}{
		{
			name:    "removable-drive",
			devName: "sda",
			want:    true,
		},
		{
			name:    "removable-partition",
			devName: "sda1",
			want:    true,
		},
		{
			name:    "non-removable-drive",
			devName: "nvme0n1",
			want:    false,
		},
		{
			name:      "non-existent",
			devName:   "sdb",
			wantError: true,
		},
		{
			name:      "broken-drive",
			devName:   "broken",
			wantError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isRemovable(&block.BlockDev{Name: tt.devName})
			if (err != nil) != tt.wantError {
				t.Errorf("isRemovable() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("isRemovable() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ParseISOFiles - Failure cases
func TestParseISOFiles(t *testing.T) {
	tmp := t.TempDir()
	fakeDevices(t, tmp)

	mountDir := filepath.Join(tmp, "mnt")
	if err := os.MkdirAll(mountDir, 0o755); err != nil {
		t.Fatal(err)
	}

	noPermDir := filepath.Join(tmp, "restricted")
	if err := os.MkdirAll(noPermDir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(noPermDir, 0o755) })

	// Dummy iso
	isoPath := filepath.Join(mountDir, "test.iso")
	if err := os.WriteFile(isoPath, []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}

	isoMnt := filepath.Join(tmp, "iso-mnt")
	if err := os.MkdirAll(filepath.Join(isoMnt, "boot/grub"), 0o755); err != nil {
		t.Fatal(err)
	}
	grubCfg := `
menuentry "Generic Linux" {
	linux /vmlinuz quiet splash
	initrd /initrd
}
`
	if err := os.WriteFile(filepath.Join(isoMnt, "boot/grub/grub.cfg"), []byte(grubCfg), 0o644); err != nil {
		t.Fatal(err)
	}

	// Mock mountISO
	old := mountISO
	t.Cleanup(func() { mountISO = old })
	l := ulogtest.Logger{TB: t}

	for _, tt := range []struct {
		name      string
		devName   string
		mountDir  string
		isoMnt    string
		wantError bool
	}{
		{
			name:      "non-removable-drive",
			devName:   "nvme0n1",
			wantError: false,
		},
		{
			name:      "non-existant-drive",
			devName:   "sdb",
			wantError: true,
		},
		{
			name:      "non-existant-mountdir",
			devName:   "sda",
			mountDir:  filepath.Join(tmp, "no-mountdir"),
			wantError: true,
		},
		{
			name:      "non-existant-isomount",
			devName:   "sda",
			mountDir:  mountDir,
			isoMnt:    filepath.Join(tmp, "no-isomount"),
			wantError: false,
		},
		{
			name:      "permission-denied-mountdir",
			devName:   "sda",
			mountDir:  noPermDir,
			wantError: false,
		},
		{
			name:      "unparsable-iso",
			devName:   "sda",
			mountDir:  mountDir,
			isoMnt:    isoMnt,
			wantError: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mountISO = func(path string, pool *mount.Pool) (string, error) {
				return tt.isoMnt, nil
			}

			got, err := ParseISOFiles(l, tt.mountDir, &block.BlockDev{Name: tt.devName}, &mount.Pool{})
			if len(got) != 0 {
				t.Errorf("ParseISOFiles() = %v, want nil", got)
			}
			if (err != nil) != tt.wantError {
				t.Errorf("ParseISOFiles() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
