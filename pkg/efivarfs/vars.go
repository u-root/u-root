// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package efivarfs allows interaction with efivarfs of the
// linux kernel.
package efivarfs

import (
	"bytes"
	"os"
	"strings"

	guid "github.com/google/uuid"
	"golang.org/x/sys/unix"
)

// VariableAttributes is an uint32 identifying the variables attributes.
type VariableAttributes uint32

const (
	// AttributeNonVolatile indicates a variable is non volatile.
	AttributeNonVolatile VariableAttributes = 0x00000001
	// AttributeBootserviceAccess indicates a variable is accessible during boot service.
	AttributeBootserviceAccess VariableAttributes = 0x00000002
	// AttributeRuntimeAccess indicates a variable is accessible during runtime.
	AttributeRuntimeAccess VariableAttributes = 0x00000004
	// AttributeHardwareErrorRecord indicates a variable holds hardware error records.
	AttributeHardwareErrorRecord VariableAttributes = 0x00000008
	// AttributeAuthenticatedWriteAccess indicates a variable needs authentication before write access.
	AttributeAuthenticatedWriteAccess VariableAttributes = 0x00000010
	// AttributeTimeBasedAuthenticatedWriteAccess indicates a variable needs time based authentication before write access.
	AttributeTimeBasedAuthenticatedWriteAccess VariableAttributes = 0x00000020
	// AttributeAppendWrite indicates data written to this variable is appended.
	AttributeAppendWrite VariableAttributes = 0x00000040
	// AttributeEnhancedAuthenticatedAccess indicate a variable uses the new authentication format.
	AttributeEnhancedAuthenticatedAccess VariableAttributes = 0x00000080
)

// VariableDescriptor contains the name and GUID identifying a variable
type VariableDescriptor struct {
	Name string
	GUID *guid.UUID
}

// ReadVariable calls get() on the current efivarfs backend.
func ReadVariable(desc VariableDescriptor) (VariableAttributes, []byte, error) {
	e, err := probeAndReturn()
	if err != nil {
		return 0, nil, err
	}
	return e.get(desc)
}

// SimpleReadVariable is like ReadVariables but takes the combined name and guid string
// of the form name-guid and returns a bytes.Reader instead of a []byte.
func SimpleReadVariable(v string) (VariableAttributes, *bytes.Reader, error) {
	e, err := probeAndReturn()
	if err != nil {
		return 0, nil, err
	}
	vs := strings.SplitN(v, "-", 2)
	g, err := guid.Parse(vs[1])
	if err != nil {
		return 0, nil, err
	}
	attrs, data, err := e.get(
		VariableDescriptor{
			Name: vs[0],
			GUID: &g,
		},
	)
	return attrs, bytes.NewReader(data), err
}

// WriteVariable calls set() on the current efivarfs backend.
func WriteVariable(desc VariableDescriptor, attrs VariableAttributes, data []byte) error {
	e, err := probeAndReturn()
	if err != nil {
		return err
	}
	return e.set(desc, attrs, data)
}

// SimpleWriteVariable is like WriteVariables but takes the combined name and guid string
// of the form name-guid and returns a bytes.Buffer instead of a []byte.
func SimpleWriteVariable(v string, attrs VariableAttributes, data *bytes.Buffer) error {
	e, err := probeAndReturn()
	if err != nil {
		return err
	}
	vs := strings.SplitN(v, "-", 2)
	g, err := guid.Parse(vs[1])
	if err != nil {
		return err
	}
	return e.set(
		VariableDescriptor{
			Name: vs[0],
			GUID: &g,
		}, attrs, data.Bytes(),
	)
}

// RemoveVariable calls remove() on the current efivarfs backend.
func RemoveVariable(desc VariableDescriptor) error {
	e, err := probeAndReturn()
	if err != nil {
		return err
	}
	return e.remove(desc)
}

// SimpleRemoveVariable is like RemoveVariable but takes the combined name and guid string
// of the form name-guid.
func SimpleRemoveVariable(v string) error {
	e, err := probeAndReturn()
	if err != nil {
		return err
	}
	vs := strings.SplitN(v, "-", 2)
	g, err := guid.Parse(vs[1])
	if err != nil {
		return err
	}
	return e.remove(
		VariableDescriptor{
			Name: vs[0],
			GUID: &g,
		},
	)
}

// ListVariables calls list() on the current efivarfs backend.
func ListVariables() ([]VariableDescriptor, error) {
	e, err := probeAndReturn()
	if err != nil {
		return nil, err
	}
	return e.list()
}

// SimpleListVariables is like ListVariables but returns a []string instead of a []VariableDescriptor.
func SimpleListVariables() ([]string, error) {
	e, err := probeAndReturn()
	if err != nil {
		return nil, err
	}
	list, err := e.list()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, v := range list {
		out = append(out, v.Name+"-"+v.GUID.String())
	}
	return out, nil
}

// getInodeFlags returns the extended attributes of a file.
func getInodeFlags(f *os.File) (int, error) {
	// If I knew how unix.Getxattr works I'd use that...
	flags, err := unix.IoctlGetInt(int(f.Fd()), unix.FS_IOC_GETFLAGS)
	if err != nil {
		return 0, &os.PathError{Op: "ioctl", Path: f.Name(), Err: err}
	}
	return flags, nil
}

// setInodeFlags sets the extended attributes of a file.
func setInodeFlags(f *os.File, flags int) error {
	// If I knew how unix.Setxattr works I'd use that...
	if err := unix.IoctlSetPointerInt(int(f.Fd()), unix.FS_IOC_SETFLAGS, flags); err != nil {
		return &os.PathError{Op: "ioctl", Path: f.Name(), Err: err}
	}
	return nil
}

// makeMutable will change a files xattrs so that
// the immutable flag is removed and return a restore
// function which can reset the flag for that filee.
func makeMutable(f *os.File) (restore func(), err error) {
	flags, err := getInodeFlags(f)
	if err != nil {
		return nil, err
	}
	if flags&unix.STATX_ATTR_IMMUTABLE == 0 {
		return func() {}, nil
	}

	if err := setInodeFlags(f, flags&^unix.STATX_ATTR_IMMUTABLE); err != nil {
		return nil, err
	}
	return func() {
		if err := setInodeFlags(f, flags); err != nil {
			// If setting the immutable did
			// not work it's alright to do nothing
			// because after a reboot the flag is
			// automatically reapplied
			return
		}
	}, nil
}
