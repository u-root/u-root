package dhcpv6

import (
	"fmt"
)

// OptServerID represents a Server Identifier option as defined by RFC 3315
// Section 22.1.
func OptServerID(d DUID) Option {
	return &optServerID{d}
}

type optServerID struct {
	DUID
}

func (*optServerID) Code() OptionCode {
	return OptionServerID
}

func (op *optServerID) String() string {
	return fmt.Sprintf("%s: %v", op.Code(), op.DUID)
}

// parseOptServerID builds an optServerID structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func parseOptServerID(data []byte) (*optServerID, error) {
	sid, err := DUIDFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &optServerID{sid}, nil
}
