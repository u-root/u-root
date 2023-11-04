package pvh

const (
	/*
		Start low ram range
	*/
	// LowRAMStart (start: 0, length: 640KiB).
	LowRAMStart = 0x0

	// Location of EBDA address.
	EBDAPointer = 0x40e

	// Initial GDT/IDT.
	BootGDTStart = 0x500
	BootIDTStart = 0x520

	// Address of the pvh_info struct.
	PVHInfoStart = 0x6000

	// Address of hvm_modlist_entry type.
	PVHModlistStart = 0x6040

	// Address of memory map table for PVH boot.
	PVHMemMapStart = 0x7000

	// Kernel command line.
	KernelCmdLine        = 0x2_0000
	KernelCmdLineSizeMax = 0x1_0000

	// MPTable describing vcpus.
	MPTableStart = 0x9_FC00

	/*
		End low ram range.
	*/

	// EDBA reserved area (start: 640KiB, length: 384KiB).
	EBDAStart = 0xA_0000

	// RSDPPointer in EDBA area.
	RSDPPointer = EBDAStart

	// SMBIOSStart first location possible for SMBIOS.
	SMBIOSStart = 0xF_0000

	/*
		Start high ram range.
	*/

	// HighRAMStart (start: 1MiB, length: 3071MiB).
	HighRAMStart = 0x10_0000

	// 32Bit reserved area (start: 3GiB, length: 896MiB).
	Mem32BitReservedStart = 0xC000_0000
	Mem32BitReservedSize  = PCIMMConfigSize + Mem32BitDeviceSize

	Mem32BitDeviceStart = Mem32BitReservedStart
	Mem32BitDeviceSize  = 640 << 20

	// PCI Memory Mapped Config Space.
	PCIMMConfigStart            = Mem32BitDeviceStart + Mem32BitDeviceSize
	PCIMMConfigSize             = 256 << 20
	PCIMMIOConfigSizePerSegment = 4096 * 256

	// TSS is 3 page after PCI MMConfig space.
	KVMTSSStart = PCIMMConfigStart + PCIMMConfigSize
	KVMTSSSize  = (3 * 4) << 10

	// Identity map is one page region after TSS.
	KVMIdentityMapStart = KVMTSSStart + KVMTSSSize
	KVMIdentityMapSize  = 4 << 10

	// IOAPIC.
	IOAPICStart = 0xFEC0_0000
	IOAPICSize  = 0x20

	// APIC.
	APICStart = 0xFEE0_0000

	// 64bit address space start.
	RAM64BitStart = 0x1_0000_0000
)

const (
	// Reserve 1 MiB for platform MMIO devices (e.g. ACPI control devices).
	PlatformDeviceAreaSize = 1 << 20
)
