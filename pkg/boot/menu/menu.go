// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package menu displays a Terminal UI based text menu to choose boot options
// from.
package menu

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/sh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
)

var (
	initialTimeout    = 10 * time.Second
	subsequentTimeout = 60 * time.Second
)

// Entry is a menu entry.
type Entry interface {
	// Label is the string displayed to the user in the menu.
	Label() string

	// Load is called when the entry is chosen, but does not transfer
	// execution to another process or kernel.
	Load(ctx context.Context) error

	// Exec transfers execution to another process or kernel.
	//
	// Exec either returns an error or does not return at all.
	Exec() error

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
	fmt.Println("\r")

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
					fmt.Printf("BUG: Please report: Terminal read error: %v.\r\n", err)
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
				fmt.Printf("%s is not a valid entry number: %v.\r\n", choice, err)
				continue
			}
			if num-1 < 0 || num > len(entries) {
				fmt.Printf("%s is not a valid entry number.\r\n", choice)
				continue
			}
			boot <- entries[num-1]
			return
		}
	}()

	select {
	case entry := <-boot:
		if entry != nil {
			fmt.Printf("Chosen option %s.\r\n\r\n", entry.Label())
		}
		return entry

	case <-t.C:
		return nil
	}
}

func loadEntry(ctx context.Context, entry Entry) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// TODO: should the parent context really be passed in here? Should there even be one?
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-sigChan
		cancel()
	}()

	err := entry.Load(ctx)

	// Stop the signals before closing the channel to avoid writes to
	// sigChan.
	signal.Stop(sigChan)
	close(sigChan)

	// Closing the channel should lead to cancel and Done being called.
	wg.Wait()

	return err
}

// ShowMenuAndLoad lets the user choose one of entries and loads it. If no
// entry is chosen by the user, an entry whose IsDefault() is true will be
// returned.
//
// The user is left to call Entry.Exec when this function returns.
func ShowMenuAndLoad(ctx context.Context, input *os.File, entries ...Entry) Entry {
	// Clear the screen (ANSI terminal escape code for screen clear).
	fmt.Printf("\033[1;1H\033[2J\n\n")
	fmt.Printf("Welcome to NERF's Boot Menu\n\n")
	fmt.Printf("Enter a number to boot a kernel. Press Ctrl+C at any point to come back to this screen (if possible).\n")

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

		if err := loadEntry(ctx, entry); err != nil {
			log.Printf("Failed to load %s: %v", entry.Label(), err)
			continue
		}
		// Entry was successfully loaded. Leave it to the caller to
		// exec, so the caller can clean up the OS before rebooting or
		// kexecing (e.g. unmount file systems).
		return entry
	}

	fmt.Println("")

	// We only get one shot at actually booting, so boot the first kernel
	// that can be loaded correctly.
	for _, e := range entries {
		// Only perform actions that are default actions. I.e. don't
		// drop to shell.
		if e.IsDefault() {
			fmt.Printf("Attempting to load %s.\n\n", e)

			if err := loadEntry(ctx, e); err != nil {
				log.Printf("Failed to load %s: %v", e.Label(), err)
				continue
			}

			// Entry was successfully loaded. Leave it to the
			// caller to exec, so the caller can clean up the OS
			// before rebooting or kexecing (e.g. unmount file
			// systems).
			return e
		}
	}
	return nil
}

// OSImages returns menu entries for the given OSImages.
func OSImages(verbose bool, imgs ...boot.OSImage) []Entry {
	var menu []Entry
	for _, img := range imgs {
		menu = append(menu, &OSImageAction{
			OSImage: img,
			Verbose: verbose,
		})
	}
	return menu
}

// OSImageAction is a menu.Entry that boots an OSImage.
type OSImageAction struct {
	boot.OSImage
	Verbose bool
}

// Load implements Entry.Load by loading the OS image into memory.
func (oia OSImageAction) Load(ctx context.Context) error {
	if err := oia.OSImage.Load(ctx, oia.Verbose); err != nil {
		return fmt.Errorf("could not load image %s: %v", oia.OSImage, err)
	}
	return nil
}

// Exec executes the loaded image.
func (oia OSImageAction) Exec() error {
	return boot.Execute()
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

// Load does nothing.
func (StartShell) Load(context.Context) error {
	return nil
}

// Exec implements Entry.Exec by running /bin/defaultsh.
func (StartShell) Exec() error {
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

// Load does nothing.
func (Reboot) Load(context.Context) error {
	return nil
}

// Exec reboots the machine using sys_reboot.
func (Reboot) Exec() error {
	unix.Sync()
	return unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}

// IsDefault indicates that this should not be run as a default action.
func (Reboot) IsDefault() bool { return false }

// TimedTestEntry is an entry that returns a specific error after a specific duration of Loading.
type TimedTestEntry struct {
	LoadDur time.Duration
	LoadErr error
	ExecErr error
}

// Label is the label to show the user.
func (tte TimedTestEntry) Label() string {
	return fmt.Sprintf("Timed Test Entry (%s)", tte.LoadDur)
}

// Load times out after LoadDur or when ctx exits.
func (tte TimedTestEntry) Load(ctx context.Context) error {
	select {
	case <-time.After(tte.LoadDur):
		return tte.LoadErr

	case <-ctx.Done():
		return ctx.Err()
	}
}

// Exec returns ExecErr.
func (tte TimedTestEntry) Exec() error {
	return tte.ExecErr
}

// IsDefault returns false.
func (tte TimedTestEntry) IsDefault() bool {
	return false
}
