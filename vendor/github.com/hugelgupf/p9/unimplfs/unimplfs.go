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

// Package unimplfs provides an p9.File template that returns ENOSYS for all
// methods.
//
// This can be used to leave some methods unimplemented in incomplete p9.File
// implementations.
package unimplfs

import (
	"syscall"

	"github.com/hugelgupf/p9/p9"
)

// NoopFile is a p9.File that returns ENOSYS for every method.
type NoopFile struct {
	p9.DefaultWalkGetAttr
}

var (
	_ p9.File = &NoopFile{}
)

// Walk implements p9.File.Walk.
func (NoopFile) Walk(names []string) ([]p9.QID, p9.File, error) {
	return nil, nil, syscall.ENOSYS
}

// StatFS implements p9.File.StatFS.
//
// Not implemented.
func (NoopFile) StatFS() (p9.FSStat, error) {
	return p9.FSStat{}, syscall.ENOSYS
}

// FSync implements p9.File.FSync.
func (NoopFile) FSync() error {
	return syscall.ENOSYS
}

// GetAttr implements p9.File.GetAttr.
func (NoopFile) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	return p9.QID{}, p9.AttrMask{}, p9.Attr{}, syscall.ENOSYS
}

// SetAttr implements p9.File.SetAttr.
func (NoopFile) SetAttr(valid p9.SetAttrMask, attr p9.SetAttr) error {
	return syscall.ENOSYS
}

// Remove implements p9.File.Remove.
func (NoopFile) Remove() error {
	return syscall.ENOSYS
}

// Rename implements p9.File.Rename.
func (NoopFile) Rename(directory p9.File, name string) error {
	return syscall.ENOSYS
}

// Close implements p9.File.Close.
func (NoopFile) Close() error {
	return nil
}

// Open implements p9.File.Open.
func (NoopFile) Open(mode p9.OpenFlags) (p9.QID, uint32, error) {
	return p9.QID{}, 0, syscall.ENOSYS
}

// Read implements p9.File.Read.
func (NoopFile) ReadAt(p []byte, offset uint64) (int, error) {
	return 0, syscall.ENOSYS
}

// Write implements p9.File.Write.
func (NoopFile) WriteAt(p []byte, offset uint64) (int, error) {
	return 0, syscall.ENOSYS
}

// Create implements p9.File.Create.
func (NoopFile) Create(name string, mode p9.OpenFlags, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.File, p9.QID, uint32, error) {
	return nil, p9.QID{}, 0, syscall.ENOSYS
}

// Mkdir implements p9.File.Mkdir.
func (NoopFile) Mkdir(name string, permissions p9.FileMode, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, syscall.ENOSYS
}

// Symlink implements p9.File.Symlink.
func (NoopFile) Symlink(oldname string, newname string, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, syscall.ENOSYS
}

// Link implements p9.File.Link.
func (NoopFile) Link(target p9.File, newname string) error {
	return syscall.ENOSYS
}

// Mknod implements p9.File.Mknod.
func (NoopFile) Mknod(name string, mode p9.FileMode, major uint32, minor uint32, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, syscall.ENOSYS
}

// RenameAt implements p9.File.RenameAt.
func (NoopFile) RenameAt(oldname string, newdir p9.File, newname string) error {
	return syscall.ENOSYS
}

// UnlinkAt implements p9.File.UnlinkAt.
func (NoopFile) UnlinkAt(name string, flags uint32) error {
	return syscall.ENOSYS
}

// Readdir implements p9.File.Readdir.
func (NoopFile) Readdir(offset uint64, count uint32) ([]p9.Dirent, error) {
	return nil, syscall.ENOSYS
}

// Readlink implements p9.File.Readlink.
func (NoopFile) Readlink() (string, error) {
	return "", syscall.ENOSYS
}

// Flush implements p9.File.Flush.
func (NoopFile) Flush() error {
	return nil
}

// Renamed implements p9.File.Renamed.
func (NoopFile) Renamed(parent p9.File, newName string) {}

// Allocate implements p9.File.Allocate.
func (NoopFile) Allocate(mode p9.AllocateMode, offset, length uint64) error {
	return syscall.ENOSYS
}
