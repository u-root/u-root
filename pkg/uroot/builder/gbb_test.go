// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

func TestGBBBuild(t *testing.T) {
	dir := t.TempDir()

	opts := Opts{
		Env: golang.Default(),
		Packages: []string{
			"../test/foo",
			"../../../cmds/core/elvish",
		},
		TempDir:   dir,
		BinaryDir: "bbin",
		BuildOpts: &gbbgolang.BuildOpts{},
	}
	af := initramfs.NewFiles()
	var gbb GBBBuilder
	if err := gbb.Build(ulogtest.Logger{t}, af, opts); err != nil {
		t.Fatalf("Build(%v, %v); %v != nil", af, opts, err)
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
