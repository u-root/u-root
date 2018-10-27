package bootconfig

import (
	"encoding/json"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/kexecbin"
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
	// kexec: try the kexecbin executable first, and if it fails, use the native
	// Go implementation of kexec from u-root
	log.Printf("Trying KexecBin on %+v", bc)
	if err := kexecbin.KexecBin(bc.Kernel, bc.KernelArgs, bc.Initramfs, bc.DeviceTree); err != nil {
		if os.IsNotExist(err) {
			log.Printf("BootConfig: KexecBin failed, trying pure-Go kexec. Error: %v", err)
		} else {
			return err
		}
	}
	kernel, err := os.Open(bc.Kernel)
	if err != nil {
		return err
	}
	var initramfs *os.File
	if bc.Initramfs != "" {
		initramfs, err = os.Open(bc.Kernel)
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
	return kexec.FileLoad(kernel, initramfs, bc.KernelArgs)
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
