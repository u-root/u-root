// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kexecbin offers a kexec API via a callout to kexec-tools.
//
// u-root's kexec implementation currently covers less use-cases than the one
// from kexec-tools.  The user has to embed a kexecbin program in the
// initramfs, and make it available in the PATH, then call the `KexecBin`
// function.  Please note that adding an external kexec implementation will
// increase the ramfs size more than the pure-Go implementation from u-root.
package kexecbin

import (
	"os"
	"os/exec"
)

var (
	// DeviceTreePaths is the virtual fs path for accessing device-tree
	// through Linux
	DeviceTreePaths = []string{"/sys/firmware/fdt", "/proc/device-tree"}
)

// KexecBin uses kexec-tools binary and runtime architecture detection
// to execute abritary files.
func KexecBin(kernelFilePath string, kernelCommandline string, initrdFilePath string, dtFilePath string) error {
	baseCmd, err := exec.LookPath("kexecbin")
	if err != nil {
		return err
	}

	var loadCommands []string
	loadCommands = append(loadCommands, "-l")
	loadCommands = append(loadCommands, kernelFilePath)

	if kernelCommandline != "" {
		loadCommands = append(loadCommands, "--command-line="+kernelCommandline)
	} else {
		loadCommands = append(loadCommands, "--reuse-cmdline")
	}

	if initrdFilePath != "" {
		loadCommands = append(loadCommands, "--initrd="+initrdFilePath)
	}

	if dtFilePath != "" {
		loadCommands = append(loadCommands, "--dtb="+dtFilePath)
	} else {
		for _, dtFilePath := range DeviceTreePaths {
			if _, err := os.Stat(dtFilePath); err == nil {
				loadCommands = append(loadCommands, "--dtb="+dtFilePath)
				break
			}
		}
	}

	// Load data into physical non reserved memory regions
	cmdLoad := exec.Command(baseCmd, loadCommands...)
	if err := cmdLoad.Run(); err != nil {
		return err
	}

	// Execute into new kernel
	cmdExec := exec.Command(baseCmd, "-e")
	return cmdExec.Run()
}
