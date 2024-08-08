package tc

import (
	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Chain represents a collection of filter
type Chain struct {
	Tc
}

// Chain allows to read and alter chains
func (tc *Tc) Chain() *Chain {
	return &Chain{*tc}
}

// Add creates a new chain
func (c *Chain) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(unix.RTM_NEWCHAIN, info)
	if err != nil {
		return err
	}
	return c.action(unix.RTM_NEWCHAIN, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Delete removes a chain
func (c *Chain) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(unix.RTM_DELCHAIN, info)
	if err != nil {
		return err
	}
	return c.action(unix.RTM_DELCHAIN, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches chains
func (c *Chain) Get(i *Msg) ([]Object, error) {
	if i == nil {
		return []Object{}, ErrNoArg
	}
	return c.get(unix.RTM_GETCHAIN, i)
}
