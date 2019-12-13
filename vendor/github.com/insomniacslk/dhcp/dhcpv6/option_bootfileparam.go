package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptBootFileParam implements the OptionBootfileParam option
//
// This module defines the OPT_BOOTFILE_PARAM structure.
// https://www.ietf.org/rfc/rfc5970.txt (section 3.2)
type OptBootFileParam []string

var _ Option = OptBootFileParam(nil)

// Code returns the option code
func (op OptBootFileParam) Code() OptionCode {
	return OptionBootfileParam
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op OptBootFileParam) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, param := range op {
		if len(param) >= 1<<16 {
			// TODO: say something here instead of silently ignoring a parameter
			continue
		}
		buf.Write16(uint16(len(param)))
		buf.WriteBytes([]byte(param))
		/*if err := buf.Error(); err != nil {
			// TODO: description of `WriteBytes` says it could return
			// an error via `buf.Error()`. But a quick look into implementation of
			// `WriteBytes` at the moment of this comment showed it does not set any
			// errors to `Error()` output. It's required to make a decision:
			// to fix `WriteBytes` or it's description or
			// to find a way to handle an error here.
		}*/
	}
	return buf.Data()
}

func (op OptBootFileParam) String() string {
	return fmt.Sprintf("OptBootFileParam(%v)", ([]string)(op))
}

// ParseOptBootFileParam builds an OptBootFileParam structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptBootFileParam(data []byte) (result OptBootFileParam, err error) {
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(2) {
		length := buf.Read16()
		result = append(result, string(buf.CopyN(int(length))))
	}
	if err := buf.FinError(); err != nil {
		return nil, err
	}
	return
}
