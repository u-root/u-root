// Copyright 2019-2022 the u-root Authors. All rights reserved
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
	"math/rand"
	"syscall"
	"time"
	"unsafe"

	"github.com/vtolstov/go-ioctl"
)

var (
	_IPMICTL_RECEIVE_MSG_TRUNC = ioctl.IOWR(_IPMI_IOC_MAGIC, 11, uintptr(unsafe.Sizeof(response{})))
	_IPMICTL_SEND_COMMAND      = ioctl.IOR(_IPMI_IOC_MAGIC, 13, uintptr(unsafe.Sizeof(request{})))

	timeout = time.Second * 10
)

// IPMI represents access to the IPMI interface.
type IPMI struct {
	*dev
}

// SendRecv sends the IPMI message, receives the response, and returns the
// response data. This is recommended for use unless the user must be able to
// specify the data pointer and length on their own.
func (i *IPMI) SendRecv(netfn NetFn, cmd Command, data []byte) ([]byte, error) {
	var dataPtr unsafe.Pointer
	if data != nil {
		dataPtr = unsafe.Pointer(&data[0])
	}
	msg := Msg{
		Netfn:   netfn,
		Cmd:     cmd,
		Data:    dataPtr,
		DataLen: uint16(len(data)),
	}
	return i.RawSendRecv(msg)
}

// RawSendRecv sends the IPMI message, receives the response, and returns the
// response data.
func (i *IPMI) RawSendRecv(msg Msg) ([]byte, error) {
	addr := &systemInterfaceAddr{
		addrType: _IPMI_SYSTEM_INTERFACE_ADDR_TYPE,
		channel:  _IPMI_BMC_CHANNEL,
	}
	req := &request{
		addr:    addr,
		addrLen: uint32(unsafe.Sizeof(addr)),
		msgid:   rand.Int63(),
		msg:     msg,
	}

	// Send request.
	for {
		switch err := i.dev.SendRequest(req); {
		case err == syscall.EINTR:
			continue
		case err != nil:
			return nil, fmt.Errorf("ioctlSetReq failed with %v", err)
		}
		break
	}

	buf := make([]byte, _IPMI_BUF_SIZE)
	recvMsg := Msg{
		Data:    unsafe.Pointer(&buf[0]),
		DataLen: _IPMI_BUF_SIZE,
	}
	recv := &response{
		addr:    req.addr,
		addrLen: req.addrLen,
		msg:     recvMsg,
	}

	return i.dev.ReceiveResponse(req.msgid, recv, buf)
}

func (i *IPMI) WatchdogRunning() (bool, error) {
	recv, err := i.SendRecv(_IPMI_NETFN_APP, BMC_GET_WATCHDOG_TIMER, nil)
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

	_, err := i.SendRecv(_IPMI_NETFN_APP, BMC_SET_WATCHDOG_TIMER, data[:6])
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
	_, err = i.SendRecv(_IPMI_NETFN_STORAGE, BMC_ADD_SEL, data)
	return err
}

func (i *IPMI) setsysinfo(data *setSystemInfoReq) error {
	msg := Msg{
		Cmd:     SET_SYSTEM_INFO_PARAMETERS,
		Netfn:   _IPMI_NETFN_APP,
		Data:    unsafe.Pointer(data),
		DataLen: 18,
	}

	if _, err := i.RawSendRecv(msg); err != nil {
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
	length := len(version)

	if length == 0 {
		return fmt.Errorf("version length is 0")
	} else if length > strlenMax {
		length = strlenMax
	}

	var data setSystemInfoReq
	var index int
	data.paramSelector = _SYSTEM_FW_VERSION
	data.setSelector = 0
	for length > index {
		if data.setSelector == 0 { // the fisrt block of string data
			data.strData[0] = _ASCII
			data.strData[1] = byte(length)
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
	data, err := i.SendRecv(_IPMI_NETFN_APP, BMC_GET_DEVICE_ID, nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(data[1:])
	mcInfo := DevID{}

	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.DeviceID); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.DeviceRevision); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.FwRev1); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.FwRev2); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.IpmiVersion); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.AdtlDeviceSupport); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.ManufacturerID); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &mcInfo.ProductID); err != nil {
		return nil, err
	}
	// In some cases we have 11 bytes, in others we may have 15 bytes. Carefully parsing the struct is important here.
	if buf.Len() > 0 {
		if err := binary.Read(buf, binary.LittleEndian, &mcInfo.AuxFwRev); err != nil {
			return nil, err
		}
	}

	return &mcInfo, nil
}

func (i *IPMI) setGlobalEnables(enables byte) error {
	_, err := i.SendRecv(_IPMI_NETFN_APP, BMC_SET_GLOBAL_ENABLES, []byte{enables})
	return err
}

func (i *IPMI) getGlobalEnables() ([]byte, error) {
	return i.SendRecv(_IPMI_NETFN_APP, BMC_GET_GLOBAL_ENABLES, nil)
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
	data, err := i.SendRecv(_IPMI_NETFN_CHASSIS, BMC_GET_CHASSIS_STATUS, nil)
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
	data, err := i.SendRecv(_IPMI_NETFN_STORAGE, BMC_GET_SEL_INFO, nil)
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

	return i.SendRecv(_IPMI_NETFN_TRANSPORT, BMC_GET_LAN_CONFIG, data[:4])
}

func (i *IPMI) RawCmd(param []byte) ([]byte, error) {
	if len(param) < 2 {
		return nil, errors.New("not enough parameters given")
	}

	msg := Msg{
		Netfn: NetFn(param[0]),
		Cmd:   Command(param[1]),
	}
	if len(param) > 2 {
		msg.Data = unsafe.Pointer(&param[2])
	}

	msg.DataLen = uint16(len(param) - 2)

	return i.RawSendRecv(msg)
}

// Close closes the file attached to ipmi
func (i *IPMI) Close() error {
	return i.dev.Close()
}
