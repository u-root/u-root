// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qnetwork

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/hugelgupf/vmtest/qemu"
)

// WithBackend adds a backend to the net device.
func WithBackend[B Backend](mods ...Modifier[B]) NetDevModifier[B] {
	return func(netdevID string, alloc *qemu.IDAllocator, opts *qemu.Options, nd *NetDevice[B]) error {
		for _, mod := range mods {
			if mod != nil {
				if err := mod(&nd.Backend); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// Backend is a network backend servicing an emulated guest device.
type Backend interface {
	NetDev(id string) string
	Validate() error
}

// Modifier modifies a backend.
type Modifier[B Backend] func(b *B) error

// UserBackend is a user mode host networking backend.
//
// It uses SLIRP to parse packets and figure out where to route them. By
// default they're injected in the host.
type UserBackend struct {
	IP4  net.IP
	Net4 *net.IPNet
	IP6  net.IP
	Net6 *net.IPNet
	Args []string
}

// WithUser is a net device modifier for UserBackend.
var WithUser = WithBackend[UserBackend]

// WithUserArg adds arbitrary args to the "-netdev user" invocation.
func WithUserArg(s ...string) Modifier[UserBackend] {
	return func(b *UserBackend) error {
		b.Args = append(b.Args, s...)
		return nil
	}
}

// WithUserCIDR sets an IPv4 or IPv6 user network.
func WithUserCIDR(cidr string) Modifier[UserBackend] {
	return func(b *UserBackend) error {
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		if ip.To4() != nil {
			if b.IP4 != nil || b.Net4 != nil {
				return fmt.Errorf("%w: IPv4 net already set", os.ErrInvalid)
			}
			b.IP4 = ip
			b.Net4 = ipnet
		} else {
			if b.IP6 != nil || b.Net6 != nil {
				return fmt.Errorf("%w: IPv6 net already set", os.ErrInvalid)
			}
			b.IP6 = ip
			b.Net6 = ipnet
		}
		return nil
	}
}

// NetDev returns the arg for "-netdev".
func (b UserBackend) NetDev(id string) string {
	s := []string{"user", "id=" + id}
	if b.Net4 != nil {
		s = append(s, "ipv4=on", fmt.Sprintf("net=%s", b.Net4))
	} else {
		s = append(s, "ipv4=off")
	}
	if b.Net6 != nil {
		s = append(s, "ipv6=on", fmt.Sprintf("ipv6-net=%s", b.Net6))
	} else {
		s = append(s, "ipv6=off")
	}
	s = append(s, b.Args...)
	return strings.Join(s, ",")
}

// Validate validates the values in UserBackend.
func (b UserBackend) Validate() error {
	if b.Net4 == nil && b.Net6 == nil {
		return fmt.Errorf("cannot create user backend without v4 or v6 network definition")
	}
	return nil
}

// HostNetwork creates a user-backed net device with the given CIDR.
func HostNetwork(cidr string, mods ...NetDevModifier[UserBackend]) qemu.Fn {
	mods = append([]NetDevModifier[UserBackend]{
		WithUser(WithUserCIDR(cidr)),
	}, mods...)
	return New(mods...)
}

// SocketBackend is a Unix domain socket backend.
type SocketBackend struct {
	Server     bool
	UnixSocket string
	Args       []string
}

// WithSocket is a net device modifier for SocketBackend.
var WithSocket = WithBackend[SocketBackend]

// NetDev returns the arg for "-netdev".
func (b SocketBackend) NetDev(id string) string {
	s := append([]string{
		"stream",
		"id=" + id,
		fmt.Sprintf("server=%t", b.Server),
		"addr.type=unix",
		"addr.path=" + b.UnixSocket,
	}, b.Args...)
	return strings.Join(s, ",")
}

// Validate validates SocketBackend values.
func (b SocketBackend) Validate() error {
	return nil
}

// IsServer sets whether the socket backend is a server or not.
func IsServer(s bool) Modifier[SocketBackend] {
	return func(b *SocketBackend) error {
		b.Server = s
		return nil
	}
}

// WithUnixSocket sets the unix domain socket address.
func WithUnixSocket(socket string) Modifier[SocketBackend] {
	return func(b *SocketBackend) error {
		b.UnixSocket = socket
		return nil
	}
}
