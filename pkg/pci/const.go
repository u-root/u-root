// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

// Constants for Linux kernel access.
const (
	// StdConfigSize is a PCI name -- standard?
	StdConfigSize = 64
	// ConfigSize is the pre-PCIE config size
	ConfigSize = 256
	// FullConfigSize is the PCIE 4k config size.
	FullConfigSize = 4096
)

// Values defining config space.
const (
	StdNumBARS = 6
)

// Config space registers
const (
	VID = 0
	DID = 2

	Cmd            = 4
	CmdIO          = 1
	CmdMem         = 2
	CmdBME         = 4
	CmdSpecial     = 8
	CmdMWINV       = 0x10
	CmdVGA         = 0x20
	CmdParity      = 0x40
	CmdStep        = 0x80
	CmdSERR        = 0x100
	CmdFastB2B     = 0x200
	CmdINTXDisable = 0x400

	ClassRevision = 8
	RevisionID    = 8
	ClassProg     = 9
	ClassDevice   = 10

	CacheLineSize = 0xc

	LatencyTimer = 0xd

	HeaderType       = 0xe
	HeaderTypeMask   = 0x7f
	HeaderTypeNormal = 0
	HeaderTypeBridge = 1

	BAR0 = 0x10
	BAR1 = 0x14
	BAR2 = 0x18
	BAR3 = 0x1c
	BAR4 = 0x20
	BAR5 = 0x24

	// The low 3 bits tell you what type of space it is.
	BARTypeMask = 7
	BARMem32    = 0
	BARIO       = 1
	BARMem64    = 4
	BARPrefetch = 8
	BARMemMask  = ^0xf
	BARIOMask   = ^3

	// Type 0 devices
	SubSystemVID   = 0x2c
	SubSystemID    = 0x2e
	ROMAddress     = 0x30
	ROMEnabled     = 1
	ROMAddressMask = ^0x7ff

	IRQLine = 0x3c
	IRQPin  = 0x3d
	MinGnt  = 0x3e
	MaxLat  = 0x3f

	// Type 1
	Primary          = 0x18 // our bus
	Secondary        = 0x19 // first bus behind bridge
	Subordinate      = 0x1a // last bus behind bridge, inclusive
	SecondaryLatency = 0x1b
	IOBase           = 0x1c
	IOLimit          = 0x1d
	IOHighBase       = 0x30
	IOHighLimit      = 0x32
	IOTypeMask       = 0xf
	IOType16         = 0
	IOType32         = 1
	IORangeMask4k    = 0xf
	IORangeMaskIntel = 3 // Intel has ever been ready to "improve" PCI

	SecStatus = 0x1e // bit 14 only

	MemBase          = 0x20
	MemLimit         = 0x22
	MemTypeMask      = 0xf
	MemMask          = ^MemTypeMask
	PrefMemBase      = 0x24
	PrefMemLimit     = 0x26
	PrefMemTypeMask  = 0xf
	PrefMemType32    = 0
	PrefMemType64    = 1
	PrefMemHighBase  = 0x28
	PrefMemHighLimit = 0x30

	BridgeRomAddress = 0x38

	BridgeControl     = 0x3e
	BridgeParity      = 1
	BridgeSERR        = 2
	BridgeISA         = 4
	BridgeVGA         = 8
	BridgeMasterAbort = 0x20
	BridgeBusReset    = 0x40
	BridgeFastB2B     = 0x80
)
