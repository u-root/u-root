package filter

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"

	"golang.org/x/net/bpf"
)

var resolver net.Resolver

// primitive implements Filter and Element
type primitive struct {
	kind        filterKind
	direction   filterDirection
	protocol    filterProtocol
	subProtocol filterSubProtocol
	negator     bool
	id          string
}

func (p primitive) IsPrimitive() bool {
	return true
}
func (p primitive) Type() ElementType {
	return Primitive
}

func (p primitive) Distill() Filter {
	return p
}

// Combine combines this primitive with another primitive, if they are combinable,
// without any loss of information. If they are not combinable, returns nil; if they
// are, returns a new primitive that represents both.
func (p primitive) Combine(o *primitive) *primitive {
	if p.Equal(o) {
		return &p
	}
	// our definition of "combinable" is: all of the fields that are set in one are either
	// set to the same value in the other, or Unset
	c := primitive{}
	switch {
	case p.kind == o.kind || o.kind == filterKindUnset:
		c.kind = p.kind
	case p.kind == filterKindUnset:
		c.kind = o.kind
	default:
		return nil
	}

	switch {
	case p.direction == o.direction || o.direction == filterDirectionUnset:
		c.direction = p.direction
	case p.direction == filterDirectionUnset:
		c.direction = o.direction
	default:
		return nil
	}

	switch {
	case p.protocol == o.protocol || o.protocol == filterProtocolUnset:
		c.protocol = p.protocol
	case p.protocol == filterProtocolUnset:
		c.protocol = o.protocol
	default:
		return nil
	}

	switch {
	case p.subProtocol == o.subProtocol || o.subProtocol == filterSubProtocolUnset:
		c.subProtocol = p.subProtocol
	case p.subProtocol == filterSubProtocolUnset:
		c.subProtocol = o.subProtocol
	default:
		return nil
	}

	switch {
	case p.id == o.id || o.id == "":
		c.id = p.id
	case p.id == "":
		c.id = o.id
	default:
		return nil
	}

	switch {
	case p.negator == o.negator:
		c.negator = p.negator
	default:
		return nil
	}

	return &c
}

func (p primitive) Compile() ([]bpf.Instruction, error) {
	// validate it
	if err := p.validate(); err != nil {
		return nil, err
	}

	// made it this far, so it must be valid. Compile it to bpf Instructions

	// calculate the total number of instructions
	// there always is at least the return packet and return none
	inst := instructions{
		inst: make([]bpf.Instruction, 0),
		size: p.Size(),
	}

	// if there are any conditions, there is a possibility of returning 0
	if p.kind == filterKindHost {
		switch p.protocol {
		case filterProtocolEther:
			inst.append(checkEtherAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolIP6:
			inst.append(loadEtherKind)
			inst.append(compareProtocolIP6(0, inst.skipToFail()))
			// ignore errors as it already has been validated
			_, a6, _ := p.getAddrs()
			inst.append(checkIP6HostAddresses(p.direction, a6[0], inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolIP:
			inst.append(loadEtherKind)
			inst.append(compareProtocolIP4(0, inst.skipToFail()))
			// ignore errors as it already has been validated
			a4, _, _ := p.getAddrs()
			inst.append(checkIP4HostAddresses(p.direction, a4[0], inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolArp:
			inst.append(loadEtherKind)
			inst.append(compareProtocolArp(0, inst.skipToFail()))
			// ignore errors as it already has been validated
			a4, _, _ := p.getAddrs()
			inst.append(checkIP4ArpAddresses(p.direction, a4[0], inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolRarp:
			inst.append(loadEtherKind)
			inst.append(compareProtocolRarp(0, inst.skipToFail()))
			// ignore errors as it already has been validated
			a4, _, _ := p.getAddrs()
			inst.append(checkIP4ArpAddresses(p.direction, a4[0], inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolUnset:
			// compare to the type

			// ignore errors as it already has been validated
			// we do, however, need to know if the addresses are ip6 or ip4
			// it takes 2 steps to check the src or dst for ip4, 8 steps for ip6
			a4, a6, _ := p.getAddrs()

			inst.append(loadEtherKind)
			if len(a4) > 0 {
				var (
					addressCheck uint8 = 2
				)
				if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
					addressCheck = 4
				}
				inst.append(compareProtocolIP4(0, addressCheck))
				// compare IP addresses
				inst.append(checkIP4HostAddresses(p.direction, a4[0], inst.skipToFail(), inst.skipToSucceed())...)
				// if Arp, go to arp addresses
				inst.append(compareProtocolArp(1, 0))
				// if not rarp, jump to next (if there is) or fail
				nextStep := inst.skipToFail()
				if len(a6) > 0 {
					nextStep = 2
					if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
						nextStep = 4
					}
				}
				inst.append(compareProtocolRarp(0, nextStep))
				inst.append(checkIP4ArpAddresses(p.direction, a4[0], inst.skipToFail(), inst.skipToSucceed())...)
			}
			if len(a6) > 0 {
				inst.append(compareProtocolIP6(0, inst.skipToFail()))
				inst.append(checkIP6HostAddresses(p.direction, a6[0], inst.skipToFail(), inst.skipToSucceed())...)
			}
		}
	}

	// port
	if p.kind == filterKindPort {
		// the port had better be valid
		portInt, err := findPort(p.id)
		if err != nil {
			return nil, err
		}

		port := uint32(portInt)
		inst.append(loadEtherKind)
		switch p.protocol {
		case filterProtocolIP6:
			inst.append(compareProtocolIP6(0, inst.skipToFail()))
			inst.append(loadIPv6Protocol)
			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareIPv6Protocol(ipProtocolTCP, 0, inst.skipToFail())...)
			case filterSubProtocolUDP:
				inst.append(compareIPv6Protocol(ipProtocolUDP, 0, inst.skipToFail())...)
			case filterSubProtocolStp:
				inst.append(compareSubProtocolSctp(0, inst.skipToFail()))
			case filterSubProtocolUnset:
				inst.append(compareSubProtocolSctp(2, 0))
				inst.append(compareSubProtocolTCP(1, 0))
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			}
			// compare IP addresses
			inst.append(checkPorts(p.direction, port, inst.skipToFail(), inst.skipToSucceed(), true)...)
		case filterProtocolIP:
			inst.append(compareProtocolIP4(0, inst.skipToFail()))
			inst.append(loadIPv4Protocol)
			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareSubProtocolTCP(0, inst.skipToFail()))
			case filterSubProtocolUDP:
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			case filterSubProtocolStp:
				inst.append(compareSubProtocolSctp(0, inst.skipToFail()))
			case filterSubProtocolUnset:
				inst.append(compareSubProtocolSctp(2, 0))
				inst.append(compareSubProtocolTCP(1, 0))
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			}
			inst.append(checkPorts(p.direction, port, inst.skipToFail(), inst.skipToSucceed(), false)...)
		case filterProtocolUnset:
			// this is a little backward, but I need to calculate how many steps in the
			// ip6 section so I can know where the ip4 section starts
			// first for loading the protocol and checking it
			var steps uint8 = 2
			if p.subProtocol == filterSubProtocolUnset {
				steps += 2
			}
			// next for loading the src and/or dst port and checking it
			steps += 2
			if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
				steps += 2
			}
			inst.append(compareProtocolIP6(0, steps))
			inst.append(loadIPv6Protocol)

			/* TODO: FIX HERE
			switch p.subProtocol {
			case filterSubProtocolUDP:
				inst.append(compareProtocolIP6(0, 5)) // size of compareIPv6Protocol
				inst.append(compareIPv6Protocol(ipProtocolUDP, inst.skipToSucceed(), inst.skipToFail())...)
				inst.append(compareProtocolIP4(0, inst.skipToFail()))
				inst.append(compareIPv4Protocol(ipProtocolUDP, 0, inst.skipToFail())...)
			case filterSubProtocolTCP:
				inst.append(compareProtocolIP6(0, 5)) // size of compareIPv6Protocol
				inst.append(compareIPv6Protocol(ipProtocolTCP, inst.skipToSucceed(), inst.skipToFail())...)
				inst.append(compareProtocolIP4(0, inst.skipToFail()))
				inst.append(compareIPv4Protocol(ipProtocolTCP, 0, inst.skipToFail())...)
			}
			*/

			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareSubProtocolTCP(0, inst.skipToFail()))
			case filterSubProtocolUDP:
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			case filterSubProtocolStp:
				inst.append(compareSubProtocolSctp(0, inst.skipToFail()))
			case filterSubProtocolUnset:
				inst.append(compareSubProtocolSctp(2, 0))
				inst.append(compareSubProtocolTCP(1, 0))
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			}
			inst.append(checkPorts(p.direction, port, inst.skipToFail(), inst.skipToSucceed(), true)...)
			inst.append(compareProtocolIP4(0, inst.skipToFail()))
			inst.append(loadIPv4Protocol)
			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareSubProtocolTCP(0, inst.skipToFail()))
			case filterSubProtocolUDP:
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			case filterSubProtocolStp:
				inst.append(compareSubProtocolSctp(0, inst.skipToFail()))
			case filterSubProtocolUnset:
				inst.append(compareSubProtocolSctp(2, 0))
				inst.append(compareSubProtocolTCP(1, 0))
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			}
			inst.append(checkPorts(p.direction, port, inst.skipToFail(), inst.skipToSucceed(), false)...)
		}
	}

	// net
	if p.kind == filterKindNet {
		switch p.protocol {
		case filterProtocolIP6:
			inst.append(loadEtherKind)
			inst.append(compareProtocolIP6(0, inst.skipToFail()))
			// ignore errors as it already has been validated
			addr, network, _ := getNetAndMask(p.id)
			inst.append(checkIP6NetAddresses(p.direction, addr, network.Mask, inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolIP:
			inst.append(loadEtherKind)
			inst.append(compareProtocolIP4(0, inst.skipToFail()))
			inst.append(checkIP4NetHostAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolArp:
			inst.append(loadEtherKind)
			inst.append(compareProtocolArp(0, inst.skipToFail()))
			inst.append(checkIP4NetArpAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolRarp:
			inst.append(loadEtherKind)
			inst.append(compareProtocolRarp(0, inst.skipToFail()))
			inst.append(checkIP4NetArpAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
		case filterProtocolUnset:
			inst.append(loadEtherKind)
			// more complicated. try each of several - if it is IP, next 4 are for the address
			// ignore error since it already was validated
			addr, network, _ := getNetAndMask(p.id)
			if addr.To4() != nil {
				var addressCheck uint8 = 2
				if !bytes.Equal(network.Mask, ip4MaskFull) {
					addressCheck++
				}
				if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
					addressCheck *= 2
				}
				inst.append(compareProtocolIP4(0, addressCheck))
				// compare IP addresses
				inst.append(checkIP4NetHostAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
				// if Arp, go to arp addresses
				inst.append(compareProtocolArp(1, 0))
				// if not rarp, nothing left
				inst.append(compareProtocolRarp(0, inst.skipToFail()))
				// compare arp/rarp addresses
				inst.append(checkIP4NetArpAddresses(p.direction, p.id, inst.skipToFail(), inst.skipToSucceed())...)
			} else {
				inst.append(compareProtocolIP6(0, inst.skipToFail()))
				inst.append(checkIP6NetAddresses(p.direction, addr, network.Mask, inst.skipToFail(), inst.skipToSucceed())...)
			}
		}
	}

	// unset
	if p.kind == filterKindUnset {
		inst.append(loadEtherKind)
		switch p.protocol {
		case filterProtocolIP:
			inst.append(compareProtocolIP4(0, inst.skipToFail()))
			inst.append(loadIPv4Protocol)
			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareSubProtocolTCP(0, inst.skipToFail()))
			case filterSubProtocolUDP:
				inst.append(compareSubProtocolUDP(0, inst.skipToFail()))
			}
		case filterProtocolIP6:
			inst.append(compareProtocolIP6(0, inst.skipToFail()))
			switch p.subProtocol {
			case filterSubProtocolTCP:
				inst.append(compareIPv6Protocol(ipProtocolTCP, 0, inst.skipToFail())...)
			case filterSubProtocolUDP:
				inst.append(compareIPv6Protocol(ipProtocolUDP, 0, inst.skipToFail())...)
			}
		case filterProtocolArp:
			inst.append(compareProtocolArp(0, inst.skipToFail()))
		case filterProtocolRarp:
			inst.append(compareProtocolRarp(0, inst.skipToFail()))
		case filterProtocolEther:
			switch p.subProtocol {
			case filterSubProtocolIP:
				inst.append(compareProtocolIP4(0, inst.skipToFail()))
			case filterSubProtocolIP6:
				inst.append(compareProtocolIP6(0, inst.skipToFail()))
			case filterSubProtocolArp:
				inst.append(compareProtocolArp(0, inst.skipToFail()))
			case filterSubProtocolRarp:
				inst.append(compareProtocolRarp(0, inst.skipToFail()))
			}
		case filterProtocolUnset:
			// kind is unset, and protocol is unset, so subprotocol must be set or it would have failed vaildation
			switch p.subProtocol {
			case filterSubProtocolUDP:
				inst.append(compareProtocolIP6(0, 5)) // size of compareIPv6Protocol
				inst.append(compareIPv6Protocol(ipProtocolUDP, inst.skipToSucceed(), inst.skipToFail())...)
				inst.append(compareProtocolIP4(0, inst.skipToFail()))
				inst.append(compareIPv4Protocol(ipProtocolUDP, 0, inst.skipToFail())...)
			case filterSubProtocolTCP:
				inst.append(compareProtocolIP6(0, 5)) // size of compareIPv6Protocol
				inst.append(compareIPv6Protocol(ipProtocolTCP, inst.skipToSucceed(), inst.skipToFail())...)
				inst.append(compareProtocolIP4(0, inst.skipToFail()))
				inst.append(compareIPv4Protocol(ipProtocolTCP, 0, inst.skipToFail())...)
			}
		}
	}

	if p.negator {
		// Add the instruction to accept packets that did not match the original condition
		inst.append(returnDrop)
		inst.append(returnKeep)
	} else {
		inst.append(returnKeep)
		inst.append(returnDrop)
	}

	return inst.inst, nil
}

func (p primitive) Equal(f Filter) bool {
	if f == nil {
		return false
	}
	var (
		o  primitive
		ok bool
	)
	if o, ok = f.(primitive); !ok {
		return false
	}
	return p.kind == o.kind &&
		p.direction == o.direction &&
		p.protocol == o.protocol &&
		p.subProtocol == o.subProtocol &&
		p.negator == o.negator &&
		p.id == o.id
}

func (p primitive) validate() error {
	switch {
	case p.subProtocol == filterSubProtocolUnknown:
		return fmt.Errorf("unknown protocol %s", p.id)
	case p.kind == filterKindHost:
		switch p.protocol {
		case filterProtocolIP, filterProtocolIP6, filterProtocolArp, filterProtocolRarp, filterProtocolUnset:
			// must be IP or valid host
			if p.id == "" {
				return fmt.Errorf("blank host")
			}
			// if it is in IP format, check the IP validity
			addr, network, _ := getNetAndMask(p.id)
			var maskFull net.IPMask
			if addr != nil && network != nil {
				if addr.To4() != nil {
					maskFull = ip4MaskFull
				} else {
					maskFull = ip6MaskFull
				}
				if !bytes.Equal(network.Mask, maskFull) {
					return fmt.Errorf("invalid host address with CIDR: %s", p.id)
				}
			}
			// if it was not a valid IP, check if it is a valid hostname
			if addr == nil {
				a4, a6, err := p.getAddrs()
				if err != nil || (len(a4)+len(a6) == 0) {
					return fmt.Errorf("unknown host: %s", p.id)
				}
				for _, a := range a4 {
					if a == nil {
						return fmt.Errorf("invalid address return in lookup: %s", a)
					}
				}
				for _, a := range a6 {
					if a == nil {
						return fmt.Errorf("invalid address return in lookup: %s", a)
					}
				}
			}
		case filterProtocolEther:
			// check that it is a valid ether host format
			if _, err := net.ParseMAC(p.id); err != nil {
				return fmt.Errorf("invalid ethernet address: %s", p.id)
			}
		}
	case p.kind == filterKindUnset && p.protocol == filterProtocolUnset && p.subProtocol == filterSubProtocolUnset:
		return fmt.Errorf("parse error")
	case p.kind == filterKindPort:
		if _, err := findPort(p.id); err != nil {
			return err
		}
	case p.kind == filterKindNet:
		// network must be one of:
		// - straight IP (v4 or v6)
		// - valid CIDR, but all bits after the mask must be 0
		addr, network, err := getNetAndMask(p.id)
		if err != nil {
			return err
		}
		masked := addr.Mask(network.Mask)
		if !addr.Equal(masked) {
			return fmt.Errorf("invalid network, network bits extend past mask bits: %s", p.id)
		}
	case p.kind == filterKindUnset && p.protocol == filterProtocolEther && p.subProtocol == filterSubProtocolUnset:
		return fmt.Errorf("parse error")
	}
	return nil
}

// Size how many instructions do we expect
func (p primitive) Size() uint8 {
	var instCount uint8
	// if there are any conditions, there is a possibility of returning 0
	switch p.kind {
	case filterKindHost:
		instCount += p.calculateStepsKindHost()
	case filterKindPort:
		instCount += p.calculateStepsKindPort()
	case filterKindUnset:
		instCount += p.calculateStepsKindUnset()
	case filterKindNet:
		instCount += p.calculateStepsKindNet()
	}

	return instCount + 2
}

// getAddrs get valid IP addresses for the provided string, whether ipv4, ipv6,
// or hostname
func (p primitive) getAddrs() ([]net.IP, []net.IP, error) {
	a6, a4, addrs := []net.IP{}, []net.IP{}, []net.IP{}
	// first see if it is a regular IP address
	if addr := net.ParseIP(p.id); addr != nil {
		addrs = append(addrs, addr)
	} else {
		// look up the host; ignore error as it already should have been done
		resolvedAddrs, _ := resolver.LookupHost(context.Background(), p.id)
		for _, a := range resolvedAddrs {
			addrs = append(addrs, net.ParseIP(a))
		}
	}
	for _, a := range addrs {
		if a.To4() != nil {
			a4 = append(a4, a)
		} else {
			a6 = append(a6, a)
		}
	}
	return a4, a6, nil
}

// calculateStepsKindHost determine the number of steps for a filter of kind host
func (p primitive) calculateStepsKindHost() uint8 {
	// do we need to use separate locations to check for the src and/or dst?
	// only if the protocol is arp/rarp *and* ip
	var (
		count, dirCount uint8
	)

	switch p.protocol {
	case filterProtocolIP6:
		// load the ether protocol
		count++
		// compare to the type
		count++
		// ignore errors as it already has been validated
		_, a6, _ := p.getAddrs()
		// it takes 8 steps to check each src or dst
		dirCount = 8 * uint8(len(a6))
	case filterProtocolIP, filterProtocolArp, filterProtocolRarp:
		// load the ether protocol
		count++
		// compare to the type
		count++
		// ignore errors as it already has been validated
		a4, _, _ := p.getAddrs()
		// it takes 2 steps to check each src or dst
		dirCount = 2 * uint8(len(a4))
	case filterProtocolUnset:
		// compare to the type
		count++
		// ignore errors as it already has been validated
		// we do, however, need to know if the addresses are ip6 or ip4
		// it takes 2 steps to check the src or dst for ip4, 8 steps for ip6
		a4, a6, _ := p.getAddrs()
		// it takes 2 steps for each src or dst in ip4
		// and then another 2 steps for each src or dst in arp/rarp
		dirCount = dirCount + uint8((2+2)*len(a4))
		// it takes 8 steps for each src or dst in ip6
		dirCount = dirCount + uint8(8*len(a6))
		if len(a4) > 0 {
			count += 3 // compare ip4, arp, rarp
		}
		if len(a6) > 0 {
			count++ // compare ip6
		}
	case filterProtocolEther:
		// it takes 4 steps to check the src or dst, since it takes 2 distinct sets to read the 6 bytes
		dirCount = 4
	}

	//
	if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
		dirCount *= 2
	}

	count += dirCount
	return count
}

// calculateStepsKindNet determine the number of steps for a filter of kind net
func (p primitive) calculateStepsKindNet() uint8 {
	// do we need to use separate locations to check for the src and/or dst?
	// only if the protocol is arp/rarp *and* ip
	var (
		count, dirCount uint8
		doubler         bool
		maskFull        net.IPMask
	)
	// no real erro handling here, and it should already have been validated
	addr, network, _ := getNetAndMask(p.id)

	switch p.protocol {
	case filterProtocolIP, filterProtocolArp, filterProtocolRarp:
		// load the ether protocol
		count++
		// compare to the type
		count++
		// it takes 2 steps to check the src or dst
		dirCount = 2
		maskFull = ip4MaskFull
	case filterProtocolIP6:
		// load the ether protocol
		count++
		// compare to the type
		count++
		dirCount += calculateIP6MaskSteps(network.Mask)
		maskFull = ip6MaskFull
	case filterProtocolUnset:
		// compare to the type
		count++
		// ignore errors as it already has been validated
		// we do, however, need to know if the addresses are ip6 or ip4
		// it takes 2 steps to check the src or dst for ip4, 8 steps for ip6
		if addr.To4() != nil {
			// it takes 2 steps for each src or dst in ip4
			dirCount += 2
			// compare to the 3 types
			count += 3
			doubler = true
			maskFull = ip4MaskFull
		} else {
			dirCount += calculateIP6MaskSteps(network.Mask)
			// compare to the one type
			count++
			maskFull = ip6MaskFull
		}
	}

	// if the netmask is not "mask full" (0xffffffff for ip4, larger for ip6), then we need to add a
	// step to each direction for netmask
	if !bytes.Equal(network.Mask, maskFull) {
		dirCount++
	}

	//
	if p.direction == filterDirectionSrcOrDst || p.direction == filterDirectionSrcAndDst {
		dirCount *= 2
	}

	if doubler {
		dirCount *= 2
	}
	count += dirCount
	return count
}

func (p primitive) calculateStepsKindPort() uint8 {
	// we will load the ether protocol, which takes 2 for ip4 or ip6, 3 for both;
	// then the ip protocol, which takes 2 for ip4 or ip6, 4 for both

	var (
		count   uint8
		doubler bool
	)
	// load the ether protocol and compare
	count += 2
	if p.protocol == filterProtocolUnset {
		count++
		doubler = true
	}

	var subProtocolCount uint8 = 2

	// port is only relevant for ip4/ip6

	// load the ip protocol and compare
	if p.subProtocol == filterSubProtocolUnset {
		subProtocolCount += 2
	}

	// checking ports on ipv6 is 2 for each of src and/or dst
	// checking ports on ipv4 is 2 for each of src and/or dst, plus 3 to calculate the location
	switch p.direction {
	case filterDirectionSrc, filterDirectionDst:
		subProtocolCount += 2
	case filterDirectionSrcOrDst, filterDirectionSrcAndDst:
		subProtocolCount += 4
	}
	if doubler {
		subProtocolCount *= 2
	}

	count += subProtocolCount

	// for ip4 (or unset, which is ip4+ip6), we need 3 more steps to find where the src/dst port are
	if p.protocol == filterProtocolIP || p.protocol == filterProtocolUnset {
		count += 3
	}
	return count
}

// calculateStepsKindUnset determine the number of steps for a filter of unset kind
func (p primitive) calculateStepsKindUnset() uint8 {
	// this already should have been validated
	var (
		count uint8
	)
	// 2 to load and compare the ether protocol
	// 2 more to load and compare the sub protocol, if provided
	count += 2
	switch {
	case p.protocol == filterProtocolUnset:
		// protocol is unset in addition to kind, so it depends on the subprotocol
		count++    // check ipv4 and ipv6
		count += 2 // 2 for ipv6 protocol check
		count += 3 // 3 for ipv6 continuation packet protocol check
		count += 2 // 2 for ipv4 protocol check
	case p.protocol != filterProtocolEther:
		count += 2 // for ether, it already was covered
	}
	return count
}

func findPort(portStr string) (int, error) {
	// check that it is either an integer, or a known and valid port
	if port, err := strconv.Atoi(portStr); err == nil {
		return port, nil
	}
	if port, err := net.LookupPort("tcp", portStr); err == nil {
		return port, nil
	}
	if port, err := net.LookupPort("udp", portStr); err == nil {
		return port, nil
	}
	return -1, fmt.Errorf("invalid port: %s", portStr)
}
