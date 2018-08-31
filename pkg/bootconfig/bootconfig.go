package bootconfig

import "encoding/json"

// BootConfig is a general-purpose boot configuration. It draws some
// characteristics from FIT but it's not compatible with it. It uses
// JSON for interoperability.
// If you add or remove fields, remember to update UnmarshalJSON.
type BootConfig struct {
	Name       string `json:"name,omitempty"`
	Kernel     string `json:"kernel"`
	Initramfs  string `json:"initramfs,omitempty"`
	KernelArgs string `json:"kernel_args,omitempty"`
	DeviceTree string `json:"devicetree,omitempty"`
}

// Validate returns true if a BootConfig object has valid content, and false
// otherwise
func (bc *BootConfig) Validate() bool {
	if bc.Kernel == "" {
		return false
	}
	return true
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
