// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipmi implements functions to communicate with
// the OpenIPMI driver interface.
package ipmi

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	_IOC_WRITE = 0x1
	_IOC_READ  = 0x2

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_NRSHIFT  = 0

	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IPMI_BMC_CHANNEL                = 0xf
	_IPMI_BUF_SIZE                   = 1024
	_IPMI_IOC_MAGIC                  = 'i'
	_IPMI_NETFN_APP                  = 0x6
	_IPMI_OPENIPMI_READ_TIMEOUT      = 15
	_IPMI_SYSTEM_INTERFACE_ADDR_TYPE = 0x0c

	_BMC_SET_WATCHDOG_TIMER = 0x24
	_BMC_GET_WATCHDOG_TIMER = 0x25

	_SET_SYSTEM_INFO_PARAMETERS = 0x58
	_GET_SYSTEM_INFO_PARAMETERS = 0x59

	_IPM_WATCHDOG_NO_ACTION    = 0x00
	_IPM_WATCHDOG_SMS_OS       = 0x04
	_IPM_WATCHDOG_CLEAR_SMS_OS = 0x10

	_SYSTEM_INFO_BLK_SZ = 16

	_SET_IN_PROGRESS   = 0
	_SYSTEM_FW_VERSION = 1

	_ASCII = 0
)

var (
	_IPMICTL_RECEIVE_MSG  = iowr(_IPMI_IOC_MAGIC, 12, int(unsafe.Sizeof(recv{})))
	_IPMICTL_SEND_COMMAND = ior(_IPMI_IOC_MAGIC, 13, int(unsafe.Sizeof(req{})))
)

type IPMI struct {
	*os.File
}

type msg struct {
	netfn   byte
	cmd     byte
	dataLen uint16
	data    unsafe.Pointer
}

type req struct {
	addr    *systemInterfaceAddr
	addrLen uint32
	msgid   int64 //nolint:structcheck
	msg     msg
}

type recv struct {
	recvType int32 //nolint:structcheck
	addr     *systemInterfaceAddr
	addrLen  uint32
	msgid    int64 //nolint:structcheck
	msg      msg
}

type systemInterfaceAddr struct {
	addrType int32
	channel  int16
	lun      byte //nolint:unused
}

type setSystemInfoReq struct {
	paramSelector byte
	setSelector   byte
	version       [_SYSTEM_INFO_BLK_SZ]byte
}

func ioc(dir int, t int, nr int, size int) int {
	return (dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT)
}

func ior(t int, nr int, size int) int {
	return ioc(_IOC_READ, t, nr, size)
}

func iowr(t int, nr int, size int) int {
	return ioc(_IOC_READ|_IOC_WRITE, t, nr, size)
}

func ioctl(fd uintptr, name int, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(name), uintptr(data))
	return err
}

func fdSet(fd uintptr, p *syscall.FdSet) {
	p.Bits[fd/64] |= 1 << (uint(fd) % 64)
}

func (i *IPMI) sendrecv(req *req) ([]byte, error) {
	addr := systemInterfaceAddr{
		addrType: _IPMI_SYSTEM_INTERFACE_ADDR_TYPE,
		channel:  _IPMI_BMC_CHANNEL,
	}

	req.addr = &addr
	req.addrLen = uint32(unsafe.Sizeof(addr))
	if err := ioctl(i.Fd(), _IPMICTL_SEND_COMMAND, unsafe.Pointer(req)); err != 0 {
		return nil, err
	}

	set := &syscall.FdSet{}
	fdSet(i.Fd(), set)
	time := &syscall.Timeval{
		Sec:  _IPMI_OPENIPMI_READ_TIMEOUT,
		Usec: 0,
	}
	if _, err := syscall.Select(int(i.Fd()+1), set, nil, nil, time); err != nil {
		return nil, err
	}

	recv := &recv{}
	recv.addr = &systemInterfaceAddr{}
	recv.addrLen = uint32(unsafe.Sizeof(addr))
	buf := make([]byte, _IPMI_BUF_SIZE)
	recv.msg.data = unsafe.Pointer(&buf[0])
	recv.msg.dataLen = _IPMI_BUF_SIZE
	if err := ioctl(i.Fd(), _IPMICTL_RECEIVE_MSG, unsafe.Pointer(recv)); err != 0 {
		return nil, err
	}

	return buf[:recv.msg.dataLen:recv.msg.dataLen], nil
}

func (i *IPMI) WatchdogRunning() (bool, error) {
	req := &req{}
	req.msg.cmd = _BMC_GET_WATCHDOG_TIMER
	req.msg.netfn = _IPMI_NETFN_APP

	recv, err := i.sendrecv(req)
	if err != nil {
		return false, err
	}

	if len(recv) > 2 && (recv[1]&0x40) != 0 {
		return true, nil
	}

	return false, nil
}

func (i *IPMI) ShutoffWatchdog() error {
	req := &req{}
	req.msg.cmd = _BMC_SET_WATCHDOG_TIMER
	req.msg.netfn = _IPMI_NETFN_APP

	var data [6]byte
	data[0] = _IPM_WATCHDOG_SMS_OS
	data[1] = _IPM_WATCHDOG_NO_ACTION
	data[2] = 0x00 // pretimeout interval
	data[3] = _IPM_WATCHDOG_CLEAR_SMS_OS
	data[4] = 0xb8 // countdown lsb (100 ms/count)
	data[5] = 0x0b // countdown msb - 5 mins
	req.msg.data = unsafe.Pointer(&data)
	req.msg.dataLen = 6

	_, err := i.sendrecv(req)
	if err != nil {
		return err
	}

	return nil
}

func (i *IPMI) waitForSetInProgressCleared() error {
	req := &req{}
	req.msg.cmd = _GET_SYSTEM_INFO_PARAMETERS
	req.msg.netfn = _IPMI_NETFN_APP
	var data [4]byte

	data[0] = 0
	data[1] = _SET_IN_PROGRESS
	data[2] = 0
	data[3] = 0
	req.msg.data = unsafe.Pointer(&data)
	req.msg.dataLen = 4

	retry := 20
	for retry > 0 {
		recv, err := i.sendrecv(req)
		if err != nil {
			return err
		}
		// response data byte 3 bit[1:0] == 00b indicates set complete
		if len(recv) == 3 && (recv[2]&0x3) == 0 {
			return nil
		}
		time.Sleep(1 * time.Millisecond)
		retry--
	}

	return fmt.Errorf("Wait for Set in Progress cleared timeout")

}

// SetSystemFWVersion sets the provided system firmware version to BMC via IPMI.
func (i *IPMI) SetSystemFWVersion(version string) error {
	strlenMax := 64 // Set 64 Bytes as the maximal version string length
	len := len(version)

	if len > strlenMax || len == 0 {
		return fmt.Errorf("Version length is 0 or longer than the suggested maximal length %d", strlenMax)
	}

	req := &req{}
	req.msg.cmd = _SET_SYSTEM_INFO_PARAMETERS
	req.msg.netfn = _IPMI_NETFN_APP
	var data setSystemInfoReq
	var copied int
	data.paramSelector = _SYSTEM_FW_VERSION
	data.setSelector = 0
	index := 0
	for len > 0 {
		if data.setSelector == 0 { // the fisrt block of string data
			data.version[0] = _ASCII
			data.version[1] = byte(len)
			copied = copy(data.version[2:], version)
			// dataLen needs to add the actual copied string data plus the first 2 bytes
			req.msg.dataLen = uint16(int(unsafe.Sizeof(data.paramSelector)) + int(unsafe.Sizeof(data.setSelector)) + copied + 2)
		} else {
			copied = copy(data.version[:], version[index:])
			req.msg.dataLen = uint16(int(unsafe.Sizeof(data.paramSelector)) + int(unsafe.Sizeof(data.setSelector)) + copied)
		}
		index += copied
		req.msg.data = unsafe.Pointer(&data)
		if err := i.waitForSetInProgressCleared(); err != nil {
			return err
		}

		if _, err := i.sendrecv(req); err != nil {
			return err
		}

		len -= copied
		data.setSelector++
		for j := range data.version { //reset to 0
			data.version[j] = 0
		}
	}

	return nil
}
