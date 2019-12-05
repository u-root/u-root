// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestRSDP(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	_, err := GetRSDP()
	if err != nil {
		t.Fatalf("GetRSDP: got %v, want nil", err)
	}
}
