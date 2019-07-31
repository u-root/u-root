// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package storage

import (
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
