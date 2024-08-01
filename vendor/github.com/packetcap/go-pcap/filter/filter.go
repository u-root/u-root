package filter

import (
	"sort"

	"golang.org/x/net/bpf"
)

// Filter constructed of a tcpdump filter expression
type Filter interface {
	Compile() ([]bpf.Instruction, error)
	Equal(o Filter) bool
	Size() uint8
	IsPrimitive() bool
	Type() ElementType
	Distill() Filter
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
