// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

type Args struct {
	Host      string
	PacketLen int
}

type Flags struct {
	AF4          bool
	AF6          bool
	ICMP         bool
	TCP          bool
	Debug        bool
	DontFrag     bool
	NoRoute      bool
	Device       string
	FirstHop     int
	MaxHops      int
	SimProbes    int
	DestPortSeq  uint
	TOS          int
	Flowlabel    int
	HereFactor   int
	NearFactor   int
	WaitSeconds  int64
	SendSeconds  int64
	ProbesPerHop int
	Source       string
	Extensions   bool
	AsLookups    bool
	Module       string
	Opts         []string
	SourcePort   int
	UDP          bool
	UDPLite      bool
	DCCP         bool
	SetRaw       bool
	MTUDiscovery bool
	Backwards    bool
	NoResolve    bool
	Gateways     []string
}
