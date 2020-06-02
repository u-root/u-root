// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
		var wantConfigs interface{}
		if err = json.Unmarshal(testJSON, &wantConfigs); err != nil {
			t.Errorf("Failed to unmarshall test json '%v':%v", test, err)
		}

		configPath := strings.TrimRight(test, ".json")
		configs, err := ParseLocalConfig(context.Background(), configPath)
		if err != nil {
			t.Fatalf("Failed to parse %s: %v", test, err)
		}

		configJSON, _ := json.MarshalIndent(configs, "", "  ")
		var gotConfigs interface{}
		if err := json.Unmarshal(configJSON, &gotConfigs); err != nil {
			t.Error(err)
		}

		if !cmp.Equal(wantConfigs, gotConfigs) {
			t.Errorf("ParseLocalConfig() mismatch(-want, +got):\n%s", cmp.Diff(wantConfigs, gotConfigs))
		}
	}
}
