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
	"fmt"

	"golang.org/x/net/bpf"
)

type labelID int32

const (
	labelInvalid labelID = -1
	labelKeep    labelID = -2
	labelFail    labelID = -3
)

type instKind uint8

const (
	instPlain instKind = iota
	instJump
	instJumpIf
	instJumpIfX
)

type pendingInst struct {
	inst    bpf.Instruction
	kind    instKind
	targetT labelID
	targetF labelID
}

// prog is a two-pass cBPF assembler: emit records instructions with symbolic
// label targets; finalize resolves labels to relative skip offsets.
type prog struct {
	layout linkLayout
	target CompileTarget
	insts  []pendingInst
	labels map[labelID]int
	next   labelID
	done   bool
}

func newProg(layout linkLayout, target CompileTarget) *prog {
	return &prog{layout: layout, target: target, labels: make(map[labelID]int)}
}

func (p *prog) newLabel() labelID {
	l := p.next
	p.next++
	return l
}

func (p *prog) bind(l labelID) {
	if l == labelKeep || l == labelFail || l == labelInvalid {
		panic(fmt.Sprintf("prog: bind of reserved label %d", l))
	}
	if _, ok := p.labels[l]; ok {
		panic(fmt.Sprintf("prog: duplicate bind for label %d", l))
	}
	p.labels[l] = len(p.insts)
}

func (p *prog) emit(in ...bpf.Instruction) {
	for _, v := range in {
		switch v.(type) {
		case bpf.Jump, bpf.JumpIf, bpf.JumpIfX:
			panic("prog: emit received jump instruction; use emitJump/emitJumpIf/emitJumpIfX")
		}
		p.insts = append(p.insts, pendingInst{inst: v, kind: instPlain, targetT: labelInvalid, targetF: labelInvalid})
	}
}

func (p *prog) emitJumpIf(cond bpf.JumpTest, val uint32, tTrue, tFalse labelID) {
	if tTrue == labelInvalid || tFalse == labelInvalid {
		panic("prog: emitJumpIf target cannot be labelInvalid")
	}
	p.insts = append(p.insts, pendingInst{
		inst: bpf.JumpIf{Cond: cond, Val: val}, kind: instJumpIf, targetT: tTrue, targetF: tFalse,
	})
}

func (p *prog) emitJumpIfX(cond bpf.JumpTest, tTrue, tFalse labelID) {
	if tTrue == labelInvalid || tFalse == labelInvalid {
		panic("prog: emitJumpIfX target cannot be labelInvalid")
	}
	p.insts = append(p.insts, pendingInst{
		inst:    bpf.JumpIfX{Cond: cond},
		kind:    instJumpIfX,
		targetT: tTrue,
		targetF: tFalse,
	})
}

func (p *prog) emitJump(target labelID) {
	if target == labelInvalid {
		panic("prog: emitJump target cannot be labelInvalid")
	}
	p.insts = append(p.insts, pendingInst{
		inst: bpf.Jump{}, kind: instJump, targetT: target, targetF: labelInvalid,
	})
}

type (
	invalidLabelError struct {
		idx   int
		field string
	}
	unboundLabelError struct {
		idx   int
		label labelID
	}
	backwardJumpError struct {
		idx       int
		target    labelID
		targetIdx int
	}
	skipOverflowError struct {
		idx    int
		target labelID
		skip   int
	}
	alreadyFinalizedError struct{}
)

func (e invalidLabelError) Error() string {
	return fmt.Sprintf("prog: inst %d: %s references labelInvalid", e.idx, e.field)
}

func (e unboundLabelError) Error() string {
	return fmt.Sprintf("prog: inst %d: unbound label %d", e.idx, e.label)
}

func (e backwardJumpError) Error() string {
	return fmt.Sprintf("prog: inst %d: backward jump to label %d (idx %d)", e.idx, e.target, e.targetIdx)
}

func (e skipOverflowError) Error() string {
	return fmt.Sprintf("prog: inst %d: skip %d to label %d exceeds uint8", e.idx, e.skip, e.target)
}
func (e alreadyFinalizedError) Error() string { return "prog: finalize already called" }
func (p *prog) finalize() ([]bpf.Instruction, error) {
	if p.done {
		return nil, alreadyFinalizedError{}
	}
	p.done = true
	nKeep := len(p.insts)
	p.insts = append(p.insts,
		pendingInst{inst: returnKeep, kind: instPlain},
		pendingInst{inst: returnDrop, kind: instPlain})
	p.labels[labelKeep] = nKeep
	p.labels[labelFail] = nKeep + 1
	labelAt := make(map[labelID]int, len(p.labels))
	for l, idx := range p.labels {
		labelAt[l] = idx
	}
	resolve := func(target labelID, fromIdx int) (uint8, error) {
		if target == labelInvalid {
			return 0, invalidLabelError{idx: fromIdx, field: "target"}
		}
		targetIdx, ok := labelAt[target]
		if !ok {
			return 0, unboundLabelError{idx: fromIdx, label: target}
		}
		if targetIdx <= fromIdx {
			return 0, backwardJumpError{idx: fromIdx, target: target, targetIdx: targetIdx}
		}
		diff := targetIdx - fromIdx - 1
		if diff > 255 {
			return 0, skipOverflowError{idx: fromIdx, target: target, skip: diff}
		}
		return uint8(diff), nil
	}
	for i, pi := range p.insts[:nKeep] {
		switch pi.kind {
		case instPlain:
			continue
		case instJump:
			skip, err := resolve(pi.targetT, i)
			if err != nil {
				return nil, err
			}
			p.insts[i].inst = bpf.Jump{Skip: uint32(skip)}
		case instJumpIf:
			st, err := resolve(pi.targetT, i)
			if err != nil {
				return nil, err
			}
			sf, err := resolve(pi.targetF, i)
			if err != nil {
				return nil, err
			}
			j := p.insts[i].inst.(bpf.JumpIf)
			j.SkipTrue, j.SkipFalse = st, sf
			p.insts[i].inst = j
		case instJumpIfX:
			st, err := resolve(pi.targetT, i)
			if err != nil {
				return nil, err
			}
			sf, err := resolve(pi.targetF, i)
			if err != nil {
				return nil, err
			}
			j := p.insts[i].inst.(bpf.JumpIfX)
			j.SkipTrue, j.SkipFalse = st, sf
			p.insts[i].inst = j
		}
	}
	out := make([]bpf.Instruction, len(p.insts))
	for i, pi := range p.insts {
		out[i] = pi.inst
	}
	return out, nil
}
