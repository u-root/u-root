// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && tinygo

package libinit

import (
	"github.com/u-root/u-root/pkg/ulog"
)

// Stub for notifying the user that network functionality is not supported on linux platforms
func linuxNetInit() {
	ulog.KernelLog.Printf("tinygo builds currently do not support network functionality on linux platforms\n")
}

func init() {
	osNetInit = linuxNetInit
}
