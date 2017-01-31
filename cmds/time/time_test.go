// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"regexp"
	"testing"
)

// Ensure one second is printed to within some degree of accuracy.
func TestFalse(t *testing.T) {
	err := exec.Command("go", "build", "time.go").Run()
	if err != nil {
		t.Fatal("Cannot build time.go:", err)
	}
	out, err := exec.Command("./time", "sleep", "1").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	wantRegexp := regexp.MustCompile(`^real 1.00\d
user 0.00\d
sys 0.00\d
$`)
	if !wantRegexp.MatchString(string(out)) {
		t.Fatalf("Want regexp:\n%s\nGot:\n%s", wantRegexp, string(out))
	}
}
