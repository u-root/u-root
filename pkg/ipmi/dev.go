// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/unix"
)

type dev struct {
	f *os.File
	syscalls
}

// SendRequest uses unix.Syscall IOCTL to send a request to the BMC.
func (d *dev) SendRequest(req *request) error {
	_, _, err := d.syscall(unix.SYS_IOCTL, d.File().Fd(), _IPMICTL_SEND_COMMAND, uintptr(unsafe.Pointer(req)))
	runtime.KeepAlive(req)
	if err != 0 {
		return fmt.Errorf("syscall failed with: %v", err)
	}
	return nil
}

// ReceiveResponse uses syscall Rawconn to read a response via unix.Syscall IOCTL from the BMC.
// It takes the message ID of the request and awaits the response with the same message ID.
func (d *dev) ReceiveResponse(msgID int64, resp *response, buf []byte) ([]byte, error) {
	var result []byte
	var rerr error
	readMsg := func(fd uintptr) bool {
		_, _, errno := d.syscall(unix.SYS_IOCTL, d.File().Fd(), _IPMICTL_RECEIVE_MSG_TRUNC, uintptr(unsafe.Pointer(resp)))
		runtime.KeepAlive(resp)
		if errno != 0 {
			rerr = fmt.Errorf("ioctlGetRecv failed with %v", errno)
			return false
		}

		if resp.msgid != msgID {
			rerr = fmt.Errorf("received wrong message")
			return false
		}

		if resp.msg.DataLen >= _IPMI_BUF_SIZE {
			rerr = fmt.Errorf("data length received too large: %d > %d", resp.msg.DataLen, _IPMI_BUF_SIZE)
		} else if buf[0] != 0 {
			rerr = fmt.Errorf("invalid response, expected first byte of response to be 0, got: %v", buf[0])
		} else {
			result = buf[:resp.msg.DataLen:resp.msg.DataLen]
			rerr = nil
		}
		return true
	}

	// Read response.
	conn, err := d.fileSyscallConn(d.File())
	if err != nil {
		return nil, fmt.Errorf("failed to get file rawconn: %v", err)
	}
	if err := d.fileSetReadDeadline(d.File(), timeout); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %v", err)
	}
	if err := d.connRead(readMsg, conn); err != nil {
		return nil, fmt.Errorf("failed to read rawconn: %v", err)
	}

	return result, rerr
}

// File returns the file of Dev
func (d *dev) File() *os.File {
	return d.f
}

// Close closes the file attached to dev.
func (d *dev) Close() error {
	return d.f.Close()
}

// newDev takes a file and returns a new Dev
func newDev(f *os.File) *dev {
	return &dev{
		f:        f,
		syscalls: &realSyscalls{},
	}
}
