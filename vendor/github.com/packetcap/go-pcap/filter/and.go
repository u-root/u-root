package filter

// and is a type that implements Element and reports if it is "and" or "or"
type and bool

func (a and) Type() ElementType {
	return Joiner
}
