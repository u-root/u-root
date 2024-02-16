// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gouinit runs Go tests in a guest VM.
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

	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/internal/json2test"
	"github.com/hugelgupf/vmtest/internal/testevent"
)

var (
	coverProfile          = flag.String("coverprofile", "", "Filename to write coverage data to")
	individualTestTimeout = flag.Duration("test_timeout", time.Minute, "timeout per Go package")
)

func walkTests(testRoot string, fn func(string, string)) error {
	return filepath.Walk(testRoot, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() || !strings.HasSuffix(path, ".test") {
			return nil
		}
		if err != nil {
			return err
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

	testEvents, err := guest.EventChannel[testevent.ErrorEvent]("/mount/9p/gotestdata/errors.json")
	if err != nil {
		return err
	}
	defer testEvents.Close()

	if err := run(testEvents); err != nil {
		_ = testEvents.Emit(testevent.ErrorEvent{
			Error: fmt.Sprintf("running tests failed: %v", err),
		})
		return err
	}
	return nil
}

func run(testEvents *guest.Emitter[testevent.ErrorEvent]) error {
	defer guest.CollectKernelCoverage()

	goTestEvents, err := guest.EventChannel[json2test.TestEvent]("/mount/9p/gotestdata/results.json")
	if err != nil {
		return err
	}
	defer goTestEvents.Close()

	return walkTests("/mount/9p/gotestdata/tests", func(path, pkgName string) {
		// Send the kill signal with a 500ms grace period.
		ctx, cancel := context.WithTimeout(context.Background(), *individualTestTimeout+500*time.Millisecond)
		defer cancel()

		r, w, err := os.Pipe()
		if err != nil {
			log.Printf("Failed to get pipe: %v", err)
			return
		}

		args := []string{"-test.v", "-test.bench=.", "-test.run=."}
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
			_ = testEvents.Emit(testevent.ErrorEvent{
				Binary: path,
				Error:  fmt.Sprintf("failed to start: %v", err),
			})
			log.Printf("Failed to start %q: %v", path, err)
			return
		}

		// The test2json is not run with a context as it does not
		// block. If we cancelled test2json with the same context as
		// the test, we may lose some of the last few lines.
		j := exec.Command("test2json", "-t", "-p", pkgName)
		j.Stdin = r
		j.Stdout, cmd.Stderr = goTestEvents, os.Stderr
		if err := j.Start(); err != nil {
			_ = testEvents.Emit(testevent.ErrorEvent{
				Binary: path,
				Error:  fmt.Sprintf("failed to start test2json: %v", err),
			})
			log.Printf("Failed to start test2json: %v", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			_ = testEvents.Emit(testevent.ErrorEvent{
				Binary: path,
				Error:  fmt.Sprintf("test exited with non-zero status: %v", err),
			})
			log.Printf("Error: test %q exited with non-zero status: %v", pkgName, err)
		}

		// Close the pipe so test2json will quit.
		if err := w.Close(); err != nil {
			log.Printf("Failed to close pipe: %v", err)
		}
		if err := j.Wait(); err != nil {
			_ = testEvents.Emit(testevent.ErrorEvent{
				Binary: path,
				Error:  fmt.Sprintf("test2json exited with non-zero status: %v", err),
			})
			log.Printf("Failed to stop test2json: %v", err)
		}

		if len(*coverProfile) > 0 {
			if err := AppendFile(coverFile, *coverProfile); err != nil {
				_ = testEvents.Emit(testevent.ErrorEvent{
					Binary: path,
					Error:  fmt.Sprintf("could not append to coverage file: %v", err),
				})
				log.Printf("Could not append to cover file: %v", err)
			}
		}
	})
}

func main() {
	flag.Parse()

	if err := runTest(); err != nil {
		log.Printf("Tests failed: %v", err)
	}
}
