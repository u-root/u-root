// +build !linux

package dhcp6

import (
	"net"
)

// HardwareType returns ErrHardwareTypeNotImplemented, because it is not
// implemented on non-Linux platforms.
func HardwareType(ifi *net.Interface) (uint16, error) {
	return 0, ErrHardwareTypeNotImplemented
}
