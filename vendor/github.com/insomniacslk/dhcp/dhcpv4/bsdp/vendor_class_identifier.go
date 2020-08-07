// +build !darwin

package bsdp

// MakeVendorClassIdentifier returns a static vendor class identifier for BSDP
// use on non-darwin hosts.
func MakeVendorClassIdentifier() (string, error) {
	return DefaultMacOSVendorClassIdentifier, nil
}
