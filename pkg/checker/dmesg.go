// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"strings"

	"golang.org/x/sys/unix"
)

func getDmesg() (string, error) {
	b := make([]byte, 256*1024)
	n, err := unix.Klogctl(unix.SYSLOG_ACTION_READ_ALL, b)
	if err != nil {
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
