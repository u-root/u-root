// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rekby/gpt"
)

var (
	// LinuxMountsPath is the standard mountpoint list path
	LinuxMountsPath = "/proc/mounts"
)

// BlockDev maps a device name to a BlockStat structure for a given block device
type BlockDev struct {
	Name   string
	Stat   BlockStat
	FsUUID string
}

// Summary prints a multiline summary of the BlockDev object
// https://www.kernel.org/doc/Documentation/block/stat.txt
func (b BlockDev) Summary() string {
	return fmt.Sprintf(`BlockStat{
    Name: %v
    ReadIOs: %v
    ReadMerges: %v
    ReadSectors: %v
    ReadTicks: %v
    WriteIOs: %v
    WriteMerges: %v
    WriteSectors: %v
    WriteTicks: %v
    InFlight: %v
    IOTicks: %v
    TimeInQueue: %v
}`,
		b.Name,
		b.Stat.ReadIOs,
		b.Stat.ReadMerges,
		b.Stat.ReadSectors,
		b.Stat.ReadTicks,
		b.Stat.WriteIOs,
		b.Stat.WriteMerges,
		b.Stat.WriteSectors,
		b.Stat.WriteTicks,
		b.Stat.InFlight,
		b.Stat.IOTicks,
		b.Stat.TimeInQueue,
	)
}

// BlockStat provides block device information as contained in
// /sys/class/block/<device_name>/stat
type BlockStat struct {
	ReadIOs      uint64
	ReadMerges   uint64
	ReadSectors  uint64
	ReadTicks    uint64
	WriteIOs     uint64
	WriteMerges  uint64
	WriteSectors uint64
	WriteTicks   uint64
	InFlight     uint64
	IOTicks      uint64
	TimeInQueue  uint64
	// Kernel 4.18 added four fields for discard tracking, see
	// https://github.com/torvalds/linux/commit/bdca3c87fb7ad1cc61d231d37eb0d8f90d001e0c
	DiscardIOs     uint64
	DiscardMerges  uint64
	DiscardSectors uint64
	DiscardTicks   uint64
}

// SystemPartitionGUID is the GUID of EFI system partitions
// EFI System partitions have GUID C12A7328-F81F-11D2-BA4B-00A0C93EC93B
var SystemPartitionGUID = gpt.Guid([...]byte{
	0x28, 0x73, 0x2a, 0xc1,
	0x1f, 0xf8,
	0xd2, 0x11,
	0xba, 0x4b,
	0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b,
})

// BlockStatFromBytes parses a block stat file and returns a BlockStat object.
// The format of the block stat file is the one defined by Linux for
// /sys/class/block/<device_name>/stat
func BlockStatFromBytes(buf []byte) (*BlockStat, error) {
	fields := strings.Fields(string(buf))
	// BlockStat has 11 fields
	if len(fields) < 11 {
		return nil, fmt.Errorf("BlockStatFromBytes: parsing %q: got %d fields(%q), want at least 11", buf, len(fields), fields)
	}
	intfields := make([]uint64, 0)
	for _, field := range fields {
		v, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return nil, err
		}
		intfields = append(intfields, v)
	}
	bs := BlockStat{
		ReadIOs:      intfields[0],
		ReadMerges:   intfields[1],
		ReadSectors:  intfields[2],
		ReadTicks:    intfields[3],
		WriteIOs:     intfields[4],
		WriteMerges:  intfields[5],
		WriteSectors: intfields[6],
		WriteTicks:   intfields[7],
		InFlight:     intfields[8],
		IOTicks:      intfields[9],
		TimeInQueue:  intfields[10],
	}
	if len(fields) >= 15 {
		bs.DiscardIOs = intfields[11]
		bs.DiscardMerges = intfields[12]
		bs.DiscardSectors = intfields[13]
		bs.DiscardTicks = intfields[14]
	}
	return &bs, nil
}

// GetBlockStats iterates over /sys/class/block entries and returns a list of
// BlockDev objects, or an error if any
func GetBlockStats() ([]BlockDev, error) {
	blockdevs := make([]BlockDev, 0)
	devnames := make([]string, 0)
	root := "/sys/class/block"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		devnames = append(devnames, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, devname := range devnames {
		fd, err := os.Open(fmt.Sprintf("%s/%s/stat", root, devname))
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		buf, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		bstat, err := BlockStatFromBytes(buf)
		if err != nil {
			return nil, err
		}
		devpath := path.Join("/dev/", devname)
		uuid := getUUID(devpath)
		blockdevs = append(blockdevs, BlockDev{Name: devname, Stat: *bstat, FsUUID: uuid})
	}
	return blockdevs, nil
}

func getUUID(devpath string) (fsuuid string) {

	fsuuid = tryVFAT(devpath)
	if fsuuid != "" {
		log.Printf("###### FsUUIS in %s: %s", devpath, fsuuid)
		return fsuuid
	}
	fsuuid = tryEXT4(devpath)
	if fsuuid != "" {
		log.Printf("###### FsUUIS in %s: %s", devpath, fsuuid)
		return fsuuid
	}
	log.Printf("###### FsUUIS in %s: NONE", devpath)
	return ""
}

//see https://www.nongnu.org/ext2-doc/ext2.html#DISK-ORGANISATION
const (
	EXT2SprblkOff       = 1024 // Offset of superblock in partition
	EXT2SprblkSize      = 512  // Actually 1024 but most of the last byters are reserved
	EXT2SprblkMagicOff  = 56   // Offset of magic number in suberblock
	EXT2SprblkMagicSize = 2
	EXT2SprblkMagic     = '\uEF53' // fixed value
	EXT2SprblkUUIDOff   = 104      // Offset of UUID in superblock
	EXT2SprblkUUIDSize  = 16
)

func tryEXT4(devname string) (uuid string) {
	log.Printf("try ext4")
	var off int64

	file, err := os.Open(devname)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Printf("%s %d\n", fileinfo.Name(), fileinfo.Size())

	// magic number
	b := make([]byte, EXT2SprblkMagicSize)
	off = EXT2SprblkOff + EXT2SprblkMagicOff
	_, err = file.ReadAt(b, off)
	if err != nil {
		log.Println(err)
		return ""
	}
	magic := uint16(b[1])<<8 + uint16(b[0])
	fmt.Printf("magic: 0x%x\n", magic)
	if magic != EXT2SprblkMagic {
		log.Printf("try ext4")
		return ""
	}

	// filesystem UUID
	b = make([]byte, EXT2SprblkUUIDSize)
	off = EXT2SprblkOff + EXT2SprblkUUIDOff
	_, err = file.ReadAt(b, off)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	uuid = fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8],
		b[9], b[10], b[11], b[12], b[13], b[14], b[15])
	fmt.Printf("UUID=%s\n", uuid)

	return uuid
}

// see https://de.wikipedia.org/wiki/File_Allocation_Table#Aufbau
const (
	FAT32MagicOff  = 82 // Offset of magic number
	FAT32MagicSize = 8
	FAT32Magic     = "FAT32   " // fixed value
	FAT32IDOff     = 67         // Offset of filesystem-ID / serielnumber. Treated as short filesystem UUID
	FAT32IDSize    = 4
)

func tryVFAT(devname string) (uuid string) {
	log.Printf("try vfat")
	var off int64

	file, err := os.Open(devname)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Printf("%s %d\n", fileinfo.Name(), fileinfo.Size())

	// magic number
	b := make([]byte, FAT32MagicSize)
	off = 0 + FAT32MagicOff
	_, err = file.ReadAt(b, off)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	magic := string(b)
	fmt.Printf("magic: %s\n", magic)
	if magic != FAT32Magic {
		log.Printf("no vfat")
		return ""
	}

	// filesystem UUID
	b = make([]byte, FAT32IDSize)
	off = 0 + FAT32IDOff
	_, err = file.ReadAt(b, off)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	uuid = fmt.Sprintf("%02x%02x-%02x%02x",
		b[3], b[2], b[1], b[0])
	fmt.Printf("UUID=%s\n", uuid)

	return uuid
}

// GetGPTTable tries to read a GPT table from the block device described by the
// passed BlockDev object, and returns a gpt.Table object, or an error if any
func GetGPTTable(device BlockDev) (*gpt.Table, error) {
	fd, err := os.Open(fmt.Sprintf("/dev/%s", device.Name))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err = fd.Seek(512, os.SEEK_SET); err != nil {
		return nil, err
	}
	table, err := gpt.ReadTable(fd, 512)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

// FilterEFISystemPartitions returns a list of BlockDev objects whose underlying
// block device is a valid EFI system partition, or an error if any
func FilterEFISystemPartitions(devices []BlockDev) ([]BlockDev, error) {
	return PartitionsByGUID(devices, SystemPartitionGUID.String())
}

// PartitionsByGUID returns a list of BlockDev objects whose underlying
// block device has the given GUID
func PartitionsByGUID(devices []BlockDev, guid string) ([]BlockDev, error) {
	partitions := make([]BlockDev, 0)
	for _, device := range devices {
		table, err := GetGPTTable(device)
		if err != nil {
			log.Printf("Skipping %s: %v", device.Name, err)
			continue
		}
		for _, part := range table.Partitions {
			if part.IsEmpty() {
				continue
			}
			if part.Type.String() == guid {
				partitions = append(partitions, device)
			}
		}
	}
	return partitions, nil
}

// PartitionsByFsUUID returns a list of BlockDev objects whose underlying
// block device has a filesystem with the given UUID
func PartitionsByFsUUID(devices []BlockDev, fsuuid string) []BlockDev {
	partitions := make([]BlockDev, 0)
	for _, device := range devices {
		if device.FsUUID == fsuuid {
			partitions = append(partitions, device)
		}
	}
	return partitions
}

// GetMountpointByDevice gets the mountpoint by given
// device name. Returns on first match
func GetMountpointByDevice(devicePath string) (*string, error) {
	file, err := os.Open(LinuxMountsPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		deviceInfo := strings.Fields(scanner.Text())
		if deviceInfo[0] == devicePath {
			return &deviceInfo[1], nil
		}
	}

	return nil, errors.New("Mountpoint not found")
}
