// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/ulog"
)

// LocalBooter implements the Booter interface for booting from local storage.
type LocalBooter struct {
	Type       string `json:"type"`
	Method     string `json:"method"`
	DeviceGUID string `json:"device_guid"`
	Kernel     string `json:"kernel,omitempty"`
	KernelArgs string `json:"kernel_args,omitempty"`
	Initramfs  string `json:"ramfs,omitempty"`
}

// NewLocalBooter parses a boot entry config and returns a Booter instance, or
// an error if any
func NewLocalBooter(config []byte, l ulog.Logger) (Booter, error) {
	/*
		The configuration format for a LocalBooter entry is a JSON with the following structure:

		{
			"type": "localboot",
			"method": "<method>",
			"device_guid": "<device GUID or empty>"
			"kernel": "<kernel path or empty>",
			"kernel_args": "<kernel args or empty>",
			"ramfs": "<ramfs path or empty>",
		}

		`type` is always set to "localboot"
		`method` can be either "grub" or "path".
		    The "grub" method will look for grub.cfg or grub2.cfg on the specified device.
		    If no device is specified, it will look on all the attached storage devices,
		    sorted alphabetically as found in /dev. The first grub configuration that is
		    found is parsed, and kernel, kernel args and ramfs are extracted. Then the
		    kernel will be kexec'ed. If this fails, the next entry will NOT be tried,
		    and no other grub configs will be scanned. In case a grub config has no
		    valid boot entries, it is ignored and the next config will be used tried.
		    The "path" method requires a device GUID and kernel path to be specified. If
		    specified, it will also use kernel args and ramfs path. This method will look
		    for the given kernel on the given device, and will kexec the kernel using the
		    given, optional, kernel args and ramfs.
		`device_guid` is the GUID of the device to look for grub config or kernel and ramfs
		`kernel` is the path, relative to the device specified by `device_guid`, of the
		    kernel to be kexec'ed
		`kernel_args` is the optional string of kernel arguments to be passed.
		`ramfs` is the path, relative to the device specified by `device_guid`, of the ramfs
		    to be used for kexec'ing into the target kernel.
	*/
	l.Printf("Trying LocalBooter...")
	l.Printf("Config: %s", string(config))
	lb := LocalBooter{}
	if err := json.Unmarshal(config, &lb); err != nil {
		return nil, err
	}
	l.Printf("LocalBooter: %+v", lb)
	if lb.Type != "localboot" {
		return nil, fmt.Errorf("%w:%q", errWrongType, lb.Type)
	}
	// the actual arguments validation is done in `Boot` to avoid duplicate code
	return &lb, nil
}

// Boot will run the boot procedure. In the case of LocalBooter, it will call
// the `localboot` command
func (lb *LocalBooter) Boot(debugEnabled bool) error {
	var bootcmd []string
	l := ulog.Null
	if debugEnabled {
		bootcmd = []string{"localboot", "-d"}
		l = ulog.Log
	} else {
		bootcmd = []string{"localboot"}
	}

	// validate arguments
	if lb.Method == "grub" {
		bootcmd = append(bootcmd, "-grub")
	} else if lb.Method == "path" {
		bootcmd = append(bootcmd, []string{"-kernel", lb.Kernel}...)
		bootcmd = append(bootcmd, []string{"-guid", lb.DeviceGUID}...)
		if lb.Initramfs != "" {
			bootcmd = append(bootcmd, []string{"-initramfs", lb.Initramfs}...)
		}
		if lb.KernelArgs != "" {
			bootcmd = append(bootcmd, []string{"-cmdline", lb.KernelArgs}...)
		}
	} else {
		return fmt.Errorf("unknown boot method %s", lb.Method)
	}

	l.Printf("Executing command: %v", bootcmd)
	cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		l.Printf("Error executing %v: %v", cmd, err)
	}
	return nil
}

// TypeName returns the name of the booter type
func (lb *LocalBooter) TypeName() string {
	return lb.Type
}
