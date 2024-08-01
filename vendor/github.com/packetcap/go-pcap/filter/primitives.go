package filter

// primitives a slice of primitive with some methods
type primitives []primitive

// combine combine multiple primitives until nothing is left to combine
func (p *primitives) combine() *primitives {
	// nothing to combine
	if p == nil || len(*p) == 0 || len(*p) == 1 {
		return p
	}
	// The simplest first cut is to have each one combine with its neighbour.
	// It isn't perfect - e.g. we will miss A combining with D - but it is a good start.
	pd := *p
	list := make(primitives, 0)
	var (
		prev, elm primitive
		lastMatch bool
		i         int
	)
	for i, elm = range pd {
		// do not bother combining with myself
		if i == 0 {
			prev = elm
			continue
		}
		if n := prev.Combine(&elm); n != nil {
			lastMatch = true
			list = append(list, *n)
			prev = *n
		} else {
			lastMatch = false
			list = append(list, prev)
			prev = elm
		}
	}
	// add the last element if it was not merged with the previous
	if !lastMatch {
		list = append(list, elm)
	}
	return &list
}

func (p *primitives) equal(o *primitives) bool {
	if o == nil {
		return false
	}
	pd, od := *p, *o
	if len(pd) != len(od) {
		return false
	}
	for i, p1 := range pd {
		if !p1.Equal(od[i]) {
			return false
		}
	}
	return true
}
