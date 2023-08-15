// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"strconv"
	"syscall"
)

// GTTY returns the TTY struct for a given fd. It is like a New in
// many packages but the name GTTY is a tradition.
func GTTY(fd int) (*TTY, error) {
	t := TTY{Opts: make(map[string]bool), CC: make(map[string]uint8)}
	// back in the day, you could have different i and o speeds.
	// since about 1975, this has not been a thing. It's still in POSIX
	// evidently. WTF?
	t.Ispeed = 115200
	t.Ospeed = 115200
	t.Row = int(24)
	t.Col = int(80)

	return &t, nil
}

// STTY uses a TTY * to set TTY settings on an fd.
// It returns a new TTY struct for the fd after the changes are made,
// and an error. It does not change the original TTY struct.
func (t *TTY) STTY(fd int) (*TTY, error) {
	return nil, syscall.ENOSYS
}

// String will stringify a TTY, including printing out the options all in the same order.
// The options are presented in the order:
// integer options as name:value
// boolean options which are set, printed as name, sorted by name
// boolean options which are clear, printed as ~name, sorted by name
// This ordering makes it a bit more readable: integer value, sorted set values, sorted clear values
func (t *TTY) String() string {
	s := fmt.Sprintf("speed:%v ", t.Ispeed)
	s += fmt.Sprintf("rows:%d cols:%d", t.Row, t.Col)
	return s
}

func intarg(s []string, bits int) (int, error) {
	if len(s) < 2 {
		return -1, fmt.Errorf("%s requires an arg", s[0])
	}
	i, err := strconv.ParseUint(s[1], 0, bits)
	if err != nil {
		return -1, fmt.Errorf("%s is not a number", s)
	}
	return int(i), nil
}

// SetOpts sets opts in a TTY given an array of key-value pairs and
// booleans. The arguments are a variety of key-value pairs and booleans.
// booleans are cleared if the first char is a -, set otherwise.
func (t *TTY) SetOpts(opts []string) error {
	var err error
	for i := 0; i < len(opts) && err == nil; i++ {
		o := opts[i]
		switch o {
		case "rows":
			t.Row, err = intarg(opts[i:], 16)
			i++
			continue
		case "cols":
			t.Col, err = intarg(opts[i:], 16)
			i++
			continue
		case "speed":
			// 32 may sound crazy but ... baud can be REALLY large
			t.Ispeed, err = intarg(opts[i:], 32)
			i++
			continue
		}
	}

	return err
}

// Raw sets a TTY into raw mode, returning a TTY struct
func Raw(fd int) (*TTY, error) {
	t, err := GTTY(fd)
	if err != nil {
		return nil, err
	}
	return t.STTY(fd)
}
