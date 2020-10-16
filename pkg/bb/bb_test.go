// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bb

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
)

func TestPackageRewriteFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bin := filepath.Join(dir, "foo")
	if err := BuildBusybox(golang.Default(), []string{"github.com/u-root/u-root/pkg/uroot/test/foo"}, false, bin); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("foo failed: %v %v", string(o), err)
	}
}
