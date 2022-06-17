// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bb

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/testutil"
)

func TestPackageRewriteFile(t *testing.T) {
	testutil.SkipInGithubActions(t)
	bin := filepath.Join(t.TempDir(), "foo")
	if err := BuildBusybox(golang.Default(), []string{"github.com/u-root/u-root/pkg/uroot/test/foo"}, false, bin); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("foo failed: %v %v", string(o), err)
	}
}
