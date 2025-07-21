// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package efivarfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	guid "github.com/google/uuid"
)

func TestProbeAndReturn(t *testing.T) {
	for _, tt := range []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "wrong magic",
			path:    "/tmp",
			wantErr: ErrNoFS,
		},
		{
			name:    "wrong directory",
			path:    "/bogus",
			wantErr: ErrNoFS,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewPath(tt.path); !errors.Is(err, tt.wantErr) {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	for _, tt := range []struct {
		name    string
		vd      VariableDescriptor
		attr    VariableAttributes
		data    []byte
		setup   func(path string, t *testing.T)
		wantErr error
	}{
		{
			name: "get var",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr: AttributeNonVolatile,
			data: []byte("testdata"),
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				var buf bytes.Buffer
				if err := binary.Write(&buf, binary.LittleEndian, AttributeNonVolatile); err != nil {
					t.Errorf("Failed writing bytes: %v", err)
				}
				if _, err := buf.Write([]byte("testdata")); err != nil {
					t.Errorf("Failed writing data: %v", err)
				}
				if _, err := buf.WriteTo(f); err != nil {
					t.Errorf("Failed writing to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: nil,
		},
		{
			name: "not exist",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr:    0,
			data:    nil,
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: ErrVarNotExist,
		},
		/* TODO: this test seems utterly broken. I have no idea why it ever seemed it might work.
		{
			name: "no permission",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() *guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return &g
				}(),
			},
			attr: 0,
			data: nil,
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				if err := f.Chmod(0222); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: ErrVarPermission,
		},
		*/
		{
			name: "var empty",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr: 0,
			data: nil,
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t).Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: ErrVarNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			tt.setup(tmp, t)
			// This setup bypasses all the tests for this fake varfs.
			e := &EFIVarFS{path: tmp}

			attr, data, err := e.Get(tt.vd)
			if errors.Is(err, ErrNoFS) {
				t.Logf("no EFIVarFS: %v; skipping this test", err)
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected: %q, got: %v", tt.wantErr, err)
			}
			if attr != tt.attr {
				t.Errorf("Want %v, Got: %v", tt.attr, attr)
			}
			if !bytes.Equal(data, tt.data) {
				t.Errorf("Want: %v, Got: %v", tt.data, data)
			}
		})
	}
}

func TestSet(t *testing.T) {
	for _, tt := range []struct {
		name    string
		vd      VariableDescriptor
		attr    VariableAttributes
		data    []byte
		setup   func(path string, t *testing.T)
		wantErr error
	}{
		{
			name: "set var",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr:    0,
			data:    nil,
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: nil,
		},
		{
			name: "append write",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr:    AttributeAppendWrite,
			data:    nil,
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: nil,
		},
		{
			name: "no read permission",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr: 0,
			data: nil,
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				if err := f.Chmod(0o222); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: ErrVarPermission,
		},
		{
			name: "var exists",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr: 0,
			data: nil,
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: errors.New("inappropriate ioctl for device"),
		},
		{
			name: "input data",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			attr:    0,
			data:    []byte("tests"),
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			tt.setup(tmp, t)
			// This setup bypasses all the tests for this fake varfs.
			e := &EFIVarFS{path: tmp}

			if err := e.Set(tt.vd, tt.attr, tt.data); err != nil {
				if !errors.Is(err, tt.wantErr) {
					// Needed as some errors include changing tmp directory names
					if !strings.Contains(err.Error(), tt.wantErr.Error()) {
						t.Errorf("Want: %v, Got: %v", tt.wantErr, err)
					}
				}
			}
		})
	}
}

func TestRemove(t *testing.T) {
	for _, tt := range []struct {
		name    string
		vd      VariableDescriptor
		setup   func(path string, t *testing.T)
		wantErr error
	}{
		{
			name: "remove var",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t).Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: errors.New("inappropriate ioctl for device"),
		},
		{
			name: "no permission",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				if err := f.Chmod(0o444); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: ErrVarPermission,
		},
		{
			name: "var not exist",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: ErrVarNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			tt.setup(tmp, t)
			// This setup bypasses all the tests for this fake varfs.
			e := &EFIVarFS{path: tmp}

			if err := e.Remove(tt.vd); err != nil {
				if !errors.Is(err, tt.wantErr) {
					// Needed as some errors include changing tmp directory names
					if !strings.Contains(err.Error(), tt.wantErr.Error()) {
						t.Errorf("Want: %v, Got: %v", tt.wantErr, err)
					}
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	for _, tt := range []struct {
		name    string
		vd      VariableDescriptor
		dir     string
		setup   func(path string, t *testing.T)
		wantErr error
	}{
		{
			name: "empty var",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir: t.TempDir(),
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t).Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: nil,
		},
		{
			name: "var with data",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir: t.TempDir(),
			setup: func(path string, t *testing.T) {
				t.Helper()
				f := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350", t)
				var buf bytes.Buffer
				if err := binary.Write(&buf, binary.LittleEndian, AttributeNonVolatile); err != nil {
					t.Errorf("Failed writing bytes: %v", err)
				}
				if _, err := buf.Write([]byte("testdata")); err != nil {
					t.Errorf("Failed writing data: %v", err)
				}
				if _, err := buf.WriteTo(f); err != nil {
					t.Errorf("Failed writing to file: %v", err)
				}
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: nil,
		},
		{
			name: "no regular files",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir: t.TempDir(),
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(path, "testdir"), 0o644); err != nil {
					t.Errorf("Failed to create directory: %v", err)
				}
			},
			wantErr: nil,
		},
		{
			name: "no permission",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir: t.TempDir(),
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := os.Chmod(path, 0o222); err != nil {
					t.Errorf("Failed changing permissions: %v", err)
				}
			},
			wantErr: ErrVarPermission,
		},
		{
			name: "no dir",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir:     "/bogus",
			setup:   func(path string, t *testing.T) { t.Helper() },
			wantErr: ErrVarNotExist,
		},
		{
			name: "malformed vars",
			vd: VariableDescriptor{
				Name: "TestVar",
				GUID: func() guid.UUID {
					g := guid.MustParse("bc54d3fb-ed45-462d-9df8-b9f736228350")
					return g
				}(),
			},
			dir: t.TempDir(),
			setup: func(path string, t *testing.T) {
				t.Helper()
				if err := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b9f7362283500000", t).Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				if err := createTestVar(path, "TestVar-bc54d3fb-ed45-462d-9df8-b", t).Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
			},
			wantErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			tt.setup(tt.dir, t)
			// This setup bypasses all the tests for this fake varfs.
			e := &EFIVarFS{path: tmp}

			if _, err := e.List(); err != nil {
				if !errors.Is(err, tt.wantErr) {
					// Needed as some errors include changing tmp directory names
					if !strings.Contains(err.Error(), tt.wantErr.Error()) {
						t.Errorf("Want: %v, Got: %v", tt.wantErr, err)
					}
				}
			}
		})
	}
}

func createTestVar(path, varFullName string, t *testing.T) *os.File {
	t.Helper()
	f, err := os.Create(filepath.Join(path, varFullName))
	if err != nil {
		t.Errorf("Failed creating test var: %v", err)
	}
	return f
}

func TestNew(t *testing.T) {
	// the EFI file system may not be available, but we call New
	// anyway to at least get some coverage.
	e, err := New()
	t.Logf("New(): %v, %v", e, err)
}
