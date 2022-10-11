// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package menu

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"golang.org/x/term"
)

type MenuTerminal interface {
	io.Writer
	ReadLine() (string, error)
	SetPrompt(string)
	SetEntryCallback(func())
	SetTimeout(time.Duration) error
	Close() error
}

var _ = MenuTerminal(&xterm{})

// xterm is a wrapper for term.Terminal following the MenuTerminal interface
type xterm struct {
	term.Terminal
	// Save variables needed to close the xterm
	fileInput *os.File
	oldState  *term.State
}

// NewTerminal opens an xTerminal using the given file input.
// Note that the xterm is in raw mode. Write \r\n whenever you would
// write a \n. When testing in qemu, it might look fine because
// there might be another tty cooking the newlines. In for
// example minicom, the behavior is different. And you would
// see something like:
//
//	Select a boot option to edit:
//	                             >
//
// Instead of:
//
//	Select a boot option to edit:
//	 >
func NewTerminal(f *os.File) *xterm {
	oldState, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		log.Printf("BUG: Please report: We cannot actually let you choose from menu (MakeRaw failed): %v", err)
	}

	if err = syscall.SetNonblock(int(f.Fd()), true); err != nil {
		log.Printf("BUG: Error setting Fd %d to nonblocking: %v", f.Fd(), err)
	}

	return &xterm{
		*term.NewTerminal(f, ""),
		f,
		oldState,
	}
}

func (t *xterm) Close() error {
	if t.oldState == nil {
		return fmt.Errorf("cannot restore terminal state to nil")
	}
	return term.Restore(int(t.fileInput.Fd()), t.oldState)
}

func (t *xterm) SetTimeout(dur time.Duration) error {
	return t.fileInput.SetDeadline(time.Now().Add(dur))
}

// Sets the timeout for the file and adds a timeout refresh on user entry
func (t *xterm) SetEntryCallback(f func()) {
	t.AutoCompleteCallback = func(line string, pos int, key rune) (string, int, bool) {
		f()
		return "", 0, false
	}
}
