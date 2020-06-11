// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package loop

import (
	"testing"

	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/mount/loop"}, nil)
}
