// Copyright (c) 2018, Google Inc. All rights reserved.
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

// Package tpmutil provides common utility functions for both TPM 1.2 and TPM 2.0 devices.
//
// Users should call either UseTPM12LengthPrefixSize or
// UseTPM20LengthPrefixSize before using this package, depending on their type
// of TPM device.
package tpmutil

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// OpenTPM opens a channel to the TPM at the given path. If the file is a
// device, then it treats it like a normal TPM device, and if the file is a
// Unix domain socket, then it opens a connection to the socket.
func OpenTPM(path string) (io.ReadWriteCloser, error) {
	// If it's a regular file, then open it
	var rwc io.ReadWriteCloser
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeDevice != 0 {
		var f *os.File
		f, err = os.OpenFile(path, os.O_RDWR, 0600)
		if err != nil {
			return nil, err
		}
		rwc = io.ReadWriteCloser(f)
	} else if fi.Mode()&os.ModeSocket != 0 {
		uc, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
		if err != nil {
			return nil, err
		}
		rwc = io.ReadWriteCloser(uc)
	} else {
		return nil, fmt.Errorf("unsupported TPM file mode %s", fi.Mode().String())
	}

	return rwc, nil
}

const maxTPMResponse = 4096

// RunCommand executes cmd with given tag and arguments. Returns TPM response
// body (without response header) and response code from the header. Returned
// error may be nil if response code is not RCSuccess; caller should check
// both.
func RunCommand(rw io.ReadWriter, tag Tag, cmd Command, in ...interface{}) ([]byte, ResponseCode, error) {
	if rw == nil {
		return nil, 0, errors.New("nil TPM handle")
	}

	ch := commandHeader{tag, 0, cmd}
	inb, err := packWithHeader(ch, in...)
	if err != nil {
		return nil, 0, err
	}

	if _, err := rw.Write(inb); err != nil {
		return nil, 0, err
	}

	outb := make([]byte, maxTPMResponse)
	outlen, err := rw.Read(outb)
	if err != nil {
		return nil, 0, err
	}
	// Resize the buffer to match the amount read from the TPM.
	outb = outb[:outlen]

	var rh responseHeader
	read, err := Unpack(outb, &rh)
	if err != nil {
		return nil, 0, err
	}
	outb = outb[read:]

	if rh.Res != RCSuccess {
		return nil, rh.Res, nil
	}

	return outb, rh.Res, nil
}
