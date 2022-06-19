// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zbi contains a parser for the Zircon boot image format.
package zbi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/align"
)

const (
	ContainerMagic      uint32 = 0x868cf7e6
	ItemMagic           uint32 = 0xb5781729
	VersionFlag         uint32 = 0x00010000
	CRC32Flag           uint32 = 0x00020000
	NoCRC32Flag         uint32 = 0x4a87e8d6
	ZBITypeKernelPrefix uint32 = 0x004e524b // KRN\0
	ZBITypeKernelMask   uint32 = 0x00FFFFFF // Mask to compare to the prefix.
)

type ZBITypeMetadata struct {
	Name      string
	Extention string
}
type ZBIType uint32

const (
	ZBITypeContainer            ZBIType = 0x544f4f42 // BOOT
	ZBITypeKernelX64            ZBIType = 0x4c4e524b // KRNL
	ZBITypeKernelArm64          ZBIType = 0x384e524b // KRN8
	ZBITypeDiscard              ZBIType = 0x50494b53 // SKIP
	ZBITypeStorageRamdisk       ZBIType = 0x4b534452 // RDSK
	ZBITypeStorageBootfs        ZBIType = 0x42534642 // BFSB
	ZBITypeStorageKernel        ZBIType = 0x5254534b // KSTR
	ZBITypeStorageBootfsFactory ZBIType = 0x46534642 // BFSF
	ZBITypeCmdline              ZBIType = 0x4c444d43 // CMDL
	ZBITypeCrashlog             ZBIType = 0x4d4f4f42 // BOOM
	ZBITypeNvram                ZBIType = 0x4c4c564e // NVLL
	ZBITypePlatformID           ZBIType = 0x44494C50 // PLID
	ZBITypeDrvBoardInfo         ZBIType = 0x4953426D // mBSI
	ZBITypeCPUConfig            ZBIType = 0x43555043 // CPUC
	ZBITypeCPUTopology          ZBIType = 0x544F504F // TOPO
	ZBITypeMemConfig            ZBIType = 0x434D454D // MEMC
	ZBITypeKernelDriver         ZBIType = 0x5652444B // KDRV
	ZBITypeAcpiRsdp             ZBIType = 0x50445352 // RSDP
	ZBITypeSMBios               ZBIType = 0x49424d53 // SMBI
	ZBITypeEFISystemTable       ZBIType = 0x53494645 // EFIS
	ZBITypeFramebuffer          ZBIType = 0x42465753 // SWFB
	ZBITypeImageArgs            ZBIType = 0x47524149 // IARG
	ZBITypeBootVersion          ZBIType = 0x53525642 // BVRS
	ZBITypeDrvMacAddress        ZBIType = 0x43414D6D // mMAC
	ZBITypeDrvPartitionMap      ZBIType = 0x5452506D // mPRT
	ZBITypeDrvBoardPrivate      ZBIType = 0x524F426D // mBOR
	ZBITypeHwRebootReason       ZBIType = 0x42525748 // HWRB
	ZBITypeSerialNumber         ZBIType = 0x4e4c5253 // SRLN
	ZBITypeBootloaderFile       ZBIType = 0x4C465442 // BTFL
	ZBITypeDevicetree           ZBIType = 0xd00dfeed
	ZBITypeSecureEntropy        ZBIType = 0x444e4152 // RAND
)

var (
	ZBITypes = map[ZBIType]ZBITypeMetadata{
		ZBITypeContainer:            {Name: "CONTAINER", Extention: ".bin"},
		ZBITypeKernelX64:            {Name: "KERNEL_X64", Extention: ".bin"},
		ZBITypeKernelArm64:          {Name: "KERNEL_ARM64", Extention: ".bin"},
		ZBITypeDiscard:              {Name: "DISCARD", Extention: ".bin"},
		ZBITypeStorageKernel:        {Name: "KERNEL", Extention: ".bin"},
		ZBITypeStorageRamdisk:       {Name: "RAMDISK", Extention: ".bin"},
		ZBITypeStorageBootfs:        {Name: "BOOTFS", Extention: ".bin"},
		ZBITypeStorageBootfsFactory: {Name: "BOOTFS_FACTORY", Extention: ".bin"},
		ZBITypeCmdline:              {Name: "CMDLINE", Extention: ".txt"},
		ZBITypeCrashlog:             {Name: "CRASHLOG", Extention: ".bin"},
		ZBITypeNvram:                {Name: "NVRAM", Extention: ".bin"},
		ZBITypePlatformID:           {Name: "PLATFORM_ID", Extention: ".bin"},
		ZBITypeCPUConfig:            {Name: "CPU_CONFIG", Extention: ".bin"},
		ZBITypeCPUTopology:          {Name: "CPU_TOPOLOGY", Extention: ".bin"},
		ZBITypeMemConfig:            {Name: "MEM_CONFIG", Extention: ".bin"},
		ZBITypeKernelDriver:         {Name: "KERNEL_DRIVER", Extention: ".bin"},
		ZBITypeAcpiRsdp:             {Name: "ACPI_RSDP", Extention: ".bin"},
		ZBITypeSMBios:               {Name: "SMBIOS", Extention: ".bin"},
		ZBITypeEFISystemTable:       {Name: "EFI_SYSTEM_TABLE", Extention: ".bin"},
		ZBITypeFramebuffer:          {Name: "FRAMEBUFFER", Extention: ".bin"},
		ZBITypeImageArgs:            {Name: "IMAGE_ARGS", Extention: ".txt"},
		ZBITypeBootVersion:          {Name: "BOOT_VERSION", Extention: ".bin"},
		ZBITypeDrvBoardInfo:         {Name: "DRV_BOARD_INFO", Extention: ".bin"},
		ZBITypeDrvMacAddress:        {Name: "DRV_MAC_ADDRESS", Extention: ".bin"},
		ZBITypeDrvPartitionMap:      {Name: "DRV_PARTITION_MAP", Extention: ""},
		ZBITypeDrvBoardPrivate:      {Name: "DRV_BOARD_PRIVATE", Extention: ""},
		ZBITypeHwRebootReason:       {Name: "HW_REBOOT_REASON", Extention: ".bin"},
		ZBITypeSerialNumber:         {Name: "SERIAL_NUMBER", Extention: ".txt"},
		ZBITypeBootloaderFile:       {Name: "BOOTLOADER_FILE", Extention: ".bin"},
		ZBITypeDevicetree:           {Name: "DEVICETREE", Extention: ".dtb"},
		ZBITypeSecureEntropy:        {Name: "ENTROPY", Extention: ".bin"},
	}
)

func (it *ZBIType) IsKernel() bool {
	return uint32(*it)&ZBITypeKernelMask == ZBITypeKernelPrefix
}

func (it *ZBIType) IsDriverMetadata() bool {
	return *it&0xFF == 0x6D // 'm'
}

func (it *ZBIType) ToString() (string, error) {
	if typeMetadata, ok := ZBITypes[*it]; ok {
		return typeMetadata.Name, nil
	}
	return "", fmt.Errorf("Can't find metadata for %#08x ZBIType", it)
}

func (it *ZBIType) MarshalJSON() ([]byte, error) {
	name, err := it.ToString()
	if err != nil {
		return nil, err
	}
	nameWithQuotes := fmt.Sprintf("%q", name)
	return []byte(nameWithQuotes), nil
}

type Header struct {
	Type      ZBIType
	Length    uint32
	Extra     uint32
	Flags     uint32
	Reserved0 uint32
	Reserved1 uint32
	Magic     uint32
	CRC32     uint32
}

type BootItem struct {
	Header         Header
	PayloadAddress uint64
}

func NewContainerHeader(length uint32) Header {
	return Header{
		Type:      ZBITypeContainer,
		Length:    length,
		Extra:     ContainerMagic,
		Flags:     VersionFlag,
		Reserved0: 0,
		Reserved1: 0,
		Magic:     ItemMagic,
		CRC32:     NoCRC32Flag,
	}
}

type Image struct {
	Header    Header
	BootItems []BootItem
	Bootable  bool
}

func (i *Image) isBootable() bool {
	if len(i.BootItems) == 0 {
		return false
	}
	return i.BootItems[0].Header.Type.IsKernel()
}

func (i *Image) readContainerHeader(f io.ReadSeeker) error {
	header := &i.Header
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := binary.Read(f, binary.LittleEndian, header); err != nil {
		return err
	}

	if header.Type != ZBITypeContainer {
		return fmt.Errorf("invalid header type, expected %#08x, got %#08x", ZBITypeContainer, header.Type)
	}

	if header.Magic != ItemMagic {
		return fmt.Errorf("invalid item magic, expected %#08x, got %#08x", ItemMagic, header.Magic)
	}

	if header.Extra != ContainerMagic {
		return fmt.Errorf("invalid container magic, expected %#08x, got %#08x", ContainerMagic, header.Extra)
	}

	return nil
}

type ZBIKernel struct {
	Entry             uint64
	ReserveMemorySize uint64
}

type ZirconKernel struct {
	HdrFile    Header
	HdrKernel  Header
	DataKernel ZBIKernel
	contents   []uint8
}

func (i *Image) readBootItems(f io.ReadSeeker) error {
	for {
		item := BootItem{}
		if err := readHeader(f, &item.Header); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		position, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		item.PayloadAddress = uint64(position)
		i.BootItems = append(i.BootItems, item)

		padding := align.Up(uint(item.Header.Length), 8)
		f.Seek(int64(padding), io.SeekCurrent)
	}
}

func Read(f io.ReadSeeker) (*Image, error) {
	image := &Image{}
	if err := image.readContainerHeader(f); err != nil {
		return nil, err
	}
	if err := image.readBootItems(f); err != nil {
		return nil, err
	}
	image.Bootable = image.isBootable()
	return image, nil
}

func Load(imagePath string) (*Image, error) {
	imageFile, err := os.Open(imagePath)
	defer imageFile.Close()

	if err != nil {
		return nil, fmt.Errorf("load ZBI image failed: %w", err)
	}

	image, err := Read(imageFile)
	if err != nil {
		return nil, fmt.Errorf("reading ZBI image failed: %w", err)
	}
	return image, nil
}

func readHeader(f io.ReadSeeker, h *Header) error {
	if err := binary.Read(f, binary.LittleEndian, h); err != nil {
		return err
	}
	return nil
}
