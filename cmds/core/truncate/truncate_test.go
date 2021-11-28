// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var truncateTests = []struct {
	flags           []string
	ret             int   // -1 for an expected error
	genTargetFile   bool  // if set, a temporary target file will be created before the test (used for -c)
	genRefFile      bool  // if set, a temporary reference file will be created before the test (used for -r)
	fileExistsAfter bool  // if set, we expect that the file will exist after the test
	size            int64 // -1 to signal we don't care for size test, early continue
	initTargetSize  int64 // only used when genTargetFile is true
	initRefSize     int64 // only used when genRefFile is true
}{
	{
		// Without args
		flags: []string{},
		ret:   -1,
	}, {
		// Invalid, valid args, but -s or -r is missing
		flags: []string{"-c"},
		ret:   -1,
	}, {
		// Invalid, invalid flag
		flags: []string{"-x"},
		ret:   -1,
	}, {
		// Invalid, invalid flag combo
		flags: []string{"-s", "1", "-r"},
		ret:   -1,
	}, {
		// Valid, file does not exist
		flags:           []string{"-s", "0"},
		ret:             0,
		genTargetFile:   false,
		fileExistsAfter: true,
		size:            0,
	}, {
		// Valid, file does exist and is smaller
		flags:           []string{"-s", "1"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  0,
		size:            1,
	}, {
		// Valid, file does exist and is bigger
		flags:           []string{"-s", "1"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  2,
		size:            1,
	}, {
		// Valid, file does exist grow
		flags:           []string{"-s", "+3K"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  2,
		size:            2 + 3*1024,
	}, {
		// Valid, file does exist shrink
		flags:           []string{"-s", "-3"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  5,
		size:            2,
	}, {
		// Valid, file does exist shrink lower than 0
		flags:           []string{"-s", "-3M"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  2,
		size:            0,
	}, {
		// Weird GNU behavior that this actual error is ignored
		flags:           []string{"-c", "-s", "2"},
		ret:             0,
		genTargetFile:   false,
		fileExistsAfter: false,
		size:            -1,
	}, {
		// Existing one
		flags:           []string{"-c", "-s", "3"},
		ret:             0,
		genTargetFile:   true,
		fileExistsAfter: true,
		initTargetSize:  0,
		size:            3,
	}, {
		// Valid ref file create
		flags:           []string{"-r"},
		ret:             0,
		genTargetFile:   false,
		genRefFile:      true,
		fileExistsAfter: true,
		initTargetSize:  0,
		initRefSize:     12,
		size:            12,
	}, {
		// Valid ref file existing, grow
		flags:           []string{"-r"},
		ret:             0,
		genTargetFile:   true,
		genRefFile:      true,
		fileExistsAfter: true,
		initTargetSize:  0,
		initRefSize:     17,
		size:            17,
	}, {
		// Valid ref file existing, shrink
		flags:           []string{"-r"},
		ret:             0,
		genTargetFile:   true,
		genRefFile:      true,
		fileExistsAfter: true,
		initTargetSize:  76,
		initRefSize:     18,
		size:            18,
	}, {
		// Invalid, ref file doesn't exist
		flags:           []string{"-r"},
		ret:             -1,
		genTargetFile:   false,
		genRefFile:      false,
		fileExistsAfter: false,
	},
}

// TestTruncate implements a table-driven test.
func TestTruncate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for i, test := range truncateTests {
		targetFile := filepath.Join(tmpDir, fmt.Sprintf("target%d", i))
		refFile := filepath.Join(tmpDir, fmt.Sprintf("ref%d", i))
		if test.genTargetFile {
			data := make([]byte, test.initTargetSize)
			if err := ioutil.WriteFile(targetFile, data, 0o600); err != nil {
				t.Errorf("Failed to create test file %s: %v", targetFile, err)
				continue
			}
		}
		if test.genRefFile {
			data := make([]byte, test.initRefSize)
			if err := ioutil.WriteFile(refFile, data, 0o600); err != nil {
				t.Errorf("Failed to create test file %s: %v", targetFile, err)
				continue
			}
		}
		// Execute truncate.go
		cmd := testutil.Command(t, append(test.flags, refFile, targetFile)...)
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
		st, err := os.Stat(targetFile)
		if err != nil && test.fileExistsAfter {
			t.Fatalf("Expected %s to exist, but os.Stat() returned error: %v\n", targetFile, err)
		}
		if s := st.Size(); s != test.size {
			t.Fatalf("Expected that %s has size: %d, but it has size: %d\n", targetFile, test.size, s)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
