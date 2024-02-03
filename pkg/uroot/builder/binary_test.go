// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"github.com/u-root/uio/ulog/ulogtest"
)

func TestBinaryBuild(t *testing.T) {
	dir := t.TempDir()

	opts := Opts{
		Env: golang.Default(golang.DisableCGO()),
		Packages: []string{
			"../test/foo",
			"../../../cmds/core/elvish",
			"github.com/u-root/u-root/cmds/core/init",
			"cmd/test2json",
		},
		TempDir:   dir,
		BinaryDir: "bbin",
		BuildOpts: &golang.BuildOpts{},
	}
	af := initramfs.NewFiles()
	var b BinaryBuilder
	if err := b.Build(ulogtest.Logger{TB: t}, af, opts); err != nil {
		t.Fatalf("Build(%v, %v); %v != nil", af, opts, err)
	}

	mustContain := []string{
		"bbin/elvish",
		"bbin/foo",
		"bbin/init",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q; archive: %v", name, af)
		}
	}
}
