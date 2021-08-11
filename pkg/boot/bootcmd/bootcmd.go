// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bootcmd handles common cleanup functions and flags that all boot
// commands should support.
package bootcmd

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/mount"
)

// ShowMenuAndBoot handles common cleanup functions and flags that all boot
// commands should support.
//
// mountPool is unmounted before kexecing. noLoad prints the list of entries
// and exits. If noLoad is false, a boot menu is shown to the user. The
// user-chosen boot entry will be kexec'd unless noExec is true.
func ShowMenuAndBoot(entries []menu.Entry, mountPool *mount.Pool, noLoad, noExec bool) {
	if noLoad {
		log.Print("Not loading menu or kernel. Options:")
		for i, entry := range entries {
			log.Printf("%d. %s", i+1, entry.Label())
			log.Printf("=> %s", entry)
		}
		os.Exit(0)
	}

	loadedEntry := menu.ShowMenuAndLoad(true, entries...)

	// Clean up.
	if mountPool != nil {
		if err := mountPool.UnmountAll(mount.MNT_DETACH); err != nil {
			log.Printf("Failed in UnmountAll: %v", err)
		}
	}
	if loadedEntry == nil {
		log.Fatalf("Nothing to boot.")
	}
	if noExec {
		log.Printf("Chosen menu entry: %s", loadedEntry)
		os.Exit(0)
	}
	// Exec should either return an error or not return at all.
	if err := loadedEntry.Exec(); err != nil {
		log.Fatalf("Failed to exec %s: %v", loadedEntry, err)
	}

	// Kexec should either return an error or not return.
	log.Fatalf("Kexec should have returned an error or not returned at all.")
}
