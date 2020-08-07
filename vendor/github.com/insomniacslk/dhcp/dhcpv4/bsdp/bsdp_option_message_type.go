package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// MessageType represents the different BSDP message types.
//
// Implements the BSDP option message type. Can be one of LIST, SELECT, or
// FAILED.
type MessageType byte

// BSDP Message types - e.g. LIST, SELECT, FAILED
const (
	MessageTypeNone   MessageType = 0
	MessageTypeList   MessageType = 1
	MessageTypeSelect MessageType = 2
	MessageTypeFailed MessageType = 3
)

// ToBytes returns a serialized stream of bytes for this option.
func (m MessageType) ToBytes() []byte {
	return []byte{byte(m)}
}

// String returns a human-friendly representation of MessageType.
func (m MessageType) String() string {
	if s, ok := messageTypeToString[m]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", m)
}

// messageTypeToString maps each BSDP message type to a human-readable string.
var messageTypeToString = map[MessageType]string{
	MessageTypeList:   "LIST",
	MessageTypeSelect: "SELECT",
	MessageTypeFailed: "FAILED",
}

// FromBytes reads data into m.
func (m *MessageType) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	*m = MessageType(buf.Read8())
	return buf.FinError()
}

// OptMessageType returns a new BSDP Message Type option.
func OptMessageType(mt MessageType) dhcpv4.Option {
	return dhcpv4.Option{
		Code:  OptionMessageType,
		Value: mt,
	}
}
