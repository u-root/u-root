package vars

type blackhole struct{}

func (blackhole) Set(interface{}) error {
	return nil
}

func (blackhole) Get() interface{} {
	// TODO: Return a special placeholder value.
	return ""
}

// NewBlackhole returns a blackhole variable. Assignments to a blackhole
// variable will be discarded, and getting a blackhole variable always returns
// an empty string.
func NewBlackhole() Var {
	return blackhole{}
}

// IsBlackhole returns whether the variable is a blackhole variable.
func IsBlackhole(v Var) bool {
	_, ok := v.(blackhole)
	return ok
}
