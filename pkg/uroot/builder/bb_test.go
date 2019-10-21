// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

func TestBBBuild(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

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

	var mustContain = []string{
		"bbin/elvish",
		"bbin/foo",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q; archive: %v", name, af)
		}
	}

}
