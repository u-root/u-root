// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && (linux || windows)

package main

import (
	"fmt"
	"net"

	"github.com/ishidawataru/sctp"
)

func listenToSCTPSocket(network, address string) (net.Listener, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve SCTP address: %w", err)
	}

	return sctp.ListenSCTP(network, sctpAddr)
}
