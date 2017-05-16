// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type (
	// tty is an os-independent version of the combined info in termios and window size structs.
	// It is used to get/set info to the termios functions as well as marshal/unmarshal data
	// in JSON formwt for dump and loading.
	tty struct {
		Ispeed int
		Ospeed int
		Row    int
		Col    int
		Line   int
		Min    int
		Time   int

		CC map[string]uint8

		Opts map[string]bool
	}

	winsize struct {
		Row, Col       uint16
		Xpixel, Ypixel uint16
	}
)

func gtty(fd uintptr) (*tty, error) {
	term, err := tiGet(fd)
	if err != nil {
		return nil, err
	}
	w, err := wsGet(fd)
	if err != nil {
		return nil, err
	}

	var t = tty{Opts: make(map[string]bool), CC: make(map[string]uint8)}
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

func stty(fd uintptr, t *tty) (*tty, error) {

	// Get a syscall.Termios
	// which we can partially fill in.
	term, err := tiGet(fd)
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

	if err := tiSet(fd, term); err != nil {
		return nil, err
	}

	w := &winsize{Row: uint16(t.Row), Col: uint16(t.Col)}
	if err := wsSet(fd, w); err != nil {
		return nil, err
	}

	return gtty(fd)
}

func pretty(w io.Writer, t *tty) {
	fmt.Printf("speed: %v ", t.Ispeed)
	for n, c := range t.CC {
		fmt.Printf("%v: %#q, ", n, c)
	}
	fmt.Fprintf(w, "%d rows, %d cols\n", t.Row, t.Col)

	for n, set := range t.Opts {
		if set {
			fmt.Fprintf(w, "%v ", n)
		} else {
			fmt.Fprintf(w, "~%v ", n)
		}
	}
	fmt.Fprintln(w)
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

// the arguments are a variety of key-value pairs and booleans.
// booleans are cleared if the first char is a -, set otherwise.
func setOpts(t *tty, opts []string) error {
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

func setRaw(fd uintptr) (*tty, error) {
	t, err := gtty(fd)
	if err != nil {
		return nil, err
	}

	setOpts(t, []string{"~ignbrk", "~brkint", "~parmrk", "~istrip", "~inlcr", "~igncr", "~icrnl", "~ixon", "~opost", "~echo", "~echonl", "~icanon", "~isig", "~iexten", "~parenb" /*"cs8", */, "min", "1", "time", "0"})

	return stty(fd, t)
}
