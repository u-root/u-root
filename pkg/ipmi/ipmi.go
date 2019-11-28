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

	_BMC_SET_WATCHDOG_TIMER     = 0x24
	_BMC_GET_WATCHDOG_TIMER     = 0x25
	_SET_SYSTEM_INFO_PARAMETERS = 0x58

	_IPM_WATCHDOG_NO_ACTION    = 0x00
	_IPM_WATCHDOG_SMS_OS       = 0x04
	_IPM_WATCHDOG_CLEAR_SMS_OS = 0x10

	_SYSTEM_INFO_BLK_SZ = 16

	_SYSTEM_FW_VERSION = 1

	_ASCII = 0

	// Set 62 Bytes (4 sets) as the maximal string length
	strlenMax = 62
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
	strData       [_SYSTEM_INFO_BLK_SZ]byte
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

func (i *IPMI) setsysinfo(data *setSystemInfoReq) error {
	req := &req{}
	req.msg.cmd = _SET_SYSTEM_INFO_PARAMETERS
	req.msg.netfn = _IPMI_NETFN_APP
	req.msg.dataLen = 18 // size of setSystemInfoReq
	req.msg.data = unsafe.Pointer(data)

	if _, err := i.sendrecv(req); err != nil {
		return err
	}

	return nil
}

func strcpyPadded(dst []byte, src string) {
	dstLen := len(dst)
	if copied := copy(dst, src); copied < dstLen {
		padding := make([]byte, dstLen-copied)
		copy(dst[copied:], padding)
	}
}

// SetSystemFWVersion sets the provided system firmware version to BMC via IPMI.
func (i *IPMI) SetSystemFWVersion(version string) error {
	len := len(version)

	if len == 0 {
		return fmt.Errorf("Version length is 0")
	} else if len > strlenMax {
		len = strlenMax
	}

	var data setSystemInfoReq
	var index int
	data.paramSelector = _SYSTEM_FW_VERSION
	data.setSelector = 0
	for len > index {
		if data.setSelector == 0 { // the fisrt block of string data
			data.strData[0] = _ASCII
			data.strData[1] = byte(len)
			strcpyPadded(data.strData[2:], version)
			index += _SYSTEM_INFO_BLK_SZ - 2
		} else {
			strcpyPadded(data.strData[:], version)
			index += _SYSTEM_INFO_BLK_SZ
		}

		if err := i.setsysinfo(&data); err != nil {
			return err
		}
		data.setSelector++
	}

	return nil
}
