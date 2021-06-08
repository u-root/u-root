// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/u-root/u-root/integration/testcmd/common"
	"golang.org/x/sys/unix"
)

const individualTestTimeout = 25 * time.Second

func walkTests(testRoot string, fn func(string, string)) error {
	return filepath.Walk(testRoot, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() || !strings.HasSuffix(path, ".test") || err != nil {
			return nil
		}
		t2, err := filepath.Rel(testRoot, path)
		if err != nil {
			return err
		}
		pkgName := filepath.Dir(t2)

		fn(path, pkgName)
		return nil
	})
}

func runTest() error {
	cleanup, err := common.MountSharedDir()
	if err != nil {
		return err
	}
	defer cleanup()
	defer common.CollectKernelCoverage()

	walkTests("/testdata/tests", func(path, pkgName string) {
		// Send the kill signal with a 500ms grace period.
		ctx, cancel := context.WithTimeout(context.Background(), individualTestTimeout+500*time.Millisecond)
		defer cancel()

		r, w, err := os.Pipe()
		if err != nil {
			log.Printf("Failed to get pipe: %v", err)
			return
		}

		cmd := exec.CommandContext(ctx, path, "-test.v", "-test.timeout", individualTestTimeout.String())
		cmd.Stdin, cmd.Stderr = os.Stdin, os.Stderr

		// Write to stdout for humans, write to w for the JSON converter.
		//
		// The test collector will gobble up JSON for statistics, and
		// print non-JSON for humans to consume.
		cmd.Stdout = io.MultiWriter(os.Stdout, w)

		// Start test in its own dir so that testdata is available as a
		// relative directory.
		cmd.Dir = filepath.Dir(path)
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start %v: %v", path, err)
			return
		}

		// The test2json is not run with a context as it does not
		// block. If we cancelled test2json with the same context as
		// the test, we may lose some of the last few lines.
		j := exec.Command("test2json", "-t", "-p", pkgName)
		j.Stdin = r
		j.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := j.Start(); err != nil {
			log.Printf("Failed to start test2json: %v", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			// Log for processing by test2json.
			fmt.Fprintf(w, "Error: test for %q exited early: %v", pkgName, err)
		}

		// Close the pipe so test2json will quit.
		w.Close()
		j.Wait()
	})
	return nil
}

// Mount a vfat volume and run the tests within.
func main() {
	if err := runTest(); err != nil {
		log.Printf("Tests failed: %v", err)
	} else {
		// The test infra is expecting this exact print.
		log.Print("TESTS PASSED MARKER")
	}

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to reboot: %v", err)
	}
}
