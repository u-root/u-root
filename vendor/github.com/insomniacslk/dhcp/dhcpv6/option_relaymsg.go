package dhcpv6

// This module defines the OptRelayMsg structure.
// https://www.ietf.org/rfc/rfc3315.txt

import (
	"fmt"
)

type OptRelayMsg struct {
	relayMessage DHCPv6
}

func (op *OptRelayMsg) Code() OptionCode {
	return OptionRelayMsg
}

func (op *OptRelayMsg) ToBytes() []byte {
	return op.relayMessage.ToBytes()
}

func (op *OptRelayMsg) RelayMessage() DHCPv6 {
	return op.relayMessage
}

func (op *OptRelayMsg) SetRelayMessage(relayMessage DHCPv6) {
	op.relayMessage = relayMessage
}

func (op *OptRelayMsg) String() string {
	return fmt.Sprintf("OptRelayMsg{relaymsg=%v}", op.relayMessage)
}

// build an OptRelayMsg structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRelayMsg(data []byte) (*OptRelayMsg, error) {
	var err error
	var opt OptRelayMsg
	opt.relayMessage, err = FromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
