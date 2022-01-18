// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/smbios"}, &vmtest.Options{
		QEMUOpts: qemu.Options{},
	})
}
