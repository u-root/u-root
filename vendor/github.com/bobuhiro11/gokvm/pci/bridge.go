package pci

import "errors"

var ErrIONotPermit = errors.New("IO is not permitted for PCI bridge")

type bridge struct{}

func (br bridge) GetDeviceHeader() DeviceHeader {
	return DeviceHeader{
		DeviceID:      0x0d57,
		VendorID:      0x8086,
		HeaderType:    1,
		SubsystemID:   0,
		InterruptLine: 0,
		InterruptPin:  0,
		BAR:           [6]uint32{},
		Command:       0,
	}
}

func (br bridge) Read(port uint64, bytes []byte) error {
	return ErrIONotPermit
}

func (br bridge) Write(port uint64, bytes []byte) error {
	return ErrIONotPermit
}

func (br bridge) IOPort() uint64 {
	return 0
}

func (br bridge) Size() uint64 {
	return 0x10
}

func NewBridge() Device {
	return &bridge{}
}
