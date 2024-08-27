// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health

import (
	"github.com/jaypipes/ghw"
)

// Kernel should hold all stats for the kernel.
// If needed, later, the tag can be a command as well.
// Since it is returned as JSON, it is no problem
// to change it over time.
type Kernel struct {
	Version string `file:"/proc/version"`
	Modules string `file:"/proc/modules"`
	Drivers string `file:"/proc/devices"`
}

type Stat struct {
	Hostname string
	Info     *ghw.HostInfo
	Kernel   Kernel
	// This should be empty, but some packages behave badly.
	Stderr string
	Err    string
}

type GatheredStats struct {
	Stats []Stat
}

type NodeList struct {
	List []string
	cmd  string
	args []string
}

type Gather struct {
	Nodes *NodeList
	cmd   string
	args  []string
}
