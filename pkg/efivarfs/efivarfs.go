// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package efivarfs is a wrapper around go-efilib and allows
// interaction with the Linux efivarfs

package efivarfs

import (
	"bytes"
	"os"

	efi "github.com/canonical/go-efilib"
	"golang.org/x/sys/unix"
)

// Attributes corresponds to efi.VariableAttributes
type Attributes uint32

const (
	AttributeNonVolatile                       Attributes = Attributes(efi.AttributeNonVolatile)
	AttributeBootserviceAccess                 Attributes = Attributes(efi.AttributeBootserviceAccess)
	AttributeRuntimeAccess                     Attributes = Attributes(efi.AttributeRuntimeAccess)
	AttributeHardwareErrorRecord               Attributes = Attributes(efi.AttributeHardwareErrorRecord)
	AttributeAuthenticatedWriteAccess          Attributes = Attributes(efi.AttributeAuthenticatedWriteAccess)
	AttributeTimeBasedAuthenticatedWriteAccess Attributes = Attributes(efi.AttributeTimeBasedAuthenticatedWriteAccess)
	AttributeAppendWrite                       Attributes = Attributes(efi.AttributeAppendWrite)
	AttributeEnhancedAuthenticatedAccess       Attributes = Attributes(efi.AttributeEnhancedAuthenticatedAccess)
)

// ListVars returns a string array of all EFI variables
func ListVars() ([]string, error) {
	vars, err := efi.ListVariables()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, v := range vars {
		out = append(out, v.Name+"-"+v.GUID.String())
	}
	return out, nil
}

// WriteVar writes to the specified EFI var using the specified data and creates it if it doesn't exist
func WriteVar(name string, guidString string, attrs Attributes, data bytes.Buffer) error {
	guid, err := efi.DecodeGUIDString(guidString)
	if err != nil {
		return err
	}
	return efi.WriteVariable(name, guid, efi.VariableAttributes(attrs), data.Bytes())
}

// ReadVar return the data and Attributes from the specified EFI var
func ReadVar(name string, guidString string) (bytes.Reader, Attributes, error) {
	guid, err := efi.DecodeGUIDString(guidString)
	if err != nil {
		return *bytes.NewReader(nil), 0, err
	}
	out, attrs, err := efi.ReadVariable(name, guid)
	return *bytes.NewReader(out), Attributes(attrs), err
}

// DeleteVar deletes the specified EFI var
func DeleteVar(name string, guidString string) error {
	path := name + "-" + guidString
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = unix.IoctlSetPointerInt(int(f.Fd()), unix.FS_IOC_SETFLAGS, 0)
	if err != nil {
		return &os.PathError{Op: "ioctl", Path: f.Name(), Err: err}
	}
	f.Close()
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
