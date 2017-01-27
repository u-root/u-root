// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var commTests = []struct {
	flags []string
	in1   string
	in2   string
	out   string
}{
	{
		flags: []string{},
		in1:   "Line1\nlIne2\nline\nline3\nline4",
		in2:   "Line1\nlIne2\nline\nline3\nline4",
		out:   "\t\tLine1\n\t\tlIne2\n\t\tline\n\t\tline3\n\t\tline4\n",
	},
}

// Table-driven test of the comm command
func TestComm(t *testing.T) {
	for _, test := range commTests {
		// Create temporary directory
		tmpDir, err := ioutil.TempDir("", "TestComm")
		if err != nil {
			t.Error(err)
			continue
		}
		defer os.RemoveAll(tmpDir)

		// Write inputs into the two files
		var files [2]string
		for i, contents := range []string{test.in1, test.in2} {
			files[i] = filepath.Join(tmpDir, fmt.Sprintf("txt%d", i))
			if err := ioutil.WriteFile(files[i], []byte(contents), 0600); err != nil {
				t.Errorf("Failed to create test file %d: %v", i, err)
				continue
			}
		}

		// Execute comm.go
		args := append(append([]string{"run", "comm.go"}, test.flags...), files[0], files[1])
		cmd := exec.Command("go", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("Comm exited with error: %v; output:\n%s", err, string(output))
		} else if string(output) != test.out {
			t.Errorf("Fail: want\n %#v\n got\n %#v", test.out, string(output))
		}
	}
}
