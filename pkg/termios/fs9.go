// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// N.B.: While these functions are only used on Plan 9,
// they can be tested on any system: they are just doing
// file IO. Until we have Plan 9 VMs to test, we can test
// them in Linux.
package termios

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func consctl(root string, fd uintptr) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "fd", fmt.Sprintf("%dctl", fd)))
	if err != nil {
		return "", err
	}

	fields := strings.Fields(strings.TrimSpace(string(data)))

	if len(fields) == 0 {
		return "", fmt.Errorf("no fields in ctl file")
	}

	return fields[len(fields)-1], nil
}

func consctlFile(root string, fd uintptr) (*os.File, error) {
	s, err := consctl(root, fd)
	if err != nil {
		return nil, err
	}

	// The plan 9 standard is that the ctl file is the device
	// name + the string "ctl"
	s = s + "ctl"

	return os.OpenFile(s, os.O_WRONLY, 0)
}

func readWinSize(n string) (uint16, uint16, error) {
	f, err := os.Open(n)
	if err != nil {
		return 0, 0, err
	}

	defer f.Close()

	// wctl is not a file, it is more like a pipe.
	// We can not read just 48 bytes, for the winsize; rio and
	// lola return "buffer too small".
	// Read all 72 bytes. If at least 48 bytes are returned,
	// that is good enough.
	var b [72]byte
	amt, err := io.ReadFull(f, b[:])
	if amt < 48 {
		return 0, 0, err
	}

	var ulx, uly, lrx, lry uint16
	amt, err = fmt.Sscanf(string(b[:]), "%d %d %d %d", &ulx, &uly, &lrx, &lry)
	if amt != 4 || err != nil {
		return 0, 0, fmt.Errorf("%q:got %d of 4 items:%w", string(b[:]), amt, err)
	}

	return lry - uly, lrx - ulx, nil
}
