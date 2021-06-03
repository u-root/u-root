// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Boot implements the Booter interface for booting from local storage.
type BootBooter struct {
	Type         string `json:"type"`
	DeviceIgnore string `json:"device_ignore"`
	KernelAppend string `json:"kernel_append"`
	KernelRemove string `json:"kernel_remove"`
	KernelReuse  string `json:"kernel_reuse"`
}

// NewBootBooter parses a boot entry config and returns a Booter instance, or
// an error if any
func NewBootBooter(config []byte) (Booter, error) {
	/* The configuration format for a Boot Booter entry is a JSON with the following structure:

		{
			"type": "boot",
			"device_ignore": "<pci devices or empty>",
			"kernel_append": "<kernel args to append or empty>",
			"kernel_remove": "<kernel args to remove or empty>",
			"kernel_reuse":  "<kernel args to reuse or empty>",
		}

	The JSON corisponds to the boot application command line variables:
	 "block": comma separated list of pci vendor and device ids to ignore (format vendor:device). E.g. 0x8086:0x1234,0x8086:0xabcd
	 "append": comma separated list of kernel params value to reuse from current kernel configuration
	 "remove": comma separated list of kernel params value to remove from parsed kernel configuration
	 "reuse": comma separated list of kernel params value to reuse from current kernel configuration
	*/

	log.Printf("Trying Boot Booter...")
	log.Printf("Config: %s", string(config))
	lb := BootBooter{}
	if err := json.Unmarshal(config, &lb); err != nil {
		return nil, err
	}
	log.Printf("BootBooter: %+v", lb)
	if lb.Type != "boot" {
		return nil, fmt.Errorf("wrong type for BootBooter: %s", lb.Type)
	}
	// the actual arguments validation is done in `Boot` to avoid duplicate code
	return &lb, nil
}

// Boot will run the boot procedure. In the case of BootBooter, it will call
// the `boot` command
func (lb *BootBooter) Boot(debugEnabled bool) error {
	var bootcmd []string

	bootcmd = []string{"boot"}

	if debugEnabled {
		bootcmd = append(bootcmd, "-v")
	}

	// validate arguments
	if lb.DeviceIgnore != "" {
		bootcmd = append(bootcmd, []string{"-block", lb.DeviceIgnore}...)
	}
	if lb.KernelAppend != "" {
		bootcmd = append(bootcmd, []string{"-append", lb.KernelAppend}...)
	}
	if lb.KernelRemove != "" {
		bootcmd = append(bootcmd, []string{"-remove", lb.KernelRemove}...)
	}
	if lb.KernelReuse != "" {
		bootcmd = append(bootcmd, []string{"-reuse", lb.KernelReuse}...)
	}

	log.Printf("Executing command: %v", bootcmd)
	cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
	}
	return nil
}

// TypeName returns the name of the booter type
func (lb *BootBooter) TypeName() string {
	return lb.Type
}
