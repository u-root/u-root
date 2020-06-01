// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package menu displays a Terminal UI based text menu to choose boot options
// from.
package menu

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/sh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
)

const (
	initialTimeout    = 10 * time.Second
	subsequentTimeout = 60 * time.Second
)

// Entry is a menu entry.
type Entry interface {
	// Label is the string displayed to the user in the menu.
	Label() string

	// Do is called when the entry is chosen.
	Do() error

	// IsDefault indicates that this action should be run by default if the
	// user didn't make an entry choice.
	IsDefault() bool
}

// Choose presents the user a menu on input to choose an entry from and returns that entry.
func Choose(input *os.File, entries ...Entry) Entry {
	fmt.Println("")
	for i, e := range entries {
		fmt.Printf("%02d. %s\n\n", i+1, e.Label())
	}
	// For some reason this has to be a non-empty empty line, so that the
	// terminal prompt doesn't override it. Or something.
	fmt.Println(" ")

	oldState, err := terminal.MakeRaw(int(input.Fd()))
	if err != nil {
		log.Printf("BUG: Please report: We cannot actually let you choose from menu (MakeRaw failed): %v", err)
		return nil
	}
	defer terminal.Restore(int(input.Fd()), oldState)

	// TODO(chrisko): reduce this timeout a la GRUB. 3 seconds, and hitting
	// any button resets the timeout. We could save 7 seconds here.
	t := time.NewTimer(initialTimeout)

	boot := make(chan Entry, 1)

	go func() {
		// Read exactly one line.
		term := terminal.NewTerminal(input, "Choose a menu option (hit enter to boot the default - 01 is the default option) > ")

		term.AutoCompleteCallback = func(line string, pos int, key rune) (string, int, bool) {
			// We ain't gonna autocomplete, but we'll reset the countdown timer when you press a key.
			t.Reset(subsequentTimeout)
			return "", 0, false
		}

		for {
			choice, err := term.ReadLine()
			if err != nil {
				if err != io.EOF {
					fmt.Printf("BUG: Please report: Terminal read error: %v. ", err)
				}
				boot <- nil
				return
			}

			if choice == "" {
				// nil will result in the default order.
				boot <- nil
				return
			}
			num, err := strconv.Atoi(choice)
			if err != nil {
				fmt.Printf("%s is not a valid entry number: %v. ", choice, err)
				continue
			}
			if num-1 < 0 || num > len(entries) {
				fmt.Printf("%s is not a valid entry number. ", choice)
				continue
			}
			boot <- entries[num-1]
			return
		}
	}()

	select {
	case entry := <-boot:
		if entry != nil {
			fmt.Printf("Chosen option %s.\n\n", entry.Label())
		}
		return entry

	case <-t.C:
		return nil
	}
}

// errStopTestOnly makes ShowMenuAndBoot return if Entry.Do returns it. This
// exists because we expect all menu entries to take over execution context if
// they succeed (e.g. exec a shell, reboot, exec a kernel). Success for Do() is
// only if it never returns.
//
// We can't test that it won't return, so we use this placeholder value instead
// to indicate "it worked".
var errStopTestOnly = errors.New("makes ShowMenuAndBoot return only in tests")

// ShowMenuAndBoot lets the user choose one of entries and boots it.
func ShowMenuAndBoot(input *os.File, entries ...Entry) {
	// Clear the screen (ANSI terminal escape code for screen clear).
	fmt.Printf("\033[1;1H\033[2J\n\n")
	fmt.Printf("Welcome to NERF's Boot Menu\n\n")
	fmt.Printf("Enter a number to boot a kernel:\n")

	for {
		// Allow the user to choose.
		entry := Choose(input, entries...)
		if entry == nil {
			// This only returns something if the user explicitly
			// entered something.
			//
			// If nothing was entered, fall back to default.
			break
		}
		if err := entry.Do(); err != nil {
			log.Printf("Failed to do %s: %v", entry.Label(), err)
		}
	}

	fmt.Println("")

	// We only get one shot at actually booting, so boot the first kernel
	// that can be loaded correctly.
	for _, e := range entries {
		// Only perform actions that are default actions. I.e. don't
		// drop to shell.
		if e.IsDefault() {
			fmt.Printf("Attempting to boot %s.\n\n", e.Label())
			if err := e.Do(); err == errStopTestOnly {
				return
			} else if err != nil {
				log.Printf("Failed to boot %s: %v", e.Label(), err)
			}
		}
	}
}

// OSImages returns menu entries for the given OSImages.
func OSImages(dryRun bool, imgs ...boot.OSImage) []Entry {
	var menu []Entry
	for _, img := range imgs {
		menu = append(menu, &OSImageAction{
			OSImage: img,
			DryRun:  dryRun,
		})
	}
	return menu
}

// OSImageAction is a menu.Entry that boots an OSImage.
type OSImageAction struct {
	boot.OSImage
	DryRun bool
}

// Do implements Entry.Do by booting the image.
func (oia OSImageAction) Do() error {
	if err := oia.OSImage.Load(oia.DryRun); err != nil {
		log.Printf("Could not load image %s: %v", oia.OSImage, err)
	}
	if oia.DryRun {
		// err should only be nil in a dry run.
		log.Printf("Loaded kernel %s.", oia.OSImage)
		os.Exit(0)
	}
	if err := boot.Execute(); err != nil {
		return err
	}
	return nil
}

// IsDefault returns true -- this action should be performed in order by
// default if the user did not choose a boot entry.
func (OSImageAction) IsDefault() bool { return true }

// StartShell is a menu.Entry that starts a LinuxBoot shell.
type StartShell struct{}

// Label is the label to show to the user.
func (StartShell) Label() string {
	return "Enter a LinuxBoot shell"
}

// Do implements Entry.Do by running /bin/defaultsh.
func (StartShell) Do() error {
	// Reset signal handler for SIGINT to enable user interrupts again
	signal.Reset(syscall.SIGINT)
	return sh.RunWithLogs("/bin/defaultsh")
}

// IsDefault indicates that this should not be run as a default action.
func (StartShell) IsDefault() bool { return false }

// Reboot is a menu.Entry that reboots the machine.
type Reboot struct{}

// Label is the label to show to the user.
func (Reboot) Label() string {
	return "Reboot"
}

// Do reboots the machine using sys_reboot.
func (Reboot) Do() error {
	unix.Sync()
	return unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}

// IsDefault indicates that this should not be run as a default action.
func (Reboot) IsDefault() bool { return false }
