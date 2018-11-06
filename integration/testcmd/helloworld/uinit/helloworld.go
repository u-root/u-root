// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// The most trivial init script.
func main() {
	fmt.Println("HELLO WORLD")

	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
