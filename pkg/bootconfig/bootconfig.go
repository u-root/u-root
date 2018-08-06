package bootconfig

import "encoding/json"

// BootConfig is a general-purpose boot configuration. It draws some
// characteristics from FIT but it's not compatible with it. It uses
// JSON for interoperability.
type BootConfig struct {
	Name       string `json:"name"`
	Kernel     string `json:"kernel"`
	Initramfs  string `json:"initramfs"`
	KernelArgs string `json:"kernel_args"`
}

func NewBootConfig(data []byte) (*BootConfig, error) {
	var bootconfig BootConfig
	if err := json.Unmarshal(data, &bootconfig); err != nil {
		return nil, err
	}
	return &bootconfig, nil
}
