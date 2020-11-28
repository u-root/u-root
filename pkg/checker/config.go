// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"encoding/json"
	"fmt"
)

// NewConfig parses a JSON configuration buffer into a Config object.
// NOTE: an empty checklist is considered valid, and the user is responsible
// to decide how to deal with that.
func NewConfig(buf []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	// TODO check version string after partial parsing
	for idx, item := range config.Checklist {
		if item.Check.Name == "" {
			return nil, fmt.Errorf("invalid check #%d: no name specified", idx+1)
		}
		if item.Check.Description == "" {
			return nil, fmt.Errorf("invalid check #%d: no description specified", idx+1)
		}
	}
	return &config, nil
}

// Config is the checker configuration structure.
type Config struct {
	Version   string       `json:"version"`
	Checklist []ConfigItem `json:"checklist"`
}

// ConfigItem wraps a check configuration and an optional remediation
// configuration.
type ConfigItem struct {
	Check           ConfigCheck        `json:"check"`
	Remediation     *ConfigRemediation `json:"remediation,omitempty"`
	ContinueOnError bool               `json:"continue_on_error"`
}

// ConfigCheck is a sub-field of ConfigItem representing a check.
type ConfigCheck struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Args        []interface{} `json:"args,omitempty"`
}

// ConfigRemediation is a sub-field of ConfigItem representing a
// remediation.
type ConfigRemediation struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Args        []interface{} `json:"args,omitempty"`
}
