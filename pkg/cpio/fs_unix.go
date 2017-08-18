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
	"time"
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

var modeMap = map[uint64]os.FileMode{
	modeSocket:  os.ModeSocket,
	modeSymlink: os.ModeSymlink,
	modeFile:    0,
	modeBlock:   os.ModeDevice,
	modeDir:     os.ModeDir,
	modeChar:    os.ModeCharDevice,
	modeFIFO:    os.ModeNamedPipe,
}

// modes sets the modes, changing the easy ones first and the harder ones last.
// In this way, we set as much as we can before bailing out. It's not an error
// to not be able to set uid and gid, at least not yet.
// For now we also ignore sticky bits.
func setModes(r Record) error {
	if err := os.Chmod(r.Name, os.FileMode(perm(r))); err != nil {
		return err
	}
	if err := os.Chtimes(r.Name, time.Time{}, time.Unix(int64(r.MTime), 0)); err != nil {
		return err
	}
	if err := os.Chown(r.Name, int(r.UID), int(r.GID)); err != nil {
		return err
	}
	// TODO: only set SUID and GUID if we can set the owner.
	return nil
}

func perm(r Record) uint32 {
	return uint32(r.Mode) & modePermissions
}

func dev(r Record) int {
	return int(r.Rmajor<<8 | r.Rminor)
}

func linuxModeToMode(m uint64) (os.FileMode, error) {
	if t, ok := modeMap[m&modeTypeMask]; ok {
		return t, nil
	}
	return 0, fmt.Errorf("Invalid file type %#o", m&modeTypeMask)
}

func CreateFile(f Record) error {
	m, err := linuxModeToMode(f.Mode)
	if err != nil {
		return err
	}

	dir, _ := filepath.Split(f.Name)
	// The problem: many cpio archives do not specify the directories
	// and hence the permissions. They just specify the whole path.
	// In order to create files in these directories, we have to make them at least
	// mode 755.
	if dir != "" {
		switch m {
		case os.FileMode(0),
			os.ModeDevice,
			os.ModeCharDevice,
			os.ModeSymlink:
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}

	switch m {
	case os.ModeSocket, os.ModeNamedPipe:
		return fmt.Errorf("%q: type %v: cannot create IPC endpoints", f.Name, m)

	case os.FileMode(0):
		nf, err := os.Create(f.Name)
		if err != nil {
			return err
		}
		defer nf.Close()
		if _, err := io.Copy(nf, f); err != nil {
			return err
		}
		return setModes(f)

	case os.ModeDir:
		if err := os.MkdirAll(f.Name, os.FileMode(perm(f))); err != nil {
			return err
		}
		return setModes(f)

	case os.ModeDevice:
		if err := syscall.Mknod(f.Name, perm(f)|syscall.S_IFBLK, dev(f)); err != nil {
			return err
		}
		return setModes(f)

	case os.ModeCharDevice:
		if err := syscall.Mknod(f.Name, perm(f)|syscall.S_IFCHR, dev(f)); err != nil {
			return err
		}
		return setModes(f)

	case os.ModeSymlink:
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		return os.Symlink(string(content), f.Name)

	default:
		return fmt.Errorf("%v: Unknown type %#o", f.Name, m)
	}
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

var (
	inodeMap = map[devInode]Info{}
	inumber  uint64
)

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
func inode(i Info) (Info, bool) {
	d := devInode{dev: i.Dev, ino: i.Ino}
	i.Dev = 0

	if d, ok := inodeMap[d]; ok {
		i.Ino = d.Ino
		return i, true
	}

	i.Ino = inumber
	inumber++
	inodeMap[d] = i

	return i, false
}

func GetRecord(path string) (Record, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return Record{}, err
	}

	sys := fi.Sys().(*syscall.Stat_t)
	info, done := inode(sysInfo(path, sys))

	switch fi.Mode() & os.ModeType {
	case 0: // Regular file.
		if done {
			return Record{Info: info}, nil
		}
		return Record{Info: info, ReadCloser: NewDeferReadCloser(path)}, nil

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
