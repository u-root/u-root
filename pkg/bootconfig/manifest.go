// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"encoding/json"
	"fmt"
)

// global variables
var (
	CurrentManifestVersion = 1
)

// Manifest is a list of BootConfig objects. The goal is to provide  multiple
// configurations to choose from.
type Manifest struct {
	// Version is a positive integer that determines the version of the Manifest
	// structure. This will be used when introducing breaking changes in the
	// Manifest interface.
	Version int          `json:"version"`
	Configs []BootConfig `json:"configs"`
}

// NewManifest returns a new empty Manifest structure with the current version
// field populated.
func NewManifest() *Manifest {
	return &Manifest{
		Version: CurrentManifestVersion,
	}
}

// ManifestFromBytes parses a manifest configuration, i.e. a list of boot
// configurations, in JSON format and returns a Manifest object.
func ManifestFromBytes(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

// GetBootConfig returns the i-th boot configuration from the manifest, or an
// error if an invalid index is passed.
func (mc *Manifest) GetBootConfig(idx int) (*BootConfig, error) {
	if idx < 0 || idx >= len(mc.Configs) {
		return nil, fmt.Errorf("Invalid index: not in range: %d", idx)
	}
	return &mc.Configs[idx], nil
}
