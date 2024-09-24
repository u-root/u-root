package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
)

// iproute2/tc/tc_core.c:tc_calc_rtable()
func generateRateTable(pol *Policy) ([]byte, error) {
	var rate [256]uint32

	if pol == nil {
		return []byte{}, fmt.Errorf("generateRateTable: %w", ErrNoArg)
	}
	mtu := pol.Mtu

	var cellLog int = -1
	if mtu == 0 {
		mtu = 2047
	}

	var linklayer, mpu uint
	var polRate uint64

	if pol.Rate.Rate != 0 {
		linklayer = uint(pol.Rate.Linklayer)
		mpu = uint(pol.Rate.Mpu)
		polRate = uint64(pol.Rate.Rate)
	} else if pol.PeakRate.Rate != 0 {
		linklayer = uint(pol.PeakRate.Linklayer)
		mpu = uint(pol.PeakRate.Mpu)
		polRate = uint64(pol.PeakRate.Rate)
	} else {
		return []byte{}, fmt.Errorf("generateRateTable: Rate or PeakRate is required: %w", ErrNoArg)
	}

	if cellLog < 0 {
		cellLog = 0
		for (mtu >> uint(cellLog)) > 255 {
			cellLog++
		}
	}

	for i := 0; i < 256; i++ {
		sz := adjustSize(uint((i+1)<<uint(cellLog)), mpu, linklayer)
		rate[i] = core.XmitTime(polRate, uint32(sz))
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, nativeEndian, rate)
	return buf.Bytes(), err
}

// iproute2/tc/tc_core.c:tc_adjust_size()
func adjustSize(sz, mpu, linklayer uint) uint32 {
	if sz < mpu {
		sz = mpu
	}

	switch linklayer {
	case unix.LINKLAYER_ATM:
		// iproute2/tc/tc_core.c:tc_align_to_atm()
		var linksize, cells uint

		cells = sz / uint(unix.ATM_CELL_PAYLOAD)
		if (sz % unix.ATM_CELL_PAYLOAD) > 0 {
			cells++
		}

		linksize = cells * unix.ATM_CELL_SIZE
		return uint32(linksize)
	case unix.LINKLAYER_ETHERNET:
		fallthrough
	default:
		return uint32(sz)
	}
}
