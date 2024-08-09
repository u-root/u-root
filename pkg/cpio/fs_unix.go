// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package cpio

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/upath"
	"github.com/u-root/uio/uio"
	"golang.org/x/sys/unix"
)

var modeMap = map[uint64]os.FileMode{
	modeSocket:  os.ModeSocket,
	modeSymlink: os.ModeSymlink,
	modeFile:    0,
	modeBlock:   os.ModeDevice,
	modeDir:     os.ModeDir,
	modeChar:    os.ModeCharDevice,
	modeFIFO:    os.ModeNamedPipe,
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
	return 0, fmt.Errorf("invalid file type %#o", m&modeTypeMask)
}

// CreateFile creates a local file for f relative to the current working
// directory.
//
// CreateFile will attempt to set all metadata for the file, including
// ownership, times, and permissions.
func CreateFile(f Record) error {
	return CreateFileInRoot(f, ".", true)
}

// CreateFileInRoot creates a local file for f relative to rootDir.
//
// It will attempt to set all metadata for the file, including ownership,
// times, and permissions. If these fail, it only returns an error if
// forcePriv is true.
//
// Block and char device creation will only return error if forcePriv is true.
func CreateFileInRoot(f Record, rootDir string, forcePriv bool) error {
	m, err := linuxModeToFileType(f.Mode)
	if err != nil {
		return err
	}

	f.Name, err = upath.SafeFilepathJoin(rootDir, f.Name)
	if err != nil {
		// The behavior is to skip files which are unsafe due to
		// zipslip, but continue extracting everything else.
		log.Printf("Warning: Skipping file %q due to: %v", f.Name, err)
		return nil
	}
	dir := filepath.Dir(f.Name)
	// The problem: many cpio archives do not specify the directories and
	// hence the permissions. They just specify the whole path.  In order
	// to create files in these directories, we have to make them at least
	// mode 755.
	if _, err := os.Stat(dir); os.IsNotExist(err) && len(dir) > 0 {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("CreateFileInRoot %q: %w", f.Name, err)
		}
	}

	switch m {
	case os.ModeSocket, os.ModeNamedPipe:
		return fmt.Errorf("%q: type %v: cannot create IPC endpoints", f.Name, m)

	case os.ModeSymlink:
		content, err := io.ReadAll(uio.Reader(f))
		if err != nil {
			return err
		}
		return os.Symlink(string(content), f.Name)

	case os.FileMode(0):
		nf, err := os.Create(f.Name)
		if err != nil {
			return err
		}
		defer nf.Close()
		if _, err := io.Copy(nf, uio.Reader(f)); err != nil {
			return err
		}

	case os.ModeDir:
		if err := os.MkdirAll(f.Name, toFileMode(f)); err != nil {
			return err
		}

	case os.ModeDevice:
		if err := mknod(f.Name, perm(f)|syscall.S_IFBLK, dev(f)); err != nil && forcePriv {
			return err
		}

	case os.ModeCharDevice:
		if err := mknod(f.Name, perm(f)|syscall.S_IFCHR, dev(f)); err != nil && forcePriv {
			return err
		}

	default:
		return fmt.Errorf("%v: Unknown type %#o", f.Name, m)
	}

	if err := setModes(f); err != nil && forcePriv {
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

// A Recorder is a structure that contains variables used to calculate
// file parameters such as inode numbers for a CPIO file. The life-time
// of a Record structure is meant to be the same as the construction of a
// single CPIO archive. Do not reuse between CPIOs if you don't know what
// you're doing.
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
// The second return value indicates whether it is a hardlink or not.
func (r *Recorder) inode(i Info) (Info, bool) {
	d := devInode{dev: i.Dev, ino: i.Ino}
	i.Dev = 0

	if d, ok := r.inodeMap[d]; ok {
		i.Ino = d.Ino
		return i, d.Name != i.Name
	}

	i.Ino = r.inumber
	r.inumber++
	r.inodeMap[d] = i

	return i, false
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
		return Record{Info: info, ReaderAt: uio.NewLazyLimitFile(path, int64(info.FileSize))}, nil

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

// NewRecorder creates a new Recorder.
//
// A recorder is a structure that contains variables used to calculate
// file parameters such as inode numbers for a CPIO file. The life-time
// of a Record structure is meant to be the same as the construction of a
// single CPIO archive. Do not reuse between CPIOs if you don't know what
// you're doing.
func NewRecorder() *Recorder {
	return &Recorder{make(map[devInode]Info), 2}
}

// LSInfoFromRecord converts a Record to be usable with the ls package for
// listing files.
func LSInfoFromRecord(rec Record) ls.FileInfo {
	var target string

	mode := modeFromLinux(rec.Mode)
	if mode&os.ModeType == os.ModeSymlink {
		if l, err := uio.ReadAll(rec); err != nil {
			target = err.Error()
		} else {
			target = string(l)
		}
	}

	return ls.FileInfo{
		Name:          rec.Name,
		Mode:          mode,
		Rdev:          unix.Mkdev(uint32(rec.Rmajor), uint32(rec.Rminor)),
		UID:           uint32(rec.UID),
		GID:           uint32(rec.GID),
		Size:          int64(rec.FileSize),
		MTime:         time.Unix(int64(rec.MTime), 0).UTC(),
		SymlinkTarget: target,
	}
}
