// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !tinygo

package libinit

import (
	"fmt"

	"github.com/u-root/u-root/pkg/ulog"
	"github.com/vishvananda/netlink"
)

// NetInit is u-root network initialization.
func linuxNetInit() {
	if err := loopbackUp(); err != nil {
		ulog.KernelLog.Printf("Failed to initialize loopback: %v", err)
	}
}

func loopbackUp() error {
	lo, err := netlink.LinkByName("lo")
	if err != nil {
		return err
	}

	if err := netlink.LinkSetUp(lo); err != nil {
		return fmt.Errorf("couldn't set link loopback up: %w", err)
	}
	return nil
}

func init() {
	osNetInit = linuxNetInit
}
