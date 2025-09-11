package client

import (
	"os"
	"time"

	"github.com/go-git/go-billy/v5"
	osfs "github.com/go-git/go-billy/v5/osfs"
)

// NewOSFS returns a billy.FileSystem for a path.
func NewOSFS(r string) billy.Filesystem {
	bfs := osfs.New(r, osfs.WithBoundOS())
	return bfs
}

// COS or OSFS + Change wraps a billy.FS to not fail the `Change` interface.
type COS struct {
	billy.Filesystem
}

// Chmod changes mode
func (fs COS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(fs.Join(fs.Root(), name), mode)
}

// Lchown changes ownership. As in lchown(2), it does not
// follow a symlink.
func (fs COS) Lchown(name string, uid, gid int) error {
	return os.Lchown(fs.Join(fs.Root(), name), uid, gid)
}

// Chown changes ownership
func (fs COS) Chown(name string, uid, gid int) error {
	return os.Chown(fs.Join(fs.Root(), name), uid, gid)
}

// Chtimes changes access time
func (fs COS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(fs.Join(fs.Root(), name), atime, mtime)
}
