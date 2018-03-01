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

type UserInputMessage struct {
	args []string
}

type StatusMessage struct {
	essid string
}

type WifiOption struct {
	Essid     string
	AuthSuite SecProto
}

var (
	UserInputChannel = make(chan UserInputMessage)
	StatusChannel    = make(chan StatusMessage)
)
