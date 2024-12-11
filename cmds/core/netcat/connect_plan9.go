// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"net"
)

func connectToSCTPSocket(network, address string) (net.Conn, error) {
	return nil, fmt.Errorf("sctp is not supported on Plan 9")
}
