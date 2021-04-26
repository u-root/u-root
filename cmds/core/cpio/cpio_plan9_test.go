// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func checkFileInfo(t *testing.T, ent *dirEnt, newFileInfo os.FileInfo) {
	t.Helper()

	newStatT := newFileInfo.Sys().(*syscall.Dir)
	statT := ent.FileInfo.Sys().(*syscall.Dir)
	t.Logf("Entry %v statT %v", ent, newFileInfo)
	if ent.FileInfo.Name() != newFileInfo.Name() ||
		ent.FileInfo.Size() != newFileInfo.Size() ||
		ent.FileInfo.Mode() != newFileInfo.Mode() ||
		ent.FileInfo.IsDir() != newFileInfo.IsDir() ||
		statT.Mode != newStatT.Mode ||
		statT.Uid != newStatT.Uid ||
		statT.Gid != newStatT.Gid {
		msg := "File has mismatched attributes:\n"
		msg += "Property |   Original |  Extracted\n"
		msg += fmt.Sprintf("Name:    | %10s | %10s\n", ent.FileInfo.Name(), newFileInfo.Name())
		msg += fmt.Sprintf("Size:    | %10d | %10d\n", ent.FileInfo.Size(), newFileInfo.Size())
		msg += fmt.Sprintf("Mode:    | %10d | %10d\n", ent.FileInfo.Mode(), newFileInfo.Mode())
		msg += fmt.Sprintf("IsDir:   | %10t | %10t\n", ent.FileInfo.IsDir(), newFileInfo.IsDir())
		msg += fmt.Sprintf("FI Mode: | %10d | %10d\n", statT.Mode, newStatT.Mode)
		msg += fmt.Sprintf("Uid:     | %10s | %10s\n", statT.Uid, newStatT.Uid)
		msg += fmt.Sprintf("Gid:     | %10s | %10s\n", statT.Gid, newStatT.Gid)
		t.Error(msg)
	}
}
