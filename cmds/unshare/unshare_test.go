// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	
	"github.com/u-root/u-root/shared/testutil"
)

var (
	testPath = "."
	// if true removeAll the testPath on the end
	remove = true
)

type test struct {
	args    []string
	expects string
}

// TODO: well, there's a lot todo. And it's easy with this test to
// really mess up your system, what with all the mounting it has to do,
// so ... not sure.
// Also, note, this is basically a wrapper for exec, and all the tests done there
// need not be repeated here. So it's not quite clear what a test should do.
func init() {
	fatal = func(s string, i ...interface{}) {
		fmt.Printf(s, i...)
	}
	fatal("HI THERE!")
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestUnshareSimple(t *testing.T) {
	tmpDir, xPath := testutil.CompileInTempDir(t)
	if remove {
		defer os.RemoveAll(tmpDir)
	}
	t.Logf("tmpDir %v, xPath %v", tmpDir, xPath)
}

