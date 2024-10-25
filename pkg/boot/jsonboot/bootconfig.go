// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonboot

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/crypto"
)

// BootConfig is a general-purpose boot configuration. It draws some
// characteristics from FIT but it's not compatible with it. It uses
// JSON for interoperability.
type BootConfig struct {
	Name          string   `json:"name,omitempty"`
	Kernel        string   `json:"kernel"`
	Initramfs     string   `json:"initramfs,omitempty"`
	KernelArgs    string   `json:"kernel_args,omitempty"`
	DeviceTree    string   `json:"devicetree,omitempty"`
	Multiboot     string   `json:"multiboot_kernel,omitempty"`
	MultibootArgs string   `json:"multiboot_args,omitempty"`
	Modules       []string `json:"multiboot_modules,omitempty"`
}

// IsValid returns true if a BootConfig object has valid content, and false
// otherwise
func (bc *BootConfig) IsValid() bool {
	return (bc.Kernel != "" && bc.Multiboot == "") || (bc.Kernel == "" && bc.Multiboot != "")
}

// ID retrurns an identifyer composed of bc's name and crc32 hash of bc.
// The ID is suitable to be used as part of a filepath.
func (bc *BootConfig) ID() string {
	id := strings.Title(strings.ToLower(bc.Name))
	id = strings.ReplaceAll(id, " ", "")
	id = strings.ReplaceAll(id, "/", "")
	id = strings.ReplaceAll(id, "\\", "")

	buf := []byte(filepath.Base(bc.Kernel))
	buf = append(buf, []byte(bc.KernelArgs)...)
	buf = append(buf, []byte(filepath.Base(bc.Initramfs))...)
	buf = append(buf, []byte(filepath.Base(bc.DeviceTree))...)
	buf = append(buf, []byte(filepath.Base(bc.Multiboot))...)
	buf = append(buf, []byte(bc.MultibootArgs)...)
	for _, mod := range bc.Modules {
		buf = append(buf, []byte(filepath.Base(mod))...)
	}
	h := crc32.ChecksumIEEE(buf)
	x := fmt.Sprintf("%x", h)

	return "BC_" + id + x
}

// FileNames returns a slice of all filenames in the bootconfig.
func (bc *BootConfig) FileNames() []string {
	var files []string
	if bc.Kernel != "" {
		files = append(files, bc.Kernel)
	}
	if bc.Initramfs != "" {
		files = append(files, bc.Initramfs)
	}
	if bc.DeviceTree != "" {
		files = append(files, bc.DeviceTree)
	}
	if bc.Multiboot != "" {
		files = append(files, bc.Multiboot)
	}
	for _, mod := range bc.Modules {
		if mod != "" {
			files = append(files, mod)
		}
	}
	return files
}

// ChangeFilePaths modifies the filepaths inside BootConfig. It replaces
// the current paths with new path leaving the last element of the path
// unchanged.
func (bc *BootConfig) ChangeFilePaths(newPath string) {
	if bc.Kernel != "" {
		bc.Kernel = filepath.Join(newPath, filepath.Base(bc.Kernel))
	}
	if bc.Initramfs != "" {
		bc.Initramfs = filepath.Join(newPath, filepath.Base(bc.Initramfs))
	}
	if bc.DeviceTree != "" {
		bc.DeviceTree = filepath.Join(newPath, filepath.Base(bc.DeviceTree))
	}
	if bc.Multiboot != "" {
		bc.Multiboot = filepath.Join(newPath, filepath.Base(bc.Multiboot))
	}
	for j, mod := range bc.Modules {
		if mod != "" {
			bc.Modules[j] = filepath.Join(newPath, filepath.Base(mod))
		}
	}
}

// SetFilePathsPrefix modifies the filepaths inside BootConfig. It appends
// prefix at the beginning of the current paths
func (bc *BootConfig) SetFilePathsPrefix(prefix string) {
	if bc.Kernel != "" {
		bc.Kernel = filepath.Join(prefix, bc.Kernel)
	}
	if bc.Initramfs != "" {
		bc.Initramfs = filepath.Join(prefix, bc.Initramfs)
	}
	if bc.DeviceTree != "" {
		bc.DeviceTree = filepath.Join(prefix, bc.DeviceTree)
	}
	if bc.Multiboot != "" {
		bc.Multiboot = filepath.Join(prefix, bc.Multiboot)
	}
	for j, mod := range bc.Modules {
		if mod != "" {
			bc.Modules[j] = filepath.Join(prefix, mod)
		}
	}
}

func (bc *BootConfig) bytestream() []byte {
	b := bc.Name + bc.Kernel + bc.Initramfs + bc.KernelArgs + bc.DeviceTree + bc.Multiboot + bc.MultibootArgs
	for _, module := range bc.Modules {
		b = b + module
	}
	return []byte(b)
}

// Boot tries to boot the kernel with optional initramfs and command line
// options. If a device-tree is specified, that will be used too
func (bc *BootConfig) Boot() error {
	crypto.TryMeasureData(crypto.BootConfigPCR, bc.bytestream(), "bootconfig")
	crypto.TryMeasureFiles(bc.FileNames()...)
	if bc.Kernel != "" {
		kernel, err := os.Open(bc.Kernel)
		if err != nil {
			return fmt.Errorf("can't open kernel file for measurement: %w", err)
		}
		var initramfs *os.File
		if bc.Initramfs != "" {
			initramfs, err = os.Open(bc.Initramfs)
			if err != nil {
				return fmt.Errorf("can't open initramfs file for measurement: %w", err)
			}
		}
		defer func() {
			// clean up
			if kernel != nil {
				if err := kernel.Close(); err != nil {
					log.Printf("Error closing kernel file descriptor: %v", err)
				}
			}
			if initramfs != nil {
				if err := initramfs.Close(); err != nil {
					log.Printf("Error closing initramfs file descriptor: %v", err)
				}
			}
		}()
		// Decompress Kernel (if compressed)
		kernelRaw, err := boot.CopyToFileIfNotRegular(util.TryGzipFilter(kernel), true)
		if err != nil {
			return err
		}

		if err := kexec.FileLoad(kernelRaw, initramfs, bc.KernelArgs); err != nil {
			return fmt.Errorf("kexec.FileLoad() failed: %w", err)
		}
	} else if bc.Multiboot != "" {
		mbkernel, err := os.Open(bc.Multiboot)
		if err != nil {
			log.Printf("Error opening multiboot kernel file: %v", err)
			return err
		}
		defer mbkernel.Close()

		// check multiboot header
		if err := multiboot.Probe(mbkernel); err != nil {
			log.Printf("Error parsing multiboot header: %v", err)
			return err
		}
		modules, err := multiboot.OpenModules(bc.Modules)
		if err != nil {
			return err
		}
		defer modules.Close()
		if err := multiboot.Load(true, mbkernel, bc.MultibootArgs, modules, nil); err != nil {
			return fmt.Errorf("kexec.Load() error: %w", err)
		}
	}
	err := kexec.Reboot()
	if err == nil {
		return errors.New("unexpectedly returned from Reboot() without error: system did not reboot")
	}
	return err
}

// NewBootConfig parses a boot configuration in JSON format and returns a
// BootConfig object.
func NewBootConfig(data []byte) (*BootConfig, error) {
	var bootconfig BootConfig
	if err := json.Unmarshal(data, &bootconfig); err != nil {
		return nil, err
	}
	return &bootconfig, nil
}
