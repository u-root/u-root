// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

func TestBBBuild(t *testing.T) {
	// Only run this test when we are not using modules.
	if os.Getenv("GO111MODULE") != "off" {
		t.Skipf("Skipping non-modular test")
	}
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
	if err := bbb.Build(ulogtest.Logger{t}, af, opts); err != nil {
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
