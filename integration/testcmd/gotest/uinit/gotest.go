// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
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

var coverProfile = flag.String("coverprofile", "", "Filename to write coverage data to")

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

// AppendFile takes two filepaths and concatenates the files at those.
func AppendFile(srcFile, targetFile string) error {
	cov, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer cov.Close()

	out, err := os.OpenFile(targetFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, cov); err != nil {
		if err := out.Close(); err != nil {
			return err
		}
		return fmt.Errorf("error appending %s to %s: %v", srcFile, targetFile, err)
	}

	return out.Close()
}

// runTest mounts a vfat or 9pfs volume and runs the tests within.
func runTest() error {
	flag.Parse()

	if err := os.MkdirAll("/testdata", 0755); err != nil {
		log.Fatalf("Couldn't create testdata directory: %v", err)
	}

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

		args := []string{"-test.v"}
		coverFile := filepath.Join(filepath.Dir(path), "coverage.txt")
		if len(*coverProfile) > 0 {
			args = append(args, "-test.coverprofile", coverFile)
		}

		cmd := exec.CommandContext(ctx, path, args...)
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
			log.Printf("Failed to start %q: %v", path, err)
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
		if err := w.Close(); err != nil {
			log.Printf("Failed to close pipe: %v", err)
		}
		if err := j.Wait(); err != nil {
			log.Printf("Failed to stop test2json: %v", err)
		}

		if err := AppendFile(coverFile, *coverProfile); err != nil {
			log.Printf("Could not append to cover file: %v", err)
		}
	})
	return nil
}

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
