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
	// ContainerMagic is Zircon image header extra magic.
	ContainerMagic uint32 = 0x868cf7e6
	// ItemMagic is Ziron image format magic.
	ItemMagic uint32 = 0xb5781729
	// VersionFlag is a default version flag that can be used when bootstrapping a Zircon image header.
	VersionFlag uint32 = 0x00010000
	// CRC32Flag is a flag to indicate performing CRC32 check.
	CRC32Flag uint32 = 0x00020000
	// NoCRC32Flag is a flag to indicate not performing CRC32 check.
	NoCRC32Flag uint32 = 0x4a87e8d6
	// ZBITypeKernelPrefix is kernel prefix.
	ZBITypeKernelPrefix uint32 = 0x004e524b // KRN\0
	// ZBITypeKernelMask is mask to extract kernel prefix bits.
	ZBITypeKernelMask uint32 = 0x00FFFFFF // Mask to compare to the prefix.
)

// ZBITypeMetadata describes a ZBI type.
type ZBITypeMetadata struct {
	Name      string
	Extention string
}

// ZBIType is a uint32.
type ZBIType uint32

const (
	// ZBITypeContainer represents BOOT type.
	ZBITypeContainer ZBIType = 0x544f4f42
	// ZBITypeKernelX64 represents KRNL type, a x86 kernel.
	ZBITypeKernelX64 ZBIType = 0x4c4e524b // KRNL
	// ZBITypeKernelArm64 represents KRN8 type, an Arm64 kernel.
	ZBITypeKernelArm64 ZBIType = 0x384e524b // KRN8
	// ZBITypeDiscard represents SKIP type.
	ZBITypeDiscard ZBIType = 0x50494b53
	// ZBITypeStorageRamdisk represents RDSK type.
	ZBITypeStorageRamdisk ZBIType = 0x4b534452
	// ZBITypeStorageBootfs represents BFSB type.
	ZBITypeStorageBootfs ZBIType = 0x42534642
	// ZBITypeStorageKernel represents KSTR type.
	ZBITypeStorageKernel ZBIType = 0x5254534b
	// ZBITypeStorageBootfsFactory represents BFSF ty
	ZBITypeStorageBootfsFactory ZBIType = 0x46534642
	// ZBITypeCmdline represents CMDL type.
	ZBITypeCmdline ZBIType = 0x4c444d43
	// ZBITypeCrashlog represents BOOM type.
	ZBITypeCrashlog ZBIType = 0x4d4f4f42
	// ZBITypeNvram represents NVLL type.
	ZBITypeNvram ZBIType = 0x4c4c564e
	// ZBITypePlatformID represents PLID type.
	ZBITypePlatformID ZBIType = 0x44494C50
	// ZBITypeDrvBoardInfo represents mBSI type.
	ZBITypeDrvBoardInfo ZBIType = 0x4953426D
	// ZBITypeCPUConfig represents CPUC type.
	ZBITypeCPUConfig ZBIType = 0x43555043
	// ZBITypeCPUTopology represents TOPO type.
	ZBITypeCPUTopology ZBIType = 0x544F504F
	// ZBITypeMemConfig represents MEMC type.
	ZBITypeMemConfig ZBIType = 0x434D454D
	// ZBITypeKernelDriver represents KDRV type.
	ZBITypeKernelDriver ZBIType = 0x5652444B
	// ZBITypeAcpiRsdp represents RSDP type.
	ZBITypeAcpiRsdp ZBIType = 0x50445352
	// ZBITypeSMBios represents SMBI type.
	ZBITypeSMBios ZBIType = 0x49424d53
	// ZBITypeEFISystemTable represents EFIS type.
	ZBITypeEFISystemTable ZBIType = 0x53494645
	// ZBITypeFramebuffer represents SWFB type.
	ZBITypeFramebuffer ZBIType = 0x42465753
	// ZBITypeImageArgs represents IARG type.
	ZBITypeImageArgs ZBIType = 0x47524149
	// ZBITypeBootVersion represents BVRS type.
	ZBITypeBootVersion ZBIType = 0x53525642
	// ZBITypeDrvMacAddress represents mMAC type.
	ZBITypeDrvMacAddress ZBIType = 0x43414D6D
	// ZBITypeDrvPartitionMap represents mPRT type.
	ZBITypeDrvPartitionMap ZBIType = 0x5452506D
	// ZBITypeDrvBoardPrivate represents mBOR type.
	ZBITypeDrvBoardPrivate ZBIType = 0x524F426D
	// ZBITypeHwRebootReason represents HWRB type.
	ZBITypeHwRebootReason ZBIType = 0x42525748
	// ZBITypeSerialNumber represents SRLN type.
	ZBITypeSerialNumber ZBIType = 0x4e4c5253
	// ZBITypeBootloaderFile represents BTFL type.
	ZBITypeBootloaderFile ZBIType = 0x4C465442
	// ZBITypeDevicetree represents device tree type.
	ZBITypeDevicetree ZBIType = 0xd00dfeed
	// ZBITypeSecureEntropy represents RAND type.
	ZBITypeSecureEntropy ZBIType = 0x444e4152
)

// ZBITypes is a ZBIType to ZBITypeMetadata mapping.
var ZBITypes = map[ZBIType]ZBITypeMetadata{
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

// IsKernel tells if current ZBIType is kernel.
func (it *ZBIType) IsKernel() bool {
	return uint32(*it)&ZBITypeKernelMask == ZBITypeKernelPrefix
}

// IsDriverMetadata tells if current ZBIType contains driver meta data.
func (it *ZBIType) IsDriverMetadata() bool {
	return *it&0xFF == 0x6D // 'm'
}

// ToString return string representation of current ZBIType.
func (it *ZBIType) ToString() (string, error) {
	if typeMetadata, ok := ZBITypes[*it]; ok {
		return typeMetadata.Name, nil
	}
	return "", fmt.Errorf("can't find metadata for %#08x ZBIType", it)
}

// MarshalJSON returns JSON bytes of current ZBIType.
func (it *ZBIType) MarshalJSON() ([]byte, error) {
	name, err := it.ToString()
	if err != nil {
		return nil, err
	}
	nameWithQuotes := fmt.Sprintf("%q", name)
	return []byte(nameWithQuotes), nil
}

// Header abstracts a Zircon image header.
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

// BootItem abstracts a bootable item.
type BootItem struct {
	Header         Header
	PayloadAddress uint64
}

// NewContainerHeader returns a new image header with given image length.
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

// Image abstracts a Zircon image.
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

// ZBIKernel abstracts an in-memory kernel entry.
type ZBIKernel struct {
	Entry             uint64
	ReserveMemorySize uint64
}

// ZirconKernel is the whole contiguous image loaded into memory by the boot loader.
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

// Read parses a Ziron Image from an io.ReadSeeker.
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

// Load loads an Image from given path.
func Load(imagePath string) (*Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("load ZBI image failed: %w", err)
	}
	defer imageFile.Close()

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
