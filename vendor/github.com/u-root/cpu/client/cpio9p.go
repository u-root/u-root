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

package client

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/hugelgupf/p9/fsimpl/templatefs"
	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/u-root/pkg/cpio"
)

// CPIO9P is a p9.Attacher.
type CPIO9P struct {
	p9.DefaultWalkGetAttr

	rr   cpio.RecordReader
	m    map[string]uint64
	recs []cpio.Record
}

// CPIO9PFID defines a FID.
// It kind of sucks because it has a pointer
// for every FID. Luckily they go away when clunked.
type CPIO9PFID struct {
	p9.DefaultWalkGetAttr
	templatefs.XattrUnimplemented
	templatefs.NilCloser
	templatefs.NilSyncer
	templatefs.NoopRenamed

	fs   *CPIO9P
	path uint64
}

// NewCPIO9P returns a CPIO9P, properly initialized, from a path.
func NewCPIO9P(c string) (*CPIO9P, error) {
	f, err := os.Open(c)
	if err != nil {
		return nil, err
	}
	return NewCPIO9PReaderAt(f)
}

// NewCPIO9PReaderAt returns a CPIO9P, properly initialized, from an io.ReaderAt.
func NewCPIO9PReaderAt(r io.ReaderAt) (*CPIO9P, error) {
	archive, err := cpio.Format("newc")
	if err != nil {
		return nil, err
	}

	rr := archive.Reader(r)

	recs, err := cpio.ReadAllRecords(rr)
	if len(recs) == 0 {
		return nil, fmt.Errorf("cpio:No records: %w", os.ErrInvalid)
	}

	if err != nil {
		return nil, err
	}

	m := map[string]uint64{}
	for i, r := range recs {
		v("put %s in %d", r.Info.Name, i)
		m[r.Info.Name] = uint64(i)
	}

	return &CPIO9P{rr: rr, recs: recs, m: m}, nil
}

// Attach implements p9.Attacher.Attach.
// Only works for root.
func (s *CPIO9P) Attach() (p9.File, error) {
	return &CPIO9PFID{fs: s, path: 0}, nil
}

var (
	_ p9.File     = &CPIO9PFID{}
	_ p9.Attacher = &CPIO9P{}
)

func (l *CPIO9PFID) rec() (*cpio.Record, error) {
	if int(l.path) > len(l.fs.recs) {
		return nil, os.ErrNotExist
	}
	v("cpio:rec for %v is %v", l, l.fs.recs[l.path])
	return &l.fs.recs[l.path], nil
}

// info constructs a QID for this file.
func (l *CPIO9PFID) info() (p9.QID, *cpio.Info, error) {
	var qid p9.QID

	r, err := l.rec()
	if err != nil {
		return qid, nil, err
	}

	fi := r.Info
	// Construct the QID type.
	var m = fs.FileMode(fi.Mode)
	qid.Type = p9.ModeFromOS(m).QIDType()
	// That above sequence should work. Does not. Always returns 0.
	// I've stared at the p9 code and I no longer understand why the
	// tests even pass.
	switch fi.Mode & 0xf000 {
	case 0xa000:
		qid.Type = p9.TypeSymlink
	case 0x4000:
		qid.Type = p9.TypeDir
	}
	// Save the path from the Ino.
	qid.Path = l.path
	return qid, &fi, nil
}

// Walk implements p9.File.Walk.
func (l *CPIO9PFID) Walk(names []string) ([]p9.QID, p9.File, error) {
	r, err := l.rec()
	if err != nil {
		return nil, nil, err
	}
	verbose("cpio:starting record for %v is %v", l, r)
	var qids []p9.QID
	last := &CPIO9PFID{path: l.path, fs: l.fs}
	// If the names are empty we return info for l
	// An extra stat is never hurtful; all servers
	// are a bundle of race conditions and there's no need
	// to make things worse.
	if len(names) == 0 {
		c := &CPIO9PFID{path: last.path, fs: l.fs}
		qid, fi, err := c.info()
		verbose("cpio:Walk to %v: %v, %v, %v", *c, qid, fi, err)
		if err != nil {
			return nil, nil, err
		}
		qids = append(qids, qid)
		verbose("cpio:Walk: return %v, %v, nil", qids, last)
		return qids, last, nil
	}
	verbose("cpio:Walk: %v", names)
	// I've messed this up a few times.
	// If you start with the QID for r, you will be adding '.', which is
	// likely not what you want. the first step is to do the lookup of the
	// first name component.
	var fullpath string
	if r.Name != "." {
		fullpath = r.Name
	}
	for _, name := range names {
		fullpath = filepath.Join(fullpath, name)
		ix, ok := l.fs.m[fullpath]
		verbose("cpio:Walk %q get %v, %v", fullpath, ix, ok)
		if !ok {
			return nil, nil, os.ErrNotExist
		}
		c := &CPIO9PFID{path: ix, fs: l.fs}
		qid, fi, err := c.info()
		verbose("cpio:Walk to %q from %v: %v, %v, %v", fullpath, r, qid, fi, ok)
		if err != nil {
			return nil, nil, err
		}
		qids = append(qids, qid)
		last.path = ix
	}
	verbose("cpio:Walk: return %v, %v, nil", qids, last)
	return qids, last, nil
}

// Open implements p9.File.Open.
func (l *CPIO9PFID) Open(mode p9.OpenFlags) (p9.QID, uint32, error) {
	qid, fi, err := l.info()
	verbose("cpio:Open %v: (%v, %v, %v", *l, qid, fi, err)
	if err != nil {
		return qid, 0, err
	}

	if mode.Mode() != p9.ReadOnly {
		return qid, 0, os.ErrPermission
	}

	// Do the actual open.
	// from DIOD
	// if iounit=0, v9fs will use msize-P9_IOHDRSZ
	verbose("cpio:Open returns %v, 0, nil", qid)
	return qid, 0, nil
}

// ReadAt implements p9.File.ReadAt.
func (l *CPIO9PFID) ReadAt(p []byte, offset int64) (int, error) {
	r, err := l.rec()
	if err != nil {
		return -1, err
	}
	return r.ReadAt(p, offset)
}

// WriteAt implements p9.File.WriteAt.
func (l *CPIO9PFID) WriteAt(p []byte, offset int64) (int, error) {
	return -1, os.ErrPermission
}

// Create implements p9.File.Create.
func (l *CPIO9PFID) Create(name string, mode p9.OpenFlags, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.File, p9.QID, uint32, error) {
	return nil, p9.QID{}, 0, os.ErrPermission
}

// Mkdir implements p9.File.Mkdir.
//
// Not properly implemented.
func (l *CPIO9PFID) Mkdir(name string, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, os.ErrPermission
}

// Symlink implements p9.File.Symlink.
//
// Not properly implemented.
func (l *CPIO9PFID) Symlink(oldname string, newname string, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, os.ErrPermission
}

// Link implements p9.File.Link.
//
// Not properly implemented.
func (l *CPIO9PFID) Link(target p9.File, newname string) error {
	return os.ErrPermission
}

func (l *CPIO9PFID) readdir() ([]uint64, error) {
	verbose("cpio:readdir at %d", l.path)
	r, err := l.rec()
	if err != nil {
		return nil, err
	}
	dn := r.Info.Name
	verbose("cpio:readdir starts from %v %v", l, r)
	// while the name is a prefix of the records we are scanning,
	// append the record.
	// This can not be returned as a range as we do not want
	// contents of all subdirs.
	var list []uint64
	for i, r := range l.fs.recs[l.path+1:] {
		// filepath.Rel fails, we're done here.
		b, err := filepath.Rel(dn, r.Name)
		if err != nil {
			verbose("cpio:r.Name %q: DONE", r.Name)
			break
		}
		dir, _ := filepath.Split(b)
		if len(dir) > 0 {
			continue
		}
		verbose("cpio:readdir: %v", i)
		list = append(list, uint64(i)+l.path+1)
	}
	return list, nil
}

// Readdir implements p9.File.Readdir.
// This is a bit of a mess in cpio, but the good news is that
// files will be in some sort of order ...
func (l *CPIO9PFID) Readdir(offset uint64, count uint32) (p9.Dirents, error) {
	qid, _, err := l.info()
	if err != nil {
		return nil, err
	}
	list, err := l.readdir()
	if err != nil {
		return nil, err
	}
	if offset > uint64(len(list)) {
		return nil, io.EOF
	}
	verbose("cpio:readdir list %v", list)
	var dirents p9.Dirents
	dirents = append(dirents, p9.Dirent{
		QID:    qid,
		Type:   qid.Type,
		Name:   ".",
		Offset: l.path,
	})
	verbose("cpio:add path %d '.'", l.path)
	//log.Printf("cpio:readdir %q returns %d entries start at offset %d", l.path, len(fi), offset)
	for _, i := range list[offset:] {
		entry := CPIO9PFID{path: i, fs: l.fs}
		qid, _, err := entry.info()
		if err != nil {
			continue
		}
		r, err := entry.rec()
		if err != nil {
			continue
		}
		verbose("cpio:add path %d %q", i, filepath.Base(r.Info.Name))
		dirents = append(dirents, p9.Dirent{
			QID:    qid,
			Type:   qid.Type,
			Name:   filepath.Base(r.Info.Name),
			Offset: i,
		})
	}

	verbose("cpio:readdir:return %v, nil", dirents)
	return dirents, nil
}

// Readlink implements p9.File.Readlink.
func (l *CPIO9PFID) Readlink() (string, error) {
	v("cpio:readlinkat:%v", l)
	r, err := l.rec()
	if err != nil {
		return "", err
	}
	link := make([]byte, r.FileSize)
	v("cpio:readlink: %d byte link", len(link))
	if n, err := r.ReadAt(link, 0); err != nil || n != len(link) {
		v("cpio:readlink: fail with (%d,%v)", n, err)
		return "", err
	}
	v("cpio:readlink: %q", string(link))
	return string(link), nil
}

// Flush implements p9.File.Flush.
func (l *CPIO9PFID) Flush() error {
	return nil
}

// UnlinkAt implements p9.File.UnlinkAt.
func (l *CPIO9PFID) UnlinkAt(name string, flags uint32) error {
	return os.ErrPermission
}

// Mknod implements p9.File.Mknod.
func (*CPIO9PFID) Mknod(name string, mode p9.FileMode, major uint32, minor uint32, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, syscall.ENOSYS
}

// Rename implements p9.File.Rename.
func (*CPIO9PFID) Rename(directory p9.File, name string) error {
	return syscall.ENOSYS
}

// RenameAt implements p9.File.RenameAt.
// There is no guarantee that there is not a zipslip issue.
func (l *CPIO9PFID) RenameAt(oldName string, newDir p9.File, newName string) error {
	return syscall.ENOSYS
}

// StatFS implements p9.File.StatFS.
//
// Not implemented.
func (*CPIO9PFID) StatFS() (p9.FSStat, error) {
	return p9.FSStat{}, syscall.ENOSYS
}

// SetAttr implements SetAttr.
func (l *CPIO9PFID) SetAttr(mask p9.SetAttrMask, attr p9.SetAttr) error {
	return os.ErrPermission
}

// Lock implements lock by doing nothing.
func (*CPIO9PFID) Lock(pid int, locktype p9.LockType, flags p9.LockFlags, start, length uint64, client string) (p9.LockStatus, error) {
	return p9.LockStatus(0), nil
}

// GetAttr implements p9.File.GetAttr.
//
// Not fully implemented.
func (l *CPIO9PFID) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	qid, fi, err := l.info()
	if err != nil {
		return qid, p9.AttrMask{}, p9.Attr{}, err
	}

	//you are not getting symlink!
	attr := p9.Attr{
		Mode:             p9.FileMode(fi.Mode),
		UID:              p9.UID(fi.UID),
		GID:              p9.GID(fi.GID),
		NLink:            p9.NLink(fi.NLink),
		RDev:             p9.Dev(fi.Dev),
		Size:             uint64(fi.FileSize),
		BlockSize:        uint64(4096),
		Blocks:           uint64(fi.FileSize / 4096),
		ATimeSeconds:     uint64(0),
		ATimeNanoSeconds: uint64(0),
		MTimeSeconds:     uint64(fi.MTime),
		MTimeNanoSeconds: uint64(0),
		CTimeSeconds:     0,
		CTimeNanoSeconds: 0,
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
