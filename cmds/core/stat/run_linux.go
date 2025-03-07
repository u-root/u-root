// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !arm && !386 && !mips && !mipsle

package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

// device inode perm nlink uid gid rdev size atime mtime
// ctime btime blksize blocks st_flags filename
const defFormat = "%d %d %s %d %d %d %d %d %q %q %q %d %d %s\n"

func fromTimespec(ts syscall.Timespec) string {
	t := time.Unix(ts.Sec, ts.Nsec)
	return t.Format("Jan  2 15:04:05 2006")
}

func run(stdout io.Writer, stderr io.Writer, files ...string) int {
	var errCode int
	for _, file := range files {
		fi, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintf(stderr, "stat: %v", err)
			errCode = 1
			continue
		}

		stat, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			fmt.Fprintf(stderr, "stat: %v", err)
			errCode = 1
			continue
		}

		fmt.Fprintf(stdout, defFormat,
			stat.Dev,
			stat.Ino,
			fi.Mode().Perm(),
			stat.Nlink,
			stat.Uid,
			stat.Gid,
			stat.Rdev,
			stat.Size,
			fromTimespec(stat.Atim),
			fromTimespec(stat.Mtim),
			fromTimespec(stat.Ctim),
			stat.Blksize,
			stat.Blocks,
			fi.Name(),
		)
	}

	return errCode
}
