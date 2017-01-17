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
	// Case sensitive
	{
		flags: []string{},
		in1:   "line1\nline2\nLine3\nlIne4\nline",
		in2:   "line1\nLine2\nline3\nlInes\nLINEZ",
		out:   "\t\tline1\n\tLine2\nline2\nLine3\nlIne4\nline\n\tline3\n\tlInes\n\tLINEZ\n",
	},
	// Case insensitive
	{
		flags: []string{"-i"},
		in1:   "line1\nline2\nLine3\nlIne4\nline",
		in2:   "line1\nLine2\nline3\nlInes\nLINEZ",
		out:   "\t\tline1\n\t\tline2\n\t\tLine3\nlIne4\n\tlInes\n\tLINEZ\nline\n",
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
		if output, err := cmd.Output(); err != nil {
			t.Errorf("Cannot get output of comm: %v", err)
		} else if string(output) != test.out {
			t.Errorf("Fail: want\n %#v\n got\n %#v", test.out, string(output))
		}
	}
}
