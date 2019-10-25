package dhcpv6

import (
	"fmt"
)

// OptBootFileURL implements the OptionBootfileURL option
//
// This module defines the OptBootFileURL structure.
// https://www.ietf.org/rfc/rfc5970.txt
type OptBootFileURL string

var _ Option = OptBootFileURL("")

// Code returns the option code
func (op OptBootFileURL) Code() OptionCode {
	return OptionBootfileURL
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op OptBootFileURL) ToBytes() []byte {
	return []byte(op)
}

func (op OptBootFileURL) String() string {
	return fmt.Sprintf("OptBootFileURL(%s)", string(op))
}

// ParseOptBootFileURL builds an OptBootFileURL structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptBootFileURL(data []byte) (OptBootFileURL, error) {
	return OptBootFileURL(string(data)), nil
}
