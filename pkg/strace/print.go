// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Changes are Copyright 2018 the u-root Authors.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

// DefaultLogMaximumSize is the default LogMaximumSize.
const DefaultLogMaximumSize = 1024

// LogMaximumSize determines the maximum display size for data blobs (read,
// write, etc.).
var LogMaximumSize uint = DefaultLogMaximumSize

// EventMaximumSize determines the maximum size for data blobs (read, write,
// etc.) sent over the event channel. Default is 0 because most clients cannot
// do anything useful with binary text dump of byte array arguments.
var EventMaximumSize uint

func dump(t Task, addr Addr, size uint, maximumBlobSize uint) string {
	origSize := size
	if size > maximumBlobSize {
		size = maximumBlobSize
	}
	if size == 0 {
		return ""
	}

	b := make([]byte, size)
	amt, err := t.Read(addr, b)
	if err != nil {
		return fmt.Sprintf("%#x (error decoding string: %s)", addr, err)
	}

	dot := ""
	if uint(amt) < origSize {
		// ... if we truncated the dump.
		dot = "..."
	}

	return fmt.Sprintf("%#x %q%s", addr, b[:amt], dot)
}

func iovecs(t Task, addr Addr, iovcnt int, printContent bool, maxBytes uint64) string {
	if iovcnt < 0 || iovcnt > 0x10 /*unix.MSG_MAXIOVLEN*/ {
		return fmt.Sprintf("%#x (error decoding iovecs: invalid iovcnt)", addr)
	}
	v := make([]iovec, iovcnt)
	_, err := t.Read(addr, v)
	if err != nil {
		return fmt.Sprintf("%#x (error decoding iovecs: %v)", addr, err)
	}

	var totalBytes uint64
	var truncated bool
	iovs := make([]string, iovcnt)
	for i, vv := range v {
		if vv.S == 0 || !printContent {
			iovs[i] = fmt.Sprintf("{base=%#x, len=%d}", vv.P, vv.S)
			continue
		}

		size := uint64(vv.S)
		if truncated || totalBytes+size > maxBytes {
			truncated = true
			size = maxBytes - totalBytes
		} else {
			totalBytes += uint64(vv.S)
		}

		b := make([]byte, size)
		amt, err := t.Read(vv.P, b)
		if err != nil {
			iovs[i] = fmt.Sprintf("{base=%#x, len=%d, %q..., error decoding string: %v}", vv.P, vv.S, b[:amt], err)
			continue
		}

		dot := ""
		if truncated {
			// Indicate truncation.
			dot = "..."
		}
		iovs[i] = fmt.Sprintf("{base=%#x, len=%d, %q%s}", vv.P, vv.S, b[:amt], dot)
	}

	return fmt.Sprintf("%#x %s", addr, strings.Join(iovs, ", "))
}

func fdpair(t Task, addr Addr) string {
	var fds [2]int32
	_, err := t.Read(addr, &fds)
	if err != nil {
		return fmt.Sprintf("%#x (error decoding fds: %s)", addr, err)
	}

	return fmt.Sprintf("%#x [%d %d]", addr, fds[0], fds[1])
}

// SaneUtsname is an utsname without the weird partially-filled []byte
type SaneUtsname struct {
	Sysname    string
	Nodename   string
	Release    string
	Version    string
	Machine    string
	Domainname string
}

// SaneUname returns a SaneUtsname, i.e. a Unix Time Sharing name without the weird
// partially filled []byte
func SaneUname(u unix.Utsname) SaneUtsname {
	return SaneUtsname{
		Sysname:    convertUname(u.Sysname),
		Nodename:   convertUname(u.Nodename),
		Release:    convertUname(u.Release),
		Version:    convertUname(u.Version),
		Machine:    convertUname(u.Machine),
		Domainname: convertUname(u.Domainname),
	}
}

func convertUname(s [65]uint8) string {
	return string(bytes.TrimRight(s[:], "\x00"))
}

func uname(t Task, addr Addr) string {
	var u unix.Utsname
	if _, err := t.Read(addr, &u); err != nil {
		return fmt.Sprintf("%#x (error decoding utsname: %s)", addr, err)
	}

	return fmt.Sprintf("%#x %#v", addr, SaneUname(u))
}
