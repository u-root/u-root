package iodev

import (
	"time"
)

const (
	indexMask   = uint8(0x7F)
	indexOffset = uint64(0x70)
	dataOffset  = uint64(0x71)
	dataLen     = uint64(128)
)

type CMOS struct {
	Index uint8
	Data  []uint8
}

func NewCMOS(memBelow4G, memAbove4G uint64) *CMOS {
	cmos := &CMOS{
		Index: 0,
		Data:  make([]uint8, dataLen),
	}

	// We assume 3G RAM at all times for now.....
	extMem := uint16(0xBC00)

	cmos.Data[0x34] = uint8(extMem)
	cmos.Data[0x35] = uint8(extMem >> 8)

	// Only valid for PVH boot of firmware with 3G RAM fixed.....
	cmos.Data[0x5b] = 0
	cmos.Data[0x5c] = 0
	cmos.Data[0x5d] = 0

	return cmos
}

func (c *CMOS) Read(base uint64, data []byte) error {
	if len(data) != 1 {
		return errDataLenInvalid
	}

	var d uint8

	// Reading CMOS RAM also requires two steps:
	// 1. OUT to port hex 70 with the CMOS address that is to be read from.
	// 2. IN from port hex 71, and the data read is returned in the AL register.
	// Ref: http://bitsavers.trailing-edge.com/pdf/ibm/pc/at/1502494_PC_AT_Technical_Reference_Mar84.pdf
	switch base {
	case indexOffset:
		data[0] = c.Index
	case dataOffset:
		dt := time.Now()
		secs := dt.Second()
		min := dt.Minute()
		hour := dt.Hour()
		weekd := dt.Weekday()
		day := dt.Day()
		month := dt.Month()
		year := dt.Year()

		switch c.Index {
		case 0x00:
			d = toBCD(uint8(secs))
		case 0x02:
			d = toBCD(uint8(min))
		case 0x04:
			d = toBCD(uint8(hour))
		case 0x06:
			d = toBCD(uint8(weekd))
		case 0x07:
			d = toBCD(uint8(day))
		case 0x08:
			d = toBCD(uint8(month))
		case 0x09:
			d = toBCD(uint8(year % 100))
		case 0x0A:
			d = 1<<5 | 0<<7 // 32kHz Clock and we assume no update in progress
		case 0x0D:
			d = 1 << 7
		case 0x32:
			d = toBCD((uint8(year+1900) / 100))
		default:
			d = c.Data[c.Index&indexMask]
		}

		data[0] = d
	}

	return nil
}

func (c *CMOS) Write(base uint64, data []byte) error {
	if len(data) != 1 {
		return errDataLenInvalid
	}

	// Writing to CMOS RAM involves two steps:
	// 1. OUT to port hex 70 with the CMOS address that will be written to.
	// 2. OUT to port hex 71 with the data to be written.
	// Ref: http://bitsavers.trailing-edge.com/pdf/ibm/pc/at/1502494_PC_AT_Technical_Reference_Mar84.pdf
	switch base {
	case indexOffset:
		c.Index = data[0]
	case dataOffset:
		if c.Index == 0x8F && data[0] == 0 {
			// CMOS reset - we ignore for now
		} else {
			c.Data[c.Index&indexMask] = data[0]
		}
	}

	return nil
}

func toBCD(v uint8) uint8 {
	return ((v / 100) << 4) | (v % 10)
}

func (c *CMOS) IOPort() uint64 {
	return 0x70
}

func (c *CMOS) Size() uint64 {
	return 0x2
}
