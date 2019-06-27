// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

// Linux mode_t bits.
const (
	modeTypeMask    = 0170000
	modeSocket      = 0140000
	modeSymlink     = 0120000
	modeFile        = 0100000
	modeBlock       = 0060000
	modeDir         = 0040000
	modeChar        = 0020000
	modeFIFO        = 0010000
	modeSUID        = 0004000
	modeSGID        = 0002000
	modeSticky      = 0001000
	modePermissions = 0000777
)

type inumber uint64

type UnixFiler struct {
	Root       string
	CreateDevs bool
	// inodes is a record of inodes we have created. We do not protect it with a
	// mutex as cpio reading is inherently serial.
	inodes map[inumber]Record
}

var modeMap = map[uint64]os.FileMode{
	modeSocket:  os.ModeSocket,
	modeSymlink: os.ModeSymlink,
	modeFile:    0,
	modeBlock:   os.ModeDevice,
	modeDir:     os.ModeDir,
	modeChar:    os.ModeCharDevice,
	modeFIFO:    os.ModeNamedPipe,
}

// NewUnixFiler returns a new Filer for Unix-like systems.
// For any given build of this package,
// for a target operating sytems, we'll only support one flavor of filer.
// We thought about having multiple filer
// types, and being able to specify them in the command line, but given
// the existence of the VFS layer in the kernel, which exists to abstract
// all kinds of file systems to a single standard API, it was hard to
// make a case for recreating a "vfs" layer in this package.
// Realistically, for cpio, we're almost certainly never going to have
// anything but Unix-like kernels.
func NewUnixFiler(opts ...func(*UnixFiler)) Filer {
	var f = &UnixFiler{Root: ".", inodes: map[inumber]Record{}}
	for _, o := range opts {
		o(f)
	}
	return f
}

// setModes sets the modes, changing the easy ones first and the harder ones last.
// In this way, we set as much as we can before bailing out.
// N.B.: if you set something with S_ISUID, then change the owner,
// the kernel (Linux, OSX, etc.) clears S_ISUID (a good idea). So, the simple thing:
// Do the chmod operations in order of difficulty, and give up as soon as we fail.
// Set the basic permissions -- not including SUID, GUID, etc.
// Set the times
// Set the owner
// Set ALL the mode bits, in case we need to do SUID, etc. If we could not
// set the owner, we won't even try this operation of course, so we won't
// have SUID incorrectly set for the wrong user.
func setModes(r Record) error {
	if err := os.Chmod(r.Name, toFileMode(r)&os.ModePerm); err != nil {
		return err
	}
	/*if err := os.Chtimes(r.Name, time.Time{}, time.Unix(int64(r.MTime), 0)); err != nil {
		return err
	}*/
	if err := os.Chown(r.Name, int(r.UID), int(r.GID)); err != nil {
		return err
	}
	if err := os.Chmod(r.Name, toFileMode(r)); err != nil {
		return err
	}
	return nil
}

func toFileMode(r Record) os.FileMode {
	m := os.FileMode(perm(r))
	if r.Mode&unix.S_ISUID != 0 {
		m |= os.ModeSetuid
	}
	if r.Mode&unix.S_ISGID != 0 {
		m |= os.ModeSetgid
	}
	if r.Mode&unix.S_ISVTX != 0 {
		m |= os.ModeSticky
	}
	return m
}

func perm(r Record) uint32 {
	return uint32(r.Mode) & modePermissions
}

func dev(r Record) int {
	return int(r.Rmajor<<8 | r.Rminor)
}

func linuxModeToFileType(m uint64) (os.FileMode, error) {
	if t, ok := modeMap[m&modeTypeMask]; ok {
		return t, nil
	}
	return 0, fmt.Errorf("Invalid file type %#o", m&modeTypeMask)
}

// CreateFileInRoot creates a local file for f relative to rootDir.
//
// It will attempt to set all metadata for the file, including ownership,
// times, and permissions. If these fail, it only returns an error if
// forcePriv is true.
//
// Block and char device creation will only return error if forcePriv is true.
func (f *UnixFiler) Create(r Record) error {
	m, err := linuxModeToFileType(r.Mode)
	if err != nil {
		return err
	}

	r.Name = filepath.Clean(filepath.Join(f.Root, r.Name))
	dir := filepath.Dir(r.Name)
	// The problem: many cpio archives do not specify the directories and
	// hence the permissions. They just specify the whole path.  In order
	// to create files in these directories, we have to make them at least
	// mode 755.
	if _, err := os.Stat(dir); os.IsNotExist(err) && len(dir) > 0 {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("CreateFileInRoot %q: %v", r.Name, err)
		}
	}

	switch m {
	case os.ModeSocket, os.ModeNamedPipe:
		return fmt.Errorf("%q: type %v: cannot create IPC endpoints", r.Name, m)

	case os.ModeSymlink:
		content, err := ioutil.ReadAll(uio.Reader(r))
		if err != nil {
			return err
		}
		return os.Symlink(string(content), r.Name)

	case os.FileMode(0):
		// Check the dev/inumber. If we have seen it before, create a hardlink
		// and return. The data associated with a set of records
		// with the same inumber will always be attached to the first
		// record, so we need not worry whether the data was written.
		// If r.Ino is 0, the archive was produced on a
		// system that does not have hard links (e.g. Plan 9); a file system
		// that does not have hard links (e.g. vfat); or by a program/pkg that
		// does not create hard link records (pkg/uroot).
		// In that case, we need not check.
		if r.Ino != 0 {
			if of, ok := f.inodes[inumber(r.Ino)]; ok {
				if err := os.Link(of.Name, r.Name); err != nil {
					return err
				}
			}
		}
		nf, err := os.Create(r.Name)
		if err != nil {
			return err
		}
		f.inodes[inumber(r.Ino)] = r
		defer nf.Close()

		if _, err := io.Copy(nf, uio.Reader(r)); err != nil {
			return err
		}

	case os.ModeDir:
		if err := os.MkdirAll(r.Name, toFileMode(r)); err != nil {
			return err
		}

	case os.ModeDevice:
		if err := syscall.Mknod(r.Name, perm(r)|syscall.S_IFBLK, dev(r)); err != nil && f.CreateDevs {
			return err
		}

	case os.ModeCharDevice:
		if err := syscall.Mknod(r.Name, perm(r)|syscall.S_IFCHR, dev(r)); err != nil && f.CreateDevs {
			return err
		}

	default:
		return fmt.Errorf("%v: Unknown type %#o", r.Name, m)
	}

	if err := setModes(r); err != nil && f.CreateDevs {
		return err
	}
	return nil
}

// Inumber and devnumbers are unique to Unix-like
// operating systems. You can not uniquely disambiguate a file in a
// Unix system with just an inumber, you need a device number too.
// To handle hard links (unique to Unix) we need to figure out if a
// given file has been seen before. To do this we see if a file has the
// same [dev,ino] tuple as one we have seen. If so, we won't bother
// reading it in.

type devInode struct {
	dev uint64
	ino uint64
}

type Recorder struct {
	inodeMap map[devInode]Info
	inumber  uint64
}

// Certain elements of the file can not be set by cpio:
// the Inode #
// the Dev
// maintaining these elements leaves us with a non-reproducible
// output stream. In this function, we figure out what inumber
// we need to use, and clear out anything we can.
// We always zero the Dev.
// We try to find the matching inode. If found, we use its inumber.
// If not, we get a new inumber for it and save the inode away.
// This eliminates two of the messier parts of creating reproducible
// output streams.
func (r *Recorder) inode(i Info) (Info, bool) {
	dvi := devInode{dev: i.Dev, ino: i.Ino}

	// The Dev has no meaning to the eventual destination;
	// for reproducibility, zero it.
	i.Dev = 0

	// The inumber has no meaning to the eventual destination.
	// Create a new one but, first, see if we are a hard link.
	if d, ok := r.inodeMap[dvi]; ok {
		i.Ino = d.Ino
		return i, true
	}

	// Make synthetic inumbers always start at 1.
	r.inumber++
	i.Ino = r.inumber
	r.inodeMap[dvi] = i

	return i, false
}

func newLazyFile(name string) io.ReaderAt {
	return uio.NewLazyOpenerAt(func() (io.ReaderAt, error) {
		return os.Open(name)
	})
}

// GetRecord returns a cpio Record for the given path on the local file system.
//
// GetRecord does not follow symlinks. If path is a symlink, the record
// returned will reflect that symlink.
func (r *Recorder) GetRecord(path string) (Record, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return Record{}, err
	}

	sys := fi.Sys().(*syscall.Stat_t)
	info, done := r.inode(sysInfo(path, sys))

	switch fi.Mode() & os.ModeType {
	case 0: // Regular file.
		if done {
			return Record{Info: info}, nil
		}
		return Record{Info: info, ReaderAt: newLazyFile(path)}, nil

	case os.ModeSymlink:
		linkname, err := os.Readlink(path)
		if err != nil {
			return Record{}, err
		}
		return StaticRecord([]byte(linkname), info), nil

	default:
		return StaticRecord(nil, info), nil
	}
}

// Create a new Recorder.
//
// A recorder is a structure that contains variables used to calculate
// file parameters such as inode numbers for a CPIO file. The life-time
// of a Record structure is meant to be the same as the construction of a
// single CPIO archive. Do not reuse between CPIOs if you don't know what
// you're doing.
func NewRecorder() *Recorder {
	return &Recorder{make(map[devInode]Info), 0}
}
