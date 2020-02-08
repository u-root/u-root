// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"log"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/termios"
)

// TestSh tests a bash-like line completer.
// It is intended to be used interactively for now
// and will exit if INTERACTIVE is not set.
func TestSh(t *testing.T) {
	if os.Getenv("INTERACTIVE") == "" {
		t.Skip()
	}
	Debug = t.Logf
	tty, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := tty.Raw()
	if err != nil {
		log.Printf("non-fatal cannot get tty: %v", err)
	}
	defer func() {
		if err := tty.Set(r); err != nil {
			t.Error(err)
		}
	}()

	f := NewFileCompleter("")
	p, err := NewPathCompleter()
	if err != nil {
		log.Fatal(err)
	}

	bin := NewMultiCompleter(NewStringCompleter([]string{"exit"}), p, f)
	l := NewNewerLineReader(bin, f)
	l.Prompt = "Prompt% "
	for !l.EOF {
		if err := l.ReadLine(tty, tty); err != nil {
			t.Logf("looperr: %v", err)
		}
		if _, err := tty.Write([]byte("\r\n")); err != nil {
			t.Error(err)
		}
	}
	t.Log("All done!")
}
