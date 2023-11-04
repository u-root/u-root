package pvh

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"

	"github.com/bobuhiro11/gokvm/kvm"
)

const (
	xenHVMstartMagicValue uint32 = 0x336ec578
	xenELFNotePhys32Entry uint32 = 18
	pvhNoteStrSz          uint32 = 4
	elfNoteSize                  = 12
)

var (
	errAlign            = errors.New("alignment is not a power of 2")
	errPVHEntryNotFound = errors.New("no pvh entry found")
)

type HVMStartInfo struct {
	Magic         uint32
	Version       uint32
	Flags         uint32
	NrModules     uint32
	ModlistPAddr  uint64
	CmdLinePAddr  uint64
	RSDPPAddr     uint64
	MemMapPAddr   uint64
	MemMapEntries uint32
	_             uint32
}

func NewStartInfo(rsdpPAddr, cmdLinePAddr uint64) *HVMStartInfo {
	return &HVMStartInfo{
		Magic:        xenHVMstartMagicValue,
		Version:      1,
		NrModules:    0,
		CmdLinePAddr: cmdLinePAddr,
		RSDPPAddr:    rsdpPAddr,
		MemMapPAddr:  PVHMemMapStart,
	}
}

func (h *HVMStartInfo) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	for _, item := range []interface{}{
		h.Magic,
		h.Version,
		h.Flags,
		h.NrModules,
		h.ModlistPAddr,
		h.CmdLinePAddr,
		h.RSDPPAddr,
		h.MemMapPAddr,
		h.MemMapEntries,
		uint32(0x0),
	} {
		if err := binary.Write(&buf, binary.LittleEndian, item); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

type HVMModListEntry struct {
	Addr        uint64
	Size        uint64
	CmdLineAddr uint64
	_           uint64
}

func NewModListEntry(addr, size, cmdaddr uint64) *HVMModListEntry {
	return &HVMModListEntry{
		Addr:        addr,
		Size:        size,
		CmdLineAddr: cmdaddr,
	}
}

func (h *HVMModListEntry) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	for _, item := range []interface{}{
		h.Addr,
		h.Size,
		h.CmdLineAddr,
		uint64(0x0),
	} {
		if err := binary.Write(&buf, binary.LittleEndian, item); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

type HVMMemMapTableEntry struct {
	Addr uint64
	Size uint64
	Type uint32
	_    uint32
}

func NewMemMapTableEntry(addr, size uint64, t uint32) *HVMMemMapTableEntry {
	return &HVMMemMapTableEntry{
		Addr: addr,
		Size: size,
		Type: t,
	}
}

func (h *HVMMemMapTableEntry) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	for _, item := range []interface{}{
		h.Addr,
		h.Size,
		h.Type,
		uint32(0x0),
	} {
		if err := binary.Write(&buf, binary.LittleEndian, item); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

type GDT [4]uint64

func CreateGDT() GDT {
	var gdtTable GDT

	gdtTable[0] = GdtEntry(0, 0, 0)               // NULL
	gdtTable[1] = GdtEntry(0xc09b, 0, 0xffffffff) // Code
	gdtTable[2] = GdtEntry(0xc093, 0, 0xffffffff) // DATA
	gdtTable[3] = GdtEntry(0x008b, 0, 0x67)       // TSS

	return gdtTable
}

func InitSRegs(vcpuFd uintptr, gdttable GDT) error {
	codeseg := SegmentFromGDT(gdttable[1], 1)
	dataseg := SegmentFromGDT(gdttable[2], 2)
	tssseg := SegmentFromGDT(gdttable[3], 3)

	// We need to write this to ....maybe create this config earlier.
	gdt := kvm.Descriptor{
		Base:  BootGDTStart,
		Limit: uint16(len(gdttable)*8) - 1, // 4 entries of 64bit (8byte) per entry
	}

	idt := kvm.Descriptor{
		Base:  BootIDTStart,
		Limit: 8,
	}

	sregs, err := kvm.GetSregs(vcpuFd)
	if err != nil {
		return err
	}

	sregs.GDT = gdt
	sregs.IDT = idt

	sregs.CS = codeseg
	sregs.DS = dataseg
	sregs.ES = dataseg
	sregs.FS = dataseg
	sregs.GS = dataseg
	sregs.SS = dataseg
	sregs.TR = tssseg

	sregs.EFER |= (0 << 17) | (0 << 9) | (0 << 8) // VM=0, IF=0, TF=0

	sregs.CR0 = 0x1
	sregs.CR4 = 0x0

	return kvm.SetSregs(vcpuFd, sregs)
}

func (gdt GDT) Bytes() []byte {
	bytes := make([]byte, binary.Size(gdt))

	for i, entry := range gdt {
		binary.LittleEndian.PutUint64(bytes[i*binary.Size(entry):], entry)
	}

	return bytes
}

func InitRegs(vcpuFd uintptr, bootIP uint64) error {
	regs, err := kvm.GetRegs(vcpuFd)
	if err != nil {
		return err
	}

	regs.RFLAGS = 0x2
	regs.RBX = PVHInfoStart
	regs.RIP = bootIP

	return kvm.SetRegs(vcpuFd, regs)
}

type elfNote struct {
	NameSize uint32
	DescSize uint32
	Type     uint32
}

func ParsePVHEntry(fwimg io.ReaderAt, phdr *elf.Prog) (uint32, error) {
	node := elfNote{}
	off := int64(phdr.Off)
	readSize := 0

	for readSize < int(phdr.Filesz) {
		nodeByte := make([]byte, 12)

		n, err := fwimg.ReadAt(nodeByte, off)
		if err != nil {
			return 0x0, err
		}

		readSize += n
		off += int64(n)

		nsb := make([]byte, 4)
		dsb := make([]byte, 4)
		tsb := make([]byte, 4)

		copy(nsb, nodeByte[:3])
		copy(dsb, nodeByte[4:7])
		copy(tsb, nodeByte[8:])

		node.NameSize = binary.LittleEndian.Uint32(nsb)
		node.DescSize = binary.LittleEndian.Uint32(dsb)
		node.Type = binary.LittleEndian.Uint32(tsb)

		if node.Type == xenELFNotePhys32Entry && node.NameSize == pvhNoteStrSz {
			buf := make([]byte, pvhNoteStrSz)

			n, err := fwimg.ReadAt(buf, off)
			if err != nil {
				return 0x0, err
			}

			off += int64(n)
			// Check the String
			if bytes.Equal(buf, []byte{'X', 'e', 'n', '\000'}) {
				break
			}
		}

		nameAlign, err := alignUp(uint64(node.NameSize))
		if err != nil {
			return 0x0, err
		}

		descAlign, err := alignUp(uint64(node.DescSize))
		if err != nil {
			return 0x0, err
		}

		readSize += int(nameAlign)
		readSize += int(descAlign)
		off = int64(phdr.Off) + int64(readSize)
	}

	if readSize >= int(phdr.Filesz) {
		// No PVH entry found. Return
		return 0x0, errPVHEntryNotFound
	}

	// off is the value we need to add aligned namesize - PVH_NOTE_STR_SZ
	nameAlign, err := alignUp(uint64(node.NameSize))
	if err != nil {
		return 0x0, err
	}

	off += (int64(nameAlign) - int64(pvhNoteStrSz))
	pvhAddrByte := make([]byte, 4) // address is 4 byte/32-bit

	if _, err := fwimg.ReadAt(pvhAddrByte, off); err != nil {
		return 0x0, err
	}

	retAddr := binary.LittleEndian.Uint32(pvhAddrByte)

	return retAddr, nil
}

func alignUp(addr uint64) (uint64, error) {
	align := uint64(4)
	if !isPowerOf2(align) {
		return addr, errAlign
	}

	alignMask := align - 1

	if addr&alignMask == 0 {
		return addr, nil
	}

	return (addr | alignMask) + 1, nil
}

func isPowerOf2(n uint64) bool {
	if n == 0 {
		return true
	}

	return (n & (n - 1)) == 0
}

func CheckPVH(kern io.ReaderAt) (bool, error) {
	elfkern, err := elf.NewFile(kern)
	if err != nil {
		return false, nil //nolint:nilerr
	}
	defer elfkern.Close()

	for _, prog := range elfkern.Progs {
		note := elfNote{}
		off := int64(prog.Off)
		readSize := 0

		for readSize < int(prog.Filesz) {
			noteByte := make([]byte, elfNoteSize)

			n, err := kern.ReadAt(noteByte, off)
			if err != nil {
				return false, err
			}

			readSize += n
			off += int64(n)

			nsb := make([]byte, 4)
			dsb := make([]byte, 4)
			tsb := make([]byte, 4)

			copy(nsb, noteByte[:3])
			copy(dsb, noteByte[4:7])
			copy(tsb, noteByte[8:])

			note.NameSize = binary.LittleEndian.Uint32(nsb)
			note.DescSize = binary.LittleEndian.Uint32(dsb)
			note.Type = binary.LittleEndian.Uint32(tsb)

			if note.Type == xenELFNotePhys32Entry && note.NameSize == pvhNoteStrSz {
				buf := make([]byte, pvhNoteStrSz)

				_, err := kern.ReadAt(buf, off)
				if err != nil {
					return false, err
				}

				if bytes.Equal(buf, []byte{'X', 'e', 'n', '\000'}) {
					return true, nil
				}
			}

			nameAlign, err := alignUp(uint64(note.NameSize))
			if err != nil {
				return false, err
			}

			descAlign, err := alignUp(uint64(note.DescSize))
			if err != nil {
				return false, err
			}

			readSize += int(nameAlign)
			readSize += int(descAlign)
			off = int64(prog.Off) + int64(readSize)
		}
	}

	return false, nil
}
