// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// mock is for testing purposes only
type mock struct {
	unixSyscall func(uintptr, uintptr, uintptr, uintptr) (uintptr, uintptr, unix.Errno)

	// fileSyscallConn gives the option to overwrite the actual function call with a dummy function when writing tests.
	fileSyscallConn func(f *os.File) (syscall.RawConn, error)

	// fileSetReadDeadline gives the option to overwrite the actual function call with a dummy function when writing tests.
	fileSetReadDeadline func(f *os.File, t time.Duration) error

	// connRead gives the option to overwrite the actual function call with a dummy function when writing tests.
	connRead func(f func(fd uintptr) bool, conn syscall.RawConn) error
}

func (m *mock) SendRequest(req *request) error {
	return nil
}

func (m *mock) ReceiveResponse(msgID int64, resp *response, buf []byte) ([]byte, error) {
	return buf, nil
}

func (m *mock) GetFile() *os.File {
	return nil
}

func GetMockIPMI() *IPMI {
	return &IPMI{Handler: &mock{}}
}
