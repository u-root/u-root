package tc

import (
	"errors"
	"log"
)

// Various errors
var (
	ErrNoArgAlter = errors.New("argument cannot be altered")
	ErrInvalidDev = errors.New("invalid device ID")

	// ErrNotImplemented is returned for not yet implemented parts.
	ErrNotImplemented = errors.New("functionality not yet implemented")

	// ErrNoArg is returned for missing arguments.
	ErrNoArg = errors.New("missing argument")

	// ErrInvalidArg is returned on invalid given arguments.
	ErrInvalidArg = errors.New("invalid argument")

	// ErrUnknownKind is returned for unknown qdisc, filter or class types.
	ErrUnknownKind = errors.New("unknown kind")
)

// Config contains options for RTNETLINK
type Config struct {
	// NetNS defines the network namespace
	NetNS int

	// Interface to log internals
	Logger *log.Logger
}

// Constants to define the direction
const (
	HandleRoot    uint32 = 0xFFFFFFFF
	HandleIngress uint32 = 0xFFFFFFF1

	HandleMinPriority uint32 = 0xFFE0
	HandleMinIngress  uint32 = 0xFFF2
	HandleMinEgress   uint32 = 0xFFF3

	// To alter filter in shared blocks, set Msg.Ifindex to MagicBlock
	MagicBlock = 0xFFFFFFFF
)

// Common flags from include/uapi/linux/pkt_cls.h
const (
	// don't offload filter to HW
	SkipHw uint32 = 1 << iota
	// don't use filter in SW
	SkipSw
	// filter is offloaded to HW
	InHw
	// filter isn't offloaded to HW
	NotInHw
	// verbose logging
	Verbose
)

const (
	// mask to differentiate between classes, qdiscs and filters
	actionMask  = 0x3c
	actionQdisc = 0x24
)
