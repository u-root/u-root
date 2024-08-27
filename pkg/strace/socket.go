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
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/align"
	"github.com/u-root/u-root/pkg/strace/internal/abi"
	sbinary "github.com/u-root/u-root/pkg/strace/internal/binary"
	"golang.org/x/sys/unix"
)

func cmsghdr(t Task, addr Addr, length uint64, maxBytes uint64) string {
	if length > maxBytes {
		return fmt.Sprintf("%#x (error decoding control: invalid length (%d))", addr, length)
	}

	buf := make([]byte, length)
	if _, err := t.Read(addr, &buf); err != nil {
		return fmt.Sprintf("%#x (error decoding control: %v)", addr, err)
	}

	var strs []string

	for i := 0; i < len(buf); {
		if i+abi.SizeOfControlMessageHeader > len(buf) {
			strs = append(strs, "{invalid control message (too short)}")
			break
		}

		var h abi.ControlMessageHeader
		sbinary.Unmarshal(buf[i:i+abi.SizeOfControlMessageHeader], binary.NativeEndian, &h)

		var skipData bool
		level := "SOL_SOCKET"
		if h.Level != unix.SOL_SOCKET {
			skipData = true
			level = fmt.Sprint(h.Level)
		}

		typ, ok := abi.ControlMessageType[h.Type]
		if !ok {
			skipData = true
			typ = fmt.Sprint(h.Type)
		}

		if h.Length > uint64(len(buf)-i) {
			strs = append(strs, fmt.Sprintf(
				"{level=%s, type=%s, length=%d, content extends beyond buffer}",
				level,
				typ,
				h.Length,
			))
			break
		}

		i += abi.SizeOfControlMessageHeader
		// TODO: uh, what
		width := archWidth
		length := int(h.Length) - abi.SizeOfControlMessageHeader

		if skipData {
			strs = append(strs, fmt.Sprintf("{level=%s, type=%s, length=%d}", level, typ, h.Length))
			i += int(align.Up(uint(length), uint(width)))
			continue
		}

		switch h.Type {
		case unix.SCM_RIGHTS:
			rightsSize := int(align.Down(uint(length), abi.SizeOfControlMessageRight))

			numRights := rightsSize / abi.SizeOfControlMessageRight
			fds := make(abi.ControlMessageRights, numRights)
			sbinary.Unmarshal(buf[i:i+rightsSize], binary.NativeEndian, &fds)

			rights := make([]string, 0, len(fds))
			for _, fd := range fds {
				rights = append(rights, fmt.Sprint(fd))
			}

			strs = append(strs, fmt.Sprintf(
				"{level=%s, type=%s, length=%d, content: %s}",
				level,
				typ,
				h.Length,
				strings.Join(rights, ","),
			))

		case unix.SCM_CREDENTIALS:
			if length < abi.SizeOfControlMessageCredentials {
				strs = append(strs, fmt.Sprintf(
					"{level=%s, type=%s, length=%d, content too short}",
					level,
					typ,
					h.Length,
				))
				break
			}

			var creds abi.ControlMessageCredentials
			sbinary.Unmarshal(buf[i:i+abi.SizeOfControlMessageCredentials], binary.LittleEndian, &creds)

			strs = append(strs, fmt.Sprintf(
				"{level=%s, type=%s, length=%d, pid: %d, uid: %d, gid: %d}",
				level,
				typ,
				h.Length,
				creds.PID,
				creds.UID,
				creds.GID,
			))

		case unix.SO_TIMESTAMP:
			if length < abi.SizeOfTimeval {
				strs = append(strs, fmt.Sprintf(
					"{level=%s, type=%s, length=%d, content too short}",
					level,
					typ,
					h.Length,
				))
				break
			}

			var tv unix.Timeval
			sbinary.Unmarshal(buf[i:i+abi.SizeOfTimeval], binary.NativeEndian, &tv)

			strs = append(strs, fmt.Sprintf(
				"{level=%s, type=%s, length=%d, Sec: %d, Usec: %d}",
				level,
				typ,
				h.Length,
				tv.Sec,
				tv.Usec,
			))

		default:
			panic("unreachable")
		}
		i += int(align.Up(uint(length), uint(width)))
	}

	return fmt.Sprintf("%#x %s", addr, strings.Join(strs, ", "))
}

func msghdr(t Task, addr Addr, printContent bool, maxBytes uint64) string {
	var msg abi.MessageHeader64
	if _, err := t.Read(addr, &msg); err != nil {
		return fmt.Sprintf("%#x (error decoding msghdr: %v)", addr, err)
	}

	s := fmt.Sprintf(
		"%#x {name=%#x, namelen=%d, iovecs=%s",
		addr,
		msg.Name,
		msg.NameLen,
		iovecs(t, Addr(msg.Iov), int(msg.IovLen), printContent, maxBytes),
	)
	if printContent {
		s = fmt.Sprintf("%s, control={%s}", s, cmsghdr(t, Addr(msg.Control), msg.ControlLen, maxBytes))
	} else {
		s = fmt.Sprintf("%s, control=%#x, control_len=%d", s, msg.Control, msg.ControlLen)
	}
	return fmt.Sprintf("%s, flags=%d}", s, msg.Flags)
}

func sockAddr(t Task, addr Addr, length uint32) string {
	if addr == 0 {
		return "null"
	}

	b, err := CaptureAddress(t, addr, length)
	if err != nil {
		return fmt.Sprintf("%#x {error reading address: %v}", addr, err)
	}

	// Extract address family.
	if len(b) < 2 {
		return fmt.Sprintf("%#x {address too short: %d bytes}", addr, len(b))
	}
	family := binary.NativeEndian.Uint16(b)

	familyStr := abi.SocketFamily.Parse(uint64(family))

	switch family {
	case unix.AF_INET, unix.AF_INET6, unix.AF_UNIX:
		fa, err := GetAddress(b)
		if err != nil {
			return fmt.Sprintf("%#x {Family: %s, error extracting address: %v}", addr, familyStr, err)
		}

		if family == unix.AF_UNIX {
			return fmt.Sprintf("%#x {Family: %s, Addr: %q}", addr, familyStr, fa.Addr)
		}

		return fmt.Sprintf("%#x {Family: %s, Addr: %#02x, Port: %d}", addr, familyStr, []byte(fa.Addr), fa.Port)
	case unix.AF_NETLINK:
		// sa, err := netlink.ExtractSockAddr(b)
		// if err != nil {
		return fmt.Sprintf("%#x {Family: %s, error extracting address: %v}", addr, familyStr, err)
		//}
		//return fmt.Sprintf("%#x {Family: %s, PortID: %d, Groups: %d}", addr, familyStr, sa.PortID, sa.Groups)
	default:
		return fmt.Sprintf("%#x {Family: %s, family addr format unknown}", addr, familyStr)
	}
}

func postSockAddr(t Task, addr Addr, lengthPtr Addr) string {
	if addr == 0 {
		return "null"
	}

	if lengthPtr == 0 {
		return fmt.Sprintf("%#x {length null}", addr)
	}

	l, err := copySockLen(t, lengthPtr)
	if err != nil {
		return fmt.Sprintf("%#x {error reading length: %v}", addr, err)
	}

	return sockAddr(t, addr, l)
}

func copySockLen(t Task, addr Addr) (uint32, error) {
	// socklen_t is 32-bits.
	var l uint32
	_, err := t.Read(addr, &l)
	return l, err
}

func sockLenPointer(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}
	l, err := copySockLen(t, addr)
	if err != nil {
		return fmt.Sprintf("%#x {error reading length: %v}", addr, err)
	}
	return fmt.Sprintf("%#x {length=%v}", addr, l)
}
