package bootconfig

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/systemboot/systemboot/pkg/crypto"
	"github.com/u-root/u-root/pkg/kexec"
)

// BootConfig is a general-purpose boot configuration. It draws some
// characteristics from FIT but it's not compatible with it. It uses
// JSON for interoperability.
type BootConfig struct {
	Name       string `json:"name,omitempty"`
	Kernel     string `json:"kernel"`
	Initramfs  string `json:"initramfs,omitempty"`
	KernelArgs string `json:"kernel_args,omitempty"`
	DeviceTree string `json:"devicetree,omitempty"`
}

// IsValid returns true if a BootConfig object has valid content, and false
// otherwise
func (bc *BootConfig) IsValid() bool {
	return bc.Kernel != ""
}

// Boot tries to boot the kernel with optional initramfs and command line
// options. If a device-tree is specified, that will be used too
func (bc *BootConfig) Boot() error {
	crypto.TryMeasureBootConfig(bc.Name, bc.Kernel, bc.Initramfs, bc.KernelArgs, bc.DeviceTree)

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

	err = kexec.Reboot()
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
