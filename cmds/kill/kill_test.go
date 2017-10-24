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
	"path/filepath"
	"testing"
	"time"
)

// Run the command, with the optional args, and return a string
// for stdout, stderr, and an error.
func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestKillProcess(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "KillTest")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	killtestpath := filepath.Join(tmpDir, "killtest.exe")
	o, e, err := run(exec.Command("go", "build", "-o", killtestpath, "."))
	if err != nil {
		t.Fatalf("go build -o %s cmds/kill: %v, %s", killtestpath, err, o+":"+e)
	}

	t.Logf("Built %v for test", killtestpath)

	// Re-exec the test binary itself to emulate "sleep 1".
	cmd := exec.Command("/bin/sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start test process: %v, output %s", err, string(e))
	}

	// from the orignal. hokey .1 second wait for the process to start. Racy.
	time.Sleep(100 * time.Millisecond)

	if o, e, err = run(exec.Command(killtestpath, "-9", fmt.Sprintf("%d", cmd.Process.Pid))); err != nil {
		t.Errorf("Could not spawn first kill: %v", err)
	}

	t.Logf("Ran kill: output :%s:, extra info :%v", o, e)
	if err = cmd.Wait(); err == nil {
		t.Errorf("Test process succeeded, but expected to fail")
	}

	// now this is a little weird. We're going to try to kill it again.
	// Arguably, this should be done in another test, but finding a process
	// you just "know" does not exist is tricky. What PID do you use?
	// So we just kill the one we just killed; it should get an error.
	// If not, something's wrong.
	if _, _, err = run(exec.Command(killtestpath, "-9", fmt.Sprintf("%d", cmd.Process.Pid))); err == nil {
		t.Fatalf("Second kill: got nil, want error")
	}
}

func TestBadInvocations(t *testing.T) {
	var (
		tab = []struct {
			a   []string
			err string
		}{
			{a: []string{"-1w34"}, err: "1w34 is not a valid signal\n"},
			{a: []string{"-s"}, err: eUsage + "\n"},
			{a: []string{"-s", "a"}, err: "a is not a valid signal\n"},
			{a: []string{"a"}, err: "Some processes could not be killed: [a: arguments must be process or job IDS]\n"},
			{a: []string{"--signal"}, err: eUsage + "\n"},
			{a: []string{"--signal", "a"}, err: "a is not a valid signal\n"},
			{a: []string{"-1", "a"}, err: "Some processes could not be killed: [a: arguments must be process or job IDS]\n"},
		}
	)

	tmpDir, err := ioutil.TempDir("", "KillTest")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	killtestpath := filepath.Join(tmpDir, "killtest.exe")
	o, e, err := run(exec.Command("go", "build", "-o", killtestpath, "."))
	if err != nil {
		t.Fatalf("go build -o %s cmds/kill: %v, %s", killtestpath, err, o+":"+e)
	}

	t.Logf("Built %v for test", killtestpath)
	for _, v := range tab {
		o, e, err := run(exec.Command(killtestpath, v.a...))
		t.Logf("%v: %s %s %v", v.a, o, e, err)
		if e != v.err {
			t.Errorf("Kill for '%v' failed: got '%s', want '%s'", v.a, e, v.err)
		}
		if err == nil {
			t.Errorf("Kill for '%v' failed: got nil, want err", v.a)
		}
	}
}
