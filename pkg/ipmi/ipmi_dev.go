// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !test

package ipmi

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Dev struct {
	*os.File
}

func (d *Dev) SendRequest(req *request) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, d.File.Fd(), _IPMICTL_SEND_COMMAND, uintptr(unsafe.Pointer(req)))
	runtime.KeepAlive(req)
	if err != 0 {
		return err
	}
	return nil
}

func (d *Dev) ReceiveResponse(msgID int64, resp *response, buf []byte) ([]byte, error) {
	var result []byte
	var rerr error
	readMsg := func(fd uintptr) bool {
		_, _, err := unix.Syscall(unix.SYS_IOCTL, d.File.Fd(), _IPMICTL_RECEIVE_MSG_TRUNC, uintptr(unsafe.Pointer(resp)))
		runtime.KeepAlive(resp)
		if err != 0 {
			rerr = fmt.Errorf("ioctlGetRecv failed with %v", err)
			return false
		}

		if resp.msgid != msgID {
			log.Printf("Received wrong message. Trying again.")
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
	conn, err := d.GetFile().SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("failed to get file rawconn: %v", err)
	}
	if err := d.GetFile().SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %v", err)
	}
	if err := conn.Read(readMsg); err != nil {
		return nil, fmt.Errorf("failed to read rawconn: %v", err)
	}

	return result, rerr
}

func (d *Dev) GetFile() *os.File {
	return d.File
}

func GetDev(f *os.File) *Dev {
	return &Dev{File: f}
}
