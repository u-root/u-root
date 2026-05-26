// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
package p9dev

import (
	"syscall"
	"testing"

	"github.com/hugelgupf/p9/p9"
)

type mockFile struct {
	p9.File
	attr p9.Attr
}

func (m *mockFile) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	return p9.QID{}, req, m.attr, nil
}

func (m *mockFile) Walk(names []string) ([]p9.QID, p9.File, error) {
	return nil, &mockFile{}, nil
}

func (m *mockFile) Create(name string, flags p9.OpenFlags, permissions p9.FileMode, uid p9.UID, gid p9.GID) (p9.File, p9.QID, uint32, error) {
	return &mockFile{}, p9.QID{}, 0, nil
}

type mockAttacher struct {
	p9.Attacher
}

func (m *mockAttacher) Attach() (p9.File, error) {
	return &mockFile{}, nil
}

func TestGetAttr(t *testing.T) {
	tests := []struct {
		name     string
		mode     p9.FileMode
		wantMode p9.FileMode
	}{
		{
			name:     "regular_file",
			mode:     p9.ModeRegular,
			wantMode: p9.ModeRegular,
		},
		{
			name:     "directory",
			mode:     p9.ModeDirectory,
			wantMode: p9.ModeDirectory,
		},
		{
			name:     "char_device",
			mode:     p9.ModeCharacterDevice,
			wantMode: p9.ModeRegular,
		},
		{
			name:     "block_device",
			mode:     p9.ModeBlockDevice,
			wantMode: p9.ModeRegular,
		},
		{
			name:     "fifo",
			mode:     p9.ModeNamedPipe,
			wantMode: p9.ModeRegular,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mf := &mockFile{attr: p9.Attr{Mode: tt.mode}}
			f := &File{File: mf}
			_, _, attr, err := f.GetAttr(p9.AttrMaskAll)
			if err != nil {
				t.Fatalf("GetAttr failed: %v", err)
			}
			if attr.Mode != tt.wantMode {
				t.Errorf("GetAttr() mode = %v, want %v", attr.Mode, tt.wantMode)
			}
		})
	}
}

func TestMknod(t *testing.T) {
	f := &File{File: &mockFile{}}
	_, err := f.Mknod("foo", p9.ModeRegular, 0, 0, 0, 0)
	if err != syscall.ENOSYS {
		t.Errorf("Mknod() error = %v, want %v", err, syscall.ENOSYS)
	}
}

func TestWalk(t *testing.T) {
	f := &File{File: &mockFile{}}
	_, newFile, err := f.Walk(nil)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
	if _, ok := newFile.(*File); !ok {
		t.Errorf("Walk() returned %T, want *File", newFile)
	}
}

func TestCreate(t *testing.T) {
	f := &File{File: &mockFile{}}
	newFile, _, _, err := f.Create("foo", p9.ReadOnly, p9.ModeRegular, 0, 0)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if _, ok := newFile.(*File); !ok {
		t.Errorf("Create() returned %T, want *File", newFile)
	}
}

func TestWalkGetAttr(t *testing.T) {
	f := &File{File: &mockFile{}}
	_, _, _, _, err := f.WalkGetAttr(nil)
	_, _, _, _, wantErr := f.defaultWalker.WalkGetAttr(nil)
	if err != wantErr {
		t.Errorf("WalkGetAttr() error = %v, want %v", err, wantErr)
	}
}

func TestNew(t *testing.T) {
	a := New(&mockAttacher{})
	f, err := a.Attach()
	if err != nil {
		t.Fatalf("Attach failed: %v", err)
	}
	if _, ok := f.(*File); !ok {
		t.Errorf("Attach() returned %T, want *File", f)
	}
}
