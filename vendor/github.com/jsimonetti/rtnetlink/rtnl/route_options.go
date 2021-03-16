package rtnl

import (
	"net"

	"github.com/jsimonetti/rtnetlink"
)

// RouteOptions is the functional options struct
type RouteOptions struct {
	Src   *net.IPNet
	Attrs rtnetlink.RouteAttributes
}

// RouteOption is the functional options func
type RouteOption func(*RouteOptions)

// DefaultRouteOptions defines the default route options.
func DefaultRouteOptions(ifc *net.Interface, dst net.IPNet, gw net.IP) *RouteOptions {
	ro := &RouteOptions{
		Src: nil,
		Attrs: rtnetlink.RouteAttributes{
			Dst:      dst.IP,
			OutIface: uint32(ifc.Index),
		},
	}

	if gw != nil {
		ro.Attrs.Gateway = gw
	}

	return ro
}

// WithRouteSrc sets the src address.
func WithRouteSrc(src *net.IPNet) RouteOption {
	return func(opts *RouteOptions) {
		opts.Src = src
	}
}

// WithRouteAttrs sets the attributes.
func WithRouteAttrs(attrs rtnetlink.RouteAttributes) RouteOption {
	return func(opts *RouteOptions) {
		opts.Attrs = attrs
	}
}
