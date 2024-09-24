// Copyright 2021-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !tinygo

package watchdog

import (
	"fmt"
)

func Usage() {
	fmt.Print(`watchdogd run [--dev DEV] [--timeout N] [--pre_timeout N] [--keep_alive N] [--monitors STRING]
	Run the watchdogd daemon in a child process (does not daemonize).
watchdogd stop
	Send a signal to arm the running watchdogd.
watchdogd continue
	Send a signal to disarm the running watchdogd.
watchdogd arm
	Send a signal to arm the running watchdogd.
watchdogd disarm
	Send a signal to disarm the running watchdogd.
`)
}
