// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/crypto"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot"
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

// FileNames returns a slice of all filenames in the bootconfig.
func (bc *BootConfig) fileNames() []string {
	str := make([]string, 0)
	str = append(str, bc.Kernel)
	str = append(str, bc.Initramfs)
	for _, module := range bc.Modules {
		str = append(str, module)
	}
	return str
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
	crypto.TryMeasureFiles(bc.fileNames()...)
	if bc.Kernel != "" {
		kernel, err := os.Open(bc.Kernel)
		if err != nil {
			return err
		}
		var initramfs *os.File
		if bc.Initramfs != "" {
			initramfs, err = os.Open(bc.Initramfs)
			if err != nil {
				return err
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
		if err := kexec.FileLoad(kernel, initramfs, bc.KernelArgs); err != nil {
			return err
		}
	} else if bc.Multiboot != "" {
		// check multiboot header
		if err := multiboot.Probe(bc.Multiboot); err != nil {
			log.Printf("Error parsing multiboot header: %v", err)
			return err
		}
		if err := multiboot.Load(true, bc.Multiboot, bc.MultibootArgs, bc.Modules, nil); err != nil {
			return fmt.Errorf("kexec.Load() error: %v", err)
		}
	}
	err := kexec.Reboot()
	if err == nil {
		return errors.New("Unexpectedly returned from Reboot() without error. The system did not reboot")
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
