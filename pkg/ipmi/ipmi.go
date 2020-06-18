// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipmi implements functions to communicate with the OpenIPMI driver
// interface.
// For a detailed description of OpenIPMI, see
// http://openipmi.sourceforge.net/IPMI.pdf
package ipmi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/vtolstov/go-ioctl"
	"golang.org/x/sys/unix"
)

const (
	_IPMI_BMC_CHANNEL                = 0xf
	_IPMI_BUF_SIZE                   = 1024
	_IPMI_IOC_MAGIC                  = 'i'
	_IPMI_OPENIPMI_READ_TIMEOUT      = 15
	_IPMI_SYSTEM_INTERFACE_ADDR_TYPE = 0x0c

	// Net functions
	_IPMI_NETFN_CHASSIS   NetFn = 0x0
	_IPMI_NETFN_APP       NetFn = 0x6
	_IPMI_NETFN_STORAGE   NetFn = 0xA
	_IPMI_NETFN_TRANSPORT NetFn = 0xC

	// IPM Device "Global" Commands
	BMC_GET_DEVICE_ID Command = 0x01

	// BMC Device and Messaging Commands
	BMC_SET_WATCHDOG_TIMER     Command = 0x24
	BMC_GET_WATCHDOG_TIMER     Command = 0x25
	BMC_SET_GLOBAL_ENABLES     Command = 0x2E
	BMC_GET_GLOBAL_ENABLES     Command = 0x2F
	SET_SYSTEM_INFO_PARAMETERS Command = 0x58
	BMC_ADD_SEL                Command = 0x44

	// Chassis Device Commands
	BMC_GET_CHASSIS_STATUS Command = 0x01

	// SEL device Commands
	BMC_GET_SEL_INFO Command = 0x40

	//LAN Device Commands
	BMC_GET_LAN_CONFIG Command = 0x02

	IPM_WATCHDOG_NO_ACTION    = 0x00
	IPM_WATCHDOG_SMS_OS       = 0x04
	IPM_WATCHDOG_CLEAR_SMS_OS = 0x10

	ADTL_SEL_DEVICE         = 0x04
	EN_SYSTEM_EVENT_LOGGING = 0x08

	// SEL
	// STD_TYPE  = 0x02
	OEM_NTS_TYPE = 0xFB

	_SYSTEM_INFO_BLK_SZ = 16

	_SYSTEM_FW_VERSION = 1

	_ASCII = 0

	// Set 62 Bytes (4 sets) as the maximal string length
	strlenMax = 62
)

var (
	_IPMICTL_RECEIVE_MSG  = ioctl.IOWR(_IPMI_IOC_MAGIC, 12, uintptr(unsafe.Sizeof(recv{})))
	_IPMICTL_SEND_COMMAND = ioctl.IOR(_IPMI_IOC_MAGIC, 13, uintptr(unsafe.Sizeof(req{})))

	timeout = 30 * time.Second
)

// IPMI represents access to the IPMI interface.
type IPMI struct {
	*os.File
}

// Command is the command code for a given message.
type Command byte

// NetFn is the network function of the class of message being sent.
type NetFn byte

// Msg is the full IPMI message to be sent.
type Msg struct {
	Netfn   NetFn
	Cmd     Command
	DataLen uint16
	Data    unsafe.Pointer
}

type req struct {
	addr    *systemInterfaceAddr
	addrLen uint32
	msgid   int64 //nolint:structcheck
	msg     Msg
}

type recv struct {
	recvType int32 //nolint:structcheck
	addr     *systemInterfaceAddr
	addrLen  uint32
	msgid    int64 //nolint:structcheck
	msg      Msg
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

type ChassisStatus struct {
	CurrentPowerState byte
	LastPowerEvent    byte
	MiscChassisState  byte
	FrontPanelButton  byte
}

type SELInfo struct {
	Version     byte
	Entries     uint16
	FreeSpace   uint16
	LastAddTime uint32
	LastDelTime uint32
	OpSupport   byte
}

func fdSet(fd uintptr, p *syscall.FdSet) {
	p.Bits[fd/64] |= 1 << (uint(fd) % 64)
}

func ioctlSetReq(fd, name uintptr, req *req) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, name, uintptr(unsafe.Pointer(req)))
	runtime.KeepAlive(req)
	if err != 0 {
		return err
	}
	return nil
}

func ioctlGetRecv(fd, name uintptr, recv *recv) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, name, uintptr(unsafe.Pointer(recv)))
	runtime.KeepAlive(recv)
	if err != 0 {
		return err
	}
	return nil
}

// SendRecvBasic sends the IPMI message, receives the response, and returns the
// response data. This is recommended for use unless the user must be able to
// specify the data pointer and length on their own.
func (i *IPMI) SendRecvBasic(netfn NetFn, cmd Command, data []byte) ([]byte, error) {
	msg := Msg{
		Netfn:   netfn,
		Cmd:     cmd,
		Data:    unsafe.Pointer(&data[0]),
		DataLen: uint16(len(data)),
	}
	return i.SendRecv(msg)
}

// SendRecv sends the IPMI message, receives the response, and returns the
// response data.
func (i *IPMI) SendRecv(msg Msg) ([]byte, error) {
	addr := &systemInterfaceAddr{
		addrType: _IPMI_SYSTEM_INTERFACE_ADDR_TYPE,
		channel:  _IPMI_BMC_CHANNEL,
	}
	req := &req{
		addr:    addr,
		addrLen: uint32(unsafe.Sizeof(addr)),
		msgid:   rand.Int63(),
		msg:     msg,
	}

	if err := ioctlSetReq(i.Fd(), _IPMICTL_SEND_COMMAND, req); err != nil {
		return nil, err
	}

	set := &syscall.FdSet{}
	fdSet(i.Fd(), set)
	timeval := &syscall.Timeval{
		Sec:  _IPMI_OPENIPMI_READ_TIMEOUT,
		Usec: 0,
	}
	if _, err := syscall.Select(int(i.Fd()+1), set, nil, nil, timeval); err != nil {
		return nil, err
	}

	buf := make([]byte, _IPMI_BUF_SIZE)
	recvMsg := Msg{
		Data:    unsafe.Pointer(&buf[0]),
		DataLen: _IPMI_BUF_SIZE,
	}
	recv := &recv{
		addr:    req.addr,
		addrLen: req.addrLen,
		msg:     recvMsg,
	}

	t := time.After(timeout)
	for {
		if err := ioctlGetRecv(i.Fd(), _IPMICTL_RECEIVE_MSG, recv); err != nil {
			return nil, err
		}
		if recv.msgid != req.msgid {
			log.Printf("Received wrong message. Trying again.")
		} else {
			break
		}
		select {
		case <-t:
			return nil, fmt.Errorf("timeout waiting for response")
		}
	}

	return buf[:recv.msg.DataLen:recv.msg.DataLen], nil
}

func (i *IPMI) WatchdogRunning() (bool, error) {
	msg := Msg{
		Cmd:   BMC_GET_WATCHDOG_TIMER,
		Netfn: _IPMI_NETFN_APP,
	}

	recv, err := i.SendRecv(msg)
	if err != nil {
		return false, err
	}

	if len(recv) > 2 && (recv[1]&0x40) != 0 {
		return true, nil
	}

	return false, nil
}

func (i *IPMI) ShutoffWatchdog() error {
	var data [6]byte
	data[0] = IPM_WATCHDOG_SMS_OS
	data[1] = IPM_WATCHDOG_NO_ACTION
	data[2] = 0x00 // pretimeout interval
	data[3] = IPM_WATCHDOG_CLEAR_SMS_OS
	data[4] = 0xb8 // countdown lsb (100 ms/count)
	data[5] = 0x0b // countdown msb - 5 mins

	msg := Msg{
		Cmd:     BMC_SET_WATCHDOG_TIMER,
		Netfn:   _IPMI_NETFN_APP,
		Data:    unsafe.Pointer(&data),
		DataLen: 6,
	}

	_, err := i.SendRecv(msg)
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
	data, err := e.marshall()

	if err != nil {
		return err
	}

	msg := Msg{
		Cmd:     BMC_ADD_SEL,
		Netfn:   _IPMI_NETFN_STORAGE,
		Data:    unsafe.Pointer(&data[0]),
		DataLen: 16,
	}

	_, err = i.SendRecv(msg)
	return err
}

func (i *IPMI) setsysinfo(data *setSystemInfoReq) error {
	msg := Msg{
		Cmd:     SET_SYSTEM_INFO_PARAMETERS,
		Netfn:   _IPMI_NETFN_APP,
		Data:    unsafe.Pointer(data),
		DataLen: 18,
	}

	if _, err := i.SendRecv(msg); err != nil {
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
	msg := Msg{
		Cmd:   BMC_GET_DEVICE_ID,
		Netfn: _IPMI_NETFN_APP,
	}

	data, err := i.SendRecv(msg)

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
	msg := Msg{
		Cmd:     BMC_SET_GLOBAL_ENABLES,
		Netfn:   _IPMI_NETFN_APP,
		Data:    unsafe.Pointer(&enables),
		DataLen: 1,
	}

	_, err := i.SendRecv(msg)
	return err
}

func (i *IPMI) getGlobalEnables() ([]byte, error) {
	msg := Msg{
		Cmd:   BMC_GET_GLOBAL_ENABLES,
		Netfn: _IPMI_NETFN_APP,
	}

	return i.SendRecv(msg)
}

func (i *IPMI) EnableSEL() (bool, error) {
	// Check if SEL device is supported or not
	mcInfo, err := i.GetDeviceID()

	if err != nil {
		return false, err
	} else if (mcInfo.AdtlDeviceSupport & ADTL_SEL_DEVICE) == 0 {
		return false, nil
	}

	data, err := i.getGlobalEnables()

	if err != nil {
		return false, err
	}

	if (data[1] & EN_SYSTEM_EVENT_LOGGING) == 0 {
		// SEL is not enabled, enable SEL
		if err = i.setGlobalEnables(data[1] | EN_SYSTEM_EVENT_LOGGING); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (i *IPMI) GetChassisStatus() (*ChassisStatus, error) {
	msg := Msg{
		Cmd:   BMC_GET_CHASSIS_STATUS,
		Netfn: _IPMI_NETFN_CHASSIS,
	}

	data, err := i.SendRecv(msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data[1:])

	var status ChassisStatus
	if err := binary.Read(buf, binary.LittleEndian, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (i *IPMI) GetSELInfo() (*SELInfo, error) {
	msg := Msg{
		Cmd:   BMC_GET_SEL_INFO,
		Netfn: _IPMI_NETFN_STORAGE,
	}

	data, err := i.SendRecv(msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data[1:])

	var info SELInfo
	if err := binary.Read(buf, binary.LittleEndian, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (i *IPMI) GetLanConfig(channel byte, param byte) ([]byte, error) {
	var data [4]byte
	data[0] = channel
	data[1] = param
	data[2] = 0
	data[3] = 0

	msg := Msg{
		Cmd:     BMC_GET_LAN_CONFIG,
		Netfn:   _IPMI_NETFN_TRANSPORT,
		Data:    unsafe.Pointer(&data[0]),
		DataLen: 4,
	}

	return i.SendRecv(msg)
}

func (i *IPMI) RawCmd(param []byte) ([]byte, error) {
	if len(param) < 2 {
		return nil, errors.New("Not enough parameters given")
	}

	msg := Msg{
		Netfn: NetFn(param[0]),
		Cmd:   Command(param[1]),
	}
	if len(param) > 2 {
		msg.Data = unsafe.Pointer(&param[2])
	}

	msg.DataLen = uint16(len(param) - 2)

	return i.SendRecv(msg)
}
