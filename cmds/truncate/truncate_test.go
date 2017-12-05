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

var truncateTests = []struct {
	flags           []string
	ret             int   // -1 for an expected error
	genFile         bool  // if set, a temporary file will be created before the test (used for -c)
	fileExistsAfter bool  // if set, we expect that the file will exist after the test
	size            int64 // -1 to signal we don't care for size test, early continue
	initSize        int64 // only used when genFile is true
}{
	{
		// Without args
		flags: []string{},
		ret:   -1,
	}, {
		// Invalid, valid args, but -s is missing
		flags: []string{"-c"},
		ret:   -1,
	}, {
		// Invalid, invalid flag
		flags: []string{"-x"},
		ret:   -1,
	}, {
		// Valid, file does not exist
		flags:           []string{"-s", "0"},
		ret:             0,
		genFile:         false,
		fileExistsAfter: true,
		size:            0,
	}, {
		// Valid, file does exist and is smaller
		flags:           []string{"-s", "1"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        0,
		size:            1,
	}, {
		// Valid, file does exist and is bigger
		flags:           []string{"-s", "1"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            1,
	}, {
		// Valid, file does exist grow
		flags:           []string{"-s", "+3K"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            2 + 3*1024,
	}, {
		// Valid, file does exist shrink
		flags:           []string{"-s", "-3"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        5,
		size:            2,
	}, {
		// Valid, file does exist shrink lower than 0
		flags:           []string{"-s", "-3M"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            0,
	}, {
		// Weird GNU behavior that this actual error is ignored
		flags:           []string{"-c", "-s", "2"},
		ret:             0,
		genFile:         false,
		fileExistsAfter: false,
		size:            -1,
	}, {
		// Existing one
		flags:           []string{"-c", "-s", "3"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        0,
		size:            3,
	},
}

// TestTruncate implements a table-driven test.
func TestTruncate(t *testing.T) {
	// Compile truncate.
	tmpDir, truncatePath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for i, test := range truncateTests {
		testfile := filepath.Join(tmpDir, fmt.Sprintf("txt%d", i))
		if test.genFile {
			data := make([]byte, test.initSize)
			if err := ioutil.WriteFile(testfile, data, 0600); err != nil {
				t.Errorf("Failed to create test file %s: %v", testfile, err)
				continue
			}
		}
		// Execute truncate.go
		args := append(append([]string{}, test.flags...), testfile)
		cmd := exec.Command(truncatePath, args...)
		err := cmd.Run()
		if err != nil {
			if test.ret == 0 {
				t.Fatalf("Truncate exited with error: %v, but return code %d expected\n", err, test.ret)
			} else if test.ret == -1 { // expected error, nothing more to see
				continue
			}
			t.Fatalf("Truncate exited with error: %v, test specified: %d, something is terribly wrong\n", err, test.ret)
		}
		if test.size == -1 {
			continue
		}
		st, err := os.Stat(testfile)
		if err != nil && test.fileExistsAfter {
			t.Fatalf("Expected %s to exist, but os.Stat() retuned error: %v\n", testfile, err)
		}
		if s := st.Size(); s != test.size {
			t.Fatalf("Expected that %s has size: %d, but it has size: %d\n", testfile, test.size, s)
		}
	}
}
