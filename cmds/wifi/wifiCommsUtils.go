// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type ServerToServiceMessage struct {
	essid string
	id    string
	pass  string
}

type ServiceToServerMessage struct {
	essid string
}

var (
	ServerToServiceChan = make(chan ServerToServiceMessage)
	ServiceToServerChan = make(chan ServiceToServerMessage)
)
