// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import "github.com/u-root/u-root/pkg/netcat"

func init() {
	osListeners[netcat.SOCKET_TYPE_SCTP] = listenToSCTPSocket
}
