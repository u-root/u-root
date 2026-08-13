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
	"strconv"

	"golang.org/x/net/bpf"
)

type probeKind uint8

const (
	probeEtherType probeKind = iota
	probeIPVersion
)

type packetCursor struct {
	probe        probeKind
	etherTypeOff uint32
	l3Offset     uint32
	mplsDepth    uint16
	vlanDepth    uint16
	hasEther     bool
}

func cursorFor(layout linkLayout) packetCursor {
	if layout.hasL2Protocols() {
		return packetCursor{
			probe:        probeEtherType,
			etherTypeOff: 12,
			l3Offset:     layout.l3Off(),
			hasEther:     true,
		}
	}
	return packetCursor{probe: probeIPVersion, l3Offset: layout.l3Off()}
}

func (c packetCursor) layout() linkLayout {
	return cursorLayout{cursor: c}
}

type cursorLayout struct {
	cursor packetCursor
}

func (l cursorLayout) l3Off() uint32 {
	return l.cursor.l3Offset
}

func (l cursorLayout) genLinkProbe() []bpf.Instruction {
	if l.cursor.probe == probeEtherType {
		return []bpf.Instruction{
			bpf.LoadAbsolute{Off: l.cursor.etherTypeOff, Size: lengthHalf},
		}
	}
	return []bpf.Instruction{
		bpf.LoadAbsolute{Off: l.cursor.l3Offset, Size: lengthByte},
		bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0xf0},
	}
}

func (l cursorLayout) genLinkType(proto uint32, skipTrue, skipFalse uint8) bpf.Instruction {
	value := proto
	if l.cursor.probe == probeIPVersion {
		switch proto {
		case etherTypeIPv4:
			value = 0x40
		case etherTypeIPv6:
			value = 0x60
		default:
			value = 0xf0
		}
	}
	return bpf.JumpIf{
		Cond:      bpf.JumpEqual,
		Val:       value,
		SkipTrue:  skipTrue,
		SkipFalse: skipFalse,
	}
}

func (l cursorLayout) linkProbeSize() uint8 {
	if l.cursor.probe == probeEtherType {
		return 1
	}
	return 2
}

func (l cursorLayout) linkCompareSize() uint8 {
	return 1
}

func (l cursorLayout) hasL2Protocols() bool {
	return l.cursor.hasEther
}

type emitContinuation func(*prog, packetCursor)

type emitBranches struct {
	onMatch emitContinuation
	onMiss  emitContinuation
}

func compileCursorFilter(f Filter, layout linkLayout, target CompileTarget) ([]bpf.Instruction, error) {
	b := newProg(layout, target)
	emitFilterWithCursor(
		f,
		b,
		cursorFor(layout),
		emitBranches{
			onMatch: func(b *prog, _ packetCursor) { b.emitJump(labelKeep) },
			onMiss:  func(b *prog, _ packetCursor) { b.emitJump(labelFail) },
		},
	)
	return b.finalize()
}

func emitFilterWithCursor(
	f Filter,
	b *prog,
	cursor packetCursor,
	branches emitBranches,
) {
	switch v := f.(type) {
	case primitive:
		emitPrimitiveWithCursor(&v, b, cursor, branches)
	case arithmeticComparison:
		emitComparisonWithCursor(v, b, cursor, branches)
	case negated:
		emitFilterWithCursor(
			v.inner,
			b,
			cursor,
			emitBranches{
				onMatch: func(b *prog, _ packetCursor) { branches.onMiss(b, cursor) },
				onMiss:  func(b *prog, _ packetCursor) { branches.onMatch(b, cursor) },
			},
		)
	case composite:
		emitCompositeWithCursor(v, b, cursor, branches)
	}
}

func emitPrimitiveWithCursor(
	p *primitive,
	b *prog,
	cursor packetCursor,
	branches emitBranches,
) {
	if shouldEmitLinuxSocketVLAN(p, b, cursor) {
		emitLinuxSocketVLANWithCursor(p, b, cursor, branches)
		return
	}

	matched := b.newLabel()
	missed := b.newLabel()
	b.layout = cursor.layout()
	if cursor.mplsDepth > 0 && p.kind != filterKindMPLS && primitiveUsesInnerNetwork(p) {
		bottomOfStack := b.newLabel()
		b.emit(bpf.LoadAbsolute{Off: cursor.l3Offset - 2, Size: lengthByte})
		b.emitJumpIf(bpf.JumpBitsSet, 0x01, bottomOfStack, missed)
		b.bind(bottomOfStack)
	}
	p.emit(b, matched, missed)

	b.bind(matched)
	next := cursor
	switch p.kind {
	case filterKindVLAN:
		next.etherTypeOff += vlanHeaderLength
		next.l3Offset += vlanHeaderLength
		next.vlanDepth++
	case filterKindMPLS:
		if next.mplsDepth == 0 {
			next.probe = probeIPVersion
		}
		next.l3Offset += mplsLabelLength
		next.mplsDepth++
	}
	branches.onMatch(b, next)

	b.bind(missed)
	branches.onMiss(b, cursor)
}

func shouldEmitLinuxSocketVLAN(p *primitive, b *prog, cursor packetCursor) bool {
	return p.kind == filterKindVLAN &&
		b.target == CompileTargetLinuxSocket &&
		cursor.hasEther &&
		cursor.probe == probeEtherType &&
		cursor.etherTypeOff == 12 &&
		cursor.mplsDepth == 0 &&
		cursor.vlanDepth == 0
}

func emitLinuxSocketVLANWithCursor(
	p *primitive,
	b *prog,
	cursor packetCursor,
	branches emitBranches,
) {
	inlinePath := b.newLabel()
	metadataPath := b.newLabel()
	inlineMatched := b.newLabel()
	metadataMatched := b.newLabel()
	missed := b.newLabel()
	done := b.newLabel()

	b.emit(bpf.LoadAbsolute{Off: linuxBPFExtVLANTagPresent, Size: lengthByte})
	b.emitJumpIf(bpf.JumpEqual, 1, metadataPath, inlinePath)

	b.bind(inlinePath)
	b.layout = cursor.layout()
	p.emit(b, inlineMatched, missed)
	b.bind(inlineMatched)
	inlineNext := cursor
	inlineNext.etherTypeOff += vlanHeaderLength
	inlineNext.l3Offset += vlanHeaderLength
	inlineNext.vlanDepth++
	branches.onMatch(b, inlineNext)
	b.emitJump(done)

	b.bind(metadataPath)
	if p.id == "" {
		b.emitJump(metadataMatched)
	} else {
		vlanID, err := strconv.ParseUint(p.id, 10, 12)
		if err != nil {
			b.emitJump(missed)
		} else {
			b.emit(bpf.LoadAbsolute{Off: linuxBPFExtVLANTag, Size: lengthHalf})
			b.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0x0fff})
			b.emitJumpIf(bpf.JumpEqual, uint32(vlanID), metadataMatched, missed)
		}
	}
	b.bind(metadataMatched)
	metadataNext := cursor
	metadataNext.vlanDepth++
	branches.onMatch(b, metadataNext)
	b.emitJump(done)

	b.bind(missed)
	branches.onMiss(b, cursor)

	b.bind(done)
}

func primitiveUsesInnerNetwork(p *primitive) bool {
	switch p.kind {
	case filterKindHost:
		return p.protocol != filterProtocolEther
	case filterKindMulticast:
		return p.protocol == filterProtocolIP || p.protocol == filterProtocolIP6
	case filterKindUnset:
		return p.protocol != filterProtocolEther && p.subProtocol != filterSubProtocolSTP
	case filterKindVLAN, filterKindMPLS:
		return false
	default:
		return true
	}
}

func emitCompositeWithCursor(
	c composite,
	b *prog,
	cursor packetCursor,
	branches emitBranches,
) {
	var emitAt func(int, packetCursor)
	emitAt = func(index int, current packetCursor) {
		if index == len(c.filters) {
			branches.onMatch(b, current)
			return
		}
		if c.and {
			emitFilterWithCursor(
				c.filters[index],
				b,
				current,
				emitBranches{
					onMatch: func(_ *prog, next packetCursor) { emitAt(index+1, next) },
					onMiss:  branches.onMiss,
				},
			)
			return
		}
		emitFilterWithCursor(
			c.filters[index],
			b,
			current,
			emitBranches{
				onMatch: branches.onMatch,
				onMiss:  func(_ *prog, _ packetCursor) { emitAt(index+1, cursor) },
			},
		)
	}
	emitAt(0, cursor)
}

func filterNeedsCursor(f Filter) bool {
	switch v := f.(type) {
	case primitive:
		return v.kind == filterKindVLAN || v.kind == filterKindMPLS
	case arithmeticComparison:
		return true
	case negated:
		return filterNeedsCursor(v.inner)
	case composite:
		for _, child := range v.filters {
			if filterNeedsCursor(child) {
				return true
			}
		}
	}
	return false
}
