// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/cpio"
)

// OSImage represents a bootable OS package.
type OSImage interface {
	fmt.Stringer

	// ExecutionInfo prints information about the OS image. A user should
	// be able to use the kexec command line tool to execute the OSImage
	// given the printed information.
	ExecutionInfo(log *log.Logger)

	// Execute kexec's the OS image: it loads the OS image into memory and
	// jumps to the kernel's entry point.
	Execute() error

	// Pack writes the OS image to the modules directory of sw and the
	// package type to package_type of sw.
	Pack(sw cpio.RecordWriter) error
}

var (
	osimageMap = map[string]func(*cpio.Archive) (OSImage, error){
		"linux": func(a *cpio.Archive) (OSImage, error) {
			return NewLinuxImageFromArchive(a)
		},
		"multiboot": newMultibootImage,
	}
)
