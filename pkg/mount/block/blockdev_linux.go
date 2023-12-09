// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package block finds, mounts, and modifies block devices on Linux systems.
package block

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/pci"
	"golang.org/x/sys/unix"
)

var (
	// LinuxMountsPath is the standard mountpoint list path
	LinuxMountsPath = "/proc/mounts"

	// Debug function to override for verbose logging.
	Debug = func(string, ...interface{}) {}

	// SystemPartitionGUID is the GUID of EFI system partitions
	// EFI System partitions have GUID C12A7328-F81F-11D2-BA4B-00A0C93EC93B
	SystemPartitionGUID = gpt.Guid([...]byte{
		0x28, 0x73, 0x2a, 0xc1,
		0x1f, 0xf8,
		0xd2, 0x11,
		0xba, 0x4b,
		0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b,
	})

	ErrListFormat = errors.New("device list needs to be of format vendor1:device1,vendor2:device2")
)

// BlockDev maps a device name to a BlockStat structure for a given block device
type BlockDev struct {
	Name   string
	FSType string
	FsUUID string
}

// Device makes sure the block device exists and returns a handle to it.
//
// maybeDevpath can be path like /dev/sda1, /sys/class/block/sda1 or just sda1.
// We will just use the last component.
func Device(maybeDevpath string) (*BlockDev, error) {
	devname := filepath.Base(maybeDevpath)
	if _, err := os.Stat(filepath.Join("/sys/class/block", devname)); err != nil {
		return nil, err
	}

	devpath := filepath.Join("/dev/", devname)
	if uuid, err := getFSUUID(devpath); err == nil {
		return &BlockDev{Name: devname, FsUUID: uuid}, nil
	}
	return &BlockDev{Name: devname}, nil
}

// String implements fmt.Stringer.
func (b *BlockDev) String() string {
	if len(b.FSType) > 0 {
		return fmt.Sprintf("BlockDevice(name=%s, fs_type=%s, fs_uuid=%s)", b.Name, b.FSType, b.FsUUID)
	}
	return fmt.Sprintf("BlockDevice(name=%s, fs_uuid=%s)", b.Name, b.FsUUID)
}

// DevicePath is the path to the actual device.
func (b BlockDev) DevicePath() string {
	return filepath.Join("/dev/", b.Name)
}

// Name implements mount.Mounter.
func (b *BlockDev) DevName() string {
	return b.Name
}

// Mount implements mount.Mounter.
func (b *BlockDev) Mount(path string, flags uintptr, opts ...func() error) (*mount.MountPoint, error) {
	devpath := filepath.Join("/dev", b.Name)
	if len(b.FSType) > 0 {
		return mount.Mount(devpath, path, b.FSType, "", flags, opts...)
	}

	return mount.TryMount(devpath, path, "", flags, opts...)
}

// GPTTable tries to read a GPT table from the block device described by the
// passed BlockDev object, and returns a gpt.Table object, or an error if any
func (b *BlockDev) GPTTable() (*gpt.Table, error) {
	fd, err := os.Open(b.DevicePath())
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	blkSize, err := b.BlockSize()
	if err != nil {
		blkSize = 512
	}

	if _, err := fd.Seek(int64(blkSize), io.SeekStart); err != nil {
		return nil, err
	}
	table, err := gpt.ReadTable(fd, uint64(blkSize))
	if err != nil {
		return nil, err
	}
	return &table, nil
}

// PhysicalBlockSize returns the physical block size.
func (b *BlockDev) PhysicalBlockSize() (int, error) {
	f, err := os.Open(b.DevicePath())
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return unix.IoctlGetInt(int(f.Fd()), unix.BLKPBSZGET)
}

// BlockSize returns the logical block size (BLKSSZGET).
func (b *BlockDev) BlockSize() (int, error) {
	f, err := os.Open(b.DevicePath())
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return unix.IoctlGetInt(int(f.Fd()), unix.BLKSSZGET)
}

// KernelBlockSize returns the soft block size used inside the kernel (BLKBSZGET).
func (b *BlockDev) KernelBlockSize() (int, error) {
	f, err := os.Open(b.DevicePath())
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return unix.IoctlGetInt(int(f.Fd()), unix.BLKBSZGET)
}

func ioctlGetUint64(fd int, req uint) (uint64, error) {
	var value uint64
	_, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&value)))
	if err != 0 {
		return 0, err
	}
	return value, nil
}

// Size returns the size in bytes.
func (b *BlockDev) Size() (uint64, error) {
	f, err := os.Open(b.DevicePath())
	if err != nil {
		return 0, err
	}
	defer f.Close()

	sz, err := ioctlGetUint64(int(f.Fd()), unix.BLKGETSIZE64)
	if err != nil {
		return 0, &os.PathError{
			Op:   "get size",
			Path: b.DevicePath(),
			Err:  os.NewSyscallError("ioctl(BLKGETSIZE64)", err),
		}
	}
	return sz, nil
}

// ReadPartitionTable prompts the kernel to re-read the partition table on this block device.
func (b *BlockDev) ReadPartitionTable() error {
	f, err := os.OpenFile(b.DevicePath(), os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return unix.IoctlSetInt(int(f.Fd()), unix.BLKRRPART, 0)
}

// PCIInfo searches sysfs for the PCI vendor and device id.
// We fill in the PCI struct with just those two elements.
func (b *BlockDev) PCIInfo() (*pci.PCI, error) {
	p, err := filepath.EvalSymlinks(filepath.Join("/sys/class/block", b.Name))
	if err != nil {
		return nil, err
	}
	// Loop through devices until we find the actual backing pci device.
	// For Example:
	// /sys/class/block/nvme0n1p1 usually resolves to something like
	// /sys/devices/pci..../.../.../nvme/nvme0/nvme0n1/nvme0n1p1. This leads us to the
	// first partition of the first namespace of the nvme0 device. In this case, the actual pci device and vendor
	// is found in nvme, three levels up. We traverse back up to the parent device
	// and we keep going until we find a device and vendor file.
	dp := filepath.Join(p, "device")
	vp := filepath.Join(p, "vendor")
	found := false
	for p != "/sys/devices" {
		// Check if there is a vendor and device file in this directory.
		if d, err := os.Stat(dp); err == nil && !d.IsDir() {
			if v, err := os.Stat(vp); err == nil && !v.IsDir() {
				found = true
				break
			}
		}
		p = filepath.Dir(p)
		dp = filepath.Join(p, "device")
		vp = filepath.Join(p, "vendor")
	}
	if !found {
		return nil, fmt.Errorf("unable to find backing pci device with device and vendor files for %v", b.Name)
	}

	return pci.OnePCI(p)
}

func getFSUUID(devpath string) (string, error) {
	file, err := os.Open(devpath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fsuuid, err := tryFAT32(file)
	if err == nil {
		return fsuuid, nil
	}
	fsuuid, err = tryFAT16(file)
	if err == nil {
		return fsuuid, nil
	}
	fsuuid, err = tryEXT4(file)
	if err == nil {
		return fsuuid, nil
	}
	fsuuid, err = tryXFS(file)
	if err == nil {
		return fsuuid, nil
	}
	return "", fmt.Errorf("unknown UUID (not vfat, ext4, nor xfs)")
}

// See https://www.nongnu.org/ext2-doc/ext2.html#DISK-ORGANISATION.
const (
	// Offset of superblock in partition.
	ext2SprblkOff = 1024

	// Offset of magic number in suberblock.
	ext2SprblkMagicOff  = 56
	ext2SprblkMagicSize = 2

	ext2SprblkMagic = 0xEF53

	// Offset of UUID in superblock.
	ext2SprblkUUIDOff  = 104
	ext2SprblkUUIDSize = 16
)

func tryEXT4(file io.ReaderAt) (string, error) {
	var off int64

	// Read magic number.
	b := make([]byte, ext2SprblkMagicSize)
	off = ext2SprblkOff + ext2SprblkMagicOff
	if _, err := file.ReadAt(b, off); err != nil {
		return "", err
	}
	magic := binary.LittleEndian.Uint16(b[:2])
	if magic != ext2SprblkMagic {
		return "", fmt.Errorf("ext4 magic not found")
	}

	// Filesystem UUID.
	b = make([]byte, ext2SprblkUUIDSize)
	off = ext2SprblkOff + ext2SprblkUUIDOff
	if _, err := file.ReadAt(b, off); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// See https://de.wikipedia.org/wiki/File_Allocation_Table#Aufbau.
const (
	fat12Magic = "FAT12   "
	fat16Magic = "FAT16   "

	// Offset of magic number.
	fat16MagicOff  = 0x36
	fat16MagicSize = 8

	// Offset of filesystem ID / serial number. Treated as short filesystem UUID.
	fat16IDOff  = 0x27
	fat16IDSize = 4
)

func tryFAT16(file io.ReaderAt) (string, error) {
	// Read magic number.
	b := make([]byte, fat16MagicSize)
	if _, err := file.ReadAt(b, fat16MagicOff); err != nil {
		return "", err
	}
	magic := string(b)
	if magic != fat16Magic && magic != fat12Magic {
		return "", fmt.Errorf("fat16 magic not found")
	}

	// Filesystem UUID.
	b = make([]byte, fat16IDSize)
	if _, err := file.ReadAt(b, fat16IDOff); err != nil {
		return "", err
	}

	return fmt.Sprintf("%02x%02x-%02x%02x", b[3], b[2], b[1], b[0]), nil
}

// See https://de.wikipedia.org/wiki/File_Allocation_Table#Aufbau.
const (
	fat32Magic = "FAT32   "

	// Offset of magic number.
	fat32MagicOff  = 0x52
	fat32MagicSize = 8

	// Offset of filesystem ID / serial number. Treated as short filesystem UUID.
	fat32IDOff  = 67
	fat32IDSize = 4
)

func tryFAT32(file io.ReaderAt) (string, error) {
	// Read magic number.
	b := make([]byte, fat32MagicSize)
	if _, err := file.ReadAt(b, fat32MagicOff); err != nil {
		return "", err
	}
	magic := string(b)
	if magic != fat32Magic {
		return "", fmt.Errorf("fat32 magic not found")
	}

	// Filesystem UUID.
	b = make([]byte, fat32IDSize)
	if _, err := file.ReadAt(b, fat32IDOff); err != nil {
		return "", err
	}

	return fmt.Sprintf("%02x%02x-%02x%02x", b[3], b[2], b[1], b[0]), nil
}

const (
	xfsMagic     = "XFSB"
	xfsMagicSize = 4
	xfsUUIDOff   = 32
	xfsUUIDSize  = 16
)

func tryXFS(file io.ReaderAt) (string, error) {
	// Read magic number.
	b := make([]byte, xfsMagicSize)
	if _, err := file.ReadAt(b, 0); err != nil {
		return "", err
	}
	magic := string(b)
	if magic != xfsMagic {
		return "", fmt.Errorf("xfs magic not found")
	}

	// Filesystem UUID.
	b = make([]byte, xfsUUIDSize)
	if _, err := file.ReadAt(b, xfsUUIDOff); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// BlockDevices is a list of block devices.
type BlockDevices []*BlockDev

// GetBlockDevices iterates over /sys/class/block entries and returns a list of
// BlockDev objects, or an error if any
func GetBlockDevices() (BlockDevices, error) {
	var blockdevs []*BlockDev
	var devnames []string

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
		dev, err := Device(rel)
		if err != nil {
			return err
		}
		blockdevs = append(blockdevs, dev)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blockdevs, nil
}

// FilterName returns a list of BlockDev objects whose underlying
// block device has a Name with the given Name
func (b BlockDevices) FilterName(name string) BlockDevices {
	partitions := make(BlockDevices, 0)
	for _, device := range b {
		if device.Name == name {
			partitions = append(partitions, device)
		}
	}
	return partitions
}

// FilterNames filters block devices by the given list of device names (e.g.
// /dev/sda1 sda2 /sys/class/block/sda3).
func (b BlockDevices) FilterNames(names ...string) BlockDevices {
	m := make(map[string]struct{})
	for _, n := range names {
		m[filepath.Base(n)] = struct{}{}
	}

	var devices BlockDevices
	for _, device := range b {
		if _, ok := m[device.Name]; ok {
			devices = append(devices, device)
		}
	}
	return devices
}

// FilterFSUUID returns a list of BlockDev objects whose underlying block
// device has a filesystem with the given FSUUID.
func (b BlockDevices) FilterFSUUID(fsuuid string) BlockDevices {
	partitions := make(BlockDevices, 0)
	for _, device := range b {
		if device.FsUUID == fsuuid {
			partitions = append(partitions, device)
		}
	}
	return partitions
}

// FilterZeroSize attempts to find block devices that have at least one block
// of content.
//
// This serves to eliminate block devices that have no backing storage, but
// appear in /sys/class/block anyway (like some loop, nbd, or ram devices).
func (b BlockDevices) FilterZeroSize() BlockDevices {
	var nb BlockDevices
	for _, device := range b {
		if n, err := device.Size(); err != nil || n == 0 {
			continue
		}
		nb = append(nb, device)
	}
	return nb
}

// FilterHavingPartitions returns BlockDevices with have the specified
// partitions. (e.g. f(1, 2) {sda, sda1, sda2, sdb} -> {sda})
func (b BlockDevices) FilterHavingPartitions(parts []int) BlockDevices {
	devices := make(BlockDevices, 0)
	for _, device := range b {
		hasParts := true
		for _, part := range parts {
			if _, err := os.Stat(filepath.Join("/sys/class/block",
				ComposePartName(device.Name, part))); err != nil {
				hasParts = false
				break
			}
		}
		if hasParts {
			devices = append(devices, device)
		}
	}
	return devices
}

// FilterPartID returns partitions with the given partition ID GUID.
func (b BlockDevices) FilterPartID(guid string) BlockDevices {
	var names []string
	for _, device := range b {
		table, err := device.GPTTable()
		if err != nil {
			continue
		}
		for i, part := range table.Partitions {
			if part.IsEmpty() {
				continue
			}
			if strings.EqualFold(part.Id.String(), guid) {
				names = append(names, ComposePartName(device.Name, i+1))
			}
		}
	}
	return b.FilterNames(names...)
}

// FilterPartType returns partitions with the given partition type GUID.
func (b BlockDevices) FilterPartType(guid string) BlockDevices {
	var names []string
	for _, device := range b {
		table, err := device.GPTTable()
		if err != nil {
			continue
		}
		for i, part := range table.Partitions {
			if part.IsEmpty() {
				continue
			}
			if strings.EqualFold(part.Type.String(), guid) {
				names = append(names, ComposePartName(device.Name, i+1))
			}
		}
	}
	return b.FilterNames(names...)
}

// FilterPartLabel returns a list of BlockDev objects whose underlying block
// device has the given partition label. The name comparison is case-insensitive.
func (b BlockDevices) FilterPartLabel(label string) BlockDevices {
	var names []string
	for _, device := range b {
		table, err := device.GPTTable()
		if err != nil {
			continue
		}
		for i, part := range table.Partitions {
			if part.IsEmpty() {
				continue
			}
			if strings.EqualFold(part.Name(), label) {
				names = append(names, ComposePartName(device.Name, i+1))
			}
		}
	}
	return b.FilterNames(names...)
}

// FilterAllowPCIString parses a string in the format vendor:device,vendor:device
// and returns a list of BlockDev objects whose backing pci devices match
// the vendor:device pairs passed in. All values are treated as hex.
// E.g. 0x8086:0xABCD,8086:0x1234
func (b BlockDevices) FilterAllowPCIString(allowlist string) (BlockDevices, error) {
	pciList, err := parsePCIList(allowlist)
	if err != nil {
		return nil, err
	}
	return b.FilterAllowPCI(pciList), nil
}

// FilterAllowPCI returns a list of BlockDev objects whose backing
// pci devices match the allowlist of PCI devices passed in.
// FilterAllowPCI discards entries which don't have a matching PCI vendor
// and device ID as an entry in the allowlist.
func (b BlockDevices) FilterAllowPCI(allowlist pci.Devices) BlockDevices {
	type mapKey struct {
		vendor, device uint16
	}
	m := make(map[mapKey]bool)

	for _, v := range allowlist {
		m[mapKey{v.Vendor, v.Device}] = true
	}
	Debug("allow map is %v", m)

	partitions := make(BlockDevices, 0)
	for _, device := range b {
		p, err := device.PCIInfo()
		if err != nil {
			// In the case of an error, we err on the safe side and choose to block it.
			// Not all block devices are backed by a pci device, for example SATA drives.
			Debug("Failed to find PCI info; %v", err)
			continue
		}
		if _, ok := m[mapKey{p.Vendor, p.Device}]; !ok {
			Debug("Blocking device %v since it doesn't appear in allowlist", device.Name)
			continue
		}
		// Included in allowlist, we're good to go
		Debug("Allowing device %v, with pci %v, in map", device, p)
		partitions = append(partitions, device)
	}
	return partitions
}

// FilterBlockPCIString parses a string in the format vendor:device,vendor:device
// and returns a list of BlockDev objects whose backing pci devices do not match
// the vendor:device pairs passed in. All values are treated as hex.
// E.g. 0x8086:0xABCD,8086:0x1234
func (b BlockDevices) FilterBlockPCIString(blocklist string) (BlockDevices, error) {
	pciList, err := parsePCIList(blocklist)
	if err != nil {
		return nil, err
	}
	return b.FilterBlockPCI(pciList), nil
}

// FilterBlockPCI returns a list of BlockDev objects whose backing
// pci devices do not match the blocklist of PCI devices passed in.
// FilterBlockPCI discards entries which have a matching PCI vendor
// and device ID as an entry in the blocklist.
func (b BlockDevices) FilterBlockPCI(blocklist pci.Devices) BlockDevices {
	type mapKey struct {
		vendor, device uint16
	}
	m := make(map[mapKey]bool)

	for _, v := range blocklist {
		m[mapKey{v.Vendor, v.Device}] = true
	}
	Debug("block map is %v", m)

	partitions := make(BlockDevices, 0)
	for _, device := range b {
		p, err := device.PCIInfo()
		if err != nil {
			// In the case of an error, we err on the safe side and choose not to block it.
			// Not all block devices are backed by a pci device, for example SATA drives.
			Debug("Failed to find PCI info; %v", err)
			partitions = append(partitions, device)
			continue
		}
		if _, ok := m[mapKey{p.Vendor, p.Device}]; !ok {
			// Not in blocklist, we're good to go
			Debug("Not blocking device %v, with pci %v, not in map", device, p)
			partitions = append(partitions, device)
		} else {
			log.Printf("Blocking device %v since it appears in blocklist", device.Name)
		}
	}
	return partitions
}

// parsePCIList parses a string in the format vendor:device,vendor:device
// and returns a list of PCI devices containing the vendor and device pairs.
func parsePCIList(parseList string) (pci.Devices, error) {
	pciList := pci.Devices{}
	bL := strings.Split(parseList, ",")
	for _, b := range bL {
		p := strings.Split(b, ":")
		if len(p) != 2 {
			return nil, fmt.Errorf("parsing device list %q: %w", parseList, ErrListFormat)
		}
		// Check that values are hex and convert them to sysfs formats
		// This accepts 0xABCD and turns it into 0xabcd
		// abcd also turns into 0xabcd
		v, err := strconv.ParseUint(strings.TrimPrefix(p[0], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("parsing pci device %q:%w", p[0], err)
		}

		d, err := strconv.ParseUint(strings.TrimPrefix(p[1], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("parsing pci device %q:%w", p[1], err)
		}

		pciList = append(pciList, &pci.PCI{Vendor: uint16(v), Device: uint16(d)})
	}
	return pciList, nil
}

// ComposePartName returns the partition name described by the parent devName
// and partNo counting from 1. It is assumed that device names ending in a
// number like nvme0n1 have partitions named like nvme0n1p1, nvme0n1p2, ...
// and devices ending in a letter like sda have partitions named like
//
//	sda1, sda2, ...
func ComposePartName(devName string, partNo int) string {
	r := []rune(devName[len(devName)-1:])
	if unicode.IsDigit(r[0]) {
		return fmt.Sprintf("%sp%d", devName, partNo)
	}
	return fmt.Sprintf("%s%d", devName, partNo)
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

	return nil, errors.New("mountpoint not found")
}
