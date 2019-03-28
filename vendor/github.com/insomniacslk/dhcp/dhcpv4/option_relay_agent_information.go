package dhcpv4

import (
	"fmt"
)

// RelayOptions is like Options, but stringifies using the Relay Agent Specific
// option space.
type RelayOptions struct {
	Options
}

var relayHumanizer = OptionHumanizer{
	ValueHumanizer: func(code OptionCode, data []byte) fmt.Stringer {
		return OptionGeneric{data}
	},
	CodeHumanizer: func(c uint8) OptionCode {
		return GenericOptionCode(c)
	},
}

// String prints the contained options using Relay Agent-specific option code parsing.
func (r RelayOptions) String() string {
	return r.Options.ToString(relayHumanizer)
}

// FromBytes parses relay agent options from data.
func (r *RelayOptions) FromBytes(data []byte) error {
	r.Options = make(Options)
	return r.Options.FromBytes(data)
}

// OptRelayAgentInfo returns a new DHCP Relay Agent Info option.
//
// The relay agent info option is described by RFC 3046.
func OptRelayAgentInfo(o ...Option) Option {
	return Option{Code: OptionRelayAgentInformation, Value: RelayOptions{OptionsFromList(o...)}}
}

type raiSubOptionCode uint8

func (o raiSubOptionCode) Code() uint8 {
	return uint8(o)
}

func (o raiSubOptionCode) String() string {
	if s, ok := raiSubOptionCodeToString[o]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", o)
}

// AgentCircuitIDSubOption as per https://tools.ietf.org/html/rfc3046#section-2.0
const AgentCircuitIDSubOption raiSubOptionCode = 1

var raiSubOptionCodeToString = map[raiSubOptionCode]string{
	AgentCircuitIDSubOption: "Agent Circuit ID Sub-option",
}
