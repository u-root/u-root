// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diskboot

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestParseEmpty(t *testing.T) {
	config := ParseConfig("/", "/grub.cfg", nil)
	if config == nil {
		t.Error("Expected non-nil config")
	}
	if len(config.Entries) > 0 {
		t.Error("Expected no entries: got", config.Entries)
	}
}

func TestConfigs(t *testing.T) {
	// find all saved configs
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		testJSON, err := ioutil.ReadFile(test)
		if err != nil {
			t.Errorf("Failed to read test json '%v':%v", test, err)
		}
		testConfigs := make([]*Config, 0)
		err = json.Unmarshal(testJSON, &testConfigs)
		if err != nil {
			t.Errorf("Failed to unmarshall test json '%v':%v", test, err)
		}

		configPath := strings.TrimRight(test, ".json")
		configs := FindConfigs(configPath)
		configJSON, _ := json.Marshal(configs)
		t.Logf("Configs for %v\n%v", test, string(configJSON))

		if diff := deep.Equal(testConfigs, configs); diff != nil {
			t.Error(diff)
		}
	}
}
