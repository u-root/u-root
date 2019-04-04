package dhcpv6

import (
	"fmt"

	"github.com/insomniacslk/dhcp/rfc1035label"
)

// OptDomainSearchList list implements a OptionDomainSearchList option
//
// This module defines the OptDomainSearchList structure.
// https://www.ietf.org/rfc/rfc3646.txt
type OptDomainSearchList struct {
	DomainSearchList *rfc1035label.Labels
}

func (op *OptDomainSearchList) Code() OptionCode {
	return OptionDomainSearchList
}

// ToBytes marshals this option to bytes.
func (op *OptDomainSearchList) ToBytes() []byte {
	return op.DomainSearchList.ToBytes()
}

func (op *OptDomainSearchList) String() string {
	return fmt.Sprintf("OptDomainSearchList{searchlist=%v}", op.DomainSearchList.Labels)
}

// ParseOptDomainSearchList builds an OptDomainSearchList structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func ParseOptDomainSearchList(data []byte) (*OptDomainSearchList, error) {
	var opt OptDomainSearchList
	var err error
	opt.DomainSearchList, err = rfc1035label.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
