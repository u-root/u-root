package main

import (
	"os"

	"github.com/u-root/u-root/pkg/kexec"
)

// BootConfig holds information to boot a kernel using kexec
type BootConfig struct {
	Kernel    *os.File
	Initramfs *os.File
	Cmdline   string
}

// IsValid returns true if the BootConfig has a valid kernel and initrd entry
func (bc BootConfig) IsValid() bool {
	return bc.Kernel != nil && bc.Initramfs != nil
}

// Boot tries to boot the kernel pointed by the BootConfig option, or returns an
// error if it cannot be booted. The kernel is loaded using kexec
func (bc BootConfig) Boot() error {
	if err := kexec.FileLoad(bc.Kernel, bc.Initramfs, bc.Cmdline); err != nil {
		return err
	}
	kexec.Reboot()
	// this should be never reached
	return nil
}

// Close will close all the open file descriptor used for kernel and initrd
func (bc *BootConfig) Close() {
	if bc.Kernel != nil {
		bc.Kernel.Close()
		bc.Kernel = nil
	}
	if bc.Initramfs != nil {
		bc.Initramfs.Close()
		bc.Initramfs = nil
	}
}
