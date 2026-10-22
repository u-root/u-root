// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

func checkFileInfo(t *testing.T, ent *dirEnt, newFileInfo os.FileInfo) {
	t.Helper()

	newStatT := newFileInfo.Sys().(*syscall.Stat_t)
	statT := ent.FileInfo.Sys().(*syscall.Stat_t)
	// Directory sizes are filesystem-specific and are not portable archive content.
	sizeMatches := ent.FileInfo.IsDir() || ent.FileInfo.Size() == newFileInfo.Size()
	t.Logf("Entry %v statT %v old ino %d new Ino %d", ent, newFileInfo, statT.Ino, newStatT.Ino)
	if ent.FileInfo.Name() != newFileInfo.Name() ||
		!sizeMatches ||
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

func TestHardLinkEscape(t *testing.T) {
	for _, tt := range []struct {
		name    string
		records []cpio.Record
		// mustNotExist is relative to the parent of the extraction directory.
		mustNotExist string
		// notVictim, if set, must not be a hard link to the victim.
		notVictim string
	}{
		{
			// A link whose new name is outside the extraction directory.
			name: "newname",
			records: []cpio.Record{
				{ReaderAt: strings.NewReader("evil"), Info: cpio.Info{Name: "payload", Ino: 1, Mode: cpio.S_IFREG | 0o644, FileSize: 4}},
				{ReaderAt: strings.NewReader(""), Info: cpio.Info{Name: "../escaped", Ino: 1, Mode: cpio.S_IFREG | 0o644}},
			},
			mustNotExist: "escaped",
		},
		{
			// A link to a file outside the extraction directory.
			name: "oldname",
			records: []cpio.Record{
				{ReaderAt: strings.NewReader("x"), Info: cpio.Info{Name: "../victim", Ino: 1, Mode: cpio.S_IFREG | 0o644, FileSize: 1}},
				{ReaderAt: strings.NewReader(""), Info: cpio.Info{Name: "leak", Ino: 1, Mode: cpio.S_IFREG | 0o644}},
			},
			notVictim: "extract/leak",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			parent := t.TempDir()
			extractDir := filepath.Join(parent, "extract")
			if err := os.Mkdir(extractDir, 0o755); err != nil {
				t.Fatal(err)
			}
			victim := filepath.Join(parent, "victim")
			if err := os.WriteFile(victim, []byte("secret"), 0o600); err != nil {
				t.Fatal(err)
			}

			archive := &bytes.Buffer{}
			rw := cpio.Newc.Writer(archive)
			if err := cpio.WriteRecords(rw, tt.records); err != nil {
				t.Fatal(err)
			}
			if err := cpio.WriteTrailer(rw); err != nil {
				t.Fatal(err)
			}
			archiveFile := filepath.Join(parent, "archive.cpio")
			if err := os.WriteFile(archiveFile, archive.Bytes(), 0o644); err != nil {
				t.Fatal(err)
			}
			f, err := os.Open(archiveFile)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(wd)
			if err := os.Chdir(extractDir); err != nil {
				t.Fatal(err)
			}

			if err := run([]string{"i"}, f, &bytes.Buffer{}, false, "newc"); err == nil {
				t.Errorf("run(i) = nil, want error for escaping hard link")
			}

			if tt.mustNotExist != "" {
				if _, err := os.Lstat(filepath.Join(parent, tt.mustNotExist)); err == nil {
					t.Errorf("hard link %q was created outside the extraction directory", tt.mustNotExist)
				}
			}
			if tt.notVictim != "" {
				fi, err := os.Stat(filepath.Join(parent, tt.notVictim))
				vi, verr := os.Stat(victim)
				if err == nil && verr == nil && os.SameFile(fi, vi) {
					t.Errorf("%q is a hard link to a file outside the extraction directory", tt.notVictim)
				}
			}
		})
	}
}
