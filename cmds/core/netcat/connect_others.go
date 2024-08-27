// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || windows

package main

import (
	"fmt"
	"net"

	"github.com/ishidawataru/sctp"
)

func connectToSCTPSocket(network, address string) (net.Conn, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve SCTP address: %w", err)
	}

	return sctp.DialSCTP(network, nil, sctpAddr)
}
