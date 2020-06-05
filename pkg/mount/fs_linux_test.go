// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetFileSystems(t *testing.T) {
	fstypes, err := internalGetFilesystems("testdata/filesystems")
	expected := []string{"ext4", "ext3", "vfat"}
	if err != nil {
		t.Errorf("InternalGetFilesystems failed with error %v", err)
	}
	if !reflect.DeepEqual(fstypes, expected) {
		t.Errorf("Expected '%q', but resulted with '%q'", expected, fstypes)
	}
}

func TestEmptyFileSystems(t *testing.T) {
	fstypes, err := internalGetFilesystems("testdata/emptyFile")
	if err != nil {
		t.Errorf("InternalGetFilesystems failed with error %v", err)
	}
	if len(fstypes) != 0 {
		t.Error("Expected no results for empty filesystem file.")
	}
}

func TestFindFileSystem(t *testing.T) {
	procfilesystems := `nodev   sysfs
nodev   rootfs
nodev   ramfs
        vfat
        btrfs
        ext3
        ext2
        ext4
`

	for _, tt := range []struct {
		name string
		err  string
	}{
		{"rootfs", "<nil>"},
		{"ext3", "<nil>"},
		{"bogusfs", "file system type \"bogusfs\" not found"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := internalFindFileSystem(procfilesystems, tt.name)
			// There has to be a better way to do this.
			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("%s: got %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}
