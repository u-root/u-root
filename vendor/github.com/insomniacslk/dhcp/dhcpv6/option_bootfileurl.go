package dhcpv6

import (
	"fmt"
)

// OptBootFileURL implements the OptionBootfileURL option
//
// This module defines the OptBootFileURL structure.
// https://www.ietf.org/rfc/rfc5970.txt
type OptBootFileURL struct {
	BootFileURL []byte
}

// Code returns the option code
func (op *OptBootFileURL) Code() OptionCode {
	return OptionBootfileURL
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptBootFileURL) ToBytes() []byte {
	return op.BootFileURL
}

func (op *OptBootFileURL) String() string {
	return fmt.Sprintf("OptBootFileURL{BootFileUrl=%s}", op.BootFileURL)
}

// ParseOptBootFileURL builds an OptBootFileURL structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptBootFileURL(data []byte) (*OptBootFileURL, error) {
	var opt OptBootFileURL
	opt.BootFileURL = append([]byte(nil), data...)
	return &opt, nil
}
