// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

type testsyscalls struct {
	forceErrno              unix.Errno
	forceConnReadErr        error
	forcesetReadDeadlineErr error
	forceConnErr            error
}

func (t *testsyscalls) syscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, unix.Errno) {
	if t.forceErrno != 0 {
		return 0, 0, t.forceErrno
	}

	if trap != unix.SYS_IOCTL {
		return 0, 0, unix.EINVAL
	}

	if a1 < 1 {
		return 0, 0, unix.EINVAL
	}

	if a3 < 1 {
		return 0, 0, unix.EINVAL
	}

	switch a2 {
	case _IPMICTL_RECEIVE_MSG_TRUNC, _IPMICTL_SEND_COMMAND:
		return 0, 0, 0
	default:
		return 0, 0, unix.EINVAL
	}
}

func (t *testsyscalls) fileSyscallConn(f *os.File) (syscall.RawConn, error) {
	if t.forceConnErr != nil {
		return nil, t.forceConnErr
	}
	return f.SyscallConn()
}

// This function only need to return nil. The real deal only works on special file descriptors.
func (t *testsyscalls) fileSetReadDeadline(f *os.File, time time.Duration) error {
	if t.forcesetReadDeadlineErr != nil {
		return t.forcesetReadDeadlineErr
	}
	return nil
}

func (t *testsyscalls) connRead(f func(fd uintptr) bool, conn syscall.RawConn) error {
	if t.forceConnReadErr != nil {
		return t.forceConnReadErr
	}
	return conn.Read(f)
}

func TestDev(t *testing.T) {
	df, err := os.CreateTemp("", "ipmi_dummy_file-")
	if err != nil {
		t.Error(err)
	}

	sc := &testsyscalls{}

	d := dev{
		f:        df,
		syscalls: sc,
	}
	defer os.RemoveAll(df.Name())

	for _, tt := range []struct {
		name        string
		req         *request
		resp        *response
		forceErrno  unix.Errno
		wantSendErr error
		wantRecvErr error
		connReadErr error
		deadlineErr error
	}{
		{
			name: "NoError",
			req:  &request{},
			resp: &response{},
		},
		{
			name:        "ForceSysCallError",
			forceErrno:  unix.Errno(1),
			wantSendErr: fmt.Errorf("syscall failed with: operation not permitted"),
			wantRecvErr: fmt.Errorf("failed to read rawconn"),
			req:         &request{},
			resp:        &response{},
		},
		{
			name:        "ForceConnError",
			req:         &request{},
			resp:        &response{},
			wantRecvErr: fmt.Errorf("failed to get file rawconn"),
			connReadErr: fmt.Errorf("Force connRead error"),
		},
		{
			name:        "ForceSetReadDeadlineError",
			req:         &request{},
			resp:        &response{},
			wantRecvErr: fmt.Errorf("failed to set read deadline"),
			deadlineErr: fmt.Errorf("force set read deadline"),
		},
		{
			name: "FailMessageID",
			req: &request{
				msgid: 1,
			},
			resp: &response{
				msgid: 2,
			},
			wantRecvErr: fmt.Errorf("failed to read rawconn: waiting for unsupported file type"),
		},
		{
			name: "ForceLongRecvMsgDataLen",
			req:  &request{},
			resp: &response{
				msg: Msg{
					DataLen: _IPMI_BUF_SIZE + 1,
				},
			},
			wantRecvErr: fmt.Errorf("data length received too large"),
		},
	} {
		sc.forceErrno = tt.forceErrno
		sc.forceConnReadErr = tt.connReadErr
		sc.forceConnErr = tt.connReadErr
		sc.forcesetReadDeadlineErr = tt.deadlineErr

		t.Run("SendRequest"+tt.name, func(t *testing.T) {
			if err := d.SendRequest(tt.req); err != nil {
				if !strings.Contains(err.Error(), tt.wantSendErr.Error()) {
					t.Errorf("d.SendRequest(tt.req) = %q, not %q", err, tt.wantSendErr)
				}
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 1)
			if _, err := d.ReceiveResponse(0, tt.resp, buf); err != nil {
				if !strings.Contains(err.Error(), tt.wantRecvErr.Error()) {
					t.Errorf("d.ReceiveResponse(0, tt.resp, buf) = _, %q, not _, %q", err, tt.wantRecvErr)
				}
			}
		})

	}
}

func TestNewDev(t *testing.T) {
	df, err := os.CreateTemp("", "ipmi_dummy_file-")
	if err != nil {
		t.Errorf(`os.CreateTemp("", "ipmi_dummy_file-") = df, %q, not df, nil`, err)
	}
	defer os.RemoveAll(df.Name())
	_ = newDev(df)
}
