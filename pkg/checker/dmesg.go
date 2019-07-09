// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"strings"
	"syscall"
	"unsafe"
)

// shamelessly copied from u-root/cmds/dmesg
const (
	_SYSLOG_ACTION_READ_ALL = 3
)

func getDmesg() (string, error) {
	level := uintptr(_SYSLOG_ACTION_READ_ALL)
	b := make([]byte, 256*1024)
	n, _, err := syscall.Syscall(syscall.SYS_SYSLOG, level, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if err != 0 {
		return "", err
	}
	return string(b[:n]), nil
}

func grep(b, pattern string) []string {
	lines := strings.Split(b, "\n")
	ret := make([]string, 0)
	for _, line := range lines {
		if strings.Contains(line, pattern) {
			ret = append(ret, line)
		}
	}
	return ret
}
