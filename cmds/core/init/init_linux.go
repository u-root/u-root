// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/kmodule"
)

// installModules installs kernel modules (.ko files) from /lib/modules.
// Useful for modules that need to be loaded for boot (ie a network
// driver needed for netboot)
func installModules() {
	modulePattern := "/lib/modules/*.ko"
	files, err := filepath.Glob(modulePattern)
	if err != nil {
		log.Printf("installModules: %v", err)
		return
	}
	if len(files) == 0 {
		// Since it is common for users to not have modules, no need to error or
		// print if there are none to install
		return
	}

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("installModules: can't open %q: %v", filename, err)
			continue
		}
		// Module flags are passed to the command line in the form modulename.flag=val
		// And must be passed to FileInit as flag=val to be installed properly
		moduleName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		flags := cmdline.FlagsForModule(moduleName)
		err = kmodule.FileInit(f, flags, 0)
		f.Close()
		if err != nil {
			log.Printf("installModules: can't install %q: %v", filename, err)
		}
	}
}

func init() {
	osInitGo = runOSInitGo
}

func runOSInitGo() {
	// Install modules before exec-ing into user mode below
	installModules()

	// systemd is "special". If we are supposed to run systemd, we're
	// going to exec, and if we're going to exec, we're done here.
	// systemd uber alles.
	initFlags := cmdline.GetInitFlagMap()

	// systemd gets upset when it discovers it isn't really process 1, so
	// we can't start it in its own namespace. I just love systemd.
	systemd, present := initFlags["systemd"]
	systemdEnabled, boolErr := strconv.ParseBool(systemd)
	if present && boolErr == nil && systemdEnabled {
		if err := syscall.Exec("/inito", []string{"/inito"}, os.Environ()); err != nil {
			log.Printf("Lucky you, systemd failed: %v", err)
		}
	}
}
