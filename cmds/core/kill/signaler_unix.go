// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

var defaultSignal = "-SIGTERM"

func kill(sig os.Signal, pids ...string) []error {
	var errs []error
	s := sig.(syscall.Signal)
	for _, p := range pids {
		pid, err := strconv.Atoi(p)
		if err != nil {
			errs = append(errs, fmt.Errorf("%v: arguments must be process or job IDS", p))
			continue
		}
		if err := syscall.Kill(pid, s); err != nil {
			errs = append(errs, err)
		}

	}
	return errs
}
