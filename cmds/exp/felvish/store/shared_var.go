package store

import (
	"errors"
	"sync"
)

type shared struct {
	sync.Mutex
	vars map[string]string
}

// ErrNoVar is returned by (*Store).GetSharedVar when there is no such variable.
var ErrNoVar = errors.New("no such variable")

// SharedVar gets the value of a shared variable.
func (vars *shared) SharedVar(n string) (string, error) {
	vars.Lock()
	defer vars.Unlock()
	v, ok := vars.vars[n]
	if !ok {
		return "", ErrNoVar
	}
	return v, nil
}

// SetSharedVar sets the value of a shared variable.
func (vars *shared) SetSharedVar(n, v string) error {
	vars.Lock()
	defer vars.Unlock()
	vars.vars[n] = v
	return nil
}

// DelSharedVar deletes a shared variable.
func (vars *shared) DelSharedVar(n string) error {
	vars.Lock()
	defer vars.Unlock()
	delete(vars.vars, n)
	return nil
}

func NewSharedVar() *shared {
	return &shared{vars: make(map[string]string)}
}
