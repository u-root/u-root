// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

func TestYes(t *testing.T) {
	type test struct {
		name        string
		in          []string
		expected    string
		closeStdout bool
	}
	tests := []test{
		{name: "noParameterCloseTest", in: []string{}, expected: "y", closeStdout: true},
		{name: "noParameterKillTest", in: []string{}, expected: "y", closeStdout: false},
		{name: "oneParameterCloseTest", in: []string{"hi"}, expected: "hi", closeStdout: true},
		{name: "oneParameterKillTest", in: []string{"hi"}, expected: "hi", closeStdout: false},
		{name: "fourParameterCloseTest", in: []string{"hi", "how", "are", "you"}, expected: "hi how are you", closeStdout: true},
		{name: "fourParameterKillTest", in: []string{"hi", "how", "are", "you"}, expected: "hi how are you", closeStdout: false},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			cmd := testutil.Command(t, v.in...)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatalf("Failed to get stdout of yes command: got %v, want nil", err)
			}
			if err := cmd.Start(); err != nil {
				t.Fatalf("Failed to start yes command: got %v, want nil", err)
			}
			in := bufio.NewScanner(stdout)
			for i := 0; i < 1000; i++ {
				if !in.Scan() {
					t.Fatalf("Could not scan: got %v, want nil", in.Err())
				}
				if text := in.Text(); text != v.expected {
					t.Errorf("Got %v at iteration %d, want %v", text, i, v.expected)
					break
				}
			}
			if v.closeStdout == true {
				if err := stdout.Close(); err != nil {
					t.Fatalf("Close standard out pipe: got %v, want nil", err)
				}
				if err := cmd.Wait(); err == nil || err.Error() != "signal: broken pipe" {
					t.Errorf("Complete the child process: got %v, want 'signal: broken pipe'", err)
				}
			} else {
				if err := cmd.Process.Kill(); err != nil {
					t.Errorf("Kill the child process: got %v, want nil", err)
				}
			}
		})
	}
}
