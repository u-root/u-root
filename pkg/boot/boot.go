// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package boot is the high-level interface for booting another operating
// system from Linux using kexec.
package boot

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/uio/ulog"
)

// DefaultLinuxImageCfgFile is the default file name in tmp directory to write loaded LinuxImage info to.
const DefaultLinuxImageCfgFile = "linux_image_cfg.json"

// LoadOption is an optional parameter to Load.
type LoadOption func(*loadOptions)

type loadOptions struct {
	logger        ulog.Logger
	verbose       bool
	callKexecLoad bool
	// linuxImageCfgFile specifies where to writes loaded linuximage info to.
	linuxImageCfgFile string
}

func defaultLoadOptions() *loadOptions {
	return &loadOptions{
		logger:            ulog.Null,
		verbose:           false,
		callKexecLoad:     true,
		linuxImageCfgFile: filepath.Join(os.TempDir(), DefaultLinuxImageCfgFile),
	}
}

// WithLogger is a LoadOption that logs verbose debug output l.
func WithLogger(l ulog.Logger) LoadOption {
	return func(o *loadOptions) {
		o.verbose = (l != nil)
		if l == nil {
			l = ulog.Null
		}
		o.logger = l
	}
}

// Verbose is a LoadOption that logs to log.Default().
var Verbose = WithLogger(ulog.Log)

// WithVerbose enables verbose logging if verbose is set to true.
func WithVerbose(verbose bool) LoadOption {
	return func(o *loadOptions) {
		o.verbose = verbose
		if verbose {
			o.logger = ulog.Log
		} else {
			o.logger = ulog.Null
		}
	}
}

// WithDryRun is a LoadOption that makes sure no kexec_load syscall is called during Load.
func WithDryRun(dryRun bool) LoadOption {
	return func(o *loadOptions) {
		o.callKexecLoad = !dryRun
	}
}

// WithLinuxImageCfgFile allows user to specify a different linux image config file path.
func WithLinuxImageCfgFile(f string) LoadOption {
	return func(o *loadOptions) {
		o.linuxImageCfgFile = f
	}
}

// OSImage represents a bootable OS package.
type OSImage interface {
	fmt.Stringer

	// Label is a name or short description for this OSImage.
	//
	// Label is intended for boot menus.
	Label() string

	// Rank the priority of the images for boot menus.
	//
	// The larger the number, the prior the image shows in the menu.
	Rank() int

	// Edit the kernel command line if possible. Must be called before
	// Load.
	Edit(func(cmdline string) string)

	// Load loads the OS image into kernel memory, ready for execution.
	//
	// After Load is called, call boot.Execute() to stop Linux and boot the
	// loaded OSImage.
	Load(opts ...LoadOption) error
}

// Execute executes a previously loaded OSImage.
//
// This will only work if OSImage.Load was called on some OSImage.
func Execute() error {
	return kexec.Reboot()
}
