package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCodelUnspec = iota
	tcaCodelTarget
	tcaCodelLimit
	tcaCodelInterval
	tcaCodelECN
	tcaCodelCEThreshold
)

// Codel contains attributes of the codel discipline
type Codel struct {
	Target      *uint32
	Limit       *uint32
	Interval    *uint32
	ECN         *uint32
	CEThreshold *uint32
}

// unmarshalCodel parses the Codel-encoded data and stores the result in the value pointed to by info.
func unmarshalCodel(data []byte, info *Codel) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaCodelTarget:
			info.Target = uint32Ptr(ad.Uint32())
		case tcaCodelLimit:
			info.Limit = uint32Ptr(ad.Uint32())
		case tcaCodelInterval:
			info.Interval = uint32Ptr(ad.Uint32())
		case tcaCodelECN:
			info.ECN = uint32Ptr(ad.Uint32())
		case tcaCodelCEThreshold:
			info.CEThreshold = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalCodel()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalCodel returns the binary encoding of Red
func marshalCodel(info *Codel) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Codel: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Target != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelTarget, Data: uint32Value(info.Target)})
	}
	if info.Limit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelLimit, Data: uint32Value(info.Limit)})
	}
	if info.Interval != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelInterval, Data: uint32Value(info.Interval)})
	}
	if info.ECN != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelECN, Data: uint32Value(info.ECN)})
	}
	if info.CEThreshold != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelCEThreshold, Data: uint32Value(info.CEThreshold)})
	}
	return marshalAttributes(options)
}
