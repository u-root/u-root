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

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

// Address is a byte slice cast as a string that represents the address of a
// network node. Or, in the case of unix endpoints, it may represent a path.
type Address string

// FullAddress is the network address and port
type FullAddress struct {
	// Addr is the network address.
	Addr Address

	// Port is the transport port.
	//
	// This may not be used by all endpoint types.
	Port uint16
}

// String implements String
func (a *FullAddress) String() string {
	if a == nil {
		return ":"
	}
	return fmt.Sprintf("%s:%#x", []byte(a.Addr), a.Port)
}

// GetAddress reads an sockaddr struct from the given address and converts it
// to the FullAddress format. It supports AF_UNIX, AF_INET and AF_INET6
// addresses.
func GetAddress(addr []byte) (*FullAddress, error) {
	r := bytes.NewBuffer(addr)

	var fam uint16
	if err := binary.Read(r, binary.NativeEndian, &fam); err != nil {
		return nil, err
	}

	// Get the rest of the fields based on the address family.
	switch fam {
	case unix.AF_UNIX:
		path := r.Bytes()

		if len(path) > unix.PathMax {
			return nil, unix.ENAMETOOLONG
		}

		// Drop the terminating NUL (if one exists) and everything after
		// it for filesystem (non-abstract) addresses.
		if n := bytes.IndexByte(path, 0); n > 0 {
			path = path[:n]
		} else {
			return nil, unix.EINVAL
		}

		return &FullAddress{
			Addr: Address(path),
		}, nil

	case unix.AF_INET:
		var a unix.RawSockaddrInet4
		r = bytes.NewBuffer(addr)
		if err := binary.Read(r, binary.BigEndian, &a); err != nil {
			return nil, unix.EFAULT
		}
		return &FullAddress{
			Addr: Address(net.IP(a.Addr[:]).String()),
			Port: uint16(a.Port),
		}, nil

	case unix.AF_INET6:
		var a unix.RawSockaddrInet6
		r = bytes.NewBuffer(addr)
		if err := binary.Read(r, binary.BigEndian, &a); err != nil {
			return nil, unix.EFAULT
		}

		return &FullAddress{
			Addr: Address(net.IP(a.Addr[:]).String()),
			Port: uint16(a.Port),
		}, nil

	default:
		return nil, unix.ENOTSUP
	}
}
