// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/u-root/u-root/pkg/wifi"

type State struct {
	NearbyWifis     []wifi.WifiOption
	ConnectingEssid string
	CurEssid        string
}

type ConnectReqChanMsg struct {
	c         chan (error) // channel that the requesting routine is listening on
	essid     string
	routineID []byte
	success   bool
}

type RefreshReqChanMsg struct {
	c chan (error)
}

var (
	// Assumption: The user shouldn't "accidentally" try to connect or refresh more than 4 times
	DefaultBufferSize = 4
	ConnectReqChan    = make(chan ConnectReqChanMsg, DefaultBufferSize)
	RefreshReqChan    = make(chan RefreshReqChanMsg, DefaultBufferSize)
)
