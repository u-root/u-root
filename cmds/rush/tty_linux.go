// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/signal"
)

const (
	sysIoctl = 16
)

// tty does whatever needs to be done to set up a tty for GOOS.
func tty() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for range sigs {
			fmt.Println("")
		}
	}()
}
