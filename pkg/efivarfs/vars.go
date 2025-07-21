// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package efivarfs allows interaction with efivarfs of the
// linux kernel.
package efivarfs

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	guid "github.com/google/uuid"
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
	GUID guid.UUID
}

// ErrBadGUID is for any errors parsing GUIDs.
var ErrBadGUID = errors.New("bad GUID")

func guidParse(v string) ([]string, *guid.UUID, error) {
	vs := strings.SplitN(v, "-", 2)
	if len(vs) < 2 {
		return nil, nil, fmt.Errorf("GUID must have name-GUID format: %w", ErrBadGUID)
	}
	g, err := guid.Parse(vs[1])
	if err != nil {
		return nil, nil, fmt.Errorf("%w:%w", ErrBadGUID, err)
	}
	return vs, &g, nil
}

// ReadVariable calls get() on the current efivarfs backend.
func ReadVariable(e EFIVar, desc VariableDescriptor) (VariableAttributes, []byte, error) {
	return e.Get(desc)
}

// SimpleReadVariable is like ReadVariable but takes the combined name and guid string
// of the form name-guid and returns a bytes.Reader instead of a []byte.
func SimpleReadVariable(e EFIVar, v string) (VariableAttributes, *bytes.Reader, error) {
	vs, g, err := guidParse(v)
	if err != nil {
		return 0, nil, err
	}
	attrs, data, err := ReadVariable(e,
		VariableDescriptor{
			Name: vs[0],
			GUID: *g,
		},
	)
	return attrs, bytes.NewReader(data), err
}

// WriteVariable calls set() on the current efivarfs backend.
func WriteVariable(e EFIVar, desc VariableDescriptor, attrs VariableAttributes, data []byte) error {
	return e.Set(desc, attrs, data)
}

// SimpleWriteVariable is like WriteVariable but takes the combined name and guid string
// of the form name-guid.
func SimpleWriteVariable(e EFIVar, v string, attrs VariableAttributes, data *bytes.Buffer) error {
	vs, g, err := guidParse(v)
	if err != nil {
		return err
	}
	return WriteVariable(e,
		VariableDescriptor{
			Name: vs[0],
			GUID: *g,
		}, attrs, data.Bytes(),
	)
}

// RemoveVariable calls remove() on the current efivarfs backend.
func RemoveVariable(e EFIVar, desc VariableDescriptor) error {
	return e.Remove(desc)
}

// SimpleRemoveVariable is like RemoveVariable but takes the combined name and guid string
// of the form name-guid.
func SimpleRemoveVariable(e EFIVar, v string) error {
	vs, g, err := guidParse(v)
	if err != nil {
		return err
	}
	return RemoveVariable(e,
		VariableDescriptor{
			Name: vs[0],
			GUID: *g,
		},
	)
}

// ListVariables calls list() on the current efivarfs backend.
func ListVariables(e EFIVar) ([]VariableDescriptor, error) {
	return e.List()
}

// SimpleListVariables is like ListVariables but returns a []string instead of a []VariableDescriptor.
func SimpleListVariables(e EFIVar) ([]string, error) {
	list, err := ListVariables(e)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, v := range list {
		out = append(out, v.Name+"-"+v.GUID.String())
	}
	return out, nil
}
