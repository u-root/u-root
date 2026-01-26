// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/uefivars/vartest"
)

// main is needed to extract the testdata from a zip to temp dir, and to clean
// up the temp dir after
func TestMain(m *testing.M) {
	var setups = []struct {
		name      string
		zippath   string
		testevd   bool
		testevfsd bool
	}{
		{
			name:      "test efivars backend",
			zippath:   "testdata/sys_fw_efi_vars.zip",
			testevd:   true,
			testevfsd: false,
		},
		{
			name:      "test efivarfs backend",
			zippath:   "testdata/sys_fw_efivarfs.zip",
			testevd:   false,
			testevfsd: true,
		},
	}
	for _, w := range setups {
		efiVarDir, cleanup, err := vartest.SetupVarZip(w.zippath)
		if err != nil {
			log.Fatalf("couldn't setup testdata for efivar tests: %v", err)
		}

		// If there is only one directory in the zip, use that as the root.
		if entries, err := os.ReadDir(efiVarDir); err == nil && len(entries) == 1 && entries[0].IsDir() {
			efiVarDir = filepath.Join(efiVarDir, entries[0].Name())
		}

		EfiVarDir = "/tmp/invalid"
		EfiVarfsDir = "/tmp/invalid"
		if w.testevd {
			EfiVarDir = efiVarDir
		}
		if w.testevfsd {
			EfiVarfsDir = efiVarDir
		}
		rc := m.Run()

		cleanup()
		if rc != 0 {
			os.Exit(rc)
		}

	}
	os.Exit(0)
}
