// Copyright 2018 the u-root Authors. All rights reserved
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

	"github.com/u-root/u-root/pkg/mount"
	"golang.org/x/sys/unix"
)

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

var (
	coverprofile = flag.String("coverprofile", "", "Filename to write coverage data to")
)

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
	defer out.Close()

	if _, err := io.Copy(out, cov); err != nil {
		return fmt.Errorf("error appending %s to %s: %v", srcFile, targetFile, err)
	}
	return nil
}

// Mount a vfat volume and run the tests within.
func main() {
	flag.Parse()

	if err := os.MkdirAll("/testdata", 0755); err != nil {
		log.Fatalf("Couldn't create testdata: %v", err)
	}

	var (
		mp  *mount.MountPoint
		err error
	)
	if os.Getenv("UROOT_USE_9P") == "1" {
		mp, err = mount.Mount("tmpdir", "/testdata", "9p", "", 0)
	} else {
		mp, err = mount.Mount("/dev/sda1", "/testdata", "vfat", "", unix.MS_RDONLY)
	}
	if err != nil {
		log.Fatalf("Failed to mount test directory: %v", err)
	}
	defer mp.Unmount(0) //nolint:errcheck

	walkTests("/testdata/tests", func(path, pkgName string) {
		ctx, cancel := context.WithTimeout(context.Background(), 25000*time.Millisecond)
		defer cancel()

		r, w, err := os.Pipe()
		if err != nil {
			log.Printf("Failed to get pipe: %v", err)
			return
		}

		args := []string{"-test.v"}
		coverfile := filepath.Join(filepath.Dir(path), "coverage.txt")
		if len(*coverprofile) > 0 {
			args = append(args, "-test.coverprofile", coverfile)
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
			log.Printf("Failed to start %v: %v", path, err)
			return
		}

		j := exec.CommandContext(ctx, "test2json", "-t", "-p", pkgName)
		j.Stdin = r
		j.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := j.Start(); err != nil {
			log.Printf("Failed to start test2json: %v", err)
			return
		}

		// Don't do anything if the test fails. The log collector will
		// deal with it. ¯\_(ツ)_/¯
		cmd.Wait()
		// Close the pipe so test2json will quit.
		w.Close()
		j.Wait()

		if err := AppendFile(coverfile, *coverprofile); err != nil {
			log.Printf("Could not append to cover file: %s", err)
		}
	})

	log.Printf("GoTest Done")

	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
