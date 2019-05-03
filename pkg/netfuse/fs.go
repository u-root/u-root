// Copyright 2015 Google Inc. All Rights Reserved.
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

//go:generate go run gen.go

package netfuse

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jacobsa/syncutil"
	"github.com/u-root/fuse"
	"github.com/u-root/fuse/fuseops"
	"github.com/u-root/fuse/fuseutil"
	"golang.org/x/sys/unix"
)

// FS defines a file system rooted at a local directory.
// It does not implement all functions and hence has an embedded
// NotImplementedFileSystem. The way FUSE works requires that it
// maintain an inode cache from a lookup to a forget.
// This inode cache is protected by a mutex.
type FS struct {
	fuseutil.NotImplementedFileSystem

	/////////////////////////
	// Mutable state
	/////////////////////////

	mu syncutil.InvariantMutex

	// Inodes for which a lookup has been done.
	//
	// All inodes are protected by the file system mutex.
	//
	// INVARIANT: For each inode in, in.CheckInvariants() does not panic.
	// INVARIANT: len(inodes) > fuseops.RootInodeID
	// INVARIANT: For all i < fuseops.RootInodeID, inodes[i] == nil
	// INVARIANT: inodes[fuseops.RootInodeID] != nil
	// INVARIANT: inodes[fuseops.RootInodeID].isDir()
	inodes map[fuseops.InodeID]*inode // GUARDED_BY(mu)
}

// A handleID is an integer and is provided by the file system on open.
// fds are guaranteed to be unique per server. This server will be as stateless as
// possible, and we will let the kernel own these IDs.
type handleID int

// FSDebug is used to print FS debug messages.
var FSDebug = func(string, ...interface{}) {}

// NewFS creates a new file system given a root path.
// The path is not checked for validity, as it can come
// into and go out of existence any time.
func NewFS(root string) *FS {
	fs := &FS{
		inodes: make(map[fuseops.InodeID]*inode),
	}

	fs.inodes[fuseops.RootInodeID] = newInode(root)
	// Set up invariant checking.
	fs.mu = syncutil.NewInvariantMutex(fs.checkInvariants)

	return fs
}

// New creates a new file system and starts a server for it.
func New(root string) fuse.Server {
	fs := NewFS(root)
	return fuseutil.NewFileSystemServer(fs)
}

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func (fs *FS) checkInvariants() {
	// Check the root inode.
	if !fs.inodes[fuseops.RootInodeID].isDir() {
		panic("Expected root to be a directory.")
	}

	// INVARIANT: For each inode in, in.CheckInvariants() does not panic.
	for _, in := range fs.inodes {
		in.CheckInvariants()
	}
}

// getInode searches for an InodeID in the cache and returns
// syscall.ENOENT if it is not found.
//
// LOCKS_REQUIRED(fs.mu)
func (fs *FS) getInode(id fuseops.InodeID) (inode *inode, err error) {
	inode, ok := fs.inodes[id]
	if !ok {
		return nil, syscall.ENOENT
	}

	return inode, nil
}

// allocateInode allocates a new inode.
//
// LOCKS_REQUIRED(fs.mu)
func (fs *FS) allocateInode(fullpath string, id fuseops.InodeID) (*inode, error) {
	if i, ok := fs.inodes[id]; ok {
		return i, nil
	}
	// Create the inode.
	in := newInode(fullpath)
	in.n = 1
	fs.inodes[id] = in
	return in, nil
}

// deallocateInode removes an Inode by InodeID.
// LOCKS_REQUIRED(fs.mu)
func (fs *FS) deallocateInode(id fuseops.InodeID) {
	delete(fs.inodes, id)
}

// newInode creates a new Inode given a name, and returns an ID and attributes, or an error.
func (fs *FS) newInode(n string) (fuseops.InodeID, *fuseops.InodeAttributes, error) {
	// Does the directory have an entry with the given name?
	childID, attrs, _, err := LookUp(n)
	if err != nil {
		FSDebug("newInode: %q returning %v", n, err)
		// You have to return fuse.ENOENT
		return 0, nil, fuse.ENOENT
	}

	// Grab the child.
	if _, err := fs.allocateInode(n, childID); err != nil {
		return 0, nil, err
	}
	return childID, attrs, nil
}

////////////////////////////////////////////////////////////////////////
// FileSystem methods
////////////////////////////////////////////////////////////////////////

// StatFS implements statfs.
func (fs *FS) StatFS(ctx context.Context, op *fuseops.StatFSOp) (err error) {
	return unix.ENOSYS
}

// LookUpInode finds an item in the file system, given a directory inode and the name.
// If the file is found, it is entered into the inode cache.
func (fs *FS) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	// Grab the parent directory.
	inode, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}

	childID, attrs, err := fs.newInode(filepath.Join(inode.FullPath, op.Name))
	if err != nil {
		return err
	}
	// Fill in the response.
	op.Entry.Child = childID
	op.Entry.Attributes = *attrs

	op.Entry.AttributesExpiration = time.Now().Add(5 * time.Minute)
	op.Entry.EntryExpiration = op.Entry.AttributesExpiration

	return nil
}

// GetInodeAttributes returns the FUSE attributes of an Inode, by converting
// os.FileInfo and Stat_t information.
func (fs *FS) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	inode, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}
	attr, err := inode.stat()
	if err != nil {
		return err
	}
	op.Attributes = *attr
	//op.Attributes.Atime = sys.Atime
	FSDebug("Inode attrs %v opattr %v", inode, op.Attributes)
	op.AttributesExpiration = time.Now().Add(5 * time.Minute)

	return nil
}

// SetInodeAttributes sets the FUSE-settable Inode attributes, including mode and owner.
func (fs *FS) SetInodeAttributes(ctx context.Context, op *fuseops.SetInodeAttributesOp) (err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	inode, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}

	inode.SetAttributes(op.Size, op.Mode, op.Mtime)

	attr, err := inode.stat()
	if err != nil {
		return err
	}
	op.Attributes = *attr

	op.AttributesExpiration = time.Now().Add(5 * time.Minute)

	return
}

// MkDir makes a directory by name relative to an Inode.
func (fs *FS) MkDir(ctx context.Context, op *fuseops.MkDirOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	parent, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}

	n := filepath.Join(parent.FullPath, op.Name)
	FSDebug("mkdir mkdir %v", n)
	if err := os.Mkdir(n, op.Mode); err != nil {
		return err
	}

	// Don't assume anything; we could be racing with something.
	// Do a full Lookup
	id, attr, _, err := LookUp(n)
	FSDebug("mkdir lookup %v", err)
	if err != nil {
		return err
	}

	_, err = fs.allocateInode(n, id)
	if err != nil {
		return err
	}

	op.Entry.Child = id
	op.Entry.Attributes = *attr
	op.Entry.AttributesExpiration = time.Now().Add(5 * time.Minute)
	op.Entry.EntryExpiration = op.Entry.AttributesExpiration

	return nil
}

// createFile is the heart of creating a file.
// LOCKS_REQUIRED(fs.mu)
func (fs *FS) createFile(parentID fuseops.InodeID, name string, mode os.FileMode) (*fuseops.ChildInodeEntry, int, error) {
	parent, err := fs.getInode(parentID)
	if err != nil {
		return nil, -1, err
	}

	n := filepath.Join(parent.FullPath, name)
	// The open mode with os.Create is implicitly O_WRONLY.
	// unix.Open with a O_RDWR|O_CREAT is what the Go runtime does,
	// and that behavior is what this FUSE package (and FUSE kernel module)
	// seem to want.
	fd, err := unix.Open(n, unix.O_CREAT|unix.O_RDWR|unix.O_TRUNC|unix.O_CLOEXEC, uint32(mode.Perm()))

	if err != nil {
		FSDebug("failed %v", err)
		return nil, -1, err
	}

	// See if it worked out.
	// Why do this instead of assuming it all worked?
	// Restrictions in our environment (umask)
	// or the underlying file system might change the attr
	// in some way. It's easiest just to ask what the attr are.
	id, attr, _, err := LookUp(n)
	if err != nil {
		FSDebug("failed %v", err)
		return nil, -1, err
	}

	if _, err := fs.allocateInode(n, id); err != nil {
		FSDebug("allocate failed %v", err)
		return nil, -1, err
	}

	e := &fuseops.ChildInodeEntry{
		Child:                id,
		Attributes:           *attr,
		AttributesExpiration: time.Now().Add(5 * time.Minute),
		EntryExpiration:      time.Now().Add(5 * time.Minute),
	}

	return e, fd, nil
}

// CreateFile create a file as defined by the CreateFileOp
func (fs *FS) CreateFile(ctx context.Context, op *fuseops.CreateFileOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	e, fd, err := fs.createFile(op.Parent, op.Name, op.Mode)
	if err != nil {
		return err
	}
	op.Entry = *e
	op.Handle = fuseops.HandleID(fd)
	return nil
}

// CreateSymlink creates a symlink relative to a directory.
func (fs *FS) CreateSymlink(ctx context.Context, op *fuseops.CreateSymlinkOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	parent, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}
	o := filepath.Join(parent.FullPath, op.Name)
	if err := os.Symlink(op.Target, o); err != nil {
		return err
	}

	childID, attrs, err := fs.newInode(o)
	if err != nil {
		return err
	}

	op.Entry.Child = childID
	op.Entry.Attributes = *attrs
	op.Entry.AttributesExpiration = time.Now().Add(5 * time.Minute)
	op.Entry.EntryExpiration = op.Entry.AttributesExpiration
	return nil
}

// CreateLink creates a hard link relative to a directory.
func (fs *FS) CreateLink(ctx context.Context, op *fuseops.CreateLinkOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	parent, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}

	// We do require that the target exist for a hard link.
	// FUSE should have looked up it beforehand.
	// It may be gone by the time we get this request.
	o, err := fs.getInode(op.Target)
	if err != nil {
		return err
	}

	attr, err := o.stat()
	if err != nil {
		return err
	}
	t := filepath.Join(parent.FullPath, op.Name)
	if err := os.Link(o.FullPath, t); err != nil {
		return err
	}

	// Subtle point for hard links.
	// Given an inode #, there is no need to remember
	// every last name the inode # might be linked to
	// in directories. So we're done.
	// I.e. we don't add to the inode table for this fs.
	// This is very different from 9p; in 9p, every QID
	// is supposed to uniquely refer to a single file.

	op.Entry.Child, op.Entry.Attributes = op.Target, *attr
	op.Entry.AttributesExpiration = time.Now().Add(5 * time.Minute)
	op.Entry.EntryExpiration = op.Entry.AttributesExpiration

	return nil
}

// Rename renames a file, using Inodes from two directories, and
// names relative to those directories.
func (fs *FS) Rename(ctx context.Context, op *fuseops.RenameOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	o, err := fs.getInode(op.OldParent)
	if err != nil {
		return err
	}

	n, err := fs.getInode(op.NewParent)
	if err != nil {
		return err
	}

	oname := filepath.Join(o.FullPath, op.OldName)
	nname := filepath.Join(n.FullPath, op.NewName)
	// If we have looked it up, we need it get its inumber
	oid, _, _, err := LookUp(oname)
	if err != nil {
		return err
	}
	of, _ := fs.getInode(oid)
	// we need the file to exist.
	FSDebug("rename at %s -> %s", oname, nname)
	if err := os.Rename(oname, nname); err != nil {
		return err
	}
	if of == nil {
		return nil
	}
	// well, problem: if the child is in our inode cache we need to rename it.
	// We used to think we could just delete it but that's not true.
	// FUSE will immediately do a GetInodeAttributes after a rename
	// and that needs this inode.
	if _, _, _, err := LookUp(nname); err == nil {
		nf, err := fs.getInode(oid)
		if err != nil {
			return nil
		}
		nf.Rename(nname)
	}
	return nil
}

// RmDir removes a directory using a path relative to another directory.
func (fs *FS) RmDir(ctx context.Context, op *fuseops.RmDirOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	parent, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(parent.FullPath, op.Name))
}

// Unlink calls os.Remove for a file/directory from a directory.
func (fs *FS) Unlink(ctx context.Context, op *fuseops.UnlinkOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	parent, err := fs.getInode(op.Parent)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(parent.FullPath, op.Name))
}

// OpenDir opens a directory. We read the entries all at once to maintain reasonable semantics
// but we also leave the fd open so we don't have fd collisions.
func (fs *FS) OpenDir(ctx context.Context, op *fuseops.OpenDirOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	var err error
	inode, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}
	// TODO: consider just doing fs.mu.Unlock here, in case we get contention.
	// defer is not always the right answer.
	fd, err := unix.Open(inode.FullPath, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}

	inode.dents, err = ioutil.ReadDir(inode.FullPath)
	op.Handle = fuseops.HandleID(fd)
	return err
}

// ReadDir returns directory entries from a directory inode.
func (fs *FS) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) (err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	inode, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}

	dents, err := inode.ReadDir()
	if err != nil || dents == nil {
		return err
	}

	lastOff := op.Offset
	op.BytesRead = 0
	for _, dent := range dents[op.Offset:] {
		de := fuseutil.Dirent{
			Offset: fuseops.DirOffset(lastOff + 1),
			Inode:  fuseops.InodeID(dent.Sys().(*syscall.Stat_t).Ino),
			Name:   dent.Name(),
			Type:   fuseutil.DirentType(os.FileMode(dent.Mode())),
		}
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], de)
		if n == 0 {
			break
		}
		lastOff++
		op.BytesRead += n
	}

	return nil
}

// OpenFile opens a file relative to the directory in op.
func (fs *FS) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	var err error

	inode, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}

	if !inode.isFile() {
		return unix.EISDIR
	}

	fd, err := unix.Open(inode.FullPath, op.Flag, 0)
	if err != nil {
		FSDebug("failed %v", err)

	}

	op.Handle = fuseops.HandleID(fd)

	return err
}

// ReadFile reads a file.
// There is no need to lock the FS; the Handle is managed in the kernel,
// entirely independently of this file system.
func (fs *FS) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	n, err := unix.Pread(int(op.Handle), op.Dst, op.Offset)
	op.BytesRead = n

	// Don't return EOF errors; we just indicate EOF to fuse using a short read.
	if err == io.EOF {
		err = nil
	}

	return err
}

// WriteFile writes a file.
func (fs *FS) WriteFile(ctx context.Context, op *fuseops.WriteFileOp) error {
	n, err := unix.Pwrite(int(op.Handle), op.Data, op.Offset)
	// for now, we show much data was written via change the data size.
	// This was not my idea :-)
	if n < 0 {
		op.Data = []byte{}
	} else {
		op.Data = op.Data[:n]
	}

	return err
}

// ReadSymlink reads a symbolic link relative to the directory.
func (fs *FS) ReadSymlink(ctx context.Context, op *fuseops.ReadSymlinkOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	in, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}
	op.Target, err = os.Readlink(in.FullPath)
	return err
}

// ForgetInode forgets an inode. If the reference count drops to less than 1,
// the inode is removed from the cache.
func (fs *FS) ForgetInode(ctx context.Context, op *fuseops.ForgetInodeOp) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	in, err := fs.getInode(op.Inode)
	if err != nil {
		return err
	}
	n := in.forget(op.N)
	if n < 1 {
		fs.deallocateInode(op.Inode)
	}
	return nil
}

// ReleaseFileHandle releases a file, which in this case means it is closed.
func (fs *FS) ReleaseFileHandle(ctx context.Context, op *fuseops.ReleaseFileHandleOp) error {
	return unix.Close(int(op.Handle))
}

// ReleaseDirHandle releases a directory handle, which in this case means it is closed.
func (fs *FS) ReleaseDirHandle(ctx context.Context, op *fuseops.ReleaseDirHandleOp) error {
	return unix.Close(int(op.Handle))
}
