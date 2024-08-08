package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaHhfUnspec = iota
	tcaHhfBacklogLimit
	tcaHhfQuantum
	tcaHhfHHFlowsLimit
	tcaHhfResetTimeout
	tcaHhfAdmitBytes
	tcaHhfEVICTTimeout
	tcaHhfNonHHWeight
)

// Hhf contains attributes of the hhf discipline
type Hhf struct {
	BacklogLimit *uint32
	Quantum      *uint32
	HHFlowsLimit *uint32
	ResetTimeout *uint32
	AdmitBytes   *uint32
	EVICTTimeout *uint32
	NonHHWeight  *uint32
}

// unmarshalHhf parses the Hhf-encoded data and stores the result in the value pointed to by info.
func unmarshalHhf(data []byte, info *Hhf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaHhfBacklogLimit:
			info.BacklogLimit = uint32Ptr(ad.Uint32())
		case tcaHhfQuantum:
			info.Quantum = uint32Ptr(ad.Uint32())
		case tcaHhfHHFlowsLimit:
			info.HHFlowsLimit = uint32Ptr(ad.Uint32())
		case tcaHhfResetTimeout:
			info.ResetTimeout = uint32Ptr(ad.Uint32())
		case tcaHhfAdmitBytes:
			info.AdmitBytes = uint32Ptr(ad.Uint32())
		case tcaHhfEVICTTimeout:
			info.EVICTTimeout = uint32Ptr(ad.Uint32())
		case tcaHhfNonHHWeight:
			info.NonHHWeight = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalHhf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalHhf returns the binary encoding of Hhf
func marshalHhf(info *Hhf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Hhf: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.BacklogLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfBacklogLimit, Data: uint32Value(info.BacklogLimit)})
	}
	if info.Quantum != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfQuantum, Data: uint32Value(info.Quantum)})
	}
	if info.HHFlowsLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfHHFlowsLimit, Data: uint32Value(info.HHFlowsLimit)})
	}
	if info.ResetTimeout != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfResetTimeout, Data: uint32Value(info.ResetTimeout)})
	}
	if info.AdmitBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfAdmitBytes, Data: uint32Value(info.AdmitBytes)})
	}
	if info.EVICTTimeout != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfEVICTTimeout, Data: uint32Value(info.EVICTTimeout)})
	}
	if info.NonHHWeight != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfNonHHWeight, Data: uint32Value(info.NonHHWeight)})
	}

	return marshalAttributes(options)
}
