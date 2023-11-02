package iodev

import (
	"encoding/binary"
	"time"
)

type ACPIPMTimer struct {
	Start time.Time
}

const (
	pmTimerFreqHz  uint64 = 3_579_545
	nanosPerSecond uint64 = 1_000_000_000
)

func NewACPIPMTimer() *ACPIPMTimer {
	return &ACPIPMTimer{
		Start: time.Now(),
	}
}

func (a *ACPIPMTimer) Read(base uint64, data []byte) error {
	if len(data) != 4 {
		return errDataLenInvalid
	}

	since := time.Since(a.Start)
	nanos := since.Nanoseconds()
	counter := (nanos * int64(pmTimerFreqHz)) / int64(nanosPerSecond)
	counter32 := uint32(counter & 0xFFFF_FFFF)
	counterbyte := make([]byte, 4)

	binary.LittleEndian.PutUint32(counterbyte, counter32)

	copy(data[0:], counterbyte)

	return nil
}

func (a *ACPIPMTimer) Write(base uint64, data []byte) error {
	return nil
}

func (a *ACPIPMTimer) IOPort() uint64 {
	return 0x608
}

func (a *ACPIPMTimer) Size() uint64 {
	return 0x4
}
