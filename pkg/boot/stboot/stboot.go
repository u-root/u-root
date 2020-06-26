// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"net"
	"strings"
)

const (
	// BootballExt is the file extension of bootballs
	BootballExt string = ".stboot"
	// DefaultBallName is the file name of the archive, which is expected to contain
	// the stboot configuration file along with the corresponding files
	DefaultBallName string = "ball.stboot"
	// ConfigName is the name of the stboot configuration file
	ConfigName string = "stconfig.json"
)

// ComposeIndividualBallPrefix returns a host specific name prefix for bootball files.
func ComposeIndividualBallPrefix(hwAddr net.HardwareAddr) string {
	prefix := hwAddr.String()
	prefix = strings.ReplaceAll(prefix, ":", "-")
	return prefix + "-"
}
