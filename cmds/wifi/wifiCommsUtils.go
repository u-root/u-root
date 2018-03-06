// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type SecProto int

const (
	NoEnc SecProto = iota
	WpaPsk
	WpaEap
	NotSupportedProto
)

type WifiOption struct {
	Essid     string
	AuthSuite SecProto
}

type State struct {
	NearbyWifis     []WifiOption
	ConnectingEssid string
	CurEssid        string
}

type ConnectReqChanMsg struct {
	c     chan (error)
	essid string
}

var (
	// Assumption: The user shouldn't "accidentally" try to connect more than 4 times
	ConnectReqChan = make(chan ConnectReqChanMsg, 4)
)
