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

// EfiVarFs is the path to the efivarfs mount point
//
// Note: This has to be a var instead of const because of
// our unit tests.
var EfiVarFs = "/sys/firmware/efi/efivars/"

var (
	// ErrFsNotMounted is caused if no vailed efivarfs magic is found
	ErrFsNotMounted = errors.New("no efivarfs magic found, is it mounted?")

	// ErrVarsUnavailable is caused by not having a valid backend
	ErrVarsUnavailable = errors.New("no variable backend is available")

	// ErrVarNotExist is caused by accessing a non-existing variable
	ErrVarNotExist = errors.New("variable does not exist")

	// ErrVarPermission is caused by not haven the right permissions either
	// because of not being root or xattrs not allowing changes
	ErrVarPermission = errors.New("permission denied")
)

// efivarfs represents the real efivarfs of the Linux kernel
// and has the relevant methods like get, set and remove, which
// will operate on the actual efi variables inside the Linux
// efivarfs backend.
type efivarfs struct{}

// probeAndReturn will probe for the efivarfs filesystem
// magic value on the expected mountpoint inside the sysfs.
// If the correct magic value was found it will return
// the a pointer to an efivarfs struct on which regular
// operations can be done. Otherwise it will return an
// error of type ErrFsNotMounted.
func probeAndReturn() (*efivarfs, error) {
	var stat unix.Statfs_t
	if err := unix.Statfs(EfiVarFs, &stat); err != nil {
		return nil, fmt.Errorf("statfs error occured: %w", ErrFsNotMounted)
	}
	if uint(stat.Type) != uint(unix.EFIVARFS_MAGIC) {
		return nil, fmt.Errorf("wrong fs type: %w", ErrFsNotMounted)
	}
	return &efivarfs{}, nil
}

// get reads the contents of an efivar if it exists and has the necessary permission
func (v *efivarfs) get(desc VariableDescriptor) (VariableAttributes, []byte, error) {
	path := filepath.Join(EfiVarFs, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
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

// set modifies a given efivar with the provided contents
func (v *efivarfs) set(desc VariableDescriptor, attrs VariableAttributes, data []byte) error {
	path := filepath.Join(EfiVarFs, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
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

	write, err := os.OpenFile(path, flags, 0644)
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

// remove makes the specified EFI var mutable and then deletes it
func (v *efivarfs) remove(desc VariableDescriptor) error {
	path := filepath.Join(EfiVarFs, fmt.Sprintf("%s-%s", desc.Name, desc.GUID.String()))
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

// list returns the VariableDescriptor for each efivar in the system
func (v *efivarfs) list() ([]VariableDescriptor, error) {
	const guidLength = 36
	f, err := os.OpenFile(EfiVarFs, os.O_RDONLY, 0)
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

		entries = append(entries, VariableDescriptor{Name: name, GUID: &guid})
	}

	sort.Slice(entries, func(i, j int) bool {
		return fmt.Sprintf("%s-%v", entries[i].Name, entries[i].GUID) < fmt.Sprintf("%s-%v", entries[j].Name, entries[j].GUID)
	})
	return entries, nil
}
