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

var completionCodeMessages = map[CompletionCode]string{
	IPMI_CC_OK:                          "command completed normally",
	IPMI_CC_NODE_BUSY:                   "node busy",
	IPMI_CC_INV_CMD:                     "invalid command",
	IPMI_CC_INV_CMD_FOR_LUN:             "command invalid for given LUN",
	IPMI_CC_TIMEOUT:                     "timeout while processing command",
	IPMI_CC_OUT_OF_SPACE:                "out of space",
	IPMI_CC_RES_CANCELED:                "reservation canceled or invalid reservation ID",
	IPMI_CC_REQ_DATA_TRUNC:              "request data truncated",
	IPMI_CC_REQ_DATA_INV_LENGTH:         "request data length invalid",
	IPMI_CC_REQ_DATA_FIELD_EXCEED:       "request data field length limit exceeded",
	IPMI_CC_PARAM_OUT_OF_RANGE:          "parameter out of range",
	IPMI_CC_CANT_RET_NUM_REQ_BYTES:      "cannot return number of requested data bytes",
	IPMI_CC_REQ_DATA_NOT_PRESENT:        "requested sensor, data, or record not present",
	IPMI_CC_INV_DATA_FIELD_IN_REQ:       "invalid data field in request",
	IPMI_CC_ILL_SENSOR_OR_RECORD:        "command illegal for specified sensor or record type",
	IPMI_CC_RESP_COULD_NOT_BE_PRV:       "command response could not be provided",
	IPMI_CC_CANT_RESP_DUPLI_REQ:         "cannot execute duplicated request",
	IPMI_CC_CANT_RESP_SDRR_UPDATE:       "command response could not be provided: SDR repository in update mode",
	IPMI_CC_CANT_RESP_FIRM_UPDATE:       "command response could not be provided: device in firmware update mode",
	IPMI_CC_CANT_RESP_BMC_INIT:          "command response could not be provided: BMC initialization or initialization agent in progress",
	IPMI_CC_DESTINATION_UNAVAILABLE:     "destination unavailable",
	IPMI_CC_INSUFFICIENT_PRIVILEGES:     "cannot execute command due to insufficient privilege level or other security-based restriction",
	IPMI_CC_NOT_SUPPORTED_PRESENT_STATE: "cannot execute command: command, or request parameter(s), not supported in present state",
	IPMI_CC_ILLEGAL_COMMAND_DISABLED:    "cannot execute command: parameter is illegal because command sub-function has been disabled or is unavailable",
	IPMI_CC_UNSPECIFIED_ERROR:           "unspecified error",
}

func (cc CompletionCode) String() string {
	if s, ok := completionCodeMessages[cc]; ok {
		return s
	}
	if 0x01 <= cc && cc <= 0x7e {
		return fmt.Sprintf("device specific (OEM) completion code: 0x%x", byte(cc))
	}
	if 0x80 <= cc && cc <= 0xbe {
		return fmt.Sprintf("command-specific completion code: 0x%x", byte(cc))
	}
	return fmt.Sprintf("unknown completion code in reserved range: 0x%x", byte(cc))
}

func (ce CompletionError) Error() string {
	return fmt.Sprintf("command completed with non-OK code: %s", CompletionCode(ce))
}

type dev struct {
	f *os.File
	syscalls
}

// SendRequest uses unix.Syscall IOCTL to send a request to the BMC.
func (d *dev) SendRequest(req *request) error {
	_, _, err := d.syscall(unix.SYS_IOCTL, d.File().Fd(), _IPMICTL_SEND_COMMAND, uintptr(unsafe.Pointer(req)))
	runtime.KeepAlive(req)
	if err != 0 {
		return fmt.Errorf("syscall failed with: %w", err)
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
			rerr = fmt.Errorf("ioctlGetRecv failed with %w", errno)
			return false
		}

		if resp.msgid != msgID {
			rerr = fmt.Errorf("received wrong message")
			return false
		}

		if resp.msg.DataLen >= _IPMI_BUF_SIZE {
			rerr = fmt.Errorf("data length received too large: %d > %d", resp.msg.DataLen, _IPMI_BUF_SIZE)
		} else if cc := CompletionCode(buf[0]); cc != IPMI_CC_OK {
			rerr = CompletionError(cc)
		} else {
			result = buf[:resp.msg.DataLen:resp.msg.DataLen]
			rerr = nil
		}
		return true
	}

	// Read response.
	conn, err := d.fileSyscallConn(d.File())
	if err != nil {
		return nil, fmt.Errorf("failed to get file rawconn: %w", err)
	}
	if err := d.fileSetReadDeadline(d.File(), timeout); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}
	if err := d.connRead(readMsg, conn); err != nil {
		return nil, fmt.Errorf("failed to read rawconn: %w", err)
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
