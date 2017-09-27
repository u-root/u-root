package dhcp6

// Packet represents a raw DHCPv6 packet, using the format described in RFC 3315,
// Section 6.
//
// The Packet type is typically only needed for low-level operations within the
// client, server, or in tests.
type Packet struct {
	// MessageType specifies the DHCP message type constant, such as
	// MessageTypeSolicit, MessageTypeAdvertise, etc.
	MessageType MessageType

	// TransactionID specifies the DHCP transaction ID.  The transaction ID must
	// be the same for all message exchanges in one DHCP transaction.
	TransactionID [3]byte

	// Options specifies a map of DHCP options.  Its methods can be used to
	// retrieve data from an incoming packet, or send data with an outgoing
	// packet.
	Options Options
}

// MarshalBinary allocates a byte slice containing the data
// from a Packet.
func (p *Packet) MarshalBinary() ([]byte, error) {
	// 1 byte: message type
	// 3 bytes: transaction ID
	// N bytes: options slice byte count
	opts := p.Options.enumerate()
	b := make([]byte, 4+opts.count())

	b[0] = byte(p.MessageType)
	copy(b[1:4], p.TransactionID[:])
	opts.write(b[4:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a Packet.
//
// If the byte slice does not contain enough data to form a valid Packet,
// ErrInvalidPacket is returned.
func (p *Packet) UnmarshalBinary(b []byte) error {
	// Packet must contain at least a message type and transaction ID
	if len(b) < 4 {
		return ErrInvalidPacket
	}
	p.MessageType = MessageType(b[0])

	txID := [3]byte{}
	copy(txID[:], b[1:4])
	p.TransactionID = txID

	options, err := parseOptions(b[4:])
	if err != nil {
		// Invalid options means an invalid packet
		return ErrInvalidPacket
	}
	p.Options = options

	return nil
}
