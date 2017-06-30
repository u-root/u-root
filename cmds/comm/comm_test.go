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

	"github.com/u-root/u-root/pkg/testutil"
)

var commTests = []struct {
	flags []string
	in1   string
	in2   string
	out   string
}{
	{
		// Empty files
		flags: []string{},
		in1:   "",
		in2:   "",
		out:   "",
	}, {
		// Equal files
		flags: []string{},
		in1:   "Line1\nlIne2\nline2\nline3",
		in2:   "Line1\nlIne2\nline2\nline3",
		out:   "\t\tLine1\n\t\tlIne2\n\t\tline2\n\t\tline3\n",
	}, {
		// Empty file 1
		flags: []string{},
		in1:   "",
		in2:   "Line1\nlIne2\n",
		out:   "\tLine1\n\tlIne2\n",
	}, {
		// Empty file 2
		flags: []string{},
		in1:   "Line1\nlIne2\n",
		in2:   "",
		out:   "Line1\nlIne2\n",
	}, {
		// Mix of matchine lines
		flags: []string{},
		in1:   "1\n3\n5\n",
		in2:   "2\n3\n4\n",
		out:   "1\n\t2\n\t\t3\n\t4\n5\n",
	},
}

// TestComm implements a table-drivent test.
func TestComm(t *testing.T) {
	// Compile comm.
	tmpDir, commPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, test := range commTests {
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
		args := append(append([]string{}, test.flags...), files[0], files[1])
		cmd := exec.Command(commPath, args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("Comm exited with error: %v; output:\n%s", err, string(output))
		} else if string(output) != test.out {
			t.Errorf("Fail: want\n %#v\n got\n %#v", test.out, string(output))
		}
	}
}
