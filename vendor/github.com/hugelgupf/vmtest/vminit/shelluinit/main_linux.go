// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command shelluinit runs commands from an elvish script.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hugelgupf/vmtest/guest"
)

func runTest() error {
	defer guest.CollectKernelCoverage()

	// Run the test script test.sh
	test := "/mount/9p/shelltest/test.sh"
	if _, err := os.Stat(test); os.IsNotExist(err) {
		return errors.New("could not find any test script to run")
	}
	cmd := exec.Command("gosh", test)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test.sh ran unsuccessfully: %v", err)
	}
	return nil
}

func main() {
	if err := runTest(); err != nil {
		log.Printf("Tests failed: %v", err)
	} else {
		log.Print("TESTS PASSED MARKER")
	}
}
