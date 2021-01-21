// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"encoding/json"
	"fmt"

	"github.com/u-root/u-root/pkg/boot/jsonboot"
)

// Stconfig contains multiple u-root BootConfig stucts and additional
// information for stboot
type Stconfig struct {
	// configs is an array of u-root BootConfigs
	BootConfigs []jsonboot.BootConfig `json:"boot_configs"`
	// rootCertPath is the path to root certificate of the signing
	RootCertPath string `json:"root_cert"`
}

// StconfigFromBytes parses a Stcinfig from a byte slice
func stconfigFromBytes(data []byte) (*Stconfig, error) {
	var config Stconfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Bytes serializes a Stconfig stuct into a byte slice
func (cfg *Stconfig) bytes() ([]byte, error) {
	buf, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// IsValid returns true if all BootConfig structs inside the config has valid
// content.
func (cfg *Stconfig) IsValid() bool {
	if len(cfg.BootConfigs) == 0 {
		return false
	}
	for _, config := range cfg.BootConfigs {
		if !config.IsValid() {
			return false
		}
	}
	return cfg.RootCertPath != ""
}

// GetBootConfig returns the i-th boot configuration from the manifest, or an
// error if an invalid index is passed.
func (cfg *Stconfig) getBootConfig(index int) (*jsonboot.BootConfig, error) {
	if index < 0 || index >= len(cfg.BootConfigs) {
		return nil, fmt.Errorf("invalid index: not in range: %d", index)
	}
	return &cfg.BootConfigs[index], nil
}
