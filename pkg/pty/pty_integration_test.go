// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package pty

import (
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Timeout: 120 * time.Second,
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/pty"}, o)
}
