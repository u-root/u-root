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

package netfuse

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/u-root/fuse/fuseops"
	"github.com/u-root/fuse/fuseutil"
)

// inode is a struct containing information about file or directory.
type inode struct {
	// The full path
	FullPath string
	// The base name
	Name string
	/////////////////////////
	// Mutable state
	/////////////////////////

	// if it's a dir, and the open has worked, this has the dents.
	dents []os.FileInfo

	// FUSE requires us to maintain a refcount and remove the inode
	// when it hits 0.
	n uint64
}

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

// Create a new inode with the supplied attributes, which need not contain
// time-related information (the inode object will take care of that).
func newInode(FullPath string) *inode {
	return &inode{
		FullPath: FullPath,
		Name:     filepath.Base(FullPath),
	}
}

func (in *inode) CheckInvariants() {
	return
}

func (in *inode) isDir() bool {
	attr, err := in.stat()
	if err != nil {
		return false
	}
	return attr.Mode&os.ModeDir == os.ModeDir
}

func (in *inode) isSymlink() bool {
	attr, err := in.stat()
	if err != nil {
		return false
	}
	return attr.Mode&os.ModeSymlink == os.ModeSymlink
}

func (in *inode) isFile() bool {
	return !(in.isDir() || in.isSymlink())
}

////////////////////////////////////////////////////////////////////////
// Public methods
////////////////////////////////////////////////////////////////////////

// Return the number of children of the directory.
//
// REQUIRES: in.isDir()
func (in *inode) Len() (n int, err error) {
	i, err := ioutil.ReadDir(in.FullPath)
	if err != nil {
		return -1, err
	}
	return len(i), nil
}

// LookUp finds an entry for the given child name and return its inode ID, attributs, and type.
func LookUp(name string) (fuseops.InodeID, *fuseops.InodeAttributes, fuseutil.DirentType, error) {
	fi, err := os.Lstat(name)
	FSDebug("Stat %v: %v %v", name, fi, err)
	if err != nil {
		return 0, nil, 0, err
	}
	sys := fi.Sys().(*syscall.Stat_t)
	id := sys.Ino
	typ := fuseutil.DirentType(fi.Mode() & os.ModeType)

	attr := &fuseops.InodeAttributes{
		Size:  uint64(fi.Size()),
		Nlink: uint32(sys.Nlink),
		Mode:  fi.Mode(),
		Atime: time.Unix(sys.Atim.Sec, sys.Atim.Nsec),
		Mtime: time.Unix(sys.Mtim.Sec, sys.Mtim.Nsec),
		Ctime: time.Unix(sys.Ctim.Sec, sys.Ctim.Nsec),
		// This is a tough one. The tests want to test this
		// but it's not really in linux. I'll go with mtime?
		Crtime: fi.ModTime(),
		Uid:    sys.Uid,
		Gid:    sys.Gid,
	}
	return fuseops.InodeID(id), attr, typ, nil
}

// Rename renames an inode.
func (in *inode) Rename(fullpath string) {
	in.FullPath = fullpath
	in.Name = filepath.Base(fullpath)
}

// SetAttributes sets attributes for some or all parameters. If a parameter is nil, it is not changed.
func (in *inode) SetAttributes(size *uint64, mode *os.FileMode, mtime *time.Time) error {

	// Truncate?
	if size != nil {
		if err := os.Truncate(in.FullPath, int64(*size)); err != nil {
			return err
		}
	}

	// Change mode?
	if mode != nil {
		if err := os.Chmod(in.FullPath, *mode); err != nil {
			return err
		}
	}

	// Change mtime?
	if mtime != nil {
		// We're going to change atime too, since that's an acdess?
		if err := os.Chtimes(in.FullPath, time.Now(), *mtime); err != nil {
			return err
		}
	}
	return nil
}

// forget decreaments the refcount of an inode.
func (in *inode) forget(N uint64) uint64 {
	in.n -= N
	return in.n
}

// ReadDir returns the []os.FileInfo for an inode.
func (in *inode) ReadDir() ([]os.FileInfo, error) {
	return in.dents, nil
}

// stat returns the InodeAttributes for an inode.
// We stat every time. It's cheap and it's important to stay on top of existence and
// changes.
func (in *inode) stat() (*fuseops.InodeAttributes, error) {
	fi, err := os.Lstat(in.FullPath)
	if err != nil {
		// it can be removed without us knowing.
		// We do not remove it from our inode cache until
		// we get a Forget
		return nil, err
	}
	sys := fi.Sys().(*syscall.Stat_t)
	// Fill in the response.
	return &fuseops.InodeAttributes{
		Size:  uint64(fi.Size()),
		Nlink: uint32(sys.Nlink),
		Mode:  fi.Mode(),
		Uid:   sys.Uid,
		Gid:   sys.Gid,
	}, nil

}
