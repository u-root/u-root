// +build !linux

package dhcp6

import (
	"testing"
)

// TestHardwareTypeOthers verifies that HardwareType always returns
// ErrHardwareTypeNotImplemented on platforms that do not have an
// implementation written.
func TestHardwareTypeOthers(t *testing.T) {
	if _, err := HardwareType(nil); err != ErrHardwareTypeNotImplemented {
		t.Fatalf("unexpected error, should be ErrHardwareTypeNotImplemented: %v", err)
	}
}
