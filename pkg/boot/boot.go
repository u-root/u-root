// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package boot is the high-level interface for booting another operating
// system.
package boot

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/kexec"
)

// OSImage represents a bootable OS package.
type OSImage interface {
	fmt.Stringer

	// ExecutionInfo prints information about the OS image. A user should
	// be able to use the kexec command line tool to execute the OSImage
	// given the printed information.
	ExecutionInfo(log *log.Logger)

	// Load loads the OS image into kernel memory, ready for execution.
	Load() error
}

// Execute executes a previously loaded OSImage.
func Execute() error {
	return kexec.Reboot()
}
