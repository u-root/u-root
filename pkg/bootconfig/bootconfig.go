package bootconfig

import (
	"encoding/json"
	"errors"
)

// BootConfig is a general-purpose boot configuration. It draws some
// characteristics from FIT but it's not compatible with it. It uses
// JSON for interoperability.
// If you add or remove fields, remember to update UnmarshalJSON.
type BootConfig struct {
	Name       string `json:"name"`
	Kernel     string `json:"kernel"`
	Initramfs  string `json:"initramfs,omitempty"`
	KernelArgs string `json:"kernel_args,omitempty"`
	DeviceTree string `json:"devicetree,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface to enforce the
// presence of required fields.
func (bc *BootConfig) UnmarshalJSON(data []byte) error {
	// Alias is needed to avoid an infinite loop of inherited and overridden
	// UnmarshalJSON
	type Alias BootConfig
	newBc := Alias{}
	err := json.Unmarshal(data, &newBc)
	if err != nil {
		return err
	}
	if newBc.Name == "" {
		return errors.New("Name field cannot be empty or missing")
	}
	if newBc.Kernel == "" {
		return errors.New("Kernel field cannot be empty or missing")
	}
	// either this, or reflection :(
	bc.Name = newBc.Name
	bc.Kernel = newBc.Kernel
	bc.Initramfs = newBc.Initramfs
	bc.KernelArgs = newBc.KernelArgs
	bc.DeviceTree = newBc.DeviceTree
	return nil
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
