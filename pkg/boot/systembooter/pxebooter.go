// Copyright 2021-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/ulog"
)

var errWrongType = errors.New("wrong Type")

// PxeBooter implements the Booter interface for booting PXE
type PxeBooter struct {
	Type   string `json:"type"`
	IPV6   string `json:"ipv6"`
	IPV4   string `json:"ipv4"`
	Server string `json:"server"`
	File   string `json:"file"`
	Cmd    string `json:"cmd"`
}

// NewPxeBooter parses a boot entry config and returns a Booter instance, or an
// error if any
func NewPxeBooter(config []byte, l ulog.Logger) (Booter, error) {
	l.Printf("Trying PxeBooter...")
	l.Printf("Config: %s", string(config))
	nb := PxeBooter{}
	if err := json.Unmarshal(config, &nb); err != nil {
		return nil, err
	}
	l.Printf("PxeBooter: %+v", nb)
	if nb.Type != "pxeboot" {
		return nil, fmt.Errorf("%w:%q", errWrongType, nb.Type)
	}
	return &nb, nil
}

// Boot will run the boot procedure. In the case of PxeBooter, it will call the
// `pxeboot` command
func (nb *PxeBooter) Boot(debugEnabled bool) error {
	var bootcmd []string
	l := ulog.Null
	bootcmd = []string{"pxeboot"}

	if debugEnabled {
		bootcmd = append(bootcmd, "-v")
		l = ulog.Log
	}

	if nb.File != "" {
		bootcmd = append(bootcmd, "-file", nb.File)
	}

	// IPV4 and IPV6 default to enabled/true. Only set the disables.
	if nb.IPV6 == "false" {
		bootcmd = append(bootcmd, "-ipv6=false")
	}
	if nb.IPV4 == "false" {
		bootcmd = append(bootcmd, "-ipv4=false")
	}

	l.Printf("Executing command: %v", bootcmd)
	cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing %v: %w", cmd, err)
	}
	return nil
}

// TypeName returns the name of the booter type
func (nb *PxeBooter) TypeName() string {
	return nb.Type
}
