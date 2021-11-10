// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rng

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/recovery"
	"github.com/u-root/u-root/pkg/testutil"
)

func TestSetAvailableTRNG(t *testing.T) {
	for _, tt := range []struct {
		name    string
		setup   func(*testing.T)
		wantErr bool
	}{
		{
			name: "set trng",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name: "no rng available",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("none"); err != nil {
					t.Error(err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: true,
		},
		{
			name:    "no available file",
			setup:   func(t *testing.T) {},
			wantErr: true,
		},
		{
			name: "no write access",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err = f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(HwRandomCurrentFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: true,
		},
		{
			name: "no read access",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(HwRandomCurrentFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path, err := os.MkdirTemp(t.TempDir(), "rng-*")
			if err != nil {
				t.Errorf("Failed to create tmp dir: %v", err)
			}
			HwRandomAvailableFile = filepath.Join(path, "rng_available")
			HwRandomCurrentFile = filepath.Join(path, "rng_current")
			tt.setup(t)
			if tt.name == "no write access" {
				//TODO(MDr164): This doesn't work in qemu, no idea why
				testutil.SkipIfInVMTest(t)
				if err := os.Chmod(HwRandomCurrentFile, 0555); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
				if err := setAvailableTRNG(); err == nil {
					t.Error("Expected error, got nil")
				}
				if err := os.Chmod(HwRandomCurrentFile, 0755); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
			} else if tt.name == "no read access" {
				//TODO(MDr164): This doesn't work in qemu, no idea why
				testutil.SkipIfInVMTest(t)
				if err := os.Chmod(HwRandomCurrentFile, 0333); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
				if err := setAvailableTRNG(); err == nil {
					t.Error("Expected error, got nil")
				}
				if err := os.Chmod(HwRandomCurrentFile, 0755); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
			} else {
				if err := setAvailableTRNG(); err != nil && !tt.wantErr {
					t.Errorf("Expected nil, got %v", err)
				} else if err == nil && tt.wantErr {
					t.Error("Expected error, got nil")
				}
			}

			HwRandomAvailableFile = "/sys/class/misc/hw_random/rng_available"
			HwRandomCurrentFile = "/sys/class/misc/hw_random/rng_current"
		})
	}
}

func TestUpdateLinuxRandomness(t *testing.T) {
	for _, tt := range []struct {
		name    string
		setup   func(*testing.T)
		wantErr bool
	}{
		{
			name: "sufficient randomness",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(HwRandomDevice)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				b := [128]byte{}
				if _, err := f.Write(b[:]); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(RandomDevice)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(RandomEntropyAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("4000\n"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name:    "no trng",
			setup:   func(t *testing.T) {},
			wantErr: true,
		},
		{
			name: "no hwrng",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: true,
		},
		{
			name: "no random",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if _, err := f.WriteString("tpm-rng"); err != nil {
					t.Errorf("Failed to write to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				f, err = os.Create(HwRandomDevice)
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path, err := os.MkdirTemp(t.TempDir(), "rng-*")
			if err != nil {
				t.Error(err)
			}
			HwRandomAvailableFile = filepath.Join(path, "rng_available")
			HwRandomCurrentFile = filepath.Join(path, "rng_current")
			RandomDevice = filepath.Join(path, "random")
			HwRandomDevice = filepath.Join(path, "hwrng")
			RandomEntropyAvailableFile = filepath.Join(path, "entropy_avail")
			tt.setup(t)

			if !tt.wantErr {
				if err := UpdateLinuxRandomness(recovery.PermissiveRecoverer{}); err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
			} else {
				if err := UpdateLinuxRandomness(recovery.PermissiveRecoverer{}); err == nil {
					t.Error("Expected error, got nil")
				}
			}

			HwRandomAvailableFile = "/sys/class/misc/hw_random/rng_available"
			HwRandomCurrentFile = "/sys/class/misc/hw_random/rng_current"
			RandomDevice = "/dev/random"
			HwRandomDevice = "/dev/hwrng"
			RandomEntropyAvailableFile = "/proc/sys/kernel/random/entropy_avail"
		})
	}
}
