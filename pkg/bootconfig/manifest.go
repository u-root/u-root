package bootconfig

import "encoding/json"

// ManifestConfig is a list of BootConfig objects. The goal is to provide
// multiple configurations to choose from.
type ManifestConfig struct {
	Configs []BootConfig `json:"configs"`
}

// NewManifestConfig parses a manifest configuration, i.e. a list of boot
// configurations, in JSON format and returns a ManifestConfig object.
func NewManifestConfig(data []byte) (*ManifestConfig, error) {
	var manifestConfig ManifestConfig
	if err := json.Unmarshal(data, &manifestConfig); err != nil {
		return nil, err
	}
	return &manifestConfig, nil
}
