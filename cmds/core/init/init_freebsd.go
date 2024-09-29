// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"

	"github.com/u-root/u-root/pkg/libinit"
)

func quiet() {
}

func osInitGo() *initCmds {
	// TOOD: get kernel command line.
	uinitArgs := libinit.WithArguments()

	// namespace setup will have done bind mounts into /bin, so name things only once.
	return &initCmds{
		cmds: []*exec.Cmd{
			// inito is (optionally) created by the u-root command when the
			// u-root initramfs is merged with an existing initramfs that
			// has a /init. The name inito means "original /init" There may
			// be an inito if we are building on an existing initramfs. All
			// initos need their own pid space.
			libinit.Command("/inito"),
			libinit.Command("/bin/uinit", uinitArgs),
			libinit.Command("/bin/sh"),
		},
	}
}
