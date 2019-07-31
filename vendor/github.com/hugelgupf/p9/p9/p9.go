// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package p9 is a 9P2000.L implementation.
//
// Servers implement Attacher and File interfaces.
//
// Clients can use Client.
package p9

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"syscall"
)

// Debug can be assigned to log.Printf to print messages received and sent.
var Debug = func(fmt string, v ...interface{}) {}

// OpenFlags is the mode passed to Open and Create operations.
//
// These correspond to bits sent over the wire.
type OpenFlags uint32

const (
	// ReadOnly is a Topen and Tcreate flag indicating read-only mode.
	ReadOnly OpenFlags = 0

	// WriteOnly is a Topen and Tcreate flag indicating write-only mode.
	WriteOnly OpenFlags = 1

	// ReadWrite is a Topen flag indicates read-write mode.
	ReadWrite OpenFlags = 2

	// OpenFlagsModeMask is a mask of valid OpenFlags mode bits.
	OpenFlagsModeMask OpenFlags = 3
)

// Mode returns only the open mode (read-only, read-write, or write-only).
func (o OpenFlags) Mode() OpenFlags {
	return o & OpenFlagsModeMask
}

// OSFlags converts a p9.OpenFlags to an int compatible with open(2).
func (o OpenFlags) OSFlags() int {
	return int(o & OpenFlagsModeMask)
}

// String implements fmt.Stringer.
func (o OpenFlags) String() string {
	switch o {
	case ReadOnly:
		return "ReadOnly"
	case WriteOnly:
		return "WriteOnly"
	case ReadWrite:
		return "ReadWrite"
	case OpenFlagsModeMask:
		return "OpenFlagsModeMask"
	default:
		return fmt.Sprintf("unknown (%#x)", uint32(o))
	}
}

// tag is a message tag.
type tag uint16

// fid is a file identifier.
type fid uint64

// FileMode are flags corresponding to file modes.
//
// These correspond to bits sent over the wire.
// These also correspond to mode_t bits.
type FileMode uint32

const (
	// FileModeMask is a mask of all the file mode bits of FileMode.
	FileModeMask FileMode = 0170000

	// ModeSocket is an (unused) mode bit for a socket.
	ModeSocket FileMode = 0140000

	// ModeSymlink is a mode bit for a symlink.
	ModeSymlink FileMode = 0120000

	// ModeRegular is a mode bit for regular files.
	ModeRegular FileMode = 0100000

	// ModeBlockDevice is a mode bit for block devices.
	ModeBlockDevice FileMode = 060000

	// ModeDirectory is a mode bit for directories.
	ModeDirectory FileMode = 040000

	// ModeCharacterDevice is a mode bit for a character device.
	ModeCharacterDevice FileMode = 020000

	// ModeNamedPipe is a mode bit for a named pipe.
	ModeNamedPipe FileMode = 010000

	// Read is a mode bit indicating read permission.
	Read FileMode = 04

	// Write is a mode bit indicating write permission.
	Write FileMode = 02

	// Exec is a mode bit indicating exec permission.
	Exec FileMode = 01

	// AllPermissions is a mask with rwx bits set for user, group and others.
	AllPermissions FileMode = 0777

	// Sticky is a mode bit indicating sticky directories.
	Sticky FileMode = 01000

	// permissionsMask is the mask to apply to FileModes for permissions. It
	// includes rwx bits for user, group and others, and sticky bit.
	permissionsMask FileMode = 01777
)

// QIDType is the most significant byte of the FileMode word, to be used as the
// Type field of p9.QID.
func (m FileMode) QIDType() QIDType {
	switch {
	case m.IsDir():
		return TypeDir
	case m.IsSocket(), m.IsNamedPipe(), m.IsCharacterDevice():
		// Best approximation.
		return TypeAppendOnly
	case m.IsSymlink():
		return TypeSymlink
	default:
		return TypeRegular
	}
}

// FileType returns the file mode without the permission bits.
func (m FileMode) FileType() FileMode {
	return m & FileModeMask
}

// Permissions returns just the permission bits of the mode.
func (m FileMode) Permissions() FileMode {
	return m & permissionsMask
}

// Writable returns the mode with write bits added.
func (m FileMode) Writable() FileMode {
	return m | 0222
}

// IsReadable returns true if m represents a file that can be read.
func (m FileMode) IsReadable() bool {
	return m&0444 != 0
}

// IsWritable returns true if m represents a file that can be written to.
func (m FileMode) IsWritable() bool {
	return m&0222 != 0
}

// IsExecutable returns true if m represents a file that can be executed.
func (m FileMode) IsExecutable() bool {
	return m&0111 != 0
}

// IsRegular returns true if m is a regular file.
func (m FileMode) IsRegular() bool {
	return m&FileModeMask == ModeRegular
}

// IsDir returns true if m represents a directory.
func (m FileMode) IsDir() bool {
	return m&FileModeMask == ModeDirectory
}

// IsNamedPipe returns true if m represents a named pipe.
func (m FileMode) IsNamedPipe() bool {
	return m&FileModeMask == ModeNamedPipe
}

// IsCharacterDevice returns true if m represents a character device.
func (m FileMode) IsCharacterDevice() bool {
	return m&FileModeMask == ModeCharacterDevice
}

// IsBlockDevice returns true if m represents a character device.
func (m FileMode) IsBlockDevice() bool {
	return m&FileModeMask == ModeBlockDevice
}

// IsSocket returns true if m represents a socket.
func (m FileMode) IsSocket() bool {
	return m&FileModeMask == ModeSocket
}

// IsSymlink returns true if m represents a symlink.
func (m FileMode) IsSymlink() bool {
	return m&FileModeMask == ModeSymlink
}

// ModeFromOS returns a FileMode from an os.FileMode.
func ModeFromOS(mode os.FileMode) FileMode {
	m := FileMode(mode.Perm())
	switch {
	case mode.IsDir():
		m |= ModeDirectory
	case mode&os.ModeSymlink != 0:
		m |= ModeSymlink
	case mode&os.ModeSocket != 0:
		m |= ModeSocket
	case mode&os.ModeNamedPipe != 0:
		m |= ModeNamedPipe
	case mode&os.ModeCharDevice != 0:
		m |= ModeCharacterDevice
	case mode&os.ModeDevice != 0:
		m |= ModeBlockDevice
	default:
		m |= ModeRegular
	}
	return m
}

// OSMode converts a p9.FileMode to an os.FileMode.
func (m FileMode) OSMode() os.FileMode {
	var osMode os.FileMode
	osMode |= os.FileMode(m.Permissions())
	switch {
	case m.IsDir():
		osMode |= os.ModeDir
	case m.IsSymlink():
		osMode |= os.ModeSymlink
	case m.IsSocket():
		osMode |= os.ModeSocket
	case m.IsNamedPipe():
		osMode |= os.ModeNamedPipe
	case m.IsCharacterDevice():
		osMode |= os.ModeCharDevice | os.ModeDevice
	case m.IsBlockDevice():
		osMode |= os.ModeDevice
	}
	return osMode
}

// UID represents a user ID.
type UID uint32

// Ok returns true if uid is not NoUID.
func (uid UID) Ok() bool {
	return uid != NoUID
}

// GID represents a group ID.
type GID uint32

// Ok returns true if gid is not NoGID.
func (gid GID) Ok() bool {
	return gid != NoGID
}

const (
	// notag is a sentinel used to indicate no valid tag.
	noTag tag = math.MaxUint16

	// Nofid is a sentinel used to indicate no valid fid.
	noFID fid = math.MaxUint32

	// NoUID is a sentinel used to indicate no valid UID.
	NoUID UID = math.MaxUint32

	// NoGID is a sentinel used to indicate no valid GID.
	NoGID GID = math.MaxUint32
)

// msgType is a type identifier.
type msgType uint8

// msgType declarations.
const (
	msgTlerror      msgType = 6
	msgRlerror              = 7
	msgTstatfs              = 8
	msgRstatfs              = 9
	msgTlopen               = 12
	msgRlopen               = 13
	msgTlcreate             = 14
	msgRlcreate             = 15
	msgTsymlink             = 16
	msgRsymlink             = 17
	msgTmknod               = 18
	msgRmknod               = 19
	msgTrename              = 20
	msgRrename              = 21
	msgTreadlink            = 22
	msgRreadlink            = 23
	msgTgetattr             = 24
	msgRgetattr             = 25
	msgTsetattr             = 26
	msgRsetattr             = 27
	msgTxattrwalk           = 30
	msgRxattrwalk           = 31
	msgTxattrcreate         = 32
	msgRxattrcreate         = 33
	msgTreaddir             = 40
	msgRreaddir             = 41
	msgTfsync               = 50
	msgRfsync               = 51
	msgTlink                = 70
	msgRlink                = 71
	msgTmkdir               = 72
	msgRmkdir               = 73
	msgTrenameat            = 74
	msgRrenameat            = 75
	msgTunlinkat            = 76
	msgRunlinkat            = 77
	msgTversion             = 100
	msgRversion             = 101
	msgTauth                = 102
	msgRauth                = 103
	msgTattach              = 104
	msgRattach              = 105
	msgTflush               = 108
	msgRflush               = 109
	msgTwalk                = 110
	msgRwalk                = 111
	msgTread                = 116
	msgRread                = 117
	msgTwrite               = 118
	msgRwrite               = 119
	msgTclunk               = 120
	msgRclunk               = 121
	msgTremove              = 122
	msgRremove              = 123
	msgTflushf              = 124
	msgRflushf              = 125
	msgTwalkgetattr         = 126
	msgRwalkgetattr         = 127
	msgTucreate             = 128
	msgRucreate             = 129
	msgTumkdir              = 130
	msgRumkdir              = 131
	msgTumknod              = 132
	msgRumknod              = 133
	msgTusymlink            = 134
	msgRusymlink            = 135
	msgTlconnect            = 136
	msgRlconnect            = 137
)

// QIDType represents the file type for QIDs.
//
// QIDType corresponds to the high 8 bits of a Plan 9 file mode.
type QIDType uint8

const (
	// TypeDir represents a directory type.
	TypeDir QIDType = 0x80

	// TypeAppendOnly represents an append only file.
	TypeAppendOnly QIDType = 0x40

	// TypeExclusive represents an exclusive-use file.
	TypeExclusive QIDType = 0x20

	// TypeMount represents a mounted channel.
	TypeMount QIDType = 0x10

	// TypeAuth represents an authentication file.
	TypeAuth QIDType = 0x08

	// TypeTemporary represents a temporary file.
	TypeTemporary QIDType = 0x04

	// TypeSymlink represents a symlink.
	TypeSymlink QIDType = 0x02

	// TypeLink represents a hard link.
	TypeLink QIDType = 0x01

	// TypeRegular represents a regular file.
	TypeRegular QIDType = 0x00
)

// QID is a unique file identifier.
//
// This may be embedded in other requests and responses.
type QID struct {
	// Type is the highest order byte of the file mode.
	Type QIDType

	// Version is an arbitrary server version number.
	Version uint32

	// Path is a unique server identifier for this path (e.g. inode).
	Path uint64
}

// String implements fmt.Stringer.
func (q QID) String() string {
	return fmt.Sprintf("QID{Type: %d, Version: %d, Path: %d}", q.Type, q.Version, q.Path)
}

// decode implements encoder.decode.
func (q *QID) decode(b *buffer) {
	q.Type = b.ReadQIDType()
	q.Version = b.Read32()
	q.Path = b.Read64()
}

// encode implements encoder.encode.
func (q *QID) encode(b *buffer) {
	b.WriteQIDType(q.Type)
	b.Write32(q.Version)
	b.Write64(q.Path)
}

// QIDGenerator is a simple generator for QIDs that atomically increments Path
// values.
type QIDGenerator struct {
	// uids is an ever increasing value that can be atomically incremented
	// to provide unique Path values for QIDs.
	uids uint64
}

// Get returns a new 9P unique ID with a unique Path given a QID type.
//
// While the 9P spec allows Version to be incremented every time the file is
// modified, we currently do not use the Version member for anything.  Hence,
// it is set to 0.
func (q *QIDGenerator) Get(t QIDType) QID {
	return QID{
		Type:    t,
		Version: 0,
		Path:    atomic.AddUint64(&q.uids, 1),
	}
}

// FSStat is used by statfs.
type FSStat struct {
	// Type is the filesystem type.
	Type uint32

	// BlockSize is the blocksize.
	BlockSize uint32

	// Blocks is the number of blocks.
	Blocks uint64

	// BlocksFree is the number of free blocks.
	BlocksFree uint64

	// BlocksAvailable is the number of blocks *available*.
	BlocksAvailable uint64

	// Files is the number of files available.
	Files uint64

	// FilesFree is the number of free file nodes.
	FilesFree uint64

	// FSID is the filesystem ID.
	FSID uint64

	// NameLength is the maximum name length.
	NameLength uint32
}

// decode implements encoder.decode.
func (f *FSStat) decode(b *buffer) {
	f.Type = b.Read32()
	f.BlockSize = b.Read32()
	f.Blocks = b.Read64()
	f.BlocksFree = b.Read64()
	f.BlocksAvailable = b.Read64()
	f.Files = b.Read64()
	f.FilesFree = b.Read64()
	f.FSID = b.Read64()
	f.NameLength = b.Read32()
}

// encode implements encoder.encode.
func (f *FSStat) encode(b *buffer) {
	b.Write32(f.Type)
	b.Write32(f.BlockSize)
	b.Write64(f.Blocks)
	b.Write64(f.BlocksFree)
	b.Write64(f.BlocksAvailable)
	b.Write64(f.Files)
	b.Write64(f.FilesFree)
	b.Write64(f.FSID)
	b.Write32(f.NameLength)
}

// AttrMask is a mask of attributes for getattr.
type AttrMask struct {
	Mode        bool
	NLink       bool
	UID         bool
	GID         bool
	RDev        bool
	ATime       bool
	MTime       bool
	CTime       bool
	INo         bool
	Size        bool
	Blocks      bool
	BTime       bool
	Gen         bool
	DataVersion bool
}

// Contains returns true if a contains all of the attributes masked as b.
func (a AttrMask) Contains(b AttrMask) bool {
	if b.Mode && !a.Mode {
		return false
	}
	if b.NLink && !a.NLink {
		return false
	}
	if b.UID && !a.UID {
		return false
	}
	if b.GID && !a.GID {
		return false
	}
	if b.RDev && !a.RDev {
		return false
	}
	if b.ATime && !a.ATime {
		return false
	}
	if b.MTime && !a.MTime {
		return false
	}
	if b.CTime && !a.CTime {
		return false
	}
	if b.INo && !a.INo {
		return false
	}
	if b.Size && !a.Size {
		return false
	}
	if b.Blocks && !a.Blocks {
		return false
	}
	if b.BTime && !a.BTime {
		return false
	}
	if b.Gen && !a.Gen {
		return false
	}
	if b.DataVersion && !a.DataVersion {
		return false
	}
	return true
}

// Empty returns true if no fields are masked.
func (a AttrMask) Empty() bool {
	return !a.Mode && !a.NLink && !a.UID && !a.GID && !a.RDev && !a.ATime && !a.MTime && !a.CTime && !a.INo && !a.Size && !a.Blocks && !a.BTime && !a.Gen && !a.DataVersion
}

// AttrMaskAll is an AttrMask with all fields masked.
var AttrMaskAll = AttrMask{
	Mode:        true,
	NLink:       true,
	UID:         true,
	GID:         true,
	RDev:        true,
	ATime:       true,
	MTime:       true,
	CTime:       true,
	INo:         true,
	Size:        true,
	Blocks:      true,
	BTime:       true,
	Gen:         true,
	DataVersion: true,
}

// String implements fmt.Stringer.
func (a AttrMask) String() string {
	var masks []string
	if a.Mode {
		masks = append(masks, "Mode")
	}
	if a.NLink {
		masks = append(masks, "NLink")
	}
	if a.UID {
		masks = append(masks, "UID")
	}
	if a.GID {
		masks = append(masks, "GID")
	}
	if a.RDev {
		masks = append(masks, "RDev")
	}
	if a.ATime {
		masks = append(masks, "ATime")
	}
	if a.MTime {
		masks = append(masks, "MTime")
	}
	if a.CTime {
		masks = append(masks, "CTime")
	}
	if a.INo {
		masks = append(masks, "INo")
	}
	if a.Size {
		masks = append(masks, "Size")
	}
	if a.Blocks {
		masks = append(masks, "Blocks")
	}
	if a.BTime {
		masks = append(masks, "BTime")
	}
	if a.Gen {
		masks = append(masks, "Gen")
	}
	if a.DataVersion {
		masks = append(masks, "DataVersion")
	}
	return fmt.Sprintf("AttrMask{with: %s}", strings.Join(masks, " "))
}

// decode implements encoder.decode.
func (a *AttrMask) decode(b *buffer) {
	mask := b.Read64()
	a.Mode = mask&0x00000001 != 0
	a.NLink = mask&0x00000002 != 0
	a.UID = mask&0x00000004 != 0
	a.GID = mask&0x00000008 != 0
	a.RDev = mask&0x00000010 != 0
	a.ATime = mask&0x00000020 != 0
	a.MTime = mask&0x00000040 != 0
	a.CTime = mask&0x00000080 != 0
	a.INo = mask&0x00000100 != 0
	a.Size = mask&0x00000200 != 0
	a.Blocks = mask&0x00000400 != 0
	a.BTime = mask&0x00000800 != 0
	a.Gen = mask&0x00001000 != 0
	a.DataVersion = mask&0x00002000 != 0
}

// encode implements encoder.encode.
func (a *AttrMask) encode(b *buffer) {
	var mask uint64
	if a.Mode {
		mask |= 0x00000001
	}
	if a.NLink {
		mask |= 0x00000002
	}
	if a.UID {
		mask |= 0x00000004
	}
	if a.GID {
		mask |= 0x00000008
	}
	if a.RDev {
		mask |= 0x00000010
	}
	if a.ATime {
		mask |= 0x00000020
	}
	if a.MTime {
		mask |= 0x00000040
	}
	if a.CTime {
		mask |= 0x00000080
	}
	if a.INo {
		mask |= 0x00000100
	}
	if a.Size {
		mask |= 0x00000200
	}
	if a.Blocks {
		mask |= 0x00000400
	}
	if a.BTime {
		mask |= 0x00000800
	}
	if a.Gen {
		mask |= 0x00001000
	}
	if a.DataVersion {
		mask |= 0x00002000
	}
	b.Write64(mask)
}

// Attr is a set of attributes for getattr.
type Attr struct {
	Mode             FileMode
	UID              UID
	GID              GID
	NLink            uint64
	RDev             uint64
	Size             uint64
	BlockSize        uint64
	Blocks           uint64
	ATimeSeconds     uint64
	ATimeNanoSeconds uint64
	MTimeSeconds     uint64
	MTimeNanoSeconds uint64
	CTimeSeconds     uint64
	CTimeNanoSeconds uint64
	BTimeSeconds     uint64
	BTimeNanoSeconds uint64
	Gen              uint64
	DataVersion      uint64
}

// String implements fmt.Stringer.
func (a Attr) String() string {
	return fmt.Sprintf("Attr{Mode: 0o%o, UID: %d, GID: %d, NLink: %d, RDev: %d, Size: %d, BlockSize: %d, Blocks: %d, ATime: {Sec: %d, NanoSec: %d}, MTime: {Sec: %d, NanoSec: %d}, CTime: {Sec: %d, NanoSec: %d}, BTime: {Sec: %d, NanoSec: %d}, Gen: %d, DataVersion: %d}",
		a.Mode, a.UID, a.GID, a.NLink, a.RDev, a.Size, a.BlockSize, a.Blocks, a.ATimeSeconds, a.ATimeNanoSeconds, a.MTimeSeconds, a.MTimeNanoSeconds, a.CTimeSeconds, a.CTimeNanoSeconds, a.BTimeSeconds, a.BTimeNanoSeconds, a.Gen, a.DataVersion)
}

// encode implements encoder.encode.
func (a *Attr) encode(b *buffer) {
	b.WriteFileMode(a.Mode)
	b.WriteUID(a.UID)
	b.WriteGID(a.GID)
	b.Write64(a.NLink)
	b.Write64(a.RDev)
	b.Write64(a.Size)
	b.Write64(a.BlockSize)
	b.Write64(a.Blocks)
	b.Write64(a.ATimeSeconds)
	b.Write64(a.ATimeNanoSeconds)
	b.Write64(a.MTimeSeconds)
	b.Write64(a.MTimeNanoSeconds)
	b.Write64(a.CTimeSeconds)
	b.Write64(a.CTimeNanoSeconds)
	b.Write64(a.BTimeSeconds)
	b.Write64(a.BTimeNanoSeconds)
	b.Write64(a.Gen)
	b.Write64(a.DataVersion)
}

// decode implements encoder.decode.
func (a *Attr) decode(b *buffer) {
	a.Mode = b.ReadFileMode()
	a.UID = b.ReadUID()
	a.GID = b.ReadGID()
	a.NLink = b.Read64()
	a.RDev = b.Read64()
	a.Size = b.Read64()
	a.BlockSize = b.Read64()
	a.Blocks = b.Read64()
	a.ATimeSeconds = b.Read64()
	a.ATimeNanoSeconds = b.Read64()
	a.MTimeSeconds = b.Read64()
	a.MTimeNanoSeconds = b.Read64()
	a.CTimeSeconds = b.Read64()
	a.CTimeNanoSeconds = b.Read64()
	a.BTimeSeconds = b.Read64()
	a.BTimeNanoSeconds = b.Read64()
	a.Gen = b.Read64()
	a.DataVersion = b.Read64()
}

// StatToAttr converts a Linux syscall stat structure to an Attr.
func StatToAttr(s *syscall.Stat_t, req AttrMask) (Attr, AttrMask) {
	attr := Attr{
		UID: NoUID,
		GID: NoGID,
	}
	if req.Mode {
		// p9.FileMode corresponds to Linux mode_t.
		attr.Mode = FileMode(s.Mode)
	}
	if req.NLink {
		attr.NLink = s.Nlink
	}
	if req.UID {
		attr.UID = UID(s.Uid)
	}
	if req.GID {
		attr.GID = GID(s.Gid)
	}
	if req.RDev {
		attr.RDev = s.Dev
	}
	if req.ATime {
		attr.ATimeSeconds = uint64(s.Atim.Sec)
		attr.ATimeNanoSeconds = uint64(s.Atim.Nsec)
	}
	if req.MTime {
		attr.MTimeSeconds = uint64(s.Mtim.Sec)
		attr.MTimeNanoSeconds = uint64(s.Mtim.Nsec)
	}
	if req.CTime {
		attr.CTimeSeconds = uint64(s.Ctim.Sec)
		attr.CTimeNanoSeconds = uint64(s.Ctim.Nsec)
	}
	if req.Size {
		attr.Size = uint64(s.Size)
	}
	if req.Blocks {
		attr.BlockSize = uint64(s.Blksize)
		attr.Blocks = uint64(s.Blocks)
	}

	// Use the req field because we already have it.
	req.BTime = false
	req.Gen = false
	req.DataVersion = false

	return attr, req
}

// SetAttrMask specifies a valid mask for setattr.
type SetAttrMask struct {
	Permissions        bool
	UID                bool
	GID                bool
	Size               bool
	ATime              bool
	MTime              bool
	CTime              bool
	ATimeNotSystemTime bool
	MTimeNotSystemTime bool
}

// IsSubsetOf returns whether s is a subset of m.
func (s SetAttrMask) IsSubsetOf(m SetAttrMask) bool {
	sb := s.bitmask()
	sm := m.bitmask()
	return sm|sb == sm
}

// String implements fmt.Stringer.
func (s SetAttrMask) String() string {
	var masks []string
	if s.Permissions {
		masks = append(masks, "Permissions")
	}
	if s.UID {
		masks = append(masks, "UID")
	}
	if s.GID {
		masks = append(masks, "GID")
	}
	if s.Size {
		masks = append(masks, "Size")
	}
	if s.ATime {
		masks = append(masks, "ATime")
	}
	if s.MTime {
		masks = append(masks, "MTime")
	}
	if s.CTime {
		masks = append(masks, "CTime")
	}
	if s.ATimeNotSystemTime {
		masks = append(masks, "ATimeNotSystemTime")
	}
	if s.MTimeNotSystemTime {
		masks = append(masks, "MTimeNotSystemTime")
	}
	return fmt.Sprintf("SetAttrMask{with: %s}", strings.Join(masks, " "))
}

// Empty returns true if no fields are masked.
func (s SetAttrMask) Empty() bool {
	return !s.Permissions && !s.UID && !s.GID && !s.Size && !s.ATime && !s.MTime && !s.CTime && !s.ATimeNotSystemTime && !s.MTimeNotSystemTime
}

// decode implements encoder.decode.
func (s *SetAttrMask) decode(b *buffer) {
	mask := b.Read32()
	s.Permissions = mask&0x00000001 != 0
	s.UID = mask&0x00000002 != 0
	s.GID = mask&0x00000004 != 0
	s.Size = mask&0x00000008 != 0
	s.ATime = mask&0x00000010 != 0
	s.MTime = mask&0x00000020 != 0
	s.CTime = mask&0x00000040 != 0
	s.ATimeNotSystemTime = mask&0x00000080 != 0
	s.MTimeNotSystemTime = mask&0x00000100 != 0
}

func (s SetAttrMask) bitmask() uint32 {
	var mask uint32
	if s.Permissions {
		mask |= 0x00000001
	}
	if s.UID {
		mask |= 0x00000002
	}
	if s.GID {
		mask |= 0x00000004
	}
	if s.Size {
		mask |= 0x00000008
	}
	if s.ATime {
		mask |= 0x00000010
	}
	if s.MTime {
		mask |= 0x00000020
	}
	if s.CTime {
		mask |= 0x00000040
	}
	if s.ATimeNotSystemTime {
		mask |= 0x00000080
	}
	if s.MTimeNotSystemTime {
		mask |= 0x00000100
	}
	return mask
}

// encode implements encoder.encode.
func (s *SetAttrMask) encode(b *buffer) {
	b.Write32(s.bitmask())
}

// SetAttr specifies a set of attributes for a setattr.
type SetAttr struct {
	Permissions      FileMode
	UID              UID
	GID              GID
	Size             uint64
	ATimeSeconds     uint64
	ATimeNanoSeconds uint64
	MTimeSeconds     uint64
	MTimeNanoSeconds uint64
}

// String implements fmt.Stringer.
func (s SetAttr) String() string {
	return fmt.Sprintf("SetAttr{Permissions: 0o%o, UID: %d, GID: %d, Size: %d, ATime: {Sec: %d, NanoSec: %d}, MTime: {Sec: %d, NanoSec: %d}}", s.Permissions, s.UID, s.GID, s.Size, s.ATimeSeconds, s.ATimeNanoSeconds, s.MTimeSeconds, s.MTimeNanoSeconds)
}

// decode implements encoder.decode.
func (s *SetAttr) decode(b *buffer) {
	s.Permissions = b.ReadPermissions()
	s.UID = b.ReadUID()
	s.GID = b.ReadGID()
	s.Size = b.Read64()
	s.ATimeSeconds = b.Read64()
	s.ATimeNanoSeconds = b.Read64()
	s.MTimeSeconds = b.Read64()
	s.MTimeNanoSeconds = b.Read64()
}

// encode implements encoder.encode.
func (s *SetAttr) encode(b *buffer) {
	b.WritePermissions(s.Permissions)
	b.WriteUID(s.UID)
	b.WriteGID(s.GID)
	b.Write64(s.Size)
	b.Write64(s.ATimeSeconds)
	b.Write64(s.ATimeNanoSeconds)
	b.Write64(s.MTimeSeconds)
	b.Write64(s.MTimeNanoSeconds)
}

// Apply applies this to the given Attr.
func (a *Attr) Apply(mask SetAttrMask, attr SetAttr) {
	if mask.Permissions {
		a.Mode = a.Mode&^permissionsMask | (attr.Permissions & permissionsMask)
	}
	if mask.UID {
		a.UID = attr.UID
	}
	if mask.GID {
		a.GID = attr.GID
	}
	if mask.Size {
		a.Size = attr.Size
	}
	if mask.ATime {
		a.ATimeSeconds = attr.ATimeSeconds
		a.ATimeNanoSeconds = attr.ATimeNanoSeconds
	}
	if mask.MTime {
		a.MTimeSeconds = attr.MTimeSeconds
		a.MTimeNanoSeconds = attr.MTimeNanoSeconds
	}
}

// Dirent is used for readdir.
type Dirent struct {
	// QID is the entry QID.
	QID QID

	// Offset is the offset in the directory.
	//
	// This will be communicated back the original caller.
	Offset uint64

	// Type is the 9P type.
	Type QIDType

	// Name is the name of the entry (i.e. basename).
	Name string
}

// String implements fmt.Stringer.
func (d Dirent) String() string {
	return fmt.Sprintf("Dirent{QID: %d, Offset: %d, Type: 0x%X, Name: %s}", d.QID, d.Offset, d.Type, d.Name)
}

// decode implements encoder.decode.
func (d *Dirent) decode(b *buffer) {
	d.QID.decode(b)
	d.Offset = b.Read64()
	d.Type = b.ReadQIDType()
	d.Name = b.ReadString()
}

// encode implements encoder.encode.
func (d *Dirent) encode(b *buffer) {
	d.QID.encode(b)
	b.Write64(d.Offset)
	b.WriteQIDType(d.Type)
	b.WriteString(d.Name)
}
