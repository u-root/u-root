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

func (cc CompletionCode) Error() string {
	switch cc {
	case IPMI_CC_OK:
		return "command completed normally"
	case IPMI_CC_NODE_BUSY:
		return "node busy"
	case IPMI_CC_INV_CMD:
		return "invalid command"
	case IPMI_CC_INV_CMD_FOR_LUN:
		return "command invalid for given LUN"
	case IPMI_CC_TIMEOUT:
		return "timeout while processing command"
	case IPMI_CC_OUT_OF_SPACE:
		return "out of space"
	case IPMI_CC_RES_CANCELED:
		return "reservation canceled or invalid reservation ID"
	case IPMI_CC_REQ_DATA_TRUNC:
		return "request data truncated"
	case IPMI_CC_REQ_DATA_INV_LENGTH:
		return "request data length invalid"
	case IPMI_CC_REQ_DATA_FIELD_EXCEED:
		return "request data field length limit exceeded"
	case IPMI_CC_PARAM_OUT_OF_RANGE:
		return "parameter out of range"
	case IPMI_CC_CANT_RET_NUM_REQ_BYTES:
		return "cannot return number of requested data bytes"
	case IPMI_CC_REQ_DATA_NOT_PRESENT:
		return "requested sensor, data, or record not present"
	case IPMI_CC_INV_DATA_FIELD_IN_REQ:
		return "invalid data field in request"
	case IPMI_CC_ILL_SENSOR_OR_RECORD:
		return "command illegal for specified sensor or record type"
	case IPMI_CC_RESP_COULD_NOT_BE_PRV:
		return "command response could not be provided"
	case IPMI_CC_CANT_RESP_DUPLI_REQ:
		return "cannot execute duplicated request"
	case IPMI_CC_CANT_RESP_SDRR_UPDATE:
		return "command response could not be provided: SDR repository in update mode"
	case IPMI_CC_CANT_RESP_FIRM_UPDATE:
		return "command response could not be provided: device in firmware update mode"
	case IPMI_CC_CANT_RESP_BMC_INIT:
		return "command response could not be provided: BMC initialization or initialization agent in progress"
	case IPMI_CC_DESTINATION_UNAVAILABLE:
		return "destination unavailable"
	case IPMI_CC_INSUFFICIENT_PRIVILEGES:
		return "cannot execute command due to insufficient privilege level or other security-based restriction"
	case IPMI_CC_NOT_SUPPORTED_PRESENT_STATE:
		return "cannot execute command: command, or request parameter(s), not supported in present state"
	case IPMI_CC_ILLEGAL_COMMAND_DISABLED:
		return "cannot execute command: parameter is illegal because command sub-function has been disabled or is unavailable"
	case IPMI_CC_UNSPECIFIED_ERROR:
		return "unspecified error"
	default:
		if 0x01 <= cc && cc <= 0x7e {
			return fmt.Sprintf("device specific (OEM) completion code: %x", byte(cc))
		}
		if 0x80 <= cc && cc <= 0xbe {
			return fmt.Sprintf("command-specific completion code: %x", byte(cc))
		}
		return fmt.Sprintf("unknown completion code in reserved range: %x", byte(cc))
	}
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
		} else if cc := CompletionCode(buf[0]); cc != IPMI_CC_OK {
			rerr = cc
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
