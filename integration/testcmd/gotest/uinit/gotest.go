// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/u-root/u-root/integration/internal/gotest"
	"github.com/u-root/u-root/pkg/sh"
	"golang.org/x/sys/unix"
)

// Mount a vfat volume and run the tests within.
func main() {
	sh.RunOrDie("mkdir", "/testdata")
	if os.Getenv("UROOT_USE_9P") == "1" {
		sh.RunOrDie("mount", "-t", "9p", "tmpdir", "/testdata")
	} else {
		sh.RunOrDie("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")
	}

	gotest.WalkTests("/testdata/tests", func(i int, path, pkgName string) {
		runMsg := fmt.Sprintf("TAP: # running %d - %s", i, pkgName)
		passMsg := fmt.Sprintf("TAP: ok %d - %s", i, pkgName)
		failMsg := fmt.Sprintf("TAP: not ok %d - %s", i, pkgName)
		log.Println(runMsg)

		ctx, cancel := context.WithTimeout(context.Background(), 25000*time.Millisecond)
		defer cancel()
		cmd := exec.CommandContext(ctx, path)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err := cmd.Run()

		if ctx.Err() == context.DeadlineExceeded {
			log.Println("TAP: # timed out")
			log.Println(failMsg)
		} else if err == nil {
			log.Println(passMsg)
		} else {
			log.Println(err)
			log.Println(failMsg)
		}
	})

	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
