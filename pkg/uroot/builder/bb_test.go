// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// Disable this until we are done switching to modules.
func testBBBuild(t *testing.T) {
	dir := t.TempDir()

	opts := Opts{
		Env: golang.Default(),
		Packages: []string{
			"github.com/u-root/u-root/pkg/uroot/test/foo",
			"github.com/u-root/u-root/cmds/core/elvish",
		},
		TempDir:   dir,
		BinaryDir: "bbin",
	}
	af := initramfs.NewFiles()
	var bbb BBBuilder
	if err := bbb.Build(af, opts); err != nil {
		t.Error(err)
	}

	mustContain := []string{
		"bbin/elvish",
		"bbin/foo",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q; archive: %v", name, af)
		}
	}
}
