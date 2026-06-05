// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cp

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

var (
	testdata = []byte("This is a test string")
)

func TestCopySimple(t *testing.T) {
	var err error
	tmpdirDst := t.TempDir()
	tmpdirSrc := t.TempDir()

	srcfiles := make([]*os.File, 2)
	dstfiles := make([]*os.File, 2)
	for iterator := range srcfiles {
		srcfiles[iterator], err = os.CreateTemp(tmpdirSrc, "file-to-copy"+fmt.Sprintf("%d", iterator))
		if err != nil {
			t.Errorf("failed to create temp file: %q", err)
		}
		if _, err = srcfiles[iterator].Write(testdata); err != nil {
			t.Errorf("failed to write testdata to file")
		}
		dstfiles[iterator], err = os.CreateTemp(tmpdirDst, "file-to-copy"+fmt.Sprintf("%d", iterator))
		if err != nil {
			t.Errorf("failed to create temp file: %q", err)
		}
	}

	sl := filepath.Join(tmpdirDst, "test-symlink")
	if err := os.Symlink(srcfiles[1].Name(), sl); err != nil {
		t.Errorf("creating symlink failed")
	}

	for _, tt := range []struct {
		name    string
		srcfile string
		dstfile string
		opt     Options
		wantErr error
	}{
		{
			name:    "Success",
			srcfile: srcfiles[0].Name(),
			dstfile: dstfiles[0].Name(),
			opt:     Default,
		},
		{
			name:    "SrcDstDirctoriesSuccess",
			srcfile: tmpdirSrc,
			dstfile: tmpdirDst,
		},
		{
			name:    "SrcNotExist",
			srcfile: "file-does-not-exist",
			dstfile: dstfiles[0].Name(),
			wantErr: fs.ErrNotExist,
		},
		{
			name:    "DstIsDirectory",
			srcfile: srcfiles[0].Name(),
			dstfile: tmpdirDst,
			wantErr: unix.EISDIR,
		},
		{
			name:    "CopySymlink",
			srcfile: sl,
			dstfile: dstfiles[1].Name(),
			opt: Options{
				NoFollowSymlinks: false,
			},
		},
		{
			name:    "CopySymlinkFollow",
			srcfile: sl,
			dstfile: filepath.Join(tmpdirDst, "followed-symlink"),
			opt: Options{
				NoFollowSymlinks: true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.srcfile, tt.dstfile)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		//After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.Copy(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("%q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		//After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.CopyTree(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		//After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := CopyTree(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		//After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}
	}
}

func TestRunMultipleSourcesWithWorkingDirDestination(t *testing.T) {
	workingDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(workingDir, "out"), 0755); err != nil {
		t.Fatalf(`Mkdir("out") = %v, want nil`, err)
	}
	for _, name := range []string{"file1", "file2"} {
		if err := os.WriteFile(filepath.Join(workingDir, name), []byte(name), 0644); err != nil {
			t.Fatalf(`WriteFile(%q) = %v, want nil`, name, err)
		}
	}

	cmd := New()
	cmd.SetWorkingDir(workingDir)
	if err := cmd.Run("file1", "file2", "out"); err != nil {
		t.Fatalf(`Run("file1", "file2", "out") = %v, want nil`, err)
	}

	for _, name := range []string{"file1", "file2"} {
		got, err := os.ReadFile(filepath.Join(workingDir, "out", name))
		if err != nil {
			t.Fatalf(`ReadFile("out/%s") = %v, want nil`, name, err)
		}
		if string(got) != name {
			t.Errorf(`ReadFile("out/%s") = %q, want %q`, name, string(got), name)
		}
	}
}
