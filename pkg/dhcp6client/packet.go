// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp6client

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/dhcp6opts"
)

// RequestIANAFrom returns a Request packet to request an IANA from the server
// that sent the given `ad` Advertisement.
//
// RFC 3315 Section 18.1.1. determines how a Request packet should be created.
func RequestIANAFrom(ad *dhcp6.Packet) (*dhcp6.Packet, error) {
	if ad.MessageType != dhcp6.MessageTypeAdvertise {
		return nil, fmt.Errorf("need advertise packet")
	}

	// Must include; RFC 3315 Section 18.1.1.
	clientID, err := dhcp6opts.GetClientID(ad.Options)
	if err != nil {
		return nil, fmt.Errorf("couldn't find client ID in %v: %v", ad, err)
	}

	serverID, err := dhcp6opts.GetServerID(ad.Options)
	if err != nil {
		return nil, fmt.Errorf("couldn't find server ID in %v: %v", ad, err)
	}

	opts := make(dhcp6.Options)
	if err := opts.Add(dhcp6.OptionClientID, clientID); err != nil {
		return nil, err
	}
	if err := opts.Add(dhcp6.OptionServerID, serverID); err != nil {
		return nil, err
	}
	if err := newRequestOptions(opts); err != nil {
		return nil, err
	}

	return NewPacket(dhcp6.MessageTypeRequest, opts), nil
}

func newRequestOptions(options dhcp6.Options) error {
	// TODO: This should be generated.
	id := [4]byte{'r', 'o', 'o', 't'}
	iana := dhcp6opts.NewIANA(id, 0, 0, nil)
	// IANA = requesting a non-temporary address.
	if err := options.Add(dhcp6.OptionIANA, iana); err != nil {
		return err
	}
	if err := options.Add(dhcp6.OptionElapsedTime, dhcp6opts.ElapsedTime(0)); err != nil {
		return err
	}

	oro := dhcp6opts.OptionRequestOption{
		dhcp6.OptionDNSServers,
		dhcp6.OptionBootFileURL,
		dhcp6.OptionBootFileParam,
	}
	// Must include; RFC 3315 Section 18.1.1.
	return options.Add(dhcp6.OptionORO, oro)
}

func newSolicitOptions(mac net.HardwareAddr) (dhcp6.Options, error) {
	options := make(dhcp6.Options)

	if err := newRequestOptions(options); err != nil {
		return nil, err
	}

	if err := options.Add(dhcp6.OptionClientID, dhcp6opts.NewDUIDLL(6, mac)); err != nil {
		return nil, err
	}
	return options, nil
}

// NewRapidSolicit returns a Solicit packet with the RapidCommit option.
func NewRapidSolicit(mac net.HardwareAddr) (*dhcp6.Packet, error) {
	p, err := NewSolicitPacket(mac)
	if err != nil {
		return nil, err
	}

	// Request an immediate Reply with an IP instead of an Advertise packet.
	if err := p.Options.Add(dhcp6.OptionRapidCommit, nil); err != nil {
		return nil, err
	}
	return p, nil
}

// NewSolicitPacket returns a Solicit packet.
//
// TODO(hugelgupf): Conform to RFC 3315 Section 17.1.1.
func NewSolicitPacket(mac net.HardwareAddr) (*dhcp6.Packet, error) {
	options, err := newSolicitOptions(mac)
	if err != nil {
		return nil, err
	}

	return NewPacket(dhcp6.MessageTypeSolicit, options), nil
}

// NewPacket creates a new DHCPv6 packet using the given message type and
// options.
//
// A transaction ID will be generated.
func NewPacket(typ dhcp6.MessageType, opts dhcp6.Options) *dhcp6.Packet {
	p := &dhcp6.Packet{
		MessageType: typ,
		Options:     opts,
	}

	// TODO: This may actually be bad news. Investigate whether we need to
	// use crypto/rand. An attacker could inject a bad response if this is
	// predictable. RFC 3315 has some words on this.
	rand.Read(p.TransactionID[:])
	return p
}
