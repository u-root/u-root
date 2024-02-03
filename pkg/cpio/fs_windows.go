// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
)

// A Recorder is a structure that contains variables used to calculate
// file parameters such as inode numbers for a CPIO file. The life-time
// of a Record structure is meant to be the same as the construction of a
// single CPIO archive. Do not reuse between CPIOs if you don't know what
// you're doing.
type Recorder struct {
	inumber uint64
}

var modeMap = map[uint64]os.FileMode{
	modeFile: 0,
	modeDir:  os.ModeDir,
}

func unixModeToFileType(m uint64) (os.FileMode, error) {
	if t, ok := modeMap[m&modeTypeMask]; ok {
		return t, nil
	}
	return 0, fmt.Errorf("invalid file type %#o", m&modeTypeMask)
}

func toFileMode(r Record) os.FileMode {
	return os.FileMode(perm(r))
}

// setModes sets the modes.
func setModes(r Record) error {
	return os.Chmod(r.Name, toFileMode(r)&os.ModePerm)
}

func perm(r Record) uint32 {
	return uint32(r.Mode) & modePermissions
}

func dev(r Record) int {
	return int(r.Rmajor<<8 | r.Rminor)
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
func CreateFileInRoot(f Record, rootDir string, forcePriv bool) error {
	m, err := unixModeToFileType(f.Mode)
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

	default:
		return fmt.Errorf("%v: Unknown type %#o", f.Name, m)
	}

	if err := setModes(f); err != nil && forcePriv {
		return err
	}
	return nil
}

func (r *Recorder) inode(i Info) Info {
	i.Ino = r.inumber
	r.inumber++
	return i
}

// GetRecord returns a cpio Record for the given path on the local file system.
//
// GetRecord does not follow symlinks. If path is a symlink, the record
// returned will reflect that symlink.
func (r *Recorder) GetRecord(path string) (Record, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Record{}, err
	}
	sys, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return Record{}, fmt.Errorf("sys is empty:%w", syscall.ENOSYS)
	}
	info := r.inode(sysInfo(path, sys))

	switch fi.Mode() & os.ModeType {
	case 0: // Regular file.
		return Record{Info: info, ReaderAt: uio.NewLazyFile(path)}, nil
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
	return &Recorder{inumber: 2}
}

// LSInfoFromRecord converts a Record to be usable with the ls package for
// listing files.
func LSInfoFromRecord(rec Record) ls.FileInfo {
	mode := modeFromLinux(rec.Mode)
	return ls.FileInfo{
		Name:  rec.Name,
		Mode:  mode,
		UID:   fmt.Sprintf("%d", rec.UID),
		Size:  int64(rec.FileSize),
		MTime: time.Unix(int64(rec.MTime), 0).UTC(),
	}
}
