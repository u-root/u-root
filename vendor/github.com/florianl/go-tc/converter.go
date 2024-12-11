package tc

import "net"

func netIPPtr(v net.IP) *net.IP {
	return &v
}

func netIPValue(v *net.IP) net.IP {
	if v != nil {
		return *v
	}
	return net.IP{}
}

func netHardwareAddrPtr(v net.HardwareAddr) *net.HardwareAddr {
	return &v
}

func netHardwareAddrValue(v *net.HardwareAddr) net.HardwareAddr {
	if v != nil {
		return *v
	}
	return net.HardwareAddr{}
}

func bytesPtr(v []byte) *[]byte {
	return &v
}

func bytesValue(v *[]byte) []byte {
	if v != nil {
		return *v
	}
	return []byte{}
}

func boolPtr(v bool) *bool {
	return &v
}

func boolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

func stringPtr(v string) *string {
	return &v
}

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func uint8Ptr(v uint8) *uint8 {
	return &v
}

func uint8Value(v *uint8) uint8 {
	if v != nil {
		return *v
	}
	return 0
}

func uint16Ptr(v uint16) *uint16 {
	return &v
}

func uint16Value(v *uint16) uint16 {
	if v != nil {
		return *v
	}
	return 0
}

func uint32Ptr(v uint32) *uint32 {
	return &v
}

func uint32Value(v *uint32) uint32 {
	if v != nil {
		return *v
	}
	return 0
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func uint64Value(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

func int8Ptr(v int8) *int8 {
	return &v
}

func int8Value(v *int8) int8 {
	if v != nil {
		return *v
	}
	return 0
}

func int16Ptr(v int16) *int16 {
	return &v
}

func int16Value(v *int16) int16 {
	if v != nil {
		return *v
	}
	return 0
}

func int32Ptr(v int32) *int32 {
	return &v
}

func int32Value(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

func int64Ptr(v int64) *int64 {
	return &v
}

func int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}
