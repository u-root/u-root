package pci

import (
	"bytes"
	"encoding/binary"
)

// Configuration Space Access Mechanism #1
//
// refs
// https://wiki.osdev.org/PCI
// http://www2.comp.ufscar.br/~helio/boot-int/pci.html
type address uint32

func (a address) getRegisterOffset() uint32 {
	return uint32(a) & 0xfc
}

func (a address) getFunctionNumber() uint32 {
	return (uint32(a) >> 8) & 0x7
}

func (a address) getDeviceNumber() uint32 {
	return (uint32(a) >> 11) & 0x1f
}

func (a address) getBusNumber() uint32 {
	return (uint32(a) >> 16) & 0xff
}

func (a address) isEnable() bool {
	return ((uint32(a) >> 31) | 0x1) == 0x1
}

// interface for a PCI device.
type Device interface {
	GetDeviceHeader() DeviceHeader
	Read(uint64, []byte) error
	Write(uint64, []byte) error

	// IO port range for this PCI device.
	// This range corresponds to IO Range in BAR0.
	IOPort() uint64
	Size() uint64
}

type DeviceHeader struct {
	VendorID      uint16
	DeviceID      uint16
	Command       uint16
	_             uint16   // status
	_             uint8    // revisonID
	_             [3]uint8 // classCode
	_             uint8    // cacheLineSize
	_             uint8    // latencyTimer
	HeaderType    uint8
	_             uint8 // bist
	BAR           [6]uint32
	_             uint32 // cardbusCISPointer
	_             uint16 // subsystemVendorID
	SubsystemID   uint16
	_             uint32   // expansionROMBaseAddress
	_             uint8    // capabilitiesPointer
	_             [7]uint8 // reserved
	InterruptLine uint8
	InterruptPin  uint8
	_             uint8 // minGnt
	_             uint8 // maxLat
}

func (h DeviceHeader) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type PCI struct {
	addr        address
	isBAR0Probe bool
	Devices     []Device
}

func New(devices ...Device) *PCI {
	return &PCI{Devices: devices}
}

func (p *PCI) PciConfDataIn(port uint64, values []byte) error {
	// offset can be obtained from many source as below:
	//        (address from IO port 0xcf8) & 0xfc + (IO port address for Data) - 0xCFC
	// see pci_conf1_read in linux/arch/x86/pci/direct.c for more detail.
	offset := int(p.addr.getRegisterOffset() + uint32(port-0xCFC))

	if !p.addr.isEnable() {
		return nil
	}

	if p.addr.getBusNumber() != 0 {
		return nil
	}

	if p.addr.getFunctionNumber() != 0 {
		return nil
	}

	slot := int(p.addr.getDeviceNumber())

	if slot >= len(p.Devices) {
		return nil
	}

	// Probing BAR0 Size
	if bar := offset/4 - 4; bar == 0 && p.isBAR0Probe {
		size := p.Devices[slot].Size()
		copy(values[:4], NumToBytes(SizeToBits(size)))

		p.isBAR0Probe = false

		return nil
	}

	b, err := p.Devices[slot].GetDeviceHeader().Bytes()
	if err != nil {
		return err
	}

	l := len(values)
	copy(values[:l], b[offset:offset+l])

	return nil
}

func (p *PCI) PciConfDataOut(port uint64, values []byte) error {
	offset := int(p.addr.getRegisterOffset() + uint32(port-0xCFC))

	if !p.addr.isEnable() {
		return nil
	}

	if p.addr.getBusNumber() != 0 {
		return nil
	}

	if p.addr.getFunctionNumber() != 0 {
		return nil
	}

	slot := int(p.addr.getDeviceNumber())

	if slot >= len(p.Devices) {
		return nil
	}

	// Probing BAR0 Size
	if bar := offset/4 - 4; bar == 0 && BytesToNum(values) == 0xffffffff {
		p.isBAR0Probe = true

		return nil
	}

	return nil
}

func (p *PCI) PciConfAddrIn(port uint64, values []byte) error {
	if len(values) != 4 {
		return nil
	}

	copy(values[:4], NumToBytes(uint32(p.addr)))

	return nil
}

func (p *PCI) PciConfAddrOut(port uint64, values []byte) error {
	if len(values) != 4 {
		return nil
	}

	p.addr = address(BytesToNum(values))

	return nil
}

func SizeToBits(size uint64) uint32 {
	if size == 0 {
		return 0
	}

	return ^uint32(1) - uint32(size-2)
}

func BytesToNum(bytes []byte) uint64 {
	res := uint64(0)

	for i, x := range bytes {
		res |= uint64(x) << (i * 8)
	}

	return res
}

func NumToBytes(x interface{}) []byte {
	res := []byte{}
	l := 0
	y := uint64(0)

	switch v := x.(type) {
	case uint8:
		l = 1
		y = uint64(v)
	case uint16:
		l = 2
		y = uint64(v)
	case uint32:
		l = 4
		y = uint64(v)
	case uint64:
		l = 8
		y = v
	default:
		return []byte{}
	}

	for i := 0; i < l; i++ {
		res = append(res, uint8(y))
		y >>= 8
	}

	return res
}
