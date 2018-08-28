package bootconfig

import (
	"encoding/json"
	"fmt"
)

// Manifest is a list of BootConfig objects. The goal is to provide  multiple
// configurations to choose from.
type Manifest struct {
	// Version is a positive integer that determines the version of the Manifest
	// structure. This will be used when introducing breaking changes in the
	// Manifest interface
	Version int          `json:"version"`
	Configs []BootConfig `json:"configs"`
}

// NewManifest parses a manifest configuration, i.e. a list of boot
// configurations, in JSON format and returns a Manifest object.
func NewManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (mc *Manifest) GetBootConfig(name string) (*BootConfig, error) {
	for _, config := range mc.Configs {
		if config.Name == name {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("No boot config found by name: %s", name)
}
