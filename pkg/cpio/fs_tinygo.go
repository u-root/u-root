// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows && tinygo

package cpio

import (
	"fmt"
	"io"
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
func setModes(name string, r Record) error {
	if err := os.Chmod(name, toFileMode(r)&os.ModePerm); err != nil {
		return err
	}
	if err := os.Chown(name, int(r.UID), int(r.GID)); err != nil {
		return err
	}
	if err := os.Chmod(name, toFileMode(r)); err != nil {
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
func CreateFile(f Record) error {
	return CreateFileInRoot(f, ".", true)
}

// CreateFileInRoot creates a local file for f relative to rootDir.
func CreateFileInRoot(f Record, rootDir string, forcePriv bool) error {
	m, err := linuxModeToFileType(f.Mode)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return err
	}

	targetName, err := upath.SafeFilepathJoin(rootDir, f.Name)
	if err != nil {
		return err
	}

	// Create parent directories
	if dir := filepath.Dir(targetName); dir != "." {
		if _, err := os.Stat(dir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
	}

	switch m {
	case os.ModeSocket, os.ModeNamedPipe:
		return fmt.Errorf("%q: type %v: cannot create IPC endpoints", targetName, m)

	case os.ModeSymlink:
		content, err := io.ReadAll(uio.Reader(f))
		if err != nil {
			return err
		}
		return os.Symlink(string(content), targetName)

	case os.FileMode(0):
		nf, err := os.OpenFile(targetName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, toFileMode(f).Perm())
		if err != nil {
			return err
		}
		defer nf.Close()
		if _, err := io.Copy(nf, uio.Reader(f)); err != nil {
			return err
		}

	case os.ModeDir:
		if err := os.MkdirAll(targetName, toFileMode(f)); err != nil {
			return err
		}

	case os.ModeDevice:
		if err := mknod(targetName, perm(f)|syscall.S_IFBLK, dev(f)); err != nil && forcePriv {
			return err
		}

	case os.ModeCharDevice:
		if err := mknod(targetName, perm(f)|syscall.S_IFCHR, dev(f)); err != nil && forcePriv {
			return err
		}

	default:
		return fmt.Errorf("%v: Unknown type %#o", targetName, m)
	}

	if err := setModes(targetName, f); err != nil && forcePriv {
		return err
	}
	return nil
}

type devInode struct {
	dev uint64
	ino uint64
}

type Recorder struct {
	inodeMap map[devInode]Info
	inumber  uint64
}

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
