package dhcpv6

import (
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

// OptRequestedOption implements the requested options option.
//
// This module defines the OptRequestedOption structure.
// https://www.ietf.org/rfc/rfc3315.txt
type OptRequestedOption struct {
	requestedOptions []OptionCode
}

func (op *OptRequestedOption) Code() OptionCode {
	return OptionORO
}

func (op *OptRequestedOption) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ro := range op.requestedOptions {
		buf.Write16(uint16(ro))
	}
	return buf.Data()
}

func (op *OptRequestedOption) RequestedOptions() []OptionCode {
	return op.requestedOptions
}

func (op *OptRequestedOption) SetRequestedOptions(opts []OptionCode) {
	op.requestedOptions = opts
}

func (op *OptRequestedOption) AddRequestedOption(opt OptionCode) {
	for _, requestedOption := range op.requestedOptions {
		if opt == requestedOption {
			fmt.Printf("Warning: option %s is already set, appending duplicate", opt)
		}
	}
	op.requestedOptions = append(op.requestedOptions, opt)
}

func (op *OptRequestedOption) String() string {
	names := make([]string, 0, len(op.requestedOptions))
	for _, code := range op.requestedOptions {
		names = append(names, code.String())
	}
	return fmt.Sprintf("OptRequestedOption{options=[%v]}", strings.Join(names, ", "))
}

// build an OptRequestedOption structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRequestedOption(data []byte) (*OptRequestedOption, error) {
	var opt OptRequestedOption
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(2) {
		opt.requestedOptions = append(opt.requestedOptions, OptionCode(buf.Read16()))
	}
	return &opt, buf.FinError()
}
