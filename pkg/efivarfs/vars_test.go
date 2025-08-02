// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package efivarfs

import (
	"bytes"
	"errors"
	"os"
	"testing"

	guid "github.com/google/uuid"
)

type fake struct {
	err error
}

func (f *fake) Get(desc VariableDescriptor) (VariableAttributes, []byte, error) {
	return VariableAttributes(0), make([]byte, 32), f.err
}

func (f *fake) Set(desc VariableDescriptor, attrs VariableAttributes, data []byte) error {
	return f.err
}

func (f *fake) Remove(desc VariableDescriptor) error {
	return f.err
}

var fakeGUID = guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")

func (f *fake) List() ([]VariableDescriptor, error) {
	return []VariableDescriptor{
		{Name: "fake", GUID: fakeGUID},
	}, f.err
}

var _ EFIVar = &fake{}

func TestReadVariableErrNoFS(t *testing.T) {
	if _, err := NewPath("/tmp"); !errors.Is(err, ErrNoFS) {
		t.Fatalf(`NewPath("/tmp"): %s != %v`, err, ErrNoFS)
	}
}

func TestSimpleReadVariable(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		err    error
		efivar EFIVar
	}{
		{
			name:   "bad variable no -",
			val:    "xy",
			err:    ErrBadGUID,
			efivar: &fake{},
		},
		{
			name:   "bad variable",
			val:    "xy-b-c",
			err:    ErrBadGUID,
			efivar: &fake{},
		},
		{
			name:   "good variable, bad get",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    os.ErrPermission,
			efivar: &fake{err: os.ErrPermission},
		},
		{
			name:   "good variable, good get",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    nil,
			efivar: &fake{},
		},
	}

	for _, tt := range tests {
		_, _, err := SimpleReadVariable(tt.efivar, tt.val)
		if !errors.Is(err, tt.err) {
			t.Errorf("SimpleReadVariable(tt.efivar, %s): %v != %v", tt.val, err, tt.err)
		}
	}
}

func TestSimpleWriteVariable(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		err    error
		efivar EFIVar
	}{
		{
			name:   "bad variable",
			val:    "xy-b-c",
			err:    ErrBadGUID,
			efivar: &fake{},
		},
		{
			name:   "good variable, bad set",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    os.ErrPermission,
			efivar: &fake{err: os.ErrPermission},
		},
		{
			name:   "good variable, good set",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    nil,
			efivar: &fake{},
		},
	}

	for _, tt := range tests {
		err := SimpleWriteVariable(tt.efivar, tt.val, VariableAttributes(0), &bytes.Buffer{})
		if !errors.Is(err, tt.err) {
			t.Errorf("SimpleWriteVariable(tt.efivar, %s): %v != %v", tt.val, err, tt.err)
		}
	}
}

func TestSimpleRemoveVariable(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		err    error
		efivar EFIVar
	}{
		{
			name:   "bad variable",
			val:    "xy-b-c",
			err:    ErrBadGUID,
			efivar: &fake{},
		},
		{
			name:   "good variable, bad Remove",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    os.ErrPermission,
			efivar: &fake{err: os.ErrPermission},
		},
		{
			name:   "good variable, good remove",
			val:    "WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			err:    nil,
			efivar: &fake{},
		},
	}

	for _, tt := range tests {
		err := SimpleRemoveVariable(tt.efivar, tt.val)
		if !errors.Is(err, tt.err) {
			t.Errorf("SimpleRemoveVariable(tt.efivar, %s): %v != %v", tt.val, err, tt.err)
		}
	}
}

func TestSimpleListVariable(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		efivar EFIVar
	}{
		{
			name:   "bad List",
			err:    os.ErrPermission,
			efivar: &fake{err: os.ErrPermission},
		},
		{
			name:   "good List",
			err:    nil,
			efivar: &fake{},
		},
	}

	for _, tt := range tests {
		_, err := SimpleListVariables(tt.efivar)
		if !errors.Is(err, tt.err) {
			t.Errorf("SimpleListVariable(tt.efivar): %v != %v", err, tt.err)
		}
	}
}
