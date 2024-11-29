// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"golang.org/x/sys/windows"
)

func setDate(_ string, _ *time.Location, _ Clock) error {
	return windows.ERROR_NOT_SUPPORTED
}
