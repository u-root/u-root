// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package smbios

import (
	"runtime"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestFromSysfs(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	if runtime.GOARCH != "amd64" {
		t.Skip("Test only supported in amd64 Qemu")
	}

	info, err := FromSysfs()
	if err != nil || info == nil {
		t.Errorf("FromSysfs() = '%q', '%v', want nil", info, err)
	}
}
