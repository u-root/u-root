// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/uefivars/vartest"
)

// main is needed to extract the testdata from a zip to temp dir, and to clean
// up the temp dir after
func TestMain(m *testing.M) {
	efiVarDir, cleanup, err := vartest.SetupVarZip("testdata/sys_fw_efi_vars.zip")
	if err != nil {
		panic(err)
	}
	EfiVarDir = efiVarDir
	rc := m.Run()
	cleanup()
	os.Exit(rc)
}
