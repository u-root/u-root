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
	"golang.org/x/net/bpf"
)

type composite struct {
	filters Filters
	and     bool
}

func (c composite) Compile(layout linkLayout) ([]bpf.Instruction, error) {
	return compileFilter(c, layout)
}

func (c composite) emit(b *prog, onMatch, onMiss labelID) {
	// Trailing always-accept (AND) / always-reject (OR) members are skipped via
	// continue below, so the static index len-1 is not necessarily the last
	// emitted filter; using it would leave the genuinely-last filter's
	// fall-through dangling onto the returnKeep sentinel instead of onMatch/onMiss.
	lastEmit := -1
	for i, f := range c.filters {
		if c.and && f.isAlwaysAccept(b.layout) {
			continue
		}
		if !c.and && f.isAlwaysReject(b.layout) {
			continue
		}
		lastEmit = i
	}
	if lastEmit < 0 {
		// All members are always-skip. Unreachable given the always-accept/
		// always-reject short-circuits in compileFilter, but jump explicitly
		// rather than rely on fall-through aliasing the returnKeep sentinel.
		if c.and {
			b.emitJump(onMatch)
		} else {
			b.emitJump(onMiss)
		}
		return
	}
	for i, f := range c.filters {
		if f.isAlwaysReject(b.layout) {
			if c.and {
				b.emitJump(onMiss)
				return
			}
			continue
		}
		if f.isAlwaysAccept(b.layout) {
			if c.and {
				continue
			}
			b.emitJump(onMatch)
			return
		}
		if i == lastEmit {
			f.emit(b, onMatch, onMiss)
		} else if c.and {
			next := b.newLabel()
			f.emit(b, next, onMiss)
			b.bind(next)
		} else {
			next := b.newLabel()
			f.emit(b, onMatch, next)
			b.bind(next)
		}
	}
}

func (c composite) isAlwaysReject(layout linkLayout) bool {
	if c.and {
		for _, f := range c.filters {
			if f.isAlwaysReject(layout) {
				return true
			}
		}
		return false
	}
	for _, f := range c.filters {
		if !f.isAlwaysReject(layout) {
			return false
		}
	}
	return true
}

func (c composite) isAlwaysAccept(layout linkLayout) bool {
	if c.and {
		for _, f := range c.filters {
			if !f.isAlwaysAccept(layout) {
				return false
			}
		}
		return true
	}
	for _, f := range c.filters {
		if f.isAlwaysAccept(layout) {
			return true
		}
	}
	return false
}

func (c composite) Equal(o Filter) bool {
	if o == nil {
		return false
	}
	oc, ok := o.(composite)
	if !ok {
		return false
	}
	return c.and == oc.and && c.filters.Equal(oc.filters)
}

func (c composite) Size(layout linkLayout) uint8 {
	if c.isAlwaysReject(layout) {
		return 2
	}
	if c.isAlwaysAccept(layout) {
		return 2
	}
	insns, err := c.Compile(layout)
	if err != nil || len(insns) == 0 {
		return 2
	}
	return uint8(len(insns))
}

func (c composite) IsPrimitive() bool { return false }
func (c composite) Type() ElementType { return Composite }

func (c composite) LastPrimitive() *primitive {
	if len(c.filters) == 0 {
		return nil
	}
	last := c.filters[len(c.filters)-1]
	switch v := last.(type) {
	case primitive:
		return &v
	case negated:
		if p, ok := v.inner.(primitive); ok {
			return &p
		}
	}
	return nil
}

func (c composite) Distill() Filter {
	list := make(Filters, 0)
	for _, f := range c.filters {
		list = append(list, f.Distill())
	}
	c.filters = list
	if len(c.filters) == 1 {
		return c.filters[0]
	}
	if !c.and {
		return c
	}
	prims := make(primitives, 0)
	others := make(Filters, 0)
	for _, f := range c.filters {
		// Only bare primitives are eligible for combine(); negated wrappers
		// (and nested composites) must survive intact, otherwise their
		// negation is lost because emit/isAlwaysReject ignore primitive.negator.
		if p, ok := f.(primitive); ok && p.kind != filterKindVLAN && p.kind != filterKindMPLS {
			prims = append(prims, p)
		} else {
			others = append(others, f)
		}
	}
	p2 := prims.combine()
	list = make(Filters, 0)
	for _, p := range *p2 {
		list = append(list, p)
	}
	list = append(list, others...)
	c.filters = list
	if len(c.filters) == 1 {
		return c.filters[0]
	}
	return c
}
