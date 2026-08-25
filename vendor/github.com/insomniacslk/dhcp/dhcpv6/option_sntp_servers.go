package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/uio/uio"
)

// OptSNTP returns a SNTP Servers option as defined by RFC 4075.
func OptSNTP(ip ...net.IP) Option {
	return &optSNTP{SNTPServers: ip}
}

type optSNTP struct {
	SNTPServers []net.IP
}

// Code returns the option code
func (op *optSNTP) Code() OptionCode {
	return OptionSNTPServerList
}

// ToBytes returns the option serialized to bytes.
func (op *optSNTP) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ns := range op.SNTPServers {
		buf.WriteBytes(ns.To16())
	}
	return buf.Data()
}

func (op *optSNTP) String() string {
	return fmt.Sprintf("%s: %v", op.Code(), op.SNTPServers)
}

// FromBytes builds an optSNTP structure from a sequence of bytes. The input
// data does not include option code and length bytes.
func (op *optSNTP) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(net.IPv6len) {
		op.SNTPServers = append(op.SNTPServers, buf.CopyN(net.IPv6len))
	}
	return buf.FinError()
}
