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

package strace

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/u-root/u-root/pkg/ubinary"
	"golang.org/x/sys/unix"
)

// Address is a byte slice cast as a string that represents the address of a
// network node. Or, in the case of unix endpoints, it may represent a path.
type Address string

type FullAddress struct {
	// Addr is the network address.
	Addr Address

	// Port is the transport port.
	//
	// This may not be used by all endpoint types.
	Port uint16
}

// GetAddress reads an sockaddr struct from the given address and converts it
// to the FullAddress format. It supports AF_UNIX, AF_INET and AF_INET6
// addresses.
func GetAddress(t Task, addr []byte) (FullAddress, error) {
	r := bytes.NewBuffer(addr[:2])
	var fam uint16
	if err := binary.Read(r, ubinary.NativeEndian, &fam); err != nil {
		return FullAddress{}, unix.EFAULT
	}

	// Get the rest of the fields based on the address family.
	switch fam {
	case unix.AF_UNIX:
		path := addr[2:]
		if len(path) > unix.PathMax {
			return FullAddress{}, unix.EINVAL
		}
		// Drop the terminating NUL (if one exists) and everything after
		// it for filesystem (non-abstract) addresses.
		if len(path) > 0 && path[0] != 0 {
			if n := bytes.IndexByte(path[1:], 0); n >= 0 {
				path = path[:n+1]
			}
		}
		return FullAddress{
			Addr: Address(path),
		}, nil

	case unix.AF_INET:
		var a unix.RawSockaddrInet4
		r = bytes.NewBuffer(addr)
		if err := binary.Read(r, binary.BigEndian, &a); err != nil {
			return FullAddress{}, unix.EFAULT
		}
		out := FullAddress{
			Addr: Address(a.Addr[:]),
			Port: uint16(a.Port),
		}
		if out.Addr == "\x00\x00\x00\x00" {
			out.Addr = ""
		}
		return out, nil

	case unix.AF_INET6:
		var a unix.RawSockaddrInet6
		r = bytes.NewBuffer(addr)
		if err := binary.Read(r, binary.BigEndian, &a); err != nil {
			return FullAddress{}, unix.EFAULT
		}

		out := FullAddress{
			Addr: Address(a.Addr[:]),
			Port: uint16(a.Port),
		}

		//if isLinkLocal(out.Addr) {
		//			out.NIC = NICID(a.Scope_id)
		//}

		if out.Addr == Address(strings.Repeat("\x00", 16)) {
			out.Addr = ""
		}
		return out, nil

	default:
		return FullAddress{}, unix.ENOTSUP
	}
}
