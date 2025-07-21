// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocp implements OCP/Facebook-specific IPMI client functions.
package ocp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/pci"
	"github.com/u-root/u-root/pkg/smbios"
)

const (
	_IPMI_FB_OEM_NET_FUNCTION1 ipmi.NetFn = 0x30
	_IPMI_FB_OEM_NET_FUNCTION2 ipmi.NetFn = 0x36

	_FB_OEM_SET_PROC_INFO       ipmi.Command = 0x10
	_FB_OEM_SET_DIMM_INFO       ipmi.Command = 0x12
	_FB_OEM_SET_BOOT_DRIVE_INFO ipmi.Command = 0x14
	_FB_OEM_SET_BIOS_BOOT_ORDER ipmi.Command = 0x52
	_FB_OEM_GET_BIOS_BOOT_ORDER ipmi.Command = 0x53
	_FB_OEM_SET_POST_END        ipmi.Command = 0x74
)

type ProcessorInfo struct {
	ManufacturerID        [3]uint8
	Index                 uint8
	ParameterSelector     uint8
	ProductName           [48]byte
	CoreNumber            uint8
	ThreadNumberLSB       uint8
	ThreadNumberMSB       uint8
	ProcessorFrequencyLSB uint8
	ProcessorFrequencyMSB uint8
	Revision1             uint8
	Revision2             uint8
}

type DimmInfo struct {
	ManufacturerID          [3]uint8
	Index                   uint8
	ParameterSelector       uint8
	DIMMPresent             uint8
	NodeNumber              uint8
	ChannelNumber           uint8
	DIMMNumber              uint8
	DIMMType                uint8
	DIMMSpeed               uint16
	DIMMSize                uint32
	ModulePartNumber        [20]byte
	ModuleSerialNumber      uint32
	ModuleManufacturerIDLSB uint8
	ModuleManufacturerIDMSB uint8
}

type BootDriveInfo struct {
	ManufacturerID    [3]uint8
	ControlType       uint8
	DriveNumber       uint8
	ParameterSelector uint8
	VendorID          uint16
	DeviceID          uint16
}

// OENMap maps OEM names to a 3 byte OEM number.
//
// OENs are typically serialized as the first 3 bytes of a request body.
var OENMap = map[string][3]uint8{
	"Wiwynn": {0x0, 0x9c, 0x9c},
}

func SendOemIpmiProcessorInfo(i *ipmi.IPMI, info []ProcessorInfo) error {
	for index := 0; index < len(info); index++ {
		for param := 1; param <= 2; param++ {
			data, err := info[index].marshall(param)
			if err != nil {
				return err
			}

			_, err = i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION2, _FB_OEM_SET_PROC_INFO, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SendOemIpmiDimmInfo(i *ipmi.IPMI, info []DimmInfo) error {
	for index := 0; index < len(info); index++ {
		for param := 1; param <= 6; param++ {
			// If DIMM is not present, only send the information of DIMM location
			if info[index].DIMMPresent != 0x01 && param >= 2 {
				continue
			}

			data, err := info[index].marshall(param)
			if err != nil {
				return err
			}
			_, err = i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION2, _FB_OEM_SET_DIMM_INFO, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SendOemIpmiBootDriveInfo(i *ipmi.IPMI, info *BootDriveInfo) error {
	var data []byte
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, *info)
	if err != nil {
		return err
	}

	data = make([]byte, 10)
	copy(data, buf.Bytes())

	_, err = i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION2, _FB_OEM_SET_BOOT_DRIVE_INFO, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProcessorInfo) marshall(param int) ([]byte, error) {
	var data []byte
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, binary.LittleEndian, *p); err != nil {
		return nil, err
	}

	buf.Bytes()[4] = byte(param)

	switch param {
	case 1:
		data = make([]byte, 53)
		copy(data[:], buf.Bytes()[:53])
	case 2:
		data = make([]byte, 12)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:12], buf.Bytes()[53:60])
	}

	return data, nil
}

func (d *DimmInfo) marshall(param int) ([]byte, error) {
	var data []byte
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, binary.LittleEndian, *d); err != nil {
		return nil, err
	}

	buf.Bytes()[4] = byte(param)

	switch param {
	case 1:
		data = make([]byte, 9)
		copy(data[:], buf.Bytes()[:9])
	case 2:
		data = make([]byte, 6)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:6], buf.Bytes()[9:10])
	case 3:
		data = make([]byte, 11)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:11], buf.Bytes()[10:16])
	case 4:
		data = make([]byte, 25)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:25], buf.Bytes()[16:36])
	case 5:
		data = make([]byte, 9)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:9], buf.Bytes()[36:40])
	case 6:
		data = make([]byte, 7)
		copy(data[0:5], buf.Bytes()[0:5])
		copy(data[5:7], buf.Bytes()[40:42])
	}

	return data, nil
}

func GetOemIpmiProcessorInfo(si *smbios.Info) ([]ProcessorInfo, error) {
	t1, err := si.GetSystemInfo()
	if err != nil {
		return nil, err
	}

	t4, err := si.GetProcessorInfo()
	if err != nil {
		return nil, err
	}

	info := make([]ProcessorInfo, len(t4))

	boardManufacturerID, ok := OENMap[t1.Manufacturer]

	for index := 0; index < len(t4); index++ {
		if ok {
			info[index].ManufacturerID = boardManufacturerID
		}

		info[index].Index = uint8(index)
		copy(info[index].ProductName[:], t4[index].Version)
		info[index].CoreNumber = uint8(t4[index].GetCoreCount())
		info[index].ThreadNumberLSB = uint8(t4[index].GetThreadCount() & 0x00ff)
		info[index].ThreadNumberMSB = uint8(t4[index].GetThreadCount() >> 8)
		info[index].ProcessorFrequencyLSB = uint8(t4[index].CurrentSpeed & 0x00ff)
		info[index].ProcessorFrequencyMSB = uint8(t4[index].CurrentSpeed >> 8)
		info[index].Revision1 = 0
		info[index].Revision2 = 0
	}

	return info, nil
}

// DIMM type: bit[7:6] for DDR3 00-Normal Voltage(1.5V), 01-Ultra Low Voltage(1.25V), 10-Low Voltage(1.35V), 11-Reserved
//
//	                    for DDR4 00~10-Reserved, 11-Normal Voltage(1.2V)
//	           bit[5:0] 0x00=SDRAM, 0x01=DDR1 RAM, 0x02-Rambus, 0x03-DDR2 RAM, 0x04-FBDIMM, 0x05-DDR3 RAM, 0x06-DDR4 RAM
//			       , 0x07-DDR5 RAM
func detectDimmType(meminfo *DimmInfo, t17 *smbios.MemoryDevice) {
	if t17.Type == smbios.MemoryDeviceTypeDDR3 {
		switch t17.ConfiguredVoltage {
		case 1500:
			meminfo.DIMMType = 0x05
		case 1250:
			meminfo.DIMMType = 0x45
		case 1350:
			meminfo.DIMMType = 0x85
		default:
			meminfo.DIMMType = 0x05
		}
	} else {
		switch t17.Type {
		case smbios.MemoryDeviceTypeSDRAM:
			meminfo.DIMMType = 0x00
		case smbios.MemoryDeviceTypeDDR:
			meminfo.DIMMType = 0x01
		case smbios.MemoryDeviceTypeRDRAM:
			meminfo.DIMMType = 0x02
		case smbios.MemoryDeviceTypeDDR2:
			meminfo.DIMMType = 0x03
		case smbios.MemoryDeviceTypeDDR2FBDIMM:
			meminfo.DIMMType = 0x04
		case smbios.MemoryDeviceTypeDDR4:
			meminfo.DIMMType = 0xC6
		case smbios.MemoryDeviceTypeDDR5:
			meminfo.DIMMType = 0xC7
		default:
			meminfo.DIMMType = 0xC6
		}
	}
}

func GetOemIpmiDimmInfo(si *smbios.Info) ([]DimmInfo, error) {
	t1, err := si.GetSystemInfo()
	if err != nil {
		return nil, err
	}

	t17, err := si.GetMemoryDevices()
	if err != nil {
		return nil, err
	}

	info := make([]DimmInfo, len(t17))

	boardManufacturerID, ok := OENMap[t1.Manufacturer]

	for index := 0; index < len(t17); index++ {
		if ok {
			info[index].ManufacturerID = boardManufacturerID
		}

		info[index].Index = uint8(index)

		if t17[index].AssetTag == "NO DIMM" {
			info[index].DIMMPresent = 0xFF // 0xFF - Not Present
		} else {
			info[index].DIMMPresent = 0x01 // 0x01 - Present
		}

		data := strings.Split(strings.TrimPrefix(t17[index].BankLocator, "_"), "_")
		dimm, _ := strconv.ParseUint(strings.TrimPrefix(data[2], "Dimm"), 16, 8)
		channel, _ := strconv.ParseUint(strings.TrimPrefix(data[1], "Channel"), 16, 8)
		node, _ := strconv.ParseUint(strings.TrimPrefix(data[0], "Node"), 16, 8)
		info[index].DIMMNumber = uint8(dimm)
		info[index].ChannelNumber = uint8(channel)
		info[index].NodeNumber = uint8(node)
		detectDimmType(&info[index], t17[index])
		info[index].DIMMSpeed = t17[index].Speed
		info[index].DIMMSize = uint32(t17[index].Size)
		copy(info[index].ModulePartNumber[:], t17[index].PartNumber)
		sn, _ := strconv.ParseUint(t17[index].SerialNumber, 16, 32)
		info[index].ModuleSerialNumber = uint32(sn)
		memoryDeviceManufacturerID, ok := smbios.MemoryDeviceManufacturer[t17[index].Manufacturer]
		if ok {
			info[index].ModuleManufacturerIDLSB = uint8(memoryDeviceManufacturerID & 0x00ff)
			info[index].ModuleManufacturerIDMSB = uint8(memoryDeviceManufacturerID >> 8)
		}
	}

	return info, nil
}

// Read type 9 from SMBIOS tables and look for SlotDesignation which contains string 'Boot_Drive',
// and get bus and device number from the matched table to read the Device ID and Vendor ID of
// the boot drive for sending the IPMI OEM command.
// This requires the BDF number is correctly set in the type 9 table.
func GetOemIpmiBootDriveInfo(si *smbios.Info) (*BootDriveInfo, error) {
	t1, err := si.GetSystemInfo()
	if err != nil {
		return nil, err
	}

	t9, err := si.GetSystemSlots()
	if err != nil {
		return nil, err
	}

	const systemPath = "/sys/bus/pci/devices/"
	const bootDriveName = "Boot_Drive"
	var info BootDriveInfo

	if boardManufacturerID, ok := OENMap[t1.Manufacturer]; ok {
		info.ManufacturerID = boardManufacturerID
	}

	for index := 0; index < len(t9); index++ {
		if !strings.Contains(t9[index].SlotDesignation, bootDriveName) {
			continue
		}

		deviceNumber := t9[index].DeviceFunctionNumber >> 3
		devicePath := fmt.Sprintf("%04d:%02x:%02x.0",
			t9[index].SegmentGroupNumber, t9[index].BusNumber, deviceNumber)
		p, err := pci.OnePCI(filepath.Join(systemPath, devicePath))
		if err != nil {
			return nil, err
		}

		info.VendorID = p.Vendor
		info.DeviceID = p.Device
		info.ControlType = 0
		info.DriveNumber = 0
		info.ParameterSelector = 2
		// There will only be one Boot_Drive
		return &info, nil
	}
	return nil, nil
}

func SetOemIpmiPostEnd(i *ipmi.IPMI) error {
	_, err := i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION1, _FB_OEM_SET_POST_END, nil)
	if err != nil {
		return err
	}
	return nil
}

// Get BIOS boot order data and check if CMOS clear bit and valid bit are both set
func IsCMOSClearSet(i *ipmi.IPMI) (bool, []byte, error) {
	recv, err := i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION1, _FB_OEM_GET_BIOS_BOOT_ORDER, nil)
	if err != nil {
		return false, nil, err
	}
	// recv[1] bit 1: CMOS clear, bit 7: valid bit, check if both are set
	if len(recv) > 6 && (recv[1]&0x82) == 0x82 {
		return true, recv[1:], nil
	}
	return false, nil, nil
}

// Set BIOS boot order with both CMOS clear and valid bits cleared
func ClearCMOSClearValidBits(i *ipmi.IPMI, data []byte) error {
	// Clear bit 1 and bit 7
	data[0] &= 0x7d

	msg := ipmi.Msg{
		Netfn:   _IPMI_FB_OEM_NET_FUNCTION1,
		Cmd:     _FB_OEM_SET_BIOS_BOOT_ORDER,
		Data:    unsafe.Pointer(&data[0]),
		DataLen: 6,
	}

	if _, err := i.RawSendRecv(msg); err != nil {
		return err
	}
	return nil
}
