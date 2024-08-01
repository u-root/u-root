package filter

import (
	"golang.org/x/net/bpf"
)

// composite implements Filter
type composite struct {
	filters Filters
	and     bool
}

func (c composite) Compile() ([]bpf.Instruction, error) {
	// first compile each one, then go through them and join with the 'and' or 'or'
	//   - if 'and', then a failure of any one is straight to fail
	//   - if 'or', then a failure of any one means to move on to the next
	// The simplest way to implement is to just have interim jump steps.
	inst := []bpf.Instruction{}
	size := uint32(c.Size())
	for i, f := range c.filters {
		finst, err := f.Compile()
		if err != nil {
			return nil, err
		}
		// remove the last two instructions, which are the returns, if we are not on the last one
		if i == len(c.filters)-1 {
			inst = append(inst, finst...)
			continue
		}
		finst = finst[:len(finst)-2]
		inst = append(inst, finst...)
		// now add the jump to the next steppf.
		// the expectation of every primitive is that the second to last is success,
		// and the last is fail. For that step.
		if c.and {
			// Each step is required, so if the previous step failed, it just fails.
			// If it succeeded, go to the next one.
			inst = append(inst, bpf.Jump{Skip: 1})
			inst = append(inst, bpf.Jump{Skip: size - uint32(len(inst)) - 2})
		} else {
			// Each step is not required, so if the previous step failed, go to next.
			// If it succeeded, return success.
			inst = append(inst, bpf.Jump{Skip: size - uint32(len(inst)) - 3})
			inst = append(inst, bpf.Jump{Skip: 0})
		}
	}
	return inst, nil
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

// Size how many elements do we expect
func (c composite) Size() uint8 {
	var size uint8
	for _, f := range c.filters {
		size += f.Size()
	}
	return size
}

func (c composite) IsPrimitive() bool {
	return false
}
func (c composite) Type() ElementType {
	return Composite
}

func (c composite) LastPrimitive() *primitive {
	if len(c.filters) == 0 {
		return nil
	}
	last := c.filters[len(c.filters)-1]
	if !last.IsPrimitive() {
		return nil
	}
	p := last.(primitive)
	return &p
}

// Distill work through all of the member filters and see if some can be combined
func (c composite) Distill() Filter {
	list := make(Filters, 0)
	// do this in 2 rounds:
	// 1. We distill all of our children, as this may convert some composites into primitives
	// 2. Distill all of the primitives in our list, if we can
	for _, f := range c.filters {
		list = append(list, f.Distill())
	}
	c.filters = list
	// if there is just one element, return that one
	if len(c.filters) == 1 {
		return c.filters[0]
	}
	// only can distill with and
	if !c.and {
		return c
	}
	// we have "and" joiner, so perhaps we can combine overlapping elements
	prims := make(primitives, 0)
	compos := make([]composite, 0)
	for _, f := range c.filters {
		if f.IsPrimitive() {
			prims = append(prims, f.(primitive))
		} else {
			compos = append(compos, f.(composite))
		}
	}

	p2 := prims.combine()
	list = make(Filters, 0)
	for _, p := range *p2 {
		list = append(list, p)
	}
	for _, c := range compos {
		list = append(list, c)
	}
	c.filters = list
	// if there is just one element, return that one
	if len(c.filters) == 1 {
		return c.filters[0]
	}

	return c
}
