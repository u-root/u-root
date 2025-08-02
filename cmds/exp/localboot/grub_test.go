// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/mount/block"
)

func FuzzParseGrubCfg(f *testing.F) {
	tmpDir := f.TempDir()

	// no log output
	log.SetOutput(io.Discard)
	log.SetFlags(0)

	// get seed corpora from testdata_new files
	seeds, err := filepath.Glob("../../../pkg/boot/grub/testdata_new/*/*/*/grub.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	seeds2, err := filepath.Glob("../../../pkg/boot/grub/testdata_new/*/*/grub.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	seeds = append(seeds, seeds2...)
	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(string(seedBytes))
	}

	f.Fuzz(func(t *testing.T, data string) {
		if len(data) > 256000 {
			return
		}

		blockDevs := block.BlockDevices{&block.BlockDev{Name: tmpDir, FSType: "test", FsUUID: "07338180-4a96-4611-aa6a-a452600e4cfe"}}
		ParseGrubCfg(grubV2, blockDevs, data, tmpDir)
	})
}
