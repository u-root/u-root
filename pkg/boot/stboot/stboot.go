// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"net"
	"strings"
)

const (
	// BallName is the file name of the archive, which is expected to contain
	// the stboot configuration file along with the corresponding files
	BallName string = "stboot.ball"
	// ConfigName is the name of the stboot configuration file
	ConfigName string = "stconfig.json"
	//HostVarsName is the name of file containing host-specific data
	HostVarsName string = "hostvars.json"
)

// ComposeIndividualBallName extends the general BallName
// with an individual hardware address.
func ComposeIndividualBallName(hwAddr net.HardwareAddr) string {
	suffix := hwAddr.String()
	suffix = strings.ReplaceAll(suffix, ":", "-")
	return BallName + "." + suffix
}
