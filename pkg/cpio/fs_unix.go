// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
		if err := os.Symlink(string(content), f.Name); err != nil {
			return err
		}
		return setModes(f)

	default:
		return fmt.Errorf("%v: Unknown type %#o", f.Name, m)
	}
}

func GetRecord(path string) (Record, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return Record{}, err
	}

	sys := fi.Sys().(*syscall.Stat_t)
	info := sysInfo(path, sys)

	switch fi.Mode() & os.ModeType {
	case 0: // Regular file.
		osfile, err := os.Open(path)
		if err != nil {
			return Record{}, err
		}
		defer osfile.Close()

		contents, err := ioutil.ReadAll(osfile)
		if err != nil {
			return Record{}, err
		}
		return Record{bytes.NewReader(contents), info}, nil

	case os.ModeSymlink:
		linkname, err := os.Readlink(path)
		if err != nil {
			return Record{}, err
		}
		return StaticRecord(linkname, info), nil

	default:
		return EmptyRecord(info), nil
	}
}
