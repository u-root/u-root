// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bootconfig"
)

// Stconfig describes the the configuration of the OS that is
// loaded by stboot.
type Stconfig struct {
	Label string `json:"label"`

	Kernel    string `json:"kernel"`
	Initramfs string `json:"initramfs"`
	Cmdline   string `json:"cmdline"`

	Tboot       string   `json:"tboot"`
	TbootArgs   string   `json:"tboot_args"`
	ACMs        []string `json:"acms"`
	AllowNonTXT bool     `json:"allow_non_txt"`
}

// StconfigFromFile parses a Stconfig from a json file
func StconfigFromFile(src string) (*Stconfig, error) {
	cfgBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	cfg, err := StconfigFromBytes(cfgBytes)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// StconfigFromBytes parses a Stconfig from a byte slice.
func StconfigFromBytes(data []byte) (*Stconfig, error) {
	var config Stconfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Write saves cfg to file named by stboot.ConfigName at a path named by dir.
func (cfg *Stconfig) Write(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	buf, err := cfg.Bytes()
	if err != nil {
		return err
	}
	dst := filepath.Join(dir, ConfigName)
	err = ioutil.WriteFile(dst, buf, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Bytes serializes a Stconfig stuct into a byte slice.
func (cfg *Stconfig) Bytes() ([]byte, error) {
	buf, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Validate returns true if cfg is valid to be booted.
func (cfg *Stconfig) Validate() error {
	if cfg.Kernel == "" {
		return errors.New("stconfig: missing kernel")
	}
	if !cfg.AllowNonTXT {
		if cfg.Tboot == "" {
			return errors.New("stconfig: TXT required but missing tboot binary")
		}
		if len(cfg.ACMs) == 0 {
			return errors.New("stconfig: TXT required but missing ACM")
		}
	}
	return nil
}

// GetBootConfig returns the i-th boot configuration from the manifest, or an
// error if an invalid index is passed.
func (cfg *Stconfig) GetBootConfig(index int) (*bootconfig.BootConfig, error) {
	// if index < 0 || index >= len(cfg.BootConfigs) {
	// 	return nil, fmt.Errorf("invalid index: not in range: %d", index)
	// }
	// return &cfg.BootConfigs[index], nil
	return nil, errors.New("GetBootConfig() is not implemented yet")
}
