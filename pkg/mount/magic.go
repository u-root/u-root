// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package mount

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

const blocksize = 65536

// These are inferred magic numbers from documents and partitions.
// Ones known to work are first, followed by a gap, followed by not
// tested ones. Please preserve this pattern.
var (
	EXT2     = []byte{0x53, 0xef}
	EXT3     = []byte{0x53, 0xef}
	EXT4     = []byte{0x53, 0xef}
	ISOFS    = []byte{1, 'C', 'D', '0', '0', '1'}
	SQUASHFS = []byte{'h', 's', 'q', 's'}
	XFS      = []byte{'X', 'F', 'S', 'B'}
	// There's no fixed magic number for the different FAT varieties
	// Usually they start with 0xEB but it's not mandatory.
	// Therefore we just list a few examples that we have seen in the wild.
	MSDOS = []byte{0xeb, 0x3c}
	VFAT  = []byte{0xeb, 0x58}
	// QEMU virtual VFAT
	VVFAT = []byte{0xeb, 0x3e}

	AAFS        = []byte{0x5a, 0x3c, 0x69, 0xf0}
	ADFS        = []byte{0xad, 0xf5}
	AFFS        = []byte{0xad, 0xff}
	AFS         = []byte{0x53, 0x46, 0x41, 0x4F}
	BDEVFS      = []byte{0x62, 0x64, 0x65, 0x76}
	BINDERFS    = []byte{0x6c, 0x6f, 0x6f, 0x70}
	BINFMTFS    = []byte{0x42, 0x49, 0x4e, 0x4d}
	BPF         = []byte{0xca, 0xfe, 0x4a, 0x11}
	BTRFS       = []byte{'_', 'B', 'H', 'R', 'f', 'S', '_', 'M'}
	CGROUP      = []byte{0x27, 0xe0, 0xeb}
	CGROUP2     = []byte{0x63, 0x67, 0x72, 0x70}
	CODA        = []byte{0x73, 0x75, 0x72, 0x45}
	CRAMFS      = []byte{0x28, 0xcd, 0x3d, 0x45}
	CRAMFSOther = []byte{0x45, 0x3d, 0xcd, 0x28}
	DAXFS       = []byte{0x64, 0x64, 0x61, 0x78}
	DEBUGFS     = []byte{0x64, 0x62, 0x67, 0x20}
	DEVPTS      = []byte{0x1c, 0xd1}
	ECRYPTFS    = []byte{0xf1, 0x5f}
	EFIVARFS    = []byte{0xde, 0x5e, 0x81, 0xe4}
	EFS         = []byte{0x41, 0x4A, 0x53}
	// EXFAT seems to be a samsung file system.
	// EXFAT       = []byte{0x53, 0xef}
	F2FS      = []byte{0xF2, 0xF5, 0x20, 0x10}
	FUSE      = []byte{0x65, 0x73, 0x55, 0x46}
	FUTEXFS   = []byte{0xBA, 0xD1, 0xDE, 0xA}
	HOSTFS    = []byte{0x00, 0xc0, 0xff, 0xee}
	HPFS      = []byte{0xf9, 0x95, 0xe8, 0x49}
	HUGETLBFS = []byte{0x95, 0x84, 0x58, 0xf6}
	JFFS2     = []byte{0x72, 0xb6}
	JFS       = []byte{0x31, 0x53, 0x46, 0x4a}
	MTD       = []byte{0x11, 0x30, 0x78, 0x54}
	NFS       = []byte{0x69, 0x69}
	NILFS     = []byte{0x34, 0x34}
	NSFS      = []byte{0x6e, 0x73, 0x66, 0x73}
	// From docs, not tested.
	NTFS       = []byte{0xeb, 0x52, 0x90, 'N', 'T', 'F', 'S', ' ', ' ', ' ', ' '}
	OCFS2      = []byte{0x74, 0x61, 0x63, 0x6f}
	OPENPROM   = []byte{0x9f, 0xa1}
	OVERLAYFS  = []byte{0x79, 0x4c, 0x76, 0x30}
	PIPEFS     = []byte{0x50, 0x49, 0x50, 0x45}
	PROC       = []byte{0x9f, 0xa0}
	PSTOREFS   = []byte{0x61, 0x65, 0x67, 0x6C}
	QNX4       = []byte{0x00, 0x2f}
	QNX6       = []byte{0x68, 0x19, 0x11, 0x22}
	RAMFS      = []byte{0x85, 0x84, 0x58, 0xf6}
	RDTGROUP   = []byte{0x76, 0x55, 0x82, 1}
	ROMFS      = []byte{0x72, 0x75}
	SECURITYFS = []byte{0x73, 0x63, 0x66, 0x73}
	SELINUX    = []byte{0xf9, 0x7c, 0xff, 0x8c}
	SMACK      = []byte{0x43, 0x41, 0x5d, 0x53}
	SMB        = []byte{0x51, 0x7B}
	SOCKFS     = []byte{0x53, 0x4F, 0x43, 0x4B}
	SYSFS      = []byte{0x62, 0x65, 0x65, 0x72}
	TMPFS      = []byte{0x01, 0x02, 0x19, 0x94}
	TRACEFS    = []byte{0x74, 0x72, 0x61, 0x63}
	UBIFS      = []byte{0x24, 0x05, 0x19, 0x05}
	UDF        = []byte{0x15, 0x01, 0x33, 0x46}
	USBDEVICE  = []byte{0x9f, 0xa2}
	V9FS       = []byte{0x01, 0x02, 0x19, 0x97}
	XENFS      = []byte{0xab, 0xba, 0x19, 0x74}
	ZONEFS     = []byte{0x5a, 0x4f, 0x46, 0x53}
	ZSMALLOC   = []byte{0x58, 0x29, 0x58, 0x29}
)

type magic struct {
	magic    []byte
	off      int64
	name     string
	flags    uintptr
	magicInt int64
}

// magics is just a list of magic structs.
// One file system in particular shares a single magic for several types.
// For that reason, and reasons of space, this is a list, not a map.
// Performance is not really an issue: it is a short list, and there are simply
// not enough block devices/file systems for it to really matter.
// The ordering for the identical magic number file systems matters: ext4 is more
// desirable than ext2, so, we want to find ext4 first.
// The order should NOT BE ALPHABETIC, therefore; it should be ordered with known systems
// first, and, to break ties, with the most desirable of those systems first.
var magics = []magic{
	// From the filesystems magic:
	// 0x438   leshort         0xEF53          Linux
	{magic: EXT4, name: "ext4", off: 0x438, magicInt: unix.EXT4_SUPER_MAGIC},
	{magic: EXT3, name: "ext3", off: 0x438, magicInt: unix.EXT3_SUPER_MAGIC},
	{magic: EXT2, name: "ext2", off: 0x438, magicInt: unix.EXT2_SUPER_MAGIC},
	// We will always mount vfat; it's backward compatible (we think?)
	{magic: MSDOS, name: "vfat", off: 0, magicInt: unix.MSDOS_SUPER_MAGIC},
	{magic: SQUASHFS, name: "squashfs", flags: MS_RDONLY, off: 0, magicInt: unix.SQUASHFS_MAGIC},
	{magic: ISOFS, name: "iso9660", flags: MS_RDONLY, off: 32768, magicInt: unix.ISOFS_SUPER_MAGIC},
	{magic: VFAT, name: "vfat", off: 0},
	{magic: VVFAT, name: "vfat", off: 0},
	{magic: XFS, name: "xfs", off: 0, magicInt: unix.XFS_SUPER_MAGIC},
	{magic: BTRFS, name: "btrfs", off: 0x10040, magicInt: unix.BTRFS_SUPER_MAGIC},
}

var unknownMagics = []magic{
	//
	// here there be dragons.
	//
	{magic: V9FS, name: "9p", off: -1, magicInt: unix.V9FS_MAGIC},
	{magic: ADFS, name: "adfs", off: -1, magicInt: unix.ADFS_SUPER_MAGIC},
	{magic: AFFS, name: "affs", off: -1, magicInt: unix.AFFS_SUPER_MAGIC},
	{magic: SMB, name: "cifs", off: -1, magicInt: unix.CIFS_SUPER_MAGIC},
	{magic: SMB, name: "smb3", off: -1},
	{magic: CODA, name: "coda", off: -1, magicInt: unix.CODA_SUPER_MAGIC},
	{magic: DEVPTS, name: "devpts", off: -1, magicInt: unix.DEVPTS_SUPER_MAGIC},
	{magic: ECRYPTFS, name: "ecryptfs", off: -1, magicInt: unix.ECRYPTFS_SUPER_MAGIC},
	{magic: EFIVARFS, name: "efivarfs", off: -1, magicInt: unix.EFIVARFS_MAGIC},
	{magic: EFS, name: "efs", off: -1, magicInt: unix.EFS_SUPER_MAGIC},
	{magic: F2FS, name: "f2fs", off: -1, magicInt: unix.F2FS_SUPER_MAGIC},
	{magic: FUSE, name: "fuse", off: -1, magicInt: unix.FUSE_SUPER_MAGIC},
	// ?? {magic: GFS2, name: "gfs2", off: -1},
	// who care ... {magic: HFSPLUS_VOLHEAD_SIG, name: "hfsplus", off: -1},
	{magic: HOSTFS, name: "hostfs", off: -1, magicInt: unix.HOSTFS_SUPER_MAGIC},
	{magic: HPFS, name: "hpfs", off: -1, magicInt: unix.HPFS_SUPER_MAGIC},
	{magic: HUGETLBFS, name: "hugetlbfs", off: -1, magicInt: unix.HUGETLBFS_MAGIC},
	{magic: JFFS2, name: "jffs2", off: -1, magicInt: unix.JFFS2_SUPER_MAGIC},
	{magic: JFS, name: "jfs", off: -1},
	{magic: NFS, name: "nfs", off: -1, magicInt: unix.NFS_SUPER_MAGIC},
	{magic: NTFS, name: "ntfs", off: -1},
	{magic: OPENPROM, name: "openpromfs", off: -1, magicInt: unix.OPENPROM_SUPER_MAGIC},
	{magic: OVERLAYFS, name: "overlay", off: -1, magicInt: unix.OVERLAYFS_SUPER_MAGIC},
	{magic: PIPEFS, name: "pipefs", off: -1, magicInt: unix.PIPEFS_MAGIC},
	{magic: PROC, name: "proc", flags: MS_RDONLY, off: -1, magicInt: unix.PROC_SUPER_MAGIC},
	{magic: PSTOREFS, name: "pstore", off: -1, magicInt: unix.PSTOREFS_MAGIC},
	{magic: QNX4, name: "qnx4", off: -1, magicInt: unix.QNX4_SUPER_MAGIC},
	{magic: QNX6, name: "qnx6", off: -1, magicInt: unix.QNX6_SUPER_MAGIC},
	{magic: RAMFS, name: "ramfs", off: -1, magicInt: unix.RAMFS_MAGIC},
	{magic: ROMFS, name: "romfs", flags: MS_RDONLY, off: -1},
	{magic: UBIFS, name: "ubifs", flags: MS_RDONLY, off: -1},
	{magic: UDF, name: "udf", off: -1, magicInt: unix.UDF_SUPER_MAGIC},
	{magic: ZONEFS, name: "zonefs", off: -1, magicInt: unix.ZONEFS_MAGIC},
}

// FindMagics finds all the magics matching a magic number.
func FindMagics(blk []byte) []magic {
	b := bytes.NewReader(blk)
	matches := []magic{}
	for _, v := range magics {
		mag := make([]byte, len(v.magic))
		if n, err := b.ReadAt(mag, v.off); err != nil || n < len(mag) {
			continue
		}
		if bytes.Equal(v.magic, mag) {
			matches = append(matches, v)
		}
	}
	return matches
}

// FSFromBlock determines the file system type of a block device.
// It returns a string and an error. The error can be for an IO operation,
// an unknown magic number, or a magic with an unsupported file system.
// There is still a question here about whether this ought to act like
// a map and return a bool, not an error, since there are so many bogus
// block devices and we don't care about most of them.
func FSFromBlock(n string) (fs string, flags uintptr, err error) {
	// Make sure we can open, read 64k, stat it, find the magic in magics,
	// and find the file system it names.
	f, err := os.Open(n)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	block := make([]byte, blocksize*2)
	if _, err := io.ReadAtLeast(f, block, len(block)); err != nil {
		return "", 0, fmt.Errorf("no suitable filesystem for %q: %w", n, err)
	}

	magics := FindMagics(block)
	if len(magics) == 0 {
		return "", 0, fmt.Errorf("no suitable filesystem for %q", n)
	}

	for _, m := range magics {
		if err := FindFileSystem(m.name); err == nil {
			return m.name, m.flags, nil
		}
	}
	return "", 0, fmt.Errorf("no suitable filesystem for %q, from magics %q", n, magics)
}

func FromStatFS(path string) (fs string, flags uintptr, err error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return "", 0, fmt.Errorf("statfs %s: %w", path, err)
	}

	for _, m := range magics {
		if m.magicInt == int64(s.Type) {
			return m.name, m.flags, nil
		}
	}

	for _, m := range unknownMagics {
		if m.magicInt == int64(s.Type) {
			return m.name, m.flags, nil
		}
	}
	return "", 0, fmt.Errorf("no suitable filesystem for %q, from magics %q", path, magics)
}

// IsTmpRamfs tells if the file path given is under a tmpfs or ramfs.
func IsTmpRamfs(path string) (bool, error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return false, err
	}
	// Force it to be int64 so that unix.RAMFS_MAGIC won't overflow on an
	// int32, which is the type for Type on some platforms.
	t := int64(s.Type)
	return t == unix.TMPFS_MAGIC || t == unix.RAMFS_MAGIC, nil
}
