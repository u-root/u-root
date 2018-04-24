package main

import (
	"os"

	"github.com/u-root/u-root/pkg/kexec"
)

// BootConfig holds information to boot a kernel using kexec
type BootConfig struct {
	Kernel *os.File
	// KernelName is used only for reference but has no effect on booting
	KernelName string
	Initrd     *os.File
	// InitrdName is used only for reference but has no effect on booting
	InitrdName string
	Cmdline    string
}

func (bc BootConfig) IsValid() bool {
	return bc.Kernel != nil && bc.Initrd != nil
}

func (bc BootConfig) Boot() error {
	if err := kexec.FileLoad(bc.Kernel, bc.Initrd, bc.Cmdline); err != nil {
		return err
	}
	kexec.Reboot()
	// this should be never reached
	return nil
}

func (bc *BootConfig) Close() {
	if bc.Kernel != nil {
		bc.Kernel.Close()
		bc.Kernel = nil
		bc.KernelName = ""
	}
	if bc.Initrd != nil {
		bc.Initrd.Close()
		bc.Initrd = nil
		bc.InitrdName = ""
	}
}
