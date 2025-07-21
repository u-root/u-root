// Copyright 2014-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/efivarfs"
)

type failingOS struct {
	err error
}

func (f *failingOS) Get(desc efivarfs.VariableDescriptor) (efivarfs.VariableAttributes, []byte, error) {
	return efivarfs.VariableAttributes(0), make([]byte, 32), f.err
}

func (f *failingOS) Set(desc efivarfs.VariableDescriptor, attrs efivarfs.VariableAttributes, data []byte) error {
	return f.err
}

func (f *failingOS) Remove(desc efivarfs.VariableDescriptor) error {
	return f.err
}

func (f *failingOS) List() ([]efivarfs.VariableDescriptor, error) {
	return make([]efivarfs.VariableDescriptor, 3), f.err
}

var _ efivarfs.EFIVar = &failingOS{}

var (
	badfs = &failingOS{err: os.ErrNotExist}
	nofs  = &failingOS{err: efivarfs.ErrNoFS}
	iofs  = &failingOS{err: syscall.EIO}
	okfs  = &failingOS{err: nil}
)

// We should not test the actual /sys varfs itself. That is done in the package.
// So it suffices to test the login in run() with a faked up EFIVarFS that points to /tmp.
func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name     string
		e        efivarfs.EFIVar
		setup    func(path string, t *testing.T) string
		list     bool
		read     string
		delete   string
		write    string
		wantErr  error
		needRoot bool
	}{
		{
			name: "list no efivarfs",
			e:    nofs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			list:    true,
			wantErr: efivarfs.ErrNoFS,
		},
		{
			name: "read no efivarfs",
			e:    badfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			read:    "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			wantErr: os.ErrNotExist,
		},
		{
			name: "read bad variable",
			e:    badfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			read:    "TestVar",
			wantErr: efivarfs.ErrBadGUID,
		},
		{
			name: "read good variable",
			e:    okfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			read:    " WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			wantErr: nil,
		},
		{
			name: "delete no efivarfs",
			e:    badfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			delete:  "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			wantErr: os.ErrNotExist,
		},
		{
			name: "write malformed var",
			e:    badfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			write:   "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350000",
			wantErr: os.ErrInvalid,
		},
		{
			name: "write no content",
			e:    badfs,
			setup: func(path string, t *testing.T) string {
				t.Helper()
				// oh fun this is what actually sets content.
				return "/bogus"
			},
			write:   "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			wantErr: os.ErrNotExist,
		},
		{
			name:  "write no guid no efivarfs",
			e:     iofs,
			write: "TestVar",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				f, err := os.Create(filepath.Join(path, "content"))
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				s := f.Name()
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				return s
			},
			wantErr: syscall.EIO,
		},
		{
			name:  "write good variable bad content",
			e:     okfs,
			write: " WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return filepath.Join(path, "content")
			},
			wantErr: os.ErrNotExist,
		},
		{
			name:  "write good variable succeeds",
			e:     okfs,
			write: " WriteOnceStatus-4b3082a3-80c6-4d7e-9cd0-583917265df1",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				f, err := os.Create(filepath.Join(path, "content"))
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				s := f.Name()
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				return s
			},
			wantErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(io.Discard, tt.e, tt.list, tt.read, tt.delete, tt.write, tt.setup(t.TempDir(), t)); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Got: %q, Want: %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestBadRunPath(t *testing.T) {
	if err := runpath(os.Stdout, "/tmp", false, "", "", "", ""); !errors.Is(err, efivarfs.ErrNoFS) {
		t.Errorf(`runpath(os.Stdout, "/tmp", false, "", "", "", "", ""): %v != %v`, err, efivarfs.ErrNoFS)
	}
}

func TestGoodRunPath(t *testing.T) {
	if _, err := os.Stat(efivarfs.DefaultVarFS); err != nil {
		t.Skipf("%q: %v, skipping test", efivarfs.DefaultVarFS, err)
	}

	if err := runpath(os.Stdout, efivarfs.DefaultVarFS, false, "", "", "", ""); err != nil {
		t.Errorf(`runpath(os.Stdout, %q, false, "", "", "", "", ""): %v != %v`, efivarfs.DefaultVarFS, err, efivarfs.ErrNoFS)
	}
}
