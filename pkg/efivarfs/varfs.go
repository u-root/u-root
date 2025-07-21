// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package efivarfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	guid "github.com/google/uuid"
	"golang.org/x/sys/unix"
)

// DefaultVarFS is the path to the efivarfs mount point
const DefaultVarFS = "/sys/firmware/efi/efivars/"

var (
	// ErrVarsUnavailable is caused by not having a valid backend
	ErrVarsUnavailable = fmt.Errorf("no variable backend is available:%w", os.ErrNotExist)

	// ErrVarNotExist is caused by accessing a non-existing variable
	ErrVarNotExist = os.ErrNotExist

	// ErrVarPermission is caused by not haven the right permissions either
	// because of not being root or xattrs not allowing changes
	ErrVarPermission = os.ErrPermission

	// ErrNoFS is returned when the file system is not available for some
	// reason.
	ErrNoFS = errors.New("varfs not available")
)

// EFIVar is the interface for EFI variables. Note that it need not use a file system,
// but typically does.
type EFIVar interface {
	Get(desc VariableDescriptor) (VariableAttributes, []byte, error)
	List() ([]VariableDescriptor, error)
	Remove(desc VariableDescriptor) error
	Set(desc VariableDescriptor, attrs VariableAttributes, data []byte) error
}

// EFIVarFS implements EFIVar
type EFIVarFS struct {
	path string
}

var _ EFIVar = &EFIVarFS{}

// NewPath returns an EFIVarFS given a path.
func NewPath(p string) (*EFIVarFS, error) {
	e := &EFIVarFS{path: p}
	if err := e.probe(); err != nil {
		return nil, err
	}
	return e, nil
}

// New returns an EFIVarFS using the default path.
func New() (*EFIVarFS, error) {
	return NewPath(DefaultVarFS)
}

// probe will probe for the EFIVarFS filesystem
// magic value on the expected mountpoint inside the sysfs.
// If the correct magic value was found it will return
// the a pointer to an EFIVarFS struct on which regular
// operations can be done. Otherwise it will return an
// error of type ErrFsNotMounted.
func (v *EFIVarFS) probe() error {
	var stat unix.Statfs_t
	if err := unix.Statfs(v.path, &stat); err != nil {
		return fmt.Errorf("%w: not mounted", ErrNoFS)
	}
	if uint(stat.Type) != uint(unix.EFIVARFS_MAGIC) {
		return fmt.Errorf("%w: wrong magic", ErrNoFS)
	}
	return nil
}

// Get reads the contents of an efivar if it exists and has the necessary permission
func (v *EFIVarFS) Get(desc VariableDescriptor) (VariableAttributes, []byte, error) {
	path := filepath.Join(v.path, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	switch {
	case os.IsNotExist(err):
		return 0, nil, ErrVarNotExist
	case os.IsPermission(err):
		return 0, nil, ErrVarPermission
	case err != nil:
		return 0, nil, err
	}
	defer f.Close()

	var attrs VariableAttributes
	if err := binary.Read(f, binary.LittleEndian, &attrs); err != nil {
		if err == io.EOF {
			return 0, nil, ErrVarNotExist
		}
		return 0, nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return 0, nil, err
	}
	return attrs, data, nil
}

// Set modifies a given efivar with the provided contents
func (v *EFIVarFS) Set(desc VariableDescriptor, attrs VariableAttributes, data []byte) error {
	path := filepath.Join(v.path, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
	flags := os.O_WRONLY | os.O_CREATE
	if attrs&AttributeAppendWrite != 0 {
		flags |= os.O_APPEND
	}

	read, err := os.OpenFile(path, os.O_RDONLY, 0)
	switch {
	case os.IsNotExist(err):
	case os.IsPermission(err):
		return ErrVarPermission
	case err != nil:
		return err
	default:
		defer read.Close()

		restoreImmutable, err := makeMutable(read)
		switch {
		case os.IsPermission(err):
			return ErrVarPermission
		case err != nil:
			return err
		}
		defer restoreImmutable()
	}

	write, err := os.OpenFile(path, flags, 0o644)
	switch {
	case os.IsNotExist(err):
		return ErrVarNotExist
	case os.IsPermission(err):
		return ErrVarPermission
	case err != nil:
		return err
	}
	defer write.Close()

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, attrs); err != nil {
		return err
	}
	if _, err := buf.Write(data); err != nil {
		return err
	}
	if _, err := buf.WriteTo(write); err != nil {
		return err
	}
	return nil
}

// Remove makes the specified EFI var mutable and then deletes it
func (v *EFIVarFS) Remove(desc VariableDescriptor) error {
	path := filepath.Join(v.path, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	switch {
	case os.IsNotExist(err):
		return ErrVarNotExist
	case os.IsPermission(err):
		return ErrVarPermission
	case err != nil:
		return err
	default:
		_, err := makeMutable(f)
		switch {
		case os.IsPermission(err):
			return ErrVarPermission
		case err != nil:
			return err
		default:
			f.Close()
		}
	}
	return os.Remove(path)
}

// List returns the VariableDescriptor for each efivar in the system
// TODO: why can't list implement
func (v *EFIVarFS) List() ([]VariableDescriptor, error) {
	const guidLength = 36
	f, err := os.OpenFile(v.path, os.O_RDONLY, 0)
	switch {
	case os.IsNotExist(err):
		return nil, ErrVarNotExist
	case os.IsPermission(err):
		return nil, ErrVarPermission
	case err != nil:
		return nil, err
	}
	defer f.Close()

	dirents, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var entries []VariableDescriptor
	for _, dirent := range dirents {
		if !dirent.Mode().IsRegular() {
			// Skip non-regular files
			continue
		}
		if len(dirent.Name()) < guidLength+1 {
			// Skip files with a basename that isn't long enough
			// to contain a GUID and a hyphen
			continue
		}
		if dirent.Name()[len(dirent.Name())-guidLength-1] != '-' {
			// Skip files where the basename doesn't contain a
			// hyphen between the name and GUID
			continue
		}
		if dirent.Size() == 0 {
			// Skip files with zero size. These are variables that
			// have been deleted by writing an empty payload
			continue
		}

		name := dirent.Name()[:len(dirent.Name())-guidLength-1]
		guid, err := guid.Parse(dirent.Name()[len(name)+1:])
		if err != nil {
			continue
		}

		entries = append(entries, VariableDescriptor{Name: name, GUID: guid})
	}

	sort.Slice(entries, func(i, j int) bool {
		return fmt.Sprintf("%s-%v", entries[i].Name, entries[i].GUID) < fmt.Sprintf("%s-%v", entries[j].Name, entries[j].GUID)
	})
	return entries, nil
}
