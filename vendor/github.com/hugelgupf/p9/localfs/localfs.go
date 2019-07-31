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

package localfs

import (
	"log"
	"os"
	"path"
	"syscall"

	"github.com/hugelgupf/p9/p9"
	"github.com/hugelgupf/p9/unimplfs"
)

// Local is a p9.Attacher.
type Local struct {
	p9.DefaultWalkGetAttr
	unimplfs.NoopFile

	path string
	file *os.File
}

// Attach implements p9.Attacher.Attach.
func (l *Local) Attach() (p9.File, error) {
	return &Local{path: "/"}, nil
}

var (
	_ p9.File     = &Local{}
	_ p9.Attacher = &Local{}
)

// info constructs a QID for this file.
func (l *Local) info() (p9.QID, os.FileInfo, error) {
	var (
		qid p9.QID
		fi  os.FileInfo
		err error
	)

	// Stat the file.
	if l.file != nil {
		fi, err = l.file.Stat()
	} else {
		fi, err = os.Lstat(l.path)
	}
	if err != nil {
		log.Printf("error stating %#v: %v", l, err)
		return qid, nil, err
	}

	// Construct the QID type.
	qid.Type = p9.ModeFromOS(fi.Mode()).QIDType()

	// Save the path from the Ino.
	qid.Path = fi.Sys().(*syscall.Stat_t).Ino
	return qid, fi, nil
}

// Walk implements p9.File.Walk.
func (l *Local) Walk(names []string) ([]p9.QID, p9.File, error) {
	var qids []p9.QID
	last := &Local{path: l.path}
	for _, name := range names {
		c := &Local{path: path.Join(last.path, name)}
		qid, _, err := c.info()
		if err != nil {
			return nil, nil, err
		}
		qids = append(qids, qid)
		last = c
	}
	return qids, last, nil
}

// FSync implements p9.File.FSync.
func (l *Local) FSync() error {
	return l.file.Sync()
}

// GetAttr implements p9.File.GetAttr.
//
// Not fully implemented.
func (l *Local) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	qid, fi, err := l.info()
	if err != nil {
		return qid, p9.AttrMask{}, p9.Attr{}, err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	attr := p9.Attr{
		Mode:             p9.FileMode(stat.Mode),
		UID:              p9.UID(stat.Uid),
		GID:              p9.GID(stat.Gid),
		NLink:            stat.Nlink,
		RDev:             stat.Rdev,
		Size:             uint64(stat.Size),
		BlockSize:        uint64(stat.Blksize),
		Blocks:           uint64(stat.Blocks),
		ATimeSeconds:     uint64(stat.Atim.Sec),
		ATimeNanoSeconds: uint64(stat.Atim.Nsec),
		MTimeSeconds:     uint64(stat.Mtim.Sec),
		MTimeNanoSeconds: uint64(stat.Mtim.Nsec),
		CTimeSeconds:     uint64(stat.Ctim.Sec),
		CTimeNanoSeconds: uint64(stat.Ctim.Nsec),
	}
	valid := p9.AttrMask{
		Mode:   true,
		UID:    true,
		GID:    true,
		NLink:  true,
		RDev:   true,
		Size:   true,
		Blocks: true,
		ATime:  true,
		MTime:  true,
		CTime:  true,
	}

	return qid, valid, attr, nil
}

// Close implements p9.File.Close.
func (l *Local) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Open implements p9.File.Open.
func (l *Local) Open(mode p9.OpenFlags) (p9.QID, uint32, error) {
	qid, _, err := l.info()
	if err != nil {
		return qid, 0, err
	}

	// Do the actual open.
	f, err := os.OpenFile(l.path, int(mode), 0)
	if err != nil {
		return qid, 0, err
	}
	l.file = f

	return qid, 4096, nil
}

// Read implements p9.File.Read.
func (l *Local) ReadAt(p []byte, offset uint64) (int, error) {
	return l.file.ReadAt(p, int64(offset))
}

// Write implements p9.File.Write.
func (l *Local) WriteAt(p []byte, offset uint64) (int, error) {
	return l.file.WriteAt(p, int64(offset))
}

// Create implements p9.File.Create.
func (l *Local) Create(name string, mode p9.OpenFlags, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.File, p9.QID, uint32, error) {
	f, err := os.OpenFile(l.path, int(mode)|syscall.O_CREAT|syscall.O_EXCL, os.FileMode(permissions))
	if err != nil {
		return nil, p9.QID{}, 0, err
	}

	l2 := &Local{path: path.Join(l.path, name), file: f}
	qid, _, err := l2.info()
	if err != nil {
		l2.Close()
		return nil, p9.QID{}, 0, err
	}

	return l2, qid, 4096, nil
}

// Mkdir implements p9.File.Mkdir.
//
// Not properly implemented.
func (l *Local) Mkdir(name string, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.QID, error) {
	if err := os.Mkdir(path.Join(l.path, name), os.FileMode(permissions)); err != nil {
		return p9.QID{}, err
	}

	// Blank QID.
	return p9.QID{}, nil
}

// Symlink implements p9.File.Symlink.
//
// Not properly implemented.
func (l *Local) Symlink(oldname string, newname string, _ p9.UID, _ p9.GID) (p9.QID, error) {
	if err := os.Symlink(oldname, path.Join(l.path, newname)); err != nil {
		return p9.QID{}, err
	}

	// Blank QID.
	return p9.QID{}, nil
}

// Link implements p9.File.Link.
//
// Not properly implemented.
func (l *Local) Link(target p9.File, newname string) error {
	return os.Link(target.(*Local).path, path.Join(l.path, newname))
}

// Readdir implements p9.File.Readdir.
func (l *Local) Readdir(offset uint64, count uint32) ([]p9.Dirent, error) {
	// We only do *all* dirents in single shot.
	const maxDirentBuffer = 1024 * 1024
	buf := make([]byte, maxDirentBuffer)
	n, err := syscall.ReadDirent(int(l.file.Fd()), buf)
	if err != nil {
		// Return zero entries.
		return nil, nil
	}

	// Parse the entries; note that we read up to offset+count here.
	_, newCount, newNames := syscall.ParseDirent(buf[:n], int(offset)+int(count), nil)
	var dirents []p9.Dirent
	for i := int(offset); i >= 0 && i < newCount; i++ {
		entry := Local{path: path.Join(l.path, newNames[i])}
		qid, _, err := entry.info()
		if err != nil {
			continue
		}
		dirents = append(dirents, p9.Dirent{
			QID:    qid,
			Type:   qid.Type,
			Name:   newNames[i],
			Offset: uint64(i + 1),
		})
	}

	return dirents, nil
}

// Readlink implements p9.File.Readlink.
//
// Not properly implemented.
func (l *Local) Readlink() (string, error) {
	return os.Readlink(l.path)
}

// Flush implements p9.File.Flush.
func (l *Local) Flush() error {
	return nil
}

// Renamed implements p9.File.Renamed.
func (l *Local) Renamed(parent p9.File, newName string) {
	l.path = path.Join(parent.(*Local).path, newName)
}

// Allocate implements p9.File.Allocate.
func (l *Local) Allocate(mode p9.AllocateMode, offset, length uint64) error {
	return syscall.Fallocate(int(l.file.Fd()), mode.ToLinux(), int64(offset), int64(length))
}
