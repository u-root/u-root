package hwapi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"unsafe"
)

//VTdRegisters represents the IOMMIO space
type VTdRegisters struct {
	Version                                 uint32 // Architecture version supported by the implementation.
	Reserved1                               uint32 // Reserved
	Capabilities                            uint64 // Hardware reporting of capabilities.
	ExtendedCapabilities                    uint64 // Hardware reporting of extended capabilities.
	GlobalCommand                           uint32 // Register controlling general functions.
	GlobalStatus                            uint32 // Register reporting general status.
	RootTableAddress                        uint64 // Register to set up location of root table.
	ContextCommand                          uint64 // Register to manage context-entry cache.
	Reserved2                               uint32 // Reserved
	FaultStatus                             uint32 // Register to report Fault/Error status
	FaultEventControl                       uint32 // Interrupt control register for fault events.
	FaultEventData                          uint32 // Interrupt message data register for fault events.
	FaultEventAddress                       uint32 // Interrupt message address register for fault event messages.
	FaultEventUpperAddress                  uint32 // Interrupt message upper address register for fault event messages.
	Reserved3                               uint64 // Reserved
	Reserved4                               uint64 // Reserved
	AdvancedFaultLog                        uint64 // Register to configure and manage advanced fault logging.
	Reserved5                               uint32 // Reserved
	ProtectedMemoryEnable                   uint32 // Register to enable DMA-protected memory region(s).
	ProtectedLowMemoryBase                  uint32 // Register pointing to base of DMA-protected low memory region.
	ProtectedLowMemoryLimit                 uint32 // Register pointing to last address (limit) of the DMA-protected low memory region.
	ProtectedHighMemoryBase                 uint64 // Register pointing to base of DMA-protected high memory region.
	ProtectedHighMemoryLimit                uint64 // Register pointing to last address (limit) of the DMA-protected high memory region.
	InvalidationQueueHead                   uint64 // Offset to the invalidation queue entry that will be read next by hardware.
	InvalidationQueueTail                   uint64 // Offset to the invalidation queue entry that will be written next by software.
	InvalidationQueueAddress                uint64 // Base address of memory-resident invalidation queue.
	Reserved6                               uint32 // Reserved
	InvalidationCompletionStatus            uint32 // Register to indicate the completion of an Invalidation Wait Descriptor with IF=1.
	InvalidationCompletionEventControl      uint32 // Register to control Invalidation Queue Events
	InvalidationCompletionEventData         uint32 // Invalidation Queue Event message data register for Invalidation Queue events.
	InvalidationCompletionEventAddress      uint32 // Invalidation Queue Event message address register for Invalidation Queue events.
	InvalidationCompletionEventUpperAddress uint32 // Invalidation Queue Event message upper address register for Invalidation Queue events.
	Reserved7                               uint64 // Reserved.
	InterruptRemappingTableAddress          uint64 // Register indicating Base Address of Interrupt Remapping Table.
	PageRequestQueueHead                    uint64 // Offset to the page request queue entry that will be processed next by software.
	PageRequestQueueTail                    uint64 // Offset to the page request queue entry that will be written next by hardware.
	PageRequestQueueAddress                 uint64 // Base address of memory-resident page request queue.
	Reserved8                               uint32 // Reserved
	PageRequestStatus                       uint32 // Register to indicate one or more pending page requests in page request queue.
	PageRequestEventControl                 uint32 // Register to control page request events.
	PageRequestEventData                    uint32 // Page request event message data register.
	PageRequestEventAddress                 uint32 // Page request event message address register
	PageRequestEventUpperAddress            uint32 // Page request event message upper address register.
	MTRRCapability                          uint64 // Register for MTRR capability reporting.
	MTRRDefaultType                         uint64 // Register to configure MTRR default type.
	FixedRangeMTRR64K00000                  uint64 // Fixed-range memory type range register for 64K range starting at 00000h.
	FixedRangeMTRR16K80000                  uint64 // Fixed-range memory type range register for 16K range starting at 80000h.
	FixedRangeMTRR16KA0000                  uint64 // Fixed-range memory type range register for 16K range starting at A0000h.
	FixedRangeMTRR4KC0000                   uint64 // Fixed-range memory type range register for 4K range starting at C0000h.
	FixedRangeMTRR4KC8000                   uint64 // Fixed-range memory type range register for 4K range starting at C8000h.
	FixedRangeMTRR4KD0000                   uint64 // Fixed-range memory type range register for 4K range starting at D0000h.
	FixedRangeMTRR4KD8000                   uint64 // Fixed-range memory type range register for 4K range starting at D8000h.
	FixedRangeMTRR4KE0000                   uint64 // Fixed-range memory type range register for 4K range starting at E0000h.
	FixedRangeMTRR4KE8000                   uint64 // Fixed-range memory type range register for 4K range starting at E8000h.
	FixedRangeMTRR4KF0000                   uint64 // Fixed-range memory type range register for 4K range starting at F0000h.
	FixedRangeMTRR4KF8000                   uint64 // Fixed-range memory type range register for 4K range starting at F8000h.
	VariableRangeMTRRBase0                  uint64 // Variable-range memory type range0 base register.
	VariableRangeMTRRMask0                  uint64 // Variable-range memory type range0 mask register.
	VariableRangeMTRRBase1                  uint64 // Variable-range memory type range1 base register.
	VariableRangeMTRRMask1                  uint64 // Variable-range memory type range1 mask register.
	VariableRangeMTRRBase2                  uint64 // Variable-range memory type range2 base register.
	VariableRangeMTRRMask2                  uint64 // Variable-range memory type range2 mask register.
	VariableRangeMTRRBase3                  uint64 // Variable-range memory type range3 base register.
	VariableRangeMTRRMask3                  uint64 // Variable-range memory type range3 mask register.
	VariableRangeMTRRBase4                  uint64 // Variable-range memory type range4 base register.
	VariableRangeMTRRMask4                  uint64 // Variable-range memory type range4 mask register.
	VariableRangeMTRRBase5                  uint64 // Variable-range memory type range5 base register.
	VariableRangeMTRRMask5                  uint64 // Variable-range memory type range5 mask register.
	VariableRangeMTRRBase6                  uint64 // Variable-range memory type range6 base register.
	VariableRangeMTRRMask6                  uint64 // Variable-range memory type range6 mask register.
	VariableRangeMTRRBase7                  uint64 // Variable-range memory type range7 base register.
	VariableRangeMTRRMask7                  uint64 // Variable-range memory type range7 mask register.
	VariableRangeMTRRBase8                  uint64 // Variable-range memory type range8 base register.
	VariableRangeMTRRMask8                  uint64 // Variable-range memory type range8 mask register.
	VariableRangeMTRRBase9                  uint64 // Variable-range memory type range9 base register.
	VariableRangeMTRRMask9                  uint64 // Variable-range memory type range9 mask register.
	VirtualCommandCapability                uint64 // Hardware reporting of commands supported by virtual-DMA Remapping hardware.
	Reserved10                              uint64 // Reserved for future expansion of Virtual Command Capability Register.
	VirtualCommand                          uint64 // Register to submit commands to virtual DMA Remapping hardware.
	Reserved11                              uint64 // Reserved for future expansion of Virtual Command Register.
	VirtualCommandResponse                  uint64 // Register to receive responses from virtual DMA Remapping hardware.
	Reserved12                              uint64 // Reserved for future expansion of Virtual Command Response Register.
}

func (t TxtAPI) readVTdRegs() (VTdRegisters, error) {
	var regs VTdRegisters

	dir, err := os.Open("/sys/class/iommu/")
	if err != nil {
		return regs, fmt.Errorf("No IOMMU found: %s", err)
	}

	subdirs, err := dir.Readdir(0)
	if err != nil {
		return regs, fmt.Errorf("No IOMMU found: %s", err)
	}

	for _, subdir := range subdirs {
		path := fmt.Sprintf("/sys/class/iommu/%s/intel-iommu/address", subdir.Name())
		addrBuf, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		addr, err := strconv.ParseUint(string(addrBuf[:len(addrBuf)-1]), 16, 64)
		if err != nil {
			continue
		}

		buf := make([]byte, unsafe.Sizeof(regs))
		err = t.ReadPhysBuf(int64(addr), buf)
		if err != nil {
			continue
		}

		reader := bytes.NewReader(buf)
		err = binary.Read(reader, binary.LittleEndian, &regs)
		if err != nil {
			continue
		}

		return regs, nil
	}

	return regs, fmt.Errorf("No IOMMU found: /sys/class/iommu/*/intel-iommu/address does not exists or is malformed")
}

//LookupIOAddress returns the address of the root Tbl
func (t TxtAPI) LookupIOAddress(addr uint64, regs VTdRegisters) ([]uint64, error) {
	rootTblAddr := regs.RootTableAddress & 0xffffffffffff000
	ttm := (regs.RootTableAddress >> 10) & 3

	if ttm == 0 {
		return t.lookupIOLegacy(addr, rootTblAddr)
	} else if ttm == 1 {
		return lookupIOScalable(addr, rootTblAddr)
	} else {
		return []uint64{}, fmt.Errorf("unsupported IOMMU Translation Table Mode")
	}
}

func (t TxtAPI) lookupIOLegacy(addr, rootTblAddr uint64) ([]uint64, error) {
	ret := []uint64{}

	for bus := int64(0); bus < 256; bus++ {
		// read root table entries
		var rootTblEnt Uint64
		err := t.ReadPhys(int64(rootTblAddr)+bus*16+8, &rootTblEnt)
		if err != nil {
			return []uint64{}, err
		}

		if rootTblEnt&1 == 0 {
			continue
		}

		// we assume 48 bit physical addresses
		ctxTblAddr := rootTblEnt & 0xffffffffffff000

		// read ctx entry
		var ctxTblEntHi Uint64
		var ctxTblEntLo Uint64

		for devfn := int64(0); devfn < 32*8; devfn++ {
			err = t.ReadPhys(int64(ctxTblAddr)+devfn*16+8, &ctxTblEntLo)
			if err != nil {
				return []uint64{}, err
			}

			if ctxTblEntLo&1 == 0 {
				continue
			}

			err = t.ReadPhys(int64(ctxTblAddr)+devfn*16, &ctxTblEntHi)
			if err != nil {
				return []uint64{}, err
			}

			tt := (ctxTblEntLo >> 2) & 3
			aw := (ctxTblEntHi >> 2) & 7
			pml5Addr := ctxTblEntLo & 0xffffffffffff000

			if tt == 2 {
				ret = append(ret, addr)
				continue
			}

			var pdptAddr Uint64
			var canRead bool

			if aw == 3 || aw == 2 {
				// Page map level 5
				var l5ent Uint64
				err = t.ReadPhys(int64(pml5Addr)+(int64(addr>>48)&0x1ff)*8, &l5ent)
				if err != nil {
					return []uint64{}, err
				}

				pml4Addr := l5ent & 0xffffffffffff000
				canRead = l5ent&1 != 0

				if !canRead {
					continue
				}

				if aw == 2 {
					// Page map level 4
					var l4ent Uint64
					err = t.ReadPhys(int64(pml4Addr)+(int64(addr>>39)&0x1ff)*8, &l4ent)
					if err != nil {
						return []uint64{}, err
					}

					pdptAddr = l4ent & 0xffffffffffff000
					canRead = l4ent&1 != 0

					if !canRead {
						continue
					}
				} else {
					pdptAddr = pml4Addr
				}
			} else {
				pdptAddr = pml5Addr
			}

			// Page directory pointer table
			var l3ent Uint64
			err = t.ReadPhys(int64(pdptAddr)+(int64(addr>>30)&0x1ff)*8, &l3ent)
			if err != nil {
				return []uint64{}, err
			}

			pdAddr := l3ent & 0xffffffffffff000
			canRead = l3ent&1 != 0

			if !canRead {
				continue
			}

			// Page directory
			var l2ent Uint64
			err = t.ReadPhys(int64(pdAddr)+(int64(addr>>21)&0x1ff)*8, &l2ent)
			if err != nil {
				return []uint64{}, err
			}

			ptAddr := l2ent & 0xffffffffffff000
			canRead = l2ent&1 != 0

			if !canRead {
				continue
			}

			// Page table
			var l1ent Uint64
			err = t.ReadPhys(int64(ptAddr)+(int64(addr>>12)&0x1ff)*8, &l1ent)
			if err != nil {
				return []uint64{}, err
			}

			pageAddr := l1ent & 0xffffffffffff000
			canRead = l1ent&1 != 0

			if !canRead {
				continue
			}

			ret = append(ret, uint64(pageAddr)+(addr&0xffff))
		}
	}

	return ret, nil
}

func lookupIOScalable(addr, rootTblAddr uint64) ([]uint64, error) {
	return []uint64{}, fmt.Errorf("Scalable IOMMU not implemented")
	// read root table entry

	// read ctx entry
	// get page table depth

	// read req-w/o-PASID directory pointer (5 level)
	// make sure 2-pass translation isnt on
}

//AddressRangesIsDMAProtected returns true if the address is DMA protected by the IOMMU
func (t TxtAPI) AddressRangesIsDMAProtected(first, end uint64) (bool, error) {
	regs, err := t.readVTdRegs()
	if err != nil {
		return false, err
	}

	loDMAprotection := regs.Capabilities&(1<<5) != 0
	hiDMAprotection := regs.Capabilities&(1<<6) != 0
	enableDMAprotection := regs.Capabilities&1 != 0
	enable2DMAprotection := regs.Capabilities&(1<<31) != 0

	if enableDMAprotection && enable2DMAprotection && loDMAprotection && uint64(regs.ProtectedLowMemoryBase) <= first && uint64(regs.ProtectedLowMemoryLimit) >= end {
		return true, err
	}

	if enableDMAprotection && enable2DMAprotection && hiDMAprotection && regs.ProtectedHighMemoryBase <= first && regs.ProtectedHighMemoryLimit >= end {
		return true, err
	}

	for addr := first & 0xffffffffffff0000; addr < end; addr += 4096 {
		vas, err := t.LookupIOAddress(addr, regs)
		if err != nil {
			return false, err
		}

		if len(vas) < 0 {
			return false, nil
		}
	}

	return true, nil
}
