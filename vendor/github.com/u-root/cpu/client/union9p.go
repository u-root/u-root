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
	"errors"
	"io"
	"os"
	"reflect"
	"syscall"

	"github.com/hugelgupf/p9/p9"
)

// Bind is a single bind
// For a given Twalk, the walk []string will be compared
// to the string slice in the Twalk. If there is match,
// the mount is called with the complete Twalk []string.
type UnionMount struct {
	walk  []string
	mount p9.File
}

// Union9P is a p9.Attacher.
type Union9P struct {
	mounts []UnionMount
}

// Union9pFID implements p9.File.
// The only operations it need implement
// are WalkGetAttr, GetAttr, Open, Walk and Readdir.
// The GetAttr is mostly a stub.
// The Open is required to properly support Readdir.
// Walk is used to walk to one of the underlying
// file systems. Readdir reads the union of the
// top level of all the underlying file systems,
// as in Plan 9. E.g., if the tables
// include home and a cpio, Readdir will return
// the top level of home and the cpio, including
// duplicates. There is no whiteout in this
// union file system.
type union9PFID struct {
	u *Union9P
	f p9.File
}

// NewUnionMount creates a new Union Mount from a
// []string and a p9.File.
func NewUnionMount(w []string, m p9.File) UnionMount {
	return UnionMount{walk: w, mount: m}
}

// NewUnion9P returns a Union9P, properly initialized,
// from a []UnionMount. Each UnionMount has a []string that defines
// a walk path.
// The []string argument is matched to each walk path in the []UnionMount
// in turn. As in Plan 9, the first match is used; if the walk to that
// server fails, the code returns the error; it does not go any further.
//
// I.e., if /home and /home/rminnich are in the table, they need
// to be in the order
// /home/rminnich
// /home
// in the case that the second mount does not include /home/rminnich.
// (it could be from a different 9p server, for example).
// Having a UnionMount with an empty []string is allowed; this will match
// any walk []string and hence acts as a default.
//
// For example, in Sidecore, the code looks like this:
// home, err := NewCPU9P(...)
// container, err := NewCPIO9(...)
// m1 := NewUnionMount([]string{"/home"}, home)
// m2 := NewUnionMount([]string{}, container)
// u := NewUnion9P([]UnionMount{home, container})
// This ensures /home matches first, and the container CPIO matches the rest.
//
// If a default is not
// desired, callers should only use Union Mount structs with non-empty []string.
// Only one Mount with an empty walk slice should be used, as the search will
// always stop there.
// It is allowed to have multiple Mounts for a single p9.File.
// E.g, give a p9.File, f, once can:
// m1 := NewUnionMount([]string{"/etc", f)
// m2 := NewUnionMount([]string{"/bin", f)
// u := NewUnion9P([]UnionMount{m1, m2})
// and no matter what other directories exist in f, only /etc and /bin will match.
//
// Again, to add a default case, using, e.g., another p9.File, one might have
// m1 := NewUnionMount([]string{"/etc}", f)
// m2 := NewUnionMount([]string{"/bin}", f)
// mdefault := NewUnionMount([]string{""}, fi2)
// u := NewUnion9P([]UnionMount{m1, m2, mdefault})
func NewUnion9P(mounts []UnionMount) (*Union9P, error) {
	// Index 0 is always the self pointer.
	// It matches /
	// Interesting that this is exactly
	// how it is done in the Plan 9 bind table!
	// Hand craft this one; it's kind of like
	// PID 1 in Unix.
	root := union9PFID{}
	root.f = &root
	u := &Union9P{
		mounts: append([]UnionMount{{walk: []string{"/"}, mount: &root}}, mounts...),
	}
	root.u = u

	return u, nil
}

// Attach implements p9.Attacher.Attach.
func (u *Union9P) Attach() (p9.File, error) {
	return &union9PFID{f: u.mounts[0].mount, u: u}, nil
}

var (
	_ p9.File     = &union9PFID{}
	_ p9.Attacher = &Union9P{}
)

// WalkGetAttr implements File.WalkGetAttr.
func (u *union9PFID) WalkGetAttr(names []string) ([]p9.QID, p9.File, p9.AttrMask, p9.Attr, error) {
	v("union9p:WalkGetAttr")
	q, f, err := u.Walk(names)
	if err != nil {
		return nil, nil, p9.AttrMask{}, p9.Attr{}, err
	}
	v("walk to %q got %v", names, q)
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

	_, m, a, err := f.GetAttr(valid)
	if err != nil {
		return q, f, p9.AttrMask{}, p9.Attr{}, err
	}
	v("union9p: walkgetattr returns QID %v", q)
	return q, f, m, a, err
}

// Walk implements p9.File.Walk.
func (u *union9PFID) Walk(names []string) ([]p9.QID, p9.File, error) {
	v("union9p: walk(%q)", names)
	if len(names) == 0 {
		v("union9p:clonewalk")
		return []p9.QID{{Type: p9.TypeDir, Path: 1, Version: 0}}, &union9PFID{u: u.u, f: u.f}, nil
	}
	ix := -1
	for x, bind := range u.u.mounts {
		if x == 0 {
			continue
		}
		v("union9p: bind.walk %q, names %q", bind.walk, names)
		i := len(names)
		if len(bind.walk) < i {
			i = len(bind.walk)
		}
		v("union9p:Check if bind.Walk %q == names %q", bind.walk[:i], names[:i])
		if !reflect.DeepEqual(bind.walk[:i], names[:i]) {
			v("union9p:no match")
			continue
		}
		ix = x
		v("union9p:ix is %d", ix)
		break
	}

	// this can happen if they fail to have a []string
	// as the last entry.
	if ix <= 0 {
		return nil, nil, os.ErrNotExist
	}
	v("union9p:Walk to %q from %v", names, u.u.mounts[ix])
	q, f, err := u.u.mounts[ix].mount.Walk(names)
	v("union9p:return(%v, %v, %v", q, f, err)
	return q, f, err
}

// FSync implements p9.File.FSync.
func (u *union9PFID) FSync() error {
	v("union9p:fsync")
	return nil
}

// Close implements p9.File.Close.
func (u *union9PFID) Close() error {
	v("union9p:close")
	return nil
}

// Open implements p9.File.Open.
// Basically a no op: nothing to do really.
func (u *union9PFID) Open(mode p9.OpenFlags) (p9.QID, uint32, error) {
	v("union9p:open")
	if mode.Mode() != p9.ReadOnly {
		return p9.QID{}, 0, os.ErrPermission
	}

	return p9.QID{}, 0, nil
}

// Read implements p9.File.ReadAt.
func (u *union9PFID) ReadAt(p []byte, offset int64) (int, error) {
	v("union9p:readat")
	return -1, os.ErrPermission
}

// Write implements p9.File.WriteAt.
func (u *union9PFID) WriteAt(p []byte, offset int64) (int, error) {
	v("union9p:writeat")
	return -1, os.ErrPermission
}

// Create implements p9.File.Create.
func (u *union9PFID) Create(name string, mode p9.OpenFlags, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.File, p9.QID, uint32, error) {
	v("union9p:create")
	return nil, p9.QID{Type: p9.TypeDir, Path: 0x09109, Version: 0x314}, 0555, os.ErrPermission
}

// Mkdir implements p9.File.Mkdir.
//
// Not properly implemented.
func (u *union9PFID) Mkdir(name string, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.QID, error) {
	v("union9p:mkdir")
	return p9.QID{}, os.ErrPermission
}

// Symlink implements p9.File.Symlink.
//
// Not properly implemented.
func (u *union9PFID) Symlink(oldname string, newname string, _ p9.UID, _ p9.GID) (p9.QID, error) {
	v("union9p:symlink")
	return p9.QID{}, os.ErrPermission
}

// Link implements p9.File.Link.
func (u *union9PFID) Link(target p9.File, newname string) error {
	v("union9p:link")
	return os.ErrPermission
}

// Readdir implements p9.File.Readdir.
func (u *union9PFID) Readdir(offset uint64, count uint32) (p9.Dirents, error) {
	v("union9p:readdir u %v", u)
	v("union9p:readdir u %v", u.u)
	v("union9p:readdir u %v", u.u.mounts)
	var errs error
	var all p9.Dirents
	// There can only be on '.'. But that is the ONLY one we elide
	var dot bool
	for _, bind := range u.u.mounts[1:] {
		_, dir, err := bind.mount.Walk([]string{})
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		// Must open and close each time. But it's cheap.
		if _, _, err := dir.Open(0); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		defer dir.Close()
		d, err := dir.Readdir(offset, count)
		if err != nil {
			errs = errors.Join(errs, err)
		}
		v("union9p:readdir %v", dir)
		for _, de := range d {
			if de.Name == "." {
				if dot {
					continue
				}
				dot = true
			}
			all = append(all, de)
		}
	}

	if offset >= uint64(len(all)) {
		return nil, io.EOF
	}
	v("union9p:%q, errs %v", all, errs)
	return all, errs
}

// Readlink implements p9.File.Readlink.
func (u *union9PFID) Readlink() (string, error) {
	v("union9p:readlink")
	return "", os.ErrPermission
}

// Flush implements p9.File.Flush.
func (u *union9PFID) Flush() error {
	v("union9p:flush")
	return nil
}

// Renamed implements p9.File.Renamed.
func (u *union9PFID) Renamed(parent p9.File, newName string) {
}

// UnlinkAt implements p9.File.UnlinkAt.
func (u *union9PFID) UnlinkAt(name string, flags uint32) error {
	v("union9p:unlinkat")
	return os.ErrPermission
}

// Mknod implements p9.File.Mknod.
func (*union9PFID) Mknod(name string, mode p9.FileMode, major uint32, minor uint32, _ p9.UID, _ p9.GID) (p9.QID, error) {
	v("union9p:mknod")
	return p9.QID{}, syscall.ENOSYS
}

// Rename implements p9.File.Rename.
func (*union9PFID) Rename(directory p9.File, name string) error {
	v("union9p:rename")
	return syscall.ENOSYS
}

// RenameAt implements p9.File.RenameAt.
func (u *union9PFID) RenameAt(oldName string, newDir p9.File, newName string) error {
	v("union9p:renameat")
	return syscall.ENOSYS
}

// StatFS implements p9.File.StatFS.
func (*union9PFID) StatFS() (p9.FSStat, error) {
	v("union9p:statfs")
	return p9.FSStat{}, syscall.ENOSYS
}

// SetAttr implements SetAttr.
func (u *union9PFID) SetAttr(mask p9.SetAttrMask, attr p9.SetAttr) error {
	v("union9p:setattr")
	return os.ErrPermission
}

// Lock implements lock by doing nothing.
func (u *union9PFID) Lock(pid int, locktype p9.LockType, flags p9.LockFlags, start, length uint64, client string) (p9.LockStatus, error) {
	return 0, nil
}

// GetAttr implements p9.File.GetAttr.
func (u *union9PFID) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	v("union9p: getattr")
	attr := p9.Attr{
		Mode:             p9.FileMode(0777) | p9.ModeDirectory,
		UID:              p9.UID(0),
		GID:              p9.GID(0),
		NLink:            p9.NLink(1 + len(u.u.mounts)),
		RDev:             p9.Dev(0),
		Size:             uint64(0),
		BlockSize:        uint64(4096),
		Blocks:           uint64(0),
		ATimeSeconds:     uint64(0),
		ATimeNanoSeconds: uint64(0),
		MTimeSeconds:     uint64(0),
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

	return p9.QID{Type: p9.TypeDir, Path: 0, Version: 0}, valid, attr, nil
}

// SetXattr implements p9.File.SetXattr
func (u *union9PFID) SetXattr(attr string, data []byte, flags p9.XattrFlags) error {
	return syscall.ENOSYS
}

// ListXattrs implements p9.File.ListXattrs
func (u *union9PFID) ListXattrs() ([]string, error) {
	return nil, syscall.ENOSYS
}

// GetXattr implements p9.File.GetXattr
func (u *union9PFID) GetXattr(attr string) ([]byte, error) {
	return nil, syscall.ENOSYS
}

// RemoveXattr implements p9.File.RemoveXattr
func (u *union9PFID) RemoveXattr(attr string) error {
	return syscall.ENOSYS
}
