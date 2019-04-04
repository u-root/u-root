package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptElapsedTime implements the Elapsed Time option.
//
// This module defines the OptElapsedTime structure.
// https://www.ietf.org/rfc/rfc3315.txt
type OptElapsedTime struct {
	ElapsedTime uint16
}

func (op *OptElapsedTime) Code() OptionCode {
	return OptionElapsedTime
}

// ToBytes marshals this option to bytes.
func (op *OptElapsedTime) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(op.ElapsedTime))
	return buf.Data()
}

func (op *OptElapsedTime) String() string {
	return fmt.Sprintf("OptElapsedTime{elapsedtime=%v}", op.ElapsedTime)
}

// build an OptElapsedTime structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptElapsedTime(data []byte) (*OptElapsedTime, error) {
	var opt OptElapsedTime
	buf := uio.NewBigEndianBuffer(data)
	opt.ElapsedTime = buf.Read16()
	return &opt, buf.FinError()
}
