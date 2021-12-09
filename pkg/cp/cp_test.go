// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cp

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestCopy(t *testing.T) {
	tmpdirDst, err := os.MkdirTemp("", "dst-directory")
	if err != nil {
		t.Errorf("failed to create tmp directorty: %q", err)
	}
	defer os.RemoveAll(tmpdirDst)

	tmpdirSrc, err := os.MkdirTemp("", "src-directory")
	if err != nil {
		t.Errorf("failed to create tmp src directorty: %q", err)
	}
	defer os.RemoveAll(tmpdirSrc)

	srcfiles := make([]*os.File, 3)
	dstfiles := make([]*os.File, 3)
	for iterator := range srcfiles {
		srcfiles[iterator], err = os.CreateTemp(tmpdirSrc, "file-to-copy"+fmt.Sprintf("%d", iterator))
		if err != nil {
			t.Errorf("failed to create temp file: %q", err)
		}
		dstfiles[iterator], err = os.CreateTemp(tmpdirDst, "file-to-copy"+fmt.Sprintf("%d", iterator))
		if err != nil {
			t.Errorf("failed to create temp file: %q", err)
		}
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
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := Copy(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.Copy(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.CopyTree(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyTree(tt.srcfile, tt.dstfile); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q, Got: %q", tt.name, tt.wantErr, err)
			}
		})
	}
}
