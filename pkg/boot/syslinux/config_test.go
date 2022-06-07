// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/boottest"
)

// Enable this to generate new configs.
func DISABLEDTestGenerateConfigs(t *testing.T) {
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			imgs, err := ParseLocalConfig(context.Background(), configPath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.ToJSONFile(imgs, test); err != nil {
				t.Errorf("failed to generate file: %v", err)
			}
		})
	}
}

func TestConfigs(t *testing.T) {
	// find all saved configs
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimSuffix(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			want, err := os.ReadFile(test)
			if err != nil {
				t.Errorf("Failed to read test json '%v':%v", test, err)
			}

			imgs, err := ParseLocalConfig(context.Background(), configPath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.CompareImagesToJSON(imgs, want); err != nil {
				t.Errorf("ParseLocalConfig(): %v", err)
			}
		})
	}
}
