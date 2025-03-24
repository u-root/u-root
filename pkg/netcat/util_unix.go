// Copyright 2024-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows && !tamago

package netcat

import (
	"io/fs"
	"os"
	"syscall"
)

// Close implements io.WriteCloser.Close.
func (swc *StdoutWriteCloser) Close() error {
	var devNull *os.File
	var err error
	var fileInfo os.FileInfo

	if devNull, err = os.OpenFile("/dev/null", os.O_WRONLY, 0); err != nil {
		return err
	}
	defer devNull.Close()

	// The close() that's internal to dup2() is silent; therefore, explicitly sync
	// regular files and block devices at first, so that the internal close() have
	// nothing left to do.
	if fileInfo, err = os.Stdout.Stat(); err != nil {
		return err
	}
	if fileInfo.Mode()&fs.ModeType == 0 ||
		fileInfo.Mode()&(fs.ModeDevice|fs.ModeCharDevice) == fs.ModeDevice {
		if err = os.Stdout.Sync(); err != nil {
			return err
		}
	}

	return dup2(int(devNull.Fd()), syscall.Stdout)
}
