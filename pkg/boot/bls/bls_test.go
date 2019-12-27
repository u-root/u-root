// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bls

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

var fsRoot = "testdata"

var blsEntries = []struct {
	entry string
	err   string
}{
	{
		entry: "entry-1.conf",
	},
	{
		entry: "entry-2.conf",
		err:   "neither linux, efi, nor multiboot present in BootLoaderSpec config",
	},
}

func TestParseBLSEntries(t *testing.T) {
	dir := filepath.Join(fsRoot, "loader/entries")

	for _, tt := range blsEntries {
		t.Run(tt.entry, func(t *testing.T) {
			image, err := parseBLSEntry(filepath.Join(dir, tt.entry), dir)
			if err != nil {
				if tt.err == "" {
					t.Fatalf("Got error %v", err)
				}
				if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("Got error %v, expected error to contain %s", err, tt.err)
				}
				return
			}
			if tt.err != "" {
				t.Fatalf("Expected error %s, got no error", tt.err)
			}
			t.Logf("Got image: %s", image.String())
		})
	}
}

func TestScanBLSEntries(t *testing.T) {
	entries, err := ScanBLSEntries(ulogtest.Logger{t}, fsRoot)
	if err != nil {
		t.Errorf("Error scanning BLS entries: %v", err)
	}

	// TODO: have a better way of checking contents
	if len(entries) < 1 {
		t.Errorf("Expected at least BLS entry, found none")
	}
}
