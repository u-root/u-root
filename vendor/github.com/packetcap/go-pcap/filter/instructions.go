package filter

import (
	"golang.org/x/net/bpf"
)

type instructions struct {
	inst []bpf.Instruction
	size uint8
}

// append add steps to the instruction slide
func (i *instructions) append(in ...bpf.Instruction) {
	i.inst = append(i.inst, in...)
}

// skipToFail how many steps the *next* step will skip to failure
func (i *instructions) skipToFail() uint8 {
	return i.size - uint8(len(i.inst)) - 2
}

// skipToSucceed how many steps the *next* step will skip to succeed
func (i *instructions) skipToSucceed() uint8 {
	return i.size - uint8(len(i.inst)) - 3
}
