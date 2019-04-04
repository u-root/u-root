package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/uio"
)

// OptDNSRecursiveNameServer represents a OptionDNSRecursiveNameServer option
//
// This module defines the OptDNSRecursiveNameServer structure.
// https://www.ietf.org/rfc/rfc3646.txt
type OptDNSRecursiveNameServer struct {
	NameServers []net.IP
}

// Code returns the option code
func (op *OptDNSRecursiveNameServer) Code() OptionCode {
	return OptionDNSRecursiveNameServer
}

// ToBytes returns the option serialized to bytes.
func (op *OptDNSRecursiveNameServer) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ns := range op.NameServers {
		buf.WriteBytes(ns.To16())
	}
	return buf.Data()
}

func (op *OptDNSRecursiveNameServer) String() string {
	return fmt.Sprintf("OptDNSRecursiveNameServer{nameservers=%v}", op.NameServers)
}

// ParseOptDNSRecursiveNameServer builds an OptDNSRecursiveNameServer structure
// from a sequence of bytes. The input data does not include option code and length
// bytes.
func ParseOptDNSRecursiveNameServer(data []byte) (*OptDNSRecursiveNameServer, error) {
	var opt OptDNSRecursiveNameServer
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(net.IPv6len) {
		opt.NameServers = append(opt.NameServers, buf.CopyN(net.IPv6len))
	}
	return &opt, buf.FinError()
}
