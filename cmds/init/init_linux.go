// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"strconv"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
)

func init() {
	osInitGo = runOSInitGo
}

func runOSInitGo() {
	// systemd is "special". If we are supposed to run systemd, we're
	// going to exec, and if we're going to exec, we're done here.
	// systemd uber alles.
	initFlags := cmdline.GetInitFlagMap()
	// systemd gets upset when it discovers it isn't really process 1, so
	// we can't start it in its own namespace. I just love systemd.
	systemd, present := initFlags["systemd"]
	systemdEnabled, boolErr := strconv.ParseBool(systemd)
	if present && boolErr == nil && systemdEnabled == true {
		v := cmdList[0]
		debug("Exec %v", v)
		if err := syscall.Exec(v, []string{v}, envs); err != nil {
			log.Printf("Lucky you, systemd failed: %v", err)
		}
		// well, what a shame.
		cmdList = cmdList[1:]
		cmdCount++
	}
}
