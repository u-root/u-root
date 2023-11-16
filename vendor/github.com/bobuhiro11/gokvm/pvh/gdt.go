package pvh

import "github.com/bobuhiro11/gokvm/kvm"

// For GDT details see arch/x86/include/asm/segment.h

func GdtEntry(flags uint16, base uint32, limit uint32) uint64 {
	return (uint64(base)&0xFF000000)<<(56-24) |
		(uint64(flags)&0x0000F0FF)<<40 |
		(uint64(limit)&0x000F0000)<<(48-16) |
		(uint64(base)&0x00FFFFFF)<<16 |
		(uint64(limit) & 0x0000FFFF)
}

func getBase(entry uint64) uint64 {
	return ((entry & 0xFF00000000000000) >> 32) | ((entry & 0x000000FF00000000) >> 16) | (entry&0x00000000FFFF0000)>>16
}

func getG(entry uint64) uint8 {
	return uint8((entry & 0x0080000000000000) >> 55)
}

func getDB(entry uint64) uint8 {
	return uint8((entry & 0x0040000000000000) >> 54)
}

func getL(entry uint64) uint8 {
	return uint8((entry & 0x0020000000000000) >> 53)
}

func getAVL(entry uint64) uint8 {
	return uint8((entry & 0x0010000000000000) >> 52)
}

func getP(entry uint64) uint8 {
	return uint8((entry & 0x0000800000000000) >> 47)
}

func getDPL(entry uint64) uint8 {
	return uint8((entry & 0x0000600000000000) >> 45)
}

func getS(entry uint64) uint8 {
	return uint8((entry & 0x0000100000000000) >> 44)
}

func getType(entry uint64) uint8 {
	return uint8((entry & 0x00000F0000000000) >> 40)
}

// Extract the segment limit from the GDT segment descriptor.
//
// In a segment descriptor, the limit field is 20 bits, so it can directly describe
// a range from 0 to 0xFFFFF (1MByte). When G flag is set (4-KByte page granularity) it
// scales the value in the limit field by a factor of 2^12 (4Kbytes), making the effective
// limit range from 0xFFF (4 KBytes) to 0xFFFF_FFFF (4 GBytes).
//
// However, the limit field in the VMCS definition is a 32 bit field, and the limit value is not
// automatically scaled using the G flag. This means that for a desired range of 4GB for a
// given segment, its limit must be specified as 0xFFFF_FFFF. Therefore the method of obtaining
// the limit from the GDT entry is not sufficient, since it only provides 20 bits when 32 bits
// are necessary. Fortunately, we can check if the G flag is set when extracting the limit since
// the full GDT entry is passed as an argument, and perform the scaling of the limit value to
// return the full 32 bit value.
//
// The scaling mentioned above is required when using PVH boot, since the guest boots in protected
// (32-bit) mode and must be able to access the entire 32-bit address space. It does not cause issues
// for the case of direct boot to 64-bit (long) mode, since in 64-bit mode the processor does not
// perform runtime limit checking on code or data segments.
func getLimit(entry uint64) uint32 {
	l := uint32(((((entry) & 0x000F000000000000) >> 32) | ((entry) & 0x000000000000FFFF)))
	g := getG(entry)

	switch g {
	case 0:
		return l
	default:
		return (l << 12) | 0xFFFF
	}
}

func SegmentFromGDT(entry uint64, tableIndex uint8) kvm.Segment {
	var unused uint8

	u := getP(entry)

	switch u {
	case 0:
		unused = 1
	default:
		unused = 0
	}

	return kvm.Segment{
		Base:     getBase(entry),
		Limit:    getLimit(entry),
		Selector: uint16(tableIndex) * 8,
		Typ:      getType(entry),
		Present:  getP(entry),
		DPL:      getDPL(entry),
		DB:       getDB(entry),
		S:        getS(entry),
		L:        getL(entry),
		G:        getG(entry),
		AVL:      getAVL(entry),
		Unusable: unused,
	}
}
