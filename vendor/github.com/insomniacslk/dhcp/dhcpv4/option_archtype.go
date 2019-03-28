package dhcpv4

import (
	"github.com/insomniacslk/dhcp/iana"
)

// OptClientArch returns a new Client System Architecture Type option.
func OptClientArch(archs ...iana.Arch) Option {
	return Option{Code: OptionClientSystemArchitectureType, Value: iana.Archs(archs)}
}
