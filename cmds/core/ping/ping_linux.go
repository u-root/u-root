// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

func setupICMPv6Socket(c *net.IPConn) error {
	file, err := c.File()
	if err != nil {
		return fmt.Errorf("net.IPConn.File failed: %w", err)
	}
	// we want the stack to return us the network error if any occurred
	if err := unix.SetsockoptInt(int(file.Fd()), unix.SOL_IPV6, unix.IPV6_RECVERR, 1); err != nil {
		return fmt.Errorf("failed to set sock opt IPV6_RECVERR: %w", err)
	}
	return nil
}
