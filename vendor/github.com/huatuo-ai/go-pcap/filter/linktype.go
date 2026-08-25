// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"errors"
	"fmt"

	"golang.org/x/net/bpf"
)

// LinkType enumerates supported data-link types.
// LinkTypeEthernet is the zero value and serves as the default.
type LinkType uint32

const (
	LinkTypeEthernet LinkType = iota // DLT_EN10MB
	LinkTypeRaw                      // DLT_RAW
)

// CompileTarget selects environment-specific cBPF lowering.
type CompileTarget uint8

const (
	// CompileTargetPortable emits packet-byte-only cBPF suitable for offline
	// matching and pcap_open_dead style compilation.
	CompileTargetPortable CompileTarget = iota
	// CompileTargetLinuxSocket may emit Linux socket BPF extensions such as
	// VLAN auxiliary metadata loads accepted by SO_ATTACH_FILTER.
	CompileTargetLinuxSocket
)

// Sentinel errors.
var (
	ErrEmptyFilter         = errors.New("filter: empty expression")
	ErrInvalidFilter       = errors.New("filter: invalid expression")
	ErrUnsupportedFeature  = errors.New("filter: unsupported feature")
	ErrUnsupportedLinkType = errors.New("filter: unsupported link type")
	ErrL2OnlyLinkType      = errors.New("filter: expression matches only L2 protocols on non-L2 link type")
	ErrHostResolution      = errors.New("filter: host resolution failed")
)

// CompileOptions supplies the packet layout and environment-dependent inputs
// used while compiling an expression. A nil Resolver uses net.DefaultResolver.
type CompileOptions struct {
	LinkType LinkType
	Resolver Resolver
	Target   CompileTarget
}

// linkLayout abstracts per-link-type packet layout.
// Concrete implementations are ethernetLayout and rawLayout, selected
// by layoutFor through an explicit switch (closed set, no registry).
type linkLayout interface {
	// l3Off returns the byte offset from packet start to the L3 (IP) header.
	l3Off() uint32

	// genLinkProbe returns instructions that load the link-layer type into A.
	genLinkProbe() []bpf.Instruction

	// genLinkType returns a JumpIf comparing A against the link-type value
	// for proto. st/sf are relative skip offsets.
	genLinkType(proto uint32, st, sf uint8) bpf.Instruction

	// linkProbeSize returns len(genLinkProbe()).
	linkProbeSize() uint8

	// linkCompareSize returns 1 (the JumpIf from genLinkType).
	linkCompareSize() uint8

	// hasL2Protocols reports whether this link type carries L2 protocols
	// (arp, rarp, ether).
	hasL2Protocols() bool
}

// ethernetLayout implements linkLayout for DLT_EN10MB.
type ethernetLayout struct{}

func (e ethernetLayout) l3Off() uint32        { return 14 }
func (e ethernetLayout) linkProbeSize() uint8 { return 1 }
func (e ethernetLayout) hasL2Protocols() bool { return true }

func (e ethernetLayout) genLinkProbe() []bpf.Instruction {
	return []bpf.Instruction{bpf.LoadAbsolute{Off: 12, Size: 2}}
}

func (e ethernetLayout) linkCompareSize() uint8 { return 1 }

func (e ethernetLayout) genLinkType(proto uint32, st, sf uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: proto, SkipTrue: st, SkipFalse: sf}
}

// rawLayout implements linkLayout for DLT_RAW.
type rawLayout struct{}

func (r rawLayout) l3Off() uint32          { return 0 }
func (r rawLayout) linkProbeSize() uint8   { return 2 }
func (r rawLayout) linkCompareSize() uint8 { return 1 }
func (r rawLayout) hasL2Protocols() bool   { return false }

func (r rawLayout) genLinkProbe() []bpf.Instruction {
	return []bpf.Instruction{
		bpf.LoadAbsolute{Off: 0, Size: 1},
		bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0xF0},
	}
}

func (r rawLayout) genLinkType(proto uint32, st, sf uint8) bpf.Instruction {
	// Map ethertype to IP version nibble (libpcap DLT_RAW behavior).
	// L2-only protocols map to 0xF0 which never matches.
	var versionNibble uint32
	switch proto {
	case etherTypeIPv4:
		versionNibble = 0x40
	case etherTypeIPv6:
		versionNibble = 0x60
	default:
		versionNibble = 0xF0
	}
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: versionNibble, SkipTrue: st, SkipFalse: sf}
}

// layoutFor returns the linkLayout for lt.
func layoutFor(lt LinkType) (linkLayout, error) {
	switch lt {
	case LinkTypeEthernet:
		return ethernetLayout{}, nil
	case LinkTypeRaw:
		return rawLayout{}, nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedLinkType, lt)
	}
}

// Compile parses a tcpdump filter expression and compiles it to cBPF
// instructions for the given link type.
func Compile(expr string, lt LinkType) ([]bpf.Instruction, error) {
	return CompileWithOptions(expr, CompileOptions{LinkType: lt})
}

// CompileWithOptions parses and compiles a tcpdump filter expression.
func CompileWithOptions(expr string, options CompileOptions) ([]bpf.Instruction, error) {
	if expr == "" {
		return nil, ErrEmptyFilter
	}
	layout, err := layoutFor(options.LinkType)
	if err != nil {
		return nil, err
	}
	f, err := parseFilterExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFilter, err)
	}
	prepared, err := prepareFilter(f, options.Resolver)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFilter, err)
	}
	insns, err := compilePreparedFilter(prepared, layout, options.Target)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFilter, err)
	}
	// If the entire compiled filter on a non-L2 link type is just
	// 2 instructions (returnDrop, returnKeep), the whole expression
	// matches only L2 protocols. Return ErrL2OnlyLinkType so callers
	// can substitute a reject-all stub.
	if !layout.hasL2Protocols() && len(insns) == 2 {
		if insns[0] == returnDrop && insns[1] == returnKeep {
			return nil, ErrL2OnlyLinkType
		}
	}
	return insns, nil
}

// Size returns the instruction count that Compile would emit for expr+lt.
func Size(expr string, lt LinkType) (int, error) {
	return SizeWithOptions(expr, CompileOptions{LinkType: lt})
}

// SizeWithOptions returns the instruction count CompileWithOptions emits.
func SizeWithOptions(expr string, options CompileOptions) (int, error) {
	insns, err := CompileWithOptions(expr, options)
	if err != nil {
		return 0, err
	}
	return len(insns), nil
}
