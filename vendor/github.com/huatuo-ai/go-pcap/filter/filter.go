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

// Package filter compiles tcpdump filter expressions to cBPF instructions.
//
// Derived from github.com/packetcap/go-pcap (Apache-2.0), extended with
// native link-type parameterization supporting EN10MB (Ethernet) and RAW (raw IP).
//
// DLT_RAW limitations: L2-only primitives (arp, rarp, ether, vlan, gateway, stp)
// are unsupported; when an entire filter matches only L2 protocols,
// ErrL2OnlyLinkType is returned so callers can substitute a reject-all stub.
package filter

import (
	"net"
	"sort"

	"golang.org/x/net/bpf"
)

// Filter is a compiled tcpdump filter expression. Compile and Size require
// the link layout so instruction offsets and counts match the target DLT.
type Filter interface {
	Compile(layout linkLayout) ([]bpf.Instruction, error)
	Equal(o Filter) bool
	Size(layout linkLayout) uint8
	IsPrimitive() bool
	Type() ElementType
	Distill() Filter
	emit(b *prog, onMatch, onMiss labelID)
	isAlwaysReject(layout linkLayout) bool
	isAlwaysAccept(layout linkLayout) bool
}

// compileFilter compiles a Filter using the new emit-based assembler. It
// handles the always-accept/always-reject short-circuits, then falls through
// to creating a prog and calling emit.
func compileFilter(f Filter, layout linkLayout) ([]bpf.Instruction, error) {
	prepared, err := prepareFilter(f, net.DefaultResolver)
	if err != nil {
		return nil, err
	}
	return compilePreparedFilter(prepared, layout, CompileTargetPortable)
}

func compilePreparedFilter(f Filter, layout linkLayout, target CompileTarget) ([]bpf.Instruction, error) {
	if f.isAlwaysReject(layout) {
		return []bpf.Instruction{returnDrop, returnKeep}, nil
	}
	if f.isAlwaysAccept(layout) {
		return []bpf.Instruction{returnKeep, returnDrop}, nil
	}
	if filterNeedsCursor(f) {
		return compileCursorFilter(f, layout, target)
	}
	b := newProg(layout, target)
	f.emit(b, labelKeep, labelFail)
	return b.finalize()
}

// negated wraps a Filter, swapping match/miss semantics. Compile delegates to
// compileFilter; emit swaps onMatch/onMiss; the always-accept/always-reject
// checks are also swapped.
type negated struct {
	inner Filter
}

func (n negated) Compile(layout linkLayout) ([]bpf.Instruction, error) {
	return compileFilter(n, layout)
}

func (n negated) emit(b *prog, onMatch, onMiss labelID) {
	n.inner.emit(b, onMiss, onMatch)
}

func (n negated) Equal(o Filter) bool {
	other, ok := o.(negated)
	if !ok {
		return false
	}
	return n.inner.Equal(other.inner)
}

func (n negated) Size(layout linkLayout) uint8 {
	if n.isAlwaysReject(layout) {
		return 2
	}
	if n.isAlwaysAccept(layout) {
		return 2
	}
	insns, err := n.Compile(layout)
	if err != nil || len(insns) == 0 {
		return 2
	}
	return uint8(len(insns))
}

func (n negated) IsPrimitive() bool {
	return n.inner.IsPrimitive()
}

func (n negated) Type() ElementType {
	return n.inner.Type()
}

func (n negated) Distill() Filter {
	d := n.inner.Distill()
	if _, ok := d.(negated); ok {
		// Double negation cancels out.
		return d.(negated).inner
	}
	return negated{inner: d}
}

func (n negated) isAlwaysReject(layout linkLayout) bool {
	return n.inner.isAlwaysAccept(layout)
}

func (n negated) isAlwaysAccept(layout linkLayout) bool {
	return n.inner.isAlwaysReject(layout)
}

type ElementType uint8

const (
	Primitive ElementType = iota
	Composite
	Joiner
)

type Element interface {
	Type() ElementType
}

type Filters []Filter

func (f Filters) Len() int {
	return len(f)
}

func (f Filters) Less(i, j int) bool {
	return false
}

func (f Filters) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Filters) Equal(o Filters) bool {
	// not matched if of the wrong length
	if len(f) != len(o) {
		return false
	}

	// copy so that our sort does not affect the original
	f1 := f[:]
	o1 := o[:]
	sort.Sort(f1)
	sort.Sort(o1)
	for i, val := range f1 {
		if !val.Equal(o1[i]) {
			return false
		}
	}
	return true
}
