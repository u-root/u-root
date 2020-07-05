// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!windows

package termios

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"

	"golang.org/x/sys/unix"
)

// GTTY returns the TTY struct for a given fd. It is like a New in
// many packages but the name GTTY is a tradition.
func GTTY(fd int) (*TTY, error) {
	term, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}
	w, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	}

	var t = TTY{Opts: make(map[string]bool), CC: make(map[string]uint8)}
	for n, b := range boolFields {
		val := uint32(reflect.ValueOf(term).Elem().Field(b.word).Uint()) & b.mask
		t.Opts[n] = val != 0
	}

	for n, c := range cc {
		t.CC[n] = term.Cc[c]
	}

	// back in the day, you could have different i and o speeds.
	// since about 1975, this has not been a thing. It's still in POSIX
	// evidently. WTF?
	t.Ispeed = int(term.Ispeed)
	t.Ospeed = int(term.Ospeed)
	t.Row = int(w.Row)
	t.Col = int(w.Col)

	return &t, nil
}

// STTY uses a TTY * to set TTY settings on an fd.
// It returns a new TTY struct for the fd after the changes are made,
// and an error. It does not change the original TTY struct.
func (t *TTY) STTY(fd int) (*TTY, error) {
	// Get a unix.Termios which we can partially fill in.
	term, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	for n, b := range boolFields {
		set := t.Opts[n]
		i := reflect.ValueOf(term).Elem().Field(b.word).Uint()
		if set {
			i |= uint64(b.mask)
		} else {
			i &= ^uint64(b.mask)
		}
		reflect.ValueOf(term).Elem().Field(b.word).SetUint(i)
	}

	for n, c := range cc {
		term.Cc[c] = t.CC[n]
	}

	term.Ispeed = uint32(t.Ispeed)
	term.Ospeed = uint32(t.Ospeed)

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, term); err != nil {
		return nil, err
	}

	w := &unix.Winsize{Row: uint16(t.Row), Col: uint16(t.Col)}
	if err := unix.IoctlSetWinsize(fd, unix.TIOCSWINSZ, w); err != nil {
		return nil, err
	}

	return GTTY(fd)
}

// String will stringify a TTY, including printing out the options all in the same order.
func (t *TTY) String() string {
	s := fmt.Sprintf("speed:%v ", t.Ispeed)
	s += fmt.Sprintf("rows:%d cols:%d", t.Row, t.Col)
	var opts []string
	for n, c := range t.CC {
		opts = append(opts, fmt.Sprintf("%v:%#02x", n, c))
	}

	for n, set := range t.Opts {
		if set {
			opts = append(opts, fmt.Sprintf("%v:1", n))
		} else {
			opts = append(opts, fmt.Sprintf("%v:0", n))
		}
	}
	sort.Strings(opts)
	for _, v := range opts {
		s += fmt.Sprintf(" %s", v)
	}
	return s
}

func intarg(s []string) int {
	if len(s) < 2 {
		log.Fatalf("%s requires an arg", s[0])
	}
	i, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatalf("%s is not a number", s)
	}
	return i
}

// SetOpts sets opts in a TTY given an array of key-value pairs and
// booleans. The arguments are a variety of key-value pairs and booleans.
// booleans are cleared if the first char is a -, set otherwise.
func (t *TTY) SetOpts(opts []string) error {
	for i := 0; i < len(opts); i++ {
		o := opts[i]
		switch o {
		case "row":
			t.Row = intarg(opts[i:])
			i++
			continue
		case "col":
			t.Col = intarg(opts[i:])
			i++
			continue
		case "speed":
			t.Ispeed = intarg(opts[i:])
			i++
			continue
		}

		// see if it's one of the control char options.
		if _, ok := cc[opts[i]]; ok {
			t.CC[opts[i]] = uint8(intarg(opts[i:]))
			i++
			continue
		}

		// At this point, it has to be one of the boolean ones
		// or we're done here.
		set := true
		if o[0] == '~' {
			set = false
			o = o[1:]
		}
		if _, ok := boolFields[o]; !ok {
			log.Fatalf("%s: unknown option", o)
		}

		t.Opts[o] = set
	}
	return nil
}

// Raw sets a TTY into raw more, returning a TTY struct
func Raw(fd int) (*TTY, error) {
	t, err := GTTY(fd)
	if err != nil {
		return nil, err
	}

	t.SetOpts([]string{"~ignbrk", "~brkint", "~parmrk", "~istrip", "~inlcr", "~igncr", "~icrnl", "~ixon", "~opost", "~echo", "~echonl", "~icanon", "~isig", "~iexten", "~parenb" /*"cs8", */, "min", "1", "time", "0"})

	return t.STTY(fd)
}
