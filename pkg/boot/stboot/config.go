package stboot

import (
	"encoding/json"
	"fmt"

	"github.com/u-root/u-root/pkg/bootconfig"
)

// Stconfig contains multiple u-root BootConfig stucts and additional
// information for stboot
type Stconfig struct {
	// Configs is an array of u-root BootConfigs
	BootConfigs []bootconfig.BootConfig `json:"boot_configs"`
	// RootCertPath is the path to root certificate of the signing
	RootCertPath string `json:"root_cert"`
}

// StconfigFromBytes parses a manifest configuration, i.e. a list of boot
// configurations, in JSON format and returns a Manifest object.
func StconfigFromBytes(data []byte) (*Stconfig, error) {
	var config Stconfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// StconfigToBytes serializes a Stconfig stuct into a byte slice
func StconfigToBytes(cfg *Stconfig) ([]byte, error) {
	buf, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// IsValid returns true if all BootConfig structs inside the config has valid
// content.
func (cfg *Stconfig) IsValid() bool {
	for _, config := range cfg.BootConfigs {
		if !config.IsValid() {
			return false
		}
	}
	if cfg.RootCertPath == "" {
		return false
	}
	return true
}

// GetBootConfig returns the i-th boot configuration from the manifest, or an
// error if an invalid index is passed.
func (cfg *Stconfig) GetBootConfig(idx int) (*bootconfig.BootConfig, error) {
	if idx < 0 || idx >= len(cfg.BootConfigs) {
		return nil, fmt.Errorf("invalid index: not in range: %d", idx)
	}
	return &cfg.BootConfigs[idx], nil
}
