// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipmi implements functions to communicate with
// the OpenIPMI driver interface.
package ipmi

import (
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

	_BMC_SET_WATCHDOG_TIMER = 0x24
	_BMC_GET_WATCHDOG_TIMER = 0x25

	_IPM_WATCHDOG_NO_ACTION    = 0x00
	_IPM_WATCHDOG_SMS_OS       = 0x04
	_IPM_WATCHDOG_CLEAR_SMS_OS = 0x10
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
