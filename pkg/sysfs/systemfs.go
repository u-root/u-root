// Package sysfs provides a file system interface that allows access to the
// entire file system, including paths outside the specified root directory.
package sysfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// SystemFS implements fs.FS without path restrictions.
// Unlike os.DirFS, it allows access to paths outside the root directory.
type SystemFS struct {
	root string
}

// NewSystemFS creates a new SystemFS rooted at the given directory.
// Unlike os.DirFS, paths can escape the root directory using "../".
func NewSystemFS(root string) *SystemFS {
	return &SystemFS{root: root}
}

// Open opens the named file for reading.
func (sfs *SystemFS) Open(name string) (fs.File, error) {
	if name == "" {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	// Join the root with the name, allowing ".." to escape
	fullPath := filepath.Join(sfs.root, name)

	return os.Open(fullPath)
}

// Stat returns file information for the named file.
func (sfs *SystemFS) Stat(name string) (fs.FileInfo, error) {
	if name == "" {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(sfs.root, name)
	return os.Stat(fullPath)
}

// ReadDir reads the named directory and returns a list of directory entries.
func (sfs *SystemFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name == "" {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(sfs.root, name)
	return os.ReadDir(fullPath)
}

// ReadFile reads the named file and returns its contents.
func (sfs *SystemFS) ReadFile(name string) ([]byte, error) {
	if name == "" {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(sfs.root, name)
	return os.ReadFile(fullPath)
}

// Glob returns the names of all files matching pattern.
func (sfs *SystemFS) Glob(pattern string) ([]string, error) {
	fullPattern := filepath.Join(sfs.root, pattern)
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return nil, err
	}

	// Convert back to relative paths
	var result []string
	for _, match := range matches {
		rel, err := filepath.Rel(sfs.root, match)
		if err != nil {
			continue
		}
		// Convert to forward slashes for fs.FS compatibility
		rel = filepath.ToSlash(rel)
		result = append(result, rel)
	}

	return result, nil
}

// Sub returns an FS corresponding to the subtree rooted at dir.
func (sfs *SystemFS) Sub(dir string) (fs.FS, error) {
	if dir == "" {
		return nil, &fs.PathError{Op: "sub", Path: dir, Err: fs.ErrInvalid}
	}

	newRoot := filepath.Join(sfs.root, dir)
	return NewSystemFS(newRoot), nil
}

// Ensure SystemFS implements the required interfaces
var (
	_ fs.FS         = (*SystemFS)(nil)
	_ fs.StatFS     = (*SystemFS)(nil)
	_ fs.ReadDirFS  = (*SystemFS)(nil)
	_ fs.ReadFileFS = (*SystemFS)(nil)
	_ fs.GlobFS     = (*SystemFS)(nil)
	_ fs.SubFS      = (*SystemFS)(nil)
)
