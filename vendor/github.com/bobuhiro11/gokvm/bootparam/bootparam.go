package bootparam

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

const (
	MagicSignature = 0x53726448

	LoadedHigh   = uint8(1 << 0)
	KeepSegments = uint8(1 << 6)
	CanUseHeap   = uint8(1 << 7)

	EddMbrSigMax = 16
	E820Max      = 128
	E820Ram      = 1
	E820Reserved = 2

	RealModeIvtBegin = 0x00000000
	EBDAStart        = 0x0009fc00
	VGARAMBegin      = 0x000a0000
	MBBIOSBegin      = 0x000f0000
	MBBIOSEnd        = 0x000fffff
)

type E820Entry struct {
	Addr uint64
	Size uint64
	Type uint32
}

// The so-called "zeropage"
// https://www.kernel.org/doc/html/latest/x86/boot.html
// https://github.com/torvalds/linux/blob/master/arch/x86/include/uapi/asm/bootparam.h
type BootParam struct {
	Padding             [0x1e8]uint8
	E820Entries         uint8
	EddbufEntries       uint8
	EddMbrSigBufEntries uint8
	KdbStatus           uint8
	Padding2            [5]uint8
	Hdr                 SetupHeader
	Padding3            [0x290 - 0x1f1 - unsafe.Sizeof(SetupHeader{})]uint8

	// Required to adjust the offset of E820Map to 0x2D0.
	Padding4 [0x3d]uint8

	EddMbrSigBuffer [EddMbrSigMax]uint8
	E820Map         [E820Max]E820Entry
}

type SetupHeader struct {
	SetupSects          uint8
	RootFlags           uint16
	SysSize             uint32
	RAMSize             uint16
	VidMode             uint16
	RootDev             uint16
	BootFlag            uint16
	Jump                uint16
	Header              uint32
	Version             uint16
	ReadModeSwitch      uint32
	StartSysSeg         uint16
	KernelVersion       uint16
	TypeOfLoader        uint8
	LoadFlags           uint8
	SetupMoveSize       uint16
	Code32Start         uint32
	RamdiskImage        uint32
	RamdiskSize         uint32
	BootsectKludge      uint32
	HeapEndPtr          uint16
	ExtLoaderVer        uint8
	ExtLoaderType       uint8
	CmdlinePtr          uint32
	InitrdAddrMax       uint32
	KernelAlignment     uint32
	RelocatableKernel   uint8
	MinAlignment        uint8
	XloadFlags          uint16
	CmdlineSize         uint32
	HardwareSubarch     uint32
	HardwareSubarchData uint64
	PayloadOffset       uint32
	PayloadLength       uint32
	SetupData           uint64
	PrefAddress         uint64
	InitSize            uint32
	HandoverOffset      uint32
	KernelInfoOffset    uint32
}

var ErrorSignatureNotMatch = errors.New("signature not match in bzImage")

var ErrorOldProtocolVersion = errors.New("old protocol version")

func New(r io.ReaderAt) (*BootParam, error) {
	b := &BootParam{}

	// In 64-bit boot protocol, the first step in loading a Linux kernel should be
	// to setup the boot parameters (struct boot_params, traditionally known as
	// "zero page"). The memory for struct boot_params could be allocated anywhere
	// (even above 4G) and initialized to all zero. Then, the setup header at
	// offset 0x01f1 of kernel image on should be loaded into struct boot_params
	// and examined.
	//
	// refs: https://www.kernel.org/doc/html/latest/x86/boot.html#id1
	reader := io.NewSectionReader(r, 0x1f1, 0x1000)
	if err := binary.Read(reader, binary.LittleEndian, &(b.Hdr)); err != nil {
		return b, err
	}

	if err := b.isValid(); err != nil {
		return b, err
	}

	return b, nil
}

func (b *BootParam) isValid() error {
	if b.Hdr.Header != MagicSignature {
		return ErrorSignatureNotMatch
	}

	// Protocol 2.06+ is required.
	if b.Hdr.Version < 0x0206 {
		return fmt.Errorf("%w: 0x%x", ErrorOldProtocolVersion, b.Hdr.Version)
	}

	return nil
}

func (b *BootParam) AddE820Entry(addr, size uint64, typ uint32) {
	i := b.E820Entries
	b.E820Map[i] = E820Entry{
		Addr: addr,
		Size: size,
		Type: typ,
	}
	b.E820Entries = i + 1
}

func (b *BootParam) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, b); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
