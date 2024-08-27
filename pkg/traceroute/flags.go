// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

type Flags struct {
	Host         string
	Proto        string
	ICMP         bool
	TCP          bool
	MaxHops      int
	DestPortSeq  uint
	TOS          int
	ProbesPerHop int
	Source       string
	Module       string
	UDP          bool
}

type Args struct {
	Host string
}
