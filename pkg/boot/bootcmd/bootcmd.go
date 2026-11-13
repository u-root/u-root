// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bootcmd handles common cleanup functions and flags that all boot
// commands should support.
package bootcmd

import (
	"errors"
	"fmt"
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
// user-chosen boot entry will be kexec'd unless noExec is true. If the entry
// is a menu.ReturnsToMenu, such as the shell, the menu is shown again when it
// returns.
func ShowMenuAndBoot(entries []menu.Entry, mountPool *mount.Pool, noLoad, noExec bool) {
	if noLoad {
		log.Print("Not loading menu or kernel. Options:")
		for i, entry := range entries {
			log.Printf("%d. %s", i+1, entry.Label())
			log.Printf("=> %s", entry)
		}
		os.Exit(0)
	}

	if err := showMenuAndBoot(menu.ShowMenuAndLoad, entries, mountPool, noExec); err != nil {
		log.Fatal(err)
	}
}

func showMenuAndBoot(load func(bool, ...menu.Entry) menu.Entry, entries []menu.Entry, mountPool *mount.Pool, noExec bool) error {
	for {
		loadedEntry := load(true, entries...)
		_, returns := loadedEntry.(menu.ReturnsToMenu)

		// Keep the file systems the other entries are loaded from if
		// the menu is shown again.
		if !returns && mountPool != nil {
			if err := mountPool.UnmountAll(mount.MNT_DETACH); err != nil {
				log.Printf("Failed in UnmountAll: %v", err)
			}
			mountPool = nil
		}
		if loadedEntry == nil {
			return errors.New("nothing to boot")
		}
		if noExec {
			log.Printf("Chosen menu entry: %s", loadedEntry)
			return nil
		}
		// If kexec fails, something is seriously wrong. Exit, so that
		// the machine is reset.
		if err := loadedEntry.Exec(); err != nil {
			return fmt.Errorf("failed to exec %s: %w", loadedEntry, err)
		}
		if !returns {
			return errors.New("kexec should have returned an error or not returned at all")
		}
	}
}
