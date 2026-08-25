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
	"math"

	"golang.org/x/net/bpf"
)

type arithmeticExpression interface {
	emit(emitter *arithmeticEmitter, onMiss labelID)
	equal(other arithmeticExpression) bool
}

type arithmeticConstant struct {
	value uint32
}

func (e arithmeticConstant) emit(emitter *arithmeticEmitter, _ labelID) {
	emitter.prog.emit(bpf.LoadConstant{Dst: bpf.RegA, Val: e.value})
}

func (e arithmeticConstant) equal(other arithmeticExpression) bool {
	value, ok := other.(arithmeticConstant)
	return ok && e.value == value.value
}

type packetLength struct{}

func (packetLength) emit(emitter *arithmeticEmitter, _ labelID) {
	emitter.prog.emit(bpf.LoadExtension{Num: bpf.ExtLen})
}

func (packetLength) equal(other arithmeticExpression) bool {
	_, ok := other.(packetLength)
	return ok
}

type arithmeticBinary struct {
	left  arithmeticExpression
	op    bpf.ALUOp
	right arithmeticExpression
}

func (e arithmeticBinary) emit(emitter *arithmeticEmitter, onMiss labelID) {
	slot := emitter.acquire()
	e.left.emit(emitter, onMiss)
	emitter.prog.emit(bpf.StoreScratch{Src: bpf.RegA, N: slot})
	e.right.emit(emitter, onMiss)
	emitter.prog.emit(bpf.TAX{})
	if e.op == bpf.ALUOpDiv || e.op == bpf.ALUOpMod {
		nonzero := emitter.prog.newLabel()
		emitter.prog.emit(bpf.TXA{})
		emitter.prog.emitJumpIf(bpf.JumpEqual, 0, onMiss, nonzero)
		emitter.prog.bind(nonzero)
	}
	emitter.prog.emit(
		bpf.LoadScratch{Dst: bpf.RegA, N: slot},
		bpf.ALUOpX{Op: e.op},
	)
	emitter.release(slot)
}

func (e arithmeticBinary) equal(other arithmeticExpression) bool {
	value, ok := other.(arithmeticBinary)
	return ok && e.op == value.op && e.left.equal(value.left) && e.right.equal(value.right)
}

type packetProtocol uint8

const (
	packetProtocolEther packetProtocol = iota
	packetProtocolIP
	packetProtocolIP6
	packetProtocolTCP
	packetProtocolUDP
)

type packetAccess struct {
	protocol packetProtocol
	offset   arithmeticExpression
	size     int
}

func (e packetAccess) emit(emitter *arithmeticEmitter, onMiss labelID) {
	switch e.protocol {
	case packetProtocolEther:
		if !emitter.cursor.hasEther {
			emitter.prog.emitJump(onMiss)
			return
		}
		emitter.emitLoad(e.offset, 0, e.size, onMiss)
	case packetProtocolIP:
		emitter.emitNetworkLoad(e.offset, e.size, false, onMiss)
	case packetProtocolIP6:
		emitter.emitNetworkLoad(e.offset, e.size, true, onMiss)
	case packetProtocolTCP:
		emitter.emitTransportLoad(e.offset, e.size, ipProtocolTCP, onMiss)
	case packetProtocolUDP:
		emitter.emitTransportLoad(e.offset, e.size, ipProtocolUDP, onMiss)
	}
}

func (e packetAccess) equal(other arithmeticExpression) bool {
	value, ok := other.(packetAccess)
	return ok && e.protocol == value.protocol && e.size == value.size && e.offset.equal(value.offset)
}

type arithmeticComparison struct {
	left  arithmeticExpression
	test  bpf.JumpTest
	right arithmeticExpression
}

func (c arithmeticComparison) Compile(layout linkLayout) ([]bpf.Instruction, error) {
	return compileFilter(c, layout)
}

func (c arithmeticComparison) Equal(other Filter) bool {
	value, ok := other.(arithmeticComparison)
	return ok && c.test == value.test && c.left.equal(value.left) && c.right.equal(value.right)
}

func (c arithmeticComparison) Size(layout linkLayout) uint8 {
	instructions, err := c.Compile(layout)
	if err != nil || len(instructions) > math.MaxUint8 {
		return math.MaxUint8
	}
	return uint8(len(instructions))
}

func (c arithmeticComparison) IsPrimitive() bool { return true }
func (c arithmeticComparison) Type() ElementType { return Primitive }
func (c arithmeticComparison) Distill() Filter   { return c }

func (c arithmeticComparison) emit(b *prog, onMatch, onMiss labelID) {
	emitter := arithmeticEmitter{prog: b, cursor: cursorFor(b.layout)}
	c.emitPredicate(&emitter, onMatch, onMiss)
}

func (c arithmeticComparison) emitPredicate(
	emitter *arithmeticEmitter,
	onMatch labelID,
	onMiss labelID,
) {
	slot := emitter.acquire()
	c.left.emit(emitter, onMiss)
	emitter.prog.emit(bpf.StoreScratch{Src: bpf.RegA, N: slot})
	c.right.emit(emitter, onMiss)
	emitter.prog.emit(
		bpf.TAX{},
		bpf.LoadScratch{Dst: bpf.RegA, N: slot},
	)
	emitter.release(slot)
	emitter.prog.emitJumpIfX(c.test, onMatch, onMiss)
}

func (c arithmeticComparison) isAlwaysReject(layout linkLayout) bool {
	if layout.hasL2Protocols() {
		return false
	}
	return arithmeticRequiresEther(c.left) || arithmeticRequiresEther(c.right)
}
func (arithmeticComparison) isAlwaysAccept(linkLayout) bool { return false }

func arithmeticRequiresEther(expression arithmeticExpression) bool {
	switch value := expression.(type) {
	case arithmeticBinary:
		return arithmeticRequiresEther(value.left) || arithmeticRequiresEther(value.right)
	case packetAccess:
		return value.protocol == packetProtocolEther || arithmeticRequiresEther(value.offset)
	default:
		return false
	}
}

func emitComparisonWithCursor(
	comparison arithmeticComparison,
	b *prog,
	cursor packetCursor,
	branches emitBranches,
) {
	matched := b.newLabel()
	missed := b.newLabel()
	b.layout = cursor.layout()
	usesInnerNetwork := arithmeticUsesInnerNetwork(comparison.left) ||
		arithmeticUsesInnerNetwork(comparison.right)
	if cursor.mplsDepth > 0 && usesInnerNetwork {
		bottomOfStack := b.newLabel()
		b.emit(bpf.LoadAbsolute{Off: cursor.l3Offset - 2, Size: lengthByte})
		b.emitJumpIf(bpf.JumpBitsSet, 0x01, bottomOfStack, missed)
		b.bind(bottomOfStack)
	}
	emitter := arithmeticEmitter{prog: b, cursor: cursor}
	comparison.emitPredicate(&emitter, matched, missed)
	b.bind(matched)
	branches.onMatch(b, cursor)
	b.bind(missed)
	branches.onMiss(b, cursor)
}

func arithmeticUsesInnerNetwork(expression arithmeticExpression) bool {
	switch value := expression.(type) {
	case arithmeticBinary:
		return arithmeticUsesInnerNetwork(value.left) || arithmeticUsesInnerNetwork(value.right)
	case packetAccess:
		usesNetwork := value.protocol != packetProtocolEther
		return usesNetwork || arithmeticUsesInnerNetwork(value.offset)
	default:
		return false
	}
}

type arithmeticEmitter struct {
	prog        *prog
	cursor      packetCursor
	nextScratch int
}

func (e *arithmeticEmitter) acquire() int {
	if e.nextScratch >= 16 {
		panic("filter: arithmetic expression requires more than 16 scratch slots")
	}
	slot := e.nextScratch
	e.nextScratch++
	return slot
}

func (e *arithmeticEmitter) release(slot int) {
	if slot != e.nextScratch-1 {
		panic("filter: arithmetic scratch slots released out of order")
	}
	e.nextScratch--
}

func (e *arithmeticEmitter) emitNetworkLoad(
	offset arithmeticExpression,
	size int,
	ipv6 bool,
	onMiss labelID,
) {
	e.prog.layout = e.cursor.layout()
	e.prog.loadEtherKind()
	protocolOK := e.prog.newLabel()
	if ipv6 {
		e.prog.compareProtocolIP6(protocolOK, onMiss)
	} else {
		e.prog.compareProtocolIP4(protocolOK, onMiss)
	}
	e.prog.bind(protocolOK)
	e.emitLoad(offset, e.cursor.l3Offset, size, onMiss)
}

func (e *arithmeticEmitter) emitTransportLoad(
	offset arithmeticExpression,
	size int,
	protocol uint32,
	onMiss labelID,
) {
	e.prog.layout = e.cursor.layout()
	e.prog.loadEtherKind()
	ipv6 := e.prog.newLabel()
	ipv4 := e.prog.newLabel()
	finished := e.prog.newLabel()
	e.prog.compareProtocolIP6(ipv6, ipv4)

	e.prog.bind(ipv6)
	ipv6ProtocolOK := e.prog.newLabel()
	e.prog.emit(loadIPv6Protocol(e.prog.layout))
	e.prog.emitJumpIf(bpf.JumpEqual, protocol, ipv6ProtocolOK, onMiss)
	e.prog.bind(ipv6ProtocolOK)
	e.emitLoad(offset, e.cursor.l3Offset+40, size, onMiss)
	e.prog.emitJump(finished)

	e.prog.bind(ipv4)
	ipv4OK := e.prog.newLabel()
	e.prog.compareProtocolIP4(ipv4OK, onMiss)
	e.prog.bind(ipv4OK)
	ipv4ProtocolOK := e.prog.newLabel()
	e.prog.emit(loadIPv4Protocol(e.prog.layout))
	e.prog.emitJumpIf(bpf.JumpEqual, protocol, ipv4ProtocolOK, onMiss)
	e.prog.bind(ipv4ProtocolOK)
	fragmentOK := e.prog.newLabel()
	e.prog.emit(bpf.LoadAbsolute{
		Off:  e.cursor.l3Offset + intraIP4HeaderFlags,
		Size: lengthHalf,
	})
	e.prog.emitJumpIf(bpf.JumpBitsSet, jumpMask, onMiss, fragmentOK)
	e.prog.bind(fragmentOK)
	e.emitIPv4TransportLoad(offset, size, onMiss)
	e.prog.emitJump(finished)

	e.prog.bind(finished)
}

func (e *arithmeticEmitter) emitIPv4TransportLoad(
	offset arithmeticExpression,
	size int,
	onMiss labelID,
) {
	slot := e.acquire()
	offset.emit(e, onMiss)
	maximum := uint32(math.MaxUint32) - e.cursor.l3Offset - 60 - uint32(size)
	valid := e.prog.newLabel()
	e.prog.emitJumpIf(bpf.JumpGreaterThan, maximum, onMiss, valid)
	e.prog.bind(valid)
	e.prog.emit(bpf.StoreScratch{Src: bpf.RegA, N: slot})
	e.prog.emit(bpf.LoadMemShift{Off: e.cursor.l3Offset + intraIP4HeaderSize})
	e.prog.emit(
		bpf.TXA{},
		bpf.LoadScratch{Dst: bpf.RegX, N: slot},
		bpf.ALUOpX{Op: bpf.ALUOpAdd},
		bpf.ALUOpConstant{Op: bpf.ALUOpAdd, Val: e.cursor.l3Offset},
	)
	e.emitIndirectAtA(size, onMiss)
	e.release(slot)
}

func (e *arithmeticEmitter) emitLoad(
	offset arithmeticExpression,
	base uint32,
	size int,
	onMiss labelID,
) {
	if constant, ok := offset.(arithmeticConstant); ok {
		if constant.value > math.MaxUint32-base-uint32(size) {
			e.prog.emitJump(onMiss)
			return
		}
		absolute := base + constant.value
		inBounds := e.prog.newLabel()
		e.prog.emit(bpf.LoadExtension{Num: bpf.ExtLen})
		e.prog.emitJumpIf(bpf.JumpGreaterOrEqual, absolute+uint32(size), inBounds, onMiss)
		e.prog.bind(inBounds)
		e.prog.emit(bpf.LoadAbsolute{Off: absolute, Size: size})
		return
	}

	offset.emit(e, onMiss)
	maximum := uint32(math.MaxUint32) - base - uint32(size)
	valid := e.prog.newLabel()
	e.prog.emitJumpIf(bpf.JumpGreaterThan, maximum, onMiss, valid)
	e.prog.bind(valid)
	e.prog.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAdd, Val: base})
	e.emitIndirectAtA(size, onMiss)
}

func (e *arithmeticEmitter) emitIndirectAtA(size int, onMiss labelID) {
	offsetSlot := e.acquire()
	e.prog.emit(bpf.StoreScratch{Src: bpf.RegA, N: offsetSlot})
	e.prog.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAdd, Val: uint32(size)})
	e.prog.emit(bpf.TAX{})
	e.prog.emit(bpf.LoadExtension{Num: bpf.ExtLen})
	inBounds := e.prog.newLabel()
	e.prog.emitJumpIfX(bpf.JumpGreaterOrEqual, inBounds, onMiss)
	e.prog.bind(inBounds)
	e.prog.emit(
		bpf.LoadScratch{Dst: bpf.RegX, N: offsetSlot},
		bpf.LoadIndirect{Off: 0, Size: size},
	)
	e.release(offsetSlot)
}

func validateArithmeticDepth(expression arithmeticExpression, depth int) error {
	if depth >= 15 {
		return fmt.Errorf("arithmetic expression is too deeply nested")
	}
	switch value := expression.(type) {
	case arithmeticBinary:
		if constant, ok := value.right.(arithmeticConstant); ok && constant.value == 0 &&
			(value.op == bpf.ALUOpDiv || value.op == bpf.ALUOpMod) {
			return fmt.Errorf("division by zero")
		}
		if err := validateArithmeticDepth(value.left, depth+1); err != nil {
			return err
		}
		return validateArithmeticDepth(value.right, depth+1)
	case packetAccess:
		return validateArithmeticDepth(value.offset, depth+1)
	default:
		return nil
	}
}

func arithmeticScratch(expression arithmeticExpression) int {
	switch value := expression.(type) {
	case arithmeticBinary:
		return 1 + max(arithmeticScratch(value.left), arithmeticScratch(value.right))
	case packetAccess:
		return 2 + arithmeticScratch(value.offset)
	default:
		return 0
	}
}
