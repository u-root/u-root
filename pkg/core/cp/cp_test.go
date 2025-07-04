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

var testdata = []byte("This is a test string")

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
			opt:     Options{},
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
			if err := tt.opt.Copy(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("%q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		// After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.CopyTree(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		// After every test with NoFollowSymlink we have to delete the created symlink
		if strings.Contains(tt.dstfile, "symlink") {
			os.Remove(tt.dstfile)
		}
	}
}
