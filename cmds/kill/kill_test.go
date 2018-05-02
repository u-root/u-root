// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestKillProcess(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "kill")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Re-exec the test binary itself to emulate "sleep 1".
	cmd := exec.Command("/bin/sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}

	// from the orignal. hokey .1 second wait for the process to start. Racy.
	time.Sleep(100 * time.Millisecond)

	if err := testutil.Command(t, "-9", fmt.Sprintf("%d", cmd.Process.Pid)).Run(); err != nil {
		t.Errorf("Could not spawn first kill: %v", err)
	}

	if err := cmd.Wait(); err == nil {
		t.Errorf("Test process succeeded, but expected to fail")
	}
}

func TestBadInvocations(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "kill")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, v := range []struct {
		args []string
		err  string
	}{
		{args: []string{"-1w34"}, err: "1w34 is not a valid signal\n"},
		{args: []string{"-s"}, err: eUsage + "\n"},
		{args: []string{"-s", "a"}, err: "a is not a valid signal\n"},
		{args: []string{"a"}, err: "Some processes could not be killed: [a: arguments must be process or job IDS]\n"},
		{args: []string{"--signal"}, err: eUsage + "\n"},
		{args: []string{"--signal", "a"}, err: "a is not a valid signal\n"},
		{args: []string{"-1", "a"}, err: "Some processes could not be killed: [a: arguments must be process or job IDS]\n"},
	} {
		cmd := testutil.Command(t, v.args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if e := stderr.String(); e != v.err {
			t.Errorf("kill %v failed: got %s, want %s", v.args, e, v.err)
		}
		if err == nil {
			t.Errorf("kill %v failed: got nil, want err", v.args)
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
