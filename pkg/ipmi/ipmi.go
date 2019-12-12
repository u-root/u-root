// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipmi implements functions to communicate with
// the OpenIPMI driver interface.
package ipmi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	_IPMI_BMC_CHANNEL                = 0xf
	_IPMI_BUF_SIZE                   = 1024
	_IPMI_IOC_MAGIC                  = 'i'
	_IPMI_NETFN_APP                  = 0x6
	_IPMI_NETFN_STORAGE              = 0xA
	_IPMI_OPENIPMI_READ_TIMEOUT      = 15
	_IPMI_SYSTEM_INTERFACE_ADDR_TYPE = 0x0c

	// IPM Device "Global" Commands
	_BMC_GET_DEVICE_ID = 0x01

	// BMC Device and Messaging Commands
	_BMC_SET_WATCHDOG_TIMER     = 0x24
	_BMC_GET_WATCHDOG_TIMER     = 0x25
	_BMC_SET_GLOBAL_ENABLES     = 0x2E
	_BMC_GET_GLOBAL_ENABLES     = 0x2F
	_SET_SYSTEM_INFO_PARAMETERS = 0x58
	_BMC_ADD_SEL                = 0x44

	_IPM_WATCHDOG_NO_ACTION    = 0x00
	_IPM_WATCHDOG_SMS_OS       = 0x04
	_IPM_WATCHDOG_CLEAR_SMS_OS = 0x10

	_ADTL_SEL_DEVICE         = 0x04
	_EN_SYSTEM_EVENT_LOGGING = 0x08

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

// StandardEvent is a standard systemevent.
//
// The data in this event should follow IPMI spec
type StandardEvent struct {
	Timestamp    uint32
	GenID        uint16
	EvMRev       uint8
	SensorType   uint8
	SensorNum    uint8
	EventTypeDir uint8
	EventData    [3]uint8
}

// OEMTsEvent is a timestamped OEM-custom event.
//
// It holds 6 bytes of OEM-defined arbitrary data.
type OEMTsEvent struct {
	Timestamp        uint32
	ManfID           [3]uint8
	OEMTsDefinedData [6]uint8
}

// OEMNonTsEvent is a non-timestamped OEM-custom event.
//
// It holds 13 bytes of OEM-defined arbitrary data.
type OEMNontsEvent struct {
	OEMNontsDefinedData [13]uint8
}

// Event is included three kinds of events, Standard, OEM timestamped and OEM non-timestamped
//
// The record type decides which event should be used
type Event struct {
	RecordID   uint16
	RecordType uint8
	StandardEvent
	OEMTsEvent
	OEMNontsEvent
}

type setSystemInfoReq struct {
	paramSelector byte
	setSelector   byte
	strData       [_SYSTEM_INFO_BLK_SZ]byte
}

type DevID struct {
	DeviceID          byte
	DeviceRevision    byte
	FwRev1            byte
	FwRev2            byte
	IpmiVersion       byte
	AdtlDeviceSupport byte
	ManufacturerID    [3]byte
	ProductID         [2]byte
	AuxFwRev          [4]byte
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

// marshall converts the Event struct to binary data and the content of returned data is based on the record type
func (e *Event) marshall() ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, binary.LittleEndian, *e); err != nil {
		return nil, err
	}

	data := make([]byte, 16)

	// system event record
	if buf.Bytes()[2] == 0x2 {
		copy(data[:], buf.Bytes()[:16])
	}

	// OEM timestamped
	if buf.Bytes()[2] >= 0xC0 && buf.Bytes()[2] <= 0xDF {
		copy(data[0:3], buf.Bytes()[0:3])
		copy(data[3:16], buf.Bytes()[16:29])
	}

	// OEM non-timestamped
	if buf.Bytes()[2] >= 0xE0 && buf.Bytes()[2] <= 0xFF {
		copy(data[0:3], buf.Bytes()[0:3])
		copy(data[3:16], buf.Bytes()[29:42])
	}

	return data, nil
}

// LogSystemEvent adds an SEL (System Event Log) entry.
func (i *IPMI) LogSystemEvent(e *Event) error {
	req := &req{}
	req.msg.cmd = _BMC_ADD_SEL
	req.msg.netfn = _IPMI_NETFN_STORAGE

	data, err := e.marshall()

	if err != nil {
		return err
	}

	req.msg.data = unsafe.Pointer(&data[0])
	req.msg.dataLen = 16

	_, err = i.sendrecv(req)

	return err
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

func (i *IPMI) GetDeviceID() (*DevID, error) {
	req := &req{}
	req.msg.netfn = _IPMI_NETFN_APP
	req.msg.cmd = _BMC_GET_DEVICE_ID

	data, err := i.sendrecv(req)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data[1:])
	mcInfo := DevID{}

	if err := binary.Read(buf, binary.LittleEndian, &mcInfo); err != nil {
		return nil, err
	}

	return &mcInfo, nil
}

func (i *IPMI) setGlobalEnables(enables byte) error {
	req := &req{}
	req.msg.netfn = _IPMI_NETFN_APP
	req.msg.cmd = _BMC_SET_GLOBAL_ENABLES
	req.msg.data = unsafe.Pointer(&enables)
	req.msg.dataLen = 1

	_, err := i.sendrecv(req)
	return err
}

func (i *IPMI) getGlobalEnables() ([]byte, error) {
	req := &req{}
	req.msg.netfn = _IPMI_NETFN_APP
	req.msg.cmd = _BMC_GET_GLOBAL_ENABLES

	return i.sendrecv(req)
}

func (i *IPMI) EnableSEL() (bool, error) {
	// Check if SEL device is supported or not
	mcInfo, err := i.GetDeviceID()

	if err != nil {
		return false, err
	} else if (mcInfo.AdtlDeviceSupport & _ADTL_SEL_DEVICE) == 0 {
		return false, nil
	}

	data, err := i.getGlobalEnables()

	if err != nil {
		return false, err
	}

	if (data[1] & _EN_SYSTEM_EVENT_LOGGING) == 0 {
		// SEL is not enabled, enable SEL
		if err = i.setGlobalEnables(data[1] | _EN_SYSTEM_EVENT_LOGGING); err != nil {
			return false, err
		}
	}

	return true, nil
}
