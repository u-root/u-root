// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	flag "github.com/spf13/pflag"
)

const (
	testDataDir = "testdata"
)

func resetFlags() {
	*flagFromDump = ""
	*flagType = nil
}

func testOutput(t *testing.T, dumpFile string, args []string, expectedOutFile string) {
	dumpFile = filepath.Join(testDataDir, dumpFile)
	expectedOutFile = filepath.Join(testDataDir, expectedOutFile)
	actualOutFile := fmt.Sprintf("%s.actual", expectedOutFile)
	os.Remove(actualOutFile)
	os.Args = []string{os.Args[0], "--from-dump", dumpFile}
	os.Args = append(os.Args, args...)
	flag.Parse()
	defer resetFlags()
	out := bytes.NewBuffer(nil)
	if _, err := dmiDecode(out); err != nil {
		t.Errorf("%+v %+v %+v: error: %s", dumpFile, args, expectedOutFile, err)
		return
	}
	actualOut := out.Bytes()
	expectedOut, err := ioutil.ReadFile(expectedOutFile)
	if err != nil {
		t.Errorf("%+v %+v %+v: failed to load %s: %s", dumpFile, args, expectedOutFile, expectedOutFile, err)
		return
	}
	if bytes.Compare(actualOut, expectedOut) != 0 {
		ioutil.WriteFile(actualOutFile, actualOut, 0644)
		t.Errorf("%+v %+v %+v: output mismatch, see %s", dumpFile, args, expectedOutFile, actualOutFile)
		diffOut, _ := exec.Command("diff", "-u", expectedOutFile, actualOutFile).CombinedOutput()
		t.Errorf("%+v %+v %+v: diff:\n%s", dumpFile, args, expectedOutFile, string(diffOut))
	}
}

func TestDMIDecode(t *testing.T) {
	testOutput(t, "UX307LA.bin", nil, "UX307LA.txt")
	testOutput(t, "UX307LA.bin", []string{"-t", "system"}, "UX307LA.system.txt")
	testOutput(t, "UX307LA.bin", []string{"-t", "1,131"}, "UX307LA.1_131.txt")
}
