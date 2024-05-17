// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

type ProcNode struct {
	Name  string
	PID   int
	Inode int
}

var procFS = "/proc"

func readProgFS() (map[int]ProcNode, error) {
	fs, err := os.ReadDir(procFS)
	if err != nil {
		return nil, err
	}

	retMap := make(map[int]ProcNode)

	for _, entry := range fs {
		if !entry.IsDir() {
			continue
		}

		// We only want numbers
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			// Most likely not a number if we have an error
			continue
		}

		fdpath := path.Join(procFS, entry.Name(), "fd", "1")
		fdlink, err := os.Readlink(fdpath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		fdnumStr, ok := strings.CutPrefix(fdlink, "socket:[")
		if !ok {
			continue
		}

		fdnumStr1, _ := strings.CutSuffix(fdnumStr, "]")

		sockInode, err := strconv.Atoi(fdnumStr1)
		if err != nil {
			return nil, err
		}

		cmdpath := path.Join(procFS, entry.Name(), "cmdline")
		cmdline, err := os.ReadFile(cmdpath)
		if err != nil {
			return nil, err
		}

		pNode := ProcNode{
			Name:  string(cmdline[:20]),
			PID:   pid,
			Inode: sockInode,
		}

		retMap[sockInode] = pNode
	}

	return retMap, nil
}
