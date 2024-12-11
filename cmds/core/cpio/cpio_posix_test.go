// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func checkFileInfo(t *testing.T, ent *dirEnt, newFileInfo os.FileInfo) {
	t.Helper()

	newStatT := newFileInfo.Sys().(*syscall.Stat_t)
	statT := ent.FileInfo.Sys().(*syscall.Stat_t)
	t.Logf("Entry %v statT %v old ino %d new Ino %d", ent, newFileInfo, statT.Ino, newStatT.Ino)
	if ent.FileInfo.Name() != newFileInfo.Name() ||
		ent.FileInfo.Size() != newFileInfo.Size() ||
		ent.FileInfo.Mode() != newFileInfo.Mode() ||
		ent.FileInfo.IsDir() != newFileInfo.IsDir() ||
		statT.Mode != newStatT.Mode ||
		statT.Uid != newStatT.Uid ||
		statT.Gid != newStatT.Gid ||
		statT.Nlink != newStatT.Nlink {
		msg := "File has mismatched attributes:\n"
		msg += "Property |   Original |  Extracted\n"
		msg += fmt.Sprintf("Name:    | %10s | %10s\n", ent.FileInfo.Name(), newFileInfo.Name())
		msg += fmt.Sprintf("Size:    | %10d | %10d\n", ent.FileInfo.Size(), newFileInfo.Size())
		msg += fmt.Sprintf("Mode:    | %10d | %10d\n", ent.FileInfo.Mode(), newFileInfo.Mode())
		msg += fmt.Sprintf("IsDir:   | %10t | %10t\n", ent.FileInfo.IsDir(), newFileInfo.IsDir())
		msg += fmt.Sprintf("FI Mode: | %10d | %10d\n", statT.Mode, newStatT.Mode)
		msg += fmt.Sprintf("Uid:     | %10d | %10d\n", statT.Uid, newStatT.Uid)
		msg += fmt.Sprintf("Gid:     | %10d | %10d\n", statT.Gid, newStatT.Gid)
		msg += fmt.Sprintf("NLink:   | %10d | %10d\n", statT.Nlink, newStatT.Nlink)
		t.Error(msg)
	}
}
